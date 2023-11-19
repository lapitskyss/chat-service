package outbox

import (
	"context"
	"fmt"
	"time"

	"github.com/lapitskyss/chat-service/internal/types"
)

func (s *Service) Put(ctx context.Context, name, payload string, availableAt time.Time) (types.JobID, error) {
	jobID, err := s.jobsRepo.CreateJob(ctx, name, payload, availableAt)
	if err != nil {
		return types.JobIDNil, fmt.Errorf("job repository: %v", err)
	}
	return jobID, nil
}
