// Package cli implements the cobra command tree of the vedo-edutrack binary.
//
// The CLI is an input adapter over the Application layer — the same pattern
// as the HTTP handler and the MCP server (ADR-DES.API.cli-interface): every
// command calls module use cases through wire providers; there is no second
// path to data.
//
// Commands:
//
//	server         long-running process (HTTP API + SPARQL + webhooks + MCP SSE)
//	mcp            MCP server over stdio for AI agents (F6.6)
//	migrate        Atlas migrations up/down/validate (drift = 0)
//	seed           RBAC role catalog + demo data
//	ontology sync  copy subgraph from VEDO Hub (F0.2)
//	route compute  compute a route (--stub | from DB) — dev/test tool
//	plan get       read plan / progress
//	gap diagnose   root-cause diagnosis of learner lag
//	report         attestation / FGOS coverage reports to file (batch)
//
// Each command builds its own minimal wire graph (per-command lazy wire) in
// wire.go. Commands are scriptable (no interactive prompts), emit
// structured audit logs via zap, and support --output json|table|csv.
package cli
