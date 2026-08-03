package routeplanning

// Module is a route-planning view of an ontology module.
type Module struct {
	ID        string
	Title     string
	Essential bool
}

// LinkType describes a route-planning edge imported from the ontology port.
type LinkType string

const (
	LinkStrictPrerequisite LinkType = "hasStrictPrerequisite"
	LinkSoftPrerequisite   LinkType = "hasSoftPrerequisite"
	LinkEnriches           LinkType = "enriches"
	LinkAppliesTo          LinkType = "appliesTo"
)

// Link is a directed graph edge from source prerequisite/context to target goal module.
type Link struct {
	SourceID string
	TargetID string
	Type     LinkType
}

// OntologyGraph is the pure-domain graph snapshot used by pathfinding.
type OntologyGraph struct {
	Modules []Module
	Links   []Link
}

// Route is an ordered learning route from current position to goal.
type Route struct {
	Steps       []RouteStep
	TotalWeight int
	Far         Horizon
	Mid         Horizon
	Near        Horizon
}

// RouteStep is one ordered module in a computed route.
type RouteStep struct {
	ModuleID string
	Order    int
	Via      LinkType
	Weight   int
}

// ModuleIDs returns route steps as module ids in learning order.
func (r Route) ModuleIDs() []string {
	ids := make([]string, 0, len(r.Steps))
	for _, step := range r.Steps {
		ids = append(ids, step.ModuleID)
	}
	return ids
}

// SatisfiesStrictPrerequisites verifies all strict source modules precede their targets when both are in the route.
func (r Route) SatisfiesStrictPrerequisites(graph OntologyGraph) bool {
	position := map[string]int{}
	for i, step := range r.Steps {
		position[step.ModuleID] = i
	}
	for _, link := range graph.Links {
		if link.Type != LinkStrictPrerequisite {
			continue
		}
		sourcePos, sourceOK := position[link.SourceID]
		targetPos, targetOK := position[link.TargetID]
		if sourceOK && targetOK && sourcePos > targetPos {
			return false
		}
	}
	return true
}
