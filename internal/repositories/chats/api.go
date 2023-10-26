package chatsrepo

import (
	"context"
	"fmt"

	"github.com/lapitskyss/chat-service/internal/store/chat"
	"github.com/lapitskyss/chat-service/internal/types"
)

func (r *Repo) CreateIfNotExists(ctx context.Context, userID types.UserID) (types.ChatID, error) {
	chatID, err := r.db.Chat(ctx).
		Create().
		SetID(types.NewChatID()).
		SetClientID(userID).
		OnConflictColumns(chat.FieldClientID).
		UpdateNewValues().
		ID(ctx)
	if err != nil {
		return types.ChatID{}, fmt.Errorf("create chat: %v", err)
	}
	return chatID, nil
}
