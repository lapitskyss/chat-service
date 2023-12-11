//go:build integration

package messagesrepo_test

import (
	"testing"

	"github.com/stretchr/testify/suite"

	messagesrepo "github.com/lapitskyss/chat-service/internal/repositories/messages"
	"github.com/lapitskyss/chat-service/internal/testingh"
	"github.com/lapitskyss/chat-service/internal/types"
)

type MsgRepoAntiFraudAPISuite struct {
	testingh.DBSuite
	repo *messagesrepo.Repo
}

func TestMsgRepoAntiFraudAPISuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, &MsgRepoAntiFraudAPISuite{DBSuite: testingh.NewDBSuite("TestMsgRepoAntiFraudAPISuite")})
}

func (s *MsgRepoAntiFraudAPISuite) SetupSuite() {
	s.DBSuite.SetupSuite()

	var err error

	s.repo, err = messagesrepo.New(messagesrepo.NewOptions(s.Database))
	s.Require().NoError(err)
}

func (s *MsgRepoAntiFraudAPISuite) TestMarkAsVisibleForManager() {
	// Arrange.
	msgID := s.createMessage()

	// Action.
	err := s.repo.MarkAsVisibleForManager(s.Ctx, msgID)
	s.Require().NoError(err)

	// Assert.
	msg := s.Database.Message(s.Ctx).GetX(s.Ctx, msgID)
	s.False(msg.IsBlocked)
	s.False(msg.IsService)
	s.False(msg.CheckedAt.IsZero())
	s.True(msg.IsVisibleForClient)
	s.True(msg.IsVisibleForManager)
}

func (s *MsgRepoAntiFraudAPISuite) TestBlockMessage() {
	// Arrange.
	msgID := s.createMessage()

	// Action.
	err := s.repo.BlockMessage(s.Ctx, msgID)
	s.Require().NoError(err)

	// Assert.
	msg := s.Database.Message(s.Ctx).GetX(s.Ctx, msgID)
	s.True(msg.IsBlocked)
	s.False(msg.IsService)
	s.False(msg.CheckedAt.IsZero())
	s.True(msg.IsVisibleForClient)
	s.False(msg.IsVisibleForManager)
}

func (s *MsgRepoAntiFraudAPISuite) createMessage() types.MessageID {
	s.T().Helper()

	authorID := types.NewUserID()
	problemID, chatID := s.createProblemAndChat(authorID)
	msgID := types.NewMessageID()

	_, err := s.Database.Message(s.Ctx).Create().
		SetID(msgID).
		SetChatID(chatID).
		SetAuthorID(authorID).
		SetProblemID(problemID).
		SetBody(msgBody).
		SetIsBlocked(false).
		SetIsVisibleForClient(true).
		SetIsVisibleForManager(false).
		SetIsService(false).
		SetInitialRequestID(types.NewRequestID()).
		Save(s.Ctx)
	s.Require().NoError(err)

	return msgID
}

func (s *MsgRepoAntiFraudAPISuite) createProblemAndChat(clientID types.UserID) (types.ProblemID, types.ChatID) {
	s.T().Helper()

	chat, err := s.Database.Chat(s.Ctx).Create().SetClientID(clientID).Save(s.Ctx)
	s.Require().NoError(err)

	problem, err := s.Database.Problem(s.Ctx).Create().SetChatID(chat.ID).Save(s.Ctx)
	s.Require().NoError(err)

	return problem.ID, chat.ID
}
