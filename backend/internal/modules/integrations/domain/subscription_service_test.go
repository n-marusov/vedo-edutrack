package integrations

import (
	"context"
	"errors"
	"strings"
	"testing"
)

// fakeSubscriptionRepo is an in-memory SubscriptionRepository for domain tests.
type fakeSubscriptionRepo struct {
	subs map[SubscriptionID]*WebhookSubscription
}

func newFakeRepo() *fakeSubscriptionRepo {
	return &fakeSubscriptionRepo{subs: map[SubscriptionID]*WebhookSubscription{}}
}

func (f *fakeSubscriptionRepo) Create(_ context.Context, sub *WebhookSubscription) error {
	f.subs[sub.ID] = sub
	return nil
}

func (f *fakeSubscriptionRepo) Update(_ context.Context, sub *WebhookSubscription) error {
	if _, ok := f.subs[sub.ID]; !ok {
		return errors.New("not found")
	}
	f.subs[sub.ID] = sub
	return nil
}

func (f *fakeSubscriptionRepo) Delete(_ context.Context, id SubscriptionID) error {
	if _, ok := f.subs[id]; !ok {
		return errors.New("not found")
	}
	delete(f.subs, id)
	return nil
}

func (f *fakeSubscriptionRepo) Get(_ context.Context, id SubscriptionID) (*WebhookSubscription, error) {
	sub, ok := f.subs[id]
	if !ok {
		return nil, errors.New("not found")
	}
	return sub, nil
}

func (f *fakeSubscriptionRepo) ListByTenant(_ context.Context, tenantID string, active *bool) ([]WebhookSubscription, error) {
	var out []WebhookSubscription
	for _, s := range f.subs {
		if s.TenantID != tenantID {
			continue
		}
		if active != nil && s.Active != *active {
			continue
		}
		out = append(out, *s)
	}
	return out, nil
}

func (f *fakeSubscriptionRepo) ListByEventType(_ context.Context, eventType EventType) ([]WebhookSubscription, error) {
	var out []WebhookSubscription
	for _, s := range f.subs {
		if !s.Active {
			continue
		}
		for _, t := range s.EventTypes {
			if t == eventType {
				out = append(out, *s)
				break
			}
		}
	}
	return out, nil
}

func (f *fakeSubscriptionRepo) CountActiveByTenant(_ context.Context, tenantID string) (int, error) {
	count := 0
	for _, s := range f.subs {
		if s.TenantID == tenantID && s.Active {
			count++
		}
	}
	return count, nil
}

func testSecret() string { return strings.Repeat("s", MinSigningSecretLength) }

func TestSubscriptionServiceCreate(t *testing.T) {
	svc := NewSubscriptionService(newFakeRepo())
	sub, err := svc.CreateSubscription(context.Background(), "tenant-1", "https://example.test/hook", []EventType{EventModuleMastered}, testSecret())
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	if sub.ID == "" || !sub.Active || sub.Failures != 0 {
		t.Fatalf("unexpected subscription: %+v", sub)
	}
	if !sub.CreatedAt.Equal(sub.UpdatedAt) {
		t.Fatal("expected created == updated on fresh subscription")
	}
}

func TestSubscriptionServiceCreateRejectsInvalid(t *testing.T) {
	svc := NewSubscriptionService(newFakeRepo())
	if _, err := svc.CreateSubscription(context.Background(), "", "https://example.test/hook", []EventType{EventModuleMastered}, testSecret()); err == nil {
		t.Fatal("expected tenant validation error")
	}
	if _, err := svc.CreateSubscription(context.Background(), "t1", "http://insecure.test/hook", []EventType{EventModuleMastered}, testSecret()); err == nil {
		t.Fatal("expected URL validation error")
	}
	if _, err := svc.CreateSubscription(context.Background(), "t1", "https://example.test/hook", []EventType{"bogus"}, testSecret()); err == nil {
		t.Fatal("expected event type validation error")
	}
}

func TestSubscriptionServiceEnforcesTenantCap(t *testing.T) {
	repo := newFakeRepo()
	svc := NewSubscriptionService(repo)
	ctx := context.Background()

	for i := 0; i < MaxSubscriptionsPerTenant; i++ {
		if _, err := svc.CreateSubscription(ctx, "tenant-1", "https://example.test/hook", []EventType{EventModuleMastered}, testSecret()); err != nil {
			t.Fatalf("create %d: %v", i, err)
		}
	}
	if _, err := svc.CreateSubscription(ctx, "tenant-1", "https://example.test/hook", []EventType{EventModuleMastered}, testSecret()); err == nil {
		t.Fatal("expected cap exceeded error")
	} else if !strings.Contains(err.Error(), "maximum") {
		t.Fatalf("expected cap message, got %v", err)
	}
}

func TestSubscriptionServiceUpdateResetsFailures(t *testing.T) {
	svc := NewSubscriptionService(newFakeRepo())
	sub, err := svc.CreateSubscription(context.Background(), "t1", "https://example.test/hook", []EventType{EventModuleMastered}, testSecret())
	if err != nil {
		t.Fatalf("create: %v", err)
	}

	// Simulate one failure.
	if _, err := svc.RecordDeliveryFailure(context.Background(), sub.ID); err != nil {
		t.Fatalf("record failure: %v", err)
	}
	if sub.Failures != 1 {
		t.Fatalf("failures=%d, want 1", sub.Failures)
	}

	updated, err := svc.UpdateSubscription(context.Background(), sub.ID, "t1", "https://example.test/v2", []EventType{EventPlanDeviated}, "")
	if err != nil {
		t.Fatalf("update: %v", err)
	}
	if updated.URL != "https://example.test/v2" || updated.Failures != 0 {
		t.Fatalf("update did not apply/reset: %+v", updated)
	}
	if len(updated.EventTypes) != 1 || updated.EventTypes[0] != EventPlanDeviated {
		t.Fatalf("event types not updated: %+v", updated.EventTypes)
	}
}

func TestSubscriptionServiceDeactivatesAfterConsecutiveFailures(t *testing.T) {
	svc := NewSubscriptionService(newFakeRepo())
	sub, err := svc.CreateSubscription(context.Background(), "t1", "https://example.test/hook", []EventType{EventModuleMastered}, testSecret())
	if err != nil {
		t.Fatalf("create: %v", err)
	}

	var deactivated bool
	for i := 0; i < MaxConsecutiveDeliveryFailures; i++ {
		deactivated, err = svc.RecordDeliveryFailure(context.Background(), sub.ID)
		if err != nil {
			t.Fatalf("record failure %d: %v", i, err)
		}
	}
	if !deactivated {
		t.Fatal("expected deactivation at the failure boundary")
	}
	if sub.Active {
		t.Fatal("expected subscription inactive after boundary")
	}
}

func TestSubscriptionServiceTenantIsolation(t *testing.T) {
	svc := NewSubscriptionService(newFakeRepo())
	sub, err := svc.CreateSubscription(context.Background(), "tenant-a", "https://example.test/hook", []EventType{EventModuleMastered}, testSecret())
	if err != nil {
		t.Fatalf("create: %v", err)
	}

	if _, err := svc.GetSubscription(context.Background(), sub.ID, "tenant-b"); err == nil {
		t.Fatal("expected cross-tenant get denied")
	}
	if err := svc.DeleteSubscription(context.Background(), sub.ID, "tenant-b"); err == nil {
		t.Fatal("expected cross-tenant delete denied")
	}
}

func TestSubscriptionServiceListFilter(t *testing.T) {
	svc := NewSubscriptionService(newFakeRepo())
	for i := 0; i < 3; i++ {
		if _, err := svc.CreateSubscription(context.Background(), "t1", "https://example.test/hook", []EventType{EventModuleMastered}, testSecret()); err != nil {
			t.Fatalf("create %d: %v", i, err)
		}
	}
	active := true
	all, err := svc.ListSubscriptions(context.Background(), "t1", nil)
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	if len(all) != 3 {
		t.Fatalf("list len=%d, want 3", len(all))
	}
	activeOnly, err := svc.ListSubscriptions(context.Background(), "t1", &active)
	if err != nil {
		t.Fatalf("list active: %v", err)
	}
	if len(activeOnly) != 3 {
		t.Fatalf("active len=%d, want 3", len(activeOnly))
	}
}
