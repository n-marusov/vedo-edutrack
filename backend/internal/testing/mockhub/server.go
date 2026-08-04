package mockhub

import (
	"net/http/httptest"
	"os"
	"testing"

	"go.uber.org/zap"
)

// NewTestServer starts an in-process mock VEDO Hub server over the ontology
// parsed from ttlPath. Used by Go integration/contract tests (M1) — the same
// handler as the hub-mock container.
//
// The caller must close the returned server.
func NewTestServer(t testing.TB, ttlPath string) *httptest.Server {
	t.Helper()

	f, err := os.Open(ttlPath) // #nosec G304
	if err != nil {
		t.Fatalf("open ontology %s: %v", ttlPath, err)
	}
	defer func() { _ = f.Close() }()

	ont, err := Parse(f)
	if err != nil {
		t.Fatalf("parse ontology %s: %v", ttlPath, err)
	}

	srv := httptest.NewServer(NewMux(ont, zap.NewNop()))
	t.Cleanup(srv.Close)
	return srv
}
