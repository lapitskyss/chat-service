package clientchat

import (
	"time"

	"github.com/lapitskyss/chat-service/internal/types"
)

type Message struct {
	ID         types.MessageID
	AuthorID   types.UserID
	Body       string
	IsService  bool
	IsBlocked  bool
	IsReceived bool
	CreatedAt  time.Time
}

func NewMessage(
	id types.MessageID,
	authorID *types.UserID,
	body string,
	isService, isBlocked, isReceived bool,
	createdAt time.Time,
) *Message {
	msg := Message{
		ID:         id,
		Body:       body,
		IsService:  isService,
		IsBlocked:  isBlocked,
		IsReceived: isReceived,
		CreatedAt:  createdAt,
	}
	if uid := authorID; uid != nil {
		msg.AuthorID = *uid
	}
	return &msg
}
