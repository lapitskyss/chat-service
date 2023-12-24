package afcverdictsprocessor

import (
	"context"
	"crypto/rsa"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"time"

	"github.com/cenkalti/backoff/v4"
	"github.com/golang-jwt/jwt"
	"github.com/segmentio/kafka-go"
	"go.uber.org/multierr"
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
	processBatchSize int `default:"1" validate:"min=1"`

	readerFactory KafkaReaderFactory `option:"mandatory" validate:"required"`
	dlqWriter     KafkaDLQWriter     `option:"mandatory" validate:"required"`

	txtor   transactor         `option:"mandatory" validate:"required"`
	msgRepo messagesRepository `option:"mandatory" validate:"required"`
	outBox  outboxService      `option:"mandatory" validate:"required"`
}

type Service struct {
	Options

	key *rsa.PublicKey
	dlq chan erroredMessage
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

	dlq := make(chan erroredMessage)

	return &Service{
		Options: opts,
		key:     key,
		dlq:     dlq,
	}, nil
}

func (s *Service) Run(ctx context.Context) error {
	eg, ctx := errgroup.WithContext(ctx)

	eg.Go(func() error {
		return s.startDLQProducer(ctx)
	})

	for i := 0; i < s.consumers; i++ {
		eg.Go(func() error {
			return s.runConsumer(ctx)
		})
	}

	err := eg.Wait()
	if err != nil {
		return fmt.Errorf("run consumer: %v", err)
	}
	return nil
}

func (s *Service) runConsumer(ctx context.Context) (errReturned error) {
	consumer := s.readerFactory(s.brokers, s.consumerGroup, s.verdictsTopic)
	defer multierr.AppendInvoke(&errReturned, multierr.Close(consumer))

	var messagesProcessed int

	for {
		msg, err := consumer.FetchMessage(ctx)
		if err != nil {
			if errors.Is(err, context.Canceled) || errors.Is(err, io.EOF) {
				return nil
			}
			return fmt.Errorf("fetch message: %v", err)
		}

		err = s.handleMessage(ctx, msg)
		if err != nil {
			go func(err error) {
				select {
				case <-ctx.Done():
				case s.dlq <- erroredMessage{msg: msg, lastErr: err}:
				}
			}(err)
		}

		messagesProcessed++
		if messagesProcessed == s.processBatchSize {
			err = consumer.CommitMessages(ctx, msg)
			if err != nil {
				return err
			}
			messagesProcessed = 0
		}
	}
}

func (s *Service) handleMessage(ctx context.Context, msg kafka.Message) error {
	verdict, err := s.parseVerdict(ctx, msg.Value)
	if err != nil {
		return err
	}

	if verdict.Status != VerdictStatusOK {
		return s.processBlockMessage(ctx, verdict)
	}

	return s.processMarkAsVisibleForManager(ctx, verdict)
}

func (s *Service) processBlockMessage(ctx context.Context, verdict *Verdict) error {
	b := backoff.WithContext(s.newProcessMessageBackOff(), ctx)

	return backoff.Retry(func() error {
		err := s.blockMessage(ctx, verdict)
		if err != nil {
			return err
		}
		return nil
	}, b)
}

func (s *Service) processMarkAsVisibleForManager(ctx context.Context, verdict *Verdict) error {
	b := backoff.WithContext(s.newProcessMessageBackOff(), ctx)

	return backoff.Retry(func() error {
		err := s.markAsVisibleForManager(ctx, verdict)
		if err != nil {
			return err
		}
		return nil
	}, b)
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

func (s *Service) newProcessMessageBackOff() *backoff.ExponentialBackOff {
	b := backoff.NewExponentialBackOff()
	b.InitialInterval = s.backoffInitialInterval
	b.MaxElapsedTime = s.backoffMaxElapsedTime
	return b
}
