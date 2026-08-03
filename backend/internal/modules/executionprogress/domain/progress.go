// Package executionprogress provides the domain model for the executionprogress bounded context.
package executionprogress

import "time"

// MasteryStatus is the learner state for a module.
type MasteryStatus string

const (
	StatusNotStarted MasteryStatus = "not_started"
	StatusInProgress MasteryStatus = "in_progress"
	StatusMastered   MasteryStatus = "mastered"
	StatusSkipped    MasteryStatus = "skipped"
)

// FixedPlan is the execution-progress view of an immutable learning plan.
type FixedPlan struct {
	ID        string
	LearnerID string
	Modules   []PlannedModule
}

// PlannedModule is one module with planned dates.
type PlannedModule struct {
	ModuleID     string
	PlannedStart time.Time
	PlannedEnd   time.Time
}

// ModuleProgress records actual learner progress for a module.
type ModuleProgress struct {
	ModuleID   string
	Status     MasteryStatus
	StartedAt  *time.Time
	MasteredAt *time.Time
	Score      float64
}

// ModuleDeviation compares planned vs actual completion.
type ModuleDeviation struct {
	ModuleID  string
	DeltaDays int
	Percent   float64
}

// PlanDeviationDetected is emitted when actual completion exceeds ±15% of planned timeline.
type PlanDeviationDetected struct {
	ModuleID string
	DeltaPct float64
}

// ComparisonResult contains plan-vs-actual deviations and divergences.
type ComparisonResult struct {
	Deviations  []ModuleDeviation
	Events      []PlanDeviationDetected
	Divergences []string
}

// ComparePlanActual compares progress against an immutable fixed plan.
func ComparePlanActual(plan FixedPlan, progress []ModuleProgress) ComparisonResult {
	planned := map[string]PlannedModule{}
	for _, module := range plan.Modules {
		planned[module.ModuleID] = module
	}
	result := ComparisonResult{}
	for _, actual := range progress {
		module, ok := planned[actual.ModuleID]
		if !ok {
			result.Divergences = append(result.Divergences, actual.ModuleID)
			continue
		}
		if actual.MasteredAt == nil {
			continue
		}
		deltaDays := int(actual.MasteredAt.Sub(module.PlannedEnd).Hours() / 24)
		plannedDays := module.PlannedEnd.Sub(module.PlannedStart).Hours() / 24
		if plannedDays <= 0 {
			plannedDays = 1
		}
		deltaPct := float64(deltaDays) / plannedDays
		if deltaPct < 0 {
			deltaPct = -deltaPct
		}
		deviation := ModuleDeviation{ModuleID: actual.ModuleID, DeltaDays: deltaDays, Percent: deltaPct}
		result.Deviations = append(result.Deviations, deviation)
		if deltaPct > 0.15 {
			result.Events = append(result.Events, PlanDeviationDetected{ModuleID: actual.ModuleID, DeltaPct: deltaPct})
		}
	}
	return result
}
