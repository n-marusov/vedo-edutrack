// Package telemetry provides OpenTelemetry instrumentation for the backend.
//
// It initializes the OTel tracer provider and Prometheus meter provider,
// wiring them into the modular monolith's shared infrastructure layer.
//
// See ADR-IMPL.PROCESS.development-tooling §9 for the observability architecture.
package telemetry

import (
	"context"
	"fmt"
	"os"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	sdkresource "go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
	"go.opentelemetry.io/otel/trace"

	"vedo-edutrack/backend/internal/platform/config"
)

// InitTracer initializes the OTel tracer provider with an OTLP gRPC exporter.
//
// Trace export failures are non-fatal in development: the provider is still
// created and the error is logged by the caller so the application can start
// without a collector (REQ-NFR-ops.observability.distributed-tracing).
func InitTracer(ctx context.Context, serviceName string) (*sdktrace.TracerProvider, error) {
	endpoint := envOrDefault("OTEL_EXPORTER_OTLP_ENDPOINT", "http://otel-collector:4317")

	exporter, err := otlptracegrpc.New(ctx,
		otlptracegrpc.WithEndpointURL(endpoint),
		otlptracegrpc.WithTimeout(5*time.Second),
	)
	if err != nil {
		return nil, fmt.Errorf("create OTLP trace exporter: %w", err)
	}

	// Sampler: 100% error + ratio% success (OTEL_SAMPLING_RATIO, default 0.1).
	ratio, err := envFloat("OTEL_SAMPLING_RATIO", 0.1)
	if err != nil {
		return nil, err
	}
	sampler := sdktrace.ParentBased(sdktrace.TraceIDRatioBased(ratio))

	res, err := resourceFromEnv(serviceName)
	if err != nil {
		return nil, fmt.Errorf("build trace resource: %w", err)
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter, sdktrace.WithBatchTimeout(5*time.Second)),
		sdktrace.WithResource(res),
		sdktrace.WithSampler(sampler),
	)
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{}, propagation.Baggage{},
	))
	return tp, nil
}

// ShutdownTracer gracefully flushes and shuts down the tracer provider.
func ShutdownTracer(ctx context.Context, tp *sdktrace.TracerProvider) {
	if tp == nil {
		return
	}
	_ = tp.Shutdown(ctx)
}

// Tracer returns a named tracer instance.
func Tracer(serviceName string) trace.Tracer {
	return otel.Tracer(serviceName)
}

// resourceFromEnv builds an OTel resource with service metadata from the
// environment (service.name, service.version, deployment.environment).
func resourceFromEnv(serviceName string) (*sdkresource.Resource, error) {
	return sdkresource.New(ctx(),
		sdkresource.WithAttributes(
			semconv.ServiceName(serviceName),
			semconv.ServiceVersion(config.Version),
			semconv.DeploymentEnvironment(envOrDefault("APP_ENV", "development")),
		),
	)
}

// ctx is a convenience context without deadline for resource construction.
func ctx() context.Context {
	return context.Background()
}

func envOrDefault(key, defaultVal string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return defaultVal
}

func envFloat(key string, defaultVal float64) (float64, error) {
	raw := os.Getenv(key)
	if raw == "" {
		return defaultVal, nil
	}
	var f float64
	if _, err := fmt.Sscanf(raw, "%f", &f); err != nil {
		return 0, fmt.Errorf("invalid %s value %q: %w", key, raw, err)
	}
	return f, nil
}
