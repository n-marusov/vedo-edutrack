package cli

import (
	"github.com/spf13/cobra"
)

// newReportCmd builds the `report` subcommand — attestation / FGOS coverage
// reports to file (batch). Stub in M0.3.
func newReportCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "report",
		Short: "Generate attestation / FGOS coverage reports to file (batch)",
		RunE:  stubNotImplemented("report"),
	}
}
