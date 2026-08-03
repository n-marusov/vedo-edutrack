// Package sqlc holds the query contract for the executionprogress repository.
//
// Row models mirror the executionprogress schema tables (what `sqlc generate`
// would produce from queries.sql). Hand-written until the sqlc toolchain is
// installed; used by the pgx adapter in the parent package.
package sqlc

import "time"

// ModuleProgressRow mirrors executionprogress.module_progress.
type ModuleProgressRow struct {
	ID             string
	LearnerID      string
	ModuleID       string
	PlanID         *string
	Status         string
	MasteredAt     *time.Time
	ActualDate     *time.Time
	PlannedDate    *time.Time
	DeviationDays  *int32
	DeviationCause *string
	Version        int32
	UpdatedAt      time.Time
}

// DeviationEventRow mirrors executionprogress.deviation_events.
type DeviationEventRow struct {
	ID             string
	EventID        string
	LearnerID      string
	PlanID         *string
	ModuleID       *string
	DeviationType  string
	DeviationValue float64
	Threshold      float64
	DetectedAt     time.Time
}

// ReadinessForecastRow mirrors executionprogress.readiness_forecasts.
type ReadinessForecastRow struct {
	ID               string
	LearnerID        string
	PlanID           *string
	Status           string
	ExpectedDate     *time.Time
	RemainingModules int32
	Velocity         float64
	KeyRisks         []string
	DataConfidence   string
	CreatedAt        time.Time
}
