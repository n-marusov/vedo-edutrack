package mockhub

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// testOntology returns a tiny ontology for handler tests.
func testOntology(t *testing.T) *Ontology {
	t.Helper()
	src := `
@prefix owl: <http://www.w3.org/2002/07/owl#> .
@prefix rdf: <http://www.w3.org/1999/02/22-rdf-syntax-ns#> .
@prefix rdfs: <http://www.w3.org/2000/01/rdf-schema#> .
@prefix tr: <https://vedo-edutrack.dev/traceability#> .

tr:Artifact a owl:Class ;
  rdfs:label "Artifact"@en ;
  rdfs:comment "Top-level class."@en .

tr:Vision a owl:Class ;
  rdfs:label "Vision"@en ;
  rdfs:subClassOf tr:Artifact .
`
	ont, err := Parse(strings.NewReader(src))
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}
	return ont
}

func graphqlRequest(t *testing.T, h http.Handler, query string) (int, map[string]interface{}) {
	t.Helper()
	body, _ := json.Marshal(map[string]string{"query": query})
	req := httptest.NewRequest(http.MethodPost, "/graphql", bytes.NewReader(body))
	req.Header.Set("Authorization", "Bearer test")
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)

	var resp map[string]interface{}
	raw, _ := io.ReadAll(rec.Body)
	_ = json.Unmarshal(raw, &resp)
	return rec.Code, resp
}

func TestGraphQLClasses(t *testing.T) {
	h := NewHandler(testOntology(t), nil)

	code, resp := graphqlRequest(t, h, `{ classes(ontologyId: "t", perPage: 5) { total items { id label } } }`)
	if code != http.StatusOK {
		t.Fatalf("status = %d, want 200; body=%v", code, resp)
	}
	data, _ := resp["data"].(map[string]interface{})
	classes, _ := data["classes"].(map[string]interface{})
	if classes == nil {
		t.Fatalf("classes missing; resp=%v", resp)
	}
	if classes["total"] != float64(2) {
		t.Errorf("total = %v, want 2", classes["total"])
	}
}

func TestGraphQLClassDetail(t *testing.T) {
	h := NewHandler(testOntology(t), nil)

	code, resp := graphqlRequest(t, h,
		`{ class(ontologyId: "t", classId: "https://vedo-edutrack.dev/traceability#Artifact") { id label parents } }`)
	if code != http.StatusOK {
		t.Fatalf("status = %d; body=%v", code, resp)
	}
	data, _ := resp["data"].(map[string]interface{})
	class, _ := data["class"].(map[string]interface{})
	if class == nil {
		t.Fatalf("class missing; resp=%v", resp)
	}
	if class["label"] != "Artifact" {
		t.Errorf("label = %v, want Artifact", class["label"])
	}
}

func TestGraphQLUnauthorized(t *testing.T) {
	h := NewHandler(testOntology(t), nil)

	body := `{"query":"{ classes(ontologyId:\"t\") { total } }"}`
	req := httptest.NewRequest(http.MethodPost, "/graphql", bytes.NewReader([]byte(body)))
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("status = %d, want 401", rec.Code)
	}
}

func TestGraphQLServiceSDL(t *testing.T) {
	h := NewHandler(testOntology(t), nil)

	code, resp := graphqlRequest(t, h, `{ _service { sdl } }`)
	if code != http.StatusOK {
		t.Fatalf("status = %d; body=%v", code, resp)
	}
	data, _ := resp["data"].(map[string]interface{})
	svc, _ := data["_service"].(map[string]interface{})
	if svc == nil {
		t.Fatalf("_service missing; resp=%v", resp)
	}
	if sdl, _ := svc["sdl"].(string); !strings.Contains(sdl, "type QueryRoot") {
		t.Errorf("sdl missing QueryRoot")
	}
}
