package sparql

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// Metrics for the read-only SPARQL endpoint (F6). Registered lazily via
// promauto — the Prometheus registry (mounted at /metrics) picks them up.
var (
	// SparqlQueryDurationSeconds times SPARQL query execution by outcome.
	SparqlQueryDurationSeconds = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "sparql_query_duration_seconds",
			Help:    "SPARQL query execution duration in seconds.",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"outcome"},
	)
)
