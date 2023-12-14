package getchats_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"

	chatsrepo "github.com/lapitskyss/chat-service/internal/repositories/chats"
	"github.com/lapitskyss/chat-service/internal/testingh"
	"github.com/lapitskyss/chat-service/internal/types"
	getchats "github.com/lapitskyss/chat-service/internal/usecases/manager/get-chats"
	getchatsmocks "github.com/lapitskyss/chat-service/internal/usecases/manager/get-chats/mocks"
)

type UseCaseSuite struct {
	testingh.ContextSuite

	ctrl     *gomock.Controller
	chatRepo *getchatsmocks.MockchatsRepository
	uCase    getchats.UseCase
}

func TestUseCaseSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(UseCaseSuite))
}

func (s *UseCaseSuite) SetupTest() {
	s.ctrl = gomock.NewController(s.T())
	s.chatRepo = getchatsmocks.NewMockchatsRepository(s.ctrl)

	var err error
	s.uCase, err = getchats.New(getchats.NewOptions(s.chatRepo))
	s.Require().NoError(err)

	s.ContextSuite.SetupTest()
}

func (s *UseCaseSuite) TearDownTest() {
	s.ctrl.Finish()

	s.ContextSuite.TearDownTest()
}

func (s *UseCaseSuite) TestRequestValidationError() {
	// Arrange.
	req := getchats.Request{}

	// Action.
	resp, err := s.uCase.Handle(s.Ctx, req)

	// Assert.
	s.Require().Error(err)
	s.ErrorIs(err, getchats.ErrInvalidRequest)
	s.Empty(resp.Chats)
}

func (s *UseCaseSuite) TestGetChatsWithOpenProblemsError() {
	// Arrange.
	reqID := types.NewRequestID()
	managerID := types.NewUserID()

	s.chatRepo.EXPECT().AllWithOpenProblemsForManager(s.Ctx, managerID).
		Return(nil, errors.New("some error"))

	req := getchats.Request{
		ID:        reqID,
		ManagerID: managerID,
	}

	// Action.
	_, err := s.uCase.Handle(s.Ctx, req)

	// Assert.
	s.Require().Error(err)
}

func (s *UseCaseSuite) TestSuccess() {
	// Arrange.
	reqID := types.NewRequestID()
	managerID := types.NewUserID()

	s.chatRepo.EXPECT().AllWithOpenProblemsForManager(s.Ctx, managerID).
		Return([]chatsrepo.Chat{}, nil)

	req := getchats.Request{
		ID:        reqID,
		ManagerID: managerID,
	}

	// Action.
	_, err := s.uCase.Handle(s.Ctx, req)

	// Assert.
	s.Require().NoError(err)
}
