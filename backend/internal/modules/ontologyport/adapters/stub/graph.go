// Package stub provides a hardcoded in-memory ontology graph for M0.3.
//
// It is a demo stand-in for the real VEDO Hub ontology reads (F0) that land
// in M1. The graph covers 5th-grade math topics with the five EduTrack link
// types (hasStrictPrerequisite, hasSoftPrerequisite, enriches, appliesTo,
// isAnalogousTo).
package stub

import (
	"fmt"
)

// LinkType is an ontology relation between topics.
type LinkType string

// The five EduTrack link types (vision.md, F0).
const (
	LinkStrictPrereq  LinkType = "hasStrictPrerequisite"
	LinkSoftPrereq    LinkType = "hasSoftPrerequisite"
	LinkEnriches      LinkType = "enriches"
	LinkAppliesTo     LinkType = "appliesTo"
	LinkIsAnalogousTo LinkType = "isAnalogousTo"
)

// Concept is a knowledge graph node (a topic).
type Concept struct {
	ID          string
	Title       string
	Description string
	Links       []Link
}

// Link is a typed edge from this concept to another.
type Link struct {
	TopicID  string
	LinkType LinkType
}

// Graph is the in-memory ontology stub.
type Graph struct {
	concepts map[string]*Concept
}

// NewGraph returns the fixed 5th-grade math stub graph.
func NewGraph() *Graph {
	g := &Graph{concepts: map[string]*Concept{}}
	for _, c := range seedConcepts() {
		g.concepts[c.ID] = c
	}
	return g
}

// seedConcepts builds the hardcoded topic graph (math, grade 5).
func seedConcepts() []*Concept {
	return []*Concept{
		{
			ID:          "math-5-1",
			Title:       "Натуральные числа",
			Description: "Natural numbers, place value, comparison.",
			Links: []Link{
				{TopicID: "math-5-2", LinkType: LinkStrictPrereq},
			},
		},
		{
			ID:          "math-5-2",
			Title:       "Сложение и вычитание натуральных чисел",
			Description: "Addition and subtraction of natural numbers.",
			Links: []Link{
				{TopicID: "math-5-3", LinkType: LinkStrictPrereq},
				{TopicID: "math-5-4", LinkType: LinkSoftPrereq},
			},
		},
		{
			ID:          "math-5-3",
			Title:       "Умножение натуральных чисел",
			Description: "Multiplication of natural numbers.",
			Links: []Link{
				{TopicID: "math-5-5", LinkType: LinkStrictPrereq},
			},
		},
		{
			ID:          "math-5-4",
			Title:       "Вычитание и его свойства",
			Description: "Subtraction and its properties.",
			Links: []Link{
				{TopicID: "math-5-5", LinkType: LinkSoftPrereq},
			},
		},
		{
			ID:          "math-5-5",
			Title:       "Деление натуральных чисел",
			Description: "Division of natural numbers, divisibility.",
			Links: []Link{
				{TopicID: "math-5-6", LinkType: LinkStrictPrereq},
				{TopicID: "math-5-7", LinkType: LinkEnriches},
			},
		},
		{
			ID:          "math-5-6",
			Title:       "Обыкновенные дроби",
			Description: "Common fractions, comparison, operations.",
			Links: []Link{
				{TopicID: "math-5-8", LinkType: LinkStrictPrereq},
			},
		},
		{
			ID:          "math-5-7",
			Title:       "Делимость чисел",
			Description: "Divisibility rules, prime numbers.",
			Links: []Link{
				{TopicID: "math-5-8", LinkType: LinkSoftPrereq},
			},
		},
		{
			ID:          "math-5-8",
			Title:       "Сложение и вычитание дробей",
			Description: "Adding and subtracting fractions.",
			Links: []Link{
				{TopicID: "math-5-9", LinkType: LinkStrictPrereq},
			},
		},
		{
			ID:          "math-5-9",
			Title:       "Умножение и деление дробей",
			Description: "Multiplying and dividing fractions.",
			Links: []Link{
				{TopicID: "math-5-10", LinkType: LinkStrictPrereq},
			},
		},
		{
			ID:          "math-5-10",
			Title:       "Десятичные дроби",
			Description: "Decimal fractions, operations with decimals.",
			Links: []Link{
				{TopicID: "math-5-11", LinkType: LinkStrictPrereq},
				{TopicID: "math-5-12", LinkType: LinkAppliesTo},
			},
		},
		{
			ID:          "math-5-11",
			Title:       "Проценты",
			Description: "Percentages, percentage problems.",
			Links: []Link{
				{TopicID: "math-5-12", LinkType: LinkStrictPrereq},
			},
		},
		{
			ID:          "math-5-12",
			Title:       "Текстовые задачи",
			Description: "Word problems: rates, proportions, percentages.",
			Links: []Link{
				{TopicID: "math-5-13", LinkType: LinkEnriches},
			},
		},
		{
			ID:          "math-5-13",
			Title:       "Геометрические фигуры",
			Description: "Basic geometric figures: segments, angles, perimeter.",
			Links: []Link{
				{TopicID: "math-5-14", LinkType: LinkStrictPrereq},
				{TopicID: "math-5-12", LinkType: LinkIsAnalogousTo},
			},
		},
		{
			ID:          "math-5-14",
			Title:       "Площадь и объём",
			Description: "Area of figures, volume of solids.",
			Links:       []Link{},
		},
	}
}

// GetConcept returns a concept by ID (nil if absent).
func (g *Graph) GetConcept(topicID string) *Concept {
	return g.concepts[topicID]
}

// GetAllConcepts returns all concepts (unordered).
func (g *Graph) GetAllConcepts() []*Concept {
	out := make([]*Concept, 0, len(g.concepts))
	for _, c := range g.concepts {
		out = append(out, c)
	}
	return out
}

// Prerequisites walks the hasStrictPrerequisite chain backwards from a topic
// to its roots (breadth-first, ordered root-first). Returns an error when the
// goal topic is unknown.
//
// Link semantics: a link {X, hasStrictPrerequisite} on concept C means
// "C is a strict prerequisite of X" — i.e. links point from prerequisite to
// consequent. Building the route therefore traverses the reverse index.
func (g *Graph) Prerequisites(goalTopicID string) ([]*Concept, error) {
	if g.concepts[goalTopicID] == nil {
		return nil, fmt.Errorf("topic %q not found in stub graph", goalTopicID)
	}

	// Reverse index: consequent → set of strict prerequisite ids.
	reverse := map[string][]string{}
	for _, c := range g.concepts {
		for _, l := range c.Links {
			if l.LinkType == LinkStrictPrereq {
				reverse[l.TopicID] = append(reverse[l.TopicID], c.ID)
			}
		}
	}

	// BFS from the goal backwards over strict prerequisites.
	var chain []*Concept
	seen := map[string]bool{}
	queue := []string{goalTopicID}
	for len(queue) > 0 {
		id := queue[0]
		queue = queue[1:]
		if seen[id] {
			continue
		}
		seen[id] = true
		c := g.concepts[id]
		if c == nil {
			continue
		}
		chain = append(chain, c)
		queue = append(queue, reverse[id]...)
	}

	// Reverse to root-first order.
	for i, j := 0, len(chain)-1; i < j; i, j = i+1, j-1 {
		chain[i], chain[j] = chain[j], chain[i]
	}
	return chain, nil
}
