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
	sendmanagermessagejob "github.com/lapitskyss/chat-service/internal/services/outbox/jobs/send-manager-message"
	"github.com/lapitskyss/chat-service/internal/testingh"
	"github.com/lapitskyss/chat-service/internal/types"
	sendmessage "github.com/lapitskyss/chat-service/internal/usecases/manager/send-message"
	sendmessagemocks "github.com/lapitskyss/chat-service/internal/usecases/manager/send-message/mocks"
)

type UseCaseSuite struct {
	testingh.ContextSuite

	ctrl        *gomock.Controller
	msgRepo     *sendmessagemocks.MockmessagesRepository
	problemRepo *sendmessagemocks.MockproblemsRepository
	txtor       *sendmessagemocks.Mocktransactor
	outBoxSvc   *sendmessagemocks.MockoutboxService
	uCase       sendmessage.UseCase
}

func TestUseCaseSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(UseCaseSuite))
}

func (s *UseCaseSuite) SetupTest() {
	s.ctrl = gomock.NewController(s.T())
	s.msgRepo = sendmessagemocks.NewMockmessagesRepository(s.ctrl)
	s.outBoxSvc = sendmessagemocks.NewMockoutboxService(s.ctrl)
	s.problemRepo = sendmessagemocks.NewMockproblemsRepository(s.ctrl)
	s.txtor = sendmessagemocks.NewMocktransactor(s.ctrl)

	var err error
	s.uCase, err = sendmessage.New(sendmessage.NewOptions(s.msgRepo, s.outBoxSvc, s.problemRepo, s.txtor))
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

func (s *UseCaseSuite) TestGetAssignedProblemIDError() {
	// Arrange.
	reqID := types.NewRequestID()
	chatID := types.NewChatID()
	managerID := types.NewUserID()

	s.problemRepo.EXPECT().GetAssignedProblemID(gomock.Any(), managerID, chatID).
		Return(types.ProblemIDNil, errors.New("unexpected"))

	req := sendmessage.Request{
		ID:          reqID,
		ManagerID:   managerID,
		ChatID:      chatID,
		MessageBody: "Hello!",
	}

	// Action.
	_, err := s.uCase.Handle(s.Ctx, req)

	// Assert.
	s.Require().Error(err)
}

func (s *UseCaseSuite) TestCreateFullVisibleError() {
	// Arrange.
	reqID := types.NewRequestID()
	chatID := types.NewChatID()
	problemID := types.NewProblemID()
	managerID := types.NewUserID()
	const msgBody = "Hello!"

	s.problemRepo.EXPECT().GetAssignedProblemID(gomock.Any(), managerID, chatID).
		Return(problemID, nil)
	s.txtor.EXPECT().RunInTx(gomock.Any(), gomock.Any()).DoAndReturn(
		func(ctx context.Context, f func(ctx context.Context) error) error {
			return f(ctx)
		})
	s.msgRepo.EXPECT().CreateFullVisible(gomock.Any(), reqID, problemID, chatID, managerID, msgBody).
		Return(nil, errors.New("unexpected"))

	req := sendmessage.Request{
		ID:          reqID,
		ManagerID:   managerID,
		ChatID:      chatID,
		MessageBody: "Hello!",
	}

	// Action.
	_, err := s.uCase.Handle(s.Ctx, req)

	// Assert.
	s.Require().Error(err)
}

func (s *UseCaseSuite) TestPutJobError() {
	// Arrange.
	reqID := types.NewRequestID()
	chatID := types.NewChatID()
	problemID := types.NewProblemID()
	managerID := types.NewUserID()
	const msgBody = "Hello!"

	s.problemRepo.EXPECT().GetAssignedProblemID(gomock.Any(), managerID, chatID).
		Return(problemID, nil)
	s.txtor.EXPECT().RunInTx(gomock.Any(), gomock.Any()).DoAndReturn(
		func(ctx context.Context, f func(ctx context.Context) error) error {
			return f(ctx)
		})
	s.msgRepo.EXPECT().CreateFullVisible(gomock.Any(), reqID, problemID, chatID, managerID, msgBody).
		Return(&messagesrepo.Message{ID: types.NewMessageID()}, nil)
	s.outBoxSvc.EXPECT().Put(gomock.Any(), sendmanagermessagejob.Name, gomock.Any(), gomock.Any()).
		Return(types.JobIDNil, errors.New("unexpected"))

	req := sendmessage.Request{
		ID:          reqID,
		ManagerID:   managerID,
		ChatID:      chatID,
		MessageBody: "Hello!",
	}

	// Action.
	_, err := s.uCase.Handle(s.Ctx, req)

	// Assert.
	s.Require().Error(err)
}

func (s *UseCaseSuite) TestTransactionError() {
	// Arrange.
	reqID := types.NewRequestID()
	chatID := types.NewChatID()
	problemID := types.NewProblemID()
	managerID := types.NewUserID()
	const msgBody = "Hello!"

	s.problemRepo.EXPECT().GetAssignedProblemID(gomock.Any(), managerID, chatID).
		Return(problemID, nil)
	s.txtor.EXPECT().RunInTx(gomock.Any(), gomock.Any()).DoAndReturn(
		func(ctx context.Context, f func(ctx context.Context) error) error {
			_ = f(ctx)
			return sql.ErrTxDone
		})
	s.msgRepo.EXPECT().CreateFullVisible(gomock.Any(), reqID, problemID, chatID, managerID, msgBody).
		Return(&messagesrepo.Message{ID: types.NewMessageID()}, nil)
	s.outBoxSvc.EXPECT().Put(gomock.Any(), sendmanagermessagejob.Name, gomock.Any(), gomock.Any()).
		Return(types.NewJobID(), nil)

	req := sendmessage.Request{
		ID:          reqID,
		ManagerID:   managerID,
		ChatID:      chatID,
		MessageBody: "Hello!",
	}

	// Action.
	res, err := s.uCase.Handle(s.Ctx, req)

	// Assert.
	s.Require().Error(err)
	s.Require().Empty(res.MessageID)
	s.Require().Empty(res.CreatedAt)
}

func (s *UseCaseSuite) TestSuccess() {
	// Arrange.
	reqID := types.NewRequestID()
	chatID := types.NewChatID()
	problemID := types.NewProblemID()
	managerID := types.NewUserID()
	messageID := types.NewMessageID()
	const msgBody = "Hello!"
	createdAt := time.Now()

	s.problemRepo.EXPECT().GetAssignedProblemID(gomock.Any(), managerID, chatID).
		Return(problemID, nil)
	s.txtor.EXPECT().RunInTx(gomock.Any(), gomock.Any()).DoAndReturn(
		func(ctx context.Context, f func(ctx context.Context) error) error {
			return f(ctx)
		})
	s.msgRepo.EXPECT().CreateFullVisible(gomock.Any(), reqID, problemID, chatID, managerID, msgBody).
		Return(&messagesrepo.Message{
			ID:                  messageID,
			ChatID:              chatID,
			AuthorID:            managerID,
			Body:                msgBody,
			CreatedAt:           createdAt,
			IsVisibleForClient:  true,
			IsVisibleForManager: true,
			IsBlocked:           false,
			IsService:           false,
		}, nil)
	s.outBoxSvc.EXPECT().Put(gomock.Any(), sendmanagermessagejob.Name, gomock.Any(), gomock.Any()).
		Return(types.NewJobID(), nil)

	req := sendmessage.Request{
		ID:          reqID,
		ManagerID:   managerID,
		ChatID:      chatID,
		MessageBody: "Hello!",
	}

	// Action.
	resp, err := s.uCase.Handle(s.Ctx, req)

	// Assert.
	s.Require().NoError(err)
	s.Require().Equal(messageID, resp.MessageID)
	s.Require().True(createdAt.Equal(resp.CreatedAt))
}
