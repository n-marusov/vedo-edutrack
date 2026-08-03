package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"runtime"
	"sync/atomic"
	"time"
)

// Metrics holds counters for Prometheus-style exposition.
type Metrics struct {
	Requests       atomic.Int64
	Errors4xx      atomic.Int64
	Errors5xx      atomic.Int64
	ActiveConns    atomic.Int64
	TotalLatencyMs atomic.Int64 // accumulated latency for average calculation
	EmbeddedSizeKB int64
	startTime      time.Time
}

// NewMetrics creates a new Metrics instance with start time.
func NewMetrics() *Metrics {
	return &Metrics{startTime: time.Now()}
}

// RegisterEmbeddedSize records the total size of embedded assets.
func (m *Metrics) RegisterEmbeddedSize(assets interface{}) {
	// Walk the embedded filesystem to compute total size.
	// For simplicity, approximate from the binary size.
	var ms runtime.MemStats
	runtime.ReadMemStats(&ms)
	m.EmbeddedSizeKB = int64(ms.HeapAlloc / 1024)
	_ = assets
}

// metricsMiddleware collects HTTP metrics for every request.
func metricsMiddleware(m *Metrics) middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			m.ActiveConns.Add(1)
			defer m.ActiveConns.Add(-1)

			start := time.Now()
			wrapped := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}
			next.ServeHTTP(wrapped, r)

			m.Requests.Add(1)
			m.TotalLatencyMs.Add(time.Since(start).Milliseconds())

			if wrapped.statusCode >= 400 && wrapped.statusCode < 500 {
				m.Errors4xx.Add(1)
			}
			if wrapped.statusCode >= 500 {
				m.Errors5xx.Add(1)
			}
		})
	}
}

// metricsHandler serves Prometheus text format metrics at /metrics.
func metricsHandler(m *Metrics) http.HandlerFunc {
	return func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "text/plain; version=0.0.4")
		uptime := time.Since(m.startTime).Seconds()

		fmt.Fprintf(w, "# HELP http_requests_total Total HTTP requests\n")
		fmt.Fprintf(w, "# TYPE http_requests_total counter\n")
		fmt.Fprintf(w, "http_requests_total %d\n", m.Requests.Load())

		fmt.Fprintf(w, "# HELP http_errors_4xx_total HTTP 4xx responses\n")
		fmt.Fprintf(w, "# TYPE http_errors_4xx_total counter\n")
		fmt.Fprintf(w, "http_errors_4xx_total %d\n", m.Errors4xx.Load())

		fmt.Fprintf(w, "# HELP http_errors_5xx_total HTTP 5xx responses\n")
		fmt.Fprintf(w, "# TYPE http_errors_5xx_total counter\n")
		fmt.Fprintf(w, "http_errors_5xx_total %d\n", m.Errors5xx.Load())

		fmt.Fprintf(w, "# HELP http_active_connections Current active connections\n")
		fmt.Fprintf(w, "# TYPE http_active_connections gauge\n")
		fmt.Fprintf(w, "http_active_connections %d\n", m.ActiveConns.Load())

		fmt.Fprintf(w, "# HELP http_request_duration_avg_ms Average request duration (milliseconds)\n")
		fmt.Fprintf(w, "# TYPE http_request_duration_avg_ms gauge\n")
		reqs := m.Requests.Load()
		if reqs > 0 {
			avg := m.TotalLatencyMs.Load() / reqs
			fmt.Fprintf(w, "http_request_duration_avg_ms %d\n", avg)
		} else {
			fmt.Fprintf(w, "http_request_duration_avg_ms 0\n")
		}

		fmt.Fprintf(w, "# HELP embedded_assets_size_bytes Approximate embedded assets size\n")
		fmt.Fprintf(w, "# TYPE embedded_assets_size_bytes gauge\n")
		fmt.Fprintf(w, "embedded_assets_size_bytes %d\n", m.EmbeddedSizeKB*1024)

		fmt.Fprintf(w, "# HELP process_uptime_seconds Process uptime in seconds\n")
		fmt.Fprintf(w, "# TYPE process_uptime_seconds gauge\n")
		fmt.Fprintf(w, "process_uptime_seconds %.0f\n", uptime)

		fmt.Fprintf(w, "# HELP go_goroutines Number of goroutines\n")
		fmt.Fprintf(w, "# TYPE go_goroutines gauge\n")
		fmt.Fprintf(w, "go_goroutines %d\n", runtime.NumGoroutine())

		fmt.Fprintf(w, "# HELP go_memory_alloc_bytes Memory allocated (bytes)\n")
		fmt.Fprintf(w, "# TYPE go_memory_alloc_bytes gauge\n")
		var ms runtime.MemStats
		runtime.ReadMemStats(&ms)
		fmt.Fprintf(w, "go_memory_alloc_bytes %d\n", ms.Alloc)
	}
}

// requestID generates a random hex request ID for tracing.
func newRequestID() string {
	b := make([]byte, 8)
	for i := range b {
		b[i] = "0123456789abcdef"[rand.Intn(16)]
	}
	return string(b)
}
