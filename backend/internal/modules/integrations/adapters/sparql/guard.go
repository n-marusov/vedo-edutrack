// Package sparql implements the read-only SPARQL endpoint guard and the
// VEDO Hub SPARQL proxy (F6).
//
// The EduTrack SPARQL surface is read-only by contract
// (REQ-FR-api.sparql.read-only): only SELECT / ASK / DESCRIBE / CONSTRUCT are
// permitted. Mutating queries (INSERT, DELETE, LOAD, CLEAR, CREATE, DROP) are
// rejected before they reach the underlying store.
package sparql

import (
	"errors"
	"strings"
)

// readOnlyPrefixes are the SPARQL query forms allowed by the read-only gate.
var readOnlyPrefixes = []string{"select", "ask", "describe", "construct"}

// MaxResultRows is the SPARQL result truncation boundary (NFR-FR-api.sparql.read-only).
const MaxResultRows = 10000

// ErrEmptyQuery is returned for an empty/whitespace query (HTTP 400).
var ErrEmptyQuery = errors.New("empty SPARQL query")

// ErrMutationNotAllowed is returned for a mutating or unsupported query form
// (HTTP 403 — the operation is forbidden by the read-only contract).
var ErrMutationNotAllowed = errors.New("mutating or unsupported SPARQL query form (only SELECT/ASK/DESCRIBE/CONSTRUCT allowed)")

// ValidateReadOnly returns an error when the query is mutating or malformed.
// The returned error is either ErrEmptyQuery or ErrMutationNotAllowed — callers
// map them to distinct HTTP statuses (400 vs 403).
func ValidateReadOnly(query string) error {
	trimmed := strings.TrimSpace(query)
	if trimmed == "" {
		return ErrEmptyQuery
	}
	lower := strings.ToLower(trimmed)
	for _, prefix := range readOnlyPrefixes {
		if strings.HasPrefix(lower, prefix) {
			return nil
		}
	}
	return ErrMutationNotAllowed
}

// IsTooManyResults reports whether a result set exceeds the truncation boundary.
func IsTooManyResults(rowCount, maxRows int) bool {
	if maxRows <= 0 {
		maxRows = MaxResultRows
	}
	return rowCount > maxRows
}

// TruncateResult clips a result set to maxRows rows and reports whether
// truncation happened.
func TruncateResult(bindings []map[string]any, maxRows int) ([]map[string]any, bool) {
	if maxRows <= 0 {
		maxRows = MaxResultRows
	}
	if len(bindings) <= maxRows {
		return bindings, false
	}
	return bindings[:maxRows], true
}
