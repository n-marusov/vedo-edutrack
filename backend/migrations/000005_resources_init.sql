-- create schema "resources"
CREATE SCHEMA IF NOT EXISTS resources;

-- create "resources" (content + enabling resource catalog)
-- Validation (REQ-FR-resource.catalog.bind-to-module):
--   type ∈ {content, enabling}
--   content format ∈ {video, text, interactive, textbook, book}
--   enabling format ∈ {tutor_mentor, lab_equipment, access_pass, money}
--   cost >= 0; source URL is http/https.
CREATE TABLE IF NOT EXISTS resources.resources (
    id uuid NOT NULL DEFAULT gen_random_uuid(),
    title text NOT NULL,
    description text NOT NULL DEFAULT '',
    type text NOT NULL,
    format text NOT NULL,
    source text NOT NULL DEFAULT '',
    style text,
    difficulty text,
    duration_minutes integer,
    cost numeric(12,2) NOT NULL DEFAULT 0.00,
    source_url text,
    version integer NOT NULL DEFAULT 1,
    created_at timestamptz NOT NULL DEFAULT now(),
    PRIMARY KEY (id),
    CONSTRAINT resources_type_check CHECK (type IN ('content', 'enabling')),
    CONSTRAINT resources_cost_check CHECK (cost >= 0)
);

-- create "resource_bindings" (resources bound to ontology modules)
-- Link type mirrors the ontology edge that produced the binding
-- (appliesTo | enriches).
CREATE TABLE IF NOT EXISTS resources.resource_bindings (
    id uuid NOT NULL DEFAULT gen_random_uuid(),
    resource_id uuid NOT NULL,
    module_id text NOT NULL,
    link_type text NOT NULL DEFAULT 'appliesTo', -- appliesTo | enriches
    PRIMARY KEY (id),
    CONSTRAINT resource_bindings_resource_id_fkey FOREIGN KEY (resource_id) REFERENCES resources.resources (id) ON DELETE CASCADE,
    CONSTRAINT resource_bindings_resource_module_key UNIQUE (resource_id, module_id, link_type)
);

-- create indexes for catalog queries
CREATE INDEX IF NOT EXISTS resources_type_format_idx ON resources.resources (type, format);
CREATE INDEX IF NOT EXISTS resource_bindings_module_id_idx ON resources.resource_bindings (module_id);

-- down: DROP TABLE IF EXISTS resources.resource_bindings;
-- down: DROP TABLE IF EXISTS resources.resources;
-- down: DROP SCHEMA IF EXISTS resources;
