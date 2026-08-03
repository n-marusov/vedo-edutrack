// Package sqlc holds the query contract for the resources repository.
//
// Row models mirror the resources schema tables (what `sqlc generate`
// would produce from queries.sql). Hand-written until the sqlc toolchain is
// installed; used by the pgx adapter in the parent package.
package sqlc

import "time"

// ResourceRow mirrors resources.resources.
type ResourceRow struct {
	ID              string
	Title           string
	Description     string
	Type            string // content | enabling
	Format          string
	Source          string
	Style           *string
	Difficulty      *string
	DurationMinutes *int32
	Cost            float64
	SourceURL       *string
	Version         int32
	CreatedAt       time.Time
}

// ResourceBindingRow mirrors resources.resource_bindings.
type ResourceBindingRow struct {
	ID         string
	ResourceID string
	ModuleID   string
	LinkType   string // appliesTo | enriches
}
