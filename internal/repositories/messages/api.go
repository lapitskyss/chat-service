package messagesrepo

import (
	"context"
	"errors"
	"fmt"

	"github.com/lapitskyss/chat-service/internal/store"
	"github.com/lapitskyss/chat-service/internal/store/message"
	"github.com/lapitskyss/chat-service/internal/types"
)

var ErrMsgNotFound = errors.New("message not found")

func (r *Repo) GetMessageByRequestID(ctx context.Context, reqID types.RequestID) (*Message, error) {
	m, err := r.db.Message(ctx).
		Query().
		Where(message.InitialRequestID(reqID)).
		Only(ctx)
	if err != nil {
		if store.IsNotFound(err) {
			return nil, ErrMsgNotFound
		}
		return nil, fmt.Errorf("get message by request id: %v", err)
	}
	msg := adaptStoreMessage(m)
	return &msg, nil
}

// CreateClientVisible creates a message that is visible only to the client.
func (r *Repo) CreateClientVisible(
	ctx context.Context,
	reqID types.RequestID,
	problemID types.ProblemID,
	chatID types.ChatID,
	authorID types.UserID,
	msgBody string,
) (*Message, error) {
	msg := newClientVisibleMessage(problemID, chatID, authorID, msgBody)

	err := r.db.Message(ctx).
		Create().
		SetID(msg.ID).
		SetChatID(msg.ChatID).
		SetProblemID(msg.ProblemID).
		SetAuthorID(msg.AuthorID).
		SetInitialRequestID(reqID).
		SetIsVisibleForClient(msg.IsVisibleForClient).
		SetIsVisibleForManager(msg.IsVisibleForManager).
		SetBody(msg.Body).
		SetCheckedAt(msg.CheckedAt).
		SetIsBlocked(msg.IsBlocked).
		SetIsService(msg.IsService).
		SetCreatedAt(msg.CreatedAt).
		Exec(ctx)
	if err != nil {
		return nil, fmt.Errorf("create message client visible: %v", err)
	}
	return msg, nil
}
