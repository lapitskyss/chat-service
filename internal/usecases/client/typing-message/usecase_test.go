package clienttypingmessage_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"

	problemsrepo "github.com/lapitskyss/chat-service/internal/repositories/problems"
	eventstream "github.com/lapitskyss/chat-service/internal/services/event-stream"
	"github.com/lapitskyss/chat-service/internal/testingh"
	"github.com/lapitskyss/chat-service/internal/types"
	clienttypingmessage "github.com/lapitskyss/chat-service/internal/usecases/client/typing-message"
	clienttypingmessagemocks "github.com/lapitskyss/chat-service/internal/usecases/client/typing-message/mocks"
)

type UseCaseSuite struct {
	testingh.ContextSuite

	ctrl        *gomock.Controller
	problemRepo *clienttypingmessagemocks.MockproblemsRepository
	eventStream *clienttypingmessagemocks.MockeventStream
	uCase       clienttypingmessage.UseCase
}

func TestUseCaseSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(UseCaseSuite))
}

func (s *UseCaseSuite) SetupTest() {
	s.ctrl = gomock.NewController(s.T())
	s.problemRepo = clienttypingmessagemocks.NewMockproblemsRepository(s.ctrl)
	s.eventStream = clienttypingmessagemocks.NewMockeventStream(s.ctrl)

	var err error
	s.uCase, err = clienttypingmessage.New(clienttypingmessage.NewOptions(s.problemRepo, s.eventStream))
	s.Require().NoError(err)

	s.ContextSuite.SetupTest()
}

func (s *UseCaseSuite) TearDownTest() {
	s.ctrl.Finish()

	s.ContextSuite.TearDownTest()
}

func (s *UseCaseSuite) TestRequestValidationError() {
	// Arrange.
	req := clienttypingmessage.Request{}

	// Action.
	err := s.uCase.Handle(s.Ctx, req)

	// Assert.
	s.Require().Error(err)
	s.ErrorIs(err, clienttypingmessage.ErrInvalidRequest)
}

func (s *UseCaseSuite) TestGetClientOpenProblemError() {
	// Arrange.
	reqID := types.NewRequestID()
	clientID := types.NewUserID()

	s.problemRepo.EXPECT().GetClientOpenProblem(gomock.Any(), clientID).
		Return(nil, errors.New("unexpected error"))

	req := clienttypingmessage.Request{
		ID:       reqID,
		ClientID: clientID,
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

	s.problemRepo.EXPECT().GetClientOpenProblem(gomock.Any(), clientID).
		Return(&problemsrepo.Problem{
			ID:        problemID,
			ChatID:    chatID,
			ManagerID: managerID,
		}, nil)
	s.eventStream.EXPECT().Publish(gomock.Any(), managerID, newTypingEventMatcher{
		TypingEvent: &eventstream.TypingEvent{
			EventID:   types.EventIDNil,
			ClientID:  clientID,
			RequestID: reqID,
		},
	}).Return(errors.New("unexpected error"))

	req := clienttypingmessage.Request{
		ID:       reqID,
		ClientID: clientID,
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

	s.problemRepo.EXPECT().GetClientOpenProblem(gomock.Any(), clientID).
		Return(&problemsrepo.Problem{
			ID:        problemID,
			ChatID:    chatID,
			ManagerID: managerID,
		}, nil)
	s.eventStream.EXPECT().Publish(gomock.Any(), managerID, newTypingEventMatcher{
		TypingEvent: &eventstream.TypingEvent{
			EventID:   types.EventIDNil,
			ClientID:  clientID,
			RequestID: reqID,
		},
	}).Return(nil)

	req := clienttypingmessage.Request{
		ID:       reqID,
		ClientID: clientID,
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
