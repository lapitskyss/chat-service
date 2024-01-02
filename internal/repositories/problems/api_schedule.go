package problemsrepo

import (
	"context"
	"errors"
	"fmt"

	"entgo.io/ent/dialect/sql"

	"github.com/lapitskyss/chat-service/internal/store"
	"github.com/lapitskyss/chat-service/internal/store/message"
	"github.com/lapitskyss/chat-service/internal/store/problem"
	"github.com/lapitskyss/chat-service/internal/types"
)

var ErrRequestIDNotFound = errors.New("problem request id not found")

func (r *Repo) AllAvailableForManager(ctx context.Context, limit int) ([]Problem, error) {
	if limit <= 0 {
		return nil, errors.New("invalid limit")
	}
	problems, err := r.db.Problem(ctx).
		Query().
		Where(
			problem.ManagerIDIsNil(),
			problem.HasMessagesWith(message.IsVisibleForManager(true)),
		).
		Order(problem.ByCreatedAt(sql.OrderAsc())).
		Limit(limit).
		All(ctx)
	if err != nil {
		return nil, fmt.Errorf("get available problems for manager: %v", err)
	}
	return adaptProblems(problems), nil
}

func (r *Repo) SetManager(ctx context.Context, problemID types.ProblemID, managerID types.UserID) error {
	err := r.db.Problem(ctx).
		UpdateOneID(problemID).
		Where(problem.ManagerIDIsNil()).
		SetManagerID(managerID).
		Exec(ctx)
	if err != nil {
		return fmt.Errorf("set manager %s to problem %s: %v", managerID, problemID, err)
	}
	return nil
}

func (r *Repo) GetProblemRequestID(ctx context.Context, problemID types.ProblemID) (types.RequestID, error) {
	p, err := r.db.Problem(ctx).Query().
		WithMessages(func(query *store.MessageQuery) {
			query.Where(message.IsVisibleForManager(true)).
				Order(message.ByCreatedAt(sql.OrderAsc())).
				Limit(1)
		}).
		Where(problem.ID(problemID)).
		Where(problem.ManagerIDIsNil()).
		First(ctx)
	if err != nil {
		return types.RequestIDNil, fmt.Errorf("get problem request id: %v", err)
	}

	if len(p.Edges.Messages) != 1 {
		return types.RequestIDNil, ErrRequestIDNotFound
	}

	return p.Edges.Messages[0].InitialRequestID, nil
}
