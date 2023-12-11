package outbox

import (
	"context"
	"errors"
	"fmt"
	"time"

	"go.uber.org/zap"

	jobsrepo "github.com/lapitskyss/chat-service/internal/repositories/jobs"
	"github.com/lapitskyss/chat-service/internal/types"
)

const serviceName = "outbox"

var ErrJobAlreadyExist = errors.New("job already exits")

type jobsRepository interface {
	CreateJob(ctx context.Context, name, payload string, availableAt time.Time) (types.JobID, error)
	FindAndReserveJob(ctx context.Context, until time.Time) (jobsrepo.Job, error)
	CreateFailedJob(ctx context.Context, name, payload, reason string) error
	DeleteJob(ctx context.Context, jobID types.JobID) error
}

type transactor interface {
	RunInTx(ctx context.Context, f func(context.Context) error) error
}

//go:generate options-gen -out-filename=service_options.gen.go -from-struct=Options
type Options struct {
	workers    int           `option:"mandatory" validate:"min=1,max=32"`
	idleTime   time.Duration `option:"mandatory" validate:"min=100ms,max=10s"`
	reserveFor time.Duration `option:"mandatory" validate:"min=1s,max=10m"`

	jobsRepo jobsRepository `option:"mandatory" validate:"required"`
	tr       transactor     `option:"mandatory" validate:"required"`
}

type Service struct {
	Options

	jobs map[string]Job
}

func New(opts Options) (*Service, error) {
	if err := opts.Validate(); err != nil {
		return nil, fmt.Errorf("validate options: %v", err)
	}
	return &Service{
		Options: opts,
		jobs:    map[string]Job{},
	}, nil
}

func (s *Service) RegisterJob(job Job) error {
	if _, exits := s.jobs[job.Name()]; exits {
		return ErrJobAlreadyExist
	}

	s.jobs[job.Name()] = job
	return nil
}

func (s *Service) MustRegisterJob(job Job) {
	err := s.RegisterJob(job)
	if err != nil {
		panic(err)
	}
}

func (s *Service) GetJob(name string) (Job, bool) {
	job, exist := s.jobs[name]
	return job, exist
}

func (s *Service) Run(ctx context.Context) error {
	for i := 0; i < s.workers; i++ {
		go s.runWorker(ctx)
	}
	return nil
}

func (s *Service) runWorker(ctx context.Context) {
	for {
		job, err := s.jobsRepo.FindAndReserveJob(ctx, time.Now().Add(s.reserveFor))
		if err != nil {
			if !errors.Is(err, jobsrepo.ErrNoJobs) {
				logError("run worker db query", err)
			}

			// Fall asleep for idleTime
			select {
			case <-ctx.Done():
				return
			case <-time.After(s.idleTime):
			}

			continue
		}

		// Find registered job
		svcJob, exist := s.GetJob(job.Name)
		if !exist {
			s.jobFailed(ctx, job, fmt.Sprintf("job with name '%s' is not registered", job.Name))
			continue
		}

		// Check is number of attempts already achieved
		if job.Attempts > svcJob.MaxAttempts() {
			s.jobFailed(ctx, job, "maximum number of job attempts achieved")
			continue
		}

		// Execute job with provided timeout
		err = s.runJob(ctx, svcJob, job.Payload)
		if err != nil {
			logError("execute job", err)

			// Check is number of attempts achieved
			if job.Attempts >= svcJob.MaxAttempts() {
				s.jobFailed(ctx, job, "maximum number of job attempts achieved")
			}

			continue
		}

		// Remove job from the queue
		s.removeJob(ctx, job)
	}
}

func (s *Service) runJob(ctx context.Context, job Job, payload string) error {
	ctx, cancel := context.WithTimeout(ctx, job.ExecutionTimeout())
	defer cancel()

	return job.Handle(ctx, payload)
}

func (s *Service) removeJob(_ context.Context, job jobsrepo.Job) {
	// Intentionally delete job with context.WithTimeout() to avoid case when job is handled,
	// but ctx is already closed before deleting.
	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()

	err := s.jobsRepo.DeleteJob(ctx, job.ID)
	if err != nil {
		logError("delete job from queue", err)
	}
}

func (s *Service) jobFailed(ctx context.Context, job jobsrepo.Job, reason string) {
	err := s.tr.RunInTx(ctx, func(ctx context.Context) error {
		err := s.jobsRepo.CreateFailedJob(ctx, job.Name, job.Payload, reason)
		if err != nil {
			return err
		}
		err = s.jobsRepo.DeleteJob(ctx, job.ID)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		logError("move job to failed", err)
	}
}

func logError(msg string, err error) {
	zap.L().Named(serviceName).Error(msg, zap.Error(err))
}
