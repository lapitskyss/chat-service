//go:build integration

package afcverdictsprocessor_test

import (
	"context"
	"crypto/rsa"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/segmentio/kafka-go"
	"github.com/stretchr/testify/suite"

	"github.com/lapitskyss/chat-service/internal/logger"
	jobsrepo "github.com/lapitskyss/chat-service/internal/repositories/jobs"
	messagesrepo "github.com/lapitskyss/chat-service/internal/repositories/messages"
	afcverdictsprocessor "github.com/lapitskyss/chat-service/internal/services/afc-verdicts-processor"
	"github.com/lapitskyss/chat-service/internal/services/outbox"
	"github.com/lapitskyss/chat-service/internal/store/message"
	"github.com/lapitskyss/chat-service/internal/testingh"
	"github.com/lapitskyss/chat-service/internal/types"
)

type ServiceIntegrationSuite struct {
	testingh.DBSuite
	ks testingh.KafkaSuite

	ConsumerGroup  string
	SignPrivateKey string
	SignPubKey     string

	verdictsTopic       string
	verdictsDLQTopic    string
	verdictsProducer    *kafka.Writer
	dlqVerdictsConsumer *kafka.Reader

	signPrivateKey *rsa.PrivateKey
	svc            *afcverdictsprocessor.Service
}

func TestServiceIntegrationSuite_NoSignature(t *testing.T) {
	t.Parallel()
	suite.Run(t, &ServiceIntegrationSuite{
		DBSuite:       testingh.NewDBSuite("TestServiceIntegrationSuite_NoSignature"),
		ConsumerGroup: "TestServiceIntegrationSuite_NoSignature",
	})
}

func TestServiceIntegrationSuite_SignedVerdicts(t *testing.T) {
	t.Parallel()
	suite.Run(t, &ServiceIntegrationSuite{
		DBSuite:        testingh.NewDBSuite("TestServiceIntegrationSuite_SignedVerdicts"),
		ConsumerGroup:  "TestServiceIntegrationSuite_SignedVerdicts",
		SignPrivateKey: privateKey,
		SignPubKey:     publicKey,
	})
}

func (s *ServiceIntegrationSuite) SetT(t *testing.T) { //nolint:thelper
	s.DBSuite.SetT(t)
	s.ks.SetT(t)
}

func (s *ServiceIntegrationSuite) SetupSuite() {
	s.DBSuite.SetupSuite()
	s.ks.SetupSuite()
}

func (s *ServiceIntegrationSuite) TearDownSuite() {
	s.DBSuite.TearDownSuite()
	s.ks.TearDownSuite()
}

func (s *ServiceIntegrationSuite) SetupTest() {
	s.DBSuite.SetupTest()
	s.ks.SetupTest()

	s.verdictsTopic = fmt.Sprintf("%s.%d", "afc.msg-verdicts", time.Now().UnixMilli())
	s.verdictsDLQTopic = fmt.Sprintf("%s.%d", "afc.msg-verdicts.dlq", time.Now().UnixMilli())
	s.ks.RecreateTopics(s.verdictsTopic, s.verdictsDLQTopic)

	s.verdictsProducer = &kafka.Writer{
		Addr:         kafka.TCP(s.ks.KafkaBrokers()...),
		Topic:        s.verdictsTopic,
		Balancer:     &kafka.CRC32Balancer{},
		BatchSize:    1,
		Async:        false,
		RequiredAcks: kafka.RequireOne,
		Logger:       logger.NewKafkaAdapted().WithServiceName(s.ConsumerGroup),
		ErrorLogger:  logger.NewKafkaAdapted().WithServiceName(s.ConsumerGroup).ForErrors(),
	}
	s.dlqVerdictsConsumer = kafka.NewReader(kafka.ReaderConfig{
		Brokers:     s.ks.KafkaBrokers(),
		Topic:       s.verdictsDLQTopic,
		GroupID:     s.ConsumerGroup + ".dlq",
		StartOffset: kafka.FirstOffset,
		Logger:      logger.NewKafkaAdapted().WithServiceName(s.ConsumerGroup + ".dlq"),
		ErrorLogger: logger.NewKafkaAdapted().WithServiceName(s.ConsumerGroup + ".dlq").ForErrors(),
	})

	if k := s.SignPrivateKey; k != "" {
		var err error
		s.signPrivateKey, err = jwt.ParseRSAPrivateKeyFromPEM([]byte(k))
		s.Require().NoError(err)
	}

	msgRepo, err := messagesrepo.New(messagesrepo.NewOptions(s.Database))
	s.Require().NoError(err)

	jobsRepo, err := jobsrepo.New(jobsrepo.NewOptions(s.Database))
	s.Require().NoError(err)

	outboxSvc, err := outbox.New(outbox.NewOptions(1, time.Second, time.Second, jobsRepo, s.Database))
	s.Require().NoError(err)

	s.svc, err = afcverdictsprocessor.New(afcverdictsprocessor.NewOptions(
		s.ks.KafkaBrokers(),
		4,
		s.ConsumerGroup,
		s.verdictsTopic,
		afcverdictsprocessor.NewKafkaReader,
		afcverdictsprocessor.NewKafkaDLQWriter(s.ks.KafkaBrokers(), s.verdictsDLQTopic),
		s.Database,
		msgRepo,
		outboxSvc,
		afcverdictsprocessor.WithVerdictsSignKey(s.SignPubKey),
		afcverdictsprocessor.WithProcessBatchSize(4),
	))
	s.Require().NoError(err)
}

func (s *ServiceIntegrationSuite) TearDownTest() {
	if p := s.verdictsProducer; p != nil {
		s.NoError(p.Close())
	}
	if c := s.dlqVerdictsConsumer; c != nil {
		s.NoError(c.Close())
	}
	s.DBSuite.TearDownTest()
	s.ks.TearDownTest()
}

func (s *ServiceIntegrationSuite) TestComplex() {
	// Arrange.
	const n = 126
	var expPassedMsgs, expBlockedMsg, expBrokenMsgs int

	messages := make([]kafka.Message, n)
	for i := 0; i < n; i++ {
		var status string
		switch {
		case i%2 == 0:
			status = "ok"
			expPassedMsgs++

		case i%3 == 0:
			status = "suspicious"
			expBlockedMsg++

		default:
			status = "abracadabra"
			expBrokenMsgs++
		}

		chat := s.Database.Chat(s.Ctx).Create().SetClientID(types.NewUserID()).SaveX(s.Ctx)
		problem := s.Database.Problem(s.Ctx).Create().SetChatID(chat.ID).SaveX(s.Ctx)

		msg := s.Database.Message(s.Ctx).Create().
			SetChatID(chat.ID).
			SetProblemID(problem.ID).
			SetAuthorID(types.NewUserID()).
			SetIsVisibleForClient(true).
			SetIsVisibleForManager(false).
			SetIsBlocked(false).
			SetInitialRequestID(types.NewRequestID()).
			SetBody(fmt.Sprintf("message %d", i)).
			SaveX(s.Ctx)

		v := verdict{
			ChatID:    msg.ChatID.String(),
			MessageID: msg.ID.String(),
			Status:    status,
		}
		data := s.encode(v)

		messages[i] = kafka.Message{
			Key:   []byte(msg.ChatID.String()),
			Value: []byte(data),
		}
	}

	// Action.
	cancel, errCh := s.runProcessor()
	defer cancel()

	err := s.verdictsProducer.WriteMessages(s.Ctx, messages...)
	s.Require().NoError(err)

	// Assert.
	for i := 0; i < expBrokenMsgs; i++ {
		func() {
			ctx, cancel := context.WithTimeout(s.Ctx, 3*time.Second)
			defer cancel()

			msg, err := s.dlqVerdictsConsumer.ReadMessage(ctx)
			s.Require().NoErrorf(err, "no expected %dth message in dlq topic", i)
			s.assertContainsHeader(msg.Headers, "LAST_ERROR")
			s.assertContainsHeader(msg.Headers, "ORIGINAL_PARTITION")
		}()
	}

	time.Sleep(time.Second) // For the last messages processing.

	passedMsgs := s.Database.Message(s.Ctx).Query().Where(message.IsVisibleForManager(true)).CountX(s.Ctx)
	blockedMsgs := s.Database.Message(s.Ctx).Query().Where(message.IsBlocked(true)).CountX(s.Ctx)
	checkedMsgs := s.Database.Message(s.Ctx).Query().Where(message.CheckedAtNotNil()).CountX(s.Ctx)
	jobs := s.Database.Job(s.Ctx).Query().CountX(s.Ctx)
	failedJobs := s.Database.FailedJob(s.Ctx).Query().CountX(s.Ctx)

	s.Equal(expPassedMsgs, passedMsgs)
	s.Equal(expBlockedMsg, blockedMsgs)
	s.Equal(expPassedMsgs+expBlockedMsg, checkedMsgs)
	s.Equal(expPassedMsgs+expBlockedMsg, jobs)
	s.Equal(0, failedJobs)

	cancel()
	s.Require().NoError(<-errCh)
}

func (s *ServiceIntegrationSuite) runProcessor() (context.CancelFunc, <-chan error) {
	s.T().Helper()

	ctx, cancel := context.WithCancel(s.Ctx)

	errCh := make(chan error)
	go func() { errCh <- s.svc.Run(ctx) }()

	// Waiting for rebalance.
	time.Sleep(2 * time.Second)

	return cancel, errCh
}

func (s *ServiceIntegrationSuite) encode(v verdict) string {
	s.T().Helper()

	if s.signPrivateKey != nil {
		result, err := jwt.NewWithClaims(jwt.SigningMethodRS256, v).SignedString(s.signPrivateKey)
		s.Require().NoError(err)
		return result
	}

	result, err := json.Marshal(v)
	s.Require().NoError(err)
	return string(result)
}

func (s *ServiceIntegrationSuite) assertContainsHeader(headers []kafka.Header, key string) {
	s.T().Helper()

	for _, h := range headers {
		if h.Key == key {
			s.NotEmpty(string(h.Value), "header=%s", key)
			return
		}
	}
	s.Failf("kafka header %s not found", key)
}
