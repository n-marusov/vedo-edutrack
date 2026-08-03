// Package stub provides the M0.3 route-computation stub.
//
// It computes a route by walking the strict-prerequisite chain of the
// in-memory ontology stub (see ontologyport/adapters/stub). Real route
// computation from the VEDO Hub ontology lands in M1 (F1).
package stub

import (
	"fmt"

	ontostub "vedo-edutrack/backend/internal/modules/ontologyport/adapters/stub"
)

// RouteTopic is one step of the computed route.
type RouteTopic struct {
	TopicID  string
	Order    int
	LinkType string
}

// Computer computes fixed routes from the stub graph.
type Computer struct {
	graph *ontostub.Graph
}

// NewComputer creates a route computer over the stub ontology graph.
func NewComputer(graph *ontostub.Graph) *Computer {
	return &Computer{graph: graph}
}

// ComputeRoute returns the strict-prerequisite chain from the root to the
// goal topic, ordered for learning (root first).
func (c *Computer) ComputeRoute(_ string, goalTopicID string) ([]RouteTopic, error) {
	chain, err := c.graph.Prerequisites(goalTopicID)
	if err != nil {
		return nil, err
	}

	out := make([]RouteTopic, 0, len(chain))
	for i, concept := range chain {
		out = append(out, RouteTopic{
			TopicID:  concept.ID,
			Order:    i,
			LinkType: string(ontostub.LinkStrictPrereq),
		})
	}
	if len(out) == 0 {
		return nil, fmt.Errorf("no route found for topic %q", goalTopicID)
	}
	return out, nil
}
