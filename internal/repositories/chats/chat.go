package chatsrepo

import (
	"github.com/lapitskyss/chat-service/internal/store"
	"github.com/lapitskyss/chat-service/internal/types"
)

type Chat struct {
	ID       types.ChatID
	ClientID types.UserID
}

func adaptChat(c *store.Chat) Chat {
	return Chat{
		ID:       c.ID,
		ClientID: c.ClientID,
	}
}

func adaptChats(cc []*store.Chat) []Chat {
	result := make([]Chat, len(cc))
	for i, c := range cc {
		result[i] = adaptChat(c)
	}
	return result
}
