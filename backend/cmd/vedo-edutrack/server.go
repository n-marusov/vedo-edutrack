// Package main — HTTP server with health endpoints.
//
// M0.2 minimal server: exposes liveness (/healthz) and readiness (/readyz)
// endpoints only. The full API surface (REST/SPARQL/webhooks/MCP) lands in M0.3
// (ADR-DES.API.cli-interface, `server` subcommand as the long-running process).
//
// Health contract (ADR-IMPL.PROCESS.development-tooling §8):
//
//	/healthz  — liveness: process is up; always 200 once the server runs.
//	/readyz   — readiness: dependencies reachable (PostgreSQL TCP dial).
//	            JWKS/Keycloak check joins with the identity module (M0.3).
package main

import (
	"encoding/json"
	"errors"
	"log"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"vedo-edutrack/backend/internal/platform/config"
)

var startTime = time.Now()

// healthResponse is the JSON payload of /healthz and /readyz.
type healthResponse struct {
	Status  string            `json:"status"` // ok | degraded
	Version string            `json:"version"`
	Env     string            `json:"env"`
	Uptime  string            `json:"uptime"` // e.g. "1m30s"
	Checks  map[string]string `json:"checks,omitempty"`
}

// newHealthMux builds the health endpoint router.
func newHealthMux(cfg *config.Config) *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/healthz", func(w http.ResponseWriter, _ *http.Request) {
		writeHealth(w, http.StatusOK, healthResponse{
			Status:  "ok",
			Version: config.Version,
			Env:     cfg.Environment,
			Uptime:  time.Since(startTime).Round(time.Second).String(),
		})
	})

	mux.HandleFunc("/readyz", func(w http.ResponseWriter, _ *http.Request) {
		resp := healthResponse{
			Status:  "ok",
			Version: config.Version,
			Env:     cfg.Environment,
			Uptime:  time.Since(startTime).Round(time.Second).String(),
			Checks:  map[string]string{},
		}
		status := http.StatusOK

		if err := checkDatabaseReachable(cfg.DatabaseURL); err != nil {
			resp.Status = "degraded"
			resp.Checks["database"] = "down"
			status = http.StatusServiceUnavailable
		} else {
			resp.Checks["database"] = "up"
		}

		writeHealth(w, status, resp)
	})

	return mux
}

// serveHTTP runs the health HTTP server on cfg.Port (blocking).
func serveHTTP(cfg *config.Config) error {
	addr := ":" + strconv.Itoa(cfg.Port)
	log.Printf("[INFO] [server] listening on %s (health: /healthz /readyz)", addr)
	return http.ListenAndServe(addr, newHealthMux(cfg))
}

func writeHealth(w http.ResponseWriter, status int, resp healthResponse) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Printf("[WARN] [server] health encode: %v", err)
	}
}

// checkDatabaseReachable dials the DATABASE_URL host (TCP) to verify the
// database is reachable without a driver dependency (pgx lands with sqlc/M0.3).
func checkDatabaseReachable(databaseURL string) error {
	u, err := url.Parse(databaseURL)
	if err != nil {
		return err
	}
	if u.Host == "" {
		return errors.New("DATABASE_URL has no host")
	}
	conn, err := net.DialTimeout("tcp", u.Host, 2*time.Second)
	if err != nil {
		return err
	}
	_ = conn.Close()
	return nil
}
