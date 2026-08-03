package gapcoverage

// FgosBinding maps a module to a FGOS requirement.
type FgosBinding struct {
	ModuleID      string
	RequirementID string
}

// Deficit is an uncovered FGOS requirement.
type Deficit struct {
	RequirementID    string
	BlockingModuleID string
}

// CoverageReport is a live FGOS coverage snapshot.
type CoverageReport struct {
	Covered  int
	Total    int
	Percent  float64
	Deficits []Deficit
}

// ComputeCoverage computes FGOS coverage from real mastered modules.
func ComputeCoverage(bindings []FgosBinding, mastery Mastery) CoverageReport {
	if len(bindings) == 0 {
		return CoverageReport{}
	}
	coveredRequirements := map[string]bool{}
	allRequirements := map[string]string{}
	for _, binding := range bindings {
		allRequirements[binding.RequirementID] = binding.ModuleID
		if mastery.Modules[binding.ModuleID] >= 0.8 {
			coveredRequirements[binding.RequirementID] = true
		}
	}
	deficits := []Deficit{}
	for requirementID, moduleID := range allRequirements {
		if !coveredRequirements[requirementID] {
			deficits = append(deficits, Deficit{RequirementID: requirementID, BlockingModuleID: moduleID})
		}
	}
	covered := len(coveredRequirements)
	total := len(allRequirements)
	return CoverageReport{Covered: covered, Total: total, Percent: float64(covered) * 100 / float64(total), Deficits: deficits}
}
