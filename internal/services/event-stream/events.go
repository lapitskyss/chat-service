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
	EventID   types.EventID   `validate:"required"`
	RequestID types.RequestID `validate:"required"`
	MessageID types.MessageID `validate:"required"`
}

func NewMessageSentEvent(
	eventID types.EventID,
	requestID types.RequestID,
	messageID types.MessageID,
) *MessageSentEvent {
	return &MessageSentEvent{
		EventID:   eventID,
		RequestID: requestID,
		MessageID: messageID,
	}
}

func (e MessageSentEvent) Validate() error {
	return validator.Validator.Struct(e)
}

// MessageBlockEvent indicates that the message was blocked by AFC.
type MessageBlockEvent struct {
	event
	EventID   types.EventID   `validate:"required"`
	RequestID types.RequestID `validate:"required"`
	MessageID types.MessageID `validate:"required"`
}

func NewMessageBlockEvent(
	eventID types.EventID,
	requestID types.RequestID,
	messageID types.MessageID,
) *MessageBlockEvent {
	return &MessageBlockEvent{
		EventID:   eventID,
		RequestID: requestID,
		MessageID: messageID,
	}
}

func (e MessageBlockEvent) Validate() error {
	return validator.Validator.Struct(e)
}

// NewMessageEvent is a signal about the appearance of a new message in the chat.
type NewMessageEvent struct {
	event
	EventID     types.EventID   `validate:"required"`
	RequestID   types.RequestID `validate:"required"`
	ChatID      types.ChatID    `validate:"required"`
	MessageID   types.MessageID `validate:"required"`
	UserID      types.UserID    `validate:"required"`
	CreatedAt   time.Time       `validate:"-"`
	MessageBody string          `validate:"required"`
	IsService   bool            `validate:"-"`
}

func NewNewMessageEvent(
	eventID types.EventID,
	requestID types.RequestID,
	chatID types.ChatID,
	messageID types.MessageID,
	userID types.UserID,
	createdAt time.Time,
	body string,
	isService bool,
) *NewMessageEvent {
	return &NewMessageEvent{
		EventID:     eventID,
		RequestID:   requestID,
		ChatID:      chatID,
		MessageID:   messageID,
		UserID:      userID,
		CreatedAt:   createdAt,
		MessageBody: body,
		IsService:   isService,
	}
}

func (e NewMessageEvent) Validate() error {
	return validator.Validator.Struct(e)
}
