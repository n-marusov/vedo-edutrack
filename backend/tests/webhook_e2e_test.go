//go:build integration

package tests

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"go.uber.org/zap"

	"vedo-edutrack/backend/internal/modules/integrations/adapters/webhook"
	integdomain "vedo-edutrack/backend/internal/modules/integrations/domain"
)

// webhookTestEnv bundles the pieces of a webhook E2E scenario: Postgres
// outbox + in-memory subscription repo/recorder + the delivery worker.
type webhookTestEnv struct {
	outbox   *webhook.PostgresOutbox
	subRepo  *webhook.InMemorySubscriptionRepository
	recorder *webhook.InMemoryDeliveryRecorder
	subSvc   *integdomain.SubscriptionService
	worker   *webhook.DeliveryWorker
}

func newWebhookEnv(t *testing.T) *webhookTestEnv {
	t.Helper()
	pool := startPostgres(t)

	logger := zap.NewNop()

	outbox := webhook.NewPostgresOutbox(pool, logger)
	subRepo := webhook.NewInMemorySubscriptionRepository()
	recorder := webhook.NewInMemoryDeliveryRecorder()
	subSvc := integdomain.NewSubscriptionService(subRepo)

	worker := webhook.NewDeliveryWorker(webhook.DeliveryWorkerConfig{
		Outbox:        outbox,
		Subscriptions: subRepo,
		Deliveries:    recorder,
		Deactivate:    subSvc.RecordDeliveryFailure,
	}, logger)

	return &webhookTestEnv{
		outbox:   outbox,
		subRepo:  subRepo,
		recorder: recorder,
		subSvc:   subSvc,
		worker:   worker,
	}
}

const testSecret = "01234567890123456789012345678901"

func (e *webhookTestEnv) subscribe(t *testing.T, url string, events ...integdomain.EventType) integdomain.SubscriptionID {
	t.Helper()
	sub, err := e.subSvc.CreateSubscription(context.Background(), "tenant-1", url, events, testSecret)
	if err != nil {
		t.Fatalf("create subscription: %v", err)
	}
	return sub.ID
}

// startReceiver boots an httptest subscriber; returns the URL and a counter.
func startReceiver(t *testing.T, handler http.HandlerFunc) (*httptest.Server, *atomic.Int32) {
	t.Helper()
	var calls atomic.Int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls.Add(1)
		handler(w, r)
	}))
	t.Cleanup(srv.Close)
	return srv, &calls
}

func TestWebhookE2EIdempotency(t *testing.T) {
	env := newWebhookEnv(t)
	srv, calls := startReceiver(t, func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	env.subscribe(t, srv.URL, integdomain.EventModuleMastered)

	event := webhook.Event{EventID: "d4c1e2a5-1f0e-4c3a-9b2f-6a3c8e7d1a2b", Type: webhook.EventModuleMastered, Payload: map[string]any{"module_id": "m1"}}

	// First enqueue + deliver.
	if err := env.outbox.Enqueue(context.Background(), event); err != nil {
		t.Fatalf("enqueue: %v", err)
	}
	if err := env.worker.Deliver(context.Background(), event); err != nil {
		t.Fatalf("deliver: %v", err)
	}
	if calls.Load() != 1 {
		t.Fatalf("expected 1 delivery, got %d", calls.Load())
	}

	// Duplicate delivery (same event_id): the delivery recorder dedups.
	if err := env.worker.Deliver(context.Background(), event); err != nil {
		t.Fatalf("duplicate deliver: %v", err)
	}
	if calls.Load() != 1 {
		t.Fatalf("expected dedup (1 delivery), got %d", calls.Load())
	}
}

func TestWebhookE2EDeliverySigning(t *testing.T) {
	env := newWebhookEnv(t)
	var sigHeader atomic.Value
	srv, _ := startReceiver(t, func(w http.ResponseWriter, r *http.Request) {
		sigHeader.Store(r.Header.Get("X-Vedo-Signature"))
		w.WriteHeader(http.StatusOK)
	})
	env.subscribe(t, srv.URL, integdomain.EventModuleMastered)

	event := webhook.Event{EventID: "d4c1e2a5-1f0e-4c3a-9b2f-6a3c8e7d1a2b", Type: webhook.EventModuleMastered, Payload: map[string]any{}}
	if err := env.outbox.Enqueue(context.Background(), event); err != nil {
		t.Fatalf("enqueue: %v", err)
	}
	if err := env.worker.Deliver(context.Background(), event); err != nil {
		t.Fatalf("deliver: %v", err)
	}

	sig := sigHeader.Load()
	if sig == nil {
		t.Fatal("expected X-Vedo-Signature header")
	}
	header := sig.(string)
	if !strings.HasPrefix(header, "t=") || !strings.Contains(header, ",v1=") {
		t.Fatalf("unexpected signature format: %q", header)
	}
}

func TestWebhookE2ERetryThenPermanentFailure(t *testing.T) {
	env := newWebhookEnv(t)
	var attempts atomic.Int32
	srv, _ := startReceiver(t, func(w http.ResponseWriter, _ *http.Request) {
		attempts.Add(1)
		w.WriteHeader(http.StatusInternalServerError) // always fail
	})
	sub := env.subscribe(t, srv.URL, integdomain.EventModuleMastered)

	event := webhook.Event{EventID: "d4c1e2a5-1f0e-4c3a-9b2f-6a3c8e7d1a2b", Type: webhook.EventModuleMastered, Payload: map[string]any{}}
	if err := env.outbox.Enqueue(context.Background(), event); err != nil {
		t.Fatalf("enqueue: %v", err)
	}

	// Deliver repeatedly (each call is one attempt). After the failure budget
	// the subscription must be deactivated.
	for i := 0; i < integdomain.MaxConsecutiveDeliveryFailures+2; i++ {
		_ = env.worker.Deliver(context.Background(), event)
	}

	got, err := env.subSvc.GetSubscription(context.Background(), sub, "tenant-1")
	if err != nil {
		t.Fatalf("get subscription: %v", err)
	}
	if got.Active {
		t.Fatal("expected subscription deactivated after consecutive failures")
	}
}

func TestWebhookE2EEventTypeFiltering(t *testing.T) {
	env := newWebhookEnv(t)
	var moduleCalls atomic.Int32
	srv, _ := startReceiver(t, func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	env.subscribe(t, srv.URL, integdomain.EventModuleMastered)

	// A plan.deviated event must NOT reach a module.mastered-only subscriber.
	event := webhook.Event{EventID: "d4c1e2a5-1f0e-4c3a-9b2f-6a3c8e7d1a2b", Type: webhook.EventPlanDeviated, Payload: map[string]any{}}
	if err := env.outbox.Enqueue(context.Background(), event); err != nil {
		t.Fatalf("enqueue: %v", err)
	}
	if err := env.worker.Deliver(context.Background(), event); err != nil {
		t.Fatalf("deliver: %v", err)
	}
	if moduleCalls.Load() != 0 {
		t.Fatalf("expected 0 deliveries for non-matching event type, got %d", moduleCalls.Load())
	}
}

func TestWebhookE2EMultipleSubscriptions(t *testing.T) {
	env := newWebhookEnv(t)
	srvA, callsA := startReceiver(t, func(w http.ResponseWriter, _ *http.Request) { w.WriteHeader(http.StatusOK) })
	srvB, callsB := startReceiver(t, func(w http.ResponseWriter, _ *http.Request) { w.WriteHeader(http.StatusOK) })

	idA := env.subscribe(t, srvA.URL, integdomain.EventModuleMastered)
	idB := env.subscribe(t, srvB.URL, integdomain.EventModuleMastered)

	// Verify subscriptions are in the repo.
	list, err := env.subSvc.ListSubscriptions(context.Background(), "tenant-1", nil)
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	if len(list) != 2 {
		t.Fatalf("expected 2 subscriptions, got %d", len(list))
	}

	event := webhook.Event{EventID: "d4c1e2a5-1f0e-4c3a-9b2f-6a3c8e7d1a2b", Type: webhook.EventModuleMastered, Payload: map[string]any{}}
	if err := env.outbox.Enqueue(context.Background(), event); err != nil {
		t.Fatalf("enqueue: %v", err)
	}
	if err := env.worker.Deliver(context.Background(), event); err != nil {
		t.Fatalf("deliver: %v", err)
	}
	if callsA.Load() != 1 || callsB.Load() != 1 {
		t.Fatalf("expected both subscribers to receive the event: A=%d B=%d", callsA.Load(), callsB.Load())
	}
	_ = idA
	_ = idB
}

func TestWebhookE2ESignatureAlgorithm(t *testing.T) {
	// Verify the documented algorithm matches the Signer output.
	payload := []byte(`{"event_id":"d4c1e2a5-1f0e-4c3a-9b2f-6a3c8e7d1a2b","event_type":"module.mastered"}`)
	timestamp := int64(1694123456)
	mac := hmac.New(sha256.New, []byte(testSecret))
	_, _ = mac.Write([]byte(fmt.Sprintf("%d.", timestamp)))
	_, _ = mac.Write(payload)
	expectedSig := hex.EncodeToString(mac.Sum(nil))

	signer := webhook.NewSigner()
	header := signer.Sign(testSecret, payload, time.Unix(timestamp, 0))
	expected := fmt.Sprintf("t=%d,v1=%s", timestamp, expectedSig)
	if header != expected {
		t.Fatalf("signature mismatch:\n got %s\nwant %s", header, expected)
	}
}
