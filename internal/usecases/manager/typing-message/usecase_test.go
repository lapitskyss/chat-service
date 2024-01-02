package managertypingmessage_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"

	chatsrepo "github.com/lapitskyss/chat-service/internal/repositories/chats"
	problemsrepo "github.com/lapitskyss/chat-service/internal/repositories/problems"
	eventstream "github.com/lapitskyss/chat-service/internal/services/event-stream"
	"github.com/lapitskyss/chat-service/internal/testingh"
	"github.com/lapitskyss/chat-service/internal/types"
	managertypingmessage "github.com/lapitskyss/chat-service/internal/usecases/manager/typing-message"
	managertypingmessagemocks "github.com/lapitskyss/chat-service/internal/usecases/manager/typing-message/mocks"
)

type UseCaseSuite struct {
	testingh.ContextSuite

	ctrl        *gomock.Controller
	chatRepo    *managertypingmessagemocks.MockchatRepository
	problemRepo *managertypingmessagemocks.MockproblemsRepository
	eventStream *managertypingmessagemocks.MockeventStream
	uCase       managertypingmessage.UseCase
}

func TestUseCaseSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(UseCaseSuite))
}

func (s *UseCaseSuite) SetupTest() {
	s.ctrl = gomock.NewController(s.T())
	s.chatRepo = managertypingmessagemocks.NewMockchatRepository(s.ctrl)
	s.problemRepo = managertypingmessagemocks.NewMockproblemsRepository(s.ctrl)
	s.eventStream = managertypingmessagemocks.NewMockeventStream(s.ctrl)

	var err error
	s.uCase, err = managertypingmessage.New(managertypingmessage.NewOptions(s.chatRepo, s.problemRepo, s.eventStream))
	s.Require().NoError(err)

	s.ContextSuite.SetupTest()
}

func (s *UseCaseSuite) TearDownTest() {
	s.ctrl.Finish()

	s.ContextSuite.TearDownTest()
}

func (s *UseCaseSuite) TestRequestValidationError() {
	// Arrange.
	req := managertypingmessage.Request{}

	// Action.
	err := s.uCase.Handle(s.Ctx, req)

	// Assert.
	s.Require().Error(err)
	s.ErrorIs(err, managertypingmessage.ErrInvalidRequest)
}

func (s *UseCaseSuite) TestGetChatOpenProblemError() {
	// Arrange.
	reqID := types.NewRequestID()
	chatID := types.NewChatID()
	managerID := types.NewUserID()

	s.problemRepo.EXPECT().GetChatOpenProblem(gomock.Any(), chatID).
		Return(nil, errors.New("unexpected error"))

	req := managertypingmessage.Request{
		ID:        reqID,
		ManagerID: managerID,
		ChatID:    chatID,
	}

	// Action.
	err := s.uCase.Handle(s.Ctx, req)

	// Assert.
	s.Require().Error(err)
}

func (s *UseCaseSuite) TestChatNotFoundError() {
	// Arrange.
	reqID := types.NewRequestID()
	chatID := types.NewChatID()
	problemID := types.NewProblemID()
	managerID := types.NewUserID()

	s.problemRepo.EXPECT().GetChatOpenProblem(gomock.Any(), chatID).
		Return(&problemsrepo.Problem{
			ID:        problemID,
			ChatID:    chatID,
			ManagerID: types.NewUserID(),
		}, nil)

	req := managertypingmessage.Request{
		ID:        reqID,
		ManagerID: managerID,
		ChatID:    chatID,
	}

	// Action.
	err := s.uCase.Handle(s.Ctx, req)

	// Assert.
	s.Require().Error(err)
	s.Require().ErrorIs(err, managertypingmessage.ErrChatNotFound)
}

func (s *UseCaseSuite) TestGetChatByIDError() {
	// Arrange.
	reqID := types.NewRequestID()
	chatID := types.NewChatID()
	problemID := types.NewProblemID()
	managerID := types.NewUserID()

	s.problemRepo.EXPECT().GetChatOpenProblem(gomock.Any(), chatID).
		Return(&problemsrepo.Problem{
			ID:        problemID,
			ChatID:    chatID,
			ManagerID: managerID,
		}, nil)
	s.chatRepo.EXPECT().GetChatByID(gomock.Any(), chatID).
		Return(nil, errors.New("unexpected error"))

	req := managertypingmessage.Request{
		ID:        reqID,
		ManagerID: managerID,
		ChatID:    chatID,
	}

	// Action.
	err := s.uCase.Handle(s.Ctx, req)

	// Assert.
	s.Require().Error(err)
}

func (s *UseCaseSuite) TestEventStreamPublishError() {
	// Arrange.
	reqID := types.NewRequestID()
	chatID := types.NewChatID()
	problemID := types.NewProblemID()
	managerID := types.NewUserID()
	clientID := types.NewUserID()

	s.problemRepo.EXPECT().GetChatOpenProblem(gomock.Any(), chatID).
		Return(&problemsrepo.Problem{
			ID:        problemID,
			ChatID:    chatID,
			ManagerID: managerID,
		}, nil)
	s.chatRepo.EXPECT().GetChatByID(gomock.Any(), chatID).
		Return(&chatsrepo.Chat{
			ID:       chatID,
			ClientID: clientID,
		}, nil)
	s.eventStream.EXPECT().Publish(gomock.Any(), clientID, newTypingEventMatcher{
		TypingEvent: &eventstream.TypingEvent{
			EventID:   types.EventIDNil,
			ClientID:  managerID,
			RequestID: reqID,
		},
	}).Return(errors.New("unexpected error"))

	req := managertypingmessage.Request{
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
	clientID := types.NewUserID()

	s.problemRepo.EXPECT().GetChatOpenProblem(gomock.Any(), chatID).
		Return(&problemsrepo.Problem{
			ID:        problemID,
			ChatID:    chatID,
			ManagerID: managerID,
		}, nil)
	s.chatRepo.EXPECT().GetChatByID(gomock.Any(), chatID).
		Return(&chatsrepo.Chat{
			ID:       chatID,
			ClientID: clientID,
		}, nil)
	s.eventStream.EXPECT().Publish(gomock.Any(), clientID, newTypingEventMatcher{
		TypingEvent: &eventstream.TypingEvent{
			EventID:   types.EventIDNil,
			ClientID:  managerID,
			RequestID: reqID,
		},
	}).Return(nil)

	req := managertypingmessage.Request{
		ID:        reqID,
		ManagerID: managerID,
		ChatID:    chatID,
	}

	// Action.
	err := s.uCase.Handle(s.Ctx, req)

	// Assert.
	s.Require().NoError(err)
}

var _ gomock.Matcher = newTypingEventMatcher{}

type newTypingEventMatcher struct {
	*eventstream.TypingEvent
}

func (m newTypingEventMatcher) Matches(x interface{}) bool {
	envelope, ok := x.(eventstream.Event)
	if !ok {
		return false
	}

	ev, ok := envelope.(*eventstream.TypingEvent)
	if !ok {
		return false
	}

	return !ev.EventID.IsZero() &&
		ev.RequestID == m.RequestID &&
		ev.ClientID == m.ClientID
}

func (m newTypingEventMatcher) String() string {
	return fmt.Sprintf("%v", m.TypingEvent)
}
