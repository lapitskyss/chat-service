package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"os/signal"
	"syscall"

	"go.uber.org/multierr"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"

	keycloakclient "github.com/lapitskyss/chat-service/internal/clients/keycloak"
	"github.com/lapitskyss/chat-service/internal/config"
	"github.com/lapitskyss/chat-service/internal/logger"
	chatsrepo "github.com/lapitskyss/chat-service/internal/repositories/chats"
	jobsrepo "github.com/lapitskyss/chat-service/internal/repositories/jobs"
	messagesrepo "github.com/lapitskyss/chat-service/internal/repositories/messages"
	problemsrepo "github.com/lapitskyss/chat-service/internal/repositories/problems"
	clientv1 "github.com/lapitskyss/chat-service/internal/server-client/v1"
	serverdebug "github.com/lapitskyss/chat-service/internal/server-debug"
	msgproducer "github.com/lapitskyss/chat-service/internal/services/msg-producer"
	"github.com/lapitskyss/chat-service/internal/services/outbox"
	sendclientmessagejob "github.com/lapitskyss/chat-service/internal/services/outbox/jobs/send-client-message"
	"github.com/lapitskyss/chat-service/internal/store"
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

	if err := logger.Init(logger.NewOptions(
		cfg.Log.Level,
		logger.WithProductionMode(cfg.Global.IsProd()),
		logger.WithSentryDSN(cfg.Sentry.DSN),
		logger.WithSentryEnv(cfg.Global.Env),
	)); err != nil {
		return fmt.Errorf("init logger: %v", err)
	}
	defer logger.Sync()

	// Clients.
	kc, err := keycloakclient.New(keycloakclient.NewOptions(
		cfg.Clients.Keycloak.BasePath,
		cfg.Clients.Keycloak.Realm,
		cfg.Clients.Keycloak.ClientID,
		cfg.Clients.Keycloak.ClientSecret,
		keycloakclient.WithDebugMode(cfg.Clients.Keycloak.DebugMode),
	))
	if err != nil {
		return fmt.Errorf("create keycloak client: %v", err)
	}
	if cfg.Global.IsProd() && cfg.Clients.Keycloak.DebugMode {
		zap.L().Warn("keycloak client in the debug mode")
	}

	// Postgres.
	storage, err := store.NewPSQLClient(store.NewPSQLOptions(
		cfg.PSQL.Address,
		cfg.PSQL.User,
		cfg.PSQL.Password,
		cfg.PSQL.Database,
		store.WithDebug(cfg.PSQL.Debug),
	))
	if err != nil {
		return fmt.Errorf("create postgres connecton: %v", err)
	}
	if cfg.Global.IsProd() && cfg.PSQL.Debug {
		zap.L().Warn("psql client in the debug mode")
	}

	defer multierr.AppendInvoke(&errReturned, multierr.Close(storage))

	// Migration.
	err = storage.Schema.Create(ctx)
	if err != nil {
		return fmt.Errorf("auto migration: %v", err)
	}

	// Repository.
	db := store.NewDatabase(storage)
	chatRepo, err := chatsrepo.New(chatsrepo.NewOptions(db))
	if err != nil {
		return fmt.Errorf("chats repository: %v", err)
	}
	msgRepo, err := messagesrepo.New(messagesrepo.NewOptions(db))
	if err != nil {
		return fmt.Errorf("messages repository: %v", err)
	}
	problemRepo, err := problemsrepo.New(problemsrepo.NewOptions(db))
	if err != nil {
		return fmt.Errorf("messages repository: %v", err)
	}
	jobsRepo, err := jobsrepo.New(jobsrepo.NewOptions(db))
	if err != nil {
		return fmt.Errorf("jobs repository: %v", err)
	}

	// Services
	msgProducer, err := msgproducer.New(msgproducer.NewOptions(
		msgproducer.NewKafkaWriter(
			cfg.Services.MsgProducer.Brokers,
			cfg.Services.MsgProducer.Topic,
			cfg.Services.MsgProducer.BatchSize,
		),
		msgproducer.WithEncryptKey(cfg.Services.MsgProducer.EncryptKey),
	))
	if err != nil {
		return fmt.Errorf("message producer service: %v", err)
	}

	outBox, err := outbox.New(outbox.NewOptions(
		cfg.Services.Outbox.Workers,
		cfg.Services.Outbox.IdleTime,
		cfg.Services.Outbox.ReserveFor,
		jobsRepo,
		db,
	))
	if err != nil {
		return fmt.Errorf("outbox service: %v", err)
	}

	sendClientMessageJob, err := sendclientmessagejob.New(sendclientmessagejob.NewOptions(
		msgProducer,
		msgRepo,
	))
	if err != nil {
		return fmt.Errorf("send client message job: %v", err)
	}
	err = outBox.RegisterJob(sendClientMessageJob)
	if err != nil {
		return fmt.Errorf("register send client message job: %v", err)
	}

	err = outBox.Run(ctx)
	if err != nil {
		return fmt.Errorf("run outbox: %v", err)
	}

	// Servers.
	clientV1Swagger, err := clientv1.GetSwagger()
	if err != nil {
		return fmt.Errorf("get client v1 swagger: %v", err)
	}

	srvClient, err := initServerClient(
		cfg.Global.IsProd(),
		cfg.Servers.Client.Addr,
		cfg.Servers.Client.AllowOrigins,
		clientV1Swagger,
		kc,
		cfg.Servers.Client.RequiredAccess.Resource,
		cfg.Servers.Client.RequiredAccess.Role,
		db,
		chatRepo,
		msgRepo,
		problemRepo,
		outBox,
	)
	if err != nil {
		return fmt.Errorf("init client server: %v", err)
	}

	srvDebug, err := serverdebug.New(serverdebug.NewOptions(
		cfg.Servers.Debug.Addr,
		clientV1Swagger,
	))
	if err != nil {
		return fmt.Errorf("init debug server: %v", err)
	}

	eg, ctx := errgroup.WithContext(ctx)

	// Run servers.
	eg.Go(func() error { return srvClient.Run(ctx) })
	eg.Go(func() error { return srvDebug.Run(ctx) })

	if err = eg.Wait(); err != nil && !errors.Is(err, context.Canceled) {
		return fmt.Errorf("wait app stop: %v", err)
	}

	return nil
}
