//go:build integration

package problemsrepo_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/suite"

	problemsrepo "github.com/lapitskyss/chat-service/internal/repositories/problems"
	"github.com/lapitskyss/chat-service/internal/testingh"
	"github.com/lapitskyss/chat-service/internal/types"
)

type ProblemsRepoScheduleAPISuite struct {
	testingh.DBSuite
	repo *problemsrepo.Repo
}

func TestProblemsRepoScheduleAPISuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, &ProblemsRepoScheduleAPISuite{DBSuite: testingh.NewDBSuite("TestProblemsRepoScheduleAPISuite")})
}

func (s *ProblemsRepoScheduleAPISuite) SetupSuite() {
	s.DBSuite.SetupSuite()

	var err error

	s.repo, err = problemsrepo.New(problemsrepo.NewOptions(s.Database))
	s.Require().NoError(err)
}

func (s *ProblemsRepoScheduleAPISuite) Test_AllAvailableForManager() {
	s.Run("invalid limit", func() {
		for _, l := range []int{-1, 0} {
			problems, err := s.repo.AllAvailableForManager(s.Ctx, l)
			s.Require().Error(err)
			s.Empty(problems)
		}
	})

	s.Run("problems does not exist", func() {
		problems, err := s.repo.AllAvailableForManager(s.Ctx, 3)
		s.Require().NoError(err)
		s.Require().Len(problems, 0)
	})

	s.Run("no open problems without manager", func() {
		clientID := types.NewUserID()
		managerID := types.NewUserID()

		// Create chat.
		chat, err := s.Database.Chat(s.Ctx).Create().SetClientID(clientID).Save(s.Ctx)
		s.Require().NoError(err)

		// Assign open problem with manager to chat.
		_, err = s.Database.Problem(s.Ctx).Create().SetChatID(chat.ID).SetManagerID(managerID).Save(s.Ctx)
		s.Require().NoError(err)

		problems, err := s.repo.AllAvailableForManager(s.Ctx, 3)
		s.Require().NoError(err)
		s.Empty(problems)
	})

	s.Run("no open problems with messages visible for manager", func() {
		clientID := types.NewUserID()
		managerID := types.NewUserID()

		// Create chat.
		chat, err := s.Database.Chat(s.Ctx).Create().SetClientID(clientID).Save(s.Ctx)
		s.Require().NoError(err)

		// Problem without manager.
		_, err = s.Database.Problem(s.Ctx).Create().SetChatID(chat.ID).Save(s.Ctx)
		s.Require().NoError(err)

		// Problem without manager-visible messages.
		p, err := s.Database.Problem(s.Ctx).Create().SetChatID(chat.ID).SetManagerID(managerID).Save(s.Ctx)
		s.Require().NoError(err)

		for i := 0; i < 3; i++ {
			_, err = s.Database.Message(s.Ctx).Create().
				SetID(types.NewMessageID()).
				SetChatID(chat.ID).
				SetAuthorID(clientID).
				SetProblemID(p.ID).
				SetBody("SMS code is 4321").
				SetIsVisibleForClient(true).
				SetIsVisibleForManager(false).
				SetIsBlocked(true).
				SetIsService(false).
				SetInitialRequestID(types.NewRequestID()).
				Save(s.Ctx)
			s.Require().NoError(err)
		}

		problems, err := s.repo.AllAvailableForManager(s.Ctx, 3)
		s.Require().NoError(err)
		s.Empty(problems)
	})

	s.Run("problems exists", func() {
		// Create chat and problem with message available for manager
		chatID1, problemID1, _ := s.createChatWithProblemWithMessage()

		time.Sleep(time.Millisecond)

		chatID2, problemID2, _ := s.createChatWithProblemWithMessage()

		// Get it.
		availableProblems, err := s.repo.AllAvailableForManager(s.Ctx, 3)
		s.Require().NoError(err)
		s.Require().NotNil(availableProblems)
		s.Require().Len(availableProblems, 2)
		s.Equal(problemID1, availableProblems[0].ID)
		s.Equal(chatID1, availableProblems[0].ChatID)
		s.Equal(problemID2, availableProblems[1].ID)
		s.Equal(chatID2, availableProblems[1].ChatID)
	})
}

func (s *ProblemsRepoScheduleAPISuite) Test_SetManager() {
	// Arrange.
	managerID := types.NewUserID()
	_, problemID, _ := s.createChatWithProblemWithMessage()

	// Action.
	err := s.repo.SetManager(s.Ctx, problemID, managerID)
	s.Require().NoError(err)

	// Assert.
	problem := s.Database.Problem(s.Ctx).GetX(s.Ctx, problemID)
	s.Equal(problemID, problem.ID)
	s.Equal(managerID, problem.ManagerID)
}

func (s *ProblemsRepoScheduleAPISuite) Test_GetProblemRequestID() {
	// Arrange.
	_, _, _ = s.createChatWithProblemWithMessage()
	_, problemID, messageID := s.createChatWithProblemWithMessage()

	// Action.
	requestID, err := s.repo.GetProblemRequestID(s.Ctx, problemID)
	s.Require().NoError(err)

	// Assert.
	message := s.Database.Message(s.Ctx).GetX(s.Ctx, messageID)
	s.Equal(message.InitialRequestID, requestID)
}

func (s *ProblemsRepoScheduleAPISuite) createChatWithProblemWithMessage() (types.ChatID, types.ProblemID, types.MessageID) {
	s.T().Helper()

	clientID := types.NewUserID()

	chat, err := s.Database.Chat(s.Ctx).Create().SetClientID(clientID).Save(s.Ctx)
	s.Require().NoError(err)

	p, err := s.Database.Problem(s.Ctx).Create().SetChatID(chat.ID).Save(s.Ctx)
	s.Require().NoError(err)

	msg, err := s.Database.Message(s.Ctx).Create().
		SetID(types.NewMessageID()).
		SetChatID(chat.ID).
		SetAuthorID(clientID).
		SetProblemID(p.ID).
		SetBody("some body").
		SetIsBlocked(false).
		SetIsVisibleForManager(true).
		SetIsVisibleForClient(true).
		SetIsService(false).
		SetInitialRequestID(types.NewRequestID()).
		Save(s.Ctx)
	s.Require().NoError(err)

	return chat.ID, p.ID, msg.ID
}
