package hub

import (
	"context"
	"fmt"

	"go.uber.org/zap"

	ontology "vedo-edutrack/backend/internal/modules/ontologyport/domain"
)

const graphNeighborhoodQuery = `query graphNeighborhood($ontologyId: ID!, $conceptId: ID!, $depth: Int!) {
  graphNeighborhood(ontologyId: $ontologyId, conceptId: $conceptId, depth: $depth) {
    ontologyId
    version
    modules {
      id
      title
      description
      subject
      grade
      version
      metadata
      fgosBindings { requirementId title level coverage }
      resources { id title kind format difficulty uri metadata }
      stories { id title uri metadata }
    }
    links { sourceModuleId targetModuleId type metadata }
  }
}`

const classDescendantsQuery = `query classDescendants($ontologyId: ID!, $classId: ID!) {
  classDescendants(ontologyId: $ontologyId, classId: $classId) {
    id
    title
    description
    parentId
    metadata
  }
}`

const classTreeQuery = `query classTree($ontologyId: ID!, $rootClassId: ID!) {
  classTree(ontologyId: $ontologyId, rootClassId: $rootClassId) {
    id
    title
    description
    parentId
    metadata
  }
}`

const propertiesQuery = `query properties($ontologyId: ID!, $entityId: ID!) {
  properties(ontologyId: $ontologyId, entityId: $entityId) {
    key
    value
  }
}`

const propertyQuery = `query property($ontologyId: ID!, $entityId: ID!, $key: String!) {
  property(ontologyId: $ontologyId, entityId: $entityId, key: $key) {
    key
    value
  }
}`

type graphNeighborhoodData struct {
	GraphNeighborhood graphNeighborhoodPayload `json:"graphNeighborhood"`
}

type graphNeighborhoodResult struct {
	Data   graphNeighborhoodData `json:"data"`
	Errors []graphQLError        `json:"errors,omitempty"`
}

type graphNeighborhoodPayload struct {
	OntologyID string      `json:"ontologyId"`
	Version    string      `json:"version"`
	Modules    []moduleDTO `json:"modules"`
	Links      []linkDTO   `json:"links"`
}

type moduleDTO struct {
	ID           string           `json:"id"`
	Title        string           `json:"title"`
	Description  string           `json:"description"`
	Subject      string           `json:"subject"`
	Grade        string           `json:"grade"`
	Version      string           `json:"version"`
	Metadata     map[string]any   `json:"metadata"`
	FgosBindings []fgosBindingDTO `json:"fgosBindings"`
	Resources    []resourceRefDTO `json:"resources"`
	Stories      []storyRefDTO    `json:"stories"`
}

type fgosBindingDTO struct {
	RequirementID string  `json:"requirementId"`
	Title         string  `json:"title"`
	Level         string  `json:"level"`
	Coverage      float64 `json:"coverage"`
}

type resourceRefDTO struct {
	ID         string         `json:"id"`
	Title      string         `json:"title"`
	Kind       string         `json:"kind"`
	Format     string         `json:"format"`
	Difficulty string         `json:"difficulty"`
	URI        string         `json:"uri"`
	Metadata   map[string]any `json:"metadata"`
}

type storyRefDTO struct {
	ID       string         `json:"id"`
	Title    string         `json:"title"`
	URI      string         `json:"uri"`
	Metadata map[string]any `json:"metadata"`
}

type linkDTO struct {
	SourceModuleID string         `json:"sourceModuleId"`
	TargetModuleID string         `json:"targetModuleId"`
	Type           string         `json:"type"`
	Metadata       map[string]any `json:"metadata"`
}

type classResult struct {
	Data struct {
		ClassDescendants []classDTO `json:"classDescendants,omitempty"`
		ClassTree        []classDTO `json:"classTree,omitempty"`
	} `json:"data"`
	Errors []graphQLError `json:"errors,omitempty"`
}

type classDTO struct {
	ID          string         `json:"id"`
	Title       string         `json:"title"`
	Description string         `json:"description"`
	ParentID    string         `json:"parentId"`
	Metadata    map[string]any `json:"metadata"`
}

type propertiesResult struct {
	Data struct {
		Properties []propertyDTO `json:"properties,omitempty"`
		Property   *propertyDTO  `json:"property,omitempty"`
	} `json:"data"`
	Errors []graphQLError `json:"errors,omitempty"`
}

type propertyDTO struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// GraphNeighborhood reads modules and typed links around a concept from VEDO Hub.
func (c *Client) GraphNeighborhood(ctx context.Context, ontologyID, conceptID string, depth int) (ontology.Subgraph, error) {
	if depth <= 0 {
		depth = 1
	}
	operation := "graphNeighborhood"
	variables := map[string]any{"ontologyId": ontologyID, "conceptId": conceptID, "depth": depth}
	var result graphNeighborhoodResult
	if err := c.execute(ctx, operation, graphNeighborhoodQuery, variables, &result); err != nil {
		return ontology.Subgraph{}, err
	}
	if err := graphQLErrors(operation, result.Errors); err != nil {
		c.logger.Error("query failed", zap.String("operation", operation), zap.Error(err))
		return ontology.Subgraph{}, err
	}
	subgraph, err := mapSubgraph(result.Data.GraphNeighborhood)
	if err != nil {
		return ontology.Subgraph{}, err
	}
	c.logger.Debug("GraphQL response", zap.String("operation", operation), zap.Int("moduleCount", len(subgraph.Modules)), zap.Int("linkCount", len(subgraph.Links)))
	return subgraph, nil
}

// ClassDescendants resolves FGOS hierarchy, pedagogy concepts, or other class descendants.
func (c *Client) ClassDescendants(ctx context.Context, ontologyID, classID string) ([]ontology.PedagogyConcept, error) {
	operation := "classDescendants"
	variables := map[string]any{"ontologyId": ontologyID, "classId": classID}
	var result classResult
	if err := c.execute(ctx, operation, classDescendantsQuery, variables, &result); err != nil {
		return nil, err
	}
	if err := graphQLErrors(operation, result.Errors); err != nil {
		return nil, err
	}
	return mapConcepts(result.Data.ClassDescendants), nil
}

// ClassTree resolves a full class tree rooted at rootClassID.
func (c *Client) ClassTree(ctx context.Context, ontologyID, rootClassID string) ([]ontology.PedagogyConcept, error) {
	operation := "classTree"
	variables := map[string]any{"ontologyId": ontologyID, "rootClassId": rootClassID}
	var result classResult
	if err := c.execute(ctx, operation, classTreeQuery, variables, &result); err != nil {
		return nil, err
	}
	if err := graphQLErrors(operation, result.Errors); err != nil {
		return nil, err
	}
	return mapConcepts(result.Data.ClassTree), nil
}

// Properties resolves all properties for a Hub ontology entity.
func (c *Client) Properties(ctx context.Context, ontologyID, entityID string) (map[string]string, error) {
	operation := "properties"
	variables := map[string]any{"ontologyId": ontologyID, "entityId": entityID}
	var result propertiesResult
	if err := c.execute(ctx, operation, propertiesQuery, variables, &result); err != nil {
		return nil, err
	}
	if err := graphQLErrors(operation, result.Errors); err != nil {
		return nil, err
	}
	out := make(map[string]string, len(result.Data.Properties))
	for _, p := range result.Data.Properties {
		out[p.Key] = p.Value
	}
	return out, nil
}

// Property resolves one named property for a Hub ontology entity.
func (c *Client) Property(ctx context.Context, ontologyID, entityID, key string) (string, error) {
	operation := "property"
	variables := map[string]any{"ontologyId": ontologyID, "entityId": entityID, "key": key}
	var result propertiesResult
	if err := c.execute(ctx, operation, propertyQuery, variables, &result); err != nil {
		return "", err
	}
	if err := graphQLErrors(operation, result.Errors); err != nil {
		return "", err
	}
	if result.Data.Property == nil {
		return "", fmt.Errorf("property %q not found for entity %q", key, entityID)
	}
	return result.Data.Property.Value, nil
}

func mapSubgraph(dto graphNeighborhoodPayload) (ontology.Subgraph, error) {
	modules := make([]ontology.OntologyModule, 0, len(dto.Modules))
	for _, m := range dto.Modules {
		modules = append(modules, mapModule(m))
	}
	links := make([]ontology.OntologyLink, 0, len(dto.Links))
	for _, l := range dto.Links {
		link := ontology.OntologyLink{
			SourceModuleID: l.SourceModuleID,
			TargetModuleID: l.TargetModuleID,
			Type:           ontology.LinkType(l.Type),
			Metadata:       metadataFromAny(l.Metadata),
		}
		if err := validateLink(link); err != nil {
			return ontology.Subgraph{}, err
		}
		links = append(links, link)
	}
	return ontology.Subgraph{OntologyID: dto.OntologyID, Version: dto.Version, Modules: modules, Links: links}, nil
}

func mapModule(m moduleDTO) ontology.OntologyModule {
	return ontology.OntologyModule{
		ID:           m.ID,
		Title:        m.Title,
		Description:  m.Description,
		Subject:      m.Subject,
		Grade:        m.Grade,
		Version:      m.Version,
		Metadata:     metadataFromAny(m.Metadata),
		FgosBindings: mapFgosBindings(m.FgosBindings),
		Resources:    mapResourceRefs(m.Resources),
		Stories:      mapStoryRefs(m.Stories),
	}
}

func mapFgosBindings(bindings []fgosBindingDTO) []ontology.FgosBinding {
	out := make([]ontology.FgosBinding, 0, len(bindings))
	for _, b := range bindings {
		out = append(out, ontology.FgosBinding{RequirementID: b.RequirementID, Title: b.Title, Level: b.Level, Coverage: b.Coverage})
	}
	return out
}

func mapResourceRefs(resources []resourceRefDTO) []ontology.ResourceRef {
	out := make([]ontology.ResourceRef, 0, len(resources))
	for _, r := range resources {
		out = append(out, ontology.ResourceRef{ID: r.ID, Title: r.Title, Kind: r.Kind, Format: r.Format, Difficulty: r.Difficulty, URI: r.URI, Metadata: metadataFromAny(r.Metadata)})
	}
	return out
}

func mapStoryRefs(stories []storyRefDTO) []ontology.StoryRef {
	out := make([]ontology.StoryRef, 0, len(stories))
	for _, s := range stories {
		out = append(out, ontology.StoryRef{ID: s.ID, Title: s.Title, URI: s.URI, Metadata: metadataFromAny(s.Metadata)})
	}
	return out
}

func mapConcepts(classes []classDTO) []ontology.PedagogyConcept {
	out := make([]ontology.PedagogyConcept, 0, len(classes))
	for _, c := range classes {
		out = append(out, ontology.PedagogyConcept{ID: c.ID, Title: c.Title, Description: c.Description, ParentID: c.ParentID, Metadata: metadataFromAny(c.Metadata)})
	}
	return out
}
