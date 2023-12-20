package managerhandler

import (
	"github.com/lapitskyss/chat-service/internal/types"
	"github.com/lapitskyss/chat-service/internal/validator"
)

type Event interface {
	eventMarker()
	Validate() error
}

type event struct{}         //
func (*event) eventMarker() {}

type ManagerTypingEvent struct {
	event
	ChatID    types.ChatID    `validate:"required"`
	RequestID types.RequestID `validate:"required"`
}

func NewManagerTypingEvent(chatID types.ChatID, requestID types.RequestID) *ManagerTypingEvent {
	return &ManagerTypingEvent{
		ChatID:    chatID,
		RequestID: requestID,
	}
}

func (e ManagerTypingEvent) Validate() error {
	return validator.Validator.Struct(e)
}
