package afcverdictsprocessor

import (
	"context"
	"io"

	"github.com/segmentio/kafka-go"

	"github.com/lapitskyss/chat-service/internal/logger"
)

//go:generate mockgen -source=$GOFILE -destination=mocks/reader_mock.gen.go -package=afcverdictsprocessormocks

type KafkaReaderFactory func(brokers []string, groupID string, topic string) KafkaReader

type KafkaReader interface {
	io.Closer
	FetchMessage(ctx context.Context) (kafka.Message, error)
	CommitMessages(ctx context.Context, msgs ...kafka.Message) error
}

func NewKafkaReader(brokers []string, groupID string, topic string) KafkaReader {
	return kafka.NewReader(kafka.ReaderConfig{
		Brokers:               brokers,
		GroupID:               groupID,
		Topic:                 topic,
		WatchPartitionChanges: true,
		ErrorLogger:           logger.NewKafkaAdapted().WithServiceName(serviceName).ForErrors(),
	})
}
