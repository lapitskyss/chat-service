package managerv1

import (
	"errors"
	"fmt"

	"github.com/labstack/echo/v4"

	"github.com/lapitskyss/chat-service/internal/middlewares"
	gethistory "github.com/lapitskyss/chat-service/internal/usecases/client/get-history"
	getchats "github.com/lapitskyss/chat-service/internal/usecases/manager/get-chats"
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
