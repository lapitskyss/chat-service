package starter

import (
	"context"
	"flag"
	"fmt"
	"os/signal"
	"syscall"

	"github.com/google/wire"
	"go.uber.org/zap"

	"github.com/lapitskyss/chat-service/internal/config"
	"github.com/lapitskyss/chat-service/internal/logger"
)

var configPath = flag.String("config", "configs/config.toml", "Path to config file")

//nolint:unused
var infrastructureSet = wire.NewSet(
	provideContext,
	provideConfig,
	provideLogger,
	NewService,
)

func provideContext() (context.Context, func()) {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)

	cleanup := func() {
		cancel()
	}
	return ctx, cleanup
}

func provideConfig() (config.Config, error) {
	cfg, err := config.ParseAndValidate(*configPath)
	if err != nil {
		return config.Config{}, fmt.Errorf("parse and validate config %q: %v", *configPath, err)
	}
	return cfg, nil
}

func provideLogger(cfg config.Config) (*zap.Logger, func(), error) {
	if err := logger.Init(logger.NewOptions(
		cfg.Log.Level,
		logger.WithProductionMode(cfg.Global.IsProd()),
		logger.WithSentryDSN(cfg.Sentry.DSN),
		logger.WithSentryEnv(cfg.Global.Env),
	)); err != nil {
		return nil, nil, fmt.Errorf("init logger: %v", err)
	}

	cleanup := func() {
		logger.Sync()
	}

	return zap.L(), cleanup, nil
}
