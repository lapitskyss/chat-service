package wsstream

import (
	"context"
	"fmt"
	"net/http"
	"time"

	gorillaws "github.com/gorilla/websocket"
	"go.uber.org/multierr"
)

//go:generate options-gen -out-filename=stream_options.gen.go -from-struct=Options
type Options struct {
	endpoint      string                                       `option:"mandatory" validate:"required,url"`
	origin        string                                       `option:"mandatory" validate:"required"`
	secWsProtocol string                                       `option:"mandatory" validate:"required"`
	authToken     string                                       `option:"mandatory" validate:"required"`
	eventHandler  func(ctx context.Context, data []byte) error `option:"mandatory" validate:"required"`
}

type Stream struct {
	Options
	dialer       *gorillaws.Dialer
	eventSignals chan struct{}
}

func New(opts Options) (*Stream, error) {
	if err := opts.Validate(); err != nil {
		return nil, fmt.Errorf("validate options: %v", err)
	}
	return &Stream{
		Options:      opts,
		dialer:       gorillaws.DefaultDialer,
		eventSignals: make(chan struct{}, 1000),
	}, nil
}

func (s *Stream) Run(ctx context.Context) (errReturned error) {
	conn, resp, err := s.dialer.DialContext(ctx, s.endpoint, http.Header{
		"Origin":                 []string{s.origin},
		"Sec-WebSocket-Protocol": []string{s.secWsProtocol + ", " + s.authToken},
	})
	if err != nil {
		return fmt.Errorf("dial: %v", err)
	}
	_ = resp.Body.Close()

	defer multierr.AppendInvoke(&errReturned, multierr.Close(conn))
	conn.SetPingHandler(nil) // Default handler.

	go func() {
		<-ctx.Done()
		_ = conn.WriteControl(
			gorillaws.CloseMessage,
			gorillaws.FormatCloseMessage(gorillaws.CloseNormalClosure, ""),
			time.Now().Add(time.Second),
		)
	}()

	for {
		_, event, err := conn.ReadMessage()
		if gorillaws.IsCloseError(err, gorillaws.CloseNormalClosure) {
			return nil
		}

		if err := s.eventHandler(ctx, event); err != nil {
			return fmt.Errorf("handle new event: %v", err)
		}

		select {
		case <-ctx.Done():
			return nil
		case s.eventSignals <- struct{}{}:
		}
	}
}

func (s *Stream) EventSignals() <-chan struct{} {
	return s.eventSignals
}
