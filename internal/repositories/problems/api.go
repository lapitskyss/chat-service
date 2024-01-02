package problemsrepo

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/lapitskyss/chat-service/internal/store"
	"github.com/lapitskyss/chat-service/internal/store/chat"
	"github.com/lapitskyss/chat-service/internal/store/problem"
	"github.com/lapitskyss/chat-service/internal/types"
	"github.com/lapitskyss/chat-service/pkg/pointer"
)

var ErrProblemNotFound = errors.New("problem not found")

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

func (r *Repo) GetChatOpenProblem(ctx context.Context, chatID types.ChatID) (*Problem, error) {
	p, err := r.db.Problem(ctx).
		Query().
		Where(
			problem.ChatID(chatID),
			problem.ResolvedAtIsNil(),
		).
		First(ctx)
	if err != nil {
		if store.IsNotFound(err) {
			return nil, ErrProblemNotFound
		}
		return nil, fmt.Errorf("get chat open problem: %v", err)
	}
	return pointer.Ptr(adaptProblem(p)), nil
}

func (r *Repo) GetClientOpenProblem(ctx context.Context, clientID types.UserID) (*Problem, error) {
	p, err := r.db.Problem(ctx).
		Query().
		Where(
			problem.HasChatWith(chat.ClientID(clientID)),
			problem.ResolvedAtIsNil(),
		).
		First(ctx)
	if err != nil {
		if store.IsNotFound(err) {
			return nil, ErrProblemNotFound
		}
		return nil, fmt.Errorf("get chat open problem: %v", err)
	}
	return pointer.Ptr(adaptProblem(p)), nil
}

func (r *Repo) GetManagerOpenProblemsCount(ctx context.Context, managerID types.UserID) (int, error) {
	n, err := r.db.Problem(ctx).
		Query().
		Where(
			problem.ManagerID(managerID),
			problem.ResolvedAtIsNil(),
		).
		Count(ctx)
	if err != nil {
		return 0, fmt.Errorf("get manager open problem count: %v", err)
	}
	return n, nil
}

func (r *Repo) GetAssignedProblemID(
	ctx context.Context,
	managerID types.UserID,
	chatID types.ChatID,
) (types.ProblemID, error) {
	p, err := r.db.Problem(ctx).
		Query().
		Where(
			problem.ManagerID(managerID),
			problem.ResolvedAtIsNil(),
			problem.ChatID(chatID),
		).
		OnlyID(ctx)
	if err != nil {
		if store.IsNotFound(err) {
			return types.ProblemIDNil, ErrProblemNotFound
		}
		return types.ProblemIDNil, fmt.Errorf("get manager problem for chat: %v", err)
	}
	return p, nil
}

func (r *Repo) ResolveProblem(ctx context.Context, problemID types.ProblemID) error {
	err := r.db.Problem(ctx).
		UpdateOneID(problemID).
		SetResolvedAt(time.Now()).
		Exec(ctx)
	if err != nil {
		return fmt.Errorf("resolve problem: %v", err)
	}
	return nil
}
