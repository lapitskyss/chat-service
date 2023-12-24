//go:build integration

package chatsrepo_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/suite"

	chatsrepo "github.com/lapitskyss/chat-service/internal/repositories/chats"
	"github.com/lapitskyss/chat-service/internal/testingh"
	"github.com/lapitskyss/chat-service/internal/types"
)

type ChatsRepoSuite struct {
	testingh.DBSuite
	repo *chatsrepo.Repo
}

func TestChatsRepoSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, &ChatsRepoSuite{DBSuite: testingh.NewDBSuite("TestChatsRepoSuite")})
}

func (s *ChatsRepoSuite) SetupSuite() {
	s.DBSuite.SetupSuite()

	var err error

	s.repo, err = chatsrepo.New(chatsrepo.NewOptions(s.Database))
	s.Require().NoError(err)
}

func (s *ChatsRepoSuite) Test_CreateIfNotExists() {
	s.Run("chat does not exist, should be created", func() {
		clientID := types.NewUserID()

		chatID, err := s.repo.CreateIfNotExists(s.Ctx, clientID)
		s.Require().NoError(err)
		s.NotEmpty(chatID)
	})

	s.Run("chat already exists", func() {
		clientID := types.NewUserID()

		// Create chat.
		chat, err := s.Database.Chat(s.Ctx).Create().SetClientID(clientID).Save(s.Ctx)
		s.Require().NoError(err)

		chatID, err := s.repo.CreateIfNotExists(s.Ctx, clientID)
		s.Require().NoError(err)
		s.Require().NotEmpty(chatID)
		s.Equal(chat.ID, chatID)
	})
}

func (s *ChatsRepoSuite) Test_GetChatByID() {
	s.Run("chat exists", func() {
		chatID := types.NewChatID()
		clientID := types.NewUserID()

		// Create chat.
		chat, err := s.Database.Chat(s.Ctx).Create().SetID(chatID).SetClientID(clientID).Save(s.Ctx)
		s.Require().NoError(err)

		// Get it.
		actualMsg, err := s.repo.GetChatByID(s.Ctx, chatID)
		s.Require().NoError(err)
		s.Require().NotNil(actualMsg)
		s.Equal(chat.ID, actualMsg.ID)
		s.Equal(chat.ClientID, actualMsg.ClientID)
	})

	s.Run("chat does not exist", func() {
		msg, err := s.repo.GetChatByID(s.Ctx, types.NewChatID())
		s.Require().Error(err)
		s.Require().Nil(msg)
	})
}

func (s *ChatsRepoSuite) Test_AllWithOpenProblemsForManager() {
	s.Run("has chats with open problems", func() {
		clientID := types.NewUserID()
		managerID := types.NewUserID()

		chatID := s.createChatWithProblemAssignedTo(clientID, managerID)

		chats, err := s.repo.AllWithOpenProblemsForManager(s.Ctx, managerID)
		s.Require().NoError(err)
		s.Len(chats, 1)
		s.Equal(chatID, chats[0].ID)
		s.Equal(clientID, chats[0].ClientID)
	})

	s.Run("has chats with closed problems", func() {
		managerID := types.NewUserID()

		s.createChatWithClosedProblemAssignedTo(types.NewUserID(), types.NewUserID())

		chats, err := s.repo.AllWithOpenProblemsForManager(s.Ctx, managerID)
		s.Require().NoError(err)
		s.Len(chats, 0)
	})

	s.Run("has chats with open problems for another manager", func() {
		managerID := types.NewUserID()

		_ = s.createChatWithProblemAssignedTo(types.NewUserID(), types.NewUserID())
		_ = s.createChatWithProblemAssignedTo(types.NewUserID(), types.NewUserID())

		chats, err := s.repo.AllWithOpenProblemsForManager(s.Ctx, managerID)
		s.Require().NoError(err)
		s.Len(chats, 0)
	})

	s.Run("has chats with closed problems for another manager", func() {
		managerID := types.NewUserID()

		s.createChatWithClosedProblemAssignedTo(types.NewUserID(), types.NewUserID())
		s.createChatWithClosedProblemAssignedTo(types.NewUserID(), types.NewUserID())

		chats, err := s.repo.AllWithOpenProblemsForManager(s.Ctx, managerID)
		s.Require().NoError(err)
		s.Len(chats, 0)
	})

	s.Run("has chats without problems", func() {
		clientID := types.NewUserID()
		managerID := types.NewUserID()

		// Create chat.
		_, err := s.Database.Chat(s.Ctx).Create().SetClientID(clientID).Save(s.Ctx)
		s.Require().NoError(err)

		chats, err := s.repo.AllWithOpenProblemsForManager(s.Ctx, managerID)
		s.Require().NoError(err)
		s.Len(chats, 0)
	})
}

func (s *ChatsRepoSuite) createChatWithProblemAssignedTo(clientID, managerID types.UserID) types.ChatID {
	s.T().Helper()

	// 1 chat can have only 1 open problem.

	chat, err := s.Database.Chat(s.Ctx).Create().SetClientID(clientID).Save(s.Ctx)
	s.Require().NoError(err)

	_, err = s.Database.Problem(s.Ctx).Create().SetChatID(chat.ID).SetManagerID(managerID).Save(s.Ctx)
	s.Require().NoError(err)

	return chat.ID
}

func (s *ChatsRepoSuite) createChatWithClosedProblemAssignedTo(clientID, managerID types.UserID) {
	s.T().Helper()

	// 1 chat can have only 1 open problem.

	chat, err := s.Database.Chat(s.Ctx).Create().SetClientID(clientID).Save(s.Ctx)
	s.Require().NoError(err)

	_, err = s.Database.
		Problem(s.Ctx).
		Create().
		SetChatID(chat.ID).
		SetManagerID(managerID).
		SetResolvedAt(time.Now()).
		Save(s.Ctx)
	s.Require().NoError(err)
}
