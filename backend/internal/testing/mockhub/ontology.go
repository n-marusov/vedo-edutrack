package mockhub

import (
	"fmt"
	"sort"
	"strings"
)

// Ontology is the in-memory model of a parsed .ttl ontology (TBox).
type Ontology struct {
	Classes    map[string]*Class
	Properties map[string]*Property
}

// Class is an OWL class (knowledge module / concept).
type Class struct {
	ID           string
	Label        string
	Comment      string
	Parents      []string
	Children     []string
	IsAbstract   bool
	IsDeprecated bool
}

// Property is an OWL object/datatype property (relation).
type Property struct {
	ID         string
	Label      string
	Comment    string
	Type       string // owl:ObjectProperty | owl:DatatypeProperty
	Domains    []string
	Ranges     []string
	XSDType    string
	Functional bool
}

// NewOntology creates an empty ontology.
func NewOntology() *Ontology {
	return &Ontology{
		Classes:    map[string]*Class{},
		Properties: map[string]*Property{},
	}
}

// Counts returns the number of classes and properties.
func (o *Ontology) Counts() (classes, properties int) {
	return len(o.Classes), len(o.Properties)
}

// Breadcrumb returns the chain of class labels from the root to the given
// class (via rdfs:subClassOf), or nil if the class is unknown.
func (o *Ontology) Breadcrumb(classID string) []string {
	seen := map[string]bool{}
	var chain []string
	cur := classID
	for cur != "" && !seen[cur] {
		c := o.Classes[cur]
		if c == nil {
			break
		}
		chain = append([]string{c.LabelOrID()}, chain...)
		seen[cur] = true
		if len(c.Parents) == 0 {
			break
		}
		cur = c.Parents[0]
	}
	return chain
}

// Tree returns an indented tree of classes rooted at the class with no
// parents (or all roots).
func (o *Ontology) Tree() string {
	roots := o.Roots()
	if len(roots) == 0 {
		return "(empty)"
	}
	var sb strings.Builder
	for _, root := range roots {
		o.writeTree(&sb, root, 0, map[string]bool{})
	}
	return sb.String()
}

// Roots returns class IDs with no parents.
func (o *Ontology) Roots() []string {
	var roots []string
	for id, c := range o.Classes {
		if len(c.Parents) == 0 {
			roots = append(roots, id)
		}
	}
	sort.Strings(roots)
	return roots
}

func (o *Ontology) writeTree(sb *strings.Builder, id string, depth int, seen map[string]bool) {
	if seen[id] {
		return
	}
	seen[id] = true
	c := o.Classes[id]
	if c == nil {
		return
	}
	sb.WriteString(strings.Repeat("  ", depth))
	sb.WriteString(c.LabelOrID())
	sb.WriteString("\n")
	children := append([]string{}, c.Children...)
	sort.Strings(children)
	for _, ch := range children {
		o.writeTree(sb, ch, depth+1, seen)
	}
}

// Descendants returns class IDs at depth ≤ maxDepth below the given class
// (excluding the class itself). maxDepth ≤ 0 means unlimited.
func (o *Ontology) Descendants(classID string, maxDepth int) []string {
	var out []string
	seen := map[string]bool{}
	var walk func(id string, depth int)
	walk = func(id string, depth int) {
		if maxDepth > 0 && depth > maxDepth {
			return
		}
		c := o.Classes[id]
		if c == nil {
			return
		}
		for _, ch := range c.Children {
			if seen[ch] {
				continue
			}
			seen[ch] = true
			out = append(out, ch)
			walk(ch, depth+1)
		}
	}
	walk(classID, 1)
	return out
}

// Autocomplete returns up to `limit` classes whose label contains the query
// (case-insensitive).
func (o *Ontology) Autocomplete(q string, limit int) []string {
	if limit <= 0 {
		limit = 20
	}
	q = strings.ToLower(q)
	var out []string
	for id, c := range o.Classes {
		if q == "" || strings.Contains(strings.ToLower(c.Label), q) {
			out = append(out, id)
			if len(out) >= limit {
				break
			}
		}
	}
	sort.Strings(out)
	return out
}

// LabelOrID returns the label if set, otherwise the (short) ID.
func (c *Class) LabelOrID() string {
	if c.Label != "" {
		return c.Label
	}
	return shortIRI(c.ID)
}

// LabelOrID returns the label if set, otherwise the (short) ID.
func (p *Property) LabelOrID() string {
	if p.Label != "" {
		return p.Label
	}
	return shortIRI(p.ID)
}

// shortIRI shortens an IRI to its local name (after # or last /).
func shortIRI(iri string) string {
	if idx := strings.LastIndexAny(iri, "#/"); idx >= 0 && idx < len(iri)-1 {
		return iri[idx+1:]
	}
	return iri
}

// ClassByName finds a class by exact ID. Used by resolvers.
func (o *Ontology) ClassByName(id string) (*Class, error) {
	c := o.Classes[id]
	if c == nil {
		return nil, fmt.Errorf("class %q not found", id)
	}
	return c, nil
}
