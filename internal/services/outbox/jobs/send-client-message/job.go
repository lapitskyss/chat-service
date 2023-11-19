package sendclientmessagejob

import (
	"context"
	"fmt"

	messagesrepo "github.com/lapitskyss/chat-service/internal/repositories/messages"
	msgproducer "github.com/lapitskyss/chat-service/internal/services/msg-producer"
	"github.com/lapitskyss/chat-service/internal/services/outbox"
	"github.com/lapitskyss/chat-service/internal/types"
)

//go:generate mockgen -source=$GOFILE -destination=mocks/job_mock.gen.go -package=sendclientmessagejobmocks

const Name = "send-client-message"

type messageProducer interface {
	ProduceMessage(ctx context.Context, message msgproducer.Message) error
}

type messageRepository interface {
	GetMessageByID(ctx context.Context, msgID types.MessageID) (*messagesrepo.Message, error)
}

//go:generate options-gen -out-filename=job_options.gen.go -from-struct=Options
type Options struct {
	msgProducer messageProducer   `option:"mandatory" validate:"required"`
	msgRepo     messageRepository `option:"mandatory" validate:"required"`
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

	produceMsg := msgproducer.Message{
		ID:         message.ID,
		ChatID:     message.ChatID,
		Body:       message.Body,
		FromClient: true,
	}

	err = j.msgProducer.ProduceMessage(ctx, produceMsg)
	if err != nil {
		return fmt.Errorf("message producer, produce message: %v", err)
	}

	return nil
}
