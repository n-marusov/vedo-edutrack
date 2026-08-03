// Package logger provides a shared structured logger (zap).
//
// All bounded contexts use this logger through the platform layer,
// ensuring consistent log format (JSON), correlation IDs (trace_id, span_id),
// and configurable verbosity via LOG_LEVEL.
//
// See ADR-IMPL.PROCESS.development-tooling §4 and
// REQ-NFR-ops.observability.log-level-config.
package logger

import (
	"fmt"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// New creates a new structured JSON logger.
//
// The log level is controlled by the LOG_LEVEL environment variable
// (debug | info | warn | error; default: info). Logs are JSON-formatted and
// routed to stdout (collected by OTel Collector in the docker-compose stack).
func New(level string) (*zap.Logger, error) {
	lvl, err := zapcore.ParseLevel(level)
	if err != nil {
		return nil, fmt.Errorf("invalid LOG_LEVEL %q: %w", level, err)
	}

	cfg := zap.Config{
		Level:            zap.NewAtomicLevelAt(lvl),
		Encoding:         "json",
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
		EncoderConfig: zapcore.EncoderConfig{
			TimeKey:        "ts",
			LevelKey:       "level",
			NameKey:        "logger",
			CallerKey:      "caller",
			MessageKey:     "msg",
			StacktraceKey:  "stacktrace",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    zapcore.LowercaseLevelEncoder,
			EncodeTime:     zapcore.ISO8601TimeEncoder,
			EncodeDuration: zapcore.SecondsDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
		},
	}

	// caller skip=1 so log calls report the caller of the public API,
	// not the internal wrapper.
	opts := []zap.Option{zap.AddCallerSkip(1)}
	if lvl <= zapcore.DebugLevel {
		opts = append(opts, zap.AddStacktrace(zapcore.ErrorLevel))
	}

	logger, err := cfg.Build(opts...)
	if err != nil {
		return nil, fmt.Errorf("build zap logger: %w", err)
	}
	return logger, nil
}

// NewNop returns a no-op logger for tests and short-lived CLI commands.
func NewNop() *zap.Logger {
	return zap.NewNop()
}

// Sync flushes buffered log entries. Call during graceful shutdown.
func Sync(logger *zap.Logger) {
	if logger == nil {
		return
	}
	_ = logger.Sync()
}
