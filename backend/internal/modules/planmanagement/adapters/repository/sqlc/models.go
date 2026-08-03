// Package sqlc holds the query contract for the planmanagement repository.
//
// These structs mirror the planmanagement schema tables and match what
// `sqlc generate` would produce from queries.sql. They are hand-written for
// now (sqlc tooling is not installed — see plan deviation note in
// planmanagement_repository.go) and used as row models by the pgx adapter.
package sqlc

import "time"

// LearningPlanRow mirrors planmanagement.learning_plans.
type LearningPlanRow struct {
	ID              string
	LearnerID       string
	GoalModuleID    string
	PedagogyConcept *string
	OntologyVersion string
	Status          string
	TimelineStart   time.Time
	TimelineEnd     time.Time
	SnapshotHash    string
	Version         int32
	CreatedAt       time.Time
	FixedAt         *time.Time
}

// PlanStepRow mirrors planmanagement.plan_steps.
type PlanStepRow struct {
	ID           string
	PlanID       string
	ModuleID     string
	Position     int32
	Horizon      string
	IsEssential  bool
	PlannedStart time.Time
	PlannedEnd   time.Time
}

// PlanCheckpointRow mirrors planmanagement.plan_checkpoints.
type PlanCheckpointRow struct {
	ID                    string
	PlanID                string
	Name                  string
	CheckpointDate        time.Time
	FrameworkID           *string
	TargetCoveragePercent float64
}
