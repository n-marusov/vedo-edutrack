-- create schema "planmanagement"
CREATE SCHEMA IF NOT EXISTS planmanagement;

-- create "learning_plans" (immutable snapshot of a fixed route)
-- Per ADR-DES.DATA.storage-strategy: plans are INSERT-only, versioned
-- snapshots. The `version` column enables optimistic locking (409 on conflict).
CREATE TABLE IF NOT EXISTS planmanagement.learning_plans (
    id uuid NOT NULL DEFAULT gen_random_uuid(),
    learner_id uuid NOT NULL,
    goal_module_id text NOT NULL,
    pedagogy_concept text,
    ontology_version text NOT NULL,
    status text NOT NULL DEFAULT 'active', -- fixed | active | superseded
    timeline_start date NOT NULL,
    timeline_end date NOT NULL,
    snapshot_hash text NOT NULL,
    version integer NOT NULL DEFAULT 1,
    created_at timestamptz NOT NULL DEFAULT now(),
    fixed_at timestamptz,
    PRIMARY KEY (id)
);

-- create "plan_steps" (one route step of a fixed plan)
CREATE TABLE IF NOT EXISTS planmanagement.plan_steps (
    id uuid NOT NULL DEFAULT gen_random_uuid(),
    plan_id uuid NOT NULL,
    module_id text NOT NULL,
    position integer NOT NULL,
    horizon text NOT NULL DEFAULT 'far', -- far | mid | near
    is_essential boolean NOT NULL DEFAULT false,
    planned_start date NOT NULL,
    planned_end date NOT NULL,
    PRIMARY KEY (id),
    CONSTRAINT plan_steps_plan_id_fkey FOREIGN KEY (plan_id) REFERENCES planmanagement.learning_plans (id) ON DELETE CASCADE,
    CONSTRAINT plan_steps_plan_position_key UNIQUE (plan_id, position)
);

-- create "plan_checkpoints" (FGOS attestation / milestone checkpoints)
CREATE TABLE IF NOT EXISTS planmanagement.plan_checkpoints (
    id uuid NOT NULL DEFAULT gen_random_uuid(),
    plan_id uuid NOT NULL,
    name text NOT NULL,
    checkpoint_date date NOT NULL,
    framework_id text,
    target_coverage_percent numeric(5,2) NOT NULL DEFAULT 100.00,
    PRIMARY KEY (id),
    CONSTRAINT plan_checkpoints_plan_id_fkey FOREIGN KEY (plan_id) REFERENCES planmanagement.learning_plans (id) ON DELETE CASCADE,
    CONSTRAINT plan_checkpoints_plan_name_key UNIQUE (plan_id, name)
);

-- create index on learner for plan lookups
CREATE INDEX IF NOT EXISTS learning_plans_learner_id_idx ON planmanagement.learning_plans (learner_id);

-- down: DROP TABLE IF EXISTS planmanagement.plan_checkpoints;
-- down: DROP TABLE IF EXISTS planmanagement.plan_steps;
-- down: DROP TABLE IF EXISTS planmanagement.learning_plans;
-- down: DROP SCHEMA IF EXISTS planmanagement;
