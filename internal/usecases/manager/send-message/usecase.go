package sendmessage

import (
	"context"
	"errors"
	"fmt"
	"time"

	messagesrepo "github.com/lapitskyss/chat-service/internal/repositories/messages"
	sendmanagermessagejob "github.com/lapitskyss/chat-service/internal/services/outbox/jobs/send-manager-message"
	"github.com/lapitskyss/chat-service/internal/types"
)

//go:generate mockgen -source=$GOFILE -destination=mocks/usecase_mock.gen.go -package=sendmessagemocks

var ErrInvalidRequest = errors.New("invalid request")

type messagesRepository interface {
	CreateFullVisible(
		ctx context.Context,
		reqID types.RequestID,
		problemID types.ProblemID,
		chatID types.ChatID,
		authorID types.UserID,
		msgBody string,
	) (*messagesrepo.Message, error)
}

type outboxService interface {
	Put(ctx context.Context, name, payload string, availableAt time.Time) (types.JobID, error)
}

type problemsRepository interface {
	GetAssignedProblemID(ctx context.Context, managerID types.UserID, chatID types.ChatID) (types.ProblemID, error)
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

func (u UseCase) Handle(ctx context.Context, req Request) (Response, error) {
	if err := req.Validate(); err != nil {
		return Response{}, fmt.Errorf("validate request: %w: %v", ErrInvalidRequest, err)
	}

	problemID, err := u.problemRepo.GetAssignedProblemID(ctx, req.ManagerID, req.ChatID)
	if err != nil {
		return Response{}, fmt.Errorf("get assigned problem id: %v", err)
	}

	var msg *messagesrepo.Message
	err = u.tr.RunInTx(ctx, func(ctx context.Context) error {
		msg, err = u.msgRepo.CreateFullVisible(ctx, req.ID, problemID, req.ChatID, req.ManagerID, req.MessageBody)
		if err != nil {
			return fmt.Errorf("msg repo, create full visible: %v", err)
		}
		payload, err := sendmanagermessagejob.MarshalPayload(msg.ID)
		if err != nil {
			return fmt.Errorf("marshal message id: %v", err)
		}
		_, err = u.outBoxSvc.Put(ctx, sendmanagermessagejob.Name, payload, time.Time{})
		if err != nil {
			return fmt.Errorf("put outbox message: %v", err)
		}
		return nil
	})
	if err != nil {
		return Response{}, fmt.Errorf("run transaction: %v", err)
	}

	return Response{
		MessageID: msg.ID,
		CreatedAt: msg.CreatedAt,
	}, nil
}
