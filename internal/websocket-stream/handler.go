package websocketstream

import (
	"context"
	"errors"
	"fmt"
	"io"
	"time"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"

	"github.com/lapitskyss/chat-service/internal/middlewares"
	eventstream "github.com/lapitskyss/chat-service/internal/services/event-stream"
	"github.com/lapitskyss/chat-service/internal/types"
)

const (
	writeTimeout = time.Second
)

type eventStream interface {
	Subscribe(ctx context.Context, userID types.UserID) (<-chan eventstream.Event, error)
}

type msgReadHandler interface {
	Handle(ctx context.Context, userID types.UserID, r io.Reader) error
}

//go:generate options-gen -out-filename=handler_options.gen.go -from-struct=Options
type Options struct {
	pingPeriod time.Duration `default:"3s" validate:"omitempty,min=100ms,max=30s"`

	logger       *zap.Logger     `option:"mandatory" validate:"required"`
	eventStream  eventStream     `option:"mandatory" validate:"required"`
	eventAdapter EventAdapter    `option:"mandatory" validate:"required"`
	eventWriter  EventWriter     `option:"mandatory" validate:"required"`
	upgrader     Upgrader        `option:"mandatory" validate:"required"`
	readHandler  msgReadHandler  `option:"mandatory" validate:"required"`
	shutdownCh   <-chan struct{} `option:"mandatory" validate:"required"`
}

type HTTPHandler struct {
	Options
}

func NewHTTPHandler(opts Options) (*HTTPHandler, error) {
	if err := opts.Validate(); err != nil {
		return nil, fmt.Errorf("validate options: %v", err)
	}
	return &HTTPHandler{Options: opts}, nil
}

func (h *HTTPHandler) Serve(c echo.Context) error {
	conn, err := h.upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		return fmt.Errorf("upgrate connection: %v", err)
	}

	closer := newWsCloser(h.logger, conn)
	defer closer.Close(websocket.CloseNormalClosure)

	ctx := c.Request().Context()
	userID := middlewares.MustUserID(c)

	events, err := h.eventStream.Subscribe(ctx, userID)
	if err != nil {
		h.logger.Error("cannot subscribe for events", zap.Error(err))
		closer.Close(websocket.CloseInternalServerErr)
		return nil
	}

	errGrp, ctx := errgroup.WithContext(ctx)
	errGrp.Go(func() error {
		return h.readLoop(ctx, conn, userID)
	})
	errGrp.Go(func() error {
		return h.writeLoop(ctx, conn, events)
	})
	errGrp.Go(func() error {
		select {
		case <-ctx.Done():
		case <-h.shutdownCh:
			closer.Close(websocket.CloseNormalClosure)
		}
		return nil
	})
	err = errGrp.Wait()
	if err != nil {
		h.logger.Error("unexpected error", zap.Error(err))
		closer.Close(websocket.CloseInternalServerErr)
		return nil
	}

	return nil
}

// readLoop listen PONGs.
func (h *HTTPHandler) readLoop(ctx context.Context, ws Websocket, userID types.UserID) error {
	pongDeadline := 2 * h.pingPeriod

	_ = ws.SetReadDeadline(time.Now().Add(pongDeadline))
	ws.SetPongHandler(func(string) error {
		_ = ws.SetReadDeadline(time.Now().Add(pongDeadline))

		h.logger.Debug("pong")

		return nil
	})

	for {
		msgType, reader, err := ws.NextReader()
		if err != nil {
			if websocket.IsCloseError(err, websocket.CloseNormalClosure,
				websocket.CloseGoingAway, websocket.CloseNoStatusReceived) {
				return nil
			}
			return fmt.Errorf("read next message: %v", err)
		}
		if msgType == websocket.TextMessage {
			err := h.readHandler.Handle(ctx, userID, reader)
			if err != nil {
				return fmt.Errorf("handle read event: %v", err)
			}
		}
	}
}

// writeLoop listen events and writes them into Websocket.
func (h *HTTPHandler) writeLoop(ctx context.Context, ws Websocket, events <-chan eventstream.Event) error {
	t := time.NewTicker(h.pingPeriod)
	defer t.Stop()

	for {
		select {
		case <-ctx.Done():
			return nil
		case <-t.C:
			err := h.writePing(ws)
			if err != nil {
				return fmt.Errorf("write ping: %v", err)
			}
		case event, ok := <-events:
			if !ok {
				return errors.New("events stream was closed")
			}
			err := h.writeEvent(ws, event)
			if err != nil {
				return fmt.Errorf("write event: %v", err)
			}
		}
	}
}

func (h *HTTPHandler) writePing(ws Websocket) error {
	err := ws.SetWriteDeadline(time.Now().Add(writeTimeout))
	if err != nil {
		return fmt.Errorf("set write deadline: %v", err)
	}
	err = ws.WriteMessage(websocket.PingMessage, nil)
	if err != nil {
		if errors.Is(err, websocket.ErrCloseSent) {
			return nil
		}
		return fmt.Errorf("send ping msg: %v", err)
	}

	h.logger.Debug("ping")

	return nil
}

func (h *HTTPHandler) writeEvent(ws Websocket, event eventstream.Event) error {
	err := ws.SetWriteDeadline(time.Now().Add(writeTimeout))
	if err != nil {
		return fmt.Errorf("set write deadline: %v", err)
	}

	w, err := ws.NextWriter(websocket.TextMessage)
	if err != nil {
		if errors.Is(err, websocket.ErrCloseSent) {
			return nil
		}
		return fmt.Errorf("get next writer: %v", err)
	}

	result, err := h.eventAdapter.Adapt(event)
	if err != nil {
		return fmt.Errorf("adapt event: %v", err)
	}

	err = h.eventWriter.Write(result, w)
	if err != nil {
		if errors.Is(err, websocket.ErrCloseSent) {
			return nil
		}
		return fmt.Errorf("write event: %v", err)
	}

	err = w.Close()
	if err != nil {
		return fmt.Errorf("close writer: %v", err)
	}

	return nil
}
