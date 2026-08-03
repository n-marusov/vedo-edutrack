-- Package repository: planmanagement persistence adapter (PostgreSQL).
-- This file is the source of truth for SQL queries (sqlc-ready, ADR-IMPL.PROCESS.repository-structure §2).
-- Queries operate on the planmanagement schema (see migrations/000002_planmanagement_init.sql).
-- Hand-written pgx implementation lives in planmanagement_repository.go.

-- name: InsertLearningPlan :one
INSERT INTO planmanagement.learning_plans (
    learner_id, goal_module_id, pedagogy_concept, ontology_version,
    status, timeline_start, timeline_end, snapshot_hash, version, fixed_at
) VALUES (
    @learner_id, @goal_module_id, @pedagogy_concept, @ontology_version,
    @status, @timeline_start, @timeline_end, @snapshot_hash, @version, @fixed_at
)
RETURNING id, created_at;

-- name: GetLearningPlan :one
SELECT id, learner_id, goal_module_id, pedagogy_concept, ontology_version,
       status, timeline_start, timeline_end, snapshot_hash, version, created_at, fixed_at
FROM planmanagement.learning_plans
WHERE id = @id;

-- name: ListLearningPlansByLearner :many
SELECT id, learner_id, goal_module_id, pedagogy_concept, ontology_version,
       status, timeline_start, timeline_end, snapshot_hash, version, created_at, fixed_at
FROM planmanagement.learning_plans
WHERE learner_id = @learner_id
ORDER BY created_at DESC;

-- name: InsertPlanStep :one
INSERT INTO planmanagement.plan_steps (
    plan_id, module_id, position, horizon, is_essential, planned_start, planned_end
) VALUES (
    @plan_id, @module_id, @position, @horizon, @is_essential, @planned_start, @planned_end
)
RETURNING id;

-- name: ListPlanSteps :many
SELECT id, plan_id, module_id, position, horizon, is_essential, planned_start, planned_end
FROM planmanagement.plan_steps
WHERE plan_id = @plan_id
ORDER BY position;

-- name: InsertPlanCheckpoint :one
INSERT INTO planmanagement.plan_checkpoints (
    plan_id, name, checkpoint_date, framework_id, target_coverage_percent
) VALUES (
    @plan_id, @name, @checkpoint_date, @framework_id, @target_coverage_percent
)
RETURNING id;

-- name: ListPlanCheckpoints :many
SELECT id, plan_id, name, checkpoint_date, framework_id, target_coverage_percent
FROM planmanagement.plan_checkpoints
WHERE plan_id = @plan_id
ORDER BY checkpoint_date;
