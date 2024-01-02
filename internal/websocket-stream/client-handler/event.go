package clienthandler

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

type ClientTypingEvent struct {
	event
	RequestID types.RequestID `validate:"required"`
}

func NewClientTypingEvent(requestID types.RequestID) *ClientTypingEvent {
	return &ClientTypingEvent{
		RequestID: requestID,
	}
}

func (e ClientTypingEvent) Validate() error {
	return validator.Validator.Struct(e)
}
