package webhook

import (
	"testing"
)

func TestOutboxDedupByIdempotencyKey(t *testing.T) {
	outbox := NewOutbox()
	event := Event{EventID: "d4c1e2a5-1f0e-4c3a-9b2f-6a3c8e7d1a2b", Type: "module.mastered", Payload: map[string]any{"module_id": "m1"}}

	ok, duplicate, err := outbox.Enqueue(event)
	if err != nil || !ok || duplicate {
		t.Fatalf("first enqueue: ok=%t duplicate=%t err=%v", ok, duplicate, err)
	}
	ok, duplicate, err = outbox.Enqueue(event)
	if err != nil || ok || !duplicate {
		t.Fatalf("duplicate enqueue: ok=%t duplicate=%t err=%v (want ok=false, duplicate=true)", ok, duplicate, err)
	}
	if outbox.Len() != 1 {
		t.Fatalf("expected 1 event in outbox, got %d", outbox.Len())
	}
}

func TestOutboxRejectsInvalidEventID(t *testing.T) {
	outbox := NewOutbox()
	_, _, err := outbox.Enqueue(Event{EventID: "not-a-uuid", Type: "module.mastered"})
	if err == nil {
		t.Fatal("expected invalid event id rejected")
	}
}

func TestOutboxEventTypes(t *testing.T) {
	outbox := NewOutbox()
	types := []EventType{EventModuleMastered, EventPlanDeviated, EventRouteRecalculated}
	for _, eventType := range types {
		ok, _, err := outbox.Enqueue(Event{EventID: newUUID(), Type: eventType, Payload: map[string]any{}})
		if err != nil || !ok {
			t.Fatalf("enqueue %s: ok=%t err=%v", eventType, ok, err)
		}
	}
	if outbox.Len() != 3 {
		t.Fatalf("expected 3 events, got %d", outbox.Len())
	}
}
