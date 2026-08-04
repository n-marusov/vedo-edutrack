package mockhub

import (
	"encoding/json"
	"net/http"
	"sort"
	"strings"

	"go.uber.org/zap"
)

// SparqlHandler serves a minimal read-only SPARQL JSON results endpoint over
// the in-memory TBox ontology (classes only — the mock has no individuals).
//
// The mock answers any read-only query form (SELECT/ASK/DESCRIBE/CONSTRUCT)
// with the class list as bindings (variable ?s = class IRI, ?label = class
// label), so integration tests exercise the full EduTrack → Hub SPARQL proxy
// round trip. Mutating forms are rejected with 403 — mirroring the read-only
// boundary of the real Hub.
//
// Deferred mock channel (RESEARCH: "REST/SPARQL/MCP mock channels: deferred —
// add only when the adapter calls them") — now added because the M4 SPARQL
// proxy adapter calls it.
type SparqlHandler struct {
	ont    *Ontology
	logger *zap.Logger
}

// NewSparqlHandler builds the SPARQL mock handler.
func NewSparqlHandler(ont *Ontology, logger *zap.Logger) http.Handler {
	if logger == nil {
		logger = zap.NewNop()
	}
	return &SparqlHandler{ont: ont, logger: logger}
}

// ServeHTTP handles GET /sparql?query=... and POST /sparql (form body).
func (h *SparqlHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet && r.Method != http.MethodPost {
		http.Error(w, `{"error":"method_not_allowed"}`, http.StatusMethodNotAllowed)
		return
	}
	// Auth is optional on the mock (dev flow has no VEDO_HUB_BEARER_TOKEN);
	// when present it is accepted as-is.
	if auth := r.Header.Get("Authorization"); auth != "" && !strings.HasPrefix(auth, "Bearer ") {
		h.writeJSON(w, http.StatusUnauthorized, map[string]any{
			"errors": []any{map[string]any{"message": "invalid authorization header"}},
		})
		return
	}

	query := r.URL.Query().Get("query")
	if query == "" && r.Method == http.MethodPost {
		query = r.PostFormValue("query")
	}
	query = strings.TrimSpace(query)
	if query == "" {
		h.writeJSON(w, http.StatusBadRequest, map[string]any{
			"errors": []any{map[string]any{"message": "query is required"}},
		})
		return
	}

	// Read-only gate: only SELECT/ASK/DESCRIBE/CONSTRUCT.
	if !isReadOnlySparql(query) {
		h.writeJSON(w, http.StatusForbidden, map[string]any{
			"errors": []any{map[string]any{"message": "mutating or unsupported query form"}},
		})
		return
	}

	w.Header().Set("Content-Type", "application/sparql-results+json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(h.results())
}

// isReadOnlySparql reports whether the query starts with a read-only form.
func isReadOnlySparql(query string) bool {
	lower := strings.ToLower(query)
	for _, prefix := range []string{"select", "ask", "describe", "construct"} {
		if strings.HasPrefix(lower, prefix) {
			return true
		}
	}
	return false
}

// results builds a SPARQL JSON results document over all classes.
func (h *SparqlHandler) results() map[string]any {
	ids := make([]string, 0, len(h.ont.Classes))
	for id := range h.ont.Classes {
		ids = append(ids, id)
	}
	sort.Strings(ids)

	bindings := make([]any, 0, len(ids))
	for _, id := range ids {
		c := h.ont.Classes[id]
		bindings = append(bindings, map[string]any{
			"s":     map[string]any{"type": "uri", "value": id},
			"label": map[string]any{"type": "literal", "value": c.LabelOrID()},
		})
	}
	return map[string]any{
		"head": map[string]any{
			"vars": []string{"s", "label"},
		},
		"results": map[string]any{
			"bindings": bindings,
		},
	}
}

func (h *SparqlHandler) writeJSON(w http.ResponseWriter, status int, body any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(body)
}
