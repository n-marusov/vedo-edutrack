package api

import (
	"context"

	"vedo-edutrack/backend/internal/modules/integrations/adapters/webhook"
	integapp "vedo-edutrack/backend/internal/modules/integrations/application/commands"
)

// outboxBridge adapts the application-layer OutboxEvent to the webhook
// adapter's Event so command use cases can enqueue into the real outbox
// without importing the adapter directly.
type outboxBridge struct {
	outbox *webhook.PostgresOutbox
}

// NewOutboxBridge builds the bridge over an outbox (nil = no-op enqueue).
func NewOutboxBridge(outbox *webhook.PostgresOutbox) outboxBridge {
	return outboxBridge{outbox: outbox}
}

// Enqueue implements integapp.OutboxEnqueuer.
func (b outboxBridge) Enqueue(ctx context.Context, event integapp.OutboxEvent) error {
	if b.outbox == nil {
		// No outbox configured (in-memory wiring): acknowledge silently — the
		// ping event is best-effort at this layer.
		return nil
	}
	return b.outbox.Enqueue(ctx, webhook.Event{
		EventID: event.EventID,
		Type:    webhook.EventType(event.Type),
		Payload: event.Payload,
	})
}
