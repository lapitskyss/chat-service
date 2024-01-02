package managertypingmessage

import (
	"github.com/lapitskyss/chat-service/internal/types"
	"github.com/lapitskyss/chat-service/internal/validator"
)

type Request struct {
	ID        types.RequestID `validate:"required"`
	ManagerID types.UserID    `validate:"required"`
	ChatID    types.ChatID    `validate:"required"`
}

func (r Request) Validate() error {
	return validator.Validator.Struct(r)
}
