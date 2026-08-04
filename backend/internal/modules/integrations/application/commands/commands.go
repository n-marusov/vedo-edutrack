// Package commands implements webhook subscription command use cases (F6):
// create / update / delete subscriptions through the domain service, plus
// the outbox verification ping enqueued on creation.
package commands

import (
	"context"
	"errors"

	"go.uber.org/zap"

	domain "vedo-edutrack/backend/internal/modules/integrations/domain"
)

// ErrSubscriptionNotFound is returned when the subscription does not exist.
var ErrSubscriptionNotFound = errors.New("webhook subscription not found")

// CommandResult is the outcome of a subscription command.
type CommandResult struct {
	Subscription *domain.WebhookSubscription
}

// SubscriptionCommands implements webhook subscription command use cases.
type SubscriptionCommands struct {
	service *domain.SubscriptionService
	outbox  OutboxEnqueuer
	logger  *zap.Logger
}

// OutboxEnqueuer is the port for enqueueing outbox events (implemented by the
// webhook outbox adapter). Kept minimal so commands don't depend on the
// adapter package.
type OutboxEnqueuer interface {
	Enqueue(ctx context.Context, event OutboxEvent) error
}

// OutboxEvent is the outbox event wire shape used by the command layer.
type OutboxEvent struct {
	EventID string
	Type    string
	Payload map[string]any
}

// NewSubscriptionCommands builds the webhook subscription command service.
func NewSubscriptionCommands(service *domain.SubscriptionService, outbox OutboxEnqueuer, logger *zap.Logger) *SubscriptionCommands {
	if logger == nil {
		logger = zap.NewNop()
	}
	return &SubscriptionCommands{service: service, outbox: outbox, logger: logger.Named("integrations.app")}
}

// CreateWebhookSubscriptionInput is the create-subscription payload.
type CreateWebhookSubscriptionInput struct {
	TenantID      string
	URL           string
	EventTypes    []domain.EventType
	SigningSecret string
}

// CreateWebhookSubscription validates, creates the subscription via the
// domain service and enqueues a verification ping event.
func (c *SubscriptionCommands) CreateWebhookSubscription(ctx context.Context, in CreateWebhookSubscriptionInput) (*domain.WebhookSubscription, error) {
	sub, err := c.service.CreateSubscription(ctx, in.TenantID, in.URL, in.EventTypes, in.SigningSecret)
	if err != nil {
		return nil, err
	}

	c.logger.Info("SubscriptionCreated", zap.String("subscriptionID", sub.ID.String()), zap.String("url", sub.URL), zap.Strings("events", eventTypeStrings(sub.EventTypes)))

	// Enqueue a verification ping so the worker confirms the endpoint.
	if c.outbox != nil {
		if err := c.outbox.Enqueue(ctx, OutboxEvent{
			EventID: sub.ID.String() + "-ping",
			Type:    "webhook.ping",
			Payload: map[string]any{"subscription_id": sub.ID.String()},
		}); err != nil {
			c.logger.Warn("verification ping enqueue failed", zap.Error(err))
		}
	}
	return sub, nil
}

// UpdateWebhookSubscriptionInput is the update-subscription payload.
type UpdateWebhookSubscriptionInput struct {
	SubscriptionID domain.SubscriptionID
	TenantID       string
	URL            string
	EventTypes     []domain.EventType
	SigningSecret  string
}

// UpdateWebhookSubscription applies URL/events/secret updates and resets the
// failure counter.
func (c *SubscriptionCommands) UpdateWebhookSubscription(ctx context.Context, in UpdateWebhookSubscriptionInput) (*domain.WebhookSubscription, error) {
	sub, err := c.service.UpdateSubscription(ctx, in.SubscriptionID, in.TenantID, in.URL, in.EventTypes, in.SigningSecret)
	if err != nil {
		return nil, err
	}
	c.logger.Info("SubscriptionUpdated", zap.String("subscriptionID", sub.ID.String()), zap.String("url", sub.URL))
	return sub, nil
}

// DeleteWebhookSubscriptionInput is the delete-subscription payload.
type DeleteWebhookSubscriptionInput struct {
	SubscriptionID domain.SubscriptionID
	TenantID       string
}

// DeleteWebhookSubscription soft-deletes a subscription.
func (c *SubscriptionCommands) DeleteWebhookSubscription(ctx context.Context, in DeleteWebhookSubscriptionInput) error {
	if err := c.service.DeleteSubscription(ctx, in.SubscriptionID, in.TenantID); err != nil {
		return err
	}
	c.logger.Info("SubscriptionDeleted", zap.String("subscriptionID", in.SubscriptionID.String()))
	return nil
}

// eventTypeStrings converts domain event types to their wire strings.
func eventTypeStrings(types []domain.EventType) []string {
	out := make([]string, 0, len(types))
	for _, t := range types {
		out = append(out, t.String())
	}
	return out
}
