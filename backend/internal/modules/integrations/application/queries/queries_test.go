package queries

import (
	"context"
	"strings"
	"testing"

	"go.uber.org/zap"

	domain "vedo-edutrack/backend/internal/modules/integrations/domain"
)

type fakeRepo struct {
	subs map[domain.SubscriptionID]*domain.WebhookSubscription
}

func (f *fakeRepo) Create(_ context.Context, sub *domain.WebhookSubscription) error {
	f.subs[sub.ID] = sub
	return nil
}

func (f *fakeRepo) Update(_ context.Context, sub *domain.WebhookSubscription) error {
	f.subs[sub.ID] = sub
	return nil
}

func (f *fakeRepo) Delete(_ context.Context, id domain.SubscriptionID) error {
	delete(f.subs, id)
	return nil
}

func (f *fakeRepo) Get(_ context.Context, id domain.SubscriptionID) (*domain.WebhookSubscription, error) {
	if s, ok := f.subs[id]; ok {
		return s, nil
	}
	return nil, domain.ErrNotFound
}

func (f *fakeRepo) ListByTenant(_ context.Context, tenantID string, active *bool) ([]domain.WebhookSubscription, error) {
	var out []domain.WebhookSubscription
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

func (f *fakeRepo) ListByEventType(_ context.Context, eventType domain.EventType) ([]domain.WebhookSubscription, error) {
	var out []domain.WebhookSubscription
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

func (f *fakeRepo) CountActiveByTenant(_ context.Context, tenantID string) (int, error) {
	n := 0
	for _, s := range f.subs {
		if s.TenantID == tenantID && s.Active {
			n++
		}
	}
	return n, nil
}

type fakeDeliveries struct {
	rows []domain.WebhookDelivery
}

func (f *fakeDeliveries) ListBySubscription(_ context.Context, _ domain.SubscriptionID, _ int, _ int) ([]domain.WebhookDelivery, error) {
	return f.rows, nil
}

func TestListWebhookSubscriptions(t *testing.T) {
	repo := &fakeRepo{subs: map[domain.SubscriptionID]*domain.WebhookSubscription{}}
	svc := domain.NewSubscriptionService(repo)
	sub, err := svc.CreateSubscription(context.Background(), "t1", "https://example.test/hook",
		[]domain.EventType{domain.EventModuleMastered}, strings.Repeat("s", domain.MinSigningSecretLength))
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	_ = sub

	q := NewSubscriptionQueries(svc, nil, zap.NewNop())
	active := true
	list, err := q.ListWebhookSubscriptions(context.Background(), "t1", &active)
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	if len(list) != 1 {
		t.Fatalf("len=%d, want 1", len(list))
	}
}

func TestGetWebhookDeliveryHistoryScopesByTenant(t *testing.T) {
	repo := &fakeRepo{subs: map[domain.SubscriptionID]*domain.WebhookSubscription{}}
	svc := domain.NewSubscriptionService(repo)
	sub, err := svc.CreateSubscription(context.Background(), "t1", "https://example.test/hook",
		[]domain.EventType{domain.EventModuleMastered}, strings.Repeat("s", domain.MinSigningSecretLength))
	if err != nil {
		t.Fatalf("create: %v", err)
	}

	deliveries := &fakeDeliveries{rows: []domain.WebhookDelivery{
		{ID: "d1", SubscriptionID: sub.ID, EventID: "e1", Attempt: 1, Status: domain.DeliverySent},
	}}
	q := NewSubscriptionQueries(svc, deliveries, zap.NewNop())

	// Same tenant → history returned.
	history, err := q.GetWebhookDeliveryHistory(context.Background(), sub.ID, "t1", 0, 10)
	if err != nil {
		t.Fatalf("history: %v", err)
	}
	if len(history) != 1 {
		t.Fatalf("history len=%d, want 1", len(history))
	}

	// Wrong tenant → scoped out with error.
	if _, err := q.GetWebhookDeliveryHistory(context.Background(), sub.ID, "other-tenant", 0, 10); err == nil {
		t.Fatal("expected cross-tenant history denied")
	}
}
