package sendmessage

import (
	"time"

	"github.com/lapitskyss/chat-service/internal/types"
	"github.com/lapitskyss/chat-service/internal/validator"
)

type Request struct {
	ID          types.RequestID `validate:"required"`
	ClientID    types.UserID    `validate:"required"`
	MessageBody string          `validate:"required,min=1,max=3000"`
}

func (r Request) Validate() error {
	return validator.Validator.Struct(r)
}

type Response struct {
	AuthorID  types.UserID
	MessageID types.MessageID
	CreatedAt time.Time
}
