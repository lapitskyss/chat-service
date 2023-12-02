package afcverdictsprocessor

import (
	"context"
	"crypto/rsa"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strconv"
	"time"

	"github.com/cenkalti/backoff/v4"
	"github.com/golang-jwt/jwt"
	"github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/protocol"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"

	clientmessageblockedjob "github.com/lapitskyss/chat-service/internal/services/outbox/jobs/client-message-blocked"
	clientmessagesentjob "github.com/lapitskyss/chat-service/internal/services/outbox/jobs/client-message-sent"
	"github.com/lapitskyss/chat-service/internal/types"
)

//go:generate mockgen -source=$GOFILE -destination=mocks/service_mock.gen.go -package=afcverdictsprocessormocks

type messagesRepository interface {
	MarkAsVisibleForManager(ctx context.Context, msgID types.MessageID) error
	BlockMessage(ctx context.Context, msgID types.MessageID) error
}

type outboxService interface {
	Put(ctx context.Context, name, payload string, availableAt time.Time) (types.JobID, error)
}

type transactor interface {
	RunInTx(ctx context.Context, f func(context.Context) error) error
}

//go:generate options-gen -out-filename=service_options.gen.go -from-struct=Options
type Options struct {
	backoffInitialInterval time.Duration `default:"100ms" validate:"min=50ms,max=1s"`
	backoffMaxElapsedTime  time.Duration `default:"5s" validate:"min=500ms,max=1m"`

	brokers          []string `option:"mandatory" validate:"min=1"`
	consumers        int      `option:"mandatory" validate:"min=1,max=16"`
	consumerGroup    string   `option:"mandatory" validate:"required"`
	verdictsTopic    string   `option:"mandatory" validate:"required"`
	verdictsSignKey  string
	processBatchSize int

	readerFactory KafkaReaderFactory `option:"mandatory" validate:"required"`
	dlqWriter     KafkaDLQWriter     `option:"mandatory" validate:"required"`

	txtor   transactor         `option:"mandatory" validate:"required"`
	msgRepo messagesRepository `option:"mandatory" validate:"required"`
	outBox  outboxService      `option:"mandatory" validate:"required"`
}

type Service struct {
	Options

	key     *rsa.PublicKey
	backoff backoff.BackOff
}

func New(opts Options) (*Service, error) {
	if err := opts.Validate(); err != nil {
		return nil, fmt.Errorf("validate options: %v", err)
	}

	var key *rsa.PublicKey
	if opts.verdictsSignKey != "" {
		var err error
		key, err = jwt.ParseRSAPublicKeyFromPEM([]byte(opts.verdictsSignKey))
		if err != nil {
			return nil, fmt.Errorf("parse public key: %v", err)
		}
	}

	retry := backoff.NewExponentialBackOff()
	retry.InitialInterval = 100 * time.Millisecond
	retry.RandomizationFactor = 0
	retry.MaxElapsedTime = 3 * time.Second

	return &Service{
		Options: opts,
		key:     key,
		backoff: retry,
	}, nil
}

func (s *Service) Run(ctx context.Context) error {
	defer s.dlqWriter.Close()

	errGrp, ctx := errgroup.WithContext(ctx)

	for i := 0; i < s.consumers; i++ {
		errGrp.Go(func() error {
			return s.runConsumer(ctx)
		})
	}

	err := errGrp.Wait()
	if err != nil {
		return fmt.Errorf("run consumer: %v", err)
	}
	return nil
}

func (s *Service) runConsumer(ctx context.Context) error {
	reader := s.readerFactory(s.brokers, s.consumerGroup, s.verdictsTopic)

	defer reader.Close()

	for {
		msg, err := reader.FetchMessage(ctx)
		if err != nil {
			if errors.Is(err, context.Canceled) || errors.Is(err, io.EOF) {
				return nil
			}
			return fmt.Errorf("fetch message: %v", err)
		}

		s.handleMessage(ctx, msg)

		err = reader.CommitMessages(ctx, msg)
		if err != nil {
			return fmt.Errorf("commit messages: %v", err)
		}
	}
}

func (s *Service) handleMessage(ctx context.Context, msg kafka.Message) {
	verdict, err := s.parseVerdict(ctx, msg.Value)
	if err != nil {
		s.writeToDLQ(ctx, msg, err.Error())
		return
	}

	if verdict.Status != VerdictStatusOK {
		err = s.processBlockMessage(ctx, verdict)
		if err != nil {
			s.writeToDLQ(ctx, msg, err.Error())
			return
		}
		return
	}

	err = s.processMarkAsVisibleForManager(ctx, verdict)
	if err != nil {
		s.writeToDLQ(ctx, msg, err.Error())
		return
	}
}

func (s *Service) processBlockMessage(ctx context.Context, verdict *Verdict) error {
	return backoff.Retry(func() error {
		err := s.blockMessage(ctx, verdict)
		if err != nil {
			return err
		}
		return nil
	}, backoff.WithContext(s.backoff, ctx))
}

func (s *Service) processMarkAsVisibleForManager(ctx context.Context, verdict *Verdict) error {
	return backoff.Retry(func() error {
		err := s.markAsVisibleForManager(ctx, verdict)
		if err != nil {
			return err
		}
		return nil
	}, backoff.WithContext(s.backoff, ctx))
}

func (s *Service) blockMessage(ctx context.Context, verdict *Verdict) error {
	err := s.txtor.RunInTx(ctx, func(ctx context.Context) error {
		err := s.msgRepo.BlockMessage(ctx, verdict.MessageID)
		if err != nil {
			return fmt.Errorf("block message: %v", err)
		}
		payload, err := clientmessageblockedjob.MarshalPayload(verdict.MessageID)
		if err != nil {
			return fmt.Errorf("marshal message id: %v", err)
		}
		_, err = s.outBox.Put(ctx, clientmessageblockedjob.Name, payload, time.Time{})
		if err != nil {
			return fmt.Errorf("put outbox message: %v", err)
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("run transaction: %v", err)
	}
	return nil
}

func (s *Service) markAsVisibleForManager(ctx context.Context, verdict *Verdict) error {
	err := s.txtor.RunInTx(ctx, func(ctx context.Context) error {
		err := s.msgRepo.MarkAsVisibleForManager(ctx, verdict.MessageID)
		if err != nil {
			return fmt.Errorf("block message: %v", err)
		}
		payload, err := clientmessagesentjob.MarshalPayload(verdict.MessageID)
		if err != nil {
			return fmt.Errorf("marshal message id: %v", err)
		}
		_, err = s.outBox.Put(ctx, clientmessagesentjob.Name, payload, time.Time{})
		if err != nil {
			return fmt.Errorf("put outbox message: %v", err)
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("run transaction: %v", err)
	}
	return nil
}

func (s *Service) parseVerdict(_ context.Context, data []byte) (*Verdict, error) {
	if s.key != nil {
		jwtToken, err := jwt.ParseWithClaims(string(data), &Verdict{}, func(token *jwt.Token) (interface{}, error) {
			_, ok := token.Method.(*jwt.SigningMethodRSA)
			if !ok {
				return nil, fmt.Errorf("incorrect signing method")
			}
			return s.key, nil
		})
		if err != nil {
			return nil, fmt.Errorf("parse token: %v", err)
		}

		payload, ok := jwtToken.Claims.(*Verdict)
		if !ok {
			return nil, fmt.Errorf("incorrect claims type")
		}

		return payload, nil
	}

	var verdict Verdict
	err := json.Unmarshal(data, &verdict)
	if err != nil {
		return nil, fmt.Errorf("unmarshal verdict: %v", err)
	}
	err = verdict.Valid()
	if err != nil {
		return nil, fmt.Errorf("validate verdict: %v", err)
	}
	return &verdict, nil
}

func (s *Service) writeToDLQ(ctx context.Context, msg kafka.Message, errMsg string) {
	lastError := protocol.Header{
		Key:   "LAST_ERROR",
		Value: []byte(errMsg),
	}

	originalPartition := protocol.Header{
		Key:   "ORIGINAL_PARTITION",
		Value: []byte(strconv.Itoa(msg.Partition)),
	}

	msg.Headers = append(msg.Headers, lastError, originalPartition)
	msg.Topic = ""

	err := s.dlqWriter.WriteMessages(ctx, msg)
	if err != nil {
		logError("write message to dlq", err)
	}
}

func logError(msg string, err error) {
	zap.L().Named(serviceName).Error(msg, zap.Error(err))
}
