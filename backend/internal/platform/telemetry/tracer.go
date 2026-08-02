// Package telemetry provides OpenTelemetry instrumentation for the backend.
//
// It initializes the OTel tracer provider and Prometheus meter provider,
// wiring them into the modular monolith's shared infrastructure layer.
//
// See ADR-IMPL.PROCESS.development-tooling §9 for the observability architecture.
package telemetry

import (
	"context"
	"log"
)

// InitTracer initializes the OTel tracer provider.
//
// TODO: Wire go.opentelemetry.io/otel TracerProvider with OTLP exporter.
// See ADR-IMPL.PROCESS.development-tooling §9: OTel SDK → OTLP → Collector.
func InitTracer(ctx context.Context, serviceName string) error {
	log.Printf("[INFO] [telemetry.InitTracer] initializing tracer for service %q", serviceName)
	_ = ctx
	_ = serviceName
	return nil
}
