package cli

import (
	"context"

	"go.uber.org/zap"

	"vedo-edutrack/backend/internal/api"
	"vedo-edutrack/backend/internal/modules/integrations/adapters/webhook"
	integapp "vedo-edutrack/backend/internal/modules/integrations/application/commands"
	integquery "vedo-edutrack/backend/internal/modules/integrations/application/queries"
	integdomain "vedo-edutrack/backend/internal/modules/integrations/domain"
	"vedo-edutrack/backend/internal/platform/eventbus"
)

// newWebhookServices builds the shared webhook services: in-memory subscription
// repository + delivery recorder, application commands/queries, and the
// delivery worker that polls the Postgres outbox and delivers with HMAC
// signing. The worker is nil when the database is unavailable at startup
// (non-fatal in dev — readiness reports database: down).
func newWebhookServices(logger *zap.Logger) *api.WebhookServices {
	if logger == nil {
		logger = zap.NewNop()
	}

	subRepo := webhook.NewInMemorySubscriptionRepository()
	recorder := webhook.NewInMemoryDeliveryRecorder()
	subService := integdomain.NewSubscriptionService(subRepo)

	// Outbox: PostgreSQL-backed (integrations.outbox_events). The worker polls
	// it; the bridge lets the application commands enqueue ping events.
	var outbox *webhook.PostgresOutbox
	if dbPool != nil {
		outbox = webhook.NewPostgresOutbox(dbPool, logger)
	}

	cmds := integapp.NewSubscriptionCommands(subService, api.NewOutboxBridge(outbox), logger)
	queries := integquery.NewSubscriptionQueries(subService, recorder, logger)

	var worker *webhook.DeliveryWorker
	var bus *eventbus.Bus
	if dbPool != nil {
		worker = webhook.NewDeliveryWorker(webhook.DeliveryWorkerConfig{
			Outbox:        outbox,
			Subscriptions: subRepo,
			Deliveries:    recorder,
			Deactivate:    subService.RecordDeliveryFailure,
		}, logger)
	}

	// In-process event bus (F6.4): domain events from other bounded contexts
	// are mapped to outbox webhook events by the subscriber.
	bus = eventbus.New(logger)
	subscriber := webhook.NewEventSubscriber(bus, outbox, logger)
	subscriber.Register()

	return &api.WebhookServices{
		Cmds:     cmds,
		Queries:  queries,
		Repo:     subRepo,
		Recorder: recorder,
		Worker:   worker,
		Bus:      bus,
	}
}

// startWebhookWorker builds the webhook services and starts the delivery
// worker (when a database pool is available). The returned stop func cancels
// the worker's context; call it on shutdown.
func startWebhookWorker(logger *zap.Logger) (*api.WebhookServices, func()) {
	webhooks := newWebhookServices(logger)
	workerCtx, stopWorker := context.WithCancel(context.Background())
	if webhooks.Worker != nil {
		go webhooks.Worker.Run(workerCtx)
		logger.Info("webhook delivery worker started")
	}
	return webhooks, stopWorker
}
