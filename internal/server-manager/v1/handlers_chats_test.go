package managerv1_test

import (
	"errors"
	"net/http"

	managerv1 "github.com/lapitskyss/chat-service/internal/server-manager/v1"
	"github.com/lapitskyss/chat-service/internal/types"
	getchats "github.com/lapitskyss/chat-service/internal/usecases/manager/get-chats"
)

func (s *HandlersSuite) TestGetChats_Usecase_Error() {
	// Arrange.
	reqID := types.NewRequestID()
	resp, eCtx := s.newEchoCtx(reqID, "/v1/getChats", "")
	s.getChatsUseCase.EXPECT().Handle(eCtx.Request().Context(), getchats.Request{
		ID:        reqID,
		ManagerID: s.managerID,
	}).Return(getchats.Response{}, errors.New("something went wrong"))

	// Action.
	err := s.handlers.PostGetChats(eCtx, managerv1.PostGetChatsParams{XRequestID: reqID})

	// Assert.
	s.Require().Error(err)
	s.Empty(resp.Body)
}

func (s *HandlersSuite) TestGetChats_Usecase_Success() {
	// Arrange.
	reqID := types.NewRequestID()
	resp, eCtx := s.newEchoCtx(reqID, "/v1/getChats", "")
	s.getChatsUseCase.EXPECT().Handle(eCtx.Request().Context(), getchats.Request{
		ID:        reqID,
		ManagerID: s.managerID,
	}).Return(getchats.Response{Chats: []getchats.Chat{
		{
			ID:       types.MustParse[types.ChatID]("88b5e7a1-cfdd-4823-b694-a971fbf0d289"),
			ClientID: types.MustParse[types.UserID]("bd12c8fc-c9e4-4d41-b533-bd8704307c71"),
		},
	}}, nil)

	// Action.
	err := s.handlers.PostGetChats(eCtx, managerv1.PostGetChatsParams{XRequestID: reqID})

	// Assert.
	s.Require().NoError(err)
	s.Equal(http.StatusOK, resp.Code)
	s.JSONEq(`
{
    "data":
    {
        "chats":
        [
            {
                "chatId": "88b5e7a1-cfdd-4823-b694-a971fbf0d289",
                "clientId": "bd12c8fc-c9e4-4d41-b533-bd8704307c71"
            }
        ]
    }
}`, resp.Body.String())
}
