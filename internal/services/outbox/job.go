package outbox

import (
	"context"
	"time"
)

type Job interface {
	Name() string

	Handle(ctx context.Context, payload string) error

	// ExecutionTimeout is the time given to the queue handler to execute the task.
	// If the ExecutionTimeout is exceeded, the execution is aborted, the attempt is counted,
	// and the repetition will be performed.
	ExecutionTimeout() time.Duration

	// MaxAttempts is the maximum number of attempts to run the task.
	// An attempt is counted if the task was not completed due to an unknown error.
	// When MaxAttempts() is exceeded, the task moves to the dlq (dead letter queue) table.
	MaxAttempts() int
}
