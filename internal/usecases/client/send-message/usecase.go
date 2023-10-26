package sendmessage

import (
	"context"
	"errors"
	"fmt"

	messagesrepo "github.com/lapitskyss/chat-service/internal/repositories/messages"
	"github.com/lapitskyss/chat-service/internal/types"
)

//go:generate mockgen -source=$GOFILE -destination=mocks/usecase_mock.gen.go -package=sendmessagemocks

var (
	ErrInvalidRequest    = errors.New("invalid request")
	ErrChatNotCreated    = errors.New("chat not created")
	ErrProblemNotCreated = errors.New("problem not created")
)

type chatsRepository interface {
	CreateIfNotExists(ctx context.Context, userID types.UserID) (types.ChatID, error)
}

type messagesRepository interface {
	GetMessageByRequestID(ctx context.Context, reqID types.RequestID) (*messagesrepo.Message, error)
	CreateClientVisible(
		ctx context.Context,
		reqID types.RequestID,
		problemID types.ProblemID,
		chatID types.ChatID,
		authorID types.UserID,
		msgBody string,
	) (*messagesrepo.Message, error)
}

type problemsRepository interface {
	CreateIfNotExists(ctx context.Context, chatID types.ChatID) (types.ProblemID, error)
}

type transactor interface {
	RunInTx(ctx context.Context, f func(context.Context) error) error
}

//go:generate options-gen -out-filename=usecase_options.gen.go -from-struct=Options
type Options struct {
	chatRepo    chatsRepository    `option:"mandatory" validate:"required"`
	msgRepo     messagesRepository `option:"mandatory" validate:"required"`
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
		return Response{}, ErrInvalidRequest
	}

	var msg *messagesrepo.Message

	err := u.tr.RunInTx(ctx, func(ctx context.Context) error {
		var err error
		msg, err = u.msgRepo.GetMessageByRequestID(ctx, req.ID)
		if err == nil {
			return nil
		}
		if err != nil && !errors.Is(err, messagesrepo.ErrMsgNotFound) {
			return fmt.Errorf("message repo: %v", err)
		}
		chatID, err := u.chatRepo.CreateIfNotExists(ctx, req.ClientID)
		if err != nil {
			return fmt.Errorf("chat repo: %v: %w", err, ErrChatNotCreated)
		}
		problemID, err := u.problemRepo.CreateIfNotExists(ctx, chatID)
		if err != nil {
			return fmt.Errorf("problem repo: %v: %w", err, ErrProblemNotCreated)
		}
		msg, err = u.msgRepo.CreateClientVisible(ctx, req.ID, problemID, chatID, req.ClientID, req.MessageBody)
		if err != nil {
			return fmt.Errorf("message repo: %v", err)
		}
		return nil
	})
	if err != nil {
		return Response{}, err
	}

	return Response{
		AuthorID:  msg.AuthorID,
		MessageID: msg.ID,
		CreatedAt: msg.CreatedAt,
	}, nil
}
