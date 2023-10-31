package logger

import (
	"errors"
	"fmt"
	stdlog "log"
	"os"
	"syscall"

	"github.com/TheZeroSlave/zapsentry"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/lapitskyss/chat-service/internal/buildinfo"
)

var Level zap.AtomicLevel

//go:generate options-gen -out-filename=logger_options.gen.go -from-struct=Options
type Options struct {
	level          string `option:"mandatory" validate:"required,oneof=debug info warn error"`
	productionMode bool
	sentryDSN      string `validate:"omitempty,url"`
	sentryEnv      string `validate:"omitempty,oneof=dev stage prod"`
}

func (c *Options) IsSentryEnabled() bool {
	return c.sentryDSN != ""
}

func MustInit(opts Options) {
	if err := Init(opts); err != nil {
		panic(err)
	}
}

func Init(opts Options) error {
	if err := opts.Validate(); err != nil {
		return fmt.Errorf("validate options: %v", err)
	}

	// Init stdout logs
	stdoutCore, err := stdoutZapcore(opts)
	if err != nil {
		return fmt.Errorf("zapcore stdout: %v", err)
	}
	cores := []zapcore.Core{stdoutCore}

	// Init sentry logs
	if opts.IsSentryEnabled() {
		sentryCore, err := sentryZapcore(opts)
		if err != nil {
			return fmt.Errorf("zapcore sentry: %v", err)
		}
		cores = append(cores, sentryCore)
	}

	l := zap.New(zapcore.NewTee(cores...))
	zap.ReplaceGlobals(l)

	return nil
}

func stdoutZapcore(opts Options) (zapcore.Core, error) {
	var err error
	Level, err = zap.ParseAtomicLevel(opts.level)
	if err != nil {
		return nil, fmt.Errorf("parse level: %v", err)
	}

	cfg := zap.NewProductionEncoderConfig()
	cfg.NameKey = "component"
	cfg.TimeKey = "T"
	cfg.EncodeTime = zapcore.ISO8601TimeEncoder

	var encoder zapcore.Encoder
	if opts.productionMode {
		cfg.EncodeLevel = zapcore.CapitalLevelEncoder
		encoder = zapcore.NewJSONEncoder(cfg)
	} else {
		cfg.EncodeLevel = zapcore.CapitalColorLevelEncoder
		encoder = zapcore.NewConsoleEncoder(cfg)
	}
	return zapcore.NewCore(encoder, os.Stdout, Level), nil
}

func sentryZapcore(opts Options) (zapcore.Core, error) {
	cfg := zapsentry.Configuration{
		Level: zapcore.WarnLevel,
	}
	env := "dev"
	if opts.productionMode {
		env = "prod"
	}
	client, err := NewSentryClient(
		opts.sentryDSN,
		env,
		buildinfo.DepsVersion("github.com/getsentry/sentry-go"),
	)
	if err != nil {
		return nil, fmt.Errorf("new sentry client: %v", err)
	}
	return zapsentry.NewCore(cfg, zapsentry.NewSentryClientFromClient(client))
}

func Sync() {
	if err := zap.L().Sync(); err != nil && !errors.Is(err, syscall.ENOTTY) {
		stdlog.Printf("cannot sync logger: %v", err)
	}
}
