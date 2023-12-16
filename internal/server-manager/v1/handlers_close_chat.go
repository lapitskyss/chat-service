package managerv1

import (
	"errors"
	"fmt"

	"github.com/labstack/echo/v4"

	"github.com/lapitskyss/chat-service/internal/middlewares"
	closechat "github.com/lapitskyss/chat-service/internal/usecases/manager/close-chat"
)

func (h Handlers) PostCloseChat(c echo.Context, params PostCloseChatParams) error {
	ctx := c.Request().Context()
	managerID := middlewares.MustUserID(c)

	var req CloseChatRequest
	if err := c.Bind(&req); err != nil {
		return fmt.Errorf("bind request: %w", err)
	}
	err := h.closeChat.Handle(ctx, closechat.Request{
		ID:        params.XRequestID,
		ManagerID: managerID,
		ChatID:    req.ChatId,
	})
	if err != nil {
		if errors.Is(err, closechat.ErrInvalidRequest) {
			return ErrBadRequest("invalid request", err)
		}
		if errors.Is(err, closechat.ErrNoActiveProblemInChat) {
			return ErrServer(ErrorCodeNoActiveProblemInChat, "no active problem in chat", err)
		}
		return fmt.Errorf("handle `close chat` use case: %v", err)
	}

	var data interface{}
	return Success(c, CloseChatResponse{
		Data: &data,
	})
}
