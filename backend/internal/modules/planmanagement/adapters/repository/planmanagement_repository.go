// Package repository implements the planmanagement persistence adapter.
//
// Plan deviation note (T2): the M1 plan assumed an existing sqlc toolchain
// ("follow the identity-access sqlc pattern") — that toolchain does not
// exist in the repo (no sqlc.yaml, no sqlc in go.mod). The identity-access
// module uses raw pgx (see internal/cli/seed.go). This adapter therefore
// implements the query surface with hand-written pgx, keeping the SQL
// contract in sqlc/queries.sql so `sqlc generate` can adopt it later with
// zero query rewrites. Domain entities land in T7/T8; until then the
// repository operates on sqlc row models.
package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"vedo-edutrack/backend/internal/modules/planmanagement/adapters/repository/sqlc"
)

// ErrNotFound is returned when a requested entity does not exist.
var ErrNotFound = errors.New("not found")

// PlanRepository persists LearningPlan snapshots and their steps/checkpoints.
type PlanRepository struct {
	pool *pgxpool.Pool
}

// NewPlanRepository builds the repository over a shared connection pool.
func NewPlanRepository(pool *pgxpool.Pool) *PlanRepository {
	return &PlanRepository{pool: pool}
}

// SavePlan inserts an immutable plan snapshot (plan + steps + checkpoints)
// in a single transaction (ADR-DES.DATA.storage-strategy: INSERT-only).
func (r *PlanRepository) SavePlan(ctx context.Context, plan sqlc.LearningPlanRow, steps []sqlc.PlanStepRow, checkpoints []sqlc.PlanCheckpointRow) (string, error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return "", fmt.Errorf("begin save plan tx: %w", err)
	}
	defer func() { _ = tx.Rollback(ctx) }()

	var planID string
	err = tx.QueryRow(ctx,
		`INSERT INTO planmanagement.learning_plans (
			learner_id, goal_module_id, pedagogy_concept, ontology_version,
			status, timeline_start, timeline_end, snapshot_hash, version, fixed_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		RETURNING id`,
		plan.LearnerID, plan.GoalModuleID, plan.PedagogyConcept, plan.OntologyVersion,
		plan.Status, plan.TimelineStart, plan.TimelineEnd, plan.SnapshotHash,
		plan.Version, plan.FixedAt,
	).Scan(&planID)
	if err != nil {
		return "", fmt.Errorf("insert learning plan: %w", err)
	}

	for _, s := range steps {
		if _, err := tx.Exec(ctx,
			`INSERT INTO planmanagement.plan_steps (
				plan_id, module_id, position, horizon, is_essential, planned_start, planned_end
			) VALUES ($1, $2, $3, $4, $5, $6, $7)`,
			planID, s.ModuleID, s.Position, s.Horizon, s.IsEssential, s.PlannedStart, s.PlannedEnd,
		); err != nil {
			return "", fmt.Errorf("insert plan step %d: %w", s.Position, err)
		}
	}

	for _, c := range checkpoints {
		if _, err := tx.Exec(ctx,
			`INSERT INTO planmanagement.plan_checkpoints (
				plan_id, name, checkpoint_date, framework_id, target_coverage_percent
			) VALUES ($1, $2, $3, $4, $5)`,
			planID, c.Name, c.CheckpointDate, c.FrameworkID, c.TargetCoveragePercent,
		); err != nil {
			return "", fmt.Errorf("insert plan checkpoint %q: %w", c.Name, err)
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return "", fmt.Errorf("commit save plan tx: %w", err)
	}
	return planID, nil
}

// GetPlan loads a plan snapshot together with its steps and checkpoints.
func (r *PlanRepository) GetPlan(ctx context.Context, planID string) (sqlc.LearningPlanRow, []sqlc.PlanStepRow, []sqlc.PlanCheckpointRow, error) {
	plan, err := r.getPlanRow(ctx, planID)
	if err != nil {
		return sqlc.LearningPlanRow{}, nil, nil, err
	}
	steps, err := r.listSteps(ctx, planID)
	if err != nil {
		return sqlc.LearningPlanRow{}, nil, nil, err
	}
	checkpoints, err := r.listCheckpoints(ctx, planID)
	if err != nil {
		return sqlc.LearningPlanRow{}, nil, nil, err
	}
	return plan, steps, checkpoints, nil
}

// ListPlansByLearner returns the learner's plan snapshots, newest first.
func (r *PlanRepository) ListPlansByLearner(ctx context.Context, learnerID string) ([]sqlc.LearningPlanRow, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, learner_id, goal_module_id, pedagogy_concept, ontology_version,
		        status, timeline_start, timeline_end, snapshot_hash, version, created_at, fixed_at
		 FROM planmanagement.learning_plans
		 WHERE learner_id = $1
		 ORDER BY created_at DESC`, learnerID)
	if err != nil {
		return nil, fmt.Errorf("list learning plans: %w", err)
	}
	defer rows.Close()

	var plans []sqlc.LearningPlanRow
	for rows.Next() {
		p, err := scanPlanRow(rows)
		if err != nil {
			return nil, err
		}
		plans = append(plans, p)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate learning plans: %w", err)
	}
	return plans, nil
}

func (r *PlanRepository) getPlanRow(ctx context.Context, planID string) (sqlc.LearningPlanRow, error) {
	row := r.pool.QueryRow(ctx,
		`SELECT id, learner_id, goal_module_id, pedagogy_concept, ontology_version,
		        status, timeline_start, timeline_end, snapshot_hash, version, created_at, fixed_at
		 FROM planmanagement.learning_plans WHERE id = $1`, planID)
	plan, err := scanPlanRow(row)
	if errors.Is(err, pgx.ErrNoRows) {
		return sqlc.LearningPlanRow{}, fmt.Errorf("%w: plan %s", ErrNotFound, planID)
	}
	if err != nil {
		return sqlc.LearningPlanRow{}, fmt.Errorf("get learning plan %s: %w", planID, err)
	}
	return plan, nil
}

func (r *PlanRepository) listSteps(ctx context.Context, planID string) ([]sqlc.PlanStepRow, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, plan_id, module_id, position, horizon, is_essential, planned_start, planned_end
		 FROM planmanagement.plan_steps WHERE plan_id = $1 ORDER BY position`, planID)
	if err != nil {
		return nil, fmt.Errorf("list plan steps: %w", err)
	}
	defer rows.Close()

	var steps []sqlc.PlanStepRow
	for rows.Next() {
		var s sqlc.PlanStepRow
		if err := rows.Scan(&s.ID, &s.PlanID, &s.ModuleID, &s.Position, &s.Horizon, &s.IsEssential, &s.PlannedStart, &s.PlannedEnd); err != nil {
			return nil, fmt.Errorf("scan plan step: %w", err)
		}
		steps = append(steps, s)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate plan steps: %w", err)
	}
	return steps, nil
}

func (r *PlanRepository) listCheckpoints(ctx context.Context, planID string) ([]sqlc.PlanCheckpointRow, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, plan_id, name, checkpoint_date, framework_id, target_coverage_percent
		 FROM planmanagement.plan_checkpoints WHERE plan_id = $1 ORDER BY checkpoint_date`, planID)
	if err != nil {
		return nil, fmt.Errorf("list plan checkpoints: %w", err)
	}
	defer rows.Close()

	var cps []sqlc.PlanCheckpointRow
	for rows.Next() {
		var c sqlc.PlanCheckpointRow
		if err := rows.Scan(&c.ID, &c.PlanID, &c.Name, &c.CheckpointDate, &c.FrameworkID, &c.TargetCoveragePercent); err != nil {
			return nil, fmt.Errorf("scan plan checkpoint: %w", err)
		}
		cps = append(cps, c)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate plan checkpoints: %w", err)
	}
	return cps, nil
}

// rowScanner abstracts pgx.Row and pgx.Rows for scanPlanRow.
type rowScanner interface {
	Scan(dest ...any) error
}

func scanPlanRow(row rowScanner) (sqlc.LearningPlanRow, error) {
	var p sqlc.LearningPlanRow
	err := row.Scan(&p.ID, &p.LearnerID, &p.GoalModuleID, &p.PedagogyConcept, &p.OntologyVersion,
		&p.Status, &p.TimelineStart, &p.TimelineEnd, &p.SnapshotHash, &p.Version, &p.CreatedAt, &p.FixedAt)
	return p, err
}
