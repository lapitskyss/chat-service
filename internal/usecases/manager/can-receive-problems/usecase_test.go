package canreceiveproblems_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"

	"github.com/lapitskyss/chat-service/internal/testingh"
	"github.com/lapitskyss/chat-service/internal/types"
	canreceiveproblems "github.com/lapitskyss/chat-service/internal/usecases/manager/can-receive-problems"
	canreceiveproblemsmocks "github.com/lapitskyss/chat-service/internal/usecases/manager/can-receive-problems/mocks"
)

type UseCaseSuite struct {
	testingh.ContextSuite

	ctrl      *gomock.Controller
	mLoadMock *canreceiveproblemsmocks.MockmanagerLoadService
	mPoolMock *canreceiveproblemsmocks.MockmanagerPool
	uCase     canreceiveproblems.UseCase
}

func TestUseCaseSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(UseCaseSuite))
}

func (s *UseCaseSuite) SetupTest() {
	s.ctrl = gomock.NewController(s.T())
	s.mLoadMock = canreceiveproblemsmocks.NewMockmanagerLoadService(s.ctrl)
	s.mPoolMock = canreceiveproblemsmocks.NewMockmanagerPool(s.ctrl)

	var err error
	s.uCase, err = canreceiveproblems.New(canreceiveproblems.NewOptions(s.mLoadMock, s.mPoolMock))
	s.Require().NoError(err)

	s.ContextSuite.SetupTest()
}

func (s *UseCaseSuite) TearDownTest() {
	s.ctrl.Finish()
	s.ContextSuite.TearDownTest()
}

func (s *UseCaseSuite) TestRequestValidationError() {
	// Arrange.
	req := canreceiveproblems.Request{}

	// Action.
	_, err := s.uCase.Handle(s.Ctx, req)

	// Assert.
	s.Require().Error(err)
	s.ErrorIs(err, canreceiveproblems.ErrInvalidRequest)
}

func (s *UseCaseSuite) TestManagerPoolContainsError() {
	// Arrange.
	reqID := types.NewRequestID()
	managerID := types.NewUserID()

	s.mPoolMock.EXPECT().Contains(gomock.Any(), managerID).Return(false, errors.New("unexpected"))

	req := canreceiveproblems.Request{
		ID:        reqID,
		ManagerID: managerID,
	}

	// Action.
	_, err := s.uCase.Handle(s.Ctx, req)

	// Assert.
	s.Require().Error(err)
}

func (s *UseCaseSuite) TestManagerAlreadyInPool() {
	// Arrange.
	reqID := types.NewRequestID()
	managerID := types.NewUserID()

	s.mPoolMock.EXPECT().Contains(gomock.Any(), managerID).Return(true, nil)

	req := canreceiveproblems.Request{
		ID:        reqID,
		ManagerID: managerID,
	}

	// Action.
	result, err := s.uCase.Handle(s.Ctx, req)

	// Assert.
	s.Require().NoError(err)
	s.Require().False(result.Result)
}

func (s *UseCaseSuite) TestCanManagerTakeProblemError() {
	// Arrange.
	reqID := types.NewRequestID()
	managerID := types.NewUserID()

	s.mPoolMock.EXPECT().Contains(gomock.Any(), managerID).Return(false, nil)
	s.mLoadMock.EXPECT().CanManagerTakeProblem(gomock.Any(), managerID).
		Return(false, errors.New("unexpected"))

	req := canreceiveproblems.Request{
		ID:        reqID,
		ManagerID: managerID,
	}

	// Action.
	_, err := s.uCase.Handle(s.Ctx, req)

	// Assert.
	s.Require().Error(err)
}

func (s *UseCaseSuite) TestCanManagerTakeProblemNegativeResult() {
	// Arrange.
	reqID := types.NewRequestID()
	managerID := types.NewUserID()

	s.mPoolMock.EXPECT().Contains(gomock.Any(), managerID).Return(false, nil)
	s.mLoadMock.EXPECT().CanManagerTakeProblem(gomock.Any(), managerID).
		Return(false, nil)

	req := canreceiveproblems.Request{
		ID:        reqID,
		ManagerID: managerID,
	}

	// Action.
	result, err := s.uCase.Handle(s.Ctx, req)

	// Assert.
	s.Require().NoError(err)
	s.Require().False(result.Result)
}

func (s *UseCaseSuite) TestCanManagerTakeProblemPositiveResult() {
	// Arrange.
	reqID := types.NewRequestID()
	managerID := types.NewUserID()

	s.mPoolMock.EXPECT().Contains(gomock.Any(), managerID).Return(false, nil)
	s.mLoadMock.EXPECT().CanManagerTakeProblem(gomock.Any(), managerID).
		Return(true, nil)

	req := canreceiveproblems.Request{
		ID:        reqID,
		ManagerID: managerID,
	}

	// Action.
	result, err := s.uCase.Handle(s.Ctx, req)

	// Assert.
	s.Require().NoError(err)
	s.Require().True(result.Result)
}
