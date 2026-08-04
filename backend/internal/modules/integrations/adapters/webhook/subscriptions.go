package webhook

import (
	"context"
	"fmt"
	"sort"
	"sync"

	domain "vedo-edutrack/backend/internal/modules/integrations/domain"
)

// InMemorySubscriptionRepository is an in-memory SubscriptionRepository used
// by tests and the M4 handler wiring until the PostgreSQL-backed repository
// lands (Task 8+). It is concurrency-safe.
type InMemorySubscriptionRepository struct {
	mu   sync.Mutex
	subs map[domain.SubscriptionID]*domain.WebhookSubscription
}

// NewInMemorySubscriptionRepository builds an empty in-memory repository.
func NewInMemorySubscriptionRepository() *InMemorySubscriptionRepository {
	return &InMemorySubscriptionRepository{subs: map[domain.SubscriptionID]*domain.WebhookSubscription{}}
}

// Create stores a subscription.
func (r *InMemorySubscriptionRepository) Create(_ context.Context, sub *domain.WebhookSubscription) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, exists := r.subs[sub.ID]; exists {
		return fmt.Errorf("subscription %s already exists", sub.ID)
	}
	r.subs[sub.ID] = sub
	return nil
}

// Update replaces a stored subscription.
func (r *InMemorySubscriptionRepository) Update(_ context.Context, sub *domain.WebhookSubscription) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, exists := r.subs[sub.ID]; !exists {
		return domain.ErrNotFound
	}
	r.subs[sub.ID] = sub
	return nil
}

// Delete removes a subscription.
func (r *InMemorySubscriptionRepository) Delete(_ context.Context, id domain.SubscriptionID) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, exists := r.subs[id]; !exists {
		return domain.ErrNotFound
	}
	delete(r.subs, id)
	return nil
}

// Get returns one subscription (ErrNotFound when absent).
func (r *InMemorySubscriptionRepository) Get(_ context.Context, id domain.SubscriptionID) (*domain.WebhookSubscription, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	sub, ok := r.subs[id]
	if !ok {
		return nil, domain.ErrNotFound
	}
	copy := *sub
	return &copy, nil
}

// ListByTenant returns the tenant's subscriptions, optionally active-filtered.
func (r *InMemorySubscriptionRepository) ListByTenant(_ context.Context, tenantID string, active *bool) ([]domain.WebhookSubscription, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	out := make([]domain.WebhookSubscription, 0, len(r.subs))
	for _, s := range r.subs {
		if s.TenantID != tenantID {
			continue
		}
		if active != nil && s.Active != *active {
			continue
		}
		out = append(out, *s)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].CreatedAt.Before(out[j].CreatedAt) })
	return out, nil
}

// ListByEventType returns active subscriptions subscribed to the event type.
func (r *InMemorySubscriptionRepository) ListByEventType(_ context.Context, eventType domain.EventType) ([]domain.WebhookSubscription, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	out := make([]domain.WebhookSubscription, 0, len(r.subs))
	for _, s := range r.subs {
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
	sort.Slice(out, func(i, j int) bool { return out[i].CreatedAt.Before(out[j].CreatedAt) })
	return out, nil
}

// CountActiveByTenant counts the tenant's active subscriptions.
func (r *InMemorySubscriptionRepository) CountActiveByTenant(_ context.Context, tenantID string) (int, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	n := 0
	for _, s := range r.subs {
		if s.TenantID == tenantID && s.Active {
			n++
		}
	}
	return n, nil
}
