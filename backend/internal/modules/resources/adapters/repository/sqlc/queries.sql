-- Package repository: resources persistence adapter (PostgreSQL).
-- Source of truth for SQL queries (sqlc-ready, ADR-IMPL.PROCESS.repository-structure §2).
-- Schema: resources (see migrations/000005_resources_init.sql).

-- name: CreateResource :one
INSERT INTO resources.resources (
    title, description, type, format, source, style, difficulty,
    duration_minutes, cost, source_url, version, created_at
) VALUES (
    @title, @description, @type, @format, @source, @style, @difficulty,
    @duration_minutes, @cost, @source_url, @version, now()
)
RETURNING id, created_at;

-- name: GetResource :one
SELECT id, title, description, type, format, source, style, difficulty,
       duration_minutes, cost, source_url, version, created_at
FROM resources.resources
WHERE id = @id;

-- name: ListResources :many
SELECT id, title, description, type, format, source, style, difficulty,
       duration_minutes, cost, source_url, version, created_at
FROM resources.resources
WHERE (@type::text IS NULL OR type = @type)
  AND (@format::text IS NULL OR format = @format)
  AND (@source::text IS NULL OR source = @source)
ORDER BY title
LIMIT @limit_ OFFSET @offset_;

-- name: BindResourceToModule :one
INSERT INTO resources.resource_bindings (resource_id, module_id, link_type)
VALUES (@resource_id, @module_id, @link_type)
ON CONFLICT (resource_id, module_id, link_type) DO NOTHING
RETURNING id;

-- name: ListResourcesByModule :many
SELECT r.id, r.title, r.description, r.type, r.format, r.source, r.style,
       r.difficulty, r.duration_minutes, r.cost, r.source_url, r.version, r.created_at,
       b.link_type
FROM resources.resource_bindings b
JOIN resources.resources r ON r.id = b.resource_id
WHERE b.module_id = @module_id
ORDER BY r.title;
