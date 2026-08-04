// Package queries implements webhook subscription query use cases (F6):
// listing subscriptions for a tenant and paginated delivery history.
package queries

import (
	"context"

	"go.uber.org/zap"

	domain "vedo-edutrack/backend/internal/modules/integrations/domain"
)

// DeliveryRepository is the port for webhook delivery history reads.
type DeliveryRepository interface {
	// ListBySubscription returns deliveries for a subscription (newest first),
	// paginated (page/limit).
	ListBySubscription(ctx context.Context, subscriptionID domain.SubscriptionID, page, limit int) ([]domain.WebhookDelivery, error)
}

// SubscriptionQueries implements webhook subscription query use cases.
type SubscriptionQueries struct {
	service    *domain.SubscriptionService
	deliveries DeliveryRepository
	logger     *zap.Logger
}

// NewSubscriptionQueries builds the webhook subscription query service.
func NewSubscriptionQueries(service *domain.SubscriptionService, deliveries DeliveryRepository, logger *zap.Logger) *SubscriptionQueries {
	if logger == nil {
		logger = zap.NewNop()
	}
	return &SubscriptionQueries{service: service, deliveries: deliveries, logger: logger.Named("integrations.app")}
}

// ListWebhookSubscriptions returns the tenant's subscriptions, optionally
// filtered by active status.
func (q *SubscriptionQueries) ListWebhookSubscriptions(ctx context.Context, tenantID string, active *bool) ([]domain.WebhookSubscription, error) {
	return q.service.ListSubscriptions(ctx, tenantID, active)
}

// GetWebhookDeliveryHistory returns a subscription's delivery history
// (tenant-scoped), paginated.
func (q *SubscriptionQueries) GetWebhookDeliveryHistory(ctx context.Context, subscriptionID domain.SubscriptionID, tenantID string, page, limit int) ([]domain.WebhookDelivery, error) {
	// Tenant-scope the read: the service verifies ownership first.
	if _, err := q.service.GetSubscription(ctx, subscriptionID, tenantID); err != nil {
		return nil, err
	}
	if q.deliveries == nil {
		return []domain.WebhookDelivery{}, nil
	}
	if page < 0 {
		page = 0
	}
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	return q.deliveries.ListBySubscription(ctx, subscriptionID, page, limit)
}
