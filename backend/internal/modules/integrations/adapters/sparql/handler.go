package sparql

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"time"

	"go.uber.org/zap"

	"vedo-edutrack/backend/internal/platform/circuitbreaker"
	"vedo-edutrack/backend/internal/platform/ratelimit"
)

// Handler is the production SPARQL endpoint handler (F6): read-only gate,
// per-user rate limiting, circuit-breaker-wrapped Hub proxy, result
// truncation, and the full error contract
// (REQ-FR-api.sparql.read-only, REQ-NFR-security.compliance.owasp-application-security).
type Handler struct {
	client  *Client
	breaker *circuitbreaker.Breaker
	limiter *ratelimit.Limiter
	logger  *zap.Logger
}

// NewHandler builds the SPARQL endpoint handler. When limiter is nil, rate
// limiting is disabled (tests); when breaker is nil, calls are not guarded.
func NewHandler(client *Client, breaker *circuitbreaker.Breaker, limiter *ratelimit.Limiter, logger *zap.Logger) *Handler {
	if logger == nil {
		logger = zap.NewNop()
	}
	return &Handler{client: client, breaker: breaker, limiter: limiter, logger: logger.Named("sparql")}
}

// QueryResult is the truncated result document returned to the client.
type QueryResult struct {
	ResultSet
	Truncated bool
}

// Query handles a SPARQL query request. userID is the authenticated user
// ("" when unauthenticated — the auth middleware normally rejects those with
// 401 before reaching the handler).
func (h *Handler) Query(w http.ResponseWriter, r *http.Request, query string, userID string) {
	h.logger.Debug("QueryReceived", zap.String("userID", userID), zap.Int("queryLength", len(query)))

	// Read-only gate: empty query → 400; mutation → 403.
	if err := ValidateReadOnly(query); err != nil {
		switch {
		case errors.Is(err, ErrMutationNotAllowed):
			h.logger.Warn("QueryRejected", zap.String("userID", userID), zap.String("reason", "mutation"))
			writeJSONError(w, http.StatusForbidden, "sparql_mutation_not_allowed", "only SELECT/ASK/DESCRIBE/CONSTRUCT queries are permitted (read-only endpoint)")
			return
		default:
			h.logger.Warn("QueryRejected", zap.String("userID", userID), zap.String("reason", "empty_or_malformed"))
			writeJSONError(w, http.StatusBadRequest, "invalid_sparql", "query is required and must be a non-empty string")
			return
		}
	}

	// Rate limit per user (10 req/min for SPARQL, burst 2 — configured at
	// construction via the limiter).
	if h.limiter != nil {
		key := userID
		if key == "" {
			key = "anonymous"
		}
		allowed, retryAfter := h.limiter.Allow(key)
		if !allowed {
			h.logger.Warn("QueryRejected", zap.String("userID", userID), zap.String("reason", "rate_limit"))
			w.Header().Set("Retry-After", strconv.Itoa(retryAfter))
			writeJSONError(w, http.StatusTooManyRequests, "rate_limited", "too many requests")
			return
		}
	}

	// Execute against the Hub through the circuit breaker. Timeout (30s) maps
	// to 504; a permanently-open circuit maps to 503.
	started := time.Now()
	result, err := h.execute(r.Context(), query)
	if err != nil {
		switch {
		case errors.Is(err, circuitbreaker.ErrOpen):
			SparqlQueryDurationSeconds.WithLabelValues("hub_unavailable").Observe(time.Since(started).Seconds())
			h.logger.Warn("QueryRejected", zap.String("userID", userID), zap.String("reason", "hub_down"))
			writeJSONError(w, http.StatusServiceUnavailable, "hub_unavailable", "ontology service is temporarily unavailable")
			return
		case errors.Is(err, errHubTimeout):
			SparqlQueryDurationSeconds.WithLabelValues("timeout").Observe(time.Since(started).Seconds())
			h.logger.Warn("QueryRejected", zap.String("userID", userID), zap.String("reason", "timeout"))
			writeJSONError(w, http.StatusGatewayTimeout, "sparql_timeout", "query execution timed out")
			return
		default:
			SparqlQueryDurationSeconds.WithLabelValues("error").Observe(time.Since(started).Seconds())
			h.logger.Error("HubUnreachable", zap.Error(err))
			writeJSONError(w, http.StatusBadGateway, "sparql_execution_failed", "failed to execute query against the ontology service")
			return
		}
	}
	SparqlQueryDurationSeconds.WithLabelValues("success").Observe(time.Since(started).Seconds())

	h.logger.Info("QueryExecuted",
		zap.String("userID", userID),
		zap.Duration("executionTime", time.Since(started)),
		zap.Int("resultRows", len(result.Results.Bindings)),
		zap.Bool("truncated", result.Truncated),
	)

	w.Header().Set("Content-Type", "application/sparql-results+json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(result.toWire())
}

// errHubTimeout is the internal timeout signal (wrapped client timeout).
var errHubTimeout = errors.New("sparql query timed out")

// toWire converts the result to the OpenAPI SparqlResponse shape.
func (q *QueryResult) toWire() map[string]any {
	vars := q.Head.Vars
	bindings := q.Results.Bindings
	if vars == nil {
		vars = []string{}
	}
	if bindings == nil {
		bindings = []map[string]any{}
	}
	return map[string]any{
		"head":      map[string]any{"vars": vars},
		"results":   map[string]any{"bindings": bindings},
		"truncated": q.Truncated,
	}
}

// execute runs the query through the breaker and truncates the result.
func (h *Handler) execute(ctx context.Context, query string) (*QueryResult, error) {
	if h.client == nil {
		return nil, errors.New("sparql client is not configured")
	}
	if h.breaker != nil {
		allowed, err := h.breaker.Allow()
		if err != nil {
			return nil, err
		}
		if !allowed {
			return nil, circuitbreaker.ErrOpen
		}
	}

	result, err := h.client.Query(ctx, query)
	if err != nil {
		if h.breaker != nil {
			h.breaker.Failure()
		}
		if errors.Is(err, ErrTimeout) || ctx.Err() != nil {
			return nil, errHubTimeout
		}
		return nil, err
	}
	if h.breaker != nil {
		h.breaker.Success()
	}

	bindings, truncated := TruncateResult(result.Results.Bindings, MaxResultRows)
	result.Results.Bindings = bindings
	return &QueryResult{ResultSet: *result, Truncated: truncated}, nil
}

// writeJSONError writes an ErrorResponse-shaped JSON error.
func writeJSONError(w http.ResponseWriter, status int, code, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(map[string]string{
		"error":   code,
		"message": message,
	})
}
