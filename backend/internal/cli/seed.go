package cli

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/spf13/cobra"
	"go.uber.org/zap"

	"vedo-edutrack/backend/internal/platform/auth"
	"vedo-edutrack/backend/internal/platform/config"
	platformpostgres "vedo-edutrack/backend/internal/platform/postgres"
)

// newSeedCmd builds the `seed` subcommand — RBAC role catalog + demo data.
func newSeedCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "seed",
		Short: "Insert RBAC role catalog and demo data (idempotent)",
		RunE: func(cmd *cobra.Command, _ []string) error {
			integrationDemo, _ := cmd.Flags().GetBool("integration-demo")
			return runSeed(integrationDemo)
		},
	}
	cmd.Flags().Bool("integration-demo", false, "also create integration sandbox demo data (webhook subscriptions)")
	return cmd
}

// roleSeed is one role instance from the role catalog
// (REQ-NFR-security.compliance.role-catalog).
type roleSeed struct {
	name        string
	archetype   string
	description string
}

// roleCatalog is the initial EduTrack role catalog (10 instances, 7 archetypes).
var roleCatalog = []roleSeed{
	{name: "learner", archetype: string(auth.ArchetypeSelf), description: "Learner (self): own profile, route, plan, progress"},
	{name: "employee", archetype: string(auth.ArchetypeSelf), description: "Corporate employee (self)"},
	{name: "parent", archetype: string(auth.ArchetypeDependentsOwner), description: "Parent of a family-schooled learner"},
	{name: "teacher", archetype: string(auth.ArchetypeStaff), description: "Teacher / group manager"},
	{name: "methodologist", archetype: string(auth.ArchetypeManagement), description: "Methodologist: coverage, plans, reporting"},
	{name: "school-director", archetype: string(auth.ArchetypeManagement), description: "Private school director"},
	{name: "hr-manager", archetype: string(auth.ArchetypeManagement), description: "Enterprise HR / L&D lead"},
	{name: "platform-integrator", archetype: string(auth.ArchetypeIntegration), description: "EdTech platform integrator (API + webhooks)"},
	{name: "admin", archetype: string(auth.ArchetypeAdmin), description: "System administrator (full CRUD)"},
	{name: "ops", archetype: string(auth.ArchetypeOps), description: "Operations / infrastructure (JIT, outside product resources)"},
}

// runSeed inserts the RBAC role catalog and the permission matrix rows,
// then creates the default admin user. Idempotent (ON CONFLICT DO NOTHING).
// When integrationDemo is true, sandbox webhook subscriptions are created too.
func runSeed(integrationDemo bool) error {
	zapLogger.Info("seed started")

	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	pool, err := platformpostgres.Connect(ctx, cfg.DatabaseURL)
	if err != nil {
		return fmt.Errorf("connect to database: %w", err)
	}
	defer platformpostgres.Close(pool)

	tx, err := pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer func() { _ = tx.Rollback(ctx) }()

	// 1. Roles (ON CONFLICT DO NOTHING — idempotent re-runs).
	inserted, err := insertRoles(ctx, tx)
	if err != nil {
		return err
	}

	// 2. Permission matrix rows per role archetype.
	perms, err := insertPermissions(ctx, tx)
	if err != nil {
		return err
	}

	// 3. Default admin user with the admin role.
	if err := insertAdminUser(ctx, tx); err != nil {
		return err
	}

	// 4. Integration sandbox demo data (--integration-demo).
	if integrationDemo {
		if err := insertIntegrationDemoData(ctx, tx); err != nil {
			return err
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit seed tx: %w", err)
	}

	zapLogger.Info("seed completed",
		zap.Int("roles_inserted", inserted),
		zap.Int("permissions_inserted", perms),
	)
	return nil
}

// insertRoles upserts the role catalog. Returns the number of inserted rows.
func insertRoles(ctx context.Context, tx pgx.Tx) (int, error) {
	var inserted int
	for _, r := range roleCatalog {
		tag, err := tx.Exec(ctx,
			`INSERT INTO identity_access.roles (name, archetype, description)
			 VALUES ($1, $2, $3)
			 ON CONFLICT (name) DO NOTHING`,
			r.name, r.archetype, r.description,
		)
		if err != nil {
			return 0, fmt.Errorf("insert role %q: %w", r.name, err)
		}
		inserted += int(tag.RowsAffected())
	}
	zapLogger.Info("roles upserted", zap.Int("inserted", inserted))
	return inserted, nil
}

// insertPermissions materializes the permission matrix: for every role, the
// granted permissions of its archetype (deny-by-default — only grants written).
func insertPermissions(ctx context.Context, tx pgx.Tx) (int, error) {
	var inserted int
	for _, r := range roleCatalog {
		arch := auth.Archetype(r.archetype)
		perms := auth.PermissionsFor(arch)

		// Resolve the role id (inserted above or pre-existing).
		var roleID string
		err := tx.QueryRow(ctx,
			`SELECT id FROM identity_access.roles WHERE name = $1`, r.name,
		).Scan(&roleID)
		if err != nil {
			return 0, fmt.Errorf("resolve role id %q: %w", r.name, err)
		}

		for _, p := range perms {
			tag, err := tx.Exec(ctx,
				`INSERT INTO identity_access.role_permissions (role_id, permission, scope)
				 VALUES ($1, $2, $3)
				 ON CONFLICT (role_id, permission, scope) DO NOTHING`,
				roleID, string(p), defaultScope(arch),
			)
			if err != nil {
				return 0, fmt.Errorf("insert permission %q for %q: %w", p, r.name, err)
			}
			inserted += int(tag.RowsAffected())
		}
	}
	zapLogger.Info("permissions upserted", zap.Int("inserted", inserted))
	return inserted, nil
}

// defaultScope returns the archetype's default scope (role-catalog table).
func defaultScope(arch auth.Archetype) string {
	switch arch {
	case auth.ArchetypeSelf:
		return "own"
	case auth.ArchetypeDependentsOwner:
		return "dependents"
	case auth.ArchetypeStaff, auth.ArchetypeManagement:
		return "unit"
	default:
		return "all"
	}
}

// insertAdminUser creates the default admin user (admin@edutrack.local).
// Idempotent — the demo users table lives in the identity_access schema
// (created by the first migration).
func insertAdminUser(ctx context.Context, tx pgx.Tx) error {
	_, err := tx.Exec(ctx,
		`INSERT INTO identity_access.users (email, display_name)
		 VALUES ('admin@edutrack.local', 'System Administrator')
		 ON CONFLICT (email) DO NOTHING`)
	if err != nil {
		return fmt.Errorf("insert admin user: %w", err)
	}
	zapLogger.Info("default admin user ensured", zap.String("email", "admin@edutrack.local"))
	return nil
}

// insertIntegrationDemoData creates sandbox webhook subscriptions for the
// integration sandbox (--integration-demo). The demo uses the admin user's
// tenant and a local webhook receiver URL; idempotent by tenant+url.
func insertIntegrationDemoData(ctx context.Context, tx pgx.Tx) error {
	// Demo subscriptions keyed on (tenant_id, url): admin demo receiver.
	demoSubs := []struct {
		tenantID   string
		url        string
		eventTypes []string
	}{
		{tenantID: "admin@edutrack.local", url: "http://localhost:9099/hooks/edutrack", eventTypes: []string{"module.mastered", "plan.deviated", "route.recalculated"}},
		{tenantID: "admin@edutrack.local", url: "http://localhost:9099/hooks/edutrack-sparql", eventTypes: []string{"standard.risk_detected"}},
	}

	// Secret: fixed demo secret (32+ chars) — sandbox only, not production.
	const demoSecret = "vedo-edutrack-integration-demo-secret-0123456789" // #nosec G101

	for _, s := range demoSubs {
		var exists bool
		// Check the subscription doesn't already exist for this tenant+url.
		err := tx.QueryRow(ctx,
			`SELECT EXISTS(SELECT 1 FROM integrations.webhook_subscriptions WHERE tenant_id = $1 AND url = $2)`,
			s.tenantID, s.url,
		).Scan(&exists)
		if err != nil {
			return fmt.Errorf("check demo subscription %s: %w", s.url, err)
		}
		if exists {
			continue
		}

		_, err = tx.Exec(ctx,
			`INSERT INTO integrations.webhook_subscriptions (tenant_id, url, event_types, signing_secret)
			 VALUES ($1, $2, $3, $4)`,
			s.tenantID, s.url, s.eventTypes, demoSecret,
		)
		if err != nil {
			return fmt.Errorf("insert demo subscription %s: %w", s.url, err)
		}
		zapLogger.Info("integration demo subscription created", zap.String("tenant", s.tenantID), zap.String("url", s.url), zap.Strings("events", s.eventTypes))
	}
	zapLogger.Info("integration demo data ready (webhook subscriptions)")
	return nil
}
