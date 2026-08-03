// Package api implements the HTTP API layer generated from the OpenAPI spec
// (backend/api/openapi/v1.yaml). Handlers are thin input adapters over the
// Application layer (ADR-DES.API.communication-patterns, OpenAPI-first).
package api

import (
	"encoding/json"
	"net/http"

	"go.uber.org/zap"

	ontostub "vedo-edutrack/backend/internal/modules/ontologyport/adapters/stub"
	routestub "vedo-edutrack/backend/internal/modules/routeplanning/adapters/stub"
	"vedo-edutrack/backend/internal/platform/auth"
)

// StubHandler implements ServerInterface. Ontology and route endpoints use
// the M0.3 stub implementations; the rest answer 501 Not Implemented.
type StubHandler struct {
	Ontology *ontostub.Graph
	Routes   *routestub.Computer
	Logger   *zap.Logger
}

// NewStubHandler returns a handler wired to the stub graph and route computer.
func NewStubHandler() *StubHandler {
	graph := ontostub.NewGraph()
	return &StubHandler{
		Ontology: graph,
		Routes:   routestub.NewComputer(graph),
		Logger:   zap.NewNop(),
	}
}

// notImplemented writes the standard 501 JSON body for a stub endpoint.
func notImplemented(w http.ResponseWriter, endpoint string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNotImplemented)
	msg := "endpoint not yet implemented (M0.3 scaffold)"
	_ = json.NewEncoder(w).Encode(ErrorResponse{
		Error:    "not_implemented",
		Message:  &msg,
		Endpoint: &endpoint,
	})
}

func (h *StubHandler) GetJwks(w http.ResponseWriter, _ *http.Request) {
	notImplemented(w, "GET /.well-known/jwks.json")
}

func (h *StubHandler) HealthCheck(w http.ResponseWriter, _ *http.Request) {
	notImplemented(w, "GET /healthz")
}

func (h *StubHandler) GetMe(w http.ResponseWriter, r *http.Request) {
	// Current user from the validated JWT context (T5).
	userID := ""
	roles := []string{}
	claims := claimsFromRequest(r)
	if claims != nil {
		userID = claims.UserID
		roles = claims.Roles
	}
	if userID == "" {
		writeAPIError(w, http.StatusUnauthorized, "unauthorized", "missing token claims")
		return
	}
	writeJSONBody(w, http.StatusOK, UserInfo{UserId: userID, Roles: roles})
}

func (h *StubHandler) GetOntologyConcept(w http.ResponseWriter, r *http.Request, params GetOntologyConceptParams) {
	topicID := params.TopicId

	h.Logger.Info("ontology concept query", zap.String("topic_id", topicID))
	concept := h.Ontology.GetConcept(topicID)
	if concept == nil {
		h.Logger.Warn("ontology topic not found", zap.String("topic_id", topicID))
		msg := "concept not found"
		writeJSONBody(w, http.StatusNotFound, ErrorResponse{Error: "not_found", Message: &msg})
		return
	}

	links := make([]ConceptLink, 0, len(concept.Links))
	for _, l := range concept.Links {
		links = append(links, ConceptLink{
			TopicId:  l.TopicID,
			LinkType: ConceptLinkLinkType(l.LinkType),
		})
	}

	resp := ConceptResponse{
		Concept: Concept{
			Id:          concept.ID,
			Title:       concept.Title,
			Description: &concept.Description,
			Links:       &links,
		},
	}
	writeJSONBody(w, http.StatusOK, resp)
}

func (h *StubHandler) ReadinessCheck(w http.ResponseWriter, _ *http.Request) {
	notImplemented(w, "GET /readyz")
}

func (h *StubHandler) ComputeRoute(w http.ResponseWriter, r *http.Request) {
	var req RouteComputeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		msg := "invalid request body"
		writeJSONBody(w, http.StatusBadRequest, ErrorResponse{Error: "invalid_request", Message: &msg})
		return
	}

	h.Logger.Info("route computation",
		zap.String("learner_id", req.LearnerId),
		zap.String("goal_topic_id", req.GoalTopicId),
	)

	route, err := h.Routes.ComputeRoute(req.LearnerId, req.GoalTopicId)
	if err != nil {
		h.Logger.Warn("route computation failed", zap.Error(err))
		msg := err.Error()
		writeJSONBody(w, http.StatusNotFound, ErrorResponse{Error: "topic_not_found", Message: &msg})
		return
	}

	topics := make([]RouterTopic, 0, len(route))
	for _, t := range route {
		lt := t.LinkType
		topics = append(topics, RouterTopic{
			TopicId:  t.TopicID,
			Order:    t.Order,
			LinkType: &lt,
		})
	}

	h.Logger.Info("route computed", zap.Int("topics", len(topics)))
	writeJSONBody(w, http.StatusOK, RouteComputeResponse{Route: topics})
}

// claimsFromRequest extracts the validated JWT claims from the request
// context (injected by the auth middleware, T5).
func claimsFromRequest(r *http.Request) *auth.Claims {
	return auth.ClaimsFrom(r.Context())
}

// writeJSONBody encodes any value as JSON with the given status.
func writeJSONBody(w http.ResponseWriter, status int, body interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(body)
}

// writeAPIError is a shorthand for an ErrorResponse with a message.
func writeAPIError(w http.ResponseWriter, status int, code, message string) {
	writeJSONBody(w, status, ErrorResponse{Error: code, Message: &message})
}
