//go:build integration

package messagesrepo_test

import (
	"testing"

	"github.com/stretchr/testify/suite"

	messagesrepo "github.com/lapitskyss/chat-service/internal/repositories/messages"
	"github.com/lapitskyss/chat-service/internal/testingh"
	"github.com/lapitskyss/chat-service/internal/types"
)

const (
	msgBody = "whatever"
)

type MsgRepoAPISuite struct {
	testingh.DBSuite
	repo *messagesrepo.Repo
}

func TestMsgRepoAPISuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, &MsgRepoAPISuite{DBSuite: testingh.NewDBSuite("TestMsgRepoAPISuite")})
}

func (s *MsgRepoAPISuite) SetupSuite() {
	s.DBSuite.SetupSuite()

	var err error

	s.repo, err = messagesrepo.New(messagesrepo.NewOptions(s.Database))
	s.Require().NoError(err)
}

func (s *MsgRepoAPISuite) Test_GetMessageByRequestID() {
	s.Run("message exists", func() {
		authorID := types.NewUserID()

		// Create chat and problem.
		problemID, chatID := s.createProblemAndChat(authorID)

		msgID := types.NewMessageID()
		msgInitialRequestID := types.NewRequestID()

		// Create message.
		expectedMsg, err := s.Database.Message(s.Ctx).Create().
			SetID(msgID).
			SetChatID(chatID).
			SetAuthorID(authorID).
			SetProblemID(problemID).
			SetBody(msgBody).
			SetIsBlocked(true).
			SetIsService(true).
			SetInitialRequestID(msgInitialRequestID).
			Save(s.Ctx)
		s.Require().NoError(err)

		// Get it.
		actualMsg, err := s.repo.GetMessageByRequestID(s.Ctx, msgInitialRequestID)
		s.Require().NoError(err)
		s.Require().NotNil(actualMsg)
		s.Equal(expectedMsg.ID, actualMsg.ID)
		s.Equal(expectedMsg.ChatID, actualMsg.ChatID)
		s.Equal(expectedMsg.AuthorID, actualMsg.AuthorID)
		s.Equal(expectedMsg.Body, actualMsg.Body)
		s.Equal(expectedMsg.CreatedAt.Unix(), actualMsg.CreatedAt.Unix())
		s.Equal(expectedMsg.IsVisibleForClient, actualMsg.IsVisibleForClient)
		s.Equal(expectedMsg.IsVisibleForManager, actualMsg.IsVisibleForManager)
		s.Equal(expectedMsg.IsBlocked, actualMsg.IsBlocked)
		s.Equal(expectedMsg.IsService, actualMsg.IsService)
	})

	s.Run("message does not exist", func() {
		msg, err := s.repo.GetMessageByRequestID(s.Ctx, types.NewRequestID())
		s.Require().ErrorIs(err, messagesrepo.ErrMsgNotFound)
		s.Require().Nil(msg)
	})
}

func (s *MsgRepoAPISuite) Test_CreateClientVisible() {
	authorID := types.NewUserID()

	// Create chat and problem.
	problemID, chatID := s.createProblemAndChat(authorID)
	initialRequestID := types.NewRequestID()

	// Check message was created.
	msg, err := s.repo.CreateClientVisible(s.Ctx, initialRequestID, problemID, chatID, authorID, msgBody)
	s.Require().NoError(err)
	s.Require().NotNil(msg)
	s.NotEmpty(msg.ID)
	s.Equal(chatID, msg.ChatID)
	s.Equal(authorID, msg.AuthorID)
	s.Equal(msgBody, msg.Body)
	s.False(msg.CreatedAt.IsZero())
	s.True(msg.IsVisibleForClient)
	s.False(msg.IsVisibleForManager)
	s.False(msg.IsBlocked)
	s.False(msg.IsService)

	{
		dbMsg, err := s.Database.Message(s.Ctx).Get(s.Ctx, msg.ID)
		s.Require().NoError(err)
		s.Require().NotNil(dbMsg)

		s.Run("message is visible for client and invisible for manager", func() {
			s.True(dbMsg.IsVisibleForClient)
			s.False(dbMsg.IsVisibleForManager)
		})

		s.Run("initial_request_id is set correctly", func() {
			s.Equal(initialRequestID, dbMsg.InitialRequestID)
		})
	}
}

func (s *MsgRepoAPISuite) Test_CreateClientVisible_DuplicationError() {
	authorID := types.NewUserID()

	// Create chat and problem.
	problemID, chatID := s.createProblemAndChat(authorID)
	initialRequestID := types.NewRequestID()

	// Check message was created.
	_, err := s.repo.CreateClientVisible(s.Ctx, initialRequestID, problemID, chatID, authorID, msgBody)
	s.Require().NoError(err)

	// Retry message creation.
	_, err = s.repo.CreateClientVisible(s.Ctx, initialRequestID, problemID, chatID, authorID, msgBody)
	s.Require().Error(err)
}

func (s *MsgRepoAPISuite) createProblemAndChat(clientID types.UserID) (types.ProblemID, types.ChatID) {
	s.T().Helper()

	chat, err := s.Database.Chat(s.Ctx).Create().SetClientID(clientID).Save(s.Ctx)
	s.Require().NoError(err)

	problem, err := s.Database.Problem(s.Ctx).Create().SetChatID(chat.ID).Save(s.Ctx)
	s.Require().NoError(err)

	return problem.ID, chat.ID
}
