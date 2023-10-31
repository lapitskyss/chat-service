package messagesrepo

import (
	"time"

	"github.com/lapitskyss/chat-service/internal/store"
	"github.com/lapitskyss/chat-service/internal/types"
)

type Message struct {
	ID                  types.MessageID
	ChatID              types.ChatID
	AuthorID            types.UserID
	IsVisibleForClient  bool
	IsVisibleForManager bool
	Body                string
	IsBlocked           bool
	IsService           bool
	CreatedAt           time.Time
}

func adaptStoreMessage(m *store.Message) Message {
	return Message{
		ID:                  m.ID,
		ChatID:              m.ChatID,
		AuthorID:            m.AuthorID,
		IsVisibleForClient:  m.IsVisibleForClient,
		IsVisibleForManager: m.IsVisibleForManager,
		Body:                m.Body,
		IsBlocked:           m.IsBlocked,
		IsService:           m.IsService,
		CreatedAt:           m.CreatedAt,
	}
}

func adaptStoreMessages(mm []*store.Message) []Message {
	result := make([]Message, len(mm))
	for i, m := range mm {
		result[i] = adaptStoreMessage(m)
	}
	return result
}
