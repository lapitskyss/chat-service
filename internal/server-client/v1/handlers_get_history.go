package clientv1

import (
	"errors"
	"fmt"

	"github.com/labstack/echo/v4"

	"github.com/lapitskyss/chat-service/internal/middlewares"
	gethistory "github.com/lapitskyss/chat-service/internal/usecases/client/get-history"
	"github.com/lapitskyss/chat-service/pkg/pointer"
)

func (h Handlers) PostGetHistory(c echo.Context, params PostGetHistoryParams) error {
	ctx := c.Request().Context()
	clientID := middlewares.MustUserID(c)

	var req GetHistoryRequest
	if err := c.Bind(&req); err != nil {
		return fmt.Errorf("bind request: %w", err)
	}

	response, err := h.getHistory.Handle(ctx, gethistory.Request{
		ID:       params.XRequestID,
		ClientID: clientID,
		PageSize: pointer.Indirect(req.PageSize),
		Cursor:   pointer.Indirect(req.Cursor),
	})
	if err != nil {
		if errors.Is(err, gethistory.ErrInvalidRequest) {
			return ErrBadRequest(err)
		}
		if errors.Is(err, gethistory.ErrInvalidCursor) {
			return ErrBadRequest(err)
		}
		return err
	}

	return Success(c, GetHistoryResponse{
		Data: &MessagesPage{
			Messages: mapGetHistoryResponseMessages(response.Messages),
			Next:     response.NextCursor,
		},
	})
}

func mapGetHistoryResponseMessages(messages []gethistory.Message) []Message {
	result := make([]Message, len(messages))
	for i, m := range messages {
		message := Message{
			Body:       m.Body,
			CreatedAt:  m.CreatedAt,
			Id:         m.ID,
			IsBlocked:  m.IsBlocked,
			IsReceived: m.IsReceived,
			IsService:  m.IsService,
		}
		if !m.AuthorID.IsZero() {
			message.AuthorId = pointer.Ptr(m.AuthorID)
		}
		result[i] = message
	}
	return result
}
