// Package tests holds cross-module integration tests (Go, testcontainers).
// See ADR-IMPL.PROCESS.repository-structure §5 — backend/tests runs against a
// real PostgreSQL via testcontainers and exercises domain events/transactions
// across bounded contexts.
package tests
