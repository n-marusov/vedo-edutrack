package mockhub

import (
	"net/http"

	"go.uber.org/zap"
)

// NewMux routes the mock VEDO Hub HTTP surface:
//
//	GET  /healthz  → liveness (container healthcheck)
//	POST /graphql  → ontology GraphQL service (F0)
//	GET/POST /sparql → read-only SPARQL JSON results (F6, M4)
//
// The same mux is used by the hub-mock container and the in-process test
// server, so dev/test/CI behavior is identical.
func NewMux(ont *Ontology, logger *zap.Logger) http.Handler {
	if logger == nil {
		logger = zap.NewNop()
	}
	mux := http.NewServeMux()
	mux.Handle("/graphql", NewHandler(ont, logger))
	mux.Handle("/sparql", NewSparqlHandler(ont, logger))
	mux.Handle("/healthz", http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"status":"ok"}`))
	}))
	return mux
}
