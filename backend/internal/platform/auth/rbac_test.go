package auth

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"go.uber.org/zap"
)

func TestHasPermissionMatrix(t *testing.T) {
	cases := []struct {
		name  string
		roles []string
		perm  Permission
		want  bool
	}{
		{"learner can read route", []string{"learner"}, PermRouteCompute, true},
		{"learner cannot manage plan", []string{"learner"}, PermPlanManage, false},
		{"parent can read gap", []string{"parent"}, PermGapRead, true},
		{"teacher can manage plan", []string{"teacher"}, PermPlanManage, true},
		{"teacher cannot configure webhooks", []string{"teacher"}, PermWebhookConfig, false},
		{"methodologist can manage resources", []string{"methodologist"}, PermResourceMng, true},
		{"integration can configure webhooks", []string{"platform-integrator"}, PermWebhookConfig, true},
		{"integration cannot manage users", []string{"platform-integrator"}, PermUserManage, false},
		{"admin full access", []string{"admin"}, PermUserManage, true},
		{"admin can manage plan", []string{"admin"}, PermPlanManage, true},
		{"ops has no product permissions", []string{"ops"}, PermRouteCompute, false},
		{"ops cannot manage users", []string{"ops"}, PermUserManage, false},
		{"unknown role denied", []string{"ghost"}, PermRouteCompute, false},
		{"deny by default", []string{"learner"}, PermRouteWrite, false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := HasPermission(tc.roles, tc.perm); got != tc.want {
				t.Errorf("HasPermission(%v, %s) = %v, want %v", tc.roles, tc.perm, got, tc.want)
			}
		})
	}
}

// TestRequirePermissionMiddleware verifies the chi middleware enforces the
// permission: denied → 403, granted → next handler runs.
func TestRequirePermissionMiddleware(t *testing.T) {
	logger := zap.NewNop()
	mw := RequirePermission(PermUserManage, logger)

	next := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	// Denied: learner without user:manage.
	req := httptest.NewRequest(http.MethodGet, "/admin", nil)
	req = req.WithContext(WithClaims(req.Context(), &Claims{UserID: "learner1", Roles: []string{"learner"}}))
	rec := httptest.NewRecorder()
	mw(next).ServeHTTP(rec, req)
	if rec.Code != http.StatusForbidden {
		t.Errorf("denied status = %d, want 403", rec.Code)
	}

	// Allowed: admin with user:manage.
	req2 := httptest.NewRequest(http.MethodGet, "/admin", nil)
	req2 = req2.WithContext(WithClaims(req2.Context(), &Claims{UserID: "admin1", Roles: []string{"admin"}}))
	rec2 := httptest.NewRecorder()
	mw(next).ServeHTTP(rec2, req2)
	if rec2.Code != http.StatusOK {
		t.Errorf("allowed status = %d, want 200", rec2.Code)
	}
}
