// Package main is the entry point of the single vedo-edutrack binary.
//
// The binary exposes a Cobra command tree (see internal/cli): the `server`
// subcommand is the long-running process (HTTP API + SPARQL + webhooks + MCP
// over SSE), `mcp` serves MCP over stdio for AI agents, and the remaining
// subcommands (migrate, seed, ontology sync, route compute, plan get,
// gap diagnose, report) are dev/support/testing tools that reuse the same
// Application layer — see ADR-DES.API.cli-interface.
//
// Each subcommand wires its own minimal dependency graph via google/wire
// (per-command lazy wire); main stays a thin composition root.
//
// See ADR-DES.INFRA.modular-monolith-approach,
// ADR-DES.API.cli-interface and ADR-IMPL.PROCESS.repository-structure.
package main

import (
	"fmt"
	"log"
	"os"

	"vedo-edutrack/backend/internal/platform/config"
)

func main() {
	// TODO(M0.3): replace with cli.Execute() — cobra root (internal/cli).
	// Commands: server, mcp, migrate, seed, ontology sync, route compute,
	// plan get, gap diagnose, report.
	if err := run(os.Args[1:]); err != nil {
		log.Fatalf("[FATAL] [vedo-edutrack] %v", err)
	}
}

// run is the placeholder for the CLI dispatch (cobra Execute).
// Kept separate from main so it can be unit-tested.
//
// Functional subcommands at M0.2:
//
//	version  — print build version (injected via ldflags,
//	           ADR-DES.INFRA.dynamic-config-injection)
//	server   — HTTP server with health endpoints (/healthz, /readyz)
//	health   — liveness probe for container HEALTHCHECK (distroless: no curl)
func run(args []string) error {
	if len(args) == 0 {
		log.Printf("[INFO] [vedo-edutrack] no subcommand (CLI dispatch pending M0.3)")
		return nil
	}

	switch args[0] {
	case "version", "--version", "-v":
		fmt.Printf("vedo-edutrack %s (env %s)\n", config.Version, envOrDefault("APP_ENV", "development"))
		return nil

	case "server":
		cfg, err := config.Load()
		if err != nil {
			return err
		}
		return serveHTTP(cfg)

	case "health":
		return checkHealth()
	}

	log.Printf("[INFO] [vedo-edutrack] args=%v (CLI dispatch pending M0.3)", args)
	return nil
}

// envOrDefault mirrors platform/config's helper for the version banner
// (kept local to avoid an import cycle in the placeholder CLI).
func envOrDefault(key, defaultVal string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return defaultVal
}
