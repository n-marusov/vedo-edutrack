package hub

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"go.uber.org/zap"
)

// BenchmarkGraphNeighborhood10kModules measures subgraph copy throughput on a
// 10k-module neighborhood (NFR-FR-api.hub.copy-subgraph AC-3: <=5s).
func BenchmarkGraphNeighborhood10kModules(b *testing.B) {
	modules := make([]map[string]any, 0, 10000)
	links := make([]map[string]any, 0, 10000)
	for i := 0; i < 10000; i++ {
		modules = append(modules, map[string]any{
			"id": "m" + itoa(i), "title": "Module", "description": "d",
			"subject": "math", "grade": "5", "version": "v1",
			"metadata": map[string]any{"source": "mock"}, "fgosBindings": []map[string]any{},
			"resources": []map[string]any{}, "stories": []map[string]any{},
		})
		if i > 0 {
			links = append(links, map[string]any{"sourceModuleId": "m" + itoa(i-1), "targetModuleId": "m" + itoa(i), "type": "hasStrictPrerequisite", "metadata": map[string]any{}})
		}
	}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]any{
			"data": map[string]any{
				"graphNeighborhood": map[string]any{"ontologyId": "mvp", "version": "v1", "modules": modules, "links": links},
			},
		})
	}))
	defer server.Close()

	client, err := NewClient(Config{BaseURL: server.URL, Timeout: 3 * time.Second}, zap.NewNop())
	if err != nil {
		b.Fatalf("init client: %v", err)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, err := client.GraphNeighborhood(b.Context(), "mvp", "m0", 1); err != nil {
			b.Fatalf("graph neighborhood: %v", err)
		}
	}
}

func itoa(i int) string {
	if i == 0 {
		return "0"
	}
	buf := [20]byte{}
	pos := len(buf)
	for i > 0 {
		pos--
		buf[pos] = byte('0' + i%10)
		i /= 10
	}
	return string(buf[pos:])
}
