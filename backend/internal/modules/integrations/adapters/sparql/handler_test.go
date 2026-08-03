package sparql

import (
	"testing"
)

func TestReadOnlyQueryGateRejectsMutations(t *testing.T) {
	mutations := []string{
		"INSERT DATA { <http://x> <http://y> \"z\" }",
		"DELETE WHERE { ?s ?p ?o }",
		"CREATE GRAPH <http://g>",
		"DROP GRAPH <http://g>",
		"LOAD <http://file>",
		"CLEAR ALL",
	}
	for _, query := range mutations {
		if err := ValidateReadOnly(query); err == nil {
			t.Fatalf("expected mutation rejected: %q", query)
		}
	}
}

func TestReadOnlyQueryGateAllowsSelectAskDescribe(t *testing.T) {
	queries := []string{
		"SELECT ?s WHERE { ?s ?p ?o }",
		"ASK WHERE { ?s ?p ?o }",
		"DESCRIBE <http://example.org/s>",
		"SELECT DISTINCT ?c WHERE { ?c a <http://example.org/Class> } LIMIT 10",
		"select ?s where { ?s ?p ?o }", // case-insensitive
	}
	for _, query := range queries {
		if err := ValidateReadOnly(query); err != nil {
			t.Fatalf("expected read-only query allowed: %q: %v", query, err)
		}
	}
}

func TestReadOnlyQueryGateTruncationFlag(t *testing.T) {
	if !IsTooManyResults(10001, 10000) {
		t.Fatal("expected truncated flag for >10k rows")
	}
	if IsTooManyResults(9999, 10000) {
		t.Fatal("expected no truncation under 10k rows")
	}
}
