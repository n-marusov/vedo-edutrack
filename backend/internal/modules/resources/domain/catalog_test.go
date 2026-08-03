package resources

import "testing"

func TestResourceCatalogFiltersPrecisionAndEmptyResult(t *testing.T) {
	catalog, err := NewCatalog([]Resource{
		{ID: "r1", Type: ResourceTypeContent, Format: "video", Source: "school", Difficulty: "basic", DurationMinutes: 10},
		{ID: "r2", Type: ResourceTypeContent, Format: "text", Source: "school", Difficulty: "advanced", DurationMinutes: 20},
		{ID: "r3", Type: ResourceTypeEnabling, Format: "lab", Source: "partner", Difficulty: "basic", DurationMinutes: 60},
	})
	if err != nil {
		t.Fatalf("new catalog: %v", err)
	}
	got := catalog.Filter(ResourceFilter{Format: "video"})
	if len(got.Items) != 1 || got.Items[0].ID != "r1" {
		t.Fatalf("filter precision failed: %+v", got.Items)
	}
	empty := catalog.Filter(ResourceFilter{Format: "interactive"})
	if len(empty.Items) != 0 || empty.Total != 0 {
		t.Fatalf("expected empty result, got %+v", empty)
	}
}

func TestResourceValidation(t *testing.T) {
	if err := (Resource{ID: "bad", Type: "unknown", Cost: 0}).Validate(); err == nil {
		t.Fatal("expected invalid type error")
	}
	if err := (Resource{ID: "bad", Type: ResourceTypeContent, Cost: -1}).Validate(); err == nil {
		t.Fatal("expected negative cost error")
	}
}

func TestModuleBindingPreventsOrphans(t *testing.T) {
	catalog, err := NewCatalog([]Resource{{ID: "r1", Type: ResourceTypeContent}})
	if err != nil {
		t.Fatalf("new catalog: %v", err)
	}
	if err := catalog.BindToModule(ResourceBinding{ResourceID: "missing", ModuleID: "m1", LinkType: LinkAppliesTo}); err == nil {
		t.Fatal("expected orphan resource binding error")
	}
	if err := catalog.BindToModule(ResourceBinding{ResourceID: "r1", ModuleID: "m1", LinkType: LinkAppliesTo}); err != nil {
		t.Fatalf("bind resource: %v", err)
	}
	if len(catalog.ResourcesForModule("m1")) != 1 {
		t.Fatalf("expected one bound resource")
	}
}
