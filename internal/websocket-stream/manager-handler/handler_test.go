package managerhandler_test

import (
	"bytes"
	"errors"
	"testing"

	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"

	"github.com/lapitskyss/chat-service/internal/testingh"
	"github.com/lapitskyss/chat-service/internal/types"
	managerhandler "github.com/lapitskyss/chat-service/internal/websocket-stream/manager-handler"
	managerhandlermocks "github.com/lapitskyss/chat-service/internal/websocket-stream/manager-handler/mocks"
)

type UseCaseSuite struct {
	testingh.ContextSuite

	ctrl                 *gomock.Controller
	typingMessageUseCase *managerhandlermocks.MockTypingMessageUseCase
	handler              managerhandler.Handler
}

func TestUseCaseSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(UseCaseSuite))
}

func (s *UseCaseSuite) SetupTest() {
	s.ctrl = gomock.NewController(s.T())
	s.typingMessageUseCase = managerhandlermocks.NewMockTypingMessageUseCase(s.ctrl)

	var err error
	s.handler, err = managerhandler.New(managerhandler.NewOptions(
		managerhandler.JSONEventReader{},
		s.typingMessageUseCase,
	))
	s.Require().NoError(err)

	s.ContextSuite.SetupTest()
}

func (s *UseCaseSuite) TearDownTest() {
	s.ctrl.Finish()

	s.ContextSuite.TearDownTest()
}

func (s *UseCaseSuite) TestEventReaderError() {
	// Arrange.
	userID := types.NewUserID()
	buf := bytes.NewBuffer([]byte(`{"eventType"`))

	// Action.
	err := s.handler.Handle(s.Ctx, userID, buf)

	// Assert.
	s.Require().Error(err)
}

func (s *UseCaseSuite) TestValidationError() {
	// Arrange.
	userID := types.NewUserID()
	buf := bytes.NewBuffer([]byte(`
		{
			"eventType": "ManagerTypingEvent",
			"chatId": "",
			"requestId": "d85154f7-867b-44f1-a930-f37151568132"
		}
	`))

	// Action.
	err := s.handler.Handle(s.Ctx, userID, buf)

	// Assert.
	s.Require().Error(err)
}

func (s *UseCaseSuite) TestTypingMessage_Usecase_Error() {
	// Arrange.
	userID := types.NewUserID()
	buf := bytes.NewBuffer([]byte(`
		{
			"eventType": "ManagerTypingEvent",
			"chatId": "c920d118-f3ea-423b-b710-4c4da7610808",
			"requestId": "d85154f7-867b-44f1-a930-f37151568132"
		}
	`))

	s.typingMessageUseCase.EXPECT().Handle(gomock.Any(), gomock.Any()).
		Return(errors.New("unexpected error"))

	// Action.
	err := s.handler.Handle(s.Ctx, userID, buf)

	// Assert.
	s.Require().Error(err)
}

func (s *UseCaseSuite) TestTypingMessage_Usecase_Success() {
	// Arrange.
	userID := types.NewUserID()
	buf := bytes.NewBuffer([]byte(`
		{
			"eventType": "ManagerTypingEvent",
			"chatId": "c920d118-f3ea-423b-b710-4c4da7610808",
			"requestId": "d85154f7-867b-44f1-a930-f37151568132"
		}
	`))

	s.typingMessageUseCase.EXPECT().Handle(gomock.Any(), gomock.Any()).
		Return(nil)

	// Action.
	err := s.handler.Handle(s.Ctx, userID, buf)

	// Assert.
	s.Require().NoError(err)
}
