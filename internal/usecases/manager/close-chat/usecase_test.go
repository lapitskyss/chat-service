package closechat_test

import (
	"context"
	"database/sql"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"

	messagesrepo "github.com/lapitskyss/chat-service/internal/repositories/messages"
	managerclosechatjob "github.com/lapitskyss/chat-service/internal/services/outbox/jobs/manager-close-chat"
	"github.com/lapitskyss/chat-service/internal/testingh"
	"github.com/lapitskyss/chat-service/internal/types"
	closechat "github.com/lapitskyss/chat-service/internal/usecases/manager/close-chat"
	closechatmocks "github.com/lapitskyss/chat-service/internal/usecases/manager/close-chat/mocks"
)

type UseCaseSuite struct {
	testingh.ContextSuite

	ctrl        *gomock.Controller
	msgRepo     *closechatmocks.MockmessagesRepository
	outBoxSvc   *closechatmocks.MockoutboxService
	problemRepo *closechatmocks.MockproblemsRepository
	txtor       *closechatmocks.Mocktransactor
	uCase       closechat.UseCase
}

func TestUseCaseSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(UseCaseSuite))
}

func (s *UseCaseSuite) SetupTest() {
	s.ctrl = gomock.NewController(s.T())
	s.msgRepo = closechatmocks.NewMockmessagesRepository(s.ctrl)
	s.outBoxSvc = closechatmocks.NewMockoutboxService(s.ctrl)
	s.problemRepo = closechatmocks.NewMockproblemsRepository(s.ctrl)
	s.txtor = closechatmocks.NewMocktransactor(s.ctrl)

	var err error
	s.uCase, err = closechat.New(closechat.NewOptions(s.msgRepo, s.outBoxSvc, s.problemRepo, s.txtor))
	s.Require().NoError(err)

	s.ContextSuite.SetupTest()
}

func (s *UseCaseSuite) TearDownTest() {
	s.ctrl.Finish()

	s.ContextSuite.TearDownTest()
}

func (s *UseCaseSuite) TestRequestValidationError() {
	// Arrange.
	req := closechat.Request{}

	// Action.
	err := s.uCase.Handle(s.Ctx, req)

	// Assert.
	s.Require().Error(err)
	s.ErrorIs(err, closechat.ErrInvalidRequest)
}

func (s *UseCaseSuite) TestGetAssignedProblemIDError() {
	// Arrange.
	reqID := types.NewRequestID()
	chatID := types.NewChatID()
	managerID := types.NewUserID()

	s.txtor.EXPECT().RunInTx(gomock.Any(), gomock.Any()).DoAndReturn(
		func(ctx context.Context, f func(ctx context.Context) error) error {
			return f(ctx)
		})
	s.problemRepo.EXPECT().GetAssignedProblemID(gomock.Any(), managerID, chatID).
		Return(types.ProblemIDNil, errors.New("unexpected"))

	req := closechat.Request{
		ID:        reqID,
		ManagerID: managerID,
		ChatID:    chatID,
	}

	// Action.
	err := s.uCase.Handle(s.Ctx, req)

	// Assert.
	s.Require().Error(err)
}

func (s *UseCaseSuite) TestResolveProblemError() {
	// Arrange.
	reqID := types.NewRequestID()
	chatID := types.NewChatID()
	problemID := types.NewProblemID()
	managerID := types.NewUserID()

	s.txtor.EXPECT().RunInTx(gomock.Any(), gomock.Any()).DoAndReturn(
		func(ctx context.Context, f func(ctx context.Context) error) error {
			return f(ctx)
		})
	s.problemRepo.EXPECT().GetAssignedProblemID(gomock.Any(), managerID, chatID).
		Return(problemID, nil)
	s.problemRepo.EXPECT().ResolveProblem(gomock.Any(), problemID).
		Return(errors.New("unexpected"))

	req := closechat.Request{
		ID:        reqID,
		ManagerID: managerID,
		ChatID:    chatID,
	}

	// Action.
	err := s.uCase.Handle(s.Ctx, req)

	// Assert.
	s.Require().Error(err)
}

func (s *UseCaseSuite) TestCreateServiceMsgError() {
	// Arrange.
	reqID := types.NewRequestID()
	chatID := types.NewChatID()
	problemID := types.NewProblemID()
	managerID := types.NewUserID()
	const msgBody = "Your question has been marked as resolved.\nThank you for being with us!"

	s.txtor.EXPECT().RunInTx(gomock.Any(), gomock.Any()).DoAndReturn(
		func(ctx context.Context, f func(ctx context.Context) error) error {
			return f(ctx)
		})
	s.problemRepo.EXPECT().GetAssignedProblemID(gomock.Any(), managerID, chatID).
		Return(problemID, nil)
	s.problemRepo.EXPECT().ResolveProblem(gomock.Any(), problemID).
		Return(nil)
	s.msgRepo.EXPECT().CreateServiceMsg(gomock.Any(), reqID, problemID, chatID, msgBody, true, false).
		Return(nil, errors.New("unexpected"))

	req := closechat.Request{
		ID:        reqID,
		ManagerID: managerID,
		ChatID:    chatID,
	}

	// Action.
	err := s.uCase.Handle(s.Ctx, req)

	// Assert.
	s.Require().Error(err)
}

func (s *UseCaseSuite) TestPutJobError() {
	// Arrange.
	reqID := types.NewRequestID()
	chatID := types.NewChatID()
	problemID := types.NewProblemID()
	managerID := types.NewUserID()
	messageID := types.NewMessageID()
	const msgBody = "Your question has been marked as resolved.\nThank you for being with us!"
	createdAt := time.Now()

	s.txtor.EXPECT().RunInTx(gomock.Any(), gomock.Any()).DoAndReturn(
		func(ctx context.Context, f func(ctx context.Context) error) error {
			return f(ctx)
		})
	s.problemRepo.EXPECT().GetAssignedProblemID(gomock.Any(), managerID, chatID).
		Return(problemID, nil)
	s.problemRepo.EXPECT().ResolveProblem(gomock.Any(), problemID).
		Return(nil)
	s.msgRepo.EXPECT().CreateServiceMsg(gomock.Any(), reqID, problemID, chatID, msgBody, true, false).
		Return(&messagesrepo.Message{
			ID:                  messageID,
			ChatID:              chatID,
			RequestID:           reqID,
			IsVisibleForClient:  true,
			IsVisibleForManager: true,
			Body:                msgBody,
			CreatedAt:           createdAt,
		}, nil)
	s.outBoxSvc.EXPECT().Put(gomock.Any(), managerclosechatjob.Name, gomock.Any(), gomock.Any()).
		Return(types.JobIDNil, errors.New("unexpected"))

	req := closechat.Request{
		ID:        reqID,
		ManagerID: managerID,
		ChatID:    chatID,
	}

	// Action.
	err := s.uCase.Handle(s.Ctx, req)

	// Assert.
	s.Require().Error(err)
}

func (s *UseCaseSuite) TestTransactionError() {
	// Arrange.
	reqID := types.NewRequestID()
	chatID := types.NewChatID()
	problemID := types.NewProblemID()
	managerID := types.NewUserID()
	messageID := types.NewMessageID()
	const msgBody = "Your question has been marked as resolved.\nThank you for being with us!"
	createdAt := time.Now()

	s.txtor.EXPECT().RunInTx(gomock.Any(), gomock.Any()).DoAndReturn(
		func(ctx context.Context, f func(ctx context.Context) error) error {
			_ = f(ctx)
			return sql.ErrTxDone
		})
	s.problemRepo.EXPECT().GetAssignedProblemID(gomock.Any(), managerID, chatID).
		Return(problemID, nil)
	s.problemRepo.EXPECT().ResolveProblem(gomock.Any(), problemID).
		Return(nil)
	s.msgRepo.EXPECT().CreateServiceMsg(gomock.Any(), reqID, problemID, chatID, msgBody, true, false).
		Return(&messagesrepo.Message{
			ID:                  messageID,
			ChatID:              chatID,
			RequestID:           reqID,
			IsVisibleForClient:  true,
			IsVisibleForManager: true,
			Body:                msgBody,
			CreatedAt:           createdAt,
		}, nil)
	s.outBoxSvc.EXPECT().Put(gomock.Any(), managerclosechatjob.Name, gomock.Any(), gomock.Any()).
		Return(types.NewJobID(), nil)

	req := closechat.Request{
		ID:        reqID,
		ManagerID: managerID,
		ChatID:    chatID,
	}

	// Action.
	err := s.uCase.Handle(s.Ctx, req)

	// Assert.
	s.Require().Error(err)
}

func (s *UseCaseSuite) TestSuccess() {
	// Arrange.
	reqID := types.NewRequestID()
	chatID := types.NewChatID()
	problemID := types.NewProblemID()
	managerID := types.NewUserID()
	messageID := types.NewMessageID()
	const msgBody = "Your question has been marked as resolved.\nThank you for being with us!"
	createdAt := time.Now()

	s.txtor.EXPECT().RunInTx(gomock.Any(), gomock.Any()).DoAndReturn(
		func(ctx context.Context, f func(ctx context.Context) error) error {
			return f(ctx)
		})
	s.problemRepo.EXPECT().GetAssignedProblemID(gomock.Any(), managerID, chatID).
		Return(problemID, nil)
	s.problemRepo.EXPECT().ResolveProblem(gomock.Any(), problemID).
		Return(nil)
	s.msgRepo.EXPECT().CreateServiceMsg(gomock.Any(), reqID, problemID, chatID, msgBody, true, false).
		Return(&messagesrepo.Message{
			ID:                  messageID,
			ChatID:              chatID,
			RequestID:           reqID,
			IsVisibleForClient:  true,
			IsVisibleForManager: true,
			Body:                msgBody,
			CreatedAt:           createdAt,
		}, nil)
	s.outBoxSvc.EXPECT().Put(gomock.Any(), managerclosechatjob.Name, gomock.Any(), gomock.Any()).
		Return(types.NewJobID(), nil)

	req := closechat.Request{
		ID:        reqID,
		ManagerID: managerID,
		ChatID:    chatID,
	}

	// Action.
	err := s.uCase.Handle(s.Ctx, req)

	// Assert.
	s.Require().NoError(err)
}
