package routeplanning

import (
	"math/rand"
	"testing"
)

func TestPathfinderShortestPathDeterministic(t *testing.T) {
	graph := OntologyGraph{
		Modules: []Module{{ID: "start"}, {ID: "soft"}, {ID: "strict"}, {ID: "goal"}},
		Links: []Link{
			{SourceID: "start", TargetID: "soft", Type: LinkSoftPrerequisite},
			{SourceID: "soft", TargetID: "goal", Type: LinkSoftPrerequisite},
			{SourceID: "start", TargetID: "strict", Type: LinkStrictPrerequisite},
			{SourceID: "strict", TargetID: "goal", Type: LinkStrictPrerequisite},
		},
	}
	pf := NewPathfinder(DefaultWeightProfile())
	var first Route
	for i := 0; i < 10; i++ {
		route, err := pf.Compute(graph, ComputeRequest{PositionID: "start", GoalID: "goal"})
		if err != nil {
			t.Fatalf("compute route: %v", err)
		}
		if got, want := route.ModuleIDs(), []string{"start", "strict", "goal"}; !equalStrings(got, want) {
			t.Fatalf("route=%v, want %v", got, want)
		}
		if i == 0 {
			first = route
			continue
		}
		if !equalStrings(route.ModuleIDs(), first.ModuleIDs()) {
			t.Fatalf("non-deterministic route on run %d: %v vs %v", i, route.ModuleIDs(), first.ModuleIDs())
		}
	}
}

func TestPathfinderStrictPrerequisitesAppearBeforeGoal(t *testing.T) {
	graph := OntologyGraph{
		Modules: []Module{{ID: "percent"}, {ID: "solutions"}, {ID: "chemistry"}},
		Links: []Link{
			{SourceID: "percent", TargetID: "solutions", Type: LinkStrictPrerequisite},
			{SourceID: "solutions", TargetID: "chemistry", Type: LinkStrictPrerequisite},
		},
	}
	route, err := NewPathfinder(DefaultWeightProfile()).Compute(graph, ComputeRequest{PositionID: "percent", GoalID: "chemistry"})
	if err != nil {
		t.Fatalf("compute route: %v", err)
	}
	if !route.SatisfiesStrictPrerequisites(graph) {
		t.Fatalf("route violates strict prerequisites: %v", route.ModuleIDs())
	}
}

func TestPathfinderConflictResolutionPrefersStrictThenSoftOverEnrich(t *testing.T) {
	graph := OntologyGraph{Modules: []Module{{ID: "a"}, {ID: "b"}, {ID: "c"}, {ID: "goal"}}, Links: []Link{
		{SourceID: "a", TargetID: "goal", Type: LinkEnriches},
		{SourceID: "a", TargetID: "b", Type: LinkSoftPrerequisite},
		{SourceID: "b", TargetID: "goal", Type: LinkSoftPrerequisite},
		{SourceID: "a", TargetID: "c", Type: LinkStrictPrerequisite},
		{SourceID: "c", TargetID: "goal", Type: LinkStrictPrerequisite},
	}}
	route, err := NewPathfinder(DefaultWeightProfile()).Compute(graph, ComputeRequest{PositionID: "a", GoalID: "goal"})
	if err != nil {
		t.Fatalf("compute route: %v", err)
	}
	if got, want := route.ModuleIDs(), []string{"a", "c", "goal"}; !equalStrings(got, want) {
		t.Fatalf("route=%v, want %v", got, want)
	}
}

func TestPathfinderUnreachableGoalReturnsMissingPrerequisitesAndIntermediate(t *testing.T) {
	graph := OntologyGraph{Modules: []Module{{ID: "start"}, {ID: "island"}, {ID: "goal"}}, Links: []Link{
		{SourceID: "island", TargetID: "goal", Type: LinkStrictPrerequisite},
	}}
	route, err := NewPathfinder(DefaultWeightProfile()).Compute(graph, ComputeRequest{PositionID: "start", GoalID: "goal"})
	if err == nil {
		t.Fatal("expected unreachable error")
	}
	unreachable, ok := err.(*UnreachableGoalError)
	if !ok {
		t.Fatalf("error type=%T, want *UnreachableGoalError", err)
	}
	if len(unreachable.MissingPrerequisites) == 0 || unreachable.IntermediateGoalID == "" {
		t.Fatalf("expected missing prerequisites and intermediate goal, got %+v", unreachable)
	}
	if len(route.Steps) == 0 {
		t.Fatalf("expected non-empty fallback route")
	}
}

func TestPathfinderEdgeCases(t *testing.T) {
	pf := NewPathfinder(DefaultWeightProfile())
	if _, err := pf.Compute(OntologyGraph{}, ComputeRequest{PositionID: "a", GoalID: "b"}); err == nil {
		t.Fatal("expected empty ontology error")
	}
	single := OntologyGraph{Modules: []Module{{ID: "same"}}}
	route, err := pf.Compute(single, ComputeRequest{PositionID: "same", GoalID: "same"})
	if err != nil {
		t.Fatalf("goal=position: %v", err)
	}
	if got, want := route.ModuleIDs(), []string{"same"}; !equalStrings(got, want) {
		t.Fatalf("route=%v, want %v", got, want)
	}
}

func TestPathfinderMatchesBruteForceOnRandomDAGs(t *testing.T) {
	pf := NewPathfinder(DefaultWeightProfile())
	for seed := int64(0); seed < 100; seed++ {
		graph := randomDAG(seed, 8)
		route, err := pf.Compute(graph, ComputeRequest{PositionID: "m0", GoalID: "m7"})
		if err != nil {
			continue
		}
		wantCost := bruteForceShortestCost(graph, "m0", "m7", DefaultWeightProfile())
		if route.TotalWeight != wantCost {
			t.Fatalf("seed=%d cost=%d, want %d route=%v", seed, route.TotalWeight, wantCost, route.ModuleIDs())
		}
	}
}

func BenchmarkPathfinder5kModules(b *testing.B) {
	graph := randomDAG(42, 5000)
	pf := NewPathfinder(DefaultWeightProfile())
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = pf.Compute(graph, ComputeRequest{PositionID: "m0", GoalID: "m4999"})
	}
}

func randomDAG(seed int64, n int) OntologyGraph {
	r := rand.New(rand.NewSource(seed))
	modules := make([]Module, n)
	links := make([]Link, 0, n*2)
	for i := 0; i < n; i++ {
		modules[i] = Module{ID: moduleID(i)}
		if i > 0 {
			links = append(links, Link{SourceID: moduleID(i - 1), TargetID: moduleID(i), Type: LinkStrictPrerequisite})
		}
	}
	for i := 0; i < n; i++ {
		for j := i + 2; j < n && j < i+5; j++ {
			if r.Intn(3) == 0 {
				links = append(links, Link{SourceID: moduleID(i), TargetID: moduleID(j), Type: LinkSoftPrerequisite})
			}
		}
	}
	return OntologyGraph{Modules: modules, Links: links}
}

func moduleID(i int) string { return "m" + itoa(i) }

func itoa(i int) string {
	if i == 0 {
		return "0"
	}
	buf := [20]byte{}
	pos := len(buf)
	for i > 0 {
		pos--
		buf[pos] = byte('0' + i%10)
		i /= 10
	}
	return string(buf[pos:])
}

func bruteForceShortestCost(graph OntologyGraph, start, goal string, weights WeightProfile) int {
	adj := graph.adjacency(weights)
	best := 1 << 30
	var walk func(string, int, map[string]bool)
	walk = func(id string, cost int, seen map[string]bool) {
		if cost >= best {
			return
		}
		if id == goal {
			best = cost
			return
		}
		for _, edge := range adj[id] {
			if seen[edge.TargetID] {
				continue
			}
			seen[edge.TargetID] = true
			walk(edge.TargetID, cost+edge.Weight, seen)
			delete(seen, edge.TargetID)
		}
	}
	walk(start, 0, map[string]bool{start: true})
	return best
}

func equalStrings(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
