package chatsrepo

import (
	"context"
	"fmt"

	"github.com/lapitskyss/chat-service/internal/store/chat"
	"github.com/lapitskyss/chat-service/internal/store/problem"
	"github.com/lapitskyss/chat-service/internal/types"
)

func (r *Repo) CreateIfNotExists(ctx context.Context, userID types.UserID) (types.ChatID, error) {
	chatID, err := r.db.Chat(ctx).
		Create().
		SetClientID(userID).
		OnConflictColumns(chat.FieldClientID).
		Ignore().
		ID(ctx)
	if err != nil {
		return types.ChatID{}, fmt.Errorf("create chat: %v", err)
	}
	return chatID, nil
}

func (r *Repo) AllWithOpenProblemsForManager(ctx context.Context, managerID types.UserID) ([]Chat, error) {
	chats, err := r.db.Chat(ctx).
		Query().
		Where(
			chat.HasProblemsWith(
				problem.ManagerID(managerID),
				problem.ResolvedAtIsNil(),
			),
		).
		All(ctx)
	if err != nil {
		return nil, fmt.Errorf("get chat with open problems: %v", err)
	}
	return adaptChats(chats), nil
}
