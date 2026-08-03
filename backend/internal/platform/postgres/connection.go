// Package postgres provides shared PostgreSQL connection and migration utilities.
//
// It serves as a platform adapter (Clean Architecture layer) — individual
// bounded contexts use this package to obtain database connections, while
// repository implementations live inside each module's adapters/repository/.
//
// Migrations are managed via Atlas (see ADR-IMPL.PROCESS.development-tooling §3).
package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Connect opens a connection pool to PostgreSQL.
//
// The pool is sized for the M0.3 scaffold: max 20 connections, min 2 idle.
// A ping verifies reachability before returning.
func Connect(ctx context.Context, databaseURL string) (*pgxpool.Pool, error) {
	cfg, err := pgxpool.ParseConfig(databaseURL)
	if err != nil {
		return nil, fmt.Errorf("parse DATABASE_URL: %w", err)
	}

	cfg.MaxConns = 20
	cfg.MinConns = 2
	cfg.MaxConnLifetime = time.Hour
	cfg.HealthCheckPeriod = 30 * time.Second

	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		return nil, fmt.Errorf("create pgx pool: %w", err)
	}

	pingCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	if err := pool.Ping(pingCtx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("ping postgres: %w", err)
	}

	return pool, nil
}

// Close closes the connection pool. Call during graceful shutdown.
func Close(pool *pgxpool.Pool) {
	if pool == nil {
		return
	}
	pool.Close()
}

// Ping verifies database reachability for readiness checks.
func Ping(ctx context.Context, pool *pgxpool.Pool) error {
	if pool == nil {
		return fmt.Errorf("database pool not initialized")
	}
	pingCtx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()
	return pool.Ping(pingCtx)
}
