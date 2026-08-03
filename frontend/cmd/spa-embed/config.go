package main

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

// Config holds all runtime configuration sourced from environment variables.
// Implements ADR-DES.INFRA.dynamic-config-injection: one build, many environments.
type Config struct {
	Port    string
	Version string
	Env     string

	// API base URL injected into /config.js at runtime.
	APIBaseURL string

	// Rate limiting
	RateLimitPerSec int
	RateLimitBurst  int

	// Timeouts
	ReadTimeout       time.Duration
	ReadHeaderTimeout time.Duration
	WriteTimeout      time.Duration
	IdleTimeout       time.Duration

	// Metrics
	MetricsEnabled bool
}

// LoadConfig reads configuration from environment variables with sensible defaults.
func LoadConfig() Config {
	return Config{
		Port:              envOr("PORT", "8080"),
		Version:           envOr("APP_VERSION", "dev"),
		Env:               envOr("APP_ENV", "production"),
		APIBaseURL:        envOr("API_BASE_URL", "/api"),
		RateLimitPerSec:   envInt("RATE_LIMIT_PER_SEC", 100),
		RateLimitBurst:    envInt("RATE_LIMIT_BURST", 200),
		ReadTimeout:       envDuration("READ_TIMEOUT", 5*time.Second),
		ReadHeaderTimeout: envDuration("READ_HEADER_TIMEOUT", 3*time.Second),
		WriteTimeout:      envDuration("WRITE_TIMEOUT", 10*time.Second),
		IdleTimeout:       envDuration("IDLE_TIMEOUT", 120*time.Second),
		MetricsEnabled:    envBool("METRICS_ENABLED", true),
	}
}

func envOr(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

func envInt(key string, def int) int {
	if v := os.Getenv(key); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			return n
		}
	}
	return def
}

func envDuration(key string, def time.Duration) time.Duration {
	if v := os.Getenv(key); v != "" {
		if d, err := time.ParseDuration(v); err == nil {
			return d
		}
	}
	return def
}

func envBool(key string, def bool) bool {
	if v := os.Getenv(key); v != "" {
		switch v {
		case "1", "true", "yes", "on":
			return true
		case "0", "false", "no", "off":
			return false
		}
	}
	return def
}

// ConfigJS returns the JavaScript runtime config served at /config.js.
func (c Config) ConfigJS() string {
	return fmt.Sprintf(
		"window.APP_CONFIG = { version: %q, env: %q, apiBaseUrl: %q };\n",
		c.Version, c.Env, c.APIBaseURL,
	)
}
