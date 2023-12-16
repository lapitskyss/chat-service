package managerv1

import (
	"errors"
	"fmt"

	"github.com/labstack/echo/v4"

	"github.com/lapitskyss/chat-service/internal/middlewares"
	gethistory "github.com/lapitskyss/chat-service/internal/usecases/client/get-history"
	getchathistory "github.com/lapitskyss/chat-service/internal/usecases/manager/get-chat-history"
	getchats "github.com/lapitskyss/chat-service/internal/usecases/manager/get-chats"
	"github.com/lapitskyss/chat-service/pkg/pointer"
)

func (h Handlers) PostGetChats(c echo.Context, params PostGetChatsParams) error {
	ctx := c.Request().Context()
	managerID := middlewares.MustUserID(c)

	response, err := h.getChats.Handle(ctx, getchats.Request{
		ID:        params.XRequestID,
		ManagerID: managerID,
	})
	if err != nil {
		if errors.Is(err, gethistory.ErrInvalidRequest) {
			return ErrBadRequest("invalid request", err)
		}
		return fmt.Errorf("handle `get chats` use case: %v", err)
	}

	return Success(c, GetChatsResponse{
		Data: &ChatList{
			Chats: adaptChats(response.Chats),
		},
	})
}

func (h Handlers) PostGetChatHistory(c echo.Context, params PostGetChatHistoryParams) error {
	ctx := c.Request().Context()
	managerID := middlewares.MustUserID(c)

	var req GetChatHistoryRequest
	if err := c.Bind(&req); err != nil {
		return fmt.Errorf("bind request: %w", err)
	}

	response, err := h.getChatHistory.Handle(ctx, getchathistory.Request{
		ID:        params.XRequestID,
		ManagerID: managerID,
		ChatID:    req.ChatId,
		PageSize:  pointer.Indirect(req.PageSize),
		Cursor:    pointer.Indirect(req.Cursor),
	})
	if err != nil {
		if errors.Is(err, gethistory.ErrInvalidRequest) {
			return ErrBadRequest("invalid request", err)
		}
		if errors.Is(err, gethistory.ErrInvalidCursor) {
			return ErrBadRequest("invalid cursor", err)
		}
		return fmt.Errorf("handle `get chat history` use case: %v", err)
	}

	return Success(c, GetChatHistoryResponse{
		Data: &MessagesPage{
			Messages: adaptGetChatHistoryResponseMessages(response.Messages),
			Next:     response.NextCursor,
		},
	})
}

func adaptChats(cc []getchats.Chat) []Chat {
	result := make([]Chat, len(cc))
	for i, c := range cc {
		result[i] = Chat{
			ChatId:   c.ID,
			ClientId: c.ClientID,
		}
	}
	return result
}

func adaptGetChatHistoryResponseMessages(mm []getchathistory.Message) []Message {
	result := make([]Message, len(mm))
	for i, m := range mm {
		result[i] = Message{
			AuthorId:  m.AuthorID,
			Body:      m.Body,
			CreatedAt: m.CreatedAt,
			Id:        m.ID,
		}
	}
	return result
}
