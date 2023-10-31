package problemsrepo

import (
	"context"
	"fmt"

	"github.com/lapitskyss/chat-service/internal/store"
	"github.com/lapitskyss/chat-service/internal/store/chat"
	"github.com/lapitskyss/chat-service/internal/store/problem"
	"github.com/lapitskyss/chat-service/internal/types"
)

func (r *Repo) Create(ctx context.Context, chatID types.ChatID) (types.ProblemID, error) {
	p, err := r.db.Problem(ctx).
		Create().
		SetChatID(chatID).
		Save(ctx)
	if err != nil {
		return types.ProblemIDNil, fmt.Errorf("create problem: %v", err)
	}
	return p.ID, nil
}

func (r *Repo) CreateIfNotExists(ctx context.Context, chatID types.ChatID) (types.ProblemID, error) {
	problemID, err := r.db.Problem(ctx).
		Query().
		Unique(false).
		Where(
			problem.HasChatWith(chat.ID(chatID)),
			problem.ResolvedAtIsNil(),
		).
		OnlyID(ctx)
	if err != nil {
		if store.IsNotFound(err) {
			return r.Create(ctx, chatID)
		}
		return types.ProblemIDNil, fmt.Errorf("get problem: %v", err)
	}

	return problemID, nil
}
