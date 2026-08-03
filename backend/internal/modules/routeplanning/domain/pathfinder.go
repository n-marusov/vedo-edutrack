package routeplanning

import (
	"container/heap"
	"fmt"
	"sort"
)

// ComputeRequest is the pure-domain input for pathfinding.
type ComputeRequest struct {
	PositionID        string
	GoalID            string
	PedagogyConceptID string
	MidHorizonModuleN int
}

// Pathfinder computes deterministic shortest learning paths over ontology snapshots.
type Pathfinder struct {
	weights WeightProfile
}

// NewPathfinder creates a route pathfinder.
func NewPathfinder(weights WeightProfile) *Pathfinder {
	return &Pathfinder{weights: weights}
}

// UnreachableGoalError reports missing prerequisites and a non-empty fallback target.
type UnreachableGoalError struct {
	PositionID           string
	GoalID               string
	MissingPrerequisites []string
	IntermediateGoalID   string
}

func (e *UnreachableGoalError) Error() string {
	return fmt.Sprintf("goal %q is unreachable from %q; missing prerequisites: %v; intermediate goal: %q", e.GoalID, e.PositionID, e.MissingPrerequisites, e.IntermediateGoalID)
}

type weightedEdge struct {
	TargetID string
	Type     LinkType
	Weight   int
}

func (g OntologyGraph) moduleSet() map[string]Module {
	out := make(map[string]Module, len(g.Modules))
	for _, module := range g.Modules {
		out[module.ID] = module
	}
	return out
}

func (g OntologyGraph) adjacency(weights WeightProfile) map[string][]weightedEdge {
	adj := map[string][]weightedEdge{}
	for _, link := range g.Links {
		weight, ok := weights.weight(link.Type)
		if !ok {
			continue
		}
		adj[link.SourceID] = append(adj[link.SourceID], weightedEdge{TargetID: link.TargetID, Type: link.Type, Weight: weight})
	}
	for id := range adj {
		sort.SliceStable(adj[id], func(i, j int) bool {
			if adj[id][i].Weight != adj[id][j].Weight {
				return adj[id][i].Weight < adj[id][j].Weight
			}
			if adj[id][i].Type != adj[id][j].Type {
				return adj[id][i].Type < adj[id][j].Type
			}
			return adj[id][i].TargetID < adj[id][j].TargetID
		})
	}
	return adj
}

// Compute returns the deterministic shortest route. For unreachable goals it returns a non-empty fallback route and *UnreachableGoalError.
func (p *Pathfinder) Compute(graph OntologyGraph, req ComputeRequest) (Route, error) {
	modules := graph.moduleSet()
	if err := validateRequest(modules, req); err != nil {
		return Route{}, err
	}
	if req.PositionID == req.GoalID {
		return routeFromPath([]string{req.PositionID}, nil, 0, req.MidHorizonModuleN), nil
	}

	adj := graph.adjacency(p.weights)
	dijkstra := runDijkstra(adj, modules, req)
	if dijkstra.found {
		return routeFromPath(dijkstra.path, dijkstra.via, dijkstra.dist[req.GoalID], req.MidHorizonModuleN), nil
	}

	fallbackPath := reconstructPath(req.PositionID, dijkstra.lastReachable, dijkstra.prev)
	if len(fallbackPath) == 0 {
		fallbackPath = []string{req.PositionID}
	}
	missing := strictPrerequisitesForGoal(graph, req.GoalID, dijkstra.visited)
	if len(missing) == 0 {
		missing = []string{req.GoalID}
	}
	return routeFromPath(fallbackPath, dijkstra.via, dijkstra.dist[dijkstra.lastReachable], req.MidHorizonModuleN), &UnreachableGoalError{PositionID: req.PositionID, GoalID: req.GoalID, MissingPrerequisites: missing, IntermediateGoalID: dijkstra.lastReachable}
}

// validateRequest checks the compute inputs against the module set.
func validateRequest(modules map[string]Module, req ComputeRequest) error {
	if len(modules) == 0 {
		return fmt.Errorf("ontology graph is empty")
	}
	if req.PositionID == "" || req.GoalID == "" {
		return fmt.Errorf("position and goal are required")
	}
	if _, ok := modules[req.PositionID]; !ok {
		return fmt.Errorf("position module %q not found", req.PositionID)
	}
	if _, ok := modules[req.GoalID]; !ok {
		return fmt.Errorf("goal module %q not found", req.GoalID)
	}
	return nil
}

// dijkstraResult carries the outcome of a single Dijkstra run.
type dijkstraResult struct {
	found         bool
	path          []string
	via           map[string]weightedEdge
	dist          map[string]int
	prev          map[string]string
	visited       map[string]bool
	lastReachable string
}

// runDijkstra executes shortest-path search from position toward goal.
func runDijkstra(adj map[string][]weightedEdge, modules map[string]Module, req ComputeRequest) dijkstraResult {
	dist := map[string]int{req.PositionID: 0}
	prev := map[string]string{}
	via := map[string]weightedEdge{}
	pq := &nodeQueue{}
	heap.Init(pq)
	heap.Push(pq, queueNode{ID: req.PositionID, Cost: 0})
	visited := map[string]bool{}
	lastReachable := req.PositionID

	for pq.Len() > 0 {
		current := heap.Pop(pq).(queueNode)
		if visited[current.ID] {
			continue
		}
		visited[current.ID] = true
		lastReachable = current.ID
		if current.ID == req.GoalID {
			return dijkstraResult{found: true, path: reconstructPath(req.PositionID, req.GoalID, prev), via: via, dist: dist, prev: prev, visited: visited, lastReachable: lastReachable}
		}
		for _, edge := range adj[current.ID] {
			if _, ok := modules[edge.TargetID]; !ok {
				continue
			}
			nextCost := current.Cost + edge.Weight
			oldCost, seen := dist[edge.TargetID]
			if !seen || nextCost < oldCost || (nextCost == oldCost && current.ID < prev[edge.TargetID]) {
				dist[edge.TargetID] = nextCost
				prev[edge.TargetID] = current.ID
				via[edge.TargetID] = edge
				heap.Push(pq, queueNode{ID: edge.TargetID, Cost: nextCost})
			}
		}
	}
	return dijkstraResult{found: false, via: via, dist: dist, prev: prev, visited: visited, lastReachable: lastReachable}
}

func reconstructPath(start, goal string, prev map[string]string) []string {
	path := []string{goal}
	for path[len(path)-1] != start {
		p, ok := prev[path[len(path)-1]]
		if !ok {
			if goal == start {
				return []string{start}
			}
			return nil
		}
		path = append(path, p)
	}
	for i, j := 0, len(path)-1; i < j; i, j = i+1, j-1 {
		path[i], path[j] = path[j], path[i]
	}
	return path
}

func routeFromPath(path []string, via map[string]weightedEdge, totalWeight int, midLimit int) Route {
	steps := make([]RouteStep, 0, len(path))
	for i, id := range path {
		step := RouteStep{ModuleID: id, Order: i}
		if i > 0 {
			edge := via[id]
			step.Via = edge.Type
			step.Weight = edge.Weight
		}
		steps = append(steps, step)
	}
	far, mid, near := buildHorizons(path, midLimit)
	return Route{Steps: steps, TotalWeight: totalWeight, Far: far, Mid: mid, Near: near}
}

func strictPrerequisitesForGoal(graph OntologyGraph, goalID string, reached map[string]bool) []string {
	missing := []string{}
	seen := map[string]bool{}
	var walk func(string)
	walk = func(target string) {
		for _, link := range graph.Links {
			if link.Type != LinkStrictPrerequisite || link.TargetID != target || seen[link.SourceID] {
				continue
			}
			seen[link.SourceID] = true
			if !reached[link.SourceID] {
				missing = append(missing, link.SourceID)
			}
			walk(link.SourceID)
		}
	}
	walk(goalID)
	sort.Strings(missing)
	return missing
}

type queueNode struct {
	ID   string
	Cost int
}

type nodeQueue []queueNode

func (q nodeQueue) Len() int { return len(q) }
func (q nodeQueue) Less(i, j int) bool {
	if q[i].Cost != q[j].Cost {
		return q[i].Cost < q[j].Cost
	}
	return q[i].ID < q[j].ID
}
func (q nodeQueue) Swap(i, j int) { q[i], q[j] = q[j], q[i] }
func (q *nodeQueue) Push(x any) {
	*q = append(*q, x.(queueNode))
}
func (q *nodeQueue) Pop() any {
	old := *q
	n := len(old)
	item := old[n-1]
	*q = old[:n-1]
	return item
}
