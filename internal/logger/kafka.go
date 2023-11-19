package logger

import (
	"fmt"

	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var _ kafka.Logger = (*KafkaAdapted)(nil)

type KafkaAdapted struct {
	serviceName string
	lvl         zapcore.Level
}

func NewKafkaAdapted() *KafkaAdapted {
	return &KafkaAdapted{
		serviceName: "",
		lvl:         zapcore.InfoLevel,
	}
}

func (k *KafkaAdapted) Printf(s string, i ...interface{}) {
	zap.L().Named(k.serviceName).Log(k.lvl, fmt.Sprintf(s, i...))
}

func (k *KafkaAdapted) WithServiceName(serviceName string) *KafkaAdapted {
	k.serviceName = serviceName
	return k
}

func (k *KafkaAdapted) ForErrors() *KafkaAdapted {
	k.lvl = zapcore.ErrorLevel
	return k
}
