// Package webhook implements the webhook outbox infrastructure (F6).
//
// Outbound webhook events (module.mastered, plan.deviated,
// route.recalculated) are enqueued into a transactional outbox and delivered
// with idempotency guarantees: a duplicate event_id is acknowledged without
// re-dispatch (REQ-FR-api.webhooks.idempotency, ADR-§3).
//
// The outbox is PostgreSQL-backed (integrations.outbox_events) so enqueue
// happens in the same DB transaction as the business change that produced the
// event — the outbox and the domain data cannot diverge
// (ADR-DES.INFRA.modular-monolith-approach: outbox for external delivery).
package webhook

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

// EventType identifies the webhook event contract.
type EventType string

const (
	// EventModuleMastered fires when a learner masters a module.
	EventModuleMastered EventType = "module.mastered"
	// EventPlanDeviated fires when a plan deviates beyond the threshold.
	EventPlanDeviated EventType = "plan.deviated"
	// EventRouteRecalculated fires after a route recomputation.
	EventRouteRecalculated EventType = "route.recalculated"
	// EventStandardRiskDetected fires when a standard deficit is diagnosed.
	EventStandardRiskDetected EventType = "standard.risk_detected"
)

// Event is an outbound webhook payload.
type Event struct {
	EventID string
	Type    EventType
	Payload map[string]any
}

// MaxRetryAttempts is the delivery retry boundary (exponential backoff
// 1s..16s, then permanent failure). Mirrors the domain constant so the worker
// can decide when to stop retrying without importing the domain package.
const MaxRetryAttempts = 5

// defaultPollInterval is how often the outbox worker polls for pending events.
const defaultPollInterval = time.Second

// OutboxRepository is the persistence contract for the webhook outbox.
type OutboxRepository interface {
	// Enqueue adds an event to the outbox. Duplicate event_ids are ignored
	// (unique constraint on outbox_events.event_id).
	Enqueue(ctx context.Context, event Event) error
	// Dequeue returns up to limit pending events ordered by created_at,
	// marking them in-flight so concurrent workers don't double-deliver.
	Dequeue(ctx context.Context, limit int) ([]Event, error)
	// MarkDelivered records a successful delivery (delivered_at set).
	MarkDelivered(ctx context.Context, eventID string) error
	// MarkFailed records a delivery failure with a reason and increments the
	// retry counter. Returns ErrPermanentFailure when the retry boundary is
	// reached (the caller should stop retrying).
	MarkFailed(ctx context.Context, eventID, reason string) error
}

// ErrPermanentFailure is returned by MarkFailed when the event exhausted its
// retry budget and is marked permanently failed.
var ErrPermanentFailure = errors.New("event permanently failed (retry budget exhausted)")

// PostgresOutbox implements OutboxRepository over integrations.outbox_events.
type PostgresOutbox struct {
	pool   *pgxpool.Pool
	logger *zap.Logger
}

// NewPostgresOutbox builds a Postgres-backed outbox over a shared pool.
func NewPostgresOutbox(pool *pgxpool.Pool, logger *zap.Logger) *PostgresOutbox {
	if logger == nil {
		logger = zap.NewNop()
	}
	return &PostgresOutbox{pool: pool, logger: logger.Named("webhook.outbox")}
}

// Enqueue inserts an event into the outbox. Idempotent on duplicate event_id
// (unique constraint): a duplicate is acknowledged without error.
func (o *PostgresOutbox) Enqueue(ctx context.Context, event Event) error {
	if err := validateEvent(event); err != nil {
		return err
	}
	o.logger.Debug("Enqueue", zap.String("eventID", event.EventID), zap.String("eventType", string(event.Type)))
	_, err := o.pool.Exec(ctx,
		`INSERT INTO integrations.outbox_events (event_id, event_type, payload)
		 VALUES ($1, $2, $3)
		 ON CONFLICT (event_id) DO NOTHING`,
		event.EventID, string(event.Type), event.Payload,
	)
	if err != nil {
		return fmt.Errorf("enqueue outbox event %s: %w", event.EventID, err)
	}
	return nil
}

// EnqueueTx inserts an event inside the caller's transaction — used when the
// business change and the outbox write must commit atomically (outbox pattern).
func (o *PostgresOutbox) EnqueueTx(ctx context.Context, tx pgx.Tx, event Event) error {
	if err := validateEvent(event); err != nil {
		return err
	}
	o.logger.Debug("EnqueueTx", zap.String("eventID", event.EventID), zap.String("eventType", string(event.Type)))
	_, err := tx.Exec(ctx,
		`INSERT INTO integrations.outbox_events (event_id, event_type, payload)
		 VALUES ($1, $2, $3)
		 ON CONFLICT (event_id) DO NOTHING`,
		event.EventID, string(event.Type), event.Payload,
	)
	if err != nil {
		return fmt.Errorf("enqueue outbox event %s (tx): %w", event.EventID, err)
	}
	return nil
}

// Dequeue claims up to limit pending events, ordered by created_at. Each
// claim increments retry_count (the delivery-attempt counter) and marks the
// event in_flight so concurrent workers don't re-claim it. The worker uses
// retry_count to decide when to stop retrying (see MarkFailed). Claimed events
// are returned to the caller; concurrent workers skip them via SKIP LOCKED.
func (o *PostgresOutbox) Dequeue(ctx context.Context, limit int) ([]Event, error) {
	if limit <= 0 {
		limit = 10
	}
	o.logger.Debug("Dequeue", zap.Int("batchSize", limit))
	rows, err := o.pool.Query(ctx,
		`UPDATE integrations.outbox_events SET retry_count = retry_count + 1, status = 'in_flight'
		 WHERE id IN (
			SELECT id FROM integrations.outbox_events
			WHERE status = 'pending'
			ORDER BY created_at
			LIMIT $1
			FOR UPDATE SKIP LOCKED
		 )
		 RETURNING event_id, event_type, payload, retry_count`,
		limit,
	)
	if err != nil {
		return nil, fmt.Errorf("dequeue outbox events: %w", err)
	}
	defer rows.Close()

	events := make([]Event, 0, limit)
	for rows.Next() {
		var e Event
		var eventType string
		if err := rows.Scan(&e.EventID, &eventType, &e.Payload, new(int)); err != nil {
			return nil, fmt.Errorf("scan dequeued event: %w", err)
		}
		e.Type = EventType(eventType)
		events = append(events, e)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate dequeued events: %w", err)
	}
	o.logger.Debug("Dequeue", zap.Int("count", len(events)))
	return events, nil
}

// MarkDelivered marks an event delivered.
func (o *PostgresOutbox) MarkDelivered(ctx context.Context, eventID string) error {
	o.logger.Info("Delivered", zap.String("eventID", eventID))
	tag, err := o.pool.Exec(ctx,
		`UPDATE integrations.outbox_events
		 SET status = 'delivered', delivered_at = now()
		 WHERE event_id = $1`, eventID)
	if err != nil {
		return fmt.Errorf("mark outbox event %s delivered: %w", eventID, err)
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("mark delivered: outbox event %s not found", eventID)
	}
	return nil
}

// MarkFailed records a delivery failure. retry_count already counts the
// delivery attempt (incremented by Dequeue); when it reaches the boundary the
// event is marked permanently failed and ErrPermanentFailure is returned.
// Otherwise the event stays pending for the next poll (retry).
func (o *PostgresOutbox) MarkFailed(ctx context.Context, eventID, reason string) error {
	var retryCount int
	err := o.pool.QueryRow(ctx,
		`SELECT retry_count FROM integrations.outbox_events WHERE event_id = $1`,
		eventID,
	).Scan(&retryCount)
	if err != nil {
		return fmt.Errorf("read retry count for outbox event %s: %w", eventID, err)
	}

	if retryCount >= MaxRetryAttempts {
		o.logger.Error("PermanentFailure", zap.String("eventID", eventID), zap.Int("attempts", retryCount), zap.String("reason", reason))
		if _, err := o.pool.Exec(ctx,
			`UPDATE integrations.outbox_events SET status = 'failed' WHERE event_id = $1`, eventID); err != nil {
			return fmt.Errorf("mark outbox event %s permanently failed: %w", eventID, err)
		}
		return ErrPermanentFailure
	}

	o.logger.Warn("DeliveryFailed", zap.String("eventID", eventID), zap.Int("attempt", retryCount), zap.String("reason", reason))
	if _, err := o.pool.Exec(ctx,
		`UPDATE integrations.outbox_events SET status = 'pending' WHERE event_id = $1`, eventID); err != nil {
		return fmt.Errorf("reset outbox event %s to pending: %w", eventID, err)
	}
	return nil
}

// validateEvent enforces the outbox event contract (moved from the in-memory
// prototype; keep behavior identical).
func validateEvent(event Event) error {
	if event.EventID == "" {
		return errors.New("event id is required")
	}
	if event.Type == "" {
		return errors.New("event type is required")
	}
	return nil
}

// Worker polls the outbox on an interval and dispatches events to a handler.
// The handler decides what "delivered" means (HTTP POST in Task 8).
type Worker struct {
	repo         OutboxRepository
	handler      func(ctx context.Context, event Event) error
	batchSize    int
	pollInterval time.Duration
	logger       *zap.Logger
	mu           sync.Mutex
	running      bool
}

// WorkerOption configures the outbox worker.
type WorkerOption func(*Worker)

// WithPollInterval overrides the default 1s poll interval (tests use a short
// interval).
func WithPollInterval(d time.Duration) WorkerOption {
	return func(w *Worker) { w.pollInterval = d }
}

// WithBatchSize overrides the default batch size (10).
func WithBatchSize(n int) WorkerOption {
	return func(w *Worker) { w.batchSize = n }
}

// NewWorker builds an outbox polling worker. The handler is invoked for every
// dequeued event; returning nil marks the event delivered, an error marks it
// failed (with retry accounting).
func NewWorker(repo OutboxRepository, handler func(ctx context.Context, event Event) error, logger *zap.Logger, opts ...WorkerOption) *Worker {
	if logger == nil {
		logger = zap.NewNop()
	}
	w := &Worker{
		repo:         repo,
		handler:      handler,
		batchSize:    10,
		pollInterval: defaultPollInterval,
		logger:       logger.Named("webhook.worker"),
	}
	for _, opt := range opts {
		opt(w)
	}
	return w
}

// Run polls the outbox until the context is canceled. The handler is invoked
// sequentially (one event at a time) so retry/backoff accounting stays simple;
// a concurrent fan-out is added in the hardening pass (Task 15).
func (w *Worker) Run(ctx context.Context) {
	w.mu.Lock()
	w.running = true
	w.mu.Unlock()
	defer func() {
		w.mu.Lock()
		w.running = false
		w.mu.Unlock()
	}()

	w.logger.Info("worker started", zap.Int("batchSize", w.batchSize), zap.Duration("pollInterval", w.pollInterval))
	ticker := time.NewTicker(w.pollInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			w.logger.Info("worker stopped")
			return
		case <-ticker.C:
			w.pollOnce(ctx)
		}
	}
}

// pollOnce dequeues one batch and dispatches each event.
func (w *Worker) pollOnce(ctx context.Context) {
	w.logger.Debug("Polling", zap.Int("batchSize", w.batchSize))
	events, err := w.repo.Dequeue(ctx, w.batchSize)
	if err != nil {
		w.logger.Warn("dequeue failed", zap.Error(err))
		return
	}
	for _, event := range events {
		w.handleOne(ctx, event)
	}
}

// handleOne dispatches a single event through the handler and records the
// delivery verdict.
func (w *Worker) handleOne(ctx context.Context, event Event) {
	if err := w.handler(ctx, event); err != nil {
		w.logger.Warn("DeliveryRetry", zap.String("eventID", event.EventID), zap.Error(err))
		if markErr := w.repo.MarkFailed(ctx, event.EventID, err.Error()); markErr != nil {
			if errors.Is(markErr, ErrPermanentFailure) {
				w.logger.Error("PermanentFailure", zap.String("eventID", event.EventID), zap.Int("attempts", MaxRetryAttempts))
				return
			}
			w.logger.Error("mark failed", zap.String("eventID", event.EventID), zap.Error(markErr))
		}
		return
	}
	if err := w.repo.MarkDelivered(ctx, event.EventID); err != nil {
		w.logger.Error("mark delivered", zap.String("eventID", event.EventID), zap.Error(err))
	}
}
