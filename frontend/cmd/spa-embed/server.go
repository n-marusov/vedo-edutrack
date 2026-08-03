package main

import (
	"context"
	"io/fs"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"embed"
)

//go:embed all:dist
var assets embed.FS

// newServer creates and configures the HTTP server with all middleware.
func newServer(cfg Config, logger *slog.Logger, metrics *Metrics) *http.Server {
	sub, err := fs.Sub(assets, "dist")
	if err != nil {
		logger.Error("fs.Sub failed", "error", err)
		os.Exit(1)
	}
	fileServer := http.FileServer(http.FS(sub))

	mux := http.NewServeMux()

	// Health endpoints — liveness/readiness/detail.
	mux.HandleFunc("GET /livez", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	mux.HandleFunc("GET /readyz", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	mux.HandleFunc("GET /healthz", healthHandler(cfg))

	// Runtime config — generated from env vars at startup.
	mux.HandleFunc("GET /config.js", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/javascript")
		w.Header().Set("Cache-Control", "no-store")
		w.Write([]byte(cfg.ConfigJS()))
	})

	// Prometheus metrics.
	if cfg.MetricsEnabled {
		mux.HandleFunc("GET /metrics", metricsHandler(metrics))
	}

	// SPA fallback: serve index.html for unknown paths (client-side routing).
	spaHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if _, err := fs.Stat(sub, r.URL.Path[1:]); err != nil {
			r.URL.Path = "/"
		}
		fileServer.ServeHTTP(w, r)
	})
	mux.Handle("/", spaHandler)

	// Middleware chain (outermost first).
	handler := chain(mux,
		recoveryMiddleware(logger),                                   // 1. catch panics
		requestIDMiddleware(logger),                                  // 2. X-Request-ID
		loggingMiddleware(logger),                                    // 3. structured request log
		metricsMiddleware(metrics),                                   // 4. collect metrics
		securityHeadersMiddleware(),                                  // 5. OWASP security headers
		rateLimitMiddleware(cfg.RateLimitPerSec, cfg.RateLimitBurst), // 6. rate limit
		cacheMiddleware(),                                            // 7. Cache-Control
	)

	return &http.Server{
		Addr:              ":" + cfg.Port,
		Handler:           handler,
		ReadTimeout:       cfg.ReadTimeout,
		ReadHeaderTimeout: cfg.ReadHeaderTimeout,
		WriteTimeout:      cfg.WriteTimeout,
		IdleTimeout:       cfg.IdleTimeout,
		MaxHeaderBytes:    1 << 20, // 1 MB
	}
}

// Run starts the server with graceful shutdown.
func Run(cfg Config) {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level:     slog.LevelInfo,
		AddSource: true,
	}))

	metrics := NewMetrics()
	metrics.RegisterEmbeddedSize(assets)

	srv := newServer(cfg, logger, metrics)

	go func() {
		logger.Info("server started",
			"port", cfg.Port,
			"version", cfg.Version,
			"env", cfg.Env,
			"pid", os.Getpid(),
		)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("server listen failed", "error", err)
			os.Exit(1)
		}
	}()

	// Graceful shutdown on SIGTERM/SIGINT.
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	logger.Info("shutting down gracefully...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Error("shutdown error", "error", err)
		os.Exit(1)
	}
	logger.Info("server stopped")
}

// healthHandler returns a detailed health response.
func healthHandler(cfg Config) http.HandlerFunc {
	return func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status":"ok","version":"` + cfg.Version + `","env":"` + cfg.Env + `"}` + "\n"))
	}
}

// HealthCheck performs a self-health HTTP probe (used for container HEALTHCHECK).
func HealthCheck(port string) {
	client := &http.Client{Timeout: 2 * time.Second}
	resp, err := client.Get("http://127.0.0.1:" + port + "/livez")
	if err != nil {
		log.Fatalf("[FATAL] health probe failed: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		log.Fatalf("[FATAL] health probe: status %d", resp.StatusCode)
	}
}
