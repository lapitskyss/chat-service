package managerscheduler

import (
	"context"
	"errors"
	"fmt"
	"time"

	"go.uber.org/zap"

	messagesrepo "github.com/lapitskyss/chat-service/internal/repositories/messages"
	problemsrepo "github.com/lapitskyss/chat-service/internal/repositories/problems"
	managerpool "github.com/lapitskyss/chat-service/internal/services/manager-pool"
	managerassignedtoproblemjob "github.com/lapitskyss/chat-service/internal/services/outbox/jobs/manager-assigned-to-problem"
	"github.com/lapitskyss/chat-service/internal/types"
)

const serviceName = "manager-scheduler"

type outboxService interface {
	Put(ctx context.Context, name, payload string, availableAt time.Time) (types.JobID, error)
}

type messagesRepository interface {
	CreateServiceMsg(
		ctx context.Context,
		reqID types.RequestID,
		problemID types.ProblemID,
		chatID types.ChatID,
		msgBody string,
		visibleForClient bool,
		visibleForManager bool,
	) (*messagesrepo.Message, error)
}

type problemsRepository interface {
	AllAvailableForManager(ctx context.Context, limit int) ([]problemsrepo.Problem, error)
	SetManager(ctx context.Context, problemID types.ProblemID, managerID types.UserID) error
	GetProblemRequestID(ctx context.Context, problemID types.ProblemID) (types.RequestID, error)
}

type transactor interface {
	RunInTx(ctx context.Context, f func(context.Context) error) error
}

//go:generate options-gen -out-filename=service_options.gen.go -from-struct=Options
type Options struct {
	period time.Duration `option:"mandatory" validate:"min=100ms,max=1m"`

	mngrPool managerpool.Pool `option:"mandatory" validate:"required"`

	msgRepo     messagesRepository `option:"mandatory" validate:"required"`
	outBox      outboxService      `option:"mandatory" validate:"required"`
	problemRepo problemsRepository `option:"mandatory" validate:"required"`
	txtor       transactor         `option:"mandatory" validate:"required"`
}

type Service struct {
	Options
	lg *zap.Logger
}

func New(opts Options) (*Service, error) {
	if err := opts.Validate(); err != nil {
		return nil, fmt.Errorf("validate options: %v", err)
	}
	return &Service{
		Options: opts,
		lg:      zap.L().Named(serviceName),
	}, nil
}

func (s *Service) Run(ctx context.Context) error {
	for {
		err := s.run(ctx)
		if err != nil {
			return fmt.Errorf("run manager scheduler: %v", err)
		}

		select {
		case <-time.NewTimer(s.period).C:
		case <-ctx.Done():
			return nil
		}
	}
}

func (s *Service) run(ctx context.Context) error {
	limit := s.mngrPool.Size()
	if limit == 0 {
		s.lg.Debug("no available managers")
		return nil
	}

	problems, err := s.problemRepo.AllAvailableForManager(ctx, limit)
	if err != nil {
		return fmt.Errorf("get problems without manager: %v", err)
	}
	if len(problems) == 0 {
		s.lg.Debug("no new problems")
		return nil
	}

	for _, problem := range problems {
		err = s.assignManagerToProblem(ctx, problem)
		if err != nil {
			return fmt.Errorf("assign manager to problem %s: %v", problem.ID, err)
		}
	}
	return nil
}

func (s *Service) assignManagerToProblem(ctx context.Context, problem problemsrepo.Problem) (errReturned error) {
	managerID, err := s.mngrPool.Get(ctx)
	if err != nil {
		if errors.Is(err, managerpool.ErrNoAvailableManagers) {
			return nil
		}
		return fmt.Errorf("get manager from pool: %v", err)
	}
	defer func() {
		if errReturned != nil {
			if err := s.mngrPool.Put(ctx, managerID); err != nil {
				s.lg.Error("cannot put manager back in the pool",
					zap.Error(err),
					zap.Stringer("manager_id", managerID))
			}
		}
	}()

	err = s.txtor.RunInTx(ctx, func(ctx context.Context) error {
		requestID, err := s.problemRepo.GetProblemRequestID(ctx, problem.ID)
		if err != nil {
			return fmt.Errorf("get problem request id: %v", err)
		}

		err = s.problemRepo.SetManager(ctx, problem.ID, managerID)
		if err != nil {
			return fmt.Errorf("assign manager to problem: %v", err)
		}

		msg, err := s.msgRepo.CreateServiceMsg(
			ctx,
			requestID,
			problem.ID,
			problem.ChatID,
			fmt.Sprintf("Manager %s will answer you", managerID),
			true,
			false,
		)
		if err != nil {
			return fmt.Errorf("create service message: %v", err)
		}

		payload, err := managerassignedtoproblemjob.MarshalPayload(msg.ID)
		if err != nil {
			return fmt.Errorf("marshal message id: %v", err)
		}
		_, err = s.outBox.Put(ctx, managerassignedtoproblemjob.Name, payload, time.Time{})
		if err != nil {
			return fmt.Errorf("put outbox message: %v", err)
		}

		s.lg.Info("set manager for problem",
			zap.Stringer("manager_id", managerID),
			zap.Stringer("problem_id", problem.ID),
		)
		return nil
	})
	if err != nil {
		return fmt.Errorf("assign manager to problem transaction: %v", err)
	}
	return nil
}
