// Package webhook implements the webhook outbox infrastructure (F6).
//
// Outbound webhook events (module.mastered, plan.deviated,
// route.recalculated) are enqueued into an outbox and delivered with
// idempotency guarantees: a duplicate event_id is acknowledged without
// re-dispatch (REQ-FR-api.webhooks.idempotency).
package webhook

import (
	"fmt"
	"sync"

	"github.com/google/uuid"
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
)

// Event is an outbound webhook payload.
type Event struct {
	EventID string
	Type    EventType
	Payload map[string]any
}

// Outbox is an in-memory idempotent webhook outbox.
type Outbox struct {
	mu     sync.Mutex
	events map[string]Event
	order  []string
}

// NewOutbox creates an empty outbox.
func NewOutbox() *Outbox {
	return &Outbox{events: map[string]Event{}}
}

// Enqueue adds an event. Returns (created, duplicate, error).
func (o *Outbox) Enqueue(event Event) (bool, bool, error) {
	if event.EventID == "" {
		return false, false, fmt.Errorf("event id is required")
	}
	if _, err := uuid.Parse(event.EventID); err != nil {
		return false, false, fmt.Errorf("event id %q is not a UUID: %w", event.EventID, err)
	}
	if event.Type == "" {
		return false, false, fmt.Errorf("event type is required")
	}
	o.mu.Lock()
	defer o.mu.Unlock()
	if _, exists := o.events[event.EventID]; exists {
		return false, true, nil
	}
	o.events[event.EventID] = event
	o.order = append(o.order, event.EventID)
	return true, false, nil
}

// Len returns the number of unique enqueued events.
func (o *Outbox) Len() int {
	o.mu.Lock()
	defer o.mu.Unlock()
	return len(o.events)
}

// newUUID returns a random v4 UUID (test helper).
func newUUID() string {
	return uuid.NewString()
}
