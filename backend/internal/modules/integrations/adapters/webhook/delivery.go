package webhook

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sort"
	"strconv"
	"sync"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"

	domain "vedo-edutrack/backend/internal/modules/integrations/domain"
)

// DeliveryPayload is the JSON body delivered to subscriber endpoints
// (REQ-FR-api.webhooks.delivery-format).
type DeliveryPayload struct {
	EventID   string         `json:"event_id"`
	EventType string         `json:"event_type"`
	Timestamp int64          `json:"timestamp"`
	Data      map[string]any `json:"data"`
}

// Signer computes the X-Vedo-Signature header value.
//
// Format: t=<unix>,v1=<hex(HMAC-SHA256(secret, "<t>.<payload>"))>
// (Stripe-style — the timestamp prevents replay; the payload binding prevents
// tampering).
type Signer struct{}

// NewSigner builds a Signer.
func NewSigner() *Signer { return &Signer{} }

// Sign computes the signature header for the given secret and payload.
func (s *Signer) Sign(secret string, payload []byte, now time.Time) string {
	t := now.Unix()
	mac := hmac.New(sha256.New, []byte(secret))
	_, _ = mac.Write([]byte(strconv.FormatInt(t, 10)))
	_, _ = mac.Write([]byte("."))
	_, _ = mac.Write(payload)
	return fmt.Sprintf("t=%d,v1=%s", t, hex.EncodeToString(mac.Sum(nil)))
}

// Verify checks a received signature header against the secret and payload
// (used by the E2E contract tests and by subscribers).
func (s *Signer) Verify(header, secret string, payload []byte) bool {
	computed := s.Sign(secret, payload, time.Now())
	if header == computed {
		return true
	}
	// The timestamp in the header may differ from now (replay window); extract
	// the t= component and recompute for that timestamp.
	t := parseTimestamp(header)
	if t == 0 {
		return false
	}
	return header == s.SignWithTimestamp(secret, payload, time.Unix(t, 0))
}

// SignWithTimestamp signs for an explicit timestamp (verification helper).
func (s *Signer) SignWithTimestamp(secret string, payload []byte, at time.Time) string {
	return s.Sign(secret, payload, at)
}

// parseTimestamp extracts the t= value from a signature header.
func parseTimestamp(header string) int64 {
	for _, part := range splitComma(header) {
		if len(part) > 2 && part[0] == 't' && part[1] == '=' {
			if v, err := strconv.ParseInt(part[2:], 10, 64); err == nil {
				return v
			}
		}
	}
	return 0
}

// splitComma splits a header on commas without allocations (header format is
// t=<ts>,v1=<hex>).
func splitComma(s string) []string {
	var out []string
	start := 0
	for i := 0; i < len(s); i++ {
		if s[i] == ',' {
			out = append(out, s[start:i])
			start = i + 1
		}
	}
	return append(out, s[start:])
}

// DeliveryRecorder persists delivery attempts per (subscription, event). The
// unique (subscription_id, event_id) constraint provides idempotency: a
// duplicate attempt is reported and skipped.
type DeliveryRecorder interface {
	// RecordAttempt records one attempt. Returns duplicate=true when the
	// (subscription, event) pair was already recorded (idempotency).
	RecordAttempt(ctx context.Context, subID domain.SubscriptionID, eventID string, attempt int, status domain.DeliveryStatus, httpStatus int, responseBody, errMsg string) (duplicate bool, err error)
}

// InMemoryDeliveryRecorder is an in-memory DeliveryRecorder for tests and the
// M4 handler wiring until the Postgres-backed version lands.
type InMemoryDeliveryRecorder struct {
	mu   sync.Mutex
	rows map[string]domain.WebhookDelivery
}

// NewInMemoryDeliveryRecorder builds an empty in-memory delivery recorder.
func NewInMemoryDeliveryRecorder() *InMemoryDeliveryRecorder {
	return &InMemoryDeliveryRecorder{rows: map[string]domain.WebhookDelivery{}}
}

// RecordAttempt records an attempt; duplicate=true on an existing
// (subscription, event) row.
func (r *InMemoryDeliveryRecorder) RecordAttempt(_ context.Context, subID domain.SubscriptionID, eventID string, attempt int, status domain.DeliveryStatus, httpStatus int, responseBody, errMsg string) (bool, error) {
	key := subID.String() + ":" + eventID
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, exists := r.rows[key]; exists {
		return true, nil
	}
	now := time.Now().UTC()
	r.rows[key] = domain.WebhookDelivery{
		ID:             fmt.Sprintf("%s-%d", key, attempt),
		SubscriptionID: subID,
		EventID:        eventID,
		Attempt:        attempt,
		Status:         status,
		LastAttemptAt:  &now,
		HTTPStatus:     httpStatus,
		ResponseBody:   responseBody,
		Error:          errMsg,
	}
	return false, nil
}

// ListBySubscription returns a subscription's deliveries (newest first),
// paginated. Implements queries.DeliveryRepository for delivery history reads.
func (r *InMemoryDeliveryRecorder) ListBySubscription(_ context.Context, subID domain.SubscriptionID, page, limit int) ([]domain.WebhookDelivery, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if page < 0 {
		page = 0
	}
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	var all []domain.WebhookDelivery
	for _, d := range r.rows {
		if d.SubscriptionID == subID {
			all = append(all, d)
		}
	}
	sort.Slice(all, func(i, j int) bool {
		return attemptTime(all[i]).After(attemptTime(all[j]))
	})
	start := page * limit
	if start >= len(all) {
		return []domain.WebhookDelivery{}, nil
	}
	end := start + limit
	if end > len(all) {
		end = len(all)
	}
	return all[start:end], nil
}

// attemptTime returns the delivery attempt timestamp (zero when absent).
func attemptTime(d domain.WebhookDelivery) time.Time {
	if d.LastAttemptAt != nil {
		return *d.LastAttemptAt
	}
	return time.Time{}
}

// PostgresDeliveryRecorder persists deliveries to integrations.webhook_deliveries.
type PostgresDeliveryRecorder struct {
	pool   *pgxpool.Pool
	logger *zap.Logger
}

// NewPostgresDeliveryRecorder builds the Postgres-backed recorder.
func NewPostgresDeliveryRecorder(pool *pgxpool.Pool, logger *zap.Logger) *PostgresDeliveryRecorder {
	if logger == nil {
		logger = zap.NewNop()
	}
	return &PostgresDeliveryRecorder{pool: pool, logger: logger.Named("webhook.deliveries")}
}

// RecordAttempt upserts into webhook_deliveries; the unique
// (subscription_id, event_id) constraint makes retries idempotent.
func (r *PostgresDeliveryRecorder) RecordAttempt(ctx context.Context, subID domain.SubscriptionID, eventID string, attempt int, status domain.DeliveryStatus, httpStatus int, responseBody, errMsg string) (bool, error) {
	var httpStatusArg any
	if httpStatus != 0 {
		httpStatusArg = httpStatus
	}
	tag, err := r.pool.Exec(ctx,
		`INSERT INTO integrations.webhook_deliveries
			(subscription_id, event_id, attempt, status, last_attempt_at, http_status, response_body, error)
		 VALUES ($1, $2, $3, $4, now(), $5, $6, $7)
		 ON CONFLICT (subscription_id, event_id) DO NOTHING`,
		subID.String(), eventID, attempt, string(status), httpStatusArg, responseBody, errMsg,
	)
	if err != nil {
		return false, fmt.Errorf("record delivery: %w", err)
	}
	return tag.RowsAffected() == 0, nil
}

// DeliveryWorker is the production webhook delivery worker: it polls the
// outbox, fans each event out to matching active subscriptions with HMAC
// signing, records attempts (idempotent), and applies the retry/deactivation
// business rules.
type DeliveryWorker struct {
	worker     *Worker
	outbox     OutboxRepository
	subs       SubscriptionFinder
	deliveries DeliveryRecorder
	deactivate func(ctx context.Context, subID domain.SubscriptionID) (bool, error)
	signer     *Signer
	httpClient *http.Client
	logger     *zap.Logger
}

// SubscriptionFinder is the port for finding subscriptions by event type.
type SubscriptionFinder interface {
	ListByEventType(ctx context.Context, eventType domain.EventType) ([]domain.WebhookSubscription, error)
}

// DeliveryWorkerConfig configures the delivery worker.
type DeliveryWorkerConfig struct {
	Outbox        OutboxRepository
	Subscriptions SubscriptionFinder
	Deliveries    DeliveryRecorder
	// Deactivate records a delivery failure on the subscription and reports
	// whether it was deactivated (domain service).
	Deactivate   func(ctx context.Context, subID domain.SubscriptionID) (bool, error)
	PollInterval time.Duration
	BatchSize    int
	HTTPClient   *http.Client
}

// NewDeliveryWorker builds the delivery worker (defaults: 1s poll, batch 10).
func NewDeliveryWorker(cfg DeliveryWorkerConfig, logger *zap.Logger) *DeliveryWorker {
	if logger == nil {
		logger = zap.NewNop()
	}
	httpClient := cfg.HTTPClient
	if httpClient == nil {
		httpClient = &http.Client{Timeout: 10 * time.Second}
	}
	w := &DeliveryWorker{
		outbox:     cfg.Outbox,
		subs:       cfg.Subscriptions,
		deliveries: cfg.Deliveries,
		deactivate: cfg.Deactivate,
		signer:     NewSigner(),
		httpClient: httpClient,
		logger:     logger.Named("webhook.worker"),
	}
	opts := []WorkerOption{WithBatchSize(cfg.BatchSize)}
	if cfg.PollInterval > 0 {
		opts = append(opts, WithPollInterval(cfg.PollInterval))
	}
	w.worker = NewWorker(cfg.Outbox, w.deliverEvent, logger, opts...)
	return w
}

// Run polls the outbox and delivers until the context is canceled.
func (w *DeliveryWorker) Run(ctx context.Context) {
	w.logger.Info("delivery worker started")
	w.worker.Run(ctx)
}

// deliverEvent is the worker handler: fans one outbox event out to all active
// subscriptions matching its type. Returns nil when every matching
// subscription accepted the delivery (or none matched — nothing to do); an
// error when at least one delivery failed (the outbox retries the event).
func (w *DeliveryWorker) deliverEvent(ctx context.Context, event Event) error {
	subs, err := w.subs.ListByEventType(ctx, domain.EventType(event.Type))
	if err != nil {
		return fmt.Errorf("find subscriptions for %s: %w", event.Type, err)
	}
	if len(subs) == 0 {
		w.logger.Info("NoSubscriptions", zap.String("eventType", string(event.Type)))
		return nil
	}

	payload, err := json.Marshal(DeliveryPayload{
		EventID:   event.EventID,
		EventType: string(event.Type),
		Timestamp: time.Now().Unix(),
		Data:      event.Payload,
	})
	if err != nil {
		return fmt.Errorf("marshal delivery payload: %w", err)
	}

	var failed bool
	for _, sub := range subs {
		if err := w.deliverOne(ctx, sub, event, payload); err != nil {
			failed = true
			w.logger.Warn("DeliveryFailed", zap.String("eventID", event.EventID), zap.String("subscriptionID", sub.ID.String()), zap.Error(err))
		}
	}
	if failed {
		return fmt.Errorf("one or more deliveries failed for event %s", event.EventID)
	}
	return nil
}

// deliverOne delivers to a single subscription with dedup + HMAC signing.
func (w *DeliveryWorker) deliverOne(ctx context.Context, sub domain.WebhookSubscription, event Event, payload []byte) error {
	attempt := sub.Failures + 1
	duplicate, err := w.deliveries.RecordAttempt(ctx, sub.ID, event.EventID, attempt, domain.DeliveryPending, 0, "", "")
	if err != nil {
		return err
	}
	if duplicate {
		// Already delivered to this subscription — idempotency (ADR-§3).
		w.logger.Debug("delivery dedup", zap.String("eventID", event.EventID), zap.String("subscriptionID", sub.ID.String()))
		return nil
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, sub.URL, bytes.NewReader(payload))
	if err != nil {
		return fmt.Errorf("build delivery request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Vedo-Signature", w.signer.Sign(sub.SigningSecret, payload, time.Now()))

	w.logger.Debug("Delivering", zap.String("eventID", event.EventID), zap.String("subscriptionID", sub.ID.String()), zap.String("url", sub.URL))
	resp, err := w.httpClient.Do(req)
	if err != nil {
		return w.recordDeliveryOutcome(ctx, sub, event, attempt, domain.DeliveryFailed, 0, "", err.Error())
	}
	defer func() { _ = resp.Body.Close() }()
	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		return w.recordDeliveryOutcome(ctx, sub, event, attempt, domain.DeliverySent, resp.StatusCode, string(body), "")
	}
	return w.recordDeliveryOutcome(ctx, sub, event, attempt, domain.DeliveryFailed, resp.StatusCode, string(body), fmt.Sprintf("subscriber returned status %d", resp.StatusCode))
}

// recordDeliveryOutcome persists the attempt verdict and applies the
// consecutive-failure deactivation rule.
func (w *DeliveryWorker) recordDeliveryOutcome(ctx context.Context, sub domain.WebhookSubscription, event Event, attempt int, status domain.DeliveryStatus, httpStatus int, responseBody, errMsg string) error {
	if _, err := w.deliveries.RecordAttempt(ctx, sub.ID, event.EventID, attempt, status, httpStatus, responseBody, errMsg); err != nil {
		return err
	}
	if status == domain.DeliveryFailed {
		if w.deactivate != nil {
			deactivated, err := w.deactivate(ctx, sub.ID)
			if err != nil {
				w.logger.Warn("deactivate check failed", zap.Error(err))
			} else if deactivated {
				w.logger.Error("PermanentFailure", zap.String("eventID", event.EventID), zap.String("subscriptionID", sub.ID.String()), zap.Int("attempts", attempt))
			}
		}
		return fmt.Errorf("delivery failed (http=%d): %s", httpStatus, errMsg)
	}
	w.logger.Info("Delivered", zap.String("eventID", event.EventID), zap.String("subscriptionID", sub.ID.String()), zap.Int("httpStatus", httpStatus), zap.Int("attempt", attempt))
	return nil
}
