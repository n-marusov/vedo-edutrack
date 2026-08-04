package commands

import (
	"context"
	"errors"
	"strings"
	"testing"

	"go.uber.org/zap"

	domain "vedo-edutrack/backend/internal/modules/integrations/domain"
)

// fakeRepo implements domain.SubscriptionRepository in-memory.
type fakeRepo struct {
	subs map[domain.SubscriptionID]*domain.WebhookSubscription
}

func newFakeRepo() *fakeRepo {
	return &fakeRepo{subs: map[domain.SubscriptionID]*domain.WebhookSubscription{}}
}

func (f *fakeRepo) Create(_ context.Context, sub *domain.WebhookSubscription) error {
	f.subs[sub.ID] = sub
	return nil
}

func (f *fakeRepo) Update(_ context.Context, sub *domain.WebhookSubscription) error {
	if _, ok := f.subs[sub.ID]; !ok {
		return errors.New("not found")
	}
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

// fakeOutbox records enqueued events.
type fakeOutbox struct {
	events []OutboxEvent
}

func (f *fakeOutbox) Enqueue(_ context.Context, event OutboxEvent) error {
	f.events = append(f.events, event)
	return nil
}

func secret() string { return strings.Repeat("s", domain.MinSigningSecretLength) }

func newCommands() (*SubscriptionCommands, *fakeRepo, *fakeOutbox) {
	repo := newFakeRepo()
	outbox := &fakeOutbox{}
	svc := domain.NewSubscriptionService(repo)
	return NewSubscriptionCommands(svc, outbox, zap.NewNop()), repo, outbox
}

func TestCreateWebhookSubscription(t *testing.T) {
	cmds, _, outbox := newCommands()
	sub, err := cmds.CreateWebhookSubscription(context.Background(), CreateWebhookSubscriptionInput{
		TenantID:      "t1",
		URL:           "https://example.test/hook",
		EventTypes:    []domain.EventType{domain.EventModuleMastered},
		SigningSecret: secret(),
	})
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	if sub.ID == "" || !sub.Active {
		t.Fatalf("unexpected subscription: %+v", sub)
	}
	// Verification ping enqueued.
	if len(outbox.events) != 1 || outbox.events[0].Type != "webhook.ping" {
		t.Fatalf("expected verification ping, got %+v", outbox.events)
	}
}

func TestCreateWebhookSubscriptionRejectsInvalidURL(t *testing.T) {
	cmds, _, _ := newCommands()
	_, err := cmds.CreateWebhookSubscription(context.Background(), CreateWebhookSubscriptionInput{
		TenantID:      "t1",
		URL:           "http://insecure.test/hook",
		EventTypes:    []domain.EventType{domain.EventModuleMastered},
		SigningSecret: secret(),
	})
	if err == nil {
		t.Fatal("expected invalid URL rejected")
	}
}

func TestUpdateWebhookSubscription(t *testing.T) {
	cmds, _, _ := newCommands()
	sub, err := cmds.CreateWebhookSubscription(context.Background(), CreateWebhookSubscriptionInput{
		TenantID: "t1", URL: "https://example.test/hook",
		EventTypes: []domain.EventType{domain.EventModuleMastered}, SigningSecret: secret(),
	})
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	updated, err := cmds.UpdateWebhookSubscription(context.Background(), UpdateWebhookSubscriptionInput{
		SubscriptionID: sub.ID, TenantID: "t1",
		URL: "https://example.test/v2", EventTypes: []domain.EventType{domain.EventPlanDeviated},
	})
	if err != nil {
		t.Fatalf("update: %v", err)
	}
	if updated.URL != "https://example.test/v2" || len(updated.EventTypes) != 1 || updated.EventTypes[0] != domain.EventPlanDeviated {
		t.Fatalf("update not applied: %+v", updated)
	}
}

func TestDeleteWebhookSubscription(t *testing.T) {
	cmds, repo, _ := newCommands()
	sub, err := cmds.CreateWebhookSubscription(context.Background(), CreateWebhookSubscriptionInput{
		TenantID: "t1", URL: "https://example.test/hook",
		EventTypes: []domain.EventType{domain.EventModuleMastered}, SigningSecret: secret(),
	})
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	if err := cmds.DeleteWebhookSubscription(context.Background(), DeleteWebhookSubscriptionInput{SubscriptionID: sub.ID, TenantID: "t1"}); err != nil {
		t.Fatalf("delete: %v", err)
	}
	if _, err := repo.Get(context.Background(), sub.ID); err == nil {
		t.Fatal("expected subscription removed after delete")
	}
}
