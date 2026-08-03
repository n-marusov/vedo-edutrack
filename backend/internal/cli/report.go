package cli

import (
	"fmt"

	"github.com/spf13/cobra"

	gapapp "vedo-edutrack/backend/internal/modules/gapcoverage/application"
	gapdomain "vedo-edutrack/backend/internal/modules/gapcoverage/domain"
)

// newReportCmd builds the `report` subcommand — attestation / FGOS coverage
// reports to stdout or JSON.
func newReportCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "report",
		Short: "Generate attestation / FGOS coverage reports",
		RunE: func(cmd *cobra.Command, _ []string) error {
			return runReport(cmd)
		},
	}
	cmd.Flags().String("input", "", "JSON input file with bindings and mastery")
	cmd.Flags().String("learner", "", "learner id for audit logging")
	cmd.Flags().String("type", "coverage", "report type: coverage|attestation")
	cmd.Flags().Float64("threshold", 80, "attestation readiness coverage threshold percent")
	return cmd
}

func runReport(cmd *cobra.Command) error {
	inputPath, _ := cmd.Flags().GetString("input")
	reportType, _ := cmd.Flags().GetString("type")
	threshold, _ := cmd.Flags().GetFloat64("threshold")
	output, _ := cmd.Root().PersistentFlags().GetString("output")
	if inputPath == "" {
		return fmt.Errorf("--input is required until DB-backed reports land")
	}
	var input reportInput
	if err := readJSONFile(inputPath, &input); err != nil {
		return err
	}
	service := gapapp.NewCoverageService(zapLogger)
	switch reportType {
	case "coverage":
		report := service.Coverage(input.Bindings, input.Mastery)
		return writeCLIResult(cmd, output, report, func() { printCoverage(cmd, report) })
	case "attestation":
		report := service.Attestation(input.Bindings, input.Mastery, threshold)
		return writeCLIResult(cmd, output, report, func() {
			printCoverage(cmd, report.Coverage)
			fmt.Fprintf(cmd.OutOrStdout(), "ready: %t\n", report.Ready)
		})
	default:
		return fmt.Errorf("unsupported report type %q (expected coverage|attestation)", reportType)
	}
}

type reportInput struct {
	Bindings []gapdomain.FgosBinding `json:"bindings"`
	Mastery  gapdomain.Mastery       `json:"mastery"`
}

func printCoverage(cmd *cobra.Command, report gapdomain.CoverageReport) {
	fmt.Fprintf(cmd.OutOrStdout(), "coverage: %.1f%% (%d/%d)\n", report.Percent, report.Covered, report.Total)
	fmt.Fprintln(cmd.OutOrStdout(), "deficits:")
	for _, deficit := range report.Deficits {
		fmt.Fprintf(cmd.OutOrStdout(), "  %s blocked_by=%s\n", deficit.RequirementID, deficit.BlockingModuleID)
	}
}
