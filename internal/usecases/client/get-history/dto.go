package gethistory

import (
	"errors"
	"time"

	messagesrepo "github.com/lapitskyss/chat-service/internal/repositories/messages"
	"github.com/lapitskyss/chat-service/internal/types"
	"github.com/lapitskyss/chat-service/internal/validator"
)

type Request struct {
	ID       types.RequestID `validate:"required"`
	ClientID types.UserID    `validate:"required"`
	PageSize int             `validate:"omitempty,gte=10,lte=100"`
	Cursor   string          `validate:"omitempty,base64url"`
}

type Response struct {
	Messages   []Message
	NextCursor string
}

type Message struct {
	ID         types.MessageID
	AuthorID   types.UserID
	Body       string
	IsBlocked  bool
	IsService  bool
	IsReceived bool
	CreatedAt  time.Time
}

func mapMassages(messages []messagesrepo.Message) []Message {
	result := make([]Message, len(messages))
	for i, m := range messages {
		result[i] = Message{
			ID:         m.ID,
			AuthorID:   m.AuthorID,
			Body:       m.Body,
			IsBlocked:  m.IsBlocked,
			IsService:  m.IsService,
			IsReceived: m.IsVisibleForManager && !m.IsBlocked,
			CreatedAt:  m.CreatedAt,
		}
	}
	return result
}

func (r Request) Validate() error {
	if r.Cursor == "" && r.PageSize == 0 {
		return errors.New("either cursor or page size must be specified")
	}
	if r.Cursor != "" && r.PageSize != 0 {
		return errors.New("either cursor or page size must be specified, not both")
	}
	return validator.Validator.Struct(r)
}
