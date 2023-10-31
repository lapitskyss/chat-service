package clientv1_test

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	internalerrors "github.com/lapitskyss/chat-service/internal/errors"
	clientv1 "github.com/lapitskyss/chat-service/internal/server-client/v1"
	"github.com/lapitskyss/chat-service/internal/types"
	sendmessage "github.com/lapitskyss/chat-service/internal/usecases/client/send-message"
)

func (s *HandlersSuite) TestSendMessage_BindRequestError() {
	// Arrange.
	reqID := types.NewRequestID()
	resp, eCtx := s.newEchoCtx(reqID, "/v1/sendMessage", `{"messageBody": "Hel`)

	// Action.
	err := s.handlers.PostSendMessage(eCtx, clientv1.PostSendMessageParams{XRequestID: reqID})

	// Assert.
	s.Require().Error(err)
	s.Equal(http.StatusBadRequest, internalerrors.GetServerErrorCode(err))
	s.Empty(resp.Body)
}

func (s *HandlersSuite) TestSendMessage_Usecase_InvalidRequest() {
	// Arrange.
	reqID := types.NewRequestID()
	msgBody := strings.Repeat("!", 3001)

	resp, eCtx := s.newEchoCtx(reqID, "/v1/sendMessage", fmt.Sprintf(`{"messageBody": "%s"}`, msgBody))
	s.sendMsgUseCase.EXPECT().Handle(eCtx.Request().Context(), sendmessage.Request{
		ID:          reqID,
		ClientID:    s.clientID,
		MessageBody: msgBody,
	}).Return(sendmessage.Response{}, sendmessage.ErrInvalidRequest)

	// Action.
	err := s.handlers.PostSendMessage(eCtx, clientv1.PostSendMessageParams{XRequestID: reqID})

	// Assert.
	s.Require().Error(err)
	s.Equal(http.StatusBadRequest, internalerrors.GetServerErrorCode(err))
	s.Empty(resp.Body)
}

func (s *HandlersSuite) TestSendMessage_Usecase_ChatNotCreatedError() {
	// Arrange.
	reqID := types.NewRequestID()
	resp, eCtx := s.newEchoCtx(reqID, "/v1/sendMessage", `{"messageBody": "Hello!"}`)
	s.sendMsgUseCase.EXPECT().Handle(eCtx.Request().Context(), sendmessage.Request{
		ID:          reqID,
		ClientID:    s.clientID,
		MessageBody: "Hello!",
	}).Return(sendmessage.Response{}, sendmessage.ErrChatNotCreated)

	// Action.
	err := s.handlers.PostSendMessage(eCtx, clientv1.PostSendMessageParams{XRequestID: reqID})

	// Assert.
	s.Require().Error(err)
	s.EqualValues(clientv1.ErrorCodeCreateChatError, internalerrors.GetServerErrorCode(err))
	s.Empty(resp.Body)
}

func (s *HandlersSuite) TestSendMessage_Usecase_ProblemNotCreatedError() {
	// Arrange.
	reqID := types.NewRequestID()
	resp, eCtx := s.newEchoCtx(reqID, "/v1/sendMessage", `{"messageBody": "Hello!"}`)
	s.sendMsgUseCase.EXPECT().Handle(eCtx.Request().Context(), sendmessage.Request{
		ID:          reqID,
		ClientID:    s.clientID,
		MessageBody: "Hello!",
	}).Return(sendmessage.Response{}, sendmessage.ErrProblemNotCreated)

	// Action.
	err := s.handlers.PostSendMessage(eCtx, clientv1.PostSendMessageParams{XRequestID: reqID})

	// Assert.
	s.Require().Error(err)
	s.EqualValues(clientv1.ErrorCodeCreateProblemError, internalerrors.GetServerErrorCode(err))
	s.Empty(resp.Body)
}

func (s *HandlersSuite) TestSendMessage_Usecase_Success() {
	// Arrange.
	reqID := types.NewRequestID()
	msgID := types.NewMessageID()

	resp, eCtx := s.newEchoCtx(reqID, "/v1/sendMessage", `{"messageBody": "Hello!"}`)
	s.sendMsgUseCase.EXPECT().Handle(eCtx.Request().Context(), sendmessage.Request{
		ID:          reqID,
		ClientID:    s.clientID,
		MessageBody: "Hello!",
	}).Return(sendmessage.Response{
		AuthorID:  s.clientID,
		MessageID: msgID,
		CreatedAt: time.Unix(1, 1).UTC(),
	}, nil)

	// Action.
	err := s.handlers.PostSendMessage(eCtx, clientv1.PostSendMessageParams{XRequestID: reqID})

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
}`, s.clientID, msgID), resp.Body.String())
}
