package logger_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	"github.com/lapitskyss/chat-service/internal/logger"
)

func TestInit(t *testing.T) {
	err := logger.Init(logger.NewOptions("error", logger.WithProductionMode(true)))
	require.NoError(t, err)

	zap.L().Named("user-cache").Error("inconsistent state", zap.String("uid", "1234"))
	// {"level":"ERROR","T":"2022-10-09T13:56:47.626+0300","component":"user-cache","msg":"inconsistent state","uid":"1234"}
}
