// Package spa embeds the built frontend SPA into the backend binary.
//
// The frontend build output (frontend/dist) is copied into
// backend/internal/platform/spa/frontend/dist before `go build`
// (make build-frontend), so a single binary serves both the API and the SPA
// from one port — the M0.3 embed unification (see deploy/README.md).
package spa

import (
	"embed"
	"io/fs"
)

//go:embed all:frontend/dist
var embeddedFS embed.FS

// distFS is the embedded frontend/dist subtree.
var distFS, _ = fs.Sub(embeddedFS, "frontend/dist")

// EmbeddedFS returns the embedded SPA filesystem, or nil when the frontend
// was not built (dev mode without `make build-frontend`).
func EmbeddedFS() (fs.FS, bool) {
	if distFS == nil {
		return nil, false
	}
	// Probe for index.html to distinguish "no embed" from an empty embed.
	if _, err := fs.Stat(distFS, "index.html"); err != nil {
		return nil, false
	}
	return distFS, true
}

// FileCount returns the number of embedded files (for startup logging).
func FileCount() int {
	if distFS == nil {
		return 0
	}
	count := 0
	_ = fs.WalkDir(distFS, ".", func(_ string, d fs.DirEntry, err error) error {
		if err == nil && !d.IsDir() {
			count++
		}
		return nil
	})
	return count
}
