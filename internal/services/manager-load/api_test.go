package managerload_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"

	managerload "github.com/lapitskyss/chat-service/internal/services/manager-load"
	managerloadmocks "github.com/lapitskyss/chat-service/internal/services/manager-load/mocks"
	"github.com/lapitskyss/chat-service/internal/testingh"
	"github.com/lapitskyss/chat-service/internal/types"
)

const maxProblemAtTime = 30

type ServiceSuite struct {
	testingh.ContextSuite

	ctrl *gomock.Controller

	problemsRepo *managerloadmocks.MockproblemsRepository
	managerLoad  *managerload.Service
}

func TestServiceSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(ServiceSuite))
}

func (s *ServiceSuite) SetupTest() {
	s.ctrl = gomock.NewController(s.T())
	s.problemsRepo = managerloadmocks.NewMockproblemsRepository(s.ctrl)

	var err error
	s.managerLoad, err = managerload.New(managerload.NewOptions(maxProblemAtTime, s.problemsRepo))
	s.Require().NoError(err)

	s.ContextSuite.SetupTest()
}

func (s *ServiceSuite) TearDownTest() {
	s.ctrl.Finish()

	s.ContextSuite.TearDownTest()
}

func (s *ServiceSuite) TestCanManagerTakeProblem_RepoError() {
	// Arrange.
	managerID := types.NewUserID()
	errExpected := errors.New("any error")

	s.problemsRepo.EXPECT().GetManagerOpenProblemsCount(s.Ctx, managerID).
		Return(0, errExpected)

	// Action.
	resp, err := s.managerLoad.CanManagerTakeProblem(s.Ctx, managerID)

	// Assert.
	s.Require().Error(err)
	s.Require().Equal(false, resp)
}

func (s *ServiceSuite) TestCanManagerTakeProblem_MaxError() {
	// Arrange.
	managerID := types.NewUserID()

	s.problemsRepo.EXPECT().GetManagerOpenProblemsCount(s.Ctx, managerID).
		Return(maxProblemAtTime, nil)

	// Action.
	resp, err := s.managerLoad.CanManagerTakeProblem(s.Ctx, managerID)

	// Assert.
	s.Require().NoError(err)
	s.Require().Equal(false, resp)
}

func (s *ServiceSuite) TestCanManagerTakeProblem_Success() {
	// Arrange.
	managerID := types.NewUserID()

	s.problemsRepo.EXPECT().GetManagerOpenProblemsCount(s.Ctx, managerID).
		Return(maxProblemAtTime-1, nil)

	// Action.
	resp, err := s.managerLoad.CanManagerTakeProblem(s.Ctx, managerID)

	// Assert.
	s.Require().NoError(err)
	s.Require().Equal(true, resp)
}
