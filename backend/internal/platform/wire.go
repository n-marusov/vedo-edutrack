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
	"log"

	platformconfig "vedo-edutrack/backend/internal/platform/config"
	platformlogger "vedo-edutrack/backend/internal/platform/logger"
	platformpostgres "vedo-edutrack/backend/internal/platform/postgres"
	platformtelemetry "vedo-edutrack/backend/internal/platform/telemetry"
)

// InitPlatform initializes all platform-level adapters.
//
// Called once at startup from cmd/vedo-edutrack/main.go (the composition root).
// Returns an error if any platform component fails to initialize.
func InitPlatform(ctx context.Context, serviceName string) error {
	log.Println("[INFO] [platform.InitPlatform] initializing platform adapters")

	cfg, err := platformconfig.Load()
	if err != nil {
		log.Printf("[ERROR] [platform.InitPlatform] failed to load config: %v", err)
		return err
	}

	if err := platformlogger.New(cfg.LogLevel); err != nil {
		log.Printf("[ERROR] [platform.InitPlatform] failed to init logger: %v", err)
		return err
	}

	if err := platformtelemetry.InitTracer(ctx, serviceName); err != nil {
		log.Printf("[ERROR] [platform.InitPlatform] failed to init tracer: %v", err)
		return err
	}

	if err := platformtelemetry.InitMeter(ctx); err != nil {
		log.Printf("[ERROR] [platform.InitPlatform] failed to init meter: %v", err)
		return err
	}

	if err := platformpostgres.Connect(ctx, cfg.DatabaseURL); err != nil {
		log.Printf("[ERROR] [platform.InitPlatform] failed to connect to database: %v", err)
		return err
	}

	log.Println("[INFO] [platform.InitPlatform] all platform adapters initialized")
	return nil
}
