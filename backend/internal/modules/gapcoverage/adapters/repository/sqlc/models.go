// Package sqlc holds the query contract for the gapcoverage repository.
//
// Row models mirror the gapcoverage schema tables (what `sqlc generate`
// would produce from queries.sql). Hand-written until the sqlc toolchain is
// installed; used by the pgx adapter in the parent package.
package sqlc

import "time"

// GapAnalysisRow mirrors gapcoverage.gap_diagnostics.
type GapAnalysisRow struct {
	ID              string
	LearnerID       string
	PlanID          *string
	OntologyVersion string
	AnalyzedAt      time.Time
}

// GapRootRow mirrors gapcoverage.gap_roots.
type GapRootRow struct {
	ID                  string
	DiagnosticID        string
	RootModuleID        string
	ChainModuleIDs      []string
	BlockedModulesCount int32
	BlockedSubjects     []string
	CascadeRank         int32
}

// CoverageSnapshotRow mirrors gapcoverage.coverage_snapshots.
type CoverageSnapshotRow struct {
	ID                  string
	LearnerID           string
	FrameworkID         string
	FrameworkVersion    string
	OntologyVersion     string
	CoveragePercent     float64
	CoveredRequirements int32
	TotalRequirements   int32
	ComputedAt          time.Time
	Version             int32
}

// CoverageDeficitRow mirrors gapcoverage.coverage_deficits.
type CoverageDeficitRow struct {
	ID              string
	SnapshotID      string
	RequirementID   string
	Priority        string
	Status          string
	LinkedModuleIDs []string
	EffortEstimate  *string
}
