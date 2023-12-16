//go:build integration

package problemsrepo_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/suite"

	problemsrepo "github.com/lapitskyss/chat-service/internal/repositories/problems"
	storeproblem "github.com/lapitskyss/chat-service/internal/store/problem"
	"github.com/lapitskyss/chat-service/internal/testingh"
	"github.com/lapitskyss/chat-service/internal/types"
)

type ProblemsRepoSuite struct {
	testingh.DBSuite
	repo *problemsrepo.Repo
}

func TestProblemsRepoSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, &ProblemsRepoSuite{DBSuite: testingh.NewDBSuite("TestProblemsRepoSuite")})
}

func (s *ProblemsRepoSuite) SetupSuite() {
	s.DBSuite.SetupSuite()

	var err error

	s.repo, err = problemsrepo.New(problemsrepo.NewOptions(s.Database))
	s.Require().NoError(err)
}

func (s *ProblemsRepoSuite) Test_CreateIfNotExists() {
	s.Run("problem does not exist, should be created", func() {
		clientID := types.NewUserID()

		// Create chat.
		chat, err := s.Database.Chat(s.Ctx).Create().SetClientID(clientID).Save(s.Ctx)
		s.Require().NoError(err)

		problemID, err := s.repo.CreateIfNotExists(s.Ctx, chat.ID)
		s.Require().NoError(err)
		s.NotEmpty(problemID)

		problem, err := s.Database.Problem(s.Ctx).Get(s.Ctx, problemID)
		s.Require().NoError(err)
		s.Equal(problemID, problem.ID)
		s.Equal(chat.ID, problem.ChatID)
	})

	s.Run("resolved problem already exists, should be created", func() {
		clientID := types.NewUserID()

		// Create chat.
		chat, err := s.Database.Chat(s.Ctx).Create().SetClientID(clientID).Save(s.Ctx)
		s.Require().NoError(err)

		// Create problem.
		problem, err := s.Database.Problem(s.Ctx).Create().
			SetChatID(chat.ID).
			SetManagerID(types.NewUserID()).
			SetResolvedAt(time.Now()).Save(s.Ctx)
		s.Require().NoError(err)

		problemID, err := s.repo.CreateIfNotExists(s.Ctx, chat.ID)
		s.Require().NoError(err)
		s.NotEmpty(problemID)
		s.NotEqual(problem.ID, problemID)
	})

	s.Run("problem already exists", func() {
		clientID := types.NewUserID()

		// Create chat.
		chat, err := s.Database.Chat(s.Ctx).Create().SetClientID(clientID).Save(s.Ctx)
		s.Require().NoError(err)

		// Create problem.
		problem, err := s.Database.Problem(s.Ctx).Create().SetChatID(chat.ID).Save(s.Ctx)
		s.Require().NoError(err)

		problemID, err := s.repo.CreateIfNotExists(s.Ctx, chat.ID)
		s.Require().NoError(err)
		s.NotEmpty(problemID)
		s.Equal(problem.ID, problemID)
	})
}

func (s *ProblemsRepoSuite) Test_GetManagerOpenProblemsCount() {
	s.Run("manager has no open problems", func() {
		managerID := types.NewUserID()

		count, err := s.repo.GetManagerOpenProblemsCount(s.Ctx, managerID)
		s.Require().NoError(err)
		s.Empty(count)
	})

	s.Run("manager has open problems", func() {
		const (
			problemsCount         = 20
			resolvedProblemsCount = 3
		)

		managerID := types.NewUserID()
		problems := make([]types.ProblemID, 0, problemsCount)

		for i := 0; i < problemsCount; i++ {
			pID := s.createChatWithProblemAssignedTo(managerID)
			problems = append(problems, pID)
		}

		// Create problems for other managers.
		for i := 0; i < problemsCount; i++ {
			s.createChatWithProblemAssignedTo(types.NewUserID())
		}

		count, err := s.repo.GetManagerOpenProblemsCount(s.Ctx, managerID)
		s.Require().NoError(err)
		s.Equal(problemsCount, count)

		// Resolve some problems.
		for i := 0; i < resolvedProblemsCount; i++ {
			pID := problems[i*resolvedProblemsCount]
			_, err := s.Database.Problem(s.Ctx).
				Update().
				Where(storeproblem.ID(pID)).
				SetResolvedAt(time.Now()).
				Save(s.Ctx)
			s.Require().NoError(err)
		}

		count, err = s.repo.GetManagerOpenProblemsCount(s.Ctx, managerID)
		s.Require().NoError(err)
		s.Equal(problemsCount-resolvedProblemsCount, count)
	})
}

func (s *ProblemsRepoSuite) Test_GetAssignedProblemID() {
	s.Run("manager has problem for chat", func() {
		clientID := types.NewUserID()
		managerID := types.NewUserID()

		// Create chat.
		chat, err := s.Database.Chat(s.Ctx).Create().SetClientID(clientID).Save(s.Ctx)
		s.Require().NoError(err)

		// Create problem.
		problem, err := s.Database.Problem(s.Ctx).Create().SetChatID(chat.ID).SetManagerID(managerID).Save(s.Ctx)
		s.Require().NoError(err)

		problemID, err := s.repo.GetAssignedProblemID(s.Ctx, managerID, chat.ID)
		s.Require().NoError(err)
		s.NotEmpty(problemID)
		s.Equal(problem.ID, problemID)
	})

	s.Run("manager has problem is resolved", func() {
		clientID := types.NewUserID()
		managerID := types.NewUserID()

		// Create chat.
		chat, err := s.Database.Chat(s.Ctx).Create().SetClientID(clientID).Save(s.Ctx)
		s.Require().NoError(err)

		// Create problem.
		_, err = s.Database.
			Problem(s.Ctx).
			Create().
			SetChatID(chat.ID).
			SetManagerID(managerID).
			SetResolvedAt(time.Now()).
			Save(s.Ctx)
		s.Require().NoError(err)

		problemID, err := s.repo.GetAssignedProblemID(s.Ctx, managerID, chat.ID)
		s.Require().Error(err)
		s.Empty(problemID)
	})
}

func (s *ProblemsRepoSuite) Test_ResolveProblem() {
	// Arrange.
	s.Run("resolve problem", func() {
		// Arrange.
		clientID := types.NewUserID()
		managerID := types.NewUserID()

		chat, err := s.Database.Chat(s.Ctx).Create().SetClientID(clientID).Save(s.Ctx)
		s.Require().NoError(err)
		problem, err := s.Database.Problem(s.Ctx).Create().SetChatID(chat.ID).SetManagerID(managerID).Save(s.Ctx)
		s.Require().NoError(err)

		// Action.
		err = s.repo.ResolveProblem(s.Ctx, problem.ID)
		s.Require().NoError(err)

		// Assert.
		p := s.Database.Problem(s.Ctx).GetX(s.Ctx, problem.ID)
		s.NotEmpty(p.ResolvedAt)
	})
}

func (s *ProblemsRepoSuite) createChatWithProblemAssignedTo(managerID types.UserID) types.ProblemID {
	s.T().Helper()

	// 1 chat can have only 1 open problem.

	chat, err := s.Database.Chat(s.Ctx).Create().SetClientID(types.NewUserID()).Save(s.Ctx)
	s.Require().NoError(err)

	p, err := s.Database.Problem(s.Ctx).Create().SetChatID(chat.ID).SetManagerID(managerID).Save(s.Ctx)
	s.Require().NoError(err)

	return p.ID
}
