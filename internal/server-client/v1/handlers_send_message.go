package clientv1

import (
	"errors"
	"fmt"

	"github.com/labstack/echo/v4"

	"github.com/lapitskyss/chat-service/internal/middlewares"
	sendmessage "github.com/lapitskyss/chat-service/internal/usecases/client/send-message"
	"github.com/lapitskyss/chat-service/pkg/pointer"
)

func (h Handlers) PostSendMessage(c echo.Context, params PostSendMessageParams) error {
	ctx := c.Request().Context()
	clientID := middlewares.MustUserID(c)

	var req SendMessageRequest
	if err := c.Bind(&req); err != nil {
		return fmt.Errorf("bind request: %w", err)
	}

	response, err := h.sendMessage.Handle(ctx, sendmessage.Request{
		ID:          params.XRequestID,
		ClientID:    clientID,
		MessageBody: req.MessageBody,
	})
	if err != nil {
		if errors.Is(err, sendmessage.ErrInvalidRequest) {
			return ErrBadRequest("invalid request", err)
		}
		if errors.Is(err, sendmessage.ErrChatNotCreated) {
			return ErrServer(ErrorCodeCreateChatError, "create chat error", err)
		}
		if errors.Is(err, sendmessage.ErrProblemNotCreated) {
			return ErrServer(ErrorCodeCreateProblemError, "create problem error", err)
		}
		return fmt.Errorf("handle `send message` use case: %v", err)
	}

	return Success(c, SendMessageResponse{
		Data: mapPostSendMessageMessageHeader(response),
	})
}

func mapPostSendMessageMessageHeader(response sendmessage.Response) *MessageHeader {
	result := &MessageHeader{
		Id:        response.MessageID,
		CreatedAt: response.CreatedAt,
	}
	if !response.AuthorID.IsZero() {
		result.AuthorId = pointer.Ptr(response.AuthorID)
	}
	return result
}
