package main

import (
	"os"
)

// main is the entry point for the spa-embed server.
//
// Subcommands:
//
//	spa-embed          — start the HTTP server (serves SPA + /healthz + /config.js)
//	spa-embed health   — self-health HTTP probe (container HEALTHCHECK, distroless)
//
// Build:
//
//	CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o spa-embed .
func main() {
	cfg := LoadConfig()

	if len(os.Args) > 1 && os.Args[1] == "health" {
		HealthCheck(cfg.Port)
		return
	}

	Run(cfg)
}
