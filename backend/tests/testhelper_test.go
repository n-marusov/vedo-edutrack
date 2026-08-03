//go:build integration

package tests

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	tc "github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"

	"vedo-edutrack/backend/internal/platform/migrate"
	"vedo-edutrack/backend/migrations"
)

// startPostgres boots a postgres:16-alpine testcontainer, applies the embedded
// migrations and returns a ready connection pool. Shared by all integration
// tests (T27a pattern from integration_test.go).
func startPostgres(t *testing.T) *pgxpool.Pool {
	t.Helper()

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
	t.Cleanup(cancel)

	postgresC, err := tc.GenericContainer(ctx, tc.GenericContainerRequest{
		ContainerRequest: tc.ContainerRequest{
			Image:        "postgres:16-alpine",
			ExposedPorts: []string{"5432/tcp"},
			Env: map[string]string{
				"POSTGRES_USER":     "edutrack",
				"POSTGRES_PASSWORD": "edutrack",
				"POSTGRES_DB":       "edutrack",
			},
			WaitingFor: wait.ForLog("database system is ready to accept connections").WithStartupTimeout(60 * time.Second),
		},
		Started: true,
	})
	if err != nil {
		t.Fatalf("start postgres testcontainer: %v", err)
	}
	t.Cleanup(func() { _ = postgresC.Terminate(context.Background()) })

	host, err := postgresC.Host(ctx)
	if err != nil {
		t.Fatalf("postgres host: %v", err)
	}
	port, err := postgresC.MappedPort(ctx, "5432")
	if err != nil {
		t.Fatalf("postgres port: %v", err)
	}
	dsn := fmt.Sprintf("postgres://edutrack:edutrack@%s:%s/edutrack?sslmode=disable", host, port.Port())

	pool, err := connectWithRetry(ctx, dsn)
	if err != nil {
		t.Fatalf("connect to test postgres: %v", err)
	}
	t.Cleanup(pool.Close)

	runner := migrate.NewRunner(migrations.FS, pool)
	applied, err := runner.Up(ctx)
	if err != nil {
		t.Fatalf("apply migrations: %v", err)
	}
	if applied == 0 {
		t.Fatal("expected at least one migration applied")
	}

	return pool
}
