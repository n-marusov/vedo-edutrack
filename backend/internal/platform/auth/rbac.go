package auth

import (
	"encoding/json"
	"net/http"

	"go.uber.org/zap"
)

// Archetype is a role archetype from the RBAC model
// (REQ-NFR-security.compliance.role-catalog, ADR-DES.SECURITY.rbac-model).
type Archetype string

// Archetypes from the RBAC model.
const (
	ArchetypeSelf            Archetype = "self"
	ArchetypeDependentsOwner Archetype = "dependents-owner"
	ArchetypeStaff           Archetype = "staff"
	ArchetypeManagement      Archetype = "management"
	ArchetypeIntegration     Archetype = "integration"
	ArchetypeAdmin           Archetype = "admin"
	ArchetypeOps             Archetype = "ops"
)

// Permission is a functional-area × action pair from the permission matrix.
type Permission string

// Permissions (functional area × action), see
// REQ-NFR-security.compliance.permission-matrix.
const (
	PermRouteCompute  Permission = "route:read" // route-compute is read-only (route is a function)
	PermRouteWrite    Permission = "route:write"
	PermPlanRead      Permission = "plan:read"
	PermPlanManage    Permission = "plan:manage"
	PermProgressRead  Permission = "progress:read"
	PermGapRead       Permission = "gap:read"
	PermCoverageRead  Permission = "coverage:read"
	PermOntologyRead  Permission = "ontology:read"
	PermResourceRead  Permission = "resource:read"
	PermResourceMng   Permission = "resource:manage"
	PermUserRead      Permission = "user:read"
	PermUserManage    Permission = "user:manage"
	PermWebhookConfig Permission = "webhook:configure"
)

// permissionMatrix maps archetype → granted permissions. Deny-by-default:
// anything not granted here is denied.
//
// Materialized from REQ-NFR-security.compliance.permission-matrix (read-only
// for product archetypes on route/plan/progress/gap/coverage/ontology;
// plan-manage/resource-manage/user-manage/webhook-configure only for
// managing archetypes; ops has no product permissions).
var permissionMatrix = map[Archetype]map[Permission]bool{
	ArchetypeSelf: {
		PermRouteCompute: true, PermPlanRead: true, PermProgressRead: true,
		PermGapRead: true, PermCoverageRead: true, PermResourceRead: true,
		PermUserRead: true, PermOntologyRead: true,
	},
	ArchetypeDependentsOwner: {
		PermRouteCompute: true, PermPlanRead: true, PermProgressRead: true,
		PermGapRead: true, PermCoverageRead: true, PermResourceRead: true,
		PermUserRead: true, PermOntologyRead: true,
	},
	ArchetypeStaff: {
		PermRouteCompute: true, PermPlanRead: true, PermPlanManage: true,
		PermProgressRead: true, PermGapRead: true, PermCoverageRead: true,
		PermUserRead: true, PermOntologyRead: true,
	},
	ArchetypeManagement: {
		PermRouteCompute: true, PermPlanRead: true, PermPlanManage: true,
		PermProgressRead: true, PermGapRead: true, PermCoverageRead: true,
		PermResourceRead: true, PermResourceMng: true, PermUserRead: true,
		PermOntologyRead: true,
	},
	ArchetypeIntegration: {
		PermRouteCompute: true, PermPlanRead: true, PermProgressRead: true,
		PermGapRead: true, PermCoverageRead: true, PermOntologyRead: true,
		PermWebhookConfig: true,
	},
	ArchetypeAdmin: {
		PermRouteCompute: true, PermPlanRead: true, PermPlanManage: true,
		PermProgressRead: true, PermGapRead: true, PermCoverageRead: true,
		PermResourceRead: true, PermResourceMng: true, PermUserRead: true,
		PermUserManage: true, PermOntologyRead: true, PermWebhookConfig: true,
		PermRouteWrite: true,
	},
	// ops: no product permissions (infrastructure only, JIT + 2-person rule).
	ArchetypeOps: {},
}

// roleArchetype maps role instance names to their archetype. Seeded roles
// inherit their archetype's permissions (role-catalog instances,
// REQ-NFR-security.compliance.role-catalog).
var roleArchetype = map[string]Archetype{
	"learner":             ArchetypeSelf,
	"employee":            ArchetypeSelf,
	"parent":              ArchetypeDependentsOwner,
	"teacher":             ArchetypeStaff,
	"methodologist":       ArchetypeManagement,
	"school-director":     ArchetypeManagement,
	"hr-manager":          ArchetypeManagement,
	"platform-integrator": ArchetypeIntegration,
	"admin":               ArchetypeAdmin,
	"ops":                 ArchetypeOps,
}

// HasPermission reports whether any of the given roles grants the permission.
// O(1) map lookup, deny-by-default.
func HasPermission(roles []string, perm Permission) bool {
	for _, role := range roles {
		arch, ok := roleArchetype[role]
		if !ok {
			continue
		}
		if permissionMatrix[arch][perm] {
			return true
		}
	}
	return false
}

// RequirePermission returns a chi middleware that enforces the permission for
// authenticated requests (roles from the auth context, T5). Denied requests
// receive 403; unauthenticated requests are rejected by the auth middleware
// earlier (401).
func RequirePermission(perm Permission, logger *zap.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			roles := GetRoles(r.Context())
			userID := GetUserID(r.Context())

			if !HasPermission(roles, perm) {
				if logger != nil {
					logger.Warn("access denied",
						zap.String("user_id", userID),
						zap.String("permission", string(perm)),
						zap.Strings("roles", roles),
						zap.String("path", r.URL.Path),
					)
				}
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusForbidden)
				_ = json.NewEncoder(w).Encode(map[string]string{
					"error":      "forbidden",
					"permission": string(perm),
				})
				return
			}

			if logger != nil {
				logger.Info("access granted",
					zap.String("user_id", userID),
					zap.String("permission", string(perm)),
					zap.String("path", r.URL.Path),
				)
			}
			next.ServeHTTP(w, r)
		})
	}
}

// AllArchetypes lists every archetype (for seed materialization).
func AllArchetypes() []Archetype {
	return []Archetype{
		ArchetypeSelf, ArchetypeDependentsOwner, ArchetypeStaff,
		ArchetypeManagement, ArchetypeIntegration, ArchetypeAdmin, ArchetypeOps,
	}
}

// PermissionsFor returns the granted permissions of an archetype.
func PermissionsFor(arch Archetype) []Permission {
	perms := permissionMatrix[arch]
	out := make([]Permission, 0, len(perms))
	for p := range perms {
		out = append(out, p)
	}
	return out
}
