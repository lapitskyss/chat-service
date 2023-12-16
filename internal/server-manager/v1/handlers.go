package managerv1

import (
	"context"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"

	svcerr "github.com/lapitskyss/chat-service/internal/errors"
	canreceiveproblems "github.com/lapitskyss/chat-service/internal/usecases/manager/can-receive-problems"
	freehands "github.com/lapitskyss/chat-service/internal/usecases/manager/free-hands"
	getchathistory "github.com/lapitskyss/chat-service/internal/usecases/manager/get-chat-history"
	getchats "github.com/lapitskyss/chat-service/internal/usecases/manager/get-chats"
)

var _ ServerInterface = (*Handlers)(nil)

//go:generate mockgen -source=$GOFILE -destination=mocks/handlers_mocks.gen.go -package=managerv1mocks

type canReceiveProblemsUseCase interface {
	Handle(ctx context.Context, req canreceiveproblems.Request) (canreceiveproblems.Response, error)
}

type freeHandsUseCase interface {
	Handle(ctx context.Context, req freehands.Request) error
}

type getChatHistoryUseCase interface {
	Handle(ctx context.Context, req getchathistory.Request) (getchathistory.Response, error)
}

type getChatsUseCase interface {
	Handle(ctx context.Context, req getchats.Request) (getchats.Response, error)
}

//go:generate options-gen -out-filename=handlers.gen.go -from-struct=Options
type Options struct {
	canReceiveProblems canReceiveProblemsUseCase `option:"mandatory" validate:"required"`
	freeHands          freeHandsUseCase          `option:"mandatory" validate:"required"`
	getChatHistory     getChatHistoryUseCase     `option:"mandatory" validate:"required"`
	getChats           getChatsUseCase           `option:"mandatory" validate:"required"`
}

type Handlers struct {
	Options
}

func NewHandlers(opts Options) (Handlers, error) {
	if err := opts.Validate(); err != nil {
		return Handlers{}, fmt.Errorf("validate options: %v", err)
	}
	return Handlers{Options: opts}, nil
}

func Success(c echo.Context, data any) error {
	return c.JSON(http.StatusOK, data)
}

func ErrBadRequest(msg string, err error) error {
	return svcerr.NewServerError(http.StatusBadRequest, msg, err)
}

func ErrServer(code ErrorCode, msg string, err error) error {
	return svcerr.NewServerError(int(code), msg, err)
}
