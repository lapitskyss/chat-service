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
	RequestID           types.RequestID
	IsVisibleForClient  bool
	IsVisibleForManager bool
	Body                string
	IsBlocked           bool
	IsService           bool
	CreatedAt           time.Time
}

type ServiceMessage struct {
	ID                  types.MessageID
	ChatID              types.ChatID
	ClientID            types.UserID
	ManagerID           types.UserID
	RequestID           types.RequestID
	IsVisibleForClient  bool
	IsVisibleForManager bool
	Body                string
	IsBlocked           bool
	CreatedAt           time.Time
}

func adaptMessage(m *store.Message) Message {
	return Message{
		ID:                  m.ID,
		ChatID:              m.ChatID,
		AuthorID:            m.AuthorID,
		RequestID:           m.InitialRequestID,
		IsVisibleForClient:  m.IsVisibleForClient,
		IsVisibleForManager: m.IsVisibleForManager,
		Body:                m.Body,
		IsBlocked:           m.IsBlocked,
		IsService:           m.IsService,
		CreatedAt:           m.CreatedAt,
	}
}

func adaptMessages(mm []*store.Message) []Message {
	result := make([]Message, len(mm))
	for i, m := range mm {
		result[i] = adaptMessage(m)
	}
	return result
}

func adaptServiceMessage(m *store.Message) ServiceMessage {
	return ServiceMessage{
		ID:                  m.ID,
		ChatID:              m.ChatID,
		ClientID:            m.Edges.Chat.ClientID,
		ManagerID:           m.Edges.Problem.ManagerID,
		RequestID:           m.InitialRequestID,
		IsVisibleForClient:  m.IsVisibleForClient,
		IsVisibleForManager: m.IsVisibleForManager,
		Body:                m.Body,
		IsBlocked:           m.IsBlocked,
		CreatedAt:           m.CreatedAt,
	}
}
