// Package repository implements the gapcoverage persistence adapter.
//
// Plan deviation note (T2): same as planmanagement — the sqlc toolchain does
// not exist in this repo; the adapter is hand-written pgx over the query
// contract in sqlc/queries.sql (see planmanagement_repository.go for the
// full rationale). Domain entities land in T12; until then the adapter
// operates on sqlc row models.
package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"

	"vedo-edutrack/backend/internal/modules/gapcoverage/adapters/repository/sqlc"
)

// GapRepository persists gap diagnosis runs with their root causes, and FGOS
// coverage snapshots with their deficit lists.
type GapRepository struct {
	pool *pgxpool.Pool
}

// NewGapRepository builds the repository over a shared connection pool.
func NewGapRepository(pool *pgxpool.Pool) *GapRepository {
	return &GapRepository{pool: pool}
}

// SaveGapAnalysis persists a diagnosis run and all its ranked root causes in
// a single transaction. Returns the diagnostic id.
func (r *GapRepository) SaveGapAnalysis(ctx context.Context, a sqlc.GapAnalysisRow, roots []sqlc.GapRootRow) (string, error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return "", fmt.Errorf("begin save gap diagnostic tx: %w", err)
	}
	defer func() { _ = tx.Rollback(ctx) }()

	var diagnosticID string
	err = tx.QueryRow(ctx,
		`INSERT INTO gapcoverage.gap_diagnostics (
			learner_id, plan_id, ontology_version, analyzed_at
		) VALUES ($1, $2, $3, now())
		RETURNING id`,
		a.LearnerID, a.PlanID, a.OntologyVersion,
	).Scan(&diagnosticID)
	if err != nil {
		return "", fmt.Errorf("insert gap diagnostic: %w", err)
	}

	for _, root := range roots {
		if _, err := tx.Exec(ctx,
			`INSERT INTO gapcoverage.gap_roots (
				diagnostic_id, root_module_id, chain_module_ids,
				blocked_modules_count, blocked_subjects, cascade_rank
			) VALUES ($1, $2, $3, $4, $5, $6)`,
			diagnosticID, root.RootModuleID, nonNilStrings(root.ChainModuleIDs),
			root.BlockedModulesCount, nonNilStrings(root.BlockedSubjects), root.CascadeRank,
		); err != nil {
			return "", fmt.Errorf("insert gap root %q: %w", root.RootModuleID, err)
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return "", fmt.Errorf("commit save gap diagnostic tx: %w", err)
	}
	return diagnosticID, nil
}

// SaveCoverageSnapshot persists a coverage snapshot and its deficits in a
// single transaction. Returns the snapshot id.
func (r *GapRepository) SaveCoverageSnapshot(ctx context.Context, s sqlc.CoverageSnapshotRow, deficits []sqlc.CoverageDeficitRow) (string, error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return "", fmt.Errorf("begin save coverage snapshot tx: %w", err)
	}
	defer func() { _ = tx.Rollback(ctx) }()

	var snapshotID string
	err = tx.QueryRow(ctx,
		`INSERT INTO gapcoverage.coverage_snapshots (
			learner_id, framework_id, framework_version, ontology_version,
			coverage_percent, covered_requirements, total_requirements, computed_at, version
		) VALUES ($1, $2, $3, $4, $5, $6, $7, now(), $8)
		RETURNING id`,
		s.LearnerID, s.FrameworkID, s.FrameworkVersion, s.OntologyVersion,
		s.CoveragePercent, s.CoveredRequirements, s.TotalRequirements, s.Version,
	).Scan(&snapshotID)
	if err != nil {
		return "", fmt.Errorf("insert coverage snapshot: %w", err)
	}

	for _, d := range deficits {
		if _, err := tx.Exec(ctx,
			`INSERT INTO gapcoverage.coverage_deficits (
				snapshot_id, requirement_id, priority, status,
				linked_module_ids, effort_estimate
			) VALUES ($1, $2, $3, $4, $5, $6)`,
			snapshotID, d.RequirementID, d.Priority, d.Status,
			nonNilStrings(d.LinkedModuleIDs), d.EffortEstimate,
		); err != nil {
			return "", fmt.Errorf("insert coverage deficit %q: %w", d.RequirementID, err)
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return "", fmt.Errorf("commit save coverage snapshot tx: %w", err)
	}
	return snapshotID, nil
}

// ListCoverageSnapshotsByLearner returns the learner's coverage history.
func (r *GapRepository) ListCoverageSnapshotsByLearner(ctx context.Context, learnerID string) ([]sqlc.CoverageSnapshotRow, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, learner_id, framework_id, framework_version, ontology_version,
		        coverage_percent, covered_requirements, total_requirements, computed_at, version
		 FROM gapcoverage.coverage_snapshots
		 WHERE learner_id = $1 ORDER BY computed_at DESC`, learnerID)
	if err != nil {
		return nil, fmt.Errorf("list coverage snapshots: %w", err)
	}
	defer rows.Close()

	var out []sqlc.CoverageSnapshotRow
	for rows.Next() {
		var s sqlc.CoverageSnapshotRow
		if err := rows.Scan(&s.ID, &s.LearnerID, &s.FrameworkID, &s.FrameworkVersion, &s.OntologyVersion,
			&s.CoveragePercent, &s.CoveredRequirements, &s.TotalRequirements, &s.ComputedAt, &s.Version); err != nil {
			return nil, fmt.Errorf("scan coverage snapshot: %w", err)
		}
		out = append(out, s)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate coverage snapshots: %w", err)
	}
	return out, nil
}

// nonNilStrings returns an empty slice instead of nil so NOT NULL array
// columns receive an empty array rather than SQL NULL.
func nonNilStrings(s []string) []string {
	if s == nil {
		return []string{}
	}
	return s
}
