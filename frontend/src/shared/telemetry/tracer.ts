import { OTLPTraceExporter } from '@opentelemetry/exporter-trace-otlp-http';
import { resourceFromAttributes } from '@opentelemetry/resources';
import { SimpleSpanProcessor, WebTracerProvider } from '@opentelemetry/sdk-trace-web';
import { SEMRESATTRS_SERVICE_NAME } from '@opentelemetry/semantic-conventions';

let initialized = false;

/**
 * Initializes the browser OTel tracer (web → collector via OTLP/HTTP).
 * Idempotent: subsequent calls are no-ops.
 *
 * Configuration: VITE_OTEL_EXPORTER_OTLP_ENDPOINT (collector OTLP HTTP
 * endpoint; default matches the dev docker-compose collector).
 */
export function initTelemetry(): void {
  if (initialized) {
    return;
  }
  initialized = true;

  const endpoint =
    import.meta.env.VITE_OTEL_EXPORTER_OTLP_ENDPOINT ?? 'http://localhost:4318/v1/traces';

  const exporter = new OTLPTraceExporter({ url: endpoint });
  const provider = new WebTracerProvider({
    resource: resourceFromAttributes({
      [SEMRESATTRS_SERVICE_NAME]: 'vedo-edutrack-web',
    }),
    spanProcessors: [new SimpleSpanProcessor(exporter)],
  });
  provider.register();

  console.info('[telemetry] web tracer initialized', { endpoint });
}
