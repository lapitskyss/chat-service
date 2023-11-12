package managerload

import (
	"context"
	"errors"
	"fmt"

	"github.com/lapitskyss/chat-service/internal/types"
)

var ErrMaxProblemsCountReached = errors.New("manager max problem count reached")

func (s *Service) CanManagerTakeProblem(ctx context.Context, managerID types.UserID) (bool, error) {
	n, err := s.problemsRepo.GetManagerOpenProblemsCount(ctx, managerID)
	if err != nil {
		return false, fmt.Errorf("problem repo: %v", err)
	}
	if n >= s.maxProblemsAtTime {
		return false, ErrMaxProblemsCountReached
	}
	return true, nil
}
