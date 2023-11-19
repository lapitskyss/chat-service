package freehands_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"

	"github.com/lapitskyss/chat-service/internal/testingh"
	"github.com/lapitskyss/chat-service/internal/types"
	freehands "github.com/lapitskyss/chat-service/internal/usecases/manager/free-hands"
	freehandsmocks "github.com/lapitskyss/chat-service/internal/usecases/manager/free-hands/mocks"
)

type UseCaseSuite struct {
	testingh.ContextSuite

	ctrl      *gomock.Controller
	mLoadMock *freehandsmocks.MockmanagerLoadService
	mPoolMock *freehandsmocks.MockmanagerPool
	uCase     freehands.UseCase
}

func TestUseCaseSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(UseCaseSuite))
}

func (s *UseCaseSuite) SetupTest() {
	s.ctrl = gomock.NewController(s.T())
	s.mLoadMock = freehandsmocks.NewMockmanagerLoadService(s.ctrl)
	s.mPoolMock = freehandsmocks.NewMockmanagerPool(s.ctrl)

	var err error
	s.uCase, err = freehands.New(freehands.NewOptions(s.mLoadMock, s.mPoolMock))
	s.Require().NoError(err)

	s.ContextSuite.SetupTest()
}

func (s *UseCaseSuite) TearDownTest() {
	s.ctrl.Finish()
	s.ContextSuite.TearDownTest()
}

func (s *UseCaseSuite) TestRequestValidationError() {
	// Arrange.
	req := freehands.Request{}

	// Action.
	err := s.uCase.Handle(s.Ctx, req)

	// Assert.
	s.Require().Error(err)
	s.ErrorIs(err, freehands.ErrInvalidRequest)
}

func (s *UseCaseSuite) TestCanManagerTakeProblemError() {
	// Arrange.
	reqID := types.NewRequestID()
	managerID := types.NewUserID()

	s.mLoadMock.EXPECT().CanManagerTakeProblem(gomock.Any(), managerID).
		Return(false, errors.New("unexpected"))

	req := freehands.Request{
		ID:        reqID,
		ManagerID: managerID,
	}

	// Action.
	err := s.uCase.Handle(s.Ctx, req)

	// Assert.
	s.Require().Error(err)
}

func (s *UseCaseSuite) TestManagerOverloadedError() {
	// Arrange.
	reqID := types.NewRequestID()
	managerID := types.NewUserID()

	s.mLoadMock.EXPECT().CanManagerTakeProblem(gomock.Any(), managerID).
		Return(false, nil)

	req := freehands.Request{
		ID:        reqID,
		ManagerID: managerID,
	}

	// Action.
	err := s.uCase.Handle(s.Ctx, req)

	// Assert.
	s.Require().Error(err)
	s.ErrorIs(err, freehands.ErrManagerOverloaded)
}

func (s *UseCaseSuite) TestPutInManagerPoolError() {
	// Arrange.
	reqID := types.NewRequestID()
	managerID := types.NewUserID()

	s.mLoadMock.EXPECT().CanManagerTakeProblem(gomock.Any(), managerID).
		Return(true, nil)
	s.mPoolMock.EXPECT().Put(gomock.Any(), managerID).
		Return(errors.New("unexpected"))

	req := freehands.Request{
		ID:        reqID,
		ManagerID: managerID,
	}

	// Action.
	err := s.uCase.Handle(s.Ctx, req)

	// Assert.
	s.Require().Error(err)
}

func (s *UseCaseSuite) TestPutInManagerPoolSuccess() {
	// Arrange.
	reqID := types.NewRequestID()
	managerID := types.NewUserID()

	s.mLoadMock.EXPECT().CanManagerTakeProblem(gomock.Any(), managerID).
		Return(true, nil)
	s.mPoolMock.EXPECT().Put(gomock.Any(), managerID).
		Return(nil)

	req := freehands.Request{
		ID:        reqID,
		ManagerID: managerID,
	}

	// Action.
	err := s.uCase.Handle(s.Ctx, req)

	// Assert.
	s.Require().NoError(err)
}
