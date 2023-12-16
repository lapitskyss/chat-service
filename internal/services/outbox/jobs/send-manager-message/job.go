package sendmanagermessagejob

import (
	"context"
	"fmt"

	"golang.org/x/sync/errgroup"

	chatsrepo "github.com/lapitskyss/chat-service/internal/repositories/chats"
	messagesrepo "github.com/lapitskyss/chat-service/internal/repositories/messages"
	eventstream "github.com/lapitskyss/chat-service/internal/services/event-stream"
	msgproducer "github.com/lapitskyss/chat-service/internal/services/msg-producer"
	"github.com/lapitskyss/chat-service/internal/services/outbox"
	"github.com/lapitskyss/chat-service/internal/types"
)

const Name = "send-manager-message"

type messageRepository interface {
	GetMessageByID(ctx context.Context, msgID types.MessageID) (*messagesrepo.Message, error)
}

type chatRepository interface {
	GetChatByID(ctx context.Context, chatID types.ChatID) (*chatsrepo.Chat, error)
}

type eventStream interface {
	Publish(ctx context.Context, userID types.UserID, event eventstream.Event) error
}

type messageProducer interface {
	ProduceMessage(ctx context.Context, message msgproducer.Message) error
}

//go:generate options-gen -out-filename=job_options.gen.go -from-struct=Options
type Options struct {
	msgRepo     messageRepository `option:"mandatory" validate:"required"`
	chatRepo    chatRepository    `option:"mandatory" validate:"required"`
	eventStream eventStream       `option:"mandatory" validate:"required"`
	msgProducer messageProducer   `option:"mandatory" validate:"required"`
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

	msg, err := j.msgRepo.GetMessageByID(ctx, messageID)
	if err != nil {
		return fmt.Errorf("message repo, get message by id: %v", err)
	}
	chat, err := j.chatRepo.GetChatByID(ctx, msg.ChatID)
	if err != nil {
		return fmt.Errorf("chats repo, get message by id: %v", err)
	}

	err = j.msgProducer.ProduceMessage(ctx, msgproducer.Message{
		ID:         msg.ID,
		ChatID:     msg.ChatID,
		Body:       msg.Body,
		FromClient: false,
	})
	if err != nil {
		return fmt.Errorf("message producer, produce message: %v", err)
	}

	eg, ctx := errgroup.WithContext(ctx)

	eg.Go(func() error {
		err = j.eventStream.Publish(ctx, chat.ClientID, eventstream.NewNewMessageEvent(
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
			return fmt.Errorf("event stream, publish new message event to client: %v", err)
		}
		return nil
	})

	eg.Go(func() error {
		err = j.eventStream.Publish(ctx, msg.AuthorID, eventstream.NewNewMessageEvent(
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
			return fmt.Errorf("event stream, publish new message event to manager: %v", err)
		}
		return nil
	})

	err = eg.Wait()
	if err != nil {
		return fmt.Errorf("errgroup wait: %v", err)
	}

	return nil
}
