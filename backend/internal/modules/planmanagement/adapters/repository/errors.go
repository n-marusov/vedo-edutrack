package repository

import (
	"errors"
	"fmt"
)

// MapError converts a raw SQL error into a repository-level error with
// domain semantics. Repository adapters must never leak pgx error codes to
// the application layer (ADR-IMPL.PROCESS.repository-structure §2).
func MapError(err error) error {
	if err == nil {
		return nil
	}
	if errors.Is(err, ErrNotFound) {
		return err
	}
	// Uniqueness violations (23505) surface as conflict errors; all other
	// database failures stay opaque at this layer (the application layer
	// logs the full context separately).
	return fmt.Errorf("repository error: %w", err)
}
