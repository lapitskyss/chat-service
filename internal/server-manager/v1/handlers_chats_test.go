package managerv1_test

import (
	"errors"
	"net/http"
	"time"

	managerv1 "github.com/lapitskyss/chat-service/internal/server-manager/v1"
	"github.com/lapitskyss/chat-service/internal/types"
	getchathistory "github.com/lapitskyss/chat-service/internal/usecases/manager/get-chat-history"
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

func (s *HandlersSuite) TestGetChatHistory_Usecase_Error() {
	// Arrange.
	reqID := types.NewRequestID()
	chatID := types.MustParse[types.ChatID]("31b4dc06-bc31-11ed-93cc-461e464ebed8")
	resp, eCtx := s.newEchoCtx(reqID, "/v1/getChatHistory",
		`{"pageSize":10, "chatId": "31b4dc06-bc31-11ed-93cc-461e464ebed8"}`,
	)
	s.getChatHistoryUseCase.EXPECT().Handle(eCtx.Request().Context(), getchathistory.Request{
		ID:        reqID,
		ManagerID: s.managerID,
		ChatID:    chatID,
		PageSize:  10,
	}).Return(getchathistory.Response{}, errors.New("something went wrong"))

	// Action.
	err := s.handlers.PostGetChatHistory(eCtx, managerv1.PostGetChatHistoryParams{XRequestID: reqID})

	// Assert.
	s.Require().Error(err)
	s.Empty(resp.Body)
}

func (s *HandlersSuite) TestGetChatHistory_Usecase_Success() {
	// Arrange.
	reqID := types.NewRequestID()
	chatID := types.MustParse[types.ChatID]("31b4dc06-bc31-11ed-93cc-461e464ebed8")
	resp, eCtx := s.newEchoCtx(reqID, "/v1/getChatHistory",
		`{"pageSize":10, "chatId": "31b4dc06-bc31-11ed-93cc-461e464ebed8"}`,
	)
	s.getChatHistoryUseCase.EXPECT().Handle(eCtx.Request().Context(), getchathistory.Request{
		ID:        reqID,
		ManagerID: s.managerID,
		ChatID:    chatID,
		PageSize:  10,
	}).Return(getchathistory.Response{Messages: []getchathistory.Message{
		{
			ID:        types.MustParse[types.MessageID]("17562ac4-d9b4-492f-b37b-986ce20b7fa7"),
			AuthorID:  types.MustParse[types.UserID]("bd12c8fc-c9e4-4d41-b533-bd8704307c71"),
			Body:      "some message",
			CreatedAt: time.Unix(1, 1).UTC(),
		},
	}}, nil)

	// Action.
	err := s.handlers.PostGetChatHistory(eCtx, managerv1.PostGetChatHistoryParams{XRequestID: reqID})

	// Assert.
	s.Require().NoError(err)
	s.Equal(http.StatusOK, resp.Code)
	s.JSONEq(`
{
    "data":
    {
        "messages":
        [
            {
                "authorId": "bd12c8fc-c9e4-4d41-b533-bd8704307c71",
                "body": "some message",
                "createdAt": "1970-01-01T00:00:01.000000001Z",
                "id": "17562ac4-d9b4-492f-b37b-986ce20b7fa7"
            }
        ],
        "next":""
    }
}`, resp.Body.String())
}
