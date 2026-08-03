package cli

import (
	"github.com/spf13/cobra"
)

// newPlanGetCmd builds the `plan get` subcommand — reads a plan / progress.
// Stub in M0.3.
func newPlanGetCmd() *cobra.Command {
	parent := &cobra.Command{
		Use:   "plan",
		Short: "Plan operations (read plan / progress)",
	}
	parent.AddCommand(&cobra.Command{
		Use:   "get",
		Short: "Read a learning plan and its progress",
		RunE:  stubNotImplemented("plan get"),
	})
	return parent
}
