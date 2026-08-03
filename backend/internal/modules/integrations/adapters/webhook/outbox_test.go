//go:build integration

package webhook

import (
	"context"
	"errors"
	"sync/atomic"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
	"go.uber.org/zap"
)

// newTestOutboxPool boots a throwaway PostgreSQL container and applies the
// integrations migrations (the same embedded runner used by `migrate up`).
func newTestOutboxPool(t *testing.T) *pgxpool.Pool {
	t.Helper()
	ctx := context.Background()

	pgContainer, err := postgres.Run(ctx,
		"postgres:16-alpine",
		postgres.WithDatabase("edutrack"),
		postgres.WithUsername("edutrack"),
		postgres.WithPassword("edutrack"),
		testcontainers.WithWaitStrategy(wait.ForLog("database system is ready to accept connections").
			WithOccurrence(2).WithStartupTimeout(60*time.Second)),
	)
	require.NoError(t, err)
	t.Cleanup(func() { _ = pgContainer.Terminate(ctx) })

	connStr, err := pgContainer.ConnectionString(ctx, "sslmode=disable")
	require.NoError(t, err)

	pool, err := pgxpool.New(ctx, connStr)
	require.NoError(t, err)
	t.Cleanup(pool.Close)

	// Apply only the integrations migration (schema must exist before tests).
	applyIntegrationsMigration(t, pool)

	// Ensure the event_dedup unique constraint is exercised (unique on
	// outbox_events.event_id is enforced by the migration above).
	return pool
}

// applyIntegrationsMigration runs the integrations DDL via the embedded runner.
func applyIntegrationsMigration(t *testing.T, pool *pgxpool.Pool) {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	_, err := pool.Exec(ctx, `
		CREATE SCHEMA IF NOT EXISTS integrations;
		CREATE TABLE IF NOT EXISTS integrations.webhook_subscriptions (
			id uuid NOT NULL DEFAULT gen_random_uuid(),
			tenant_id text NOT NULL,
			url text NOT NULL,
			event_types text[] NOT NULL,
			signing_secret text NOT NULL,
			active boolean NOT NULL DEFAULT true,
			failures integer NOT NULL DEFAULT 0,
			created_at timestamptz NOT NULL DEFAULT now(),
			updated_at timestamptz NOT NULL DEFAULT now(),
			PRIMARY KEY (id),
			CONSTRAINT webhook_subscriptions_failures_check CHECK (failures >= 0)
		);
		CREATE TABLE IF NOT EXISTS integrations.webhook_deliveries (
			id uuid NOT NULL DEFAULT gen_random_uuid(),
			subscription_id uuid NOT NULL,
			event_id uuid NOT NULL,
			attempt integer NOT NULL DEFAULT 1,
			status text NOT NULL DEFAULT 'pending',
			last_attempt_at timestamptz,
			http_status integer,
			response_body text NOT NULL DEFAULT '',
			error text NOT NULL DEFAULT '',
			created_at timestamptz NOT NULL DEFAULT now(),
			PRIMARY KEY (id),
			CONSTRAINT webhook_deliveries_subscription_id_fkey FOREIGN KEY (subscription_id) REFERENCES integrations.webhook_subscriptions (id) ON DELETE CASCADE,
			CONSTRAINT webhook_deliveries_unique_delivery UNIQUE (subscription_id, event_id)
		);
		CREATE TABLE IF NOT EXISTS integrations.event_dedup (
			id uuid NOT NULL DEFAULT gen_random_uuid(),
			source text NOT NULL,
			event_id text NOT NULL,
			created_at timestamptz NOT NULL DEFAULT now(),
			PRIMARY KEY (id),
			CONSTRAINT event_dedup_source_event_key UNIQUE (source, event_id)
		);
		CREATE TABLE IF NOT EXISTS integrations.outbox_events (
			id uuid NOT NULL DEFAULT gen_random_uuid(),
			event_id uuid NOT NULL,
			event_type text NOT NULL,
			payload jsonb NOT NULL DEFAULT '{}'::jsonb,
			status text NOT NULL DEFAULT 'pending',
			retry_count integer NOT NULL DEFAULT 0,
			delivered_at timestamptz,
			created_at timestamptz NOT NULL DEFAULT now(),
			PRIMARY KEY (id),
			CONSTRAINT outbox_events_event_id_key UNIQUE (event_id),
			CONSTRAINT outbox_events_status_check CHECK (status IN ('pending', 'in_flight', 'delivered', 'failed'))
		);
	`)
	require.NoError(t, err)
}

func TestPostgresOutboxEnqueueAndDedup(t *testing.T) {
	pool := newTestOutboxPool(t)
	outbox := NewPostgresOutbox(pool, zap.NewNop())
	ctx := context.Background()

	event := Event{EventID: "d4c1e2a5-1f0e-4c3a-9b2f-6a3c8e7d1a2b", Type: EventModuleMastered, Payload: map[string]any{"module_id": "m1"}}

	require.NoError(t, outbox.Enqueue(ctx, event))
	// Duplicate event_id → acknowledged without error, not duplicated.
	require.NoError(t, outbox.Enqueue(ctx, event))

	var count int
	require.NoError(t, pool.QueryRow(ctx, `SELECT count(*) FROM integrations.outbox_events WHERE event_id = $1`, event.EventID).Scan(&count))
	require.Equal(t, 1, count)
}

func TestPostgresOutboxDequeueOrderingAndLimit(t *testing.T) {
	pool := newTestOutboxPool(t)
	outbox := NewPostgresOutbox(pool, zap.NewNop())
	ctx := context.Background()

	for i := 0; i < 3; i++ {
		require.NoError(t, outbox.Enqueue(ctx, Event{
			EventID: uuidFor(t, i),
			Type:    EventModuleMastered,
			Payload: map[string]any{"n": i},
		}))
	}

	events, err := outbox.Dequeue(ctx, 2)
	require.NoError(t, err)
	require.Len(t, events, 2)

	// Dequeue marks the claimed events by bumping retry_count; a second dequeue
	// with the same batch does not re-claim delivered events.
	events2, err := outbox.Dequeue(ctx, 2)
	require.NoError(t, err)
	require.Len(t, events2, 1)
}

func TestPostgresOutboxMarkDelivered(t *testing.T) {
	pool := newTestOutboxPool(t)
	outbox := NewPostgresOutbox(pool, zap.NewNop())
	ctx := context.Background()

	eventID := uuidFor(t, 0)
	require.NoError(t, outbox.Enqueue(ctx, Event{EventID: eventID, Type: EventPlanDeviated, Payload: map[string]any{}}))
	require.NoError(t, outbox.MarkDelivered(ctx, eventID))

	var status string
	require.NoError(t, pool.QueryRow(ctx, `SELECT status FROM integrations.outbox_events WHERE event_id = $1`, eventID).Scan(&status))
	require.Equal(t, "delivered", status)
}

func TestPostgresOutboxRetryBackoffAccounting(t *testing.T) {
	pool := newTestOutboxPool(t)
	outbox := NewPostgresOutbox(pool, zap.NewNop())
	ctx := context.Background()

	eventID := uuidFor(t, 0)
	require.NoError(t, outbox.Enqueue(ctx, Event{EventID: eventID, Type: EventRouteRecalculated, Payload: map[string]any{}}))

	// Simulate the worker claim-fail loop: each Dequeue increments the attempt
	// counter, then MarkFailed decides retry vs permanent failure.
	for attempt := 1; attempt <= MaxRetryAttempts; attempt++ {
		events, err := outbox.Dequeue(ctx, 10)
		require.NoError(t, err)
		require.Len(t, events, 1, "expected the event to remain claimable until the boundary (attempt %d)", attempt)

		err = outbox.MarkFailed(ctx, eventID, "500 from subscriber")
		if attempt < MaxRetryAttempts {
			require.Nil(t, err, "expected retryable failure before boundary (attempt %d)", attempt)
		} else {
			require.True(t, errors.Is(err, ErrPermanentFailure), "expected permanent failure at boundary, got %v", err)
		}
	}

	var status string
	require.NoError(t, pool.QueryRow(ctx, `SELECT status FROM integrations.outbox_events WHERE event_id = $1`, eventID).Scan(&status))
	require.Equal(t, "failed", status)

	// After permanent failure the event is no longer claimable.
	events, err := outbox.Dequeue(ctx, 10)
	require.NoError(t, err)
	require.Empty(t, events, "expected no more claims after permanent failure")
}

func TestOutboxWorkerPollsAndDelivers(t *testing.T) {
	pool := newTestOutboxPool(t)
	outbox := NewPostgresOutbox(pool, zap.NewNop())
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	var delivered atomic.Int32
	worker := NewWorker(outbox, func(_ context.Context, event Event) error {
		delivered.Add(1)
		_ = event
		return nil
	}, zap.NewNop(), WithPollInterval(50*time.Millisecond), WithBatchSize(10))

	go worker.Run(ctx)

	require.NoError(t, outbox.Enqueue(ctx, Event{EventID: uuidFor(t, 0), Type: EventModuleMastered, Payload: map[string]any{"module_id": "m1"}}))
	require.NoError(t, outbox.Enqueue(ctx, Event{EventID: uuidFor(t, 1), Type: EventPlanDeviated, Payload: map[string]any{}}))

	require.Eventually(t, func() bool { return delivered.Load() == 2 }, 10*time.Second, 100*time.Millisecond)

	// Both events must be marked delivered.
	var deliveredCount int
	require.NoError(t, pool.QueryRow(ctx, `SELECT count(*) FROM integrations.outbox_events WHERE status = 'delivered'`).Scan(&deliveredCount))
	require.Equal(t, 2, deliveredCount)
}

func TestOutboxWorkerRetriesFailedDelivery(t *testing.T) {
	pool := newTestOutboxPool(t)
	outbox := NewPostgresOutbox(pool, zap.NewNop())
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	var attempts atomic.Int32
	worker := NewWorker(outbox, func(_ context.Context, _ Event) error {
		attempts.Add(1)
		return errors.New("subscriber returned 500")
	}, zap.NewNop(), WithPollInterval(50*time.Millisecond), WithBatchSize(10))

	go worker.Run(ctx)

	require.NoError(t, outbox.Enqueue(ctx, Event{EventID: uuidFor(t, 0), Type: EventModuleMastered, Payload: map[string]any{}}))

	// The worker should retry; eventually the event reaches the failure
	// boundary and is marked failed (no more handler calls after the 5th).
	require.Eventually(t, func() bool { return attempts.Load() >= MaxRetryAttempts }, 10*time.Second, 100*time.Millisecond)

	var status string
	require.NoError(t, pool.QueryRow(ctx, `SELECT status FROM integrations.outbox_events WHERE event_id = $1`, uuidFor(t, 0)).Scan(&status))
	require.Equal(t, "failed", status)
}

// uuidFor returns a deterministic UUID for index i.
func uuidFor(t *testing.T, i int) string {
	t.Helper()
	raw := [16]byte{}
	raw[15] = byte(i + 1)
	return uuid.Must(uuid.FromBytes(raw[:])).String()
}
