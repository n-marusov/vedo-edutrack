package resources

import (
	"fmt"
	"testing"
)

// BenchmarkCatalogFilter10kResources measures catalog filtering on a 10k-item
// catalog (NFR-FR-resource.catalog.filter-by-format AC-4: <=200ms).
func BenchmarkCatalogFilter10kResources(b *testing.B) {
	items := make([]Resource, 0, 10000)
	for i := 0; i < 10000; i++ {
		format := "video"
		if i%3 == 0 {
			format = "text"
		}
		items = append(items, Resource{ID: fmt.Sprintf("res-%d", i), Title: fmt.Sprintf("Resource %d", i), Type: ResourceTypeContent, Format: format, Difficulty: "basic"})
	}
	catalog, err := NewCatalog(items)
	if err != nil {
		b.Fatalf("new catalog: %v", err)
	}
	filter := ResourceFilter{Format: "video"}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = catalog.Filter(filter)
	}
}
