package managerworkspace

import (
	"time"

	"github.com/lapitskyss/chat-service/internal/types"
)

type Message struct {
	ID        types.MessageID
	ChatID    types.ChatID
	AuthorID  types.UserID
	Body      string
	CreatedAt time.Time
}

func NewMessage(
	id types.MessageID,
	chatID types.ChatID,
	authorID types.UserID,
	body string,
	createdAt time.Time,
) *Message {
	return &Message{
		ID:        id,
		ChatID:    chatID,
		AuthorID:  authorID,
		Body:      body,
		CreatedAt: createdAt,
	}
}
