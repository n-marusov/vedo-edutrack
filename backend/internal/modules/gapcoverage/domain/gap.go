// Package gapcoverage provides the domain model for the gapcoverage bounded context.
package gapcoverage

import "sort"

// LinkType is the gap-diagnosis graph edge type.
type LinkType string

const LinkStrictPrerequisite LinkType = "hasStrictPrerequisite"

// Module is a knowledge graph node used for gap analysis.
type Module struct {
	ID    string
	Title string
}

// Link is a typed directed prerequisite edge.
type Link struct {
	SourceID string
	TargetID string
	Type     LinkType
}

// Graph is the strict-prerequisite graph snapshot.
type Graph struct {
	Modules []Module
	Links   []Link
}

// Mastery contains learner mastery levels from 0.0 to 1.0.
type Mastery struct {
	Modules map[string]float64
}

// DiagnosisStatus describes root-cause diagnosis outcome.
type DiagnosisStatus string

const (
	DiagnosisFound       DiagnosisStatus = "root-causes-found"
	DiagnosisNoRootCause DiagnosisStatus = "no-root-cause-found"
)

// RootCause is an unmastered prerequisite ranked by cascade impact.
type RootCause struct {
	ModuleID       string
	Mastery        float64
	BlockedModules int
}

// DiagnosisResult is the root-cause analysis result.
type DiagnosisResult struct {
	Status     DiagnosisStatus
	RootCauses []RootCause
}

// DiagnoseRootCause climbs strict prerequisite links and ranks unmastered roots by cascade impact.
func DiagnoseRootCause(graph Graph, mastery Mastery, lagModuleID string) DiagnosisResult {
	prereqs := strictPrerequisites(graph)
	unmastered := map[string]bool{}
	seen := map[string]bool{}
	var walkUp func(string)
	walkUp = func(target string) {
		for _, source := range prereqs[target] {
			if seen[source] {
				continue
			}
			seen[source] = true
			if mastery.Modules[source] < 0.8 {
				unmastered[source] = true
			}
			walkUp(source)
		}
	}
	walkUp(lagModuleID)
	if len(unmastered) == 0 {
		return DiagnosisResult{Status: DiagnosisNoRootCause}
	}
	causes := make([]RootCause, 0, len(unmastered))
	for moduleID := range unmastered {
		causes = append(causes, RootCause{ModuleID: moduleID, Mastery: mastery.Modules[moduleID], BlockedModules: blockedCount(graph, moduleID)})
	}
	sort.SliceStable(causes, func(i, j int) bool {
		if causes[i].BlockedModules != causes[j].BlockedModules {
			return causes[i].BlockedModules > causes[j].BlockedModules
		}
		return causes[i].ModuleID < causes[j].ModuleID
	})
	return DiagnosisResult{Status: DiagnosisFound, RootCauses: causes}
}

func strictPrerequisites(graph Graph) map[string][]string {
	out := map[string][]string{}
	for _, link := range graph.Links {
		if link.Type == LinkStrictPrerequisite {
			out[link.TargetID] = append(out[link.TargetID], link.SourceID)
		}
	}
	return out
}

func blockedCount(graph Graph, moduleID string) int {
	children := map[string][]string{}
	for _, link := range graph.Links {
		if link.Type == LinkStrictPrerequisite {
			children[link.SourceID] = append(children[link.SourceID], link.TargetID)
		}
	}
	seen := map[string]bool{}
	var walk func(string)
	walk = func(id string) {
		for _, child := range children[id] {
			if seen[child] {
				continue
			}
			seen[child] = true
			walk(child)
		}
	}
	walk(moduleID)
	return len(seen)
}
