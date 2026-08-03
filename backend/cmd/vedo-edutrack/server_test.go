package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"vedo-edutrack/backend/internal/platform/config"
)

func TestHealthzLiveness(t *testing.T) {
	oldVersion := config.Version
	config.Version = "9.9.9-test"
	t.Cleanup(func() { config.Version = oldVersion })

	cfg := &config.Config{Port: 8080, Environment: "test"}
	mux := newHealthMux(cfg)

	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("healthz status = %d, want 200", rec.Code)
	}
	var body healthResponse
	if err := json.NewDecoder(rec.Body).Decode(&body); err != nil {
		t.Fatalf("healthz decode: %v", err)
	}
	if body.Status != "ok" {
		t.Errorf("healthz status field = %q, want ok", body.Status)
	}
	if body.Version != "9.9.9-test" {
		t.Errorf("healthz version = %q, want 9.9.9-test", body.Version)
	}
	if body.Env != "test" {
		t.Errorf("healthz env = %q, want test", body.Env)
	}
	if body.Uptime == "" {
		t.Error("healthz uptime is empty")
	}
}

func TestReadyzDatabaseDown(t *testing.T) {
	// Port 1: dial fails fast → readiness must report degraded/503.
	cfg := &config.Config{DatabaseURL: "postgres://u:p@127.0.0.1:1/db?sslmode=disable"}
	mux := newHealthMux(cfg)

	req := httptest.NewRequest(http.MethodGet, "/readyz", nil)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusServiceUnavailable {
		t.Fatalf("readyz status = %d, want 503", rec.Code)
	}
	var body healthResponse
	_ = json.NewDecoder(rec.Body).Decode(&body)
	if body.Status != "degraded" {
		t.Errorf("readyz status field = %q, want degraded", body.Status)
	}
	if body.Checks["database"] != "down" {
		t.Errorf("readyz database check = %q, want down", body.Checks["database"])
	}
}

func TestReadyzDatabaseUp(t *testing.T) {
	// A reachable TCP listener stands in for PostgreSQL (dial-level check).
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {}))
	defer ts.Close()

	host := strings.TrimPrefix(ts.URL, "http://")
	u := url.URL{Scheme: "postgres", User: url.UserPassword("u", "p"), Host: host, Path: "/db"}
	cfg := &config.Config{DatabaseURL: u.String()}
	mux := newHealthMux(cfg)

	req := httptest.NewRequest(http.MethodGet, "/readyz", nil)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("readyz status = %d, want 200", rec.Code)
	}
	var body healthResponse
	_ = json.NewDecoder(rec.Body).Decode(&body)
	if body.Status != "ok" {
		t.Errorf("readyz status field = %q, want ok", body.Status)
	}
	if body.Checks["database"] != "up" {
		t.Errorf("readyz database check = %q, want up", body.Checks["database"])
	}
}
