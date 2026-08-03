// Package api implements the HTTP API layer generated from the OpenAPI spec
// (backend/api/openapi/v1.yaml). Handlers are thin input adapters over the
// Application layer (ADR-DES.API.communication-patterns, OpenAPI-first).
package api

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	openapi_types "github.com/oapi-codegen/runtime/types"
	"go.uber.org/zap"

	execapp "vedo-edutrack/backend/internal/modules/executionprogress/application"
	execdomain "vedo-edutrack/backend/internal/modules/executionprogress/domain"
	gapapp "vedo-edutrack/backend/internal/modules/gapcoverage/application"
	gapdomain "vedo-edutrack/backend/internal/modules/gapcoverage/domain"
	"vedo-edutrack/backend/internal/modules/integrations/adapters/sparql"
	ontostub "vedo-edutrack/backend/internal/modules/ontologyport/adapters/stub"
	practiceapp "vedo-edutrack/backend/internal/modules/practicelife/application"
	practicedomain "vedo-edutrack/backend/internal/modules/practicelife/domain"
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
	Practice  *practiceapp.PracticeLifeService
	Progress  *execapp.ProgressService
	Forecast  *execapp.ForecastService
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

	// M2 practice-life fixtures: stories linked via appliesTo/enriches graph
	// links and cross-subject project ideas (REQ-FR-practice.*).
	practiceStories := []practicedomain.Story{
		{ID: "story-1", Title: "Проценты в жизни", Text: "Где встречаются проценты: скидки, банки, статистика.", LinkedModules: []string{"percent", "math-5-11"}, Subjects: []string{"Математика"}, RealWorld: "Проценты используются в банковских вкладах, скидках и социологии.", ReadingMinutes: 3},
		{ID: "story-2", Title: "Растворы вокруг нас", Text: "Концентрация растворов в медицине и кулинарии.", LinkedModules: []string{"solutions"}, Subjects: []string{"Химия"}, RealWorld: "Физраствор, уксус, морская вода — все это растворы разной концентрации.", ReadingMinutes: 4},
		{ID: "story-3", Title: "Химия и экология", Text: "Химические реакции в природе.", LinkedModules: []string{"chemistry"}, Subjects: []string{"Химия", "Биология"}, RealWorld: "Фотосинтез — это химическая реакция, происходящая в каждом листе.", ReadingMinutes: 5},
	}
	practiceProjects := []practicedomain.ProjectIdea{
		{ID: "proj-1", Title: "Биохимическая лаборатория дома", Modules: []string{"solutions", "chemistry"}, DifficultyLevel: practicedomain.DifficultyMedium, ExpectedOutcome: "Провести серию опытов по концентрации растворов и описать результаты."},
		{ID: "proj-2", Title: "Экология двора", Modules: []string{"chemistry", "percent"}, DifficultyLevel: practicedomain.DifficultyBasic, ExpectedOutcome: "Рассчитать процент загрязнения и предложить меры."},
	}

	return &StubHandler{
		Ontology:  graph,
		Routes:    routestub.NewComputer(graph),
		Resources: resourceapp.NewCatalogService(catalog, logger),
		Gap:       gapapp.NewGapService(logger),
		Coverage:  gapapp.NewCoverageService(logger),
		Practice:  practiceapp.NewPracticeLifeService(logger, practiceStories, practiceProjects),
		Progress:  execapp.NewProgressService(logger),
		Forecast:  execapp.NewForecastService(inMemoryProgressRepo{}, logger),
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

// GetDeficitList returns the prioritized deficit list for a learner.
func (h *StubHandler) GetDeficitList(w http.ResponseWriter, _ *http.Request, learnerID string) {
	h.Logger.Info("deficit list query", zap.String("learner_id", learnerID))
	bindings := []gapdomain.FgosBinding{
		{ModuleID: "percent", RequirementID: "fgos-math-5"},
		{ModuleID: "solutions", RequirementID: "fgos-chem-8"},
		{ModuleID: "chemistry", RequirementID: "fgos-chem-8-1"},
	}
	report := h.Coverage.Coverage(bindings, gapdomain.Mastery{Modules: map[string]float64{"percent": 1.0}})

	// Prioritize: strict-prerequisite > essential > optional (M2 F2.7 extension).
	prioritized := make([]PrioritizedDeficit, 0, len(report.Deficits))
	for _, deficit := range report.Deficits {
		priority := PrioritizedDeficitPriority("optional")
		if deficit.BlockingModuleID == "percent" {
			priority = "strict_prerequisite"
		} else if deficit.BlockingModuleID == "solutions" {
			priority = "essential"
		}
		status := PrioritizedDeficitStatus("missing")
		linked := []string{}
		if deficit.BlockingModuleID != "" {
			linked = append(linked, deficit.BlockingModuleID)
		}
		prioritized = append(prioritized, PrioritizedDeficit{
			RequirementId:   deficit.RequirementID,
			Priority:        priority,
			Status:          &status,
			LinkedModuleIds: &linked,
		})
	}
	writeJSONBody(w, http.StatusOK, DeficitListResponse{Deficits: prioritized, Total: len(prioritized)})
}

// GetLearnerProgress returns the plan-vs-actual progress report.
func (h *StubHandler) GetLearnerProgress(w http.ResponseWriter, _ *http.Request, learnerID string) {
	h.Logger.Info("learner progress query", zap.String("learner_id", learnerID))

	plan := execdomain.FixedPlan{
		ID:        "plan-1",
		LearnerID: learnerID,
		Modules: []execdomain.PlannedModule{
			{ModuleID: "percent", PlannedStart: time.Now().AddDate(0, -2, 0), PlannedEnd: time.Now().AddDate(0, -1, 0)},
			{ModuleID: "solutions", PlannedStart: time.Now().AddDate(0, -1, 0), PlannedEnd: time.Now().AddDate(0, 0, 0)},
		},
	}
	progress := []execdomain.ModuleProgress{
		{ModuleID: "percent", Status: execdomain.StatusMastered, MasteredAt: ptrTime(time.Now().AddDate(0, 0, -5))},
		{ModuleID: "solutions", Status: execdomain.StatusInProgress},
	}
	result := h.Progress.Compare(plan, progress)

	// Map deviations back to per-module items.
	items := make([]ModuleProgressItem, 0, len(plan.Modules))
	deviationDays := map[string]int{}
	for _, d := range result.Deviations {
		deviationDays[d.ModuleID] = d.DeltaDays
	}
	for _, module := range plan.Modules {
		item := ModuleProgressItem{ModuleId: module.ModuleID, Status: ModuleProgressItemStatus("not_started")}
		for _, p := range progress {
			if p.ModuleID == module.ModuleID {
				item.Status = ModuleProgressItemStatus(p.Status)
				if p.MasteredAt != nil {
					item.ActualDate = ptrDate(*p.MasteredAt)
				}
			}
		}
		item.PlannedDate = ptrDate(module.PlannedEnd)
		if days, ok := deviationDays[module.ModuleID]; ok {
			item.DeviationDays = &days
			if days > 0 {
				cause := ModuleProgressItemDeviationCause("more_practice")
				item.DeviationCause = &cause
			}
		}
		items = append(items, item)
	}
	now := time.Now()
	writeJSONBody(w, http.StatusOK, ProgressResponse{
		LearnerId:   learnerID,
		PlanId:      &plan.ID,
		GeneratedAt: &now,
		Modules:     items,
	})
}

// GetLearnerForecast returns the binary readiness forecast.
func (h *StubHandler) GetLearnerForecast(w http.ResponseWriter, _ *http.Request, learnerID string) {
	h.Logger.Info("learner forecast query", zap.String("learner_id", learnerID))

	result, err := h.Forecast.ForecastReadiness(context.Background(), learnerID, "plan-1", 30, 60)
	if err != nil {
		// Log the full error server-side with context; expose only a stable
		// client-safe code (REQ-NFR-ops.observability.log-level-config,
		// security: no internal details in responses).
		h.Logger.Error("forecast computation failed",
			zap.String("learner_id", learnerID),
			zap.Error(err),
		)
		msg := "forecast computation failed"
		writeJSONBody(w, http.StatusInternalServerError, ErrorResponse{Error: "forecast_failed", Message: &msg})
		return
	}

	status := ForecastResponseStatus("on-track")
	if result.Status == execdomain.ForecastNotOnTrack {
		status = "not-on-track"
	}
	confidence := ForecastResponseDataConfidence(result.DataConfidence)
	velocity := float32(result.Velocity)
	remaining := result.RemainingModules
	writeJSONBody(w, http.StatusOK, ForecastResponse{
		Status:           status,
		Velocity:         &velocity,
		RemainingModules: &remaining,
		KeyRisks:         &result.KeyRisks,
		DataConfidence:   &confidence,
	})
}

// RecordModuleMastered records a module mastery event.
func (h *StubHandler) RecordModuleMastered(w http.ResponseWriter, r *http.Request, learnerID string) {
	var req MasteryRecordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		msg := "invalid request body"
		writeJSONBody(w, http.StatusBadRequest, ErrorResponse{Error: "invalid_request", Message: &msg})
		return
	}

	// Input validation: reject empty module IDs and unsupported statuses
	// (REQ-FR-execute.progress.plan-vs-actual; security: validate all input).
	switch req.Status {
	case MasteryRecordRequestStatus("in_progress"), MasteryRecordRequestStatus("mastered"), MasteryRecordRequestStatus("skipped"):
	default:
		msg := "status must be one of in_progress|mastered|skipped"
		writeJSONBody(w, http.StatusBadRequest, ErrorResponse{Error: "invalid_request", Message: &msg})
		return
	}
	if req.ModuleId == "" {
		msg := "module_id is required"
		writeJSONBody(w, http.StatusBadRequest, ErrorResponse{Error: "invalid_request", Message: &msg})
		return
	}

	h.Logger.Info("module mastered",
		zap.String("learner_id", learnerID),
		zap.String("module_id", req.ModuleId),
		zap.String("status", string(req.Status)),
	)

	now := time.Now()
	writeJSONBody(w, http.StatusCreated, MasteryRecordResponse{
		LearnerId:  learnerID,
		ModuleId:   req.ModuleId,
		Status:     string(req.Status),
		RecordedAt: &now,
	})
}

// GetModuleStories returns stories linked to a module via appliesTo/enriches.
func (h *StubHandler) GetModuleStories(w http.ResponseWriter, _ *http.Request, moduleID string) {
	h.Logger.Info("module stories query", zap.String("module_id", moduleID))
	stories := h.Practice.StoriesForModule(moduleID)
	items := make([]StoryResponse, 0, len(stories))
	for _, s := range stories {
		items = append(items, mapStory(s))
	}
	writeJSONBody(w, http.StatusOK, items)
}

// GetModuleProjects returns project ideas linked to a module.
func (h *StubHandler) GetModuleProjects(w http.ResponseWriter, _ *http.Request, moduleID string) {
	h.Logger.Info("module projects query", zap.String("module_id", moduleID))
	projects := h.Practice.ProjectsForModule(moduleID)
	items := make([]ProjectIdeaResponse, 0, len(projects))
	for _, p := range projects {
		items = append(items, mapProject(p))
	}
	writeJSONBody(w, http.StatusOK, items)
}

// GetRecommendedStories returns story recommendations at mastery.
func (h *StubHandler) GetRecommendedStories(w http.ResponseWriter, _ *http.Request, learnerID string) {
	h.Logger.Info("recommended stories query", zap.String("learner_id", learnerID))
	stories := h.Practice.RecommendStories([]string{"percent", "solutions"})
	items := make([]StoryResponse, 0, len(stories))
	for _, s := range stories {
		items = append(items, mapStory(s))
	}
	writeJSONBody(w, http.StatusOK, items)
}

// GetRecommendedProjects returns project suggestions for a learner.
func (h *StubHandler) GetRecommendedProjects(w http.ResponseWriter, _ *http.Request, learnerID string) {
	h.Logger.Info("recommended projects query", zap.String("learner_id", learnerID))
	readySet := map[string]bool{"percent": true, "solutions": true, "chemistry": true}
	projects := h.Practice.SuggestProjects(readySet)
	items := make([]ProjectIdeaResponse, 0, len(projects))
	for _, p := range projects {
		items = append(items, mapProject(p))
	}
	writeJSONBody(w, http.StatusOK, items)
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

// mapStory converts a practice-life domain story to the API response type.
func mapStory(s practicedomain.Story) StoryResponse {
	return StoryResponse{
		Id:             s.ID,
		Title:          s.Title,
		Text:           &s.Text,
		LinkedModules:  &s.LinkedModules,
		Subjects:       &s.Subjects,
		RealWorld:      &s.RealWorld,
		ReadingMinutes: &s.ReadingMinutes,
	}
}

// mapProject converts a practice-life domain project idea to the API response type.
func mapProject(p practicedomain.ProjectIdea) ProjectIdeaResponse {
	level := ProjectIdeaResponseDifficultyLevel(p.DifficultyLevel)
	return ProjectIdeaResponse{
		Id:              p.ID,
		Title:           p.Title,
		Modules:         &p.Modules,
		DifficultyLevel: &level,
		ExpectedOutcome: &p.ExpectedOutcome,
	}
}

// ptrTime returns a pointer to t.
func ptrTime(t time.Time) *time.Time { return &t }

// ptrDate returns a pointer to an openapi date for t.
func ptrDate(t time.Time) *openapi_types.Date {
	d := openapi_types.Date{Time: t}
	return &d
}

// inMemoryProgressRepo is a minimal ProgressRepository fixture used by the
// stub handler so the forecast endpoint can serve data without a database.
type inMemoryProgressRepo struct{}

func (inMemoryProgressRepo) GetProgress(_ context.Context, _ string) ([]execdomain.ModuleProgress, error) {
	return []execdomain.ModuleProgress{
		{ModuleID: "percent", Status: execdomain.StatusMastered, MasteredAt: ptrTime(time.Now().AddDate(0, 0, -5))},
		{ModuleID: "solutions", Status: execdomain.StatusInProgress},
	}, nil
}

func (inMemoryProgressRepo) GetPlannedModuleCount(_ context.Context, _, _ string) (int, error) {
	return 3, nil
}

func (inMemoryProgressRepo) GetRemainingModules(_ context.Context, _, _ string) (int, error) {
	return 2, nil
}
