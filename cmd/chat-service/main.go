package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"os/signal"
	"syscall"

	"golang.org/x/sync/errgroup"

	"github.com/lapitskyss/chat-service/internal/config"
	"github.com/lapitskyss/chat-service/internal/logger"
	clientv1 "github.com/lapitskyss/chat-service/internal/server-client/v1"
	serverdebug "github.com/lapitskyss/chat-service/internal/server-debug"
)

var configPath = flag.String("config", "configs/config.toml", "Path to config file")

func main() {
	if err := run(); err != nil {
		log.Fatalf("run app: %v", err)
	}
}

func run() (errReturned error) {
	flag.Parse()

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	cfg, err := config.ParseAndValidate(*configPath)
	if err != nil {
		return fmt.Errorf("parse and validate config %q: %v", *configPath, err)
	}

	err = logger.Init(logger.NewOptions(
		cfg.Log.Level,
		logger.WithSentryDNS(cfg.Sentry.DSN),
		logger.WithProductionMode(cfg.Global.IsProd()),
	))
	if err != nil {
		return fmt.Errorf("init logger: %v", err)
	}
	defer logger.Sync()

	srvDebug, err := serverdebug.New(serverdebug.NewOptions(cfg.Servers.Debug.Addr))
	if err != nil {
		return fmt.Errorf("init debug server: %v", err)
	}

	v1Swagger, err := clientv1.GetSwagger()
	if err != nil {
		return fmt.Errorf("init client swagger: %v", err)
	}
	svrClient, err := initServerClient(
		cfg.Global.IsProd(),
		cfg.Servers.Client.Addr,
		cfg.Servers.Client.AllowOrigins,
		cfg.Clients.Keycloak.BasePath,
		cfg.Clients.Keycloak.Realm,
		cfg.Clients.Keycloak.ClientID,
		cfg.Clients.Keycloak.ClientSecret,
		cfg.Clients.Keycloak.DebugMode,
		cfg.Servers.Client.RequiredAccess.Resource,
		cfg.Servers.Client.RequiredAccess.Role,
		v1Swagger,
	)
	if err != nil {
		return fmt.Errorf("init client server: %v", err)
	}

	eg, ctx := errgroup.WithContext(ctx)

	// Run servers.
	eg.Go(func() error { return srvDebug.Run(ctx) })
	eg.Go(func() error { return svrClient.Run(ctx) })

	if err = eg.Wait(); err != nil && !errors.Is(err, context.Canceled) {
		return fmt.Errorf("wait app stop: %v", err)
	}

	return nil
}
