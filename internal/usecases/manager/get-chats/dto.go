package getchats

import (
	chatsrepo "github.com/lapitskyss/chat-service/internal/repositories/chats"
	"github.com/lapitskyss/chat-service/internal/types"
	"github.com/lapitskyss/chat-service/internal/validator"
)

type Request struct {
	ID        types.RequestID `validate:"required"`
	ManagerID types.UserID    `validate:"required"`
}

func (r Request) Validate() error {
	return validator.Validator.Struct(r)
}

type Response struct {
	Chats []Chat
}

type Chat struct {
	ID       types.ChatID
	ClientID types.UserID
}

func adaptChats(cc []chatsrepo.Chat) []Chat {
	result := make([]Chat, len(cc))
	for i, c := range cc {
		result[i] = Chat{
			ID:       c.ID,
			ClientID: c.ClientID,
		}
	}
	return result
}
