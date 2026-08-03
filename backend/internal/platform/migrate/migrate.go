// Package migrate provides the embedded SQL migration runner for the
// vedo-edutrack binary.
//
// Design (ADR-DES.API.cli-interface, ADR-IMPL.PROCESS.development-tooling §3):
// migrations are managed as Atlas-style SQL files under backend/migrations/
// and embedded into the binary via go:embed (distroless runtime — no shell,
// no filesystem). The runner applies them against PostgreSQL in lexical
// order, tracking applied versions in a schema_migrations table.
//
// Subcommands:
//   - Up:      apply all pending migrations, each inside its own transaction
//   - Down:    revert the last applied migration (executes its `-- down:`
//     section when defined; migrations without a down section are
//     non-reversible and reported as such)
//   - Validate: drift check — applied vs embedded versions and checksum
//     mismatches (drift = 0 is the target per REQ-NFR-ops.release.deployment-verification)
package migrate

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"io/fs"
	"regexp"
	"sort"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
)

// migrationFile is one parsed migration: a lexical version prefix, the
// filename, the up SQL (everything before the `-- down:` marker) and the
// optional down SQL (everything after it).
type migrationFile struct {
	version string // "000002"
	name    string // "000002_planmanagement_init.sql"
	upSQL   string
	downSQL string // "" when no `-- down:` marker is present
	sha256  string // base64(sha256(file content)) — matches atlas h1: content
}

// downMarker is the comment line separating the up section from the optional
// down section inside a single migration file.
const downMarker = "-- down:"

// versionRe matches migration filenames: NNNNNN_<name>.sql
var versionRe = regexp.MustCompile(`^(\d{6})_.+\.sql$`)

// Runner executes embedded migrations against PostgreSQL.
type Runner struct {
	fsys fs.FS
	pool *pgxpool.Pool
}

// NewRunner builds a Runner over the embedded migration filesystem and an
// existing connection pool.
func NewRunner(fsys fs.FS, pool *pgxpool.Pool) *Runner {
	return &Runner{fsys: fsys, pool: pool}
}

// loadMigrations reads all *.sql files from the embedded FS, parses them and
// sorts them lexically by version.
func (r *Runner) loadMigrations() ([]migrationFile, error) {
	entries, err := fs.ReadDir(r.fsys, ".")
	if err != nil {
		return nil, fmt.Errorf("read embedded migrations dir: %w", err)
	}

	files := make([]migrationFile, 0, len(entries))
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".sql") {
			continue
		}
		m := versionRe.FindStringSubmatch(e.Name())
		if m == nil {
			return nil, fmt.Errorf("migration %q does not match NNNNNN_<name>.sql naming convention", e.Name())
		}

		content, err := fs.ReadFile(r.fsys, e.Name())
		if err != nil {
			return nil, fmt.Errorf("read migration %q: %w", e.Name(), err)
		}

		up, down := splitDownSection(string(content))
		sum := sha256.Sum256(content)

		files = append(files, migrationFile{
			version: m[1],
			name:    e.Name(),
			upSQL:   up,
			downSQL: down,
			sha256:  base64.StdEncoding.EncodeToString(sum[:]),
		})
	}

	sort.Slice(files, func(i, j int) bool { return files[i].version < files[j].version })
	return files, nil
}

// splitDownSection splits a migration file into up and optional down SQL.
func splitDownSection(content string) (up, down string) {
	// Find the `-- down:` marker on its own line. Everything before it is the
	// up section; everything after is the down section.
	lines := strings.Split(content, "\n")
	for i, line := range lines {
		if strings.HasPrefix(strings.TrimSpace(line), downMarker) {
			return strings.Join(lines[:i], "\n"), strings.Join(lines[i+1:], "\n")
		}
	}
	return content, ""
}

// ensureTrackingTable creates the schema_migrations tracking table.
func ensureTrackingTable(ctx context.Context, pool *pgxpool.Pool) error {
	_, err := pool.Exec(ctx, `
		CREATE TABLE IF NOT EXISTS schema_migrations (
			version    text PRIMARY KEY,
			filename   text NOT NULL,
			checksum   text NOT NULL,
			applied_at timestamptz NOT NULL DEFAULT now()
		)`)
	if err != nil {
		return fmt.Errorf("ensure schema_migrations table: %w", err)
	}
	return nil
}

// appliedVersions returns the already-applied migrations keyed by version.
func appliedVersions(ctx context.Context, pool *pgxpool.Pool) (map[string]migrationFile, error) {
	rows, err := pool.Query(ctx,
		`SELECT version, filename, checksum FROM schema_migrations ORDER BY version`)
	if err != nil {
		return nil, fmt.Errorf("query schema_migrations: %w", err)
	}
	defer rows.Close()

	applied := map[string]migrationFile{}
	for rows.Next() {
		var m migrationFile
		if err := rows.Scan(&m.version, &m.name, &m.sha256); err != nil {
			return nil, fmt.Errorf("scan schema_migrations row: %w", err)
		}
		applied[m.version] = m
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate schema_migrations rows: %w", err)
	}
	return applied, nil
}

// Up applies all pending migrations in lexical order. Each migration runs
// inside its own transaction; a failed migration aborts its transaction and
// returns an error (previously applied migrations stay applied).
func (r *Runner) Up(ctx context.Context) (int, error) {
	if err := ensureTrackingTable(ctx, r.pool); err != nil {
		return 0, err
	}

	files, err := r.loadMigrations()
	if err != nil {
		return 0, err
	}
	applied, err := appliedVersions(ctx, r.pool)
	if err != nil {
		return 0, err
	}

	appliedCount := 0
	for _, f := range files {
		if _, ok := applied[f.version]; ok {
			continue // already applied
		}
		if err := r.applyOne(ctx, f); err != nil {
			return appliedCount, err
		}
		appliedCount++
	}
	return appliedCount, nil
}

// applyOne executes a single migration's up SQL in a transaction and records
// it in schema_migrations.
func (r *Runner) applyOne(ctx context.Context, f migrationFile) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin migration tx for %s: %w", f.name, err)
	}
	defer func() { _ = tx.Rollback(ctx) }()

	if _, err := tx.Exec(ctx, f.upSQL); err != nil {
		return fmt.Errorf("apply migration %s: %w", f.name, err)
	}
	if _, err := tx.Exec(ctx,
		`INSERT INTO schema_migrations (version, filename, checksum) VALUES ($1, $2, $3)`,
		f.version, f.name, f.sha256,
	); err != nil {
		return fmt.Errorf("record migration %s: %w", f.name, err)
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit migration %s: %w", f.name, err)
	}
	return nil
}

// Down reverts the last applied migration. When the migration defines a
// `-- down:` section it is executed; otherwise the migration is reported as
// non-reversible (no state change, error returned so the operator sees it).
func (r *Runner) Down(ctx context.Context) (string, error) {
	if err := ensureTrackingTable(ctx, r.pool); err != nil {
		return "", err
	}

	applied, err := appliedVersions(ctx, r.pool)
	if err != nil {
		return "", err
	}
	lastVersion, ok := lastAppliedVersion(applied)
	if !ok {
		return "", errors.New("no applied migrations to revert")
	}
	lastApplied := applied[lastVersion]

	// Resolve the full file (need downSQL — not stored in tracking table).
	files, err := r.loadMigrations()
	if err != nil {
		return "", err
	}
	f, err := findMigration(files, lastApplied.version)
	if err != nil {
		return "", fmt.Errorf("applied migration %s not found in embedded set (was it deleted?)", lastApplied.name)
	}
	if strings.TrimSpace(f.downSQL) == "" {
		return f.name, fmt.Errorf("migration %s has no -- down: section and cannot be reverted", f.name)
	}

	return f.name, r.revertOne(ctx, *f)
}

// lastAppliedVersion returns the highest applied version (or ok=false when
// nothing has been applied yet).
func lastAppliedVersion(applied map[string]migrationFile) (string, bool) {
	last := ""
	for v := range applied {
		if v > last {
			last = v
		}
	}
	if last == "" {
		return "", false
	}
	return last, true
}

// findMigration returns the migration with the given version from the
// embedded set.
func findMigration(files []migrationFile, version string) (*migrationFile, error) {
	for i := range files {
		if files[i].version == version {
			return &files[i], nil
		}
	}
	return nil, errors.New("not found")
}

// revertOne executes a migration's down SQL and unrecords it in a single
// transaction.
func (r *Runner) revertOne(ctx context.Context, f migrationFile) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin revert tx for %s: %w", f.name, err)
	}
	defer func() { _ = tx.Rollback(ctx) }()

	if _, err := tx.Exec(ctx, f.downSQL); err != nil {
		return fmt.Errorf("revert migration %s: %w", f.name, err)
	}
	if _, err := tx.Exec(ctx,
		`DELETE FROM schema_migrations WHERE version = $1`, f.version); err != nil {
		return fmt.Errorf("unrecord migration %s: %w", f.name, err)
	}
	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit revert of %s: %w", f.name, err)
	}
	return nil
}

// Validate performs a drift check between the embedded migration set and the
// applied state. Returns a list of drift descriptions (empty = drift 0).
func (r *Runner) Validate(ctx context.Context) ([]string, error) {
	if err := ensureTrackingTable(ctx, r.pool); err != nil {
		return nil, err
	}

	files, err := r.loadMigrations()
	if err != nil {
		return nil, err
	}
	applied, err := appliedVersions(ctx, r.pool)
	if err != nil {
		return nil, err
	}

	var drift []string
	embedded := map[string]migrationFile{}
	for _, f := range files {
		embedded[f.version] = f
		if a, ok := applied[f.version]; ok {
			if a.sha256 != f.sha256 {
				drift = append(drift, fmt.Sprintf(
					"checksum mismatch for %s: applied=%s embedded=%s", f.name, a.sha256, f.sha256))
			}
		} else {
			drift = append(drift, fmt.Sprintf("pending migration not applied: %s", f.name))
		}
	}
	for v, a := range applied {
		if _, ok := embedded[v]; !ok {
			drift = append(drift, fmt.Sprintf(
				"applied migration missing from embedded set: %s", a.name))
		}
	}
	sort.Strings(drift)
	return drift, nil
}
