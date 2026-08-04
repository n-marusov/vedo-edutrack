package webhook

import (
	"context"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"vedo-edutrack/backend/internal/platform/eventbus"
)

// Domain event names (specs/ddd/domain-events.md §Events). These are the
// in-process core events published by the bounded contexts; the subscriber
// maps them to outbox webhook events.
const (
	EventNameModuleMastered          = "ModuleMastered"
	EventNamePlanDeviationDetected   = "PlanDeviationDetected"
	EventNameRouteRecalculated       = "RouteRecalculated"
	EventNameStandardDeficitDetected = "StandardDeficitDetected"
)

// EventSubscriber bridges the in-process event bus to the webhook outbox:
// every core domain event is mapped to its webhook representation and
// enqueued so subscribers receive it (specs/ddd/domain-events.md §Webhook).
type EventSubscriber struct {
	bus    *eventbus.Bus
	outbox OutboxRepository
	logger *zap.Logger
}

// NewEventSubscriber builds the bus→outbox bridge.
func NewEventSubscriber(bus *eventbus.Bus, outbox OutboxRepository, logger *zap.Logger) *EventSubscriber {
	if logger == nil {
		logger = zap.NewNop()
	}
	return &EventSubscriber{bus: bus, outbox: outbox, logger: logger.Named("webhook.events")}
}

// Register subscribes the subscriber to all mapped domain events.
func (s *EventSubscriber) Register() {
	s.bus.Subscribe(EventNameModuleMastered, s.handle(EventModuleMastered))
	s.bus.Subscribe(EventNamePlanDeviationDetected, s.handle(EventPlanDeviated))
	s.bus.Subscribe(EventNameRouteRecalculated, s.handle(EventRouteRecalculated))
	s.bus.Subscribe(EventNameStandardDeficitDetected, s.handle(EventStandardRiskDetected))
}

// handle returns a bus handler that enqueues the event into the outbox.
// The event_id is taken from the payload's event_id field when present (so a
// re-published event reuses its id and the outbox dedups it); otherwise a
// fresh UUID is generated.
func (s *EventSubscriber) handle(outboxType EventType) func(eventbus.Event) {
	return func(event eventbus.Event) {
		s.logger.Debug("EventReceived", zap.String("eventType", event.Name))
		payload, err := mapEventPayload(event)
		if err != nil {
			s.logger.Warn("EventMappingFailed", zap.String("eventType", event.Name), zap.Error(err))
			return
		}
		eventID := extractEventID(payload)
		if eventID == "" {
			eventID = uuid.NewString()
		}
		outboxEvent := Event{
			EventID: eventID,
			Type:    outboxType,
			Payload: payload,
		}
		if s.outbox == nil {
			s.logger.Warn("NoOutboxConfigured", zap.String("eventType", event.Name))
			return
		}
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := s.outbox.Enqueue(ctx, outboxEvent); err != nil {
			s.logger.Error("EventEnqueueFailed", zap.String("eventType", event.Name), zap.Error(err))
			return
		}
		s.logger.Info("EventEnqueued", zap.String("eventType", event.Name), zap.String("eventID", outboxEvent.EventID))
	}
}

// mapEventPayload converts a domain event payload to a webhook-friendly map.
// Unknown payload shapes fall back to a {"raw": ...} wrapper so no event is
// silently dropped.
func mapEventPayload(event eventbus.Event) (map[string]any, error) {
	switch p := event.Payload.(type) {
	case nil:
		return map[string]any{}, nil
	case map[string]any:
		return p, nil
	default:
		return map[string]any{"raw": p}, nil
	}
}

// extractEventID reads the event_id key from a payload map (string or UUID).
func extractEventID(payload map[string]any) string {
	raw, ok := payload["event_id"]
	if !ok {
		return ""
	}
	switch v := raw.(type) {
	case string:
		return v
	default:
		return ""
	}
}
