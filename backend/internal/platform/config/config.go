// Package config provides environment-based configuration loading.
//
// Configuration is read from environment variables (12-factor app style),
// with sensible defaults for local development. See .env.example for
// documented variables.
//
// Dynamic runtime injection (ADR-DES.INFRA.dynamic-config-injection):
//   - build version is injected at build time via
//     -ldflags "-X vedo-edutrack/backend/internal/platform/config.Version=..."
//     (the binary reports the version it was built from);
//   - environment/URLs (APP_ENV, PUBLIC_BASE_URL, VEDO_HUB_API_URL) are read
//     from the environment at runtime — one build, many environments.
//
// See REQ-NFR-ops.observability.log-level-config: LOG_LEVEL controls
// log verbosity at startup without code changes.
package config

import (
	"fmt"
	"log"
	"os"
	"strconv"
)

// Version is the build version, injected at build time via ldflags
// (-X vedo-edutrack/backend/internal/platform/config.Version=<version>).
// Defaults to "dev" for local builds (go run / go build without ldflags).
var Version = "dev"

// Config holds all application configuration loaded from the environment.
type Config struct {
	DatabaseURL string
	LogLevel    string
	Port        int
	JWKSURL     string
	// Runtime identity and URLs (dynamic env injection —
	// ADR-DES.INFRA.dynamic-config-injection).
	Environment   string // APP_ENV: development | staging | production
	PublicBaseURL string // PUBLIC_BASE_URL: externally reachable base URL (webhooks, absolute links)
	HubAPIURL     string // VEDO_HUB_API_URL: base URL of the VEDO Hub REST API (F0, read-only)
	// JWT issuer/audience for locally-issued dev tokens (T5).
	JWTIssuer   string
	JWTAudience string
}

// Load reads configuration from environment variables with defaults.
func Load() (*Config, error) {
	cfg := &Config{
		DatabaseURL:   envOrDefault("DATABASE_URL", "postgres://edutrack:edutrack@localhost:5432/edutrack?sslmode=disable"),
		LogLevel:      envOrDefault("LOG_LEVEL", "info"),
		JWKSURL:       envOrDefault("JWKS_URL", "http://localhost:8080/.well-known/jwks.json"),
		Environment:   envOrDefault("APP_ENV", "development"),
		PublicBaseURL: envOrDefault("PUBLIC_BASE_URL", "http://localhost:8080"),
		HubAPIURL:     envOrDefault("VEDO_HUB_API_URL", "http://localhost:8081"),
		JWTIssuer:     envOrDefault("JWT_ISSUER", "vedo-edutrack"),
		JWTAudience:   envOrDefault("JWT_AUDIENCE", "vedo-edutrack"),
	}

	portStr := envOrDefault("PORT", "8080")
	port, err := strconv.Atoi(portStr)
	if err != nil {
		return nil, fmt.Errorf("invalid PORT value %q: %w", portStr, err)
	}
	cfg.Port = port

	log.Printf("[INFO] [config.Load] configuration loaded (version=%s, env=%s, log_level=%s, port=%d)", Version, cfg.Environment, cfg.LogLevel, cfg.Port)
	return cfg, nil
}

func envOrDefault(key, defaultVal string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return defaultVal
}
