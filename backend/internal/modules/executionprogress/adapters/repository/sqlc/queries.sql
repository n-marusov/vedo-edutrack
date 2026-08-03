-- Package repository: executionprogress persistence adapter (PostgreSQL).
-- Source of truth for SQL queries (sqlc-ready, ADR-IMPL.PROCESS.repository-structure §2).
-- Schema: executionprogress (see migrations/000003_executionprogress_init.sql).

-- name: UpsertModuleProgress :one
INSERT INTO executionprogress.module_progress (
    learner_id, module_id, plan_id, status, mastered_at, actual_date,
    planned_date, deviation_days, deviation_cause, version, updated_at
) VALUES (
    @learner_id, @module_id, @plan_id, @status, @mastered_at, @actual_date,
    @planned_date, @deviation_days, @deviation_cause, @version, now()
)
ON CONFLICT (learner_id, module_id) DO UPDATE SET
    status = EXCLUDED.status,
    mastered_at = EXCLUDED.mastered_at,
    actual_date = EXCLUDED.actual_date,
    planned_date = EXCLUDED.planned_date,
    deviation_days = EXCLUDED.deviation_days,
    deviation_cause = EXCLUDED.deviation_cause,
    version = module_progress.version + 1,
    updated_at = now()
RETURNING id, version;

-- name: ListProgressByLearner :many
SELECT id, learner_id, module_id, plan_id, status, mastered_at, actual_date,
       planned_date, deviation_days, deviation_cause, version, updated_at
FROM executionprogress.module_progress
WHERE learner_id = @learner_id
ORDER BY updated_at DESC;

-- name: ListProgressByPlan :many
SELECT id, learner_id, module_id, plan_id, status, mastered_at, actual_date,
       planned_date, deviation_days, deviation_cause, version, updated_at
FROM executionprogress.module_progress
WHERE plan_id = @plan_id
ORDER BY planned_date;

-- name: InsertDeviationEvent :one
INSERT INTO executionprogress.deviation_events (
    event_id, learner_id, plan_id, module_id, deviation_type,
    deviation_value, threshold, detected_at
) VALUES (
    @event_id, @learner_id, @plan_id, @module_id, @deviation_type,
    @deviation_value, @threshold, @detected_at
)
ON CONFLICT (learner_id, plan_id, deviation_type, detected_at) DO NOTHING
RETURNING id;

-- name: InsertReadinessForecast :one
INSERT INTO executionprogress.readiness_forecasts (
    learner_id, plan_id, status, expected_date, remaining_modules,
    velocity, key_risks, data_confidence, created_at
) VALUES (
    @learner_id, @plan_id, @status, @expected_date, @remaining_modules,
    @velocity, @key_risks, @data_confidence, now()
)
RETURNING id;
