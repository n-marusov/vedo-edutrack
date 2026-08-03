// Package ontologyport provides the domain model for the ontologyport bounded context.
package ontologyport

import "fmt"

// LinkType is a supported EduTrack ontology relation type read from VEDO Hub.
type LinkType string

const (
	// LinkStrictPrerequisite means the source module must be mastered before the target module.
	LinkStrictPrerequisite LinkType = "hasStrictPrerequisite"
	// LinkSoftPrerequisite means the source module is recommended before the target module.
	LinkSoftPrerequisite LinkType = "hasSoftPrerequisite"
	// LinkEnriches means the source module enriches the target module with optional knowledge.
	LinkEnriches LinkType = "enriches"
	// LinkAppliesTo connects knowledge modules to applications, stories, projects, or resources.
	LinkAppliesTo LinkType = "appliesTo"
	// LinkIsAnalogousTo connects modules that can be used as analogies.
	LinkIsAnalogousTo LinkType = "isAnalogousTo"
)

// Validate checks whether the link type is part of the F0 ontology-port contract.
func (t LinkType) Validate() error {
	switch t {
	case LinkStrictPrerequisite, LinkSoftPrerequisite, LinkEnriches, LinkAppliesTo, LinkIsAnalogousTo:
		return nil
	default:
		return fmt.Errorf("unsupported ontology link type %q", t)
	}
}

// OntologyModule is a knowledge graph node imported from VEDO Hub.
type OntologyModule struct {
	ID           string
	Title        string
	Description  string
	Subject      string
	Grade        string
	Version      string
	Metadata     map[string]string
	FgosBindings []FgosBinding
	Resources    []ResourceRef
	Stories      []StoryRef
}

// OntologyLink is a typed directed edge between ontology modules.
type OntologyLink struct {
	SourceModuleID string
	TargetModuleID string
	Type           LinkType
	Metadata       map[string]string
}

// FgosBinding maps a module to a FGOS/professional-standard requirement.
type FgosBinding struct {
	RequirementID string
	Title         string
	Level         string
	Coverage      float64
}

// PedagogyConcept describes a routing/learning concept resolved through the Hub class hierarchy.
type PedagogyConcept struct {
	ID          string
	Title       string
	Description string
	ParentID    string
	Metadata    map[string]string
}

// ResourceRef points to a resource attached to a module through ontology metadata or appliesTo links.
type ResourceRef struct {
	ID         string
	Title      string
	Kind       string
	Format     string
	Difficulty string
	URI        string
	Metadata   map[string]string
}

// StoryRef points to a practical story/project idea linked to a module.
type StoryRef struct {
	ID       string
	Title    string
	URI      string
	Metadata map[string]string
}

// Subgraph is a copied Hub ontology fragment used by route/gap/resource computations.
type Subgraph struct {
	OntologyID string
	Version    string
	Modules    []OntologyModule
	Links      []OntologyLink
}
