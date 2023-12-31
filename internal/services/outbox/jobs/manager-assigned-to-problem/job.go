package managerassignedtoproblemjob

import (
	"context"
	"fmt"

	"golang.org/x/sync/errgroup"

	messagesrepo "github.com/lapitskyss/chat-service/internal/repositories/messages"
	eventstream "github.com/lapitskyss/chat-service/internal/services/event-stream"
	"github.com/lapitskyss/chat-service/internal/services/outbox"
	"github.com/lapitskyss/chat-service/internal/types"
)

const Name = "manager-assigned-to-problem"

type managerLoadService interface {
	CanManagerTakeProblem(ctx context.Context, managerID types.UserID) (bool, error)
}

type messageRepository interface {
	GetServiceMessageByID(ctx context.Context, id types.MessageID) (*messagesrepo.ServiceMessage, error)
}

type eventStream interface {
	Publish(ctx context.Context, userID types.UserID, event eventstream.Event) error
}

//go:generate options-gen -out-filename=job_options.gen.go -from-struct=Options
type Options struct {
	managerLoadSvc managerLoadService `option:"mandatory" validate:"required"`
	msgRepo        messageRepository  `option:"mandatory" validate:"required"`
	eventStream    eventStream        `option:"mandatory" validate:"required"`
}

type Job struct {
	Options
	outbox.DefaultJob
}

func New(opts Options) (*Job, error) {
	if err := opts.Validate(); err != nil {
		return &Job{}, fmt.Errorf("validate options: %v", err)
	}
	return &Job{Options: opts}, nil
}

func (j *Job) Name() string {
	return Name
}

func (j *Job) Handle(ctx context.Context, payload string) error {
	messageID, err := UnmarshalPayload(payload)
	if err != nil {
		return fmt.Errorf("unmarshal payload: %v", err)
	}

	msg, err := j.msgRepo.GetServiceMessageByID(ctx, messageID)
	if err != nil {
		return fmt.Errorf("message repo, get message by id: %v", err)
	}

	canTakeMoreProblems, err := j.managerLoadSvc.CanManagerTakeProblem(ctx, msg.ManagerID)
	if err != nil {
		return fmt.Errorf("manager load svc, can manager take problem: %v", err)
	}

	eg, ctx := errgroup.WithContext(ctx)

	eg.Go(func() error {
		err = j.eventStream.Publish(ctx, msg.ManagerID, eventstream.NewNewChatEvent(
			types.NewEventID(),
			msg.ChatID,
			msg.ClientID,
			msg.RequestID,
			canTakeMoreProblems,
		))
		if err != nil {
			return fmt.Errorf("event stream, publish new chat event: %v", err)
		}
		return nil
	})

	eg.Go(func() error {
		err = j.eventStream.Publish(ctx, msg.ClientID, eventstream.NewNewMessageEvent(
			types.NewEventID(),
			msg.RequestID,
			msg.ChatID,
			msg.ID,
			types.UserIDNil,
			msg.CreatedAt,
			msg.Body,
			true,
		))
		if err != nil {
			return fmt.Errorf("event stream, publish new message event: %v", err)
		}
		return nil
	})

	if err = eg.Wait(); err != nil {
		return fmt.Errorf("errgroup wait: %v", err)
	}

	return nil
}
