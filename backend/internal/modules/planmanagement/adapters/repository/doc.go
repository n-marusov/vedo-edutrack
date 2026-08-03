// Package repository implements the planmanagement persistence adapter.
//
// Structure (ADR-IMPL.PROCESS.repository-structure §2):
//
//	sqlc/queries.sql      query contract (sqlc-ready, source of truth)
//	sqlc/models.go        row models mirroring the schema tables
//	planmanagement_repository.go  pgx adapter (hand-written)
//	errors.go             SQL error → domain error mapping
//
// Domain entities land in T7/T8; until then the adapter operates on
// sqlc row models (mapper.go arrives with the domain types).
package repository
