package gethistory

import (
	"context"
	"errors"
	"fmt"

	"github.com/lapitskyss/chat-service/internal/cursor"
	messagesrepo "github.com/lapitskyss/chat-service/internal/repositories/messages"
	"github.com/lapitskyss/chat-service/internal/types"
)

//go:generate mockgen -source=$GOFILE -destination=mocks/usecase_mock.gen.go -package=gethistorymocks

var (
	ErrInvalidRequest = errors.New("invalid request")
	ErrInvalidCursor  = errors.New("invalid cursor")
)

type messagesRepository interface {
	GetClientChatMessages(
		ctx context.Context,
		clientID types.UserID,
		pageSize int,
		cursor *messagesrepo.Cursor,
	) ([]messagesrepo.Message, *messagesrepo.Cursor, error)
}

//go:generate options-gen -out-filename=usecase_options.gen.go -from-struct=Options
type Options struct {
	msgRepo messagesRepository `option:"mandatory" validate:"required"`
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

	cur, err := decodeCursor(req.Cursor)
	if err != nil {
		return Response{}, fmt.Errorf("decode cursor: %w: %v", ErrInvalidCursor, err)
	}

	messages, cur, err := u.msgRepo.GetClientChatMessages(ctx, req.ClientID, req.PageSize, cur)
	if err != nil {
		if errors.Is(err, messagesrepo.ErrInvalidCursor) {
			return Response{}, fmt.Errorf("get client chat messages: %w: %v", ErrInvalidCursor, err)
		}
		return Response{}, fmt.Errorf("get client chat messages: %v", err)
	}

	nextCursor, err := encodeCursor(cur)
	if err != nil {
		return Response{}, fmt.Errorf("encode cursore: %v", err)
	}

	return Response{
		Messages:   mapMassages(messages),
		NextCursor: nextCursor,
	}, nil
}

func decodeCursor(val string) (*messagesrepo.Cursor, error) {
	var cur *messagesrepo.Cursor
	if val == "" {
		return cur, nil
	}
	err := cursor.Decode(val, &cur)
	if err != nil {
		return nil, err
	}
	return cur, nil
}

func encodeCursor(val *messagesrepo.Cursor) (string, error) {
	if val == nil {
		return "", nil
	}
	nextCursor, err := cursor.Encode(val)
	if err != nil {
		return "", err
	}
	return nextCursor, nil
}
