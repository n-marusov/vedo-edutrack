package webhook

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// Metrics for the webhook delivery system (F6). Registered lazily via
// promauto — the Prometheus registry (mounted at /metrics) picks them up.
var (
	// WebhookDeliveryTotal counts webhook delivery attempts by status.
	WebhookDeliveryTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "webhook_delivery_total",
			Help: "Total webhook delivery attempts, by status.",
		},
		[]string{"status"},
	)
	// WebhookDeliveryDurationSeconds times webhook deliveries.
	WebhookDeliveryDurationSeconds = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "webhook_delivery_duration_seconds",
			Help:    "Webhook delivery duration in seconds.",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"status"},
	)
)
