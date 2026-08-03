package gapcoverage

import "testing"

func TestDiagnoseRootCauseFindsFirstUnmasteredStrictPrerequisite(t *testing.T) {
	graph := Graph{Modules: []Module{{ID: "percent", Title: "Проценты"}, {ID: "solutions"}, {ID: "chemistry"}}, Links: []Link{
		{SourceID: "percent", TargetID: "solutions", Type: LinkStrictPrerequisite},
		{SourceID: "solutions", TargetID: "chemistry", Type: LinkStrictPrerequisite},
	}}
	mastery := Mastery{Modules: map[string]float64{"percent": 0.70, "solutions": 0.20}}
	result := DiagnoseRootCause(graph, mastery, "chemistry")
	if len(result.RootCauses) == 0 || result.RootCauses[0].ModuleID != "percent" {
		t.Fatalf("expected percent root cause, got %+v", result.RootCauses)
	}
	if result.RootCauses[0].BlockedModules == 0 {
		t.Fatalf("expected cascade impact ranking, got %+v", result.RootCauses[0])
	}
}

func TestDiagnoseRootCauseEmptyChain(t *testing.T) {
	result := DiagnoseRootCause(Graph{Modules: []Module{{ID: "standalone"}}}, Mastery{Modules: map[string]float64{"standalone": 0.1}}, "standalone")
	if result.Status != DiagnosisNoRootCause {
		t.Fatalf("status=%s, want no root cause", result.Status)
	}
}

func TestCoverageAndDeficits(t *testing.T) {
	bindings := []FgosBinding{{ModuleID: "m1", RequirementID: "r1"}, {ModuleID: "m2", RequirementID: "r2"}}
	report := ComputeCoverage(bindings, Mastery{Modules: map[string]float64{"m1": 1.0}})
	if report.Percent != 50 || len(report.Deficits) != 1 || report.Deficits[0].RequirementID != "r2" {
		t.Fatalf("unexpected coverage report: %+v", report)
	}
}
