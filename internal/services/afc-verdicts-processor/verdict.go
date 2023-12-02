package afcverdictsprocessor

import (
	"github.com/lapitskyss/chat-service/internal/types"
	"github.com/lapitskyss/chat-service/internal/validator"
)

type VerdictStatus string

var (
	VerdictStatusOK         VerdictStatus = "ok"
	VerdictStatusSuspicious VerdictStatus = "suspicious"
)

type Verdict struct {
	ChatID    types.ChatID    `json:"chatId" validate:"required"`
	MessageID types.MessageID `json:"messageId" validate:"required"`
	Status    VerdictStatus   `json:"status" validate:"required,oneof=ok suspicious"`
}

func (v *Verdict) Valid() error {
	return validator.Validator.Struct(v)
}
