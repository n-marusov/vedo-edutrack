// Package api implements the HTTP API layer generated from the OpenAPI spec
// (backend/api/openapi/v1.yaml). Handlers are thin input adapters over the
// Application layer (ADR-DES.API.communication-patterns, OpenAPI-first).
package api

import (
	"encoding/json"
	"net/http"

	"go.uber.org/zap"

	gapapp "vedo-edutrack/backend/internal/modules/gapcoverage/application"
	gapdomain "vedo-edutrack/backend/internal/modules/gapcoverage/domain"
	"vedo-edutrack/backend/internal/modules/integrations/adapters/sparql"
	ontostub "vedo-edutrack/backend/internal/modules/ontologyport/adapters/stub"
	resourceapp "vedo-edutrack/backend/internal/modules/resources/application"
	resourcedomain "vedo-edutrack/backend/internal/modules/resources/domain"
	routestub "vedo-edutrack/backend/internal/modules/routeplanning/adapters/stub"
	"vedo-edutrack/backend/internal/platform/auth"
)

// StubHandler implements ServerInterface. M0.3 stub endpoints (ontology,
// routes) are kept for compatibility; M1 endpoints (resources, gap, coverage)
// are wired to real application services over in-memory domain fixtures until
// persistence lands in the DB-backed phase.
type StubHandler struct {
	Ontology  *ontostub.Graph
	Routes    *routestub.Computer
	Resources *resourceapp.CatalogService
	Gap       *gapapp.GapService
	Coverage  *gapapp.CoverageService
	Logger    *zap.Logger
	gapGraph  gapdomain.Graph
}

// NewStubHandler returns a handler wired to the stub graph, route computer,
// and real M1 application services over in-memory fixtures.
func NewStubHandler() *StubHandler {
	graph := ontostub.NewGraph()
	logger := zap.NewNop()

	catalog, _ := resourcedomain.NewCatalog([]resourcedomain.Resource{
		{ID: "res-1", Title: "Percent video", Type: resourcedomain.ResourceTypeContent, Format: "video", Source: "school", Difficulty: "basic", DurationMinutes: 10, URI: "https://example.test/res-1"},
		{ID: "res-2", Title: "Percent text", Type: resourcedomain.ResourceTypeContent, Format: "text", Source: "school", Difficulty: "basic", DurationMinutes: 30, URI: "https://example.test/res-2"},
		{ID: "res-3", Title: "Lab access", Type: resourcedomain.ResourceTypeEnabling, Format: "lab", Source: "partner", Difficulty: "advanced", DurationMinutes: 60},
	})
	_ = catalog.BindToModule(resourcedomain.ResourceBinding{ResourceID: "res-1", ModuleID: "math-5-11", LinkType: resourcedomain.LinkAppliesTo})

	gapGraph := gapdomain.Graph{Modules: []gapdomain.Module{{ID: "percent", Title: "Проценты"}, {ID: "solutions", Title: "Растворы"}, {ID: "chemistry", Title: "Химия"}}, Links: []gapdomain.Link{
		{SourceID: "percent", TargetID: "solutions", Type: gapdomain.LinkStrictPrerequisite},
		{SourceID: "solutions", TargetID: "chemistry", Type: gapdomain.LinkStrictPrerequisite},
	}}

	return &StubHandler{
		Ontology:  graph,
		Routes:    routestub.NewComputer(graph),
		Resources: resourceapp.NewCatalogService(catalog, logger),
		Gap:       gapapp.NewGapService(logger),
		Coverage:  gapapp.NewCoverageService(logger),
		Logger:    logger,
		gapGraph:  gapGraph,
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

// GetFgosCoverage returns live FGOS coverage computed from the in-memory fixture.
func (h *StubHandler) GetFgosCoverage(w http.ResponseWriter, _ *http.Request, learnerID string) {
	h.Logger.Info("FGOS coverage query", zap.String("learner_id", learnerID))
	bindings := []gapdomain.FgosBinding{
		{ModuleID: "percent", RequirementID: "fgos-math-5"},
		{ModuleID: "solutions", RequirementID: "fgos-chem-8"},
		{ModuleID: "chemistry", RequirementID: "fgos-chem-8-1"},
	}
	report := h.Coverage.Coverage(bindings, gapdomain.Mastery{Modules: map[string]float64{"percent": 1.0}})
	deficits := make([]Deficit, 0, len(report.Deficits))
	for _, deficit := range report.Deficits {
		requirementID := deficit.RequirementID
		blocking := deficit.BlockingModuleID
		deficits = append(deficits, Deficit{RequirementId: requirementID, BlockingModuleId: &blocking})
	}
	writeJSONBody(w, http.StatusOK, CoverageResponse{
		Covered:  report.Covered,
		Total:    report.Total,
		Percent:  float32(report.Percent),
		Deficits: &deficits,
	})
}

// DiagnoseGaps returns learner gap diagnosis computed over the in-memory fixture.
func (h *StubHandler) DiagnoseGaps(w http.ResponseWriter, _ *http.Request, learnerID string, params DiagnoseGapsParams) {
	h.Logger.Info("gap diagnosis query", zap.String("learner_id", learnerID), zap.String("lag_module_id", params.LagModuleId))
	mastery := gapdomain.Mastery{Modules: map[string]float64{"percent": 0.7, "solutions": 0.2}}
	result := h.Gap.Diagnose(h.gapGraph, mastery, params.LagModuleId)
	rootCauses := make([]RootCause, 0, len(result.RootCauses))
	for _, cause := range result.RootCauses {
		masteryValue := float32(cause.Mastery)
		blocked := cause.BlockedModules
		rootCauses = append(rootCauses, RootCause{ModuleId: cause.ModuleID, Mastery: &masteryValue, BlockedModules: &blocked})
	}
	status := NoRootCauseFound
	if result.Status == gapdomain.DiagnosisFound {
		status = RootCausesFound
	}
	writeJSONBody(w, http.StatusOK, GapDiagnosisResponse{Status: status, RootCauses: rootCauses})
}

// ListModuleResources returns resources bound to a module from the in-memory catalog.
func (h *StubHandler) ListModuleResources(w http.ResponseWriter, _ *http.Request, moduleID string) {
	h.Logger.Info("module resources query", zap.String("module_id", moduleID))
	result := h.Resources.Query(resourceapp.CatalogQuery{ModuleID: moduleID})
	items := make([]Resource, 0, len(result.Items))
	for _, resource := range result.Items {
		items = append(items, mapResource(resource))
	}
	writeJSONBody(w, http.StatusOK, ResourceListResponse{Items: items, Total: len(items)})
}

// ListResources returns the resource catalog filtered by query params.
func (h *StubHandler) ListResources(w http.ResponseWriter, _ *http.Request, params ListResourcesParams) {
	h.Logger.Info("resource catalog query", zap.String("format", stringOr(params.Format)), zap.String("type", stringOrType(params.Type)))
	filter := resourcedomain.ResourceFilter{}
	if params.Type != nil {
		filter.Type = resourcedomain.ResourceType(*params.Type)
	}
	if params.Format != nil {
		filter.Format = *params.Format
	}
	if params.Difficulty != nil {
		filter.Difficulty = *params.Difficulty
	}
	if params.Limit != nil {
		filter.Limit = *params.Limit
	}
	if params.Offset != nil {
		filter.Offset = *params.Offset
	}
	result := h.Resources.Query(resourceapp.CatalogQuery{Filters: filter})
	items := make([]Resource, 0, len(result.Items))
	for _, resource := range result.Items {
		items = append(items, mapResource(resource))
	}
	writeJSONBody(w, http.StatusOK, ResourceListResponse{Items: items, Total: result.Total})
}

// mapResource converts a domain resource to the API response type.
func mapResource(resource resourcedomain.Resource) Resource {
	format := resource.Format
	source := resource.Source
	difficulty := resource.Difficulty
	duration := resource.DurationMinutes
	cost := float32(resource.Cost)
	uri := resource.URI
	return Resource{
		Id:              resource.ID,
		Title:           &resource.Title,
		Type:            ResourceType(resource.Type),
		Format:          &format,
		Source:          &source,
		Difficulty:      &difficulty,
		DurationMinutes: &duration,
		Cost:            &cost,
		Uri:             &uri,
	}
}

// GetPlan returns a fixed learning plan (stub — real application wiring lands in T18).
func (h *StubHandler) GetPlan(w http.ResponseWriter, _ *http.Request, planID string) {
	h.Logger.Info("plan get query", zap.String("plan_id", planID))
	notImplemented(w, "GET /plans/{plan_id}")
}

// SparqlQuery executes a read-only SPARQL query guarded by the F6 gate.
func (h *StubHandler) SparqlQuery(w http.ResponseWriter, _ *http.Request, params SparqlQueryParams) {
	h.Logger.Info("sparql query", zap.String("query", params.Query))
	if err := sparql.ValidateReadOnly(params.Query); err != nil {
		msg := err.Error()
		writeJSONBody(w, http.StatusBadRequest, ErrorResponse{Error: "invalid_sparql", Message: &msg})
		return
	}
	// Query execution against the ontology lands with the read-only store
	// adapter (T20+). Return an empty result set so the endpoint contract is
	// stable while execution is being wired.
	vars := []string{}
	bindings := []map[string]interface{}{}
	writeJSONBody(w, http.StatusOK, SparqlResponse{
		Head: struct {
			Vars *[]string `json:"vars,omitempty"`
		}{Vars: &vars},
		Results: struct {
			Bindings *[]map[string]interface{} `json:"bindings,omitempty"`
		}{Bindings: &bindings},
	})
}

// stringOr dereferences a string pointer ("" for nil).
func stringOr(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

// stringOrType dereferences a resource type pointer ("" for nil).
func stringOrType(t *ListResourcesParamsType) string {
	if t == nil {
		return ""
	}
	return string(*t)
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
