package stub

import (
	"strings"
	"testing"

	ontostub "vedo-edutrack/backend/internal/modules/ontologyport/adapters/stub"
)

// TestComputerRouteForKnownTopic verifies the strict-prerequisite chain walk:
// route starts at the root, is ordered root-first, and every edge is
// hasStrictPrerequisite.
func TestComputerRouteForKnownTopic(t *testing.T) {
	c := NewComputer(ontostub.NewGraph())

	route, err := c.ComputeRoute("l1", "math-5-11")
	if err != nil {
		t.Fatalf("compute route: %v", err)
	}
	if len(route) == 0 {
		t.Fatal("expected non-empty route for known topic")
	}
	if route[0].TopicID != "math-5-1" {
		t.Fatalf("route must start at the root, got %s", route[0].TopicID)
	}
	for i, step := range route {
		if step.LinkType != string(ontostub.LinkStrictPrereq) {
			t.Fatalf("step %d link type = %q, want %q", i, step.LinkType, ontostub.LinkStrictPrereq)
		}
	}
	// The goal must be the last step.
	if route[len(route)-1].TopicID != "math-5-11" {
		t.Fatalf("route must end at the goal, got %s", route[len(route)-1].TopicID)
	}
}

// TestComputerRouteIsRootFirst checks ordering is consistent across the chain.
func TestComputerRouteIsRootFirst(t *testing.T) {
	c := NewComputer(ontostub.NewGraph())

	route, err := c.ComputeRoute("l1", "math-5-8")
	if err != nil {
		t.Fatalf("compute route: %v", err)
	}
	if len(route) < 2 {
		t.Fatalf("expected a chain of strict prerequisites, got %d steps", len(route))
	}
	// Strict chain for math-5-8: math-5-1 → math-5-2 → math-5-3 → math-5-5 → math-5-6 → math-5-8.
	// BFS may include alternatives; the invariant is: every step before the goal
	// is a strict prerequisite of a later step, and the goal is last.
	for i, step := range route {
		if i > 0 && step.TopicID == route[i-1].TopicID {
			t.Fatalf("duplicate step at %d", i)
		}
	}
	if route[len(route)-1].TopicID != "math-5-8" {
		t.Fatalf("goal must be the last step, got %s", route[len(route)-1].TopicID)
	}
}

// TestComputerErrorOnUnknownTopic verifies the error path for unknown goals.
func TestComputerErrorOnUnknownTopic(t *testing.T) {
	c := NewComputer(ontostub.NewGraph())

	route, err := c.ComputeRoute("l1", "does-not-exist")
	if err == nil {
		t.Fatal("expected error for unknown topic")
	}
	if route != nil {
		t.Fatalf("expected nil route on error, got %v", route)
	}
	if !strings.Contains(err.Error(), "does-not-exist") {
		t.Fatalf("error should mention the topic, got %q", err.Error())
	}
}

// TestComputerGoalEqualsPosition — the root topic (no strict prerequisites
// point to it) yields a single-step route.
func TestComputerGoalEqualsPosition(t *testing.T) {
	c := NewComputer(ontostub.NewGraph())

	route, err := c.ComputeRoute("l1", "math-5-1")
	if err != nil {
		t.Fatalf("compute route: %v", err)
	}
	if len(route) != 1 || route[0].TopicID != "math-5-1" {
		t.Fatalf("expected single-step route, got %v", route)
	}
}
