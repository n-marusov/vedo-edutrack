// Package migrations embeds the SQL migration files into the binary.
//
// The embedded filesystem is the single source of truth for schema DDL.
// The vedo-edutrack binary is self-contained (distroless runtime has no
// shell and no local filesystem for migrations), so migrations are compiled
// in via go:embed and executed by the platform/migrate runner.
//
// Migration files follow the Atlas file convention:
//
//	<NNNN>_<schema>_<desc>.sql
//
// and support an optional `-- down:` marker section for rollback statements.
// Checksums (h1: = base64(sha256)) are tracked in atlas.sum and mirrored in
// the schema_migrations tracking table by the runner.
package migrations

import "embed"

//go:embed *.sql
var FS embed.FS
