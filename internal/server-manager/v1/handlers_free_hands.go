package managerv1

import (
	"errors"
	"fmt"

	"github.com/labstack/echo/v4"

	"github.com/lapitskyss/chat-service/internal/middlewares"
	gethistory "github.com/lapitskyss/chat-service/internal/usecases/client/get-history"
	canreceiveproblems "github.com/lapitskyss/chat-service/internal/usecases/manager/can-receive-problems"
	freehands "github.com/lapitskyss/chat-service/internal/usecases/manager/free-hands"
)

func (h Handlers) PostGetFreeHandsBtnAvailability(c echo.Context, params PostGetFreeHandsBtnAvailabilityParams) error {
	ctx := c.Request().Context()
	managerID := middlewares.MustUserID(c)

	response, err := h.canReceiveProblems.Handle(ctx, canreceiveproblems.Request{
		ID:        params.XRequestID,
		ManagerID: managerID,
	})
	if err != nil {
		if errors.Is(err, gethistory.ErrInvalidRequest) {
			return ErrBadRequest("invalid request", err)
		}
		return fmt.Errorf("handle `can receive problems` use case: %v", err)
	}

	return Success(c, GetFreeHandsBtnAvailabilityResponse{
		Data: &FreeHandsBtnAvailability{
			Available: response.Result,
		},
	})
}

func (h Handlers) PostFreeHands(c echo.Context, params PostFreeHandsParams) error {
	ctx := c.Request().Context()
	managerID := middlewares.MustUserID(c)

	err := h.freeHands.Handle(ctx, freehands.Request{
		ID:        params.XRequestID,
		ManagerID: managerID,
	})
	if err != nil {
		if errors.Is(err, freehands.ErrInvalidRequest) {
			return ErrBadRequest("invalid request", err)
		}
		if errors.Is(err, freehands.ErrManagerOverloaded) {
			return ErrServer(ErrorManagerOverloaded, "manager overloaded", err)
		}
		return fmt.Errorf("handle `free hands` use case: %v", err)
	}

	var data interface{}
	return Success(c, FreeHandsResponse{
		Data: &data,
	})
}
