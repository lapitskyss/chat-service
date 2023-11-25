package eventstream

import (
	"time"

	"github.com/lapitskyss/chat-service/internal/types"
	"github.com/lapitskyss/chat-service/internal/validator"
)

type Event interface {
	eventMarker()
	Validate() error
}

type event struct{}         //
func (*event) eventMarker() {}

// MessageSentEvent indicates that the message was checked by AFC
// and was sent to the manager. Two gray ticks.
type MessageSentEvent struct {
	event
}

func (e MessageSentEvent) Validate() error { panic("not implemented") }

// NewMessageEvent is a signal about the appearance of a new message in the chat.
type NewMessageEvent struct {
	event
	EventID     types.EventID   `validate:"required"`
	RequestID   types.RequestID `validate:"required"`
	ChatID      types.ChatID    `validate:"required"`
	MessageID   types.MessageID `validate:"required"`
	UserID      types.UserID    `validate:"required"`
	Time        time.Time       `validate:"-"`
	MessageBody string          `validate:"required"`
	SomeFiled   bool            `validate:"-"`
}

func NewNewMessageEvent(
	eventID types.EventID,
	requestID types.RequestID,
	chatID types.ChatID,
	messageID types.MessageID,
	userID types.UserID,
	time time.Time,
	body string,
	someFiled bool,
) *NewMessageEvent {
	return &NewMessageEvent{
		EventID:     eventID,
		RequestID:   requestID,
		ChatID:      chatID,
		MessageID:   messageID,
		UserID:      userID,
		Time:        time,
		MessageBody: body,
		SomeFiled:   someFiled,
	}
}

func (e NewMessageEvent) Validate() error {
	return validator.Validator.Struct(e)
}
