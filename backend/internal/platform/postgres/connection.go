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
	"log"
)

// Connect opens a connection pool to PostgreSQL.
//
// TODO: Implement connection pool creation via pgxpool.
// See ADR-DES.DATA.storage-strategy for sqlc+pgx integration details.
func Connect(ctx context.Context, databaseURL string) error {
	log.Printf("[INFO] [postgres.Connect] connecting to database (URL redacted)")
	_ = ctx
	_ = databaseURL
	return nil
}
