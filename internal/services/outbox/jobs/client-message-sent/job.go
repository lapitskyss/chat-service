package clientmessagesentjob

import (
	"context"
	"fmt"

	"golang.org/x/sync/errgroup"

	messagesrepo "github.com/lapitskyss/chat-service/internal/repositories/messages"
	eventstream "github.com/lapitskyss/chat-service/internal/services/event-stream"
	"github.com/lapitskyss/chat-service/internal/services/outbox"
	"github.com/lapitskyss/chat-service/internal/types"
)

const Name = "client-message-sent"

type messageRepository interface {
	GetMessageByIDWithManager(ctx context.Context, msgID types.MessageID) (*messagesrepo.MessageWithManager, error)
}

type eventStream interface {
	Publish(ctx context.Context, userID types.UserID, event eventstream.Event) error
}

//go:generate options-gen -out-filename=job_options.gen.go -from-struct=Options
type Options struct {
	msgRepo     messageRepository `option:"mandatory" validate:"required"`
	eventStream eventStream       `option:"mandatory" validate:"required"`
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

	msg, err := j.msgRepo.GetMessageByIDWithManager(ctx, messageID)
	if err != nil {
		return fmt.Errorf("message repo, get message by id: %v", err)
	}

	eg, ctx := errgroup.WithContext(ctx)

	eg.Go(func() error {
		err = j.eventStream.Publish(ctx, msg.AuthorID, eventstream.NewMessageSentEvent(
			types.NewEventID(),
			msg.RequestID,
			msg.ID,
		))
		if err != nil {
			return fmt.Errorf("event stream, publish message sent event: %v", err)
		}
		return nil
	})

	eg.Go(func() error {
		if msg.ManagerID.IsZero() {
			return nil
		}
		err = j.eventStream.Publish(ctx, msg.ManagerID, eventstream.NewNewMessageEvent(
			types.NewEventID(),
			msg.RequestID,
			msg.ChatID,
			msg.ID,
			msg.AuthorID,
			msg.CreatedAt,
			msg.Body,
			msg.IsService,
		))
		if err != nil {
			return fmt.Errorf("event stream, publish new message event: %v", err)
		}
		return nil
	})

	err = eg.Wait()
	if err != nil {
		return fmt.Errorf("errgroup wait: %v", err)
	}

	return nil
}
