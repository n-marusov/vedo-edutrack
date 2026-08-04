package sparql

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	"go.uber.org/zap"

	"vedo-edutrack/backend/internal/platform/circuitbreaker"
	"vedo-edutrack/backend/internal/platform/ratelimit"
)

// mockHubServer serves a fake SPARQL JSON results document.
func mockHubServer(t *testing.T, payload string, status int) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Query().Get("query") == "" {
			t.Errorf("expected query param")
		}
		w.Header().Set("Content-Type", "application/sparql-results+json")
		w.WriteHeader(status)
		if payload != "" {
			_, _ = w.Write([]byte(payload))
		}
	}))
}

func newTestHandler(t *testing.T, hubURL string) *Handler {
	t.Helper()
	client, err := NewClient(ClientConfig{BaseURL: hubURL, Path: "/sparql"}, zap.NewNop())
	if err != nil {
		t.Fatalf("new client: %v", err)
	}
	breaker := circuitbreaker.New(circuitbreaker.DefaultConfig(), zap.NewNop())
	limiter := ratelimit.NewLimiter(100, 10, zap.NewNop()) // generous for tests
	return NewHandler(client, breaker, limiter, zap.NewNop())
}

func doQuery(t *testing.T, h *Handler, query, userID string) *httptest.ResponseRecorder {
	t.Helper()
	req := httptest.NewRequest(http.MethodGet, "/sparql?query="+url.QueryEscape(query), nil)
	rec := httptest.NewRecorder()
	h.Query(rec, req, query, userID)
	return rec
}

func TestHandlerValidQueryReturnsResults(t *testing.T) {
	payload := `{"head":{"vars":["s"]},"results":{"bindings":[{"s":{"type":"uri","value":"http://x"}}]}}`
	hub := mockHubServer(t, payload, http.StatusOK)
	defer hub.Close()

	h := newTestHandler(t, hub.URL)
	rec := doQuery(t, h, "SELECT ?s WHERE { ?s ?p ?o }", "user-1")

	if rec.Code != http.StatusOK {
		t.Fatalf("status=%d, want 200; body=%s", rec.Code, rec.Body.String())
	}
	if ct := rec.Header().Get("Content-Type"); !strings.Contains(ct, "application/sparql-results+json") {
		t.Fatalf("content-type=%q, want sparql-results+json", ct)
	}
	var body map[string]any
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if body["truncated"] != false {
		t.Fatalf("truncated=%v, want false", body["truncated"])
	}
}

func TestHandlerMutationRejectedWith403(t *testing.T) {
	h := newTestHandler(t, "http://127.0.0.1:1") // unreachable; mutation rejected before network
	rec := doQuery(t, h, "INSERT DATA { <http://x> <http://y> \"z\" }", "user-1")
	if rec.Code != http.StatusForbidden {
		t.Fatalf("status=%d, want 403; body=%s", rec.Code, rec.Body.String())
	}
}

func TestHandlerEmptyQueryRejectedWith400(t *testing.T) {
	h := newTestHandler(t, "http://127.0.0.1:1")
	rec := doQuery(t, h, "   ", "user-1")
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status=%d, want 400; body=%s", rec.Code, rec.Body.String())
	}
}

func TestHandlerAllowsAllReadOnlyQueryForms(t *testing.T) {
	payload := `{"head":{"vars":["s"]},"results":{"bindings":[]}}`
	for _, query := range []string{
		"SELECT ?s WHERE { ?s ?p ?o }",
		"ASK WHERE { ?s ?p ?o }",
		"DESCRIBE <http://example.org/s>",
		"CONSTRUCT { ?s ?p ?o } WHERE { ?s ?p ?o }",
	} {
		hub := mockHubServer(t, payload, http.StatusOK)
		h := newTestHandler(t, hub.URL)
		rec := doQuery(t, h, query, "user-1")
		hub.Close()
		if rec.Code != http.StatusOK {
			t.Fatalf("query %q: status=%d, want 200; body=%s", query, rec.Code, rec.Body.String())
		}
	}
}

func TestHandlerHubTimeoutReturns504(t *testing.T) {
	// A hub that sleeps longer than the client timeout → 504.
	hub := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		time.Sleep(2 * time.Second)
		w.WriteHeader(http.StatusOK)
	}))
	defer hub.Close()

	client, err := NewClient(ClientConfig{BaseURL: hub.URL, Timeout: 100 * time.Millisecond}, zap.NewNop())
	if err != nil {
		t.Fatalf("new client: %v", err)
	}
	h := NewHandler(client, nil, nil, zap.NewNop())

	rec := doQuery(t, h, "SELECT ?s WHERE { ?s ?p ?o }", "user-1")
	if rec.Code != http.StatusGatewayTimeout {
		t.Fatalf("status=%d, want 504; body=%s", rec.Code, rec.Body.String())
	}
}

func TestHandlerRateLimitReturns429(t *testing.T) {
	hub := mockHubServer(t, `{"head":{"vars":[]},"results":{"bindings":[]}}`, http.StatusOK)
	defer hub.Close()

	client, err := NewClient(ClientConfig{BaseURL: hub.URL}, zap.NewNop())
	if err != nil {
		t.Fatalf("new client: %v", err)
	}
	strict := ratelimit.NewLimiter(0.000001, 1, zap.NewNop()) // 1 token total
	h := NewHandler(client, nil, strict, zap.NewNop())

	if rec := doQuery(t, h, "SELECT ?s WHERE { ?s ?p ?o }", "user-1"); rec.Code != http.StatusOK {
		t.Fatalf("first: status=%d, want 200", rec.Code)
	}
	rec := doQuery(t, h, "SELECT ?s WHERE { ?s ?p ?o }", "user-1")
	if rec.Code != http.StatusTooManyRequests {
		t.Fatalf("second: status=%d, want 429", rec.Code)
	}
	if rec.Header().Get("Retry-After") == "" {
		t.Fatal("expected Retry-After header on 429")
	}
}

func TestHandlerTruncatesLargeResultSet(t *testing.T) {
	var sb strings.Builder
	sb.WriteString(`{"head":{"vars":["s"]},"results":{"bindings":[`)
	for i := 0; i < MaxResultRows+5; i++ {
		if i > 0 {
			sb.WriteString(",")
		}
		sb.WriteString(fmt.Sprintf(`{"s":{"type":"uri","value":"http://x/%d"}}`, i))
	}
	sb.WriteString(`]}}`)
	hub := mockHubServer(t, sb.String(), http.StatusOK)
	defer hub.Close()

	h := newTestHandler(t, hub.URL)
	rec := doQuery(t, h, "SELECT ?s WHERE { ?s ?p ?o }", "user-1")
	if rec.Code != http.StatusOK {
		t.Fatalf("status=%d, want 200", rec.Code)
	}
	var body struct {
		Truncated bool `json:"truncated"`
		Results   struct {
			Bindings []map[string]any `json:"bindings"`
		} `json:"results"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if !body.Truncated {
		t.Fatal("expected truncated=true for >10k rows")
	}
	if len(body.Results.Bindings) != MaxResultRows {
		t.Fatalf("bindings=%d, want %d", len(body.Results.Bindings), MaxResultRows)
	}
}

func TestHandlerCircuitOpenReturns503(t *testing.T) {
	hub := mockHubServer(t, `{"head":{"vars":[]},"results":{"bindings":[]}}`, http.StatusOK)
	defer hub.Close()

	client, err := NewClient(ClientConfig{BaseURL: hub.URL}, zap.NewNop())
	if err != nil {
		t.Fatalf("new client: %v", err)
	}
	breaker := circuitbreaker.New(circuitbreaker.Config{FailureThreshold: 1, Timeout: time.Minute, SuccessThreshold: 1}, zap.NewNop())
	breaker.Failure() // trip the breaker
	h := NewHandler(client, breaker, nil, zap.NewNop())

	rec := doQuery(t, h, "SELECT ?s WHERE { ?s ?p ?o }", "user-1")
	if rec.Code != http.StatusServiceUnavailable {
		t.Fatalf("status=%d, want 503; body=%s", rec.Code, rec.Body.String())
	}
}

func TestHandlerHubDownReturns502(t *testing.T) {
	// A closed listener (dial fails immediately).
	hub := httptest.NewServer(http.HandlerFunc(func(http.ResponseWriter, *http.Request) {}))
	hubURL := hub.URL
	hub.Close()

	h := newTestHandler(t, hubURL)
	rec := doQuery(t, h, "SELECT ?s WHERE { ?s ?p ?o }", "user-1")
	if rec.Code != http.StatusBadGateway {
		t.Fatalf("status=%d, want 502; body=%s", rec.Code, rec.Body.String())
	}
}

func TestTruncateResultHelper(t *testing.T) {
	rows := []map[string]any{{"a": 1}, {"b": 2}, {"c": 3}}
	clipped, truncated := TruncateResult(rows, 2)
	if !truncated || len(clipped) != 2 {
		t.Fatalf("truncated=%t len=%d, want true/2", truncated, len(clipped))
	}
	if _, truncated := TruncateResult(rows, 10); truncated {
		t.Fatal("expected no truncation under the boundary")
	}
}

func TestClientRejectsEmptyBaseURL(t *testing.T) {
	if _, err := NewClient(ClientConfig{}, zap.NewNop()); err == nil {
		t.Fatal("expected error for empty base URL")
	}
}

func TestClientTimeout(t *testing.T) {
	// A server that sleeps longer than the client timeout.
	hub := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		time.Sleep(2 * time.Second)
		w.WriteHeader(http.StatusOK)
	}))
	defer hub.Close()

	client, err := NewClient(ClientConfig{BaseURL: hub.URL, Timeout: 100 * time.Millisecond}, zap.NewNop())
	if err != nil {
		t.Fatalf("new client: %v", err)
	}
	if _, err := client.Query(httptest.NewRequest(http.MethodGet, "/", nil).Context(), "SELECT ?s WHERE {}"); err == nil {
		t.Fatal("expected timeout error")
	}
}
