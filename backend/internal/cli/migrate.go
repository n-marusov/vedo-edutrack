package cli

import (
	"context"
	"fmt"
	"time"

	"github.com/spf13/cobra"
	"go.uber.org/zap"

	"vedo-edutrack/backend/internal/platform/config"
	"vedo-edutrack/backend/internal/platform/migrate"
	platformpostgres "vedo-edutrack/backend/internal/platform/postgres"
	"vedo-edutrack/backend/migrations"
)

// newMigrateCmd builds the `migrate` subcommand — embedded Atlas-style
// migrations (ADR-DES.API.cli-interface). Subcommands: up, down, validate.
func newMigrateCmd() *cobra.Command {
	migrateCmd := &cobra.Command{
		Use:   "migrate",
		Short: "Manage database migrations (embedded, Atlas-style)",
		RunE: func(cmd *cobra.Command, _ []string) error {
			zapLogger.Info("migrate subcommand invoked (no subcommand)")
			return cmd.Help()
		},
	}

	migrateCmd.AddCommand(
		&cobra.Command{
			Use:   "up",
			Short: "Apply all pending migrations",
			RunE: func(_ *cobra.Command, _ []string) error {
				return runMigrateUp()
			},
		},
		&cobra.Command{
			Use:   "down",
			Short: "Revert the last migration",
			RunE: func(_ *cobra.Command, _ []string) error {
				return runMigrateDown()
			},
		},
		&cobra.Command{
			Use:   "validate",
			Short: "Validate migration drift (drift = 0)",
			RunE: func(_ *cobra.Command, _ []string) error {
				return runMigrateValidate()
			},
		},
	)

	return migrateCmd
}

// stubNotImplemented returns a RunE that logs and prints "not yet implemented".
// Used by commands whose implementation lands in later M1 phases.
func stubNotImplemented(cmdPath string) func(*cobra.Command, []string) error {
	return func(_ *cobra.Command, _ []string) error {
		zapLogger.Info("stub subcommand invoked", zap.String("command", cmdPath))
		fmt.Printf("%s: not yet implemented\n", cmdPath)
		return nil
	}
}

// runMigrateUp applies all pending embedded migrations.
func runMigrateUp() error {
	zapLogger.Info("[migrate] up started")

	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	pool, err := platformpostgres.Connect(ctx, cfg.DatabaseURL)
	if err != nil {
		return fmt.Errorf("connect to database: %w", err)
	}
	defer platformpostgres.Close(pool)

	runner := migrate.NewRunner(migrations.FS, pool)
	applied, err := runner.Up(ctx)
	if err != nil {
		return fmt.Errorf("apply migrations: %w", err)
	}

	zapLogger.Info("[migrate] up completed", zap.Int("applied", applied))
	fmt.Printf("migrations applied: %d\n", applied)
	return nil
}

// runMigrateDown reverts the last applied migration.
func runMigrateDown() error {
	zapLogger.Info("[migrate] down started")

	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	pool, err := platformpostgres.Connect(ctx, cfg.DatabaseURL)
	if err != nil {
		return fmt.Errorf("connect to database: %w", err)
	}
	defer platformpostgres.Close(pool)

	runner := migrate.NewRunner(migrations.FS, pool)
	reverted, err := runner.Down(ctx)
	if err != nil {
		return fmt.Errorf("revert migration: %w", err)
	}

	zapLogger.Info("[migrate] down completed", zap.String("reverted", reverted))
	fmt.Printf("migration reverted: %s\n", reverted)
	return nil
}

// runMigrateValidate checks for drift between embedded and applied migrations.
func runMigrateValidate() error {
	zapLogger.Info("[migrate] validate started")

	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	pool, err := platformpostgres.Connect(ctx, cfg.DatabaseURL)
	if err != nil {
		return fmt.Errorf("connect to database: %w", err)
	}
	defer platformpostgres.Close(pool)

	runner := migrate.NewRunner(migrations.FS, pool)
	drift, err := runner.Validate(ctx)
	if err != nil {
		return fmt.Errorf("validate migrations: %w", err)
	}

	if len(drift) == 0 {
		zapLogger.Info("[migrate] validate passed — drift = 0")
		fmt.Println("migration drift: 0")
		return nil
	}

	zapLogger.Warn("[migrate] validate failed — drift detected", zap.Int("drift_count", len(drift)))
	for _, d := range drift {
		fmt.Printf("  drift: %s\n", d)
	}
	return fmt.Errorf("migration drift detected (%d item(s))", len(drift))
}
