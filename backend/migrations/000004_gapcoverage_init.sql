-- create schema "gapcoverage"
CREATE SCHEMA IF NOT EXISTS gapcoverage;

-- create "gap_diagnostics" (one root-cause diagnosis run per learner)
CREATE TABLE IF NOT EXISTS gapcoverage.gap_diagnostics (
    id uuid NOT NULL DEFAULT gen_random_uuid(),
    learner_id uuid NOT NULL,
    plan_id uuid,
    ontology_version text NOT NULL,
    analyzed_at timestamptz NOT NULL DEFAULT now(),
    PRIMARY KEY (id),
    CONSTRAINT gap_diagnostics_learner_plan_time_key UNIQUE (learner_id, plan_id, analyzed_at)
);

-- create "gap_roots" (root causes found by climbing strict-prerequisite edges)
CREATE TABLE IF NOT EXISTS gapcoverage.gap_roots (
    id uuid NOT NULL DEFAULT gen_random_uuid(),
    diagnostic_id uuid NOT NULL,
    root_module_id text NOT NULL,
    chain_module_ids text[] NOT NULL DEFAULT '{}',
    blocked_modules_count integer NOT NULL DEFAULT 0,
    blocked_subjects text[] NOT NULL DEFAULT '{}',
    cascade_rank integer NOT NULL DEFAULT 0, -- 1 = highest cascade impact
    PRIMARY KEY (id),
    CONSTRAINT gap_roots_diagnostic_id_fkey FOREIGN KEY (diagnostic_id) REFERENCES gapcoverage.gap_diagnostics (id) ON DELETE CASCADE,
    CONSTRAINT gap_roots_diagnostic_module_key UNIQUE (diagnostic_id, root_module_id)
);

-- create "coverage_snapshots" (FGOS live coverage per learner)
CREATE TABLE IF NOT EXISTS gapcoverage.coverage_snapshots (
    id uuid NOT NULL DEFAULT gen_random_uuid(),
    learner_id uuid NOT NULL,
    framework_id text NOT NULL,
    framework_version text NOT NULL,
    ontology_version text NOT NULL,
    coverage_percent numeric(5,2) NOT NULL DEFAULT 0.00,
    covered_requirements integer NOT NULL DEFAULT 0,
    total_requirements integer NOT NULL DEFAULT 0,
    computed_at timestamptz NOT NULL DEFAULT now(),
    version integer NOT NULL DEFAULT 1,
    PRIMARY KEY (id),
    CONSTRAINT coverage_snapshots_learner_framework_time_key UNIQUE (learner_id, framework_id, computed_at)
);

-- create "coverage_deficits" (uncovered FGOS requirements with blockers)
CREATE TABLE IF NOT EXISTS gapcoverage.coverage_deficits (
    id uuid NOT NULL DEFAULT gen_random_uuid(),
    snapshot_id uuid NOT NULL,
    requirement_id text NOT NULL,
    priority text NOT NULL DEFAULT 'optional', -- strict_prerequisite > essential > optional
    status text NOT NULL DEFAULT 'missing', -- missing | partial | requires_route_extension
    linked_module_ids text[] NOT NULL DEFAULT '{}',
    effort_estimate text,
    PRIMARY KEY (id),
    CONSTRAINT coverage_deficits_snapshot_id_fkey FOREIGN KEY (snapshot_id) REFERENCES gapcoverage.coverage_snapshots (id) ON DELETE CASCADE,
    CONSTRAINT coverage_deficits_snapshot_requirement_key UNIQUE (snapshot_id, requirement_id)
);

-- create indexes for coverage lookups
CREATE INDEX IF NOT EXISTS gap_diagnostics_learner_id_idx ON gapcoverage.gap_diagnostics (learner_id);
CREATE INDEX IF NOT EXISTS coverage_snapshots_learner_id_idx ON gapcoverage.coverage_snapshots (learner_id);

-- down: DROP TABLE IF EXISTS gapcoverage.coverage_deficits;
-- down: DROP TABLE IF EXISTS gapcoverage.coverage_snapshots;
-- down: DROP TABLE IF EXISTS gapcoverage.gap_roots;
-- down: DROP TABLE IF EXISTS gapcoverage.gap_diagnostics;
-- down: DROP SCHEMA IF EXISTS gapcoverage;
