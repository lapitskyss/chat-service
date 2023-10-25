package clientv1

import (
	"context"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"

	gethistory "github.com/lapitskyss/chat-service/internal/usecases/client/get-history"
)

//go:generate mockgen -source=$GOFILE -destination=mocks/handlers_mocks.gen.go -package=clientv1mocks

type getHistoryUseCase interface {
	Handle(ctx context.Context, req gethistory.Request) (gethistory.Response, error)
}

//go:generate options-gen -out-filename=handlers.gen.go -from-struct=Options
type Options struct {
	getHistory getHistoryUseCase `option:"mandatory" validate:"required"`
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

func ErrBadRequest(data any) error {
	return echo.NewHTTPError(http.StatusBadRequest, data)
}
