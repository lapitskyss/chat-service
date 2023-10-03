package logger

import (
	"errors"
	"fmt"
	stdlog "log"
	"os"
	"syscall"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

//go:generate options-gen -out-filename=logger_options.gen.go -from-struct=Options
type Options struct {
	level          string `option:"mandatory" validate:"required,oneof=debug info warn error"`
	productionMode bool
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

	var level zap.AtomicLevel
	switch opts.level {
	case "debug":
		level = zap.NewAtomicLevelAt(zap.DebugLevel)
	case "info":
		level = zap.NewAtomicLevelAt(zap.InfoLevel)
	case "warn":
		level = zap.NewAtomicLevelAt(zap.WarnLevel)
	case "error":
		level = zap.NewAtomicLevelAt(zap.ErrorLevel)
	}

	cfg := zapcore.EncoderConfig{
		MessageKey: "msg",
		LevelKey:   "level",
		NameKey:    "component",
		TimeKey:    "T",
		EncodeTime: zapcore.ISO8601TimeEncoder,
	}

	var enc zapcore.Encoder
	if opts.productionMode {
		cfg.EncodeLevel = zapcore.CapitalLevelEncoder
		enc = zapcore.NewJSONEncoder(cfg)
	} else {
		cfg.EncodeLevel = zapcore.CapitalColorLevelEncoder
		enc = zapcore.NewConsoleEncoder(cfg)
	}

	stdout := zapcore.AddSync(os.Stdout)

	cores := []zapcore.Core{
		zapcore.NewCore(enc, stdout, level),
	}
	l := zap.New(zapcore.NewTee(cores...))
	zap.ReplaceGlobals(l)

	return nil
}

func Sync() {
	if err := zap.L().Sync(); err != nil && !errors.Is(err, syscall.ENOTTY) {
		stdlog.Printf("cannot sync logger: %v", err)
	}
}
