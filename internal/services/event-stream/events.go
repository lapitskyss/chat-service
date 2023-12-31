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
	UserID      types.UserID    `validate:"-"`
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

// NewChatEvent is a signal about the new chat is received for manager.
type NewChatEvent struct {
	event
	EventID             types.EventID   `validate:"required"`
	ChatID              types.ChatID    `validate:"required"`
	ClientID            types.UserID    `validate:"required"`
	RequestID           types.RequestID `validate:"required"`
	CanTakeMoreProblems bool            `validate:"required"`
}

func NewNewChatEvent(
	eventID types.EventID,
	chatID types.ChatID,
	clientID types.UserID,
	requestID types.RequestID,
	canTakeMoreProblems bool,
) *NewChatEvent {
	return &NewChatEvent{
		EventID:             eventID,
		ChatID:              chatID,
		ClientID:            clientID,
		RequestID:           requestID,
		CanTakeMoreProblems: canTakeMoreProblems,
	}
}

func (e NewChatEvent) Validate() error {
	return validator.Validator.Struct(e)
}

// ChatClosedEvent is a signal about chat was closed by manager.
type ChatClosedEvent struct {
	event
	EventID             types.EventID   `validate:"required"`
	ChatID              types.ChatID    `validate:"required"`
	RequestID           types.RequestID `validate:"required"`
	CanTakeMoreProblems bool            `validate:"required"`
}

func NewChatClosedEvent(
	eventID types.EventID,
	chatID types.ChatID,
	requestID types.RequestID,
	canTakeMoreProblems bool,
) *ChatClosedEvent {
	return &ChatClosedEvent{
		EventID:             eventID,
		ChatID:              chatID,
		RequestID:           requestID,
		CanTakeMoreProblems: canTakeMoreProblems,
	}
}

func (e ChatClosedEvent) Validate() error {
	return validator.Validator.Struct(e)
}

type TypingEvent struct {
	event
	EventID   types.EventID   `validate:"required"`
	ClientID  types.UserID    `validate:"required"`
	RequestID types.RequestID `validate:"required"`
}

func NewTypingEvent(
	eventID types.EventID,
	clientID types.UserID,
	requestID types.RequestID,
) *TypingEvent {
	return &TypingEvent{
		EventID:   eventID,
		ClientID:  clientID,
		RequestID: requestID,
	}
}

func (e TypingEvent) Validate() error {
	return validator.Validator.Struct(e)
}
