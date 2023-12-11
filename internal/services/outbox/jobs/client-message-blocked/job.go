package clientmessageblockedjob

import (
	"context"
	"fmt"

	messagesrepo "github.com/lapitskyss/chat-service/internal/repositories/messages"
	eventstream "github.com/lapitskyss/chat-service/internal/services/event-stream"
	"github.com/lapitskyss/chat-service/internal/services/outbox"
	"github.com/lapitskyss/chat-service/internal/types"
)

const Name = "client-message-blocked"

type messageRepository interface {
	GetMessageByID(ctx context.Context, msgID types.MessageID) (*messagesrepo.Message, error)
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

	message, err := j.msgRepo.GetMessageByID(ctx, messageID)
	if err != nil {
		return fmt.Errorf("message repo, get message by id: %v", err)
	}

	event := eventstream.NewMessageBlockEvent(
		types.NewEventID(),
		message.RequestID,
		message.ID,
	)
	err = j.eventStream.Publish(ctx, message.AuthorID, event)
	if err != nil {
		return fmt.Errorf("event stream, publish message blick event: %v", err)
	}

	return nil
}
