// Package config provides environment-based configuration loading.
//
// Configuration is read from environment variables (12-factor app style),
// with sensible defaults for local development. See .env.example for
// documented variables.
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

// Config holds all application configuration loaded from the environment.
type Config struct {
	DatabaseURL string
	LogLevel    string
	Port        int
	JWKSURL     string
}

// Load reads configuration from environment variables with defaults.
func Load() (*Config, error) {
	cfg := &Config{
		DatabaseURL: envOrDefault("DATABASE_URL", "postgres://edutrack:edutrack@localhost:5432/edutrack?sslmode=disable"),
		LogLevel:    envOrDefault("LOG_LEVEL", "info"),
		JWKSURL:     envOrDefault("JWKS_URL", "http://localhost:8080/.well-known/jwks.json"),
	}

	portStr := envOrDefault("PORT", "8080")
	port, err := strconv.Atoi(portStr)
	if err != nil {
		return nil, fmt.Errorf("invalid PORT value %q: %w", portStr, err)
	}
	cfg.Port = port

	log.Printf("[INFO] [config.Load] configuration loaded (log_level=%s, port=%d)", cfg.LogLevel, cfg.Port)
	return cfg, nil
}

func envOrDefault(key, defaultVal string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return defaultVal
}
