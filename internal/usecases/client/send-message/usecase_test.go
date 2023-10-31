package sendmessage_test

import (
	"context"
	"database/sql"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"

	messagesrepo "github.com/lapitskyss/chat-service/internal/repositories/messages"
	"github.com/lapitskyss/chat-service/internal/testingh"
	"github.com/lapitskyss/chat-service/internal/types"
	sendmessage "github.com/lapitskyss/chat-service/internal/usecases/client/send-message"
	sendmessagemocks "github.com/lapitskyss/chat-service/internal/usecases/client/send-message/mocks"
)

type UseCaseSuite struct {
	testingh.ContextSuite

	ctrl        *gomock.Controller
	chatRepo    *sendmessagemocks.MockchatsRepository
	msgRepo     *sendmessagemocks.MockmessagesRepository
	problemRepo *sendmessagemocks.MockproblemsRepository
	txtor       *sendmessagemocks.Mocktransactor
	uCase       sendmessage.UseCase
}

func TestUseCaseSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(UseCaseSuite))
}

func (s *UseCaseSuite) SetupTest() {
	s.ctrl = gomock.NewController(s.T())
	s.chatRepo = sendmessagemocks.NewMockchatsRepository(s.ctrl)
	s.msgRepo = sendmessagemocks.NewMockmessagesRepository(s.ctrl)
	s.problemRepo = sendmessagemocks.NewMockproblemsRepository(s.ctrl)
	s.txtor = sendmessagemocks.NewMocktransactor(s.ctrl)

	var err error
	s.uCase, err = sendmessage.New(sendmessage.NewOptions(s.chatRepo, s.msgRepo, s.problemRepo, s.txtor))
	s.Require().NoError(err)

	s.ContextSuite.SetupTest()
}

func (s *UseCaseSuite) TearDownTest() {
	s.ctrl.Finish()

	s.ContextSuite.TearDownTest()
}

func (s *UseCaseSuite) TestRequestValidationError() {
	// Arrange.
	req := sendmessage.Request{}

	// Action.
	_, err := s.uCase.Handle(s.Ctx, req)

	// Assert.
	s.Require().Error(err)
	s.ErrorIs(err, sendmessage.ErrInvalidRequest)
}

func (s *UseCaseSuite) TestGetMessageByRequestID_UnexpectedError() {
	// Arrange.
	reqID := types.NewRequestID()
	clientID := types.NewUserID()

	s.txtor.EXPECT().RunInTx(gomock.Any(), gomock.Any()).DoAndReturn(
		func(ctx context.Context, f func(ctx context.Context) error) error {
			return f(ctx)
		})
	s.msgRepo.EXPECT().GetMessageByRequestID(gomock.Any(), reqID).Return(nil, errors.New("unexpected"))

	req := sendmessage.Request{
		ID:          reqID,
		ClientID:    clientID,
		MessageBody: "Hello!",
	}

	// Action.
	_, err := s.uCase.Handle(s.Ctx, req)

	// Assert.
	s.Require().Error(err)
}

func (s *UseCaseSuite) TestGetMessageByRequestID_MsgFound() {
	// Arrange.
	reqID := types.NewRequestID()
	clientID := types.NewUserID()
	const msgBody = "Hello!"
	createdAt := time.Now()
	messageID := types.NewMessageID()

	s.txtor.EXPECT().RunInTx(gomock.Any(), gomock.Any()).DoAndReturn(
		func(ctx context.Context, f func(ctx context.Context) error) error {
			return f(ctx)
		})
	s.msgRepo.EXPECT().GetMessageByRequestID(gomock.Any(), reqID).
		Return(&messagesrepo.Message{
			ID:                  messageID,
			ChatID:              types.NewChatID(),
			AuthorID:            clientID,
			Body:                msgBody,
			CreatedAt:           createdAt,
			IsVisibleForClient:  true,
			IsVisibleForManager: false,
			IsBlocked:           false,
			IsService:           false,
		}, nil)

	req := sendmessage.Request{
		ID:          reqID,
		ClientID:    clientID,
		MessageBody: msgBody,
	}

	// Action.
	resp, err := s.uCase.Handle(s.Ctx, req)

	// Assert.
	s.Require().NoError(err)
	s.Require().Equal(clientID, resp.AuthorID)
	s.Require().Equal(messageID, resp.MessageID)
	s.Require().True(createdAt.Equal(resp.CreatedAt))
}

func (s *UseCaseSuite) TestCreateChatError() {
	// Arrange.
	reqID := types.NewRequestID()
	clientID := types.NewUserID()

	s.txtor.EXPECT().RunInTx(gomock.Any(), gomock.Any()).DoAndReturn(
		func(ctx context.Context, f func(ctx context.Context) error) error {
			return f(ctx)
		})
	s.msgRepo.EXPECT().GetMessageByRequestID(gomock.Any(), reqID).Return(nil, messagesrepo.ErrMsgNotFound)
	s.chatRepo.EXPECT().CreateIfNotExists(gomock.Any(), clientID).Return(types.ChatIDNil, errors.New("unexpected"))

	req := sendmessage.Request{
		ID:          reqID,
		ClientID:    clientID,
		MessageBody: "Hello!",
	}

	// Action.
	_, err := s.uCase.Handle(s.Ctx, req)

	// Assert.
	s.Require().Error(err)
	s.Require().ErrorIs(err, sendmessage.ErrChatNotCreated)
}

func (s *UseCaseSuite) TestCreateProblemError() {
	// Arrange.
	reqID := types.NewRequestID()
	clientID := types.NewUserID()
	chatID := types.NewChatID()

	s.txtor.EXPECT().RunInTx(gomock.Any(), gomock.Any()).DoAndReturn(
		func(ctx context.Context, f func(ctx context.Context) error) error {
			return f(ctx)
		})
	s.msgRepo.EXPECT().GetMessageByRequestID(gomock.Any(), reqID).Return(nil, messagesrepo.ErrMsgNotFound)
	s.chatRepo.EXPECT().CreateIfNotExists(gomock.Any(), clientID).Return(chatID, nil)
	s.problemRepo.EXPECT().CreateIfNotExists(gomock.Any(), chatID).Return(types.ProblemIDNil, errors.New("unexpected"))

	req := sendmessage.Request{
		ID:          reqID,
		ClientID:    clientID,
		MessageBody: "Hello!",
	}

	// Action.
	_, err := s.uCase.Handle(s.Ctx, req)

	// Assert.
	s.Require().Error(err)
	s.Require().ErrorIs(err, sendmessage.ErrProblemNotCreated)
}

func (s *UseCaseSuite) TestCreateMessageError() {
	// Arrange.
	reqID := types.NewRequestID()
	clientID := types.NewUserID()
	chatID := types.NewChatID()
	problemID := types.NewProblemID()
	const msgBody = "Hello!"

	s.txtor.EXPECT().RunInTx(gomock.Any(), gomock.Any()).DoAndReturn(
		func(ctx context.Context, f func(ctx context.Context) error) error {
			return f(ctx)
		})
	s.msgRepo.EXPECT().GetMessageByRequestID(gomock.Any(), reqID).Return(nil, messagesrepo.ErrMsgNotFound)
	s.chatRepo.EXPECT().CreateIfNotExists(gomock.Any(), clientID).Return(chatID, nil)
	s.problemRepo.EXPECT().CreateIfNotExists(gomock.Any(), chatID).Return(problemID, nil)
	s.msgRepo.EXPECT().CreateClientVisible(gomock.Any(), reqID, problemID, chatID, clientID, msgBody).
		Return(nil, errors.New("unexpected"))

	req := sendmessage.Request{
		ID:          reqID,
		ClientID:    clientID,
		MessageBody: msgBody,
	}

	// Action.
	_, err := s.uCase.Handle(s.Ctx, req)

	// Assert.
	s.Require().Error(err)
}

func (s *UseCaseSuite) TestTransactionError() {
	// Arrange.
	reqID := types.NewRequestID()
	clientID := types.NewUserID()
	chatID := types.NewChatID()
	problemID := types.NewProblemID()
	const msgBody = "Hello!"

	s.txtor.EXPECT().RunInTx(gomock.Any(), gomock.Any()).DoAndReturn(
		func(ctx context.Context, f func(ctx context.Context) error) error {
			_ = f(ctx)
			return sql.ErrTxDone
		})
	s.msgRepo.EXPECT().GetMessageByRequestID(gomock.Any(), reqID).Return(nil, messagesrepo.ErrMsgNotFound)
	s.chatRepo.EXPECT().CreateIfNotExists(gomock.Any(), clientID).Return(chatID, nil)
	s.problemRepo.EXPECT().CreateIfNotExists(gomock.Any(), chatID).Return(problemID, nil)
	s.msgRepo.EXPECT().CreateClientVisible(gomock.Any(), reqID, problemID, chatID, clientID, msgBody).
		Return(&messagesrepo.Message{ID: types.NewMessageID()}, nil)

	req := sendmessage.Request{
		ID:          reqID,
		ClientID:    clientID,
		MessageBody: msgBody,
	}

	// Action.
	resp, err := s.uCase.Handle(s.Ctx, req)

	// Assert.
	s.Require().Error(err)
	s.Require().Empty(resp.AuthorID)
	s.Require().Empty(resp.MessageID)
	s.Require().Empty(resp.CreatedAt)
}

func (s *UseCaseSuite) TestNewMsgCreatedSuccessfully() {
	// Arrange.
	reqID := types.NewRequestID()
	clientID := types.NewUserID()
	chatID := types.NewChatID()
	problemID := types.NewProblemID()
	const msgBody = "Hello!"
	createdAt := time.Now()
	messageID := types.NewMessageID()

	s.txtor.EXPECT().RunInTx(gomock.Any(), gomock.Any()).DoAndReturn(
		func(ctx context.Context, f func(ctx context.Context) error) error {
			return f(ctx)
		})
	s.msgRepo.EXPECT().GetMessageByRequestID(gomock.Any(), reqID).Return(nil, messagesrepo.ErrMsgNotFound)
	s.chatRepo.EXPECT().CreateIfNotExists(gomock.Any(), clientID).Return(chatID, nil)
	s.problemRepo.EXPECT().CreateIfNotExists(gomock.Any(), chatID).Return(problemID, nil)
	s.msgRepo.EXPECT().CreateClientVisible(gomock.Any(), reqID, problemID, chatID, clientID, msgBody).
		Return(&messagesrepo.Message{
			ID:                  messageID,
			ChatID:              chatID,
			AuthorID:            clientID,
			Body:                msgBody,
			CreatedAt:           createdAt,
			IsVisibleForClient:  true,
			IsVisibleForManager: false,
			IsBlocked:           false,
			IsService:           false,
		}, nil)

	req := sendmessage.Request{
		ID:          reqID,
		ClientID:    clientID,
		MessageBody: msgBody,
	}

	// Action.
	resp, err := s.uCase.Handle(s.Ctx, req)

	// Assert.
	s.Require().NoError(err)
	s.Require().Equal(clientID, resp.AuthorID)
	s.Require().Equal(messageID, resp.MessageID)
	s.Require().True(createdAt.Equal(resp.CreatedAt))
}
