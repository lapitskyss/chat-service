package managerv1_test

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	internalerrors "github.com/lapitskyss/chat-service/internal/errors"
	managerv1 "github.com/lapitskyss/chat-service/internal/server-manager/v1"
	"github.com/lapitskyss/chat-service/internal/types"
	sendmessage "github.com/lapitskyss/chat-service/internal/usecases/manager/send-message"
)

func (s *HandlersSuite) TestSendMessage_BindRequestError() {
	// Arrange.
	reqID := types.NewRequestID()
	resp, eCtx := s.newEchoCtx(reqID, "/v1/sendMessage", `{"messageBody": "Can`)

	// Action.
	err := s.handlers.PostSendMessage(eCtx, managerv1.PostSendMessageParams{XRequestID: reqID})

	// Assert.
	s.Require().Error(err)
	s.Equal(http.StatusBadRequest, internalerrors.GetServerErrorCode(err))
	s.Empty(resp.Body)
}

func (s *HandlersSuite) TestSendMessage_Usecase_UnknownError() {
	// Arrange.
	reqID := types.NewRequestID()
	chatID := types.NewChatID()

	resp, eCtx := s.newEchoCtx(reqID, "/v1/sendMessage",
		fmt.Sprintf(`{"messageBody": "Can I help you?", "chatId": %q}`, chatID))

	s.sendMessageUseCase.EXPECT().Handle(eCtx.Request().Context(), sendmessage.Request{
		ID:          reqID,
		ManagerID:   s.managerID,
		ChatID:      chatID,
		MessageBody: "Can I help you?",
	}).Return(sendmessage.Response{}, errors.New("something went wrong"))

	// Action.
	err := s.handlers.PostSendMessage(eCtx, managerv1.PostSendMessageParams{XRequestID: reqID})

	// Assert.
	s.Require().Error(err)
	s.Empty(resp.Body)
}

func (s *HandlersSuite) TestSendMessage_Usecase_Success() {
	// Arrange.
	reqID := types.NewRequestID()
	chatID := types.NewChatID()

	resp, eCtx := s.newEchoCtx(reqID, "/v1/sendMessage",
		fmt.Sprintf(`{"messageBody": "Can I help you?", "chatId": %q}`, chatID))

	msgID := types.NewMessageID()
	s.sendMessageUseCase.EXPECT().Handle(eCtx.Request().Context(), sendmessage.Request{
		ID:          reqID,
		ManagerID:   s.managerID,
		ChatID:      chatID,
		MessageBody: "Can I help you?",
	}).Return(sendmessage.Response{
		MessageID: msgID,
		CreatedAt: time.Unix(1, 1).UTC(),
	}, nil)

	// Action.
	err := s.handlers.PostSendMessage(eCtx, managerv1.PostSendMessageParams{XRequestID: reqID})

	// Assert.
	s.Require().NoError(err)
	s.Equal(http.StatusOK, resp.Code)
	s.JSONEq(fmt.Sprintf(`
{
    "data":
    {
        "authorId": "%s",
        "createdAt": "1970-01-01T00:00:01.000000001Z",
        "id": "%s"
    }
}`, s.managerID, msgID), resp.Body.String())
}
