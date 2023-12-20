package clienthandler

import (
	"context"
	"fmt"
	"io"

	"github.com/lapitskyss/chat-service/internal/types"
	clienttypingmessage "github.com/lapitskyss/chat-service/internal/usecases/client/typing-message"
)

//go:generate mockgen -source=$GOFILE -destination=mocks/handlers_mocks.gen.go -package=clienthandlermocks

type TypingMessageUseCase interface {
	Handle(ctx context.Context, req clienttypingmessage.Request) error
}

//go:generate options-gen -out-filename=handler_options.gen.go -from-struct=Options
type Options struct {
	eventReader   EventReader          `option:"mandatory" validate:"required"`
	typingMessage TypingMessageUseCase `option:"mandatory" validate:"required"`
}

type Handler struct {
	Options
}

func New(opts Options) (Handler, error) {
	if err := opts.Validate(); err != nil {
		return Handler{}, fmt.Errorf("validate options: %v", err)
	}
	return Handler{Options: opts}, nil
}

func (h Handler) Handle(ctx context.Context, userID types.UserID, r io.Reader) error {
	e, err := h.eventReader.Read(r)
	if err != nil {
		return fmt.Errorf("read event, %v", err)
	}

	if err = e.Validate(); err != nil {
		return fmt.Errorf("validate event, %v", err)
	}

	switch t := e.(type) {
	case *ClientTypingEvent:
		err = h.typingMessage.Handle(ctx, clienttypingmessage.Request{
			ID:       t.RequestID,
			ClientID: userID,
		})
		if err != nil {
			return fmt.Errorf("handle client typing message use case, %v", err)
		}
		return nil
	default:
		return fmt.Errorf("unexpecetd event type")
	}
}
