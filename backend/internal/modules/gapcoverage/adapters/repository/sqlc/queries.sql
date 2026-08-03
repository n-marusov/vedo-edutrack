-- Package repository: gapcoverage persistence adapter (PostgreSQL).
-- Source of truth for SQL queries (sqlc-ready, ADR-IMPL.PROCESS.repository-structure §2).
-- Schema: gapcoverage (see migrations/000004_gapcoverage_init.sql).

-- name: InsertGapDiagnostic :one
INSERT INTO gapcoverage.gap_diagnostics (
    learner_id, plan_id, ontology_version, analyzed_at
) VALUES (
    @learner_id, @plan_id, @ontology_version, now()
)
RETURNING id, analyzed_at;

-- name: InsertGapRoot :one
INSERT INTO gapcoverage.gap_roots (
    diagnostic_id, root_module_id, chain_module_ids,
    blocked_modules_count, blocked_subjects, cascade_rank
) VALUES (
    @diagnostic_id, @root_module_id, @chain_module_ids,
    @blocked_modules_count, @blocked_subjects, @cascade_rank
)
RETURNING id;

-- name: ListGapRootsByDiagnostic :many
SELECT id, diagnostic_id, root_module_id, chain_module_ids,
       blocked_modules_count, blocked_subjects, cascade_rank
FROM gapcoverage.gap_roots
WHERE diagnostic_id = @diagnostic_id
ORDER BY cascade_rank;

-- name: InsertCoverageSnapshot :one
INSERT INTO gapcoverage.coverage_snapshots (
    learner_id, framework_id, framework_version, ontology_version,
    coverage_percent, covered_requirements, total_requirements, computed_at, version
) VALUES (
    @learner_id, @framework_id, @framework_version, @ontology_version,
    @coverage_percent, @covered_requirements, @total_requirements, now(), @version
)
RETURNING id, computed_at;

-- name: InsertCoverageDeficit :one
INSERT INTO gapcoverage.coverage_deficits (
    snapshot_id, requirement_id, priority, status,
    linked_module_ids, effort_estimate
) VALUES (
    @snapshot_id, @requirement_id, @priority, @status,
    @linked_module_ids, @effort_estimate
)
RETURNING id;

-- name: ListCoverageDeficitsBySnapshot :many
SELECT id, snapshot_id, requirement_id, priority, status,
       linked_module_ids, effort_estimate
FROM gapcoverage.coverage_deficits
WHERE snapshot_id = @snapshot_id
ORDER BY priority DESC, linked_module_ids;

-- name: ListCoverageSnapshotsByLearner :many
SELECT id, learner_id, framework_id, framework_version, ontology_version,
       coverage_percent, covered_requirements, total_requirements, computed_at, version
FROM gapcoverage.coverage_snapshots
WHERE learner_id = @learner_id
ORDER BY computed_at DESC;
