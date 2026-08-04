package eventbus

import (
	"sync/atomic"
	"testing"
	"time"

	"go.uber.org/zap"
)

func TestPublishDeliversToSubscriber(t *testing.T) {
	bus := New(zap.NewNop())
	var got atomic.Value
	bus.Subscribe("TestEvent", func(event Event) {
		got.Store(event.Payload)
	})

	bus.Publish(Event{Name: "TestEvent", Payload: "hello"})
	requireEventually(t, func() bool { return got.Load() != nil }, "subscriber should receive the event")
	if got.Load() != "hello" {
		t.Fatalf("payload = %v, want hello", got.Load())
	}
}

func TestPublishFiltersByEventName(t *testing.T) {
	bus := New(zap.NewNop())
	var calls atomic.Int32
	bus.Subscribe("EventA", func(Event) { calls.Add(1) })

	bus.Publish(Event{Name: "EventA"})
	bus.Publish(Event{Name: "EventB"}) // should not reach the subscriber
	requireEventually(t, func() bool { return calls.Load() == 1 }, "only EventA should be delivered")
}

func TestUnsubscribeStopsDelivery(t *testing.T) {
	bus := New(zap.NewNop())
	var calls atomic.Int32
	unsub := bus.Subscribe("EventA", func(Event) { calls.Add(1) })

	bus.Publish(Event{Name: "EventA"})
	requireEventually(t, func() bool { return calls.Load() == 1 }, "initial delivery")

	unsub()
	bus.Publish(Event{Name: "EventA"})
	time.Sleep(50 * time.Millisecond)
	if calls.Load() != 1 {
		t.Fatalf("expected no delivery after unsubscribe, got %d", calls.Load())
	}
}

func TestSubscriberPanicDoesNotKillOthers(t *testing.T) {
	bus := New(zap.NewNop())
	var ok atomic.Int32
	bus.Subscribe("PanicEvent", func(Event) { panic("boom") })
	bus.Subscribe("PanicEvent", func(Event) { ok.Add(1) })

	bus.Publish(Event{Name: "PanicEvent"})
	requireEventually(t, func() bool { return ok.Load() == 1 }, "second subscriber must still run")
}

func requireEventually(t *testing.T, cond func() bool, msg string) {
	t.Helper()
	deadline := time.Now().Add(2 * time.Second)
	for time.Now().Before(deadline) {
		if cond() {
			return
		}
		time.Sleep(10 * time.Millisecond)
	}
	t.Fatal(msg)
}
