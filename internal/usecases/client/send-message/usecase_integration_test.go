//go:build integration

package sendmessage_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"

	chatsrepo "github.com/lapitskyss/chat-service/internal/repositories/chats"
	messagesrepo "github.com/lapitskyss/chat-service/internal/repositories/messages"
	problemsrepo "github.com/lapitskyss/chat-service/internal/repositories/problems"
	"github.com/lapitskyss/chat-service/internal/testingh"
	"github.com/lapitskyss/chat-service/internal/types"
	sendmessage "github.com/lapitskyss/chat-service/internal/usecases/client/send-message"
	sendmessagemocks "github.com/lapitskyss/chat-service/internal/usecases/client/send-message/mocks"
)

type UseCaseIntegrationSuite struct {
	testingh.DBSuite

	uCase                sendmessage.UseCase
	uCaseWithMsgRepoMock sendmessage.UseCase

	ctrl        *gomock.Controller
	msgRepoMock *sendmessagemocks.MockmessagesRepository
}

func TestUseCaseIntegrationSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, &UseCaseIntegrationSuite{DBSuite: testingh.NewDBSuite("TestUseCaseIntegrationSuite")})
}

func (s *UseCaseIntegrationSuite) SetupSuite() {
	s.DBSuite.SetupSuite()

	chatRepo, err := chatsrepo.New(chatsrepo.NewOptions(s.Database))
	s.Require().NoError(err)

	msgRepo, err := messagesrepo.New(messagesrepo.NewOptions(s.Database))
	s.Require().NoError(err)

	problemRepo, err := problemsrepo.New(problemsrepo.NewOptions(s.Database))
	s.Require().NoError(err)

	s.uCase, err = sendmessage.New(sendmessage.NewOptions(
		chatRepo,
		msgRepo,
		problemRepo,
		s.Database,
	))
	s.Require().NoError(err)

	s.ctrl = gomock.NewController(s.T())
	s.msgRepoMock = sendmessagemocks.NewMockmessagesRepository(s.ctrl)
	s.uCaseWithMsgRepoMock, err = sendmessage.New(sendmessage.NewOptions(
		chatRepo,
		s.msgRepoMock,
		problemRepo,
		s.Database,
	))
	s.Require().NoError(err)
}

func (s *UseCaseIntegrationSuite) TearDownSuite() {
	s.ctrl.Finish()

	s.DBSuite.TearDownSuite()
}

func (s *UseCaseIntegrationSuite) SetupTest() {
	s.DBSuite.SetupTest()
	s.Database.Message(s.Ctx).Delete().ExecX(s.Ctx)
	s.Database.Problem(s.Ctx).Delete().ExecX(s.Ctx)
	s.Database.Chat(s.Ctx).Delete().ExecX(s.Ctx)
}

func (s *UseCaseIntegrationSuite) TestPositiveScenario() {
	// Arrange.
	clientID := types.NewUserID()
	const messages = 3

	// Action.
	for i := 0; i < messages; i++ {
		resp, err := s.uCase.Handle(s.Ctx, sendmessage.Request{
			ID:          types.NewRequestID(),
			ClientID:    clientID,
			MessageBody: fmt.Sprintf("Message %d", i),
		})
		s.Require().NoError(err)
		s.Require().NotEmpty(resp)
		s.NotEmpty(resp.MessageID)
		s.NotEmpty(resp.CreatedAt)
	}

	// Assert.
	s.Equal(1, s.Database.Chat(s.Ctx).Query().CountX(s.Ctx))
	s.Equal(1, s.Database.Problem(s.Ctx).Query().CountX(s.Ctx))
	s.Equal(messages, s.Database.Message(s.Ctx).Query().CountX(s.Ctx))
}

func (s *UseCaseIntegrationSuite) TestIdempotency() {
	// Arrange.
	reqID := types.NewRequestID()
	clientID := types.NewUserID()
	const messages = 3

	// Action.
	for i := 0; i < messages; i++ {
		resp, err := s.uCase.Handle(s.Ctx, sendmessage.Request{
			ID:          reqID,
			ClientID:    clientID,
			MessageBody: fmt.Sprintf("Message %d", i),
		})
		s.Require().NoError(err)
		s.Require().NotEmpty(resp)
		s.NotEmpty(resp.MessageID)
		s.NotEmpty(resp.CreatedAt)
	}

	// Assert.
	s.Equal(1, s.Database.Chat(s.Ctx).Query().CountX(s.Ctx))
	s.Equal(1, s.Database.Problem(s.Ctx).Query().CountX(s.Ctx))
	s.Equal(1, s.Database.Message(s.Ctx).Query().CountX(s.Ctx))
}

func (s *UseCaseIntegrationSuite) TestAllOrNothing() {
	// Arrange.
	reqID := types.NewRequestID()
	clientID := types.NewUserID()
	const msgBody = "Broken"

	s.msgRepoMock.EXPECT().GetMessageByRequestID(gomock.Any(), reqID).Return(nil, messagesrepo.ErrMsgNotFound)
	s.msgRepoMock.EXPECT().CreateClientVisible(gomock.Any(), reqID, gomock.Any(), gomock.Any(), clientID, msgBody).
		Return(nil, errors.New("unexpected"))

	// Action.
	_, err := s.uCaseWithMsgRepoMock.Handle(s.Ctx, sendmessage.Request{
		ID:          reqID,
		ClientID:    clientID,
		MessageBody: msgBody,
	})

	// Assert.
	s.Require().Error(err)
	s.Equal(0, s.Database.Chat(s.Ctx).Query().CountX(s.Ctx))
	s.Equal(0, s.Database.Problem(s.Ctx).Query().CountX(s.Ctx))
	s.Equal(0, s.Database.Message(s.Ctx).Query().CountX(s.Ctx))
}
