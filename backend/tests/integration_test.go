package tests

import (
	"context"
	"testing"
	"time"
)

// TODO(M0.3): cross-module integration tests (testcontainers + PostgreSQL):
//   - plan lifecycle across planmanagement / routeplanning / executionprogress
//   - gap diagnosis events (UC-execute.gap.diagnose-root-cause, M2)
//   - ontology sync via the ontologyport ACL (F0.2)
//
// Intentionally skipped at M0.2 (T13) — scaffold only. The helpers below are
// referenced here so the scaffold stays lint-clean until real tests land.
func TestCrossModuleIntegration(t *testing.T) {
	spec := dockerContainerSpec{
		Image:    "postgres:16-alpine",
		Ports:    map[string]string{"5432/tcp": "5432"},
		Env:      map[string]string{"POSTGRES_DB": "edutrack"},
		ReadyCmd: "pg_isready -U edutrack",
	}
	startPostgresTestContainer(t, spec)
	_ = waitForDB(context.Background(), "postgres://edutrack@localhost:5432/edutrack", time.Minute)
	t.Skip("TODO: integration tests with testcontainers (backend/tests, M0.3+)")
}

// dockerContainerSpec describes a testcontainer to launch (M0.3+).
type dockerContainerSpec struct {
	Image    string
	Ports    map[string]string // containerPort -> hostPort
	Env      map[string]string
	ReadyCmd string // command run inside the container to detect readiness
}

// startPostgresTestContainer launches a PostgreSQL testcontainer (M0.3+).
func startPostgresTestContainer(t *testing.T, spec dockerContainerSpec) {
	t.Helper()
	// TODO(M0.3): implement with testcontainers-go:
	//   - start postgres:16-alpine with the given spec
	//   - wait for readiness (ReadyCmd)
	//   - register cleanup via t.Cleanup
	_ = spec
	t.Skip("TODO: implement testcontainer lifecycle (M0.3+)")
}

// waitForDB polls the DSN until it accepts connections or the timeout passes (M0.3+).
func waitForDB(ctx context.Context, dsn string, timeout time.Duration) error {
	_ = ctx
	_ = dsn
	_ = timeout
	// TODO(M0.3): pgx ping loop with context cancellation.
	return nil
}
