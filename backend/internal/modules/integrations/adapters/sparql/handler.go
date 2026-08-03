// Package sparql implements the read-only SPARQL endpoint guard (F6).
//
// The EduTrack SPARQL surface is read-only by contract
// (REQ-FR-api.sparql.read-only): only SELECT / ASK / DESCRIBE / CONSTRUCT are
// permitted. Mutating queries (INSERT, DELETE, LOAD, CLEAR, CREATE, DROP) are
// rejected before they reach the underlying store.
package sparql

import (
	"fmt"
	"strings"
)

// readOnlyPrefixes are the SPARQL query forms allowed by the read-only gate.
var readOnlyPrefixes = []string{"select", "ask", "describe", "construct"}

// MaxResultRows is the SPARQL result truncation boundary (NFR-FR-api.sparql.read-only).
const MaxResultRows = 10000

// ValidateReadOnly returns an error when the query is mutating or malformed.
func ValidateReadOnly(query string) error {
	trimmed := strings.TrimSpace(query)
	if trimmed == "" {
		return fmt.Errorf("empty SPARQL query")
	}
	lower := strings.ToLower(trimmed)
	for _, prefix := range readOnlyPrefixes {
		if strings.HasPrefix(lower, prefix) {
			return nil
		}
	}
	return fmt.Errorf("mutating or unsupported SPARQL query form (only SELECT/ASK/DESCRIBE/CONSTRUCT allowed)")
}

// IsTooManyResults reports whether a result set exceeds the truncation boundary.
func IsTooManyResults(rowCount, maxRows int) bool {
	if maxRows <= 0 {
		maxRows = MaxResultRows
	}
	return rowCount > maxRows
}
