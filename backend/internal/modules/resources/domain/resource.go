// Package resources provides the domain model for the resources bounded context.
package resources

import "fmt"

// ResourceType distinguishes content from enabling resources.
type ResourceType string

const (
	ResourceTypeContent  ResourceType = "content"
	ResourceTypeEnabling ResourceType = "enabling"
)

// LinkType is the ontology relation used for resource binding.
type LinkType string

const (
	LinkAppliesTo LinkType = "appliesTo"
	LinkEnriches  LinkType = "enriches"
)

// Resource is a catalog aggregate entry.
type Resource struct {
	ID              string
	Title           string
	Type            ResourceType
	Format          string
	Source          string
	Style           string
	Difficulty      string
	DurationMinutes int
	Cost            float64
	URI             string
}

// Validate enforces catalog invariants.
func (r Resource) Validate() error {
	if r.ID == "" {
		return fmt.Errorf("resource id is required")
	}
	if r.Type != ResourceTypeContent && r.Type != ResourceTypeEnabling {
		return fmt.Errorf("invalid resource type %q", r.Type)
	}
	if r.Cost < 0 {
		return fmt.Errorf("resource cost must be non-negative")
	}
	return nil
}

// ResourceBinding associates a resource with a route module.
type ResourceBinding struct {
	ResourceID string
	ModuleID   string
	LinkType   LinkType
}

// Validate checks binding invariants.
func (b ResourceBinding) Validate() error {
	if b.ResourceID == "" || b.ModuleID == "" {
		return fmt.Errorf("resource and module ids are required")
	}
	if b.LinkType != LinkAppliesTo && b.LinkType != LinkEnriches {
		return fmt.Errorf("invalid resource link type %q", b.LinkType)
	}
	return nil
}
