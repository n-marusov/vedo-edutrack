// Package platform provides shared DI providers for the Wire composition root.
//
// The wire.go file declares providers that are common across all modules:
// database connection pool, logger, tracer, meter, and configuration.
// Module-specific providers live in each module's own wire set.
//
// See ADR-IMPL.PROCESS.development-tooling §4 (DI via google/wire).
package platform

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	platformconfig "vedo-edutrack/backend/internal/platform/config"
	platformlogger "vedo-edutrack/backend/internal/platform/logger"
	platformpostgres "vedo-edutrack/backend/internal/platform/postgres"
	platformtelemetry "vedo-edutrack/backend/internal/platform/telemetry"
)

// Platform bundles the initialized platform adapters and their cleanup
// functions for graceful shutdown.
type Platform struct {
	Cfg      *platformconfig.Config
	Logger   *zap.Logger
	LogLevel zapcore.Level
	DB       *pgxpool.Pool
	Shutdown func(ctx context.Context)
}

// InitPlatform initializes all platform-level adapters.
//
// Called once at startup from cmd/vedo-edutrack/main.go (the composition root).
// Returns the initialized bundle or an error if a component fails.
// The returned Shutdown function closes the DB pool, flushes OTel exporters,
// and syncs the logger (safe to call multiple times).
func InitPlatform(ctx context.Context, serviceName string) (*Platform, error) {
	cfg, err := platformconfig.Load()
	if err != nil {
		return nil, err
	}

	zapLogger, err := platformlogger.New(cfg.LogLevel)
	if err != nil {
		return nil, err
	}
	zapLogger.Info("platform initialized",
		zap.String("service", serviceName),
		zap.String("env", cfg.Environment),
		zap.String("log_level", cfg.LogLevel),
	)

	tp, err := platformtelemetry.InitTracer(ctx, serviceName)
	if err != nil {
		zapLogger.Warn("OTel tracer init failed (non-fatal in dev)", zap.Error(err))
	} else {
		zapLogger.Info("OTel tracer ready")
	}

	mp, err := platformtelemetry.InitMeter(ctx, serviceName)
	if err != nil {
		zapLogger.Warn("OTel meter init failed (non-fatal in dev)", zap.Error(err))
	} else {
		zapLogger.Info("OTel meter ready")
	}

	pool, err := platformpostgres.Connect(ctx, cfg.DatabaseURL)
	if err != nil {
		zapLogger.Warn("PostgreSQL connect failed (non-fatal in dev)", zap.Error(err))
	} else {
		zapLogger.Info("PostgreSQL pool ready")
	}

	p := &Platform{
		Cfg:      cfg,
		Logger:   zapLogger,
		LogLevel: levelToZap(cfg.LogLevel),
		DB:       pool,
		Shutdown: func(ctx context.Context) {
			platformpostgres.Close(pool)
			platformtelemetry.ShutdownTracer(ctx, tp)
			platformtelemetry.ShutdownMeter(ctx, mp)
			platformlogger.Sync(zapLogger)
		},
	}
	return p, nil
}

// levelToZap converts a LOG_LEVEL string to a zapcore.Level (best-effort;
// falls back to InfoLevel on parse errors).
func levelToZap(level string) zapcore.Level {
	lvl, err := zapcore.ParseLevel(level)
	if err != nil {
		return zapcore.InfoLevel
	}
	return lvl
}
