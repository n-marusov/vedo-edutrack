package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"vedo-edutrack/backend/internal/platform/auth"
)

// newContractRouter builds the generated chi router with the M1 handler wiring.
func newContractRouter(t *testing.T) http.Handler {
	t.Helper()
	return HandlerWithOptions(NewStubHandler(nil, nil, nil), ChiServerOptions{})
}

func TestContractListResources(t *testing.T) {
	router := newContractRouter(t)
	req := httptest.NewRequest(http.MethodGet, "/resources?format=video", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status=%d, want 200; body=%s", rec.Code, rec.Body.String())
	}
	var resp ResourceListResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if len(resp.Items) == 0 || resp.Total == 0 {
		t.Fatalf("expected non-empty catalog, got %+v", resp)
	}
}

func TestContractModuleResources(t *testing.T) {
	router := newContractRouter(t)
	req := httptest.NewRequest(http.MethodGet, "/modules/math-5-11/resources", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status=%d, want 200; body=%s", rec.Code, rec.Body.String())
	}
	var resp ResourceListResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if len(resp.Items) != 1 || resp.Items[0].Id != "res-1" {
		t.Fatalf("expected bound resource res-1, got %+v", resp)
	}
}

func TestContractDiagnoseGaps(t *testing.T) {
	router := newContractRouter(t)
	req := httptest.NewRequest(http.MethodGet, "/learners/l1/gaps?lag_module_id=chemistry", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status=%d, want 200; body=%s", rec.Code, rec.Body.String())
	}
	var resp GapDiagnosisResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if resp.Status != RootCausesFound || len(resp.RootCauses) == 0 {
		t.Fatalf("expected root causes, got %+v", resp)
	}
}

func TestContractFgosCoverage(t *testing.T) {
	router := newContractRouter(t)
	req := httptest.NewRequest(http.MethodGet, "/learners/l1/coverage/fgos", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status=%d, want 200; body=%s", rec.Code, rec.Body.String())
	}
	var resp CoverageResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if resp.Total == 0 || resp.Percent < 0 || resp.Percent > 100 {
		t.Fatalf("expected valid coverage, got %+v", resp)
	}
}

func TestContractSparqlReadOnlyGate(t *testing.T) {
	router := newContractRouter(t)
	req := httptest.NewRequest(http.MethodGet, "/sparql?query=SELECT%20%3Fs%20WHERE%20%7B%3Fs%20%3Fp%20%3Fo%7D", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	// T20 makes this 200 with results; until then the endpoint must exist with
	// a well-formed error response (never 404).
	if rec.Code == http.StatusNotFound || rec.Code == http.StatusMethodNotAllowed {
		t.Fatalf("sparql route missing: status=%d", rec.Code)
	}
}

func TestContractRecordModuleMasteredValid(t *testing.T) {
	router := newContractRouter(t)
	body := `{"module_id":"percent","status":"mastered"}`
	req := httptest.NewRequest(http.MethodPost, "/learners/l1/module-mastered", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("status=%d, want 201; body=%s", rec.Code, rec.Body.String())
	}
	var resp MasteryRecordResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if resp.ModuleId != "percent" || resp.Status != "mastered" {
		t.Fatalf("unexpected response: %+v", resp)
	}
}

func TestContractRecordModuleMasteredRejectsEmptyModuleID(t *testing.T) {
	router := newContractRouter(t)
	body := `{"module_id":"","status":"mastered"}`
	req := httptest.NewRequest(http.MethodPost, "/learners/l1/module-mastered", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status=%d, want 400; body=%s", rec.Code, rec.Body.String())
	}
}

func TestContractRecordModuleMasteredRejectsInvalidStatus(t *testing.T) {
	router := newContractRouter(t)
	body := `{"module_id":"percent","status":"not_a_status"}`
	req := httptest.NewRequest(http.MethodPost, "/learners/l1/module-mastered", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status=%d, want 400; body=%s", rec.Code, rec.Body.String())
	}
}

func TestContractLearnerForecast(t *testing.T) {
	router := newContractRouter(t)
	req := httptest.NewRequest(http.MethodGet, "/learners/l1/forecast", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status=%d, want 200; body=%s", rec.Code, rec.Body.String())
	}
	var resp ForecastResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if resp.Status != ForecastResponseStatus("on-track") && resp.Status != ForecastResponseStatus("not-on-track") {
		t.Fatalf("unexpected forecast status: %s", resp.Status)
	}
}

func TestContractModuleStories(t *testing.T) {
	router := newContractRouter(t)
	req := httptest.NewRequest(http.MethodGet, "/modules/percent/stories", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status=%d, want 200; body=%s", rec.Code, rec.Body.String())
	}
	var resp []StoryResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if len(resp) == 0 {
		t.Fatal("expected at least one story for module percent")
	}
	foundMathStory := false
	for _, s := range resp {
		if s.Id == "story-math-percent" {
			foundMathStory = true
			break
		}
	}
	if !foundMathStory {
		t.Fatalf("expected story-math-percent in %+v", resp)
	}
}

func TestContractErrorFormat(t *testing.T) {
	router := newContractRouter(t)
	req := httptest.NewRequest(http.MethodGet, "/plans/missing-plan", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	var resp ErrorResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("error body must be valid ErrorResponse JSON, got %q", rec.Body.String())
	}
	if resp.Error == "" {
		t.Fatalf("expected machine-readable error code, got %+v", resp)
	}
}

func TestContractWebhookSubscriptionLifecycle(t *testing.T) {
	router := newContractRouter(t)
	tenant := "user-42"

	// Create.
	createBody := `{"url":"https://example.test/hook","event_types":["module.mastered"],"secret":"01234567890123456789012345678901"}`
	req := httptest.NewRequest(http.MethodPost, "/webhooks/subscriptions", strings.NewReader(createBody))
	req.Header.Set("Content-Type", "application/json")
	req = withUser(req, tenant)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	if rec.Code != http.StatusCreated {
		t.Fatalf("create: status=%d, want 201; body=%s", rec.Code, rec.Body.String())
	}
	var sub WebhookSubscription
	if err := json.Unmarshal(rec.Body.Bytes(), &sub); err != nil {
		t.Fatalf("decode created subscription: %v", err)
	}
	if sub.Id.String() == "" || !sub.Active || len(sub.EventTypes) != 1 {
		t.Fatalf("unexpected subscription: %+v", sub)
	}

	// List.
	req = withUser(httptest.NewRequest(http.MethodGet, "/webhooks/subscriptions", nil), tenant)
	rec = httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("list: status=%d, want 200; body=%s", rec.Code, rec.Body.String())
	}
	var list []WebhookSubscription
	if err := json.Unmarshal(rec.Body.Bytes(), &list); err != nil {
		t.Fatalf("decode list: %v", err)
	}
	if len(list) != 1 || list[0].Id.String() != sub.Id.String() {
		t.Fatalf("expected 1 subscription, got %+v", list)
	}

	// Get.
	req = withUser(httptest.NewRequest(http.MethodGet, "/webhooks/subscriptions/"+sub.Id.String(), nil), tenant)
	rec = httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("get: status=%d, want 200; body=%s", rec.Code, rec.Body.String())
	}

	// Update.
	updateBody := `{"url":"https://example.test/v2","event_types":["plan.deviated"]}`
	req = withUser(httptest.NewRequest(http.MethodPut, "/webhooks/subscriptions/"+sub.Id.String(), strings.NewReader(updateBody)), tenant)
	req.Header.Set("Content-Type", "application/json")
	rec = httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("update: status=%d, want 200; body=%s", rec.Code, rec.Body.String())
	}
	var updated WebhookSubscription
	if err := json.Unmarshal(rec.Body.Bytes(), &updated); err != nil {
		t.Fatalf("decode updated: %v", err)
	}
	if updated.Url != "https://example.test/v2" || len(updated.EventTypes) != 1 || updated.EventTypes[0] != WebhookSubscriptionEventTypesPlanDeviated {
		t.Fatalf("update not applied: %+v", updated)
	}

	// Delivery history (empty).
	req = withUser(httptest.NewRequest(http.MethodGet, "/webhooks/subscriptions/"+sub.Id.String()+"/deliveries", nil), tenant)
	rec = httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("deliveries: status=%d, want 200; body=%s", rec.Code, rec.Body.String())
	}

	// Ping.
	req = withUser(httptest.NewRequest(http.MethodPost, "/webhooks/subscriptions/"+sub.Id.String()+"/ping", nil), tenant)
	rec = httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	if rec.Code != http.StatusAccepted {
		t.Fatalf("ping: status=%d, want 202; body=%s", rec.Code, rec.Body.String())
	}

	// Delete.
	req = withUser(httptest.NewRequest(http.MethodDelete, "/webhooks/subscriptions/"+sub.Id.String(), nil), tenant)
	rec = httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	if rec.Code != http.StatusNoContent {
		t.Fatalf("delete: status=%d, want 204; body=%s", rec.Code, rec.Body.String())
	}

	// Get after delete → 404.
	req = withUser(httptest.NewRequest(http.MethodGet, "/webhooks/subscriptions/"+sub.Id.String(), nil), tenant)
	rec = httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	if rec.Code != http.StatusNotFound {
		t.Fatalf("get after delete: status=%d, want 404; body=%s", rec.Code, rec.Body.String())
	}
}

// withUser injects auth claims (user id) into the request context, mimicking
// the auth middleware for contract tests that exercise tenant-scoped routes.
// The platform-integrator role is granted so webhook CRUD tests pass RBAC.
func withUser(req *http.Request, userID string) *http.Request {
	return req.WithContext(auth.WithClaims(req.Context(), &auth.Claims{
		UserID: userID,
		Roles:  []string{"platform-integrator"},
	}))
}

func TestContractWebhookCreateRejectsInvalidPayload(t *testing.T) {
	router := newContractRouter(t)
	// http (non-localhost) URL rejected.
	body := `{"url":"http://insecure.example.test/hook","event_types":["module.mastered"]}`
	req := withUser(httptest.NewRequest(http.MethodPost, "/webhooks/subscriptions", strings.NewReader(body)), "user-42")
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("invalid url: status=%d, want 400; body=%s", rec.Code, rec.Body.String())
	}
}

func TestContractWebhookGetMissingReturns404(t *testing.T) {
	router := newContractRouter(t)
	req := httptest.NewRequest(http.MethodGet, "/webhooks/subscriptions/00000000-0000-0000-0000-000000000001", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	if rec.Code != http.StatusNotFound {
		t.Fatalf("status=%d, want 404; body=%s", rec.Code, rec.Body.String())
	}
}
