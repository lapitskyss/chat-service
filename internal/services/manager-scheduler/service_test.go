//go:build integration

package managerscheduler_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"

	jobsrepo "github.com/lapitskyss/chat-service/internal/repositories/jobs"
	messagesrepo "github.com/lapitskyss/chat-service/internal/repositories/messages"
	problemsrepo "github.com/lapitskyss/chat-service/internal/repositories/problems"
	managerpool "github.com/lapitskyss/chat-service/internal/services/manager-pool"
	inmemmanagerpool "github.com/lapitskyss/chat-service/internal/services/manager-pool/in-mem"
	managerscheduler "github.com/lapitskyss/chat-service/internal/services/manager-scheduler"
	"github.com/lapitskyss/chat-service/internal/services/outbox"
	"github.com/lapitskyss/chat-service/internal/store"
	"github.com/lapitskyss/chat-service/internal/store/problem"
	"github.com/lapitskyss/chat-service/internal/testingh"
	"github.com/lapitskyss/chat-service/internal/types"
)

const period = 100 * time.Millisecond

type ManagerSchedulerSuite struct {
	testingh.DBSuite

	mPool     managerpool.Pool
	scheduler *managerscheduler.Service
}

func TestManagerSchedulerSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, &ManagerSchedulerSuite{DBSuite: testingh.NewDBSuite("TestManagerSchedulerSuite")})
}

func (s *ManagerSchedulerSuite) SetupTest() {
	s.DBSuite.SetupTest()

	msgRepo, err := messagesrepo.New(messagesrepo.NewOptions(s.Database))
	s.Require().NoError(err)

	jobsRepo, err := jobsrepo.New(jobsrepo.NewOptions(s.Database))
	s.Require().NoError(err)

	problemRepo, err := problemsrepo.New(problemsrepo.NewOptions(s.Database))
	s.Require().NoError(err)

	outboxSvc, err := outbox.New(outbox.NewOptions(1, time.Second, time.Second, jobsRepo, s.Database))
	s.Require().NoError(err)

	s.mPool = inmemmanagerpool.New()
	s.scheduler, err = managerscheduler.New(managerscheduler.NewOptions(
		period,
		s.mPool,
		msgRepo,
		outboxSvc,
		problemRepo,
		s.Database,
	))
	s.Require().NoError(err)

	// Garbage collection.
	s.Database.Message(s.Ctx).Delete().ExecX(s.Ctx)
	s.Database.Problem(s.Ctx).Delete().ExecX(s.Ctx)
	s.Database.Chat(s.Ctx).Delete().ExecX(s.Ctx)

	s.Database.Job(s.Ctx).Delete().ExecX(s.Ctx)
	s.Database.FailedJob(s.Ctx).Delete().ExecX(s.Ctx)
}

func (s *ManagerSchedulerSuite) TestScheduling() {
	cancel, errCh := s.runScheduler()
	defer cancel()

	s.createAwaitingManagerProblem()
	s.createAwaitingManagerProblem()

	m1, m2, m3 := types.NewUserID(), types.NewUserID(), types.NewUserID()
	s.Require().NoError(s.mPool.Put(s.Ctx, m3)) // Pool: [m3]

	time.Sleep(period * 2)

	s.createAwaitingManagerProblem()
	s.Require().NoError(s.mPool.Put(s.Ctx, m2)) // Pool: [m2]
	s.Require().NoError(s.mPool.Put(s.Ctx, m1)) // Pool: [m2, m1]
	s.Require().NoError(s.mPool.Put(s.Ctx, m3)) // Pool: [m2, m1, m3]

	time.Sleep(period * 2)
	cancel()
	s.Require().NoError(<-errCh)

	problems := s.Store.Problem.Query().Order(store.Asc(problem.FieldCreatedAt)).AllX(s.Ctx)
	s.Require().Len(problems, 3)

	s.Equal(m3, problems[0].ManagerID)
	s.Equal(m2, problems[1].ManagerID)
	s.Equal(m1, problems[2].ManagerID)
	s.Equal(1, s.mPool.Size()) // Pool: [m3]

	jobsNum := s.Store.Job.Query().CountX(s.Ctx)
	s.Equal(len(problems), jobsNum)
}

func (s *ManagerSchedulerSuite) TestLessManagersThanProblems() {
	const problems = 100
	for i := 0; i < problems; i++ {
		s.createAwaitingManagerProblem()
	}

	s.Require().NoError(s.mPool.Put(s.Ctx, types.NewUserID()))

	for i := 0; i < 3; i++ {
		s.runSchedulerFor(period * 2)

		num := s.Store.Problem.Query().Where(problem.ManagerIDNotNil()).CountX(s.Ctx)
		s.Equal(1, num)
		s.Equal(0, s.mPool.Size())
	}
}

func (s *ManagerSchedulerSuite) TestMoreManagersThanProblems() {
	s.createAwaitingManagerProblem()

	const managers = 100
	for i := 0; i < managers; i++ {
		s.Require().NoError(s.mPool.Put(s.Ctx, types.NewUserID()))
	}

	for i := 0; i < 3; i++ {
		s.runSchedulerFor(period * 2)

		num := s.Store.Problem.Query().Where(problem.ManagerIDNotNil()).CountX(s.Ctx)
		s.Equal(1, num)
		s.Equal(managers-1, s.mPool.Size())
	}
}

func (s *ManagerSchedulerSuite) runSchedulerFor(timeout time.Duration) {
	s.T().Helper()

	cancel, errCh := s.runScheduler()
	defer cancel()

	time.Sleep(timeout)
	cancel()
	s.NoError(<-errCh) // No error expected because of graceful shutdown via cancel ctx.
}

func (s *ManagerSchedulerSuite) runScheduler() (context.CancelFunc, <-chan error) {
	s.T().Helper()

	ctx, cancel := context.WithCancel(s.Ctx)

	errCh := make(chan error, 1)
	go func() { errCh <- s.scheduler.Run(ctx) }()

	return cancel, errCh
}

func (s *ManagerSchedulerSuite) createAwaitingManagerProblem() {
	s.T().Helper()

	clientID := types.NewUserID()
	chat := s.Store.Chat.Create().SetClientID(clientID).SaveX(s.Ctx)
	p := s.Store.Problem.Create().SetChatID(chat.ID).SaveX(s.Ctx)
	s.Database.Message(s.Ctx).Create().
		SetID(types.NewMessageID()).
		SetChatID(chat.ID).
		SetAuthorID(clientID).
		SetProblemID(p.ID).
		SetBody("Где мои деньги?").
		SetIsVisibleForClient(true).
		SetIsVisibleForManager(true).
		SetIsBlocked(false).
		SetIsService(false).
		SetInitialRequestID(types.NewRequestID()).
		SaveX(s.Ctx)

	time.Sleep(10 * time.Millisecond)
}
