package starter

import (
	"context"
	"fmt"

	"github.com/google/wire"
	"go.uber.org/zap"

	"github.com/lapitskyss/chat-service/internal/config"
	"github.com/lapitskyss/chat-service/internal/store"
)

//nolint:unused
var storageSet = wire.NewSet(
	providePSQLClient,
	store.NewDatabase,
)

func providePSQLClient(ctx context.Context, cfg config.Config, log *zap.Logger) (*store.Client, func(), error) {
	storage, err := store.NewPSQLClient(store.NewPSQLOptions(
		cfg.PSQL.Address,
		cfg.PSQL.User,
		cfg.PSQL.Password,
		cfg.PSQL.Database,
		store.WithDebug(cfg.PSQL.Debug),
	))
	if err != nil {
		return nil, nil, fmt.Errorf("create postgres connecton: %v", err)
	}
	if cfg.Global.IsProd() && cfg.PSQL.Debug {
		zap.L().Warn("psql client in the debug mode")
	}

	// Run migration
	err = storage.Schema.Create(ctx)
	if err != nil {
		return nil, nil, fmt.Errorf("auto migration: %v", err)
	}

	cleanup := func() {
		err = storage.Close()
		if err != nil {
			log.Error("close psql client", zap.Error(err))
		}
	}

	return storage, cleanup, nil
}
