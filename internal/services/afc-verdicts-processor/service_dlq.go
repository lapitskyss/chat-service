package afcverdictsprocessor

import (
	"context"
	"fmt"
	"io"
	"strconv"

	"github.com/segmentio/kafka-go"
	"go.uber.org/multierr"

	"github.com/lapitskyss/chat-service/internal/logger"
)

const serviceName = "afc-verdicts-processor"

type erroredMessage struct {
	msg     kafka.Message
	lastErr error
}

func (s *Service) startDLQProducer(ctx context.Context) (errReturned error) {
	defer multierr.AppendInvoke(&errReturned, multierr.Close(s.dlqWriter))

	for {
		select {
		case <-ctx.Done():
			return nil

		case m, ok := <-s.dlq:
			if !ok {
				return fmt.Errorf("dlq: channel was closed")
			}

			msg := kafka.Message{
				Partition: 0,
				Key:       m.msg.Key,
				Value:     m.msg.Value,
				Headers: append(m.msg.Headers,
					kafka.Header{Key: "LAST_ERROR", Value: []byte(m.lastErr.Error())},
					kafka.Header{Key: "ORIGINAL_PARTITION", Value: []byte(strconv.Itoa(m.msg.Partition))},
				),
			}
			if err := s.dlqWriter.WriteMessages(ctx, msg); err != nil {
				return fmt.Errorf("dql: write msg: %v", err)
			}
		}
	}
}

//go:generate mockgen -source=$GOFILE -destination=mocks/dlq_writer_mock.gen.go -package=afcverdictsprocessormocks

type KafkaDLQWriter interface {
	io.Closer
	WriteMessages(ctx context.Context, msgs ...kafka.Message) error
}

func NewKafkaDLQWriter(brokers []string, topic string) KafkaDLQWriter {
	return &kafka.Writer{
		Addr:         kafka.TCP(brokers...),
		Topic:        topic,
		Balancer:     &kafka.CRC32Balancer{},
		BatchSize:    1,
		RequiredAcks: kafka.RequireOne,
		Async:        false,
		Logger:       logger.NewKafkaAdapted().WithServiceName(serviceName),
		ErrorLogger:  logger.NewKafkaAdapted().WithServiceName(serviceName).ForErrors(),
	}
}
