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
	"log"
	"os"
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
func run(args []string) error {
	log.Printf("[INFO] [vedo-edutrack] args=%v (CLI dispatch pending M0.3)", args)
	return nil
}
