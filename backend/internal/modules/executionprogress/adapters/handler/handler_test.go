//go:build integration

package handler

import "testing"

// TODO(M0.3): HTTP handler integration tests (testcontainers + PostgreSQL).
// The executionprogress HTTP handler is not implemented yet (package has only
// doc.go) — re-enable with real tests when the handler lands (T18).
func TestExecutionprogressHandlerIntegration(t *testing.T) {
	t.Skip("TODO(M0.3): executionprogress handler not implemented yet")
}
