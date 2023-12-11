package afcverdictsprocessor

import (
	"context"
	"io"

	"github.com/segmentio/kafka-go"

	"github.com/lapitskyss/chat-service/internal/logger"
)

const serviceName = "afc-verdicts-processor"

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
