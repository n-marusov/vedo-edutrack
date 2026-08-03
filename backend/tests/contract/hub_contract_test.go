package contract

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"go.uber.org/zap"

	"vedo-edutrack/backend/internal/modules/ontologyport/adapters/hub"
)

func TestHubGraphQLContractGraphNeighborhoodSchema(t *testing.T) {
	t.Parallel()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertGraphQLRequest(t, r, "graphNeighborhood", "test-token")
		writeGraphQLJSON(t, w, map[string]any{
			"data": map[string]any{
				"graphNeighborhood": map[string]any{
					"ontologyId": "mvp",
					"version":    "v1",
					"modules": []map[string]any{{
						"id":           "math-5-11",
						"title":        "Проценты",
						"description":  "Percentages.",
						"subject":      "math",
						"grade":        "5",
						"version":      "m1",
						"metadata":     map[string]any{"source": "hub-mock"},
						"fgosBindings": []map[string]any{{"requirementId": "fgos-math-5", "title": "Math 5", "level": "basic", "coverage": 0.75}},
						"resources":    []map[string]any{{"id": "res-1", "title": "Video", "kind": "content", "format": "video", "difficulty": "basic", "uri": "https://example.test/res-1", "metadata": map[string]any{"duration": "10m"}}},
						"stories":      []map[string]any{{"id": "story-1", "title": "Solution concentration", "uri": "https://example.test/story-1", "metadata": map[string]any{"domain": "chemistry"}}},
					}},
					"links": []map[string]any{{
						"sourceModuleId": "math-5-11",
						"targetModuleId": "chem-8-1",
						"type":           "appliesTo",
						"metadata":       map[string]any{"weight": 1},
					}},
				},
			},
		})
	}))
	defer server.Close()

	client := newContractClient(t, server.URL, "test-token", 3*time.Second)
	subgraph, err := client.GraphNeighborhood(context.Background(), "mvp", "math-5-11", 1)
	if err != nil {
		t.Fatalf("INFO [contract.hub] test: graph neighborhood schema — FAIL: %v", err)
	}
	if subgraph.OntologyID != "mvp" || subgraph.Version != "v1" || len(subgraph.Modules) != 1 || len(subgraph.Links) != 1 {
		t.Fatalf("ERROR [contract.hub] contract violation: unexpected subgraph: %+v", subgraph)
	}
	module := subgraph.Modules[0]
	if module.ID == "" || module.Title == "" || len(module.FgosBindings) != 1 || len(module.Resources) != 1 || len(module.Stories) != 1 {
		t.Fatalf("ERROR [contract.hub] contract violation: expected module metadata/resource/story fields, got %+v", module)
	}
	t.Log("INFO [contract.hub] test: graph neighborhood schema — PASS")
}

func TestHubGraphQLContractClassAndPropertyOperations(t *testing.T) {
	t.Parallel()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req graphQLRequest
		decodeGraphQLRequest(t, r, &req)
		switch req.OperationName {
		case "classDescendants":
			writeGraphQLJSON(t, w, map[string]any{"data": map[string]any{"classDescendants": []map[string]any{{"id": "pedagogy:essential", "title": "Essentialism", "description": "Core first", "parentId": "pedagogy", "metadata": map[string]any{"priority": "core"}}}}})
		case "classTree":
			writeGraphQLJSON(t, w, map[string]any{"data": map[string]any{"classTree": []map[string]any{{"id": "fgos:math", "title": "FGOS Math", "description": "", "parentId": "fgos", "metadata": map[string]any{}}}}})
		case "properties":
			writeGraphQLJSON(t, w, map[string]any{"data": map[string]any{"properties": []map[string]any{{"key": "format", "value": "video"}}}})
		case "property":
			writeGraphQLJSON(t, w, map[string]any{"data": map[string]any{"property": map[string]any{"key": "title", "value": "Percentages"}}})
		default:
			t.Fatalf("ERROR [contract.hub] contract violation: unexpected operation %q", req.OperationName)
		}
	}))
	defer server.Close()

	client := newContractClient(t, server.URL, "", 3*time.Second)
	classes, err := client.ClassDescendants(context.Background(), "mvp", "pedagogy")
	if err != nil || len(classes) != 1 || classes[0].ID == "" {
		t.Fatalf("INFO [contract.hub] test: classDescendants — FAIL: classes=%+v err=%v", classes, err)
	}
	tree, err := client.ClassTree(context.Background(), "mvp", "fgos")
	if err != nil || len(tree) != 1 || tree[0].Title == "" {
		t.Fatalf("INFO [contract.hub] test: classTree — FAIL: tree=%+v err=%v", tree, err)
	}
	properties, err := client.Properties(context.Background(), "mvp", "math-5-11")
	if err != nil || properties["format"] != "video" {
		t.Fatalf("INFO [contract.hub] test: properties — FAIL: properties=%+v err=%v", properties, err)
	}
	value, err := client.Property(context.Background(), "mvp", "math-5-11", "title")
	if err != nil || value != "Percentages" {
		t.Fatalf("INFO [contract.hub] test: property — FAIL: value=%q err=%v", value, err)
	}
	t.Log("INFO [contract.hub] test: class/property operations — PASS")
}

func TestHubGraphQLContractErrorsAndAuthFailure(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name       string
		statusCode int
		body       map[string]any
	}{
		{name: "malformed query", statusCode: http.StatusOK, body: map[string]any{"errors": []map[string]any{{"message": "syntax error"}}}},
		{name: "auth failure", statusCode: http.StatusUnauthorized, body: map[string]any{"errors": []map[string]any{{"message": "unauthorized"}}}},
		{name: "ontology not found", statusCode: http.StatusOK, body: map[string]any{"errors": []map[string]any{{"message": "ontology not found"}}}},
	}
	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				w.WriteHeader(tc.statusCode)
				writeGraphQLJSON(t, w, tc.body)
			}))
			defer server.Close()

			client := newContractClient(t, server.URL, "", 3*time.Second)
			_, err := client.GraphNeighborhood(context.Background(), "mvp", "missing", 1)
			if err == nil {
				t.Fatalf("ERROR [contract.hub] contract violation: expected error for %s", tc.name)
			}
			t.Logf("INFO [contract.hub] test: %s — PASS", tc.name)
		})
	}
}

func TestHubGraphQLContractTimeout(t *testing.T) {
	t.Parallel()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		time.Sleep(50 * time.Millisecond)
		writeGraphQLJSON(t, w, map[string]any{"data": map[string]any{}})
	}))
	defer server.Close()

	client := newContractClient(t, server.URL, "", 10*time.Millisecond)
	_, err := client.GraphNeighborhood(context.Background(), "mvp", "math-5-11", 1)
	if err == nil {
		t.Fatal("ERROR [contract.hub] contract violation: expected timeout error")
	}
	t.Log("INFO [contract.hub] test: timeout behavior — PASS")
}

func newContractClient(t *testing.T, url, token string, timeout time.Duration) *hub.Client {
	t.Helper()
	client, err := hub.NewClient(hub.Config{BaseURL: url, GraphQLPath: "/graphql", BearerToken: token, Timeout: timeout}, zap.NewNop())
	if err != nil {
		t.Fatalf("init client: %v", err)
	}
	return client
}

type graphQLRequest struct {
	OperationName string         `json:"operationName"`
	Variables     map[string]any `json:"variables"`
}

func assertGraphQLRequest(t *testing.T, r *http.Request, operation, token string) {
	t.Helper()
	if r.Method != http.MethodPost {
		t.Fatalf("ERROR [contract.hub] contract violation: method=%s, want POST", r.Method)
	}
	if r.URL.Path != "/graphql" {
		t.Fatalf("ERROR [contract.hub] contract violation: path=%s, want /graphql", r.URL.Path)
	}
	if token != "" && r.Header.Get("Authorization") != "Bearer "+token {
		t.Fatalf("ERROR [contract.hub] contract violation: missing bearer token auth")
	}
	var req graphQLRequest
	decodeGraphQLRequest(t, r, &req)
	if req.OperationName != operation {
		t.Fatalf("ERROR [contract.hub] contract violation: operation=%s, want %s", req.OperationName, operation)
	}
}

func decodeGraphQLRequest(t *testing.T, r *http.Request, req *graphQLRequest) {
	t.Helper()
	if !strings.Contains(r.Header.Get("Content-Type"), "application/json") {
		t.Fatalf("ERROR [contract.hub] contract violation: Content-Type=%q", r.Header.Get("Content-Type"))
	}
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		t.Fatalf("decode request: %v", err)
	}
}

func writeGraphQLJSON(t *testing.T, w http.ResponseWriter, body map[string]any) {
	t.Helper()
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(body); err != nil {
		t.Fatalf("encode response: %v", err)
	}
}
