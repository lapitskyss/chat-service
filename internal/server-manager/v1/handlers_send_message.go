package managerv1

import (
	"errors"
	"fmt"

	"github.com/labstack/echo/v4"

	"github.com/lapitskyss/chat-service/internal/middlewares"
	sendmessage "github.com/lapitskyss/chat-service/internal/usecases/manager/send-message"
)

func (h Handlers) PostSendMessage(c echo.Context, params PostSendMessageParams) error {
	ctx := c.Request().Context()
	managerID := middlewares.MustUserID(c)

	var req SendMessageRequest
	if err := c.Bind(&req); err != nil {
		return fmt.Errorf("bind request: %w", err)
	}

	res, err := h.sendMessage.Handle(ctx, sendmessage.Request{
		ID:          params.XRequestID,
		ManagerID:   managerID,
		ChatID:      req.ChatId,
		MessageBody: req.MessageBody,
	})
	if err != nil {
		if errors.Is(err, sendmessage.ErrInvalidRequest) {
			return ErrBadRequest("invalid request", err)
		}
		return fmt.Errorf("handle `send message` use case: %v", err)
	}

	return Success(c, SendMessageResponse{
		Data: &MessageWithoutBody{
			AuthorId:  managerID,
			CreatedAt: res.CreatedAt,
			Id:        res.MessageID,
		},
	})
}
