package managertypingmessage

import (
	"context"
	"errors"
	"fmt"

	chatsrepo "github.com/lapitskyss/chat-service/internal/repositories/chats"
	problemsrepo "github.com/lapitskyss/chat-service/internal/repositories/problems"
	eventstream "github.com/lapitskyss/chat-service/internal/services/event-stream"
	"github.com/lapitskyss/chat-service/internal/types"
)

//go:generate mockgen -source=$GOFILE -destination=mocks/usecase_mock.gen.go -package=managertypingmessagemocks

var (
	ErrInvalidRequest = errors.New("invalid request")
	ErrChatNotFound   = errors.New("chat not found")
)

type chatRepository interface {
	GetChatByID(ctx context.Context, chatID types.ChatID) (*chatsrepo.Chat, error)
}

type problemsRepository interface {
	GetChatOpenProblem(ctx context.Context, chatID types.ChatID) (*problemsrepo.Problem, error)
}

type eventStream interface {
	Publish(ctx context.Context, userID types.UserID, event eventstream.Event) error
}

//go:generate options-gen -out-filename=usecase_options.gen.go -from-struct=Options
type Options struct {
	chatRepo    chatRepository     `option:"mandatory" validate:"required"`
	problemRepo problemsRepository `option:"mandatory" validate:"required"`
	eventStream eventStream        `option:"mandatory" validate:"required"`
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

	problem, err := u.problemRepo.GetChatOpenProblem(ctx, req.ChatID)
	if err != nil {
		return fmt.Errorf("problem repo, get chat open problem, %v", err)
	}
	if problem.ManagerID != req.ManagerID {
		return fmt.Errorf("problem not assigned to manager, %w", ErrChatNotFound)
	}

	chat, err := u.chatRepo.GetChatByID(ctx, req.ChatID)
	if err != nil {
		return fmt.Errorf("chat repo, get chat by id, %v", err)
	}

	// Send event NewTypingEvent to client
	err = u.eventStream.Publish(ctx, chat.ClientID, eventstream.NewTypingEvent(
		types.NewEventID(),
		req.ManagerID,
		req.ID,
	))
	if err != nil {
		return fmt.Errorf("event stream, publish message blick event: %v", err)
	}

	return nil
}
