// Package repository implements the resources persistence adapter.
//
// Plan deviation note (T2): same as planmanagement — the sqlc toolchain does
// not exist in this repo; the adapter is hand-written pgx over the query
// contract in sqlc/queries.sql (see planmanagement_repository.go for the
// full rationale). Domain entities land in T14; until then the adapter
// operates on sqlc row models.
package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"vedo-edutrack/backend/internal/modules/resources/adapters/repository/sqlc"
)

// ErrNotFound is returned when a requested entity does not exist.
var ErrNotFound = errors.New("not found")

// ResourceRepository persists the resource catalog and module bindings.
type ResourceRepository struct {
	pool *pgxpool.Pool
}

// NewResourceRepository builds the repository over a shared connection pool.
func NewResourceRepository(pool *pgxpool.Pool) *ResourceRepository {
	return &ResourceRepository{pool: pool}
}

// CreateResource inserts a new catalog entry.
func (r *ResourceRepository) CreateResource(ctx context.Context, res sqlc.ResourceRow) (string, error) {
	var id string
	err := r.pool.QueryRow(ctx,
		`INSERT INTO resources.resources (
			title, description, type, format, source, style, difficulty,
			duration_minutes, cost, source_url, version, created_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, now())
		RETURNING id`,
		res.Title, res.Description, res.Type, res.Format, res.Source, res.Style,
		res.Difficulty, res.DurationMinutes, res.Cost, res.SourceURL, res.Version,
	).Scan(&id)
	if err != nil {
		return "", fmt.Errorf("insert resource %q: %w", res.Title, err)
	}
	return id, nil
}

// GetResource loads one catalog entry.
func (r *ResourceRepository) GetResource(ctx context.Context, id string) (sqlc.ResourceRow, error) {
	row := r.pool.QueryRow(ctx,
		`SELECT id, title, description, type, format, source, style, difficulty,
		        duration_minutes, cost, source_url, version, created_at
		 FROM resources.resources WHERE id = $1`, id)
	res, err := scanResourceRow(row)
	if errors.Is(err, pgx.ErrNoRows) {
		return sqlc.ResourceRow{}, fmt.Errorf("%w: resource %s", ErrNotFound, id)
	}
	if err != nil {
		return sqlc.ResourceRow{}, fmt.Errorf("get resource %s: %w", id, err)
	}
	return res, nil
}

// ResourceFilter restricts catalog queries (NULL = no constraint).
type ResourceFilter struct {
	Type   *string
	Format *string
	Source *string
	Limit  int32
	Offset int32
}

// ListResources queries the catalog with optional filters and pagination.
func (r *ResourceRepository) ListResources(ctx context.Context, f ResourceFilter) ([]sqlc.ResourceRow, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, title, description, type, format, source, style, difficulty,
		        duration_minutes, cost, source_url, version, created_at
		 FROM resources.resources
		 WHERE ($1::text IS NULL OR type = $1)
		   AND ($2::text IS NULL OR format = $2)
		   AND ($3::text IS NULL OR source = $3)
		 ORDER BY title
		 LIMIT $4 OFFSET $5`,
		f.Type, f.Format, f.Source, f.Limit, f.Offset)
	if err != nil {
		return nil, fmt.Errorf("list resources: %w", err)
	}
	defer rows.Close()

	var out []sqlc.ResourceRow
	for rows.Next() {
		res, err := scanResourceRow(rows)
		if err != nil {
			return nil, fmt.Errorf("scan resource row: %w", err)
		}
		out = append(out, res)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate resource rows: %w", err)
	}
	return out, nil
}

// BindResourceToModule creates a binding (idempotent on the triple). Returns
// (id, created, err); created=false when the binding already existed.
func (r *ResourceRepository) BindResourceToModule(ctx context.Context, b sqlc.ResourceBindingRow) (string, bool, error) {
	var id string
	err := r.pool.QueryRow(ctx,
		`INSERT INTO resources.resource_bindings (resource_id, module_id, link_type)
		 VALUES ($1, $2, $3)
		 ON CONFLICT (resource_id, module_id, link_type) DO NOTHING
		 RETURNING id`,
		b.ResourceID, b.ModuleID, b.LinkType,
	).Scan(&id)
	if errors.Is(err, pgx.ErrNoRows) {
		return "", false, nil // duplicate binding
	}
	if err != nil {
		return "", false, fmt.Errorf("bind resource %s to module %s: %w", b.ResourceID, b.ModuleID, err)
	}
	return id, true, nil
}

// ListResourcesByModule returns resources bound to a module with their link
// type.
func (r *ResourceRepository) ListResourcesByModule(ctx context.Context, moduleID string) ([]sqlc.ResourceRow, []string, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT r.id, r.title, r.description, r.type, r.format, r.source, r.style,
		        r.difficulty, r.duration_minutes, r.cost, r.source_url, r.version, r.created_at,
		        b.link_type
		 FROM resources.resource_bindings b
		 JOIN resources.resources r ON r.id = b.resource_id
		 WHERE b.module_id = $1
		 ORDER BY r.title`, moduleID)
	if err != nil {
		return nil, nil, fmt.Errorf("list resources by module: %w", err)
	}
	defer rows.Close()

	var out []sqlc.ResourceRow
	var linkTypes []string
	for rows.Next() {
		var res sqlc.ResourceRow
		var linkType string
		if err := rows.Scan(&res.ID, &res.Title, &res.Description, &res.Type, &res.Format, &res.Source,
			&res.Style, &res.Difficulty, &res.DurationMinutes, &res.Cost, &res.SourceURL, &res.Version,
			&res.CreatedAt, &linkType); err != nil {
			return nil, nil, fmt.Errorf("scan bound resource row: %w", err)
		}
		out = append(out, res)
		linkTypes = append(linkTypes, linkType)
	}
	if err := rows.Err(); err != nil {
		return nil, nil, fmt.Errorf("iterate bound resource rows: %w", err)
	}
	return out, linkTypes, nil
}

// rowScanner abstracts pgx.Row and pgx.Rows for row scanning.
type rowScanner interface {
	Scan(dest ...any) error
}

func scanResourceRow(row rowScanner) (sqlc.ResourceRow, error) {
	var res sqlc.ResourceRow
	err := row.Scan(&res.ID, &res.Title, &res.Description, &res.Type, &res.Format, &res.Source,
		&res.Style, &res.Difficulty, &res.DurationMinutes, &res.Cost, &res.SourceURL, &res.Version,
		&res.CreatedAt)
	return res, err
}
