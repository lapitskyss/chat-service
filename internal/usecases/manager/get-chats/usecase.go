package getchats

import (
	"context"
	"errors"
	"fmt"

	chatsrepo "github.com/lapitskyss/chat-service/internal/repositories/chats"
	"github.com/lapitskyss/chat-service/internal/types"
)

//go:generate mockgen -source=$GOFILE -destination=mocks/usecase_mock.gen.go -package=getchatsmocks

var ErrInvalidRequest = errors.New("invalid request")

type chatsRepository interface {
	AllWithOpenProblemsForManager(ctx context.Context, managerID types.UserID) ([]chatsrepo.Chat, error)
}

//go:generate options-gen -out-filename=usecase_options.gen.go -from-struct=Options
type Options struct {
	chatRepo chatsRepository `option:"mandatory" validate:"required"`
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

	chats, err := u.chatRepo.AllWithOpenProblemsForManager(ctx, req.ManagerID)
	if err != nil {
		return Response{}, fmt.Errorf("get all chats with open problems: %v", err)
	}

	return Response{
		Chats: adaptChats(chats),
	}, nil
}
