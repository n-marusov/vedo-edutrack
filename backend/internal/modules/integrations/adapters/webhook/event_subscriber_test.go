package webhook

import (
	"context"
	"strings"
	"sync"
	"testing"
	"time"

	"go.uber.org/zap"

	"vedo-edutrack/backend/internal/platform/eventbus"
)

// captureOutbox records enqueued events for assertions.
type captureOutbox struct {
	mu     sync.Mutex
	events []Event
}

func (c *captureOutbox) Enqueue(_ context.Context, event Event) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.events = append(c.events, event)
	return nil
}

func (c *captureOutbox) Dequeue(context.Context, int) ([]Event, error) {
	return nil, nil
}

func (c *captureOutbox) MarkDelivered(context.Context, string) error { return nil }

func (c *captureOutbox) MarkFailed(context.Context, string, string) error { return nil }

func (c *captureOutbox) list() []Event {
	c.mu.Lock()
	defer c.mu.Unlock()
	return append([]Event{}, c.events...)
}

func TestEventSubscriberMapsModuleMastered(t *testing.T) {
	bus := eventbus.New(zap.NewNop())
	outbox := &captureOutbox{}
	sub := NewEventSubscriber(bus, outbox, zap.NewNop())
	sub.Register()

	bus.Publish(eventbus.Event{
		Name: EventNameModuleMastered,
		Payload: map[string]any{
			"event_id":   "11111111-1111-1111-1111-111111111111",
			"learner_id": "l1",
			"module_id":  "m1",
		},
	})

	events := eventuallyEvents(t, outbox, 1)
	if events[0].Type != EventModuleMastered {
		t.Fatalf("type = %s, want module.mastered", events[0].Type)
	}
	if events[0].EventID != "11111111-1111-1111-1111-111111111111" {
		t.Fatalf("event_id = %s, want the payload event_id (idempotency)", events[0].EventID)
	}
}

func TestEventSubscriberMapsAllEventTypes(t *testing.T) {
	bus := eventbus.New(zap.NewNop())
	outbox := &captureOutbox{}
	sub := NewEventSubscriber(bus, outbox, zap.NewNop())
	sub.Register()

	cases := []struct {
		name     string
		wantType EventType
	}{
		{EventNameModuleMastered, EventModuleMastered},
		{EventNamePlanDeviationDetected, EventPlanDeviated},
		{EventNameRouteRecalculated, EventRouteRecalculated},
		{EventNameStandardDeficitDetected, EventStandardRiskDetected},
	}
	for i, c := range cases {
		bus.Publish(eventbus.Event{Name: c.name, Payload: map[string]any{"n": i}})
	}

	events := eventuallyEvents(t, outbox, len(cases))
	got := map[EventType]bool{}
	for _, e := range events {
		got[e.Type] = true
	}
	for _, c := range cases {
		if !got[c.wantType] {
			t.Fatalf("missing event type %s; got %v", c.wantType, got)
		}
	}
}

func TestEventSubscriberGeneratesIDWhenMissing(t *testing.T) {
	bus := eventbus.New(zap.NewNop())
	outbox := &captureOutbox{}
	sub := NewEventSubscriber(bus, outbox, zap.NewNop())
	sub.Register()

	bus.Publish(eventbus.Event{Name: EventNameModuleMastered, Payload: map[string]any{"module_id": "m1"}})

	events := eventuallyEvents(t, outbox, 1)
	if events[0].EventID == "" {
		t.Fatal("expected generated event_id")
	}
	if !strings.Contains(events[0].EventID, "-") {
		t.Fatalf("event_id = %q, want UUID format", events[0].EventID)
	}
}

func eventuallyEvents(t *testing.T, outbox *captureOutbox, want int) []Event {
	t.Helper()
	var events []Event
	for i := 0; i < 100; i++ {
		events = outbox.list()
		if len(events) >= want {
			return events
		}
		time.Sleep(10 * time.Millisecond)
	}
	t.Fatalf("expected %d events, got %d", want, len(events))
	return nil
}
