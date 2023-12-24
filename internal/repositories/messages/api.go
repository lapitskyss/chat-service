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

func (r *Repo) GetMessageByID(ctx context.Context, id types.MessageID) (*Message, error) {
	m, err := r.db.Message(ctx).
		Query().
		Unique(false).
		Where(message.ID(id)).
		Only(ctx)
	if err != nil {
		if store.IsNotFound(err) {
			return nil, fmt.Errorf("message id %v: %w", id, ErrMsgNotFound)
		}
		return nil, fmt.Errorf("get message by id: %v", err)
	}
	msg := adaptMessage(m)
	return &msg, nil
}

func (r *Repo) GetMessageByIDWithManager(ctx context.Context, id types.MessageID) (*MessageWithManager, error) {
	m, err := r.db.Message(ctx).
		Query().
		Unique(false).
		WithProblem().
		Where(message.ID(id)).
		Only(ctx)
	if err != nil {
		if store.IsNotFound(err) {
			return nil, fmt.Errorf("message id %v: %w", id, ErrMsgNotFound)
		}
		return nil, fmt.Errorf("get message by id: %v", err)
	}
	msg := adaptMessageWithManager(m)
	return &msg, nil
}

func (r *Repo) GetMessageByRequestID(ctx context.Context, reqID types.RequestID) (*Message, error) {
	m, err := r.db.Message(ctx).
		Query().
		Unique(false).
		Where(message.InitialRequestID(reqID)).
		Only(ctx)
	if err != nil {
		if store.IsNotFound(err) {
			return nil, fmt.Errorf("request id %v: %w", reqID, ErrMsgNotFound)
		}
		return nil, fmt.Errorf("get message by request id: %v", err)
	}
	msg := adaptMessage(m)
	return &msg, nil
}

func (r *Repo) GetServiceMessageByID(ctx context.Context, id types.MessageID) (*ServiceMessage, error) {
	m, err := r.db.Message(ctx).
		Query().
		WithChat().
		WithProblem().
		Unique(false).
		Where(message.ID(id)).
		Only(ctx)
	if err != nil {
		if store.IsNotFound(err) {
			return nil, fmt.Errorf("message id %v: %w", id, ErrMsgNotFound)
		}
		return nil, fmt.Errorf("get message by id: %v", err)
	}
	msg := adaptServiceMessage(m)
	return &msg, nil
}

func (r *Repo) CreateClientVisible(
	ctx context.Context,
	reqID types.RequestID,
	problemID types.ProblemID,
	chatID types.ChatID,
	authorID types.UserID,
	msgBody string,
) (*Message, error) {
	m, err := r.db.Message(ctx).Create().
		SetChatID(chatID).
		SetProblemID(problemID).
		SetAuthorID(authorID).
		SetIsVisibleForClient(true).
		SetIsVisibleForManager(false).
		SetBody(msgBody).
		SetInitialRequestID(reqID).
		Save(ctx)
	if err != nil {
		return nil, fmt.Errorf("create msg: %v", err)
	}

	mm := adaptMessage(m)
	return &mm, nil
}

func (r *Repo) CreateFullVisible(
	ctx context.Context,
	reqID types.RequestID,
	problemID types.ProblemID,
	chatID types.ChatID,
	authorID types.UserID,
	msgBody string,
) (*Message, error) {
	m, err := r.db.Message(ctx).Create().
		SetChatID(chatID).
		SetProblemID(problemID).
		SetAuthorID(authorID).
		SetIsVisibleForClient(true).
		SetIsVisibleForManager(true).
		SetBody(msgBody).
		SetInitialRequestID(reqID).
		Save(ctx)
	if err != nil {
		return nil, fmt.Errorf("create full visible msg: %v", err)
	}

	mm := adaptMessage(m)
	return &mm, nil
}

func (r *Repo) CreateServiceMsg(
	ctx context.Context,
	reqID types.RequestID,
	problemID types.ProblemID,
	chatID types.ChatID,
	msgBody string,
	visibleForClient bool,
	visibleForManager bool,
) (*Message, error) {
	m, err := r.db.Message(ctx).Create().
		SetChatID(chatID).
		SetProblemID(problemID).
		SetIsVisibleForClient(visibleForClient).
		SetIsVisibleForManager(visibleForManager).
		SetIsService(true).
		SetBody(msgBody).
		SetInitialRequestID(reqID).
		Save(ctx)
	if err != nil {
		return nil, fmt.Errorf("create msg: %v", err)
	}

	mm := adaptMessage(m)
	return &mm, nil
}
