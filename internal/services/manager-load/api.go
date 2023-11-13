package managerload

import (
	"context"
	"fmt"

	"github.com/lapitskyss/chat-service/internal/types"
)

func (s *Service) CanManagerTakeProblem(ctx context.Context, managerID types.UserID) (bool, error) {
	n, err := s.problemsRepo.GetManagerOpenProblemsCount(ctx, managerID)
	if err != nil {
		return false, fmt.Errorf("problem repo: %v", err)
	}
	if n >= s.maxProblemsAtTime {
		return false, nil
	}
	return true, nil
}
