// Package repository implements the executionprogress persistence adapter.
//
// Plan deviation note (T2): same as planmanagement — the sqlc toolchain does
// not exist in this repo; the adapter is hand-written pgx over the query
// contract in sqlc/queries.sql (see planmanagement_repository.go for the
// full rationale). Domain entities land in T11; until then the adapter
// operates on sqlc row models.
package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"vedo-edutrack/backend/internal/modules/executionprogress/adapters/repository/sqlc"
)

// ErrNotFound is returned when a requested entity does not exist.
var ErrNotFound = errors.New("not found")

// ProgressRepository persists module progress, deviation events and
// readiness forecasts.
type ProgressRepository struct {
	pool *pgxpool.Pool
}

// NewProgressRepository builds the repository over a shared connection pool.
func NewProgressRepository(pool *pgxpool.Pool) *ProgressRepository {
	return &ProgressRepository{pool: pool}
}

// UpsertModuleProgress records or updates a learner's mastery status for a
// module (optimistic lock: version increments on every update).
func (r *ProgressRepository) UpsertModuleProgress(ctx context.Context, p sqlc.ModuleProgressRow) error {
	var id string
	var version int32
	err := r.pool.QueryRow(ctx,
		`INSERT INTO executionprogress.module_progress (
			learner_id, module_id, plan_id, status, mastered_at, actual_date,
			planned_date, deviation_days, deviation_cause, version, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, now())
		ON CONFLICT (learner_id, module_id) DO UPDATE SET
			status = EXCLUDED.status,
			mastered_at = EXCLUDED.mastered_at,
			actual_date = EXCLUDED.actual_date,
			planned_date = EXCLUDED.planned_date,
			deviation_days = EXCLUDED.deviation_days,
			deviation_cause = EXCLUDED.deviation_cause,
			version = module_progress.version + 1,
			updated_at = now()
		RETURNING id, version`,
		p.LearnerID, p.ModuleID, p.PlanID, p.Status, p.MasteredAt, p.ActualDate,
		p.PlannedDate, p.DeviationDays, p.DeviationCause, p.Version,
	).Scan(&id, &version)
	if err != nil {
		return fmt.Errorf("upsert module progress (learner=%s module=%s): %w", p.LearnerID, p.ModuleID, err)
	}
	return nil
}

// ListProgressByLearner returns the learner's module progress rows.
func (r *ProgressRepository) ListProgressByLearner(ctx context.Context, learnerID string) ([]sqlc.ModuleProgressRow, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, learner_id, module_id, plan_id, status, mastered_at, actual_date,
		        planned_date, deviation_days, deviation_cause, version, updated_at
		 FROM executionprogress.module_progress
		 WHERE learner_id = $1 ORDER BY updated_at DESC`, learnerID)
	if err != nil {
		return nil, fmt.Errorf("list module progress: %w", err)
	}
	defer rows.Close()
	return scanProgressRows(rows)
}

// ListProgressByPlan returns progress rows for one fixed plan.
func (r *ProgressRepository) ListProgressByPlan(ctx context.Context, planID string) ([]sqlc.ModuleProgressRow, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, learner_id, module_id, plan_id, status, mastered_at, actual_date,
		        planned_date, deviation_days, deviation_cause, version, updated_at
		 FROM executionprogress.module_progress
		 WHERE plan_id = $1 ORDER BY planned_date`, planID)
	if err != nil {
		return nil, fmt.Errorf("list module progress by plan: %w", err)
	}
	defer rows.Close()
	return scanProgressRows(rows)
}

// InsertDeviationEvent records a PlanDeviationDetected event. Idempotent on
// the dedup key (learner_id, plan_id, deviation_type, detected_at).
// Returns the row id, or (false, nil) when the event already existed.
func (r *ProgressRepository) InsertDeviationEvent(ctx context.Context, e sqlc.DeviationEventRow) (string, bool, error) {
	var id string
	err := r.pool.QueryRow(ctx,
		`INSERT INTO executionprogress.deviation_events (
			event_id, learner_id, plan_id, module_id, deviation_type,
			deviation_value, threshold, detected_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		ON CONFLICT (learner_id, plan_id, deviation_type, detected_at) DO NOTHING
		RETURNING id`,
		e.EventID, e.LearnerID, e.PlanID, e.ModuleID, e.DeviationType,
		e.DeviationValue, e.Threshold, e.DetectedAt,
	).Scan(&id)
	if errors.Is(err, pgx.ErrNoRows) {
		return "", false, nil // duplicate — already recorded
	}
	if err != nil {
		return "", false, fmt.Errorf("insert deviation event: %w", err)
	}
	return id, true, nil
}

// SaveReadinessForecast persists a binary readiness verdict.
func (r *ProgressRepository) SaveReadinessForecast(ctx context.Context, f sqlc.ReadinessForecastRow) (string, error) {
	var id string
	err := r.pool.QueryRow(ctx,
		`INSERT INTO executionprogress.readiness_forecasts (
			learner_id, plan_id, status, expected_date, remaining_modules,
			velocity, key_risks, data_confidence, created_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, now())
		RETURNING id`,
		f.LearnerID, f.PlanID, f.Status, f.ExpectedDate, f.RemainingModules,
		f.Velocity, nonNilStringSlice(f.KeyRisks), f.DataConfidence,
	).Scan(&id)
	if err != nil {
		return "", fmt.Errorf("insert readiness forecast: %w", err)
	}
	return id, nil
}

// nonNilStringSlice returns an empty slice instead of nil so NOT NULL array
// columns receive an empty array rather than SQL NULL.
func nonNilStringSlice(s []string) []string {
	if s == nil {
		return []string{}
	}
	return s
}

// rowScanner abstracts pgx.Row and pgx.Rows for row scanning.
type rowScanner interface {
	Scan(dest ...any) error
}

func scanProgressRow(row rowScanner) (sqlc.ModuleProgressRow, error) {
	var p sqlc.ModuleProgressRow
	err := row.Scan(&p.ID, &p.LearnerID, &p.ModuleID, &p.PlanID, &p.Status, &p.MasteredAt,
		&p.ActualDate, &p.PlannedDate, &p.DeviationDays, &p.DeviationCause, &p.Version, &p.UpdatedAt)
	return p, err
}

func scanProgressRows(rows pgx.Rows) ([]sqlc.ModuleProgressRow, error) {
	defer rows.Close()
	var out []sqlc.ModuleProgressRow
	for rows.Next() {
		p, err := scanProgressRow(rows)
		if err != nil {
			return nil, fmt.Errorf("scan module progress row: %w", err)
		}
		out = append(out, p)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate module progress rows: %w", err)
	}
	return out, nil
}
