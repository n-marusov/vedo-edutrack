package auth

import (
	"context"
)

// GetUserID returns the authenticated user ID from the context.
// Returns an empty string if the context has no auth claims.
func GetUserID(ctx context.Context) string {
	claims := ClaimsFrom(ctx)
	if claims == nil {
		return ""
	}
	return claims.UserID
}

// GetRoles returns the authenticated user's roles from the context.
// Returns nil if the context has no auth claims.
func GetRoles(ctx context.Context) []string {
	claims := ClaimsFrom(ctx)
	if claims == nil {
		return nil
	}
	return claims.Roles
}
