-- create schema "executionprogress"
CREATE SCHEMA IF NOT EXISTS executionprogress;

-- create "module_progress" (per-learner mastery status per module)
-- The `cause` column holds the deviation reason from
-- {acceleration, more_practice, break, volume_change, unspecified}
-- (REQ-FR-execute.progress.plan-vs-actual AC-2).
CREATE TABLE IF NOT EXISTS executionprogress.module_progress (
    id uuid NOT NULL DEFAULT gen_random_uuid(),
    learner_id uuid NOT NULL,
    module_id text NOT NULL,
    plan_id uuid,
    status text NOT NULL DEFAULT 'not_started', -- not_started | in_progress | mastered | skipped
    mastered_at timestamptz,
    actual_date date,
    planned_date date,
    deviation_days integer,
    deviation_cause text,
    version integer NOT NULL DEFAULT 1,
    updated_at timestamptz NOT NULL DEFAULT now(),
    PRIMARY KEY (id),
    CONSTRAINT module_progress_learner_module_key UNIQUE (learner_id, module_id)
);

-- create "deviation_events" (PlanDeviationDetected records)
-- Dedup on (learner_id, plan_id, deviation_type, detected_at) — one episode
-- = one event (REQ-FR-api.webhooks.plan-deviated AC-4).
CREATE TABLE IF NOT EXISTS executionprogress.deviation_events (
    id uuid NOT NULL DEFAULT gen_random_uuid(),
    event_id uuid NOT NULL,
    learner_id uuid NOT NULL,
    plan_id uuid,
    module_id text,
    deviation_type text NOT NULL, -- schedule | volume
    deviation_value numeric(10,2) NOT NULL,
    threshold numeric(10,2) NOT NULL,
    detected_at timestamptz NOT NULL DEFAULT now(),
    PRIMARY KEY (id),
    CONSTRAINT deviation_events_event_id_key UNIQUE (event_id),
    CONSTRAINT deviation_events_dedup_key UNIQUE (learner_id, plan_id, deviation_type, detected_at)
);

-- create "readiness_forecasts" (binary on-track / not-on-track verdicts)
CREATE TABLE IF NOT EXISTS executionprogress.readiness_forecasts (
    id uuid NOT NULL DEFAULT gen_random_uuid(),
    learner_id uuid NOT NULL,
    plan_id uuid,
    status text NOT NULL, -- on-track | not-on-track
    expected_date date,
    remaining_modules integer NOT NULL DEFAULT 0,
    velocity numeric(8,2) NOT NULL DEFAULT 0, -- modules per week
    key_risks text[] NOT NULL DEFAULT '{}',
    data_confidence text NOT NULL DEFAULT 'high', -- high | low
    created_at timestamptz NOT NULL DEFAULT now(),
    PRIMARY KEY (id),
    CONSTRAINT readiness_forecasts_learner_plan_key UNIQUE (learner_id, plan_id, created_at)
);

-- create indexes for progress lookups
CREATE INDEX IF NOT EXISTS module_progress_learner_id_idx ON executionprogress.module_progress (learner_id);
CREATE INDEX IF NOT EXISTS deviation_events_learner_id_idx ON executionprogress.deviation_events (learner_id);
CREATE INDEX IF NOT EXISTS readiness_forecasts_learner_id_idx ON executionprogress.readiness_forecasts (learner_id);

-- down: DROP TABLE IF EXISTS executionprogress.readiness_forecasts;
-- down: DROP TABLE IF EXISTS executionprogress.deviation_events;
-- down: DROP TABLE IF EXISTS executionprogress.module_progress;
-- down: DROP SCHEMA IF EXISTS executionprogress;
