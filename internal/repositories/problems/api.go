package problemsrepo

import (
	"context"
	"fmt"

	"github.com/lapitskyss/chat-service/internal/store"
	"github.com/lapitskyss/chat-service/internal/store/problem"
	"github.com/lapitskyss/chat-service/internal/types"
)

func (r *Repo) Create(ctx context.Context, chatID types.ChatID) (types.ProblemID, error) {
	problemID := types.NewProblemID()
	err := r.db.Problem(ctx).
		Create().
		SetID(problemID).
		SetChatID(chatID).
		Exec(ctx)
	if err != nil {
		return types.ProblemID{}, fmt.Errorf("create problem: %v", err)
	}
	return problemID, nil
}

func (r *Repo) CreateIfNotExists(ctx context.Context, chatID types.ChatID) (types.ProblemID, error) {
	problemID, err := r.db.Problem(ctx).
		Query().
		Where(problem.ChatID(chatID)).
		Where(problem.ResolvedAtIsNil()).
		OnlyID(ctx)
	if err != nil {
		if store.IsNotFound(err) {
			return r.Create(ctx, chatID)
		}
		return types.ProblemID{}, fmt.Errorf("get problem: %v", err)
	}

	return problemID, nil
}
