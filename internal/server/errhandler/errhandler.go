package errhandler

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"

	"github.com/lapitskyss/chat-service/internal/errors"
)

var _ echo.HTTPErrorHandler = Handler{}.Handle

//go:generate options-gen -out-filename=errhandler_options.gen.go -from-struct=Options
type Options struct {
	logger          *zap.Logger                                    `option:"mandatory" validate:"required"`
	productionMode  bool                                           `option:"mandatory"`
	responseBuilder func(code int, msg string, details string) any `option:"mandatory" validate:"required"`
}

type Handler struct {
	lg              *zap.Logger
	productionMode  bool
	responseBuilder func(code int, msg string, details string) any
}

func New(opts Options) (Handler, error) {
	if err := opts.Validate(); err != nil {
		return Handler{}, fmt.Errorf("validate options: %v", err)
	}
	return Handler{
		lg:              opts.logger,
		productionMode:  opts.productionMode,
		responseBuilder: opts.responseBuilder,
	}, nil
}

func (h Handler) Handle(err error, eCtx echo.Context) {
	code, msg, details := h.processError(err)

	resp := h.responseBuilder(code, msg, details)

	if err2 := eCtx.JSON(http.StatusOK, resp); err2 != nil {
		h.lg.Error("error handler JSON", zap.Error(err))
	}
}

func (h Handler) processError(err error) (code int, msg string, details string) {
	code, msg, details = errors.ProcessServerError(err)

	// If production mode is ON method should return only code and message and hide details.
	if h.productionMode {
		details = ""
	}
	return
}
