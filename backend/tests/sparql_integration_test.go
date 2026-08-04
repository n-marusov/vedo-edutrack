//go:build integration

package tests

import (
	"context"
	"encoding/json"
	"net/http/httptest"
	"testing"
	"time"

	"go.uber.org/zap"

	"vedo-edutrack/backend/internal/modules/integrations/adapters/sparql"
	"vedo-edutrack/backend/internal/platform/circuitbreaker"
	"vedo-edutrack/backend/internal/platform/ratelimit"
	"vedo-edutrack/backend/internal/testing/mockhub"
)

// TestSparqlEndToEndWithMockHub — M4 integration flow (T27b):
//
//	mockhub (in-process, traceability.ttl) → SPARQL proxy client → handler →
//	validated application/sparql-results+json with class bindings.
//
// Exercises the full EduTrack → Hub SPARQL round trip without a real Hub
// (the mock answers any read-only query with the TBox class list).
func TestSparqlEndToEndWithMockHub(t *testing.T) {
	// 1. Boot the mock VEDO Hub (same handler as the hub-mock container).
	hub := mockhub.NewTestServer(t, "../../traceability.ttl")
	defer hub.Close()

	// 2. SPARQL proxy client pointed at the mock.
	client, err := sparql.NewClient(sparql.ClientConfig{BaseURL: hub.URL, Path: "/sparql"}, zap.NewNop())
	if err != nil {
		t.Fatalf("sparql client: %v", err)
	}

	// 3. Handler with a generous rate limit and a circuit breaker.
	breaker := circuitbreaker.New(circuitbreaker.DefaultConfig(), zap.NewNop())
	limiter := ratelimit.NewLimiter(100, 100, zap.NewNop())
	handler := sparql.NewHandler(client, breaker, limiter, zap.NewNop())

	// 4. Execute a read-only query end-to-end.
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	req := httptest.NewRequest("GET", "/sparql?query="+"SELECT+%3Fs+WHERE+%7B+%3Fs+%3Fp+%3Fo+%7D", nil)
	rec := httptest.NewRecorder()
	handler.Query(rec, req.WithContext(ctx), "SELECT ?s WHERE { ?s ?p ?o }", "user-1")

	if rec.Code != 200 {
		t.Fatalf("status=%d, want 200; body=%s", rec.Code, rec.Body.String())
	}

	// 5. Validate the response format (application/sparql-results+json).
	if ct := rec.Header().Get("Content-Type"); ct != "application/sparql-results+json" {
		t.Fatalf("content-type=%q, want application/sparql-results+json", ct)
	}

	var body struct {
		Head struct {
			Vars []string `json:"vars"`
		} `json:"head"`
		Results struct {
			Bindings []map[string]any `json:"bindings"`
		} `json:"results"`
		Truncated bool `json:"truncated"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("decode sparql response: %v", err)
	}
	if len(body.Head.Vars) == 0 {
		t.Fatal("expected head.vars populated")
	}
	if len(body.Results.Bindings) == 0 {
		t.Fatal("expected class bindings from the mock ontology")
	}
	if body.Truncated {
		t.Fatal("expected truncated=false for the small mock ontology")
	}

	// 6. The mock rejects mutating queries (read-only boundary) with 403.
	rec = httptest.NewRecorder()
	handler.Query(rec, httptest.NewRequest("GET", "/sparql", nil).WithContext(ctx), "INSERT DATA { <http://x> <http://y> \"z\" }", "user-1")
	if rec.Code != 403 {
		t.Fatalf("mutation status=%d, want 403", rec.Code)
	}
}
