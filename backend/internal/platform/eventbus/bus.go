// Package eventbus implements the in-process event bus (ADR-DES.API.
// communication-patterns: sync REST + async events via in-process bus +
// outbox + webhooks).
//
// The bus is a minimal publish/subscribe channel for domain events:
//   - Publish dispatches to subscribers asynchronously (non-blocking, the
//     publisher never waits for subscriber side effects);
//   - Subscribe registers a typed handler for an event name and returns an
//     unsubscribe function;
//   - Subscribers are isolated: a panic in one subscriber is recovered and
//     logged, it does not affect the publisher or other subscribers.
//
// Delivery is at-most-once within the process (persistence/idempotency lives
// in the outbox — see the integrations module).
package eventbus

import (
	"sync"

	"go.uber.org/zap"
)

// Event is a domain event envelope.
type Event struct {
	// Name is the stable event identifier (e.g. "ModuleMastered").
	Name string
	// Payload is the event data (typed struct by convention).
	Payload any
}

// Handler processes one event.
type Handler func(event Event)

// subscription is one registered handler with a unique id.
type subscription struct {
	id uint64
	h  Handler
}

// Bus is a concurrency-safe in-process publish/subscribe channel.
type Bus struct {
	mu       sync.RWMutex
	nextID   uint64
	handlers map[string][]subscription
	logger   *zap.Logger
}

// New creates an empty bus.
func New(logger *zap.Logger) *Bus {
	if logger == nil {
		logger = zap.NewNop()
	}
	return &Bus{handlers: map[string][]subscription{}, logger: logger.Named("eventbus")}
}

// Subscribe registers a handler for an event name. Returns an unsubscribe
// function.
func (b *Bus) Subscribe(name string, h Handler) func() {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.nextID++
	id := b.nextID
	b.handlers[name] = append(b.handlers[name], subscription{id: id, h: h})

	return func() {
		b.mu.Lock()
		defer b.mu.Unlock()
		handlers := b.handlers[name]
		for i, sub := range handlers {
			if sub.id == id {
				b.handlers[name] = append(handlers[:i], handlers[i+1:]...)
				break
			}
		}
		if len(b.handlers[name]) == 0 {
			delete(b.handlers, name)
		}
	}
}

// Publish dispatches the event to all registered handlers asynchronously.
// The call returns immediately; subscriber side effects run in goroutines.
func (b *Bus) Publish(event Event) {
	b.mu.RLock()
	handlers := make([]Handler, 0, len(b.handlers[event.Name]))
	for _, sub := range b.handlers[event.Name] {
		handlers = append(handlers, sub.h)
	}
	b.mu.RUnlock()

	for _, h := range handlers {
		go b.safeDispatch(h, event)
	}
}

// safeDispatch runs a handler with panic recovery so one broken subscriber
// cannot take down the publisher or other subscribers.
func (b *Bus) safeDispatch(h Handler, event Event) {
	defer func() {
		if rec := recover(); rec != nil {
			b.logger.Error("event subscriber panicked",
				zap.String("event", event.Name),
				zap.Any("panic", rec),
			)
		}
	}()
	h(event)
}
