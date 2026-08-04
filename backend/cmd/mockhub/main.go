// Command mockhub runs the mock VEDO Hub GraphQL server for dev/test/CI.
//
// Documented exception to the single-binary ADR (ADR-DES.API.cli-interface):
// test-only tooling, not shipped, not in SBOM
// (ADR-DES.INFRA.mock-hub-strategy, M0.3 T23).
package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"go.uber.org/zap"

	"vedo-edutrack/backend/internal/testing/mockhub"
)

func main() {
	if err := run(); err != nil {
		log.Fatalf("[FATAL] [mockhub] %v", err)
	}
}

func run() error {
	port := envInt("PORT", 8081)
	ontologyFile := env("ONTOLOGY_FILE", "../traceability.ttl")

	logger, err := zap.NewProduction()
	if err != nil {
		return fmt.Errorf("init logger: %w", err)
	}
	defer func() { _ = logger.Sync() }()

	// Load the ontology at startup.
	f, err := os.Open(ontologyFile) // #nosec G304
	if err != nil {
		return fmt.Errorf("open ontology %s: %w", ontologyFile, err)
	}
	ont, err := mockhub.Parse(f)
	_ = f.Close()
	if err != nil {
		return fmt.Errorf("parse ontology %s: %w", ontologyFile, err)
	}
	classes, props := ont.Counts()

	addr := ":" + strconv.Itoa(port)
	srv := &http.Server{
		Addr:              addr,
		Handler:           mockhub.NewMux(ont, logger),
		ReadTimeout:       30 * time.Second,
		WriteTimeout:      30 * time.Second,
		IdleTimeout:       30 * time.Second,
		ReadHeaderTimeout: 10 * time.Second,
	}

	logger.Info("mockhub starting",
		zap.String("addr", addr),
		zap.String("ontology_file", ontologyFile),
		zap.Int("classes", classes),
		zap.Int("properties", props),
	)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	errCh := make(chan error, 1)
	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errCh <- err
		}
	}()

	select {
	case err := <-errCh:
		return fmt.Errorf("mockhub: %w", err)
	case <-ctx.Done():
		logger.Info("mockhub stopping (signal received)")
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		logger.Warn("graceful shutdown failed", zap.Error(err))
	}
	logger.Info("mockhub stopped cleanly")
	return nil
}

func env(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

func envInt(key string, def int) int {
	raw := os.Getenv(key)
	if raw == "" {
		return def
	}
	n, err := strconv.Atoi(raw)
	if err != nil {
		return def
	}
	return n
}
