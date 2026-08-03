package cli

import (
	"github.com/spf13/cobra"
)

// newGapDiagnoseCmd builds the `gap diagnose` subcommand — root-cause
// diagnosis of learner lag. Stub in M0.3.
func newGapDiagnoseCmd() *cobra.Command {
	parent := &cobra.Command{
		Use:   "gap",
		Short: "Gap operations (diagnose learner lag)",
	}
	parent.AddCommand(&cobra.Command{
		Use:   "diagnose",
		Short: "Diagnose the root cause of a learner's learning gap",
		RunE:  stubNotImplemented("gap diagnose"),
	})
	return parent
}
