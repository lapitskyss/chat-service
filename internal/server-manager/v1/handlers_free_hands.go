package managerv1

import (
	"fmt"

	"github.com/labstack/echo/v4"

	"github.com/lapitskyss/chat-service/internal/middlewares"
	canreceiveproblems "github.com/lapitskyss/chat-service/internal/usecases/manager/can-receive-problems"
)

func (h Handlers) PostGetFreeHandsBtnAvailability(c echo.Context, params PostGetFreeHandsBtnAvailabilityParams) error {
	ctx := c.Request().Context()
	managerID := middlewares.MustUserID(c)

	response, err := h.canReceiveProblems.Handle(ctx, canreceiveproblems.Request{
		ID:        params.XRequestID,
		ManagerID: managerID,
	})
	if err != nil {
		return fmt.Errorf("handle `can receive problems` use case: %v", err)
	}

	return Success(c, GetFreeHandsBtnAvailabilityResponse{
		Data: &FreeHandsBtnAvailability{
			Available: response.Result,
		},
	})
}
