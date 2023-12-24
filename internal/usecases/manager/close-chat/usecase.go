package closechat

import (
	"context"
	"errors"
	"fmt"
	"time"

	messagesrepo "github.com/lapitskyss/chat-service/internal/repositories/messages"
	problemsrepo "github.com/lapitskyss/chat-service/internal/repositories/problems"
	managerclosechatjob "github.com/lapitskyss/chat-service/internal/services/outbox/jobs/manager-close-chat"
	"github.com/lapitskyss/chat-service/internal/types"
)

//go:generate mockgen -source=$GOFILE -destination=mocks/usecase_mock.gen.go -package=closechatmocks

var (
	ErrInvalidRequest        = errors.New("invalid request")
	ErrNoActiveProblemInChat = errors.New("not active problems in chat")
)

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

type outboxService interface {
	Put(ctx context.Context, name, payload string, availableAt time.Time) (types.JobID, error)
}

type problemsRepository interface {
	GetAssignedProblemID(ctx context.Context, managerID types.UserID, chatID types.ChatID) (types.ProblemID, error)
	ResolveProblem(ctx context.Context, problemID types.ProblemID) error
}

type transactor interface {
	RunInTx(ctx context.Context, f func(context.Context) error) error
}

//go:generate options-gen -out-filename=usecase_options.gen.go -from-struct=Options
type Options struct {
	msgRepo     messagesRepository `option:"mandatory" validate:"required"`
	outBoxSvc   outboxService      `option:"mandatory" validate:"required"`
	problemRepo problemsRepository `option:"mandatory" validate:"required"`
	tr          transactor         `option:"mandatory" validate:"required"`
}

type UseCase struct {
	Options
}

func New(opts Options) (UseCase, error) {
	if err := opts.Validate(); err != nil {
		return UseCase{}, fmt.Errorf("validate options: %v", err)
	}
	return UseCase{Options: opts}, nil
}

func (u UseCase) Handle(ctx context.Context, req Request) error {
	if err := req.Validate(); err != nil {
		return fmt.Errorf("validate request: %w: %v", ErrInvalidRequest, err)
	}

	err := u.tr.RunInTx(ctx, func(ctx context.Context) error {
		problemID, err := u.problemRepo.GetAssignedProblemID(ctx, req.ManagerID, req.ChatID)
		if err != nil {
			if errors.Is(err, problemsrepo.ErrProblemNotFound) {
				return ErrNoActiveProblemInChat
			}
			return fmt.Errorf("problems repo, get assigned problem id: %v", err)
		}
		err = u.problemRepo.ResolveProblem(ctx, problemID)
		if err != nil {
			return fmt.Errorf("problems repo, resolve problem: %v", err)
		}
		msgBody := "Your question has been marked as resolved.\nThank you for being with us!"
		msg, err := u.msgRepo.CreateServiceMsg(ctx,
			req.ID,
			problemID,
			req.ChatID,
			msgBody,
			true,
			false,
		)
		if err != nil {
			return fmt.Errorf("msg repo, create service message: %v", err)
		}
		payload, err := managerclosechatjob.MarshalPayload(msg.ID)
		if err != nil {
			return fmt.Errorf("marshal message id: %v", err)
		}
		_, err = u.outBoxSvc.Put(ctx, managerclosechatjob.Name, payload, time.Time{})
		if err != nil {
			return fmt.Errorf("put outbox message: %v", err)
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("close chat tx: %w", err)
	}
	return nil
}
