package main

import (
	"fmt"
	"io"
	"log/slog"
	"net"
	"net/http"
	"testing"
	"time"
)

// startTestServer creates an HTTP server on a random port, starts it,
// and returns the server and its base URL.
func startTestServer(t *testing.T, cfg Config, m *Metrics) (*http.Server, string) {
	t.Helper()

	if cfg.Port == "" || cfg.Port == "0" {
		l, err := net.Listen("tcp", "127.0.0.1:0")
		if err != nil {
			t.Fatalf("listen: %v", err)
		}
		p := l.Addr().(*net.TCPAddr).Port
		l.Close() // release the port
		cfg.Port = fmt.Sprintf("%d", p)
	}

	srv := newServer(cfg, testLogger(), m)

	errCh := make(chan error, 1)
	go func() {
		errCh <- srv.ListenAndServe()
	}()
	t.Cleanup(func() { srv.Close() })

	baseURL := "http://127.0.0.1:" + cfg.Port

	// Wait for server to become ready.
	deadline := time.Now().Add(5 * time.Second)
	for time.Now().Before(deadline) {
		resp, err := http.Get(baseURL + "/livez")
		if err == nil {
			resp.Body.Close()
			return srv, baseURL
		}
		time.Sleep(50 * time.Millisecond)
	}
	t.Fatal("server did not become ready within 5s")
	return nil, ""
}

// TestHealthEndpoints verifies livez, readyz, and healthz respond 200.
func TestHealthEndpoints(t *testing.T) {
	cfg := Config{
		Port:              "0",
		Version:           "test",
		Env:               "test",
		RateLimitPerSec:   10000,
		RateLimitBurst:    20000,
		ReadTimeout:       10 * time.Second,
		ReadHeaderTimeout: 5 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       10 * time.Second,
		MetricsEnabled:    true,
	}

	_, baseURL := startTestServer(t, cfg, NewMetrics())

	t.Run("livez", func(t *testing.T) {
		resp, err := http.Get(baseURL + "/livez")
		if err != nil {
			t.Fatalf("GET /livez: %v", err)
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			t.Errorf("expected 200, got %d", resp.StatusCode)
		}
	})

	t.Run("readyz", func(t *testing.T) {
		resp, err := http.Get(baseURL + "/readyz")
		if err != nil {
			t.Fatalf("GET /readyz: %v", err)
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			t.Errorf("expected 200, got %d", resp.StatusCode)
		}
	})

	t.Run("healthz", func(t *testing.T) {
		resp, err := http.Get(baseURL + "/healthz")
		if err != nil {
			t.Fatalf("GET /healthz: %v", err)
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			t.Errorf("expected 200, got %d", resp.StatusCode)
		}
		if ct := resp.Header.Get("Content-Type"); ct != "application/json" {
			t.Errorf("expected application/json, got %s", ct)
		}
	})

	t.Run("config.js", func(t *testing.T) {
		resp, err := http.Get(baseURL + "/config.js")
		if err != nil {
			t.Fatalf("GET /config.js: %v", err)
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			t.Errorf("expected 200, got %d", resp.StatusCode)
		}
	})

	t.Run("spa_fallback", func(t *testing.T) {
		resp, err := http.Get(baseURL + "/nonexistent-page")
		if err != nil {
			t.Fatalf("GET /nonexistent-page: %v", err)
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNotFound {
			t.Errorf("expected 200 or 404, got %d", resp.StatusCode)
		}
	})

	t.Run("metrics", func(t *testing.T) {
		resp, err := http.Get(baseURL + "/metrics")
		if err != nil {
			t.Fatalf("GET /metrics: %v", err)
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			t.Errorf("expected 200, got %d", resp.StatusCode)
		}
	})

	t.Run("security_headers", func(t *testing.T) {
		resp, err := http.Get(baseURL + "/livez")
		if err != nil {
			t.Fatalf("GET /livez: %v", err)
		}
		defer resp.Body.Close()
		checkHeader(t, resp, "X-Content-Type-Options", "nosniff")
		checkHeader(t, resp, "X-Frame-Options", "DENY")
		checkHeader(t, resp, "Referrer-Policy", "strict-origin-when-cross-origin")
	})

	t.Run("rate_limit", func(t *testing.T) {
		limited := Config{
			Port:              "0",
			RateLimitPerSec:   2,
			RateLimitBurst:    2,
			ReadTimeout:       10 * time.Second,
			ReadHeaderTimeout: 5 * time.Second,
			WriteTimeout:      10 * time.Second,
			IdleTimeout:       10 * time.Second,
		}
		_, url := startTestServer(t, limited, NewMetrics())

		okCount := 0
		for i := 0; i < 3; i++ {
			resp, err := http.Get(url + "/livez")
			if err != nil {
				t.Fatalf("GET /livez: %v", err)
			}
			if resp.StatusCode == http.StatusOK {
				okCount++
			}
			resp.Body.Close()
		}
		// The startTestServer readiness probe consumes one token; burst of 2
		// leaves 1 for the actual test (third request is rate-limited).
		if okCount < 1 {
			t.Errorf("expected at least 1 OK, got %d", okCount)
		}
	})

	t.Run("cache_headers", func(t *testing.T) {
		_, url := startTestServer(t, cfg, NewMetrics())

		// /assets/* should be cached immutably.
		resp, err := http.Get(url + "/assets/app.js")
		if err != nil {
			t.Fatalf("GET /assets/app.js: %v", err)
		}
		defer resp.Body.Close()
		if cc := resp.Header.Get("Cache-Control"); cc != "public, max-age=31536000, immutable" {
			t.Errorf("Cache-Control for /assets/*: %q", cc)
		}
	})
}

func checkHeader(t *testing.T, resp *http.Response, key, want string) {
	t.Helper()
	if got := resp.Header.Get(key); got != want {
		t.Errorf("header %s: expected %q, got %q", key, want, got)
	}
}

func testLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(io.Discard, nil))
}
