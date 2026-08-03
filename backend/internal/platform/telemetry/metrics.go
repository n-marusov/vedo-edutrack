// Package telemetry — Prometheus metrics initialization.
//
// See REQ-NFR-ops.observability.golden-signals-dashboards for the
// golden-signals dashboard requirements (latency, traffic, errors, saturation).
package telemetry

import (
	"context"
	"fmt"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
)

// InitMeter initializes the OTel meter provider with an OTLP gRPC exporter.
//
// Exported metrics feed the Grafana golden-signals dashboard
// (REQ-NFR-ops.observability.golden-signals-dashboards). Export failures are
// non-fatal in dev: the provider is still created so the app can start
// without a collector.
func InitMeter(ctx context.Context, serviceName string) (*sdkmetric.MeterProvider, error) {
	endpoint := envOrDefault("OTEL_EXPORTER_OTLP_ENDPOINT", "http://otel-collector:4317")

	exporter, err := otlpmetricgrpc.New(ctx,
		otlpmetricgrpc.WithEndpointURL(endpoint),
		otlpmetricgrpc.WithTimeout(5*time.Second),
	)
	if err != nil {
		return nil, fmt.Errorf("create OTLP metric exporter: %w", err)
	}

	res, err := resourceFromEnv(serviceName)
	if err != nil {
		return nil, fmt.Errorf("build metric resource: %w", err)
	}

	mp := sdkmetric.NewMeterProvider(
		sdkmetric.WithReader(sdkmetric.NewPeriodicReader(exporter,
			sdkmetric.WithInterval(30*time.Second))),
		sdkmetric.WithResource(res),
	)
	otel.SetMeterProvider(mp)
	return mp, nil
}

// ShutdownMeter gracefully flushes and shuts down the meter provider.
func ShutdownMeter(ctx context.Context, mp *sdkmetric.MeterProvider) {
	if mp == nil {
		return
	}
	_ = mp.Shutdown(ctx)
}
