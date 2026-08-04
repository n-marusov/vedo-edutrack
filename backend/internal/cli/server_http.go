package cli

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"net"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/spf13/cobra"
	"go.uber.org/zap"

	"vedo-edutrack/backend/internal/api"
	"vedo-edutrack/backend/internal/platform/auth"
	"vedo-edutrack/backend/internal/platform/config"
	platformpostgres "vedo-edutrack/backend/internal/platform/postgres"
	"vedo-edutrack/backend/internal/platform/spa"
)

// startTime anchors the uptime reported by health endpoints.
var startTime = time.Now()

// Health check status values (used by /readyz checks).
const (
	healthStatusOK       = "ok"
	healthStatusDegraded = "degraded"
	checkUp              = "up"
	checkDown            = "down"
)

// dbPool is the shared PostgreSQL pool, initialized when the server starts
// (used by the /readyz database check).
var dbPool *pgxpool.Pool

// platformDBPool returns the shared pool (nil when not connected).
func platformDBPool() *pgxpool.Pool {
	return dbPool
}

// healthResponse is the JSON payload of /healthz and /readyz.
type healthResponse struct {
	Status  string            `json:"status"` // ok | degraded
	Version string            `json:"version"`
	Env     string            `json:"env"`
	Uptime  string            `json:"uptime"` // e.g. "1m30s"
	Checks  map[string]string `json:"checks,omitempty"`
}

// serveHTTP builds and runs the chi HTTP server with graceful shutdown.
// It blocks until SIGINT/SIGTERM or a fatal server error.
func serveHTTP(_ *cobra.Command) error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	// Initialize the local JWT auth (auto-generates a dev key if missing).
	a, err := auth.New("", cfg.JWTIssuer, cfg.JWTAudience, cfg.JWTTokenTTL, zapLogger)
	if err != nil {
		return fmt.Errorf("init auth: %w", err)
	}

	// Connect the PostgreSQL pool for /readyz (non-fatal in dev: the server
	// still starts without a database; readiness reports database: down).
	connectCtx, connectCancel := context.WithTimeout(context.Background(), 5*time.Second)
	if pool, poolErr := platformpostgres.Connect(connectCtx, cfg.DatabaseURL); poolErr != nil {
		connectCancel()
		zapLogger.Warn("database connect failed (readiness will report down)", zap.Error(poolErr))
	} else {
		connectCancel()
		dbPool = pool
		zapLogger.Info("database pool ready")
	}

	// Webhook services (F6): shared in-memory subscription repo + outbox-based
	// delivery worker. The worker delivers outbox events to subscriber URLs
	// with HMAC signing and retry/deactivation rules (M4 Task 8).
	webhooks, stopWorker := startWebhookWorker(zapLogger)
	defer stopWorker()

	r := newRouter(cfg, a, webhooks)

	addr := ":" + strconv.Itoa(cfg.Port)
	srv := &http.Server{
		Addr:              addr,
		Handler:           r,
		ReadTimeout:       30 * time.Second,
		WriteTimeout:      30 * time.Second,
		IdleTimeout:       30 * time.Second,
		ReadHeaderTimeout: 10 * time.Second,
	}

	// Graceful shutdown: trap SIGINT/SIGTERM, drain in-flight requests for
	// up to 30s, then close the DB pool and flush OTel/loggers.
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	errCh := make(chan error, 1)
	go func() {
		zapLogger.Info("server starting", zap.String("addr", addr))
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errCh <- err
		}
	}()

	select {
	case err := <-errCh:
		return fmt.Errorf("server: %w", err)
	case <-ctx.Done():
		zapLogger.Info("server stopping (signal received)")
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		zapLogger.Warn("graceful shutdown failed", zap.Error(err))
	}

	// Close the database pool.
	if dbPool != nil {
		dbPool.Close()
		dbPool = nil
	}

	zapLogger.Info("server stopped cleanly")
	return nil
}

// securityHeaders sets baseline security headers on every response. This is
// especially important for the embedded SPA (on-prem contour) which is served
// without the Traefik edge headers; the edge may add/override more (HSTS, CSP
// tuning). CSP is deliberately permissive for the Vite dev server (inline
// styles/scripts) — tighten for production.
func securityHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h := w.Header()
		h.Set("X-Content-Type-Options", "nosniff")
		h.Set("X-Frame-Options", "DENY")
		h.Set("Referrer-Policy", "strict-origin-when-cross-origin")
		h.Set("Permissions-Policy", "camera=(), microphone=(), geolocation=()")
		h.Set("Content-Security-Policy",
			"default-src 'self'; script-src 'self' 'unsafe-inline'; style-src 'self' 'unsafe-inline'; img-src 'self' data:; font-src 'self'; connect-src 'self' http://localhost:*")
		next.ServeHTTP(w, r)
	})
}

// newRouter builds the chi router with the full M0.3 middleware stack and
// route groups. The auth bundle is optional (nil for tests that don't
// exercise authentication); webhooks is the shared webhook services bundle
// (nil → handler builds default in-memory services).
func newRouter(cfg *config.Config, a *auth.Auth, webhooks *api.WebhookServices) *chi.Mux {
	r := chi.NewRouter()

	// Middleware stack (ADR-DES.API.communication-patterns):
	// request ID, request logger, panic recovery, timeout.
	// RealIP is intentionally NOT used: chi's RealIP trusts X-Forwarded-For
	// blindly (spoofing risk, GHSA-3fxj-6jh8-hvhx); the Traefik edge sets the
	// trusted header and we rely on it only behind the edge.
	r.Use(securityHeaders)
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(30 * time.Second))

	// CORS: allow the Vite dev server and the embedded SPA origin.
	// AllowCredentials requires EXPLICIT origins — browsers reject
	// Access-Control-Allow-Origin: * for credentialed requests (CORS spec),
	// so a wildcard here would silently break dev browser auth from any
	// non-listed port.
	allowedOrigins := []string{
		"http://localhost:5173", // Vite dev server
		"http://127.0.0.1:5173",
		"http://localhost:8080", // embedded SPA (same-origin; listed for clarity)
		"http://127.0.0.1:8080",
	}
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   allowedOrigins,
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-Requested-With"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	// Health endpoints.
	r.Get("/healthz", func(w http.ResponseWriter, _ *http.Request) {
		writeHealth(w, http.StatusOK, healthResponse{
			Status:  healthStatusOK,
			Version: config.Version,
			Env:     cfg.Environment,
			Uptime:  time.Since(startTime).Round(time.Second).String(),
		})
	})
	r.Get("/readyz", readyzHandler(cfg))

	// Prometheus metrics.
	r.Handle("/metrics", promhttp.Handler())

	// Authentication (T5): JWKS endpoint, dev token issuance, and the
	// JWT validation middleware applied to the API group.
	if a != nil {
		r.Get("/.well-known/jwks.json", a.JWKSHandler())
	} else {
		r.Get("/.well-known/jwks.json", func(w http.ResponseWriter, _ *http.Request) {
			http.Error(w, `{"error":"not_implemented","endpoint":"/.well-known/jwks.json"}`, http.StatusNotImplemented)
		})
	}

	// API group — generated server interface (T8). Protected by the auth
	// middleware when a is non-nil; the token endpoint stays public.
	r.Route("/api/v1", func(r chi.Router) {
		if a != nil {
			// Public: dev token issuance.
			r.Route("/auth", func(r chi.Router) {
				r.Post("/token", a.TokenHandler())
			})
			// Protected: everything else under /api/v1.
			r = r.With(a.Middleware())
		}
		h := api.NewStubHandler(cfg, zapLogger, webhooks)
		api.HandlerWithOptions(h, api.ChiServerOptions{
			BaseRouter: r,
		})
	})

	// SPA fallback (T17): serve the embedded frontend for unknown paths;
	// index.html handles client-side routing. Graceful when not embedded.
	r.NotFound(spaFallback())

	return r
}

// readyzHandler reports readiness: database reachability and identity-provider
// reachability. Both must pass for the overall status to be "ok".
func readyzHandler(cfg *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		resp := healthResponse{
			Status:  healthStatusOK,
			Version: config.Version,
			Env:     cfg.Environment,
			Uptime:  time.Since(startTime).Round(time.Second).String(),
			Checks:  map[string]string{},
		}
		status := http.StatusOK

		// Database check: pgx pool ping; falls back to TCP dial when the pool
		// is not initialized (e.g. local run without PostgreSQL).
		if pool := platformDBPool(); pool != nil {
			if err := platformpostgres.Ping(r.Context(), pool); err != nil {
				resp.Checks["database"] = checkDown
				resp.Status = healthStatusDegraded
				status = http.StatusServiceUnavailable
				zapLogger.Warn("readyz database ping failed", zap.Error(err))
			} else {
				resp.Checks["database"] = checkUp
			}
		} else if err := checkDatabaseReachable(cfg.DatabaseURL); err != nil {
			resp.Checks["database"] = checkDown
			resp.Status = healthStatusDegraded
			status = http.StatusServiceUnavailable
			zapLogger.Warn("readyz database check failed", zap.Error(err))
		} else {
			resp.Checks["database"] = checkUp
		}

		if err := checkJWKSReachable(cfg.JWKSURL); err != nil {
			resp.Checks["identity_provider"] = checkDown
			resp.Status = healthStatusDegraded
			status = http.StatusServiceUnavailable
			zapLogger.Warn("readyz identity_provider check failed", zap.Error(err))
		} else {
			resp.Checks["identity_provider"] = checkUp
		}

		writeHealth(w, status, resp)
	}
}

// checkDatabaseReachable dials the DATABASE_URL host (TCP) to verify the
// database is reachable. A real pgx pool ping replaces this in T26.
func checkDatabaseReachable(databaseURL string) error {
	u, err := url.Parse(databaseURL)
	if err != nil {
		return err
	}
	if u.Host == "" {
		return errors.New("DATABASE_URL has no host")
	}
	return dialHost(u.Host)
}

// checkJWKSReachable dials the JWKS URL host to verify the identity provider
// is reachable.
func checkJWKSReachable(jwksURL string) error {
	u, err := url.Parse(jwksURL)
	if err != nil {
		return err
	}
	if u.Host == "" {
		return errors.New("JWKS_URL has no host")
	}
	return dialHost(u.Host)
}

// dialHost opens and closes a TCP connection to host with a 2s timeout.
func dialHost(host string) error {
	conn, err := net.DialTimeout("tcp", host, 2*time.Second)
	if err != nil {
		return err
	}
	_ = conn.Close()
	return nil
}

// writeHealth writes a healthResponse as JSON.
func writeHealth(w http.ResponseWriter, status int, resp healthResponse) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		zapLogger.Warn("health encode", zap.Error(err))
	}
}

// spaFallback serves the embedded SPA; unknown paths fall back to index.html
// (react-router handles client-side routing). When the frontend is not
// embedded (dev without build), it answers a JSON hint.
func spaFallback() http.HandlerFunc {
	fss, ok := spa.EmbeddedFS()
	if !ok {
		return func(w http.ResponseWriter, _ *http.Request) {
			zapLogger.Warn("SPA not embedded (dev mode)")
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"spa":"not embedded (dev mode)"}`))
		}
	}

	fileServer := http.FileServer(http.FS(fss))
	return func(w http.ResponseWriter, r *http.Request) {
		path := strings.TrimPrefix(r.URL.Path, "/")
		// Serve the requested file if it exists; otherwise index.html.
		if path == "" {
			http.ServeFileFS(w, r, fss, "index.html")
			return
		}
		if _, err := fs.Stat(fss, path); err == nil {
			fileServer.ServeHTTP(w, r)
			return
		}
		http.ServeFileFS(w, r, fss, "index.html")
	}
}
