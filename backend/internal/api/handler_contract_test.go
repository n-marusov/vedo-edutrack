package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// newContractRouter builds the generated chi router with the M1 handler wiring.
func newContractRouter(t *testing.T) http.Handler {
	t.Helper()
	return HandlerWithOptions(NewStubHandler(), ChiServerOptions{})
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
	if len(resp) != 1 || resp[0].Id != "story-1" {
		t.Fatalf("expected story-1 for module percent, got %+v", resp)
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
