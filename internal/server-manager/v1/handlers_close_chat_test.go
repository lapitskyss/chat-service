package managerv1_test

import (
	"errors"
	"net/http"

	internalerrors "github.com/lapitskyss/chat-service/internal/errors"
	managerv1 "github.com/lapitskyss/chat-service/internal/server-manager/v1"
	"github.com/lapitskyss/chat-service/internal/types"
	closechat "github.com/lapitskyss/chat-service/internal/usecases/manager/close-chat"
)

func (s *HandlersSuite) TestCloseChat_Usecase_Error() {
	// Arrange.
	reqID := types.NewRequestID()
	chatID := types.MustParse[types.ChatID]("31b4dc06-bc31-11ed-93cc-461e464ebed8")
	resp, eCtx := s.newEchoCtx(reqID, "/v1/closeChat",
		`{"chatId": "31b4dc06-bc31-11ed-93cc-461e464ebed8"}`)
	s.closeChatUseCase.EXPECT().Handle(eCtx.Request().Context(), closechat.Request{
		ID:        reqID,
		ManagerID: s.managerID,
		ChatID:    chatID,
	}).Return(errors.New("something went wrong"))

	// Action.
	err := s.handlers.PostCloseChat(eCtx, managerv1.PostCloseChatParams{XRequestID: reqID})

	// Assert.
	s.Require().Error(err)
	s.Empty(resp.Body)
}

func (s *HandlersSuite) TestCloseChat_Usecase_InvalidRequest() {
	// Arrange.
	reqID := types.NewRequestID()
	chatID := types.MustParse[types.ChatID]("31b4dc06-bc31-11ed-93cc-461e464ebed8")
	resp, eCtx := s.newEchoCtx(reqID, "/v1/closeChat",
		`{"chatId": "31b4dc06-bc31-11ed-93cc-461e464ebed8"}`)
	s.closeChatUseCase.EXPECT().Handle(eCtx.Request().Context(), closechat.Request{
		ID:        reqID,
		ManagerID: s.managerID,
		ChatID:    chatID,
	}).Return(closechat.ErrNoActiveProblemInChat)

	// Action.
	err := s.handlers.PostCloseChat(eCtx, managerv1.PostCloseChatParams{XRequestID: reqID})

	// Assert.
	s.Require().Error(err)
	s.EqualValues(managerv1.ErrorCodeNoActiveProblemInChat, internalerrors.GetServerErrorCode(err))
	s.Empty(resp.Body)
}

func (s *HandlersSuite) TestCloseChat_Usecase_Success() {
	// Arrange.
	reqID := types.NewRequestID()
	chatID := types.MustParse[types.ChatID]("31b4dc06-bc31-11ed-93cc-461e464ebed8")
	resp, eCtx := s.newEchoCtx(reqID, "/v1/closeChat",
		`{"chatId": "31b4dc06-bc31-11ed-93cc-461e464ebed8"}`)
	s.closeChatUseCase.EXPECT().Handle(eCtx.Request().Context(), closechat.Request{
		ID:        reqID,
		ManagerID: s.managerID,
		ChatID:    chatID,
	}).Return(nil)

	// Action.
	err := s.handlers.PostCloseChat(eCtx, managerv1.PostCloseChatParams{XRequestID: reqID})

	// Assert.
	s.Require().NoError(err)
	s.Equal(http.StatusOK, resp.Code)
	s.JSONEq(`
{
    "data": null
}`, resp.Body.String())
}
