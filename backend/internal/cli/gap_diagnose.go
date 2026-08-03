package cli

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"go.uber.org/zap"

	gapapp "vedo-edutrack/backend/internal/modules/gapcoverage/application"
	gapdomain "vedo-edutrack/backend/internal/modules/gapcoverage/domain"
)

// CLI output format constants (shared across commands).
const (
	outputJSON  = "json"
	outputTable = "table"
)

// newGapDiagnoseCmd builds the `gap diagnose` subcommand — root-cause
// diagnosis of learner lag.
func newGapDiagnoseCmd() *cobra.Command {
	parent := &cobra.Command{
		Use:   "gap",
		Short: "Gap operations (diagnose learner lag)",
	}
	cmd := &cobra.Command{
		Use:   "diagnose",
		Short: "Diagnose the root cause of a learner's learning gap",
		RunE: func(cmd *cobra.Command, _ []string) error {
			return runGapDiagnose(cmd)
		},
	}
	cmd.Flags().String("input", "", "JSON input file with graph, mastery, and lag_module_id")
	cmd.Flags().String("learner", "", "learner id for audit logging")
	cmd.Flags().String("plan", "", "optional plan id for future DB-backed diagnosis")
	parent.AddCommand(cmd)
	return parent
}

func runGapDiagnose(cmd *cobra.Command) error {
	inputPath, _ := cmd.Flags().GetString("input")
	learnerID, _ := cmd.Flags().GetString("learner")
	output, _ := cmd.Root().PersistentFlags().GetString("output")
	if inputPath == "" {
		return fmt.Errorf("--input is required until DB-backed gap diagnosis lands")
	}
	var input gapDiagnoseInput
	if err := readJSONFile(inputPath, &input); err != nil {
		return err
	}
	service := gapapp.NewGapService(zapLogger)
	result := service.Diagnose(input.Graph, input.Mastery, input.LagModuleID)
	if learnerID != "" {
		zapLogger.Info("gap diagnosis completed", zap.String("learner_id", learnerID), zap.Int("root_causes", len(result.RootCauses)))
	}
	return writeCLIResult(cmd, output, result, func() {
		fmt.Fprintf(cmd.OutOrStdout(), "status: %s\n", result.Status)
		for i, cause := range result.RootCauses {
			fmt.Fprintf(cmd.OutOrStdout(), "%2d  %s  mastery=%.2f blocked=%d\n", i+1, cause.ModuleID, cause.Mastery, cause.BlockedModules)
		}
	})
}

type gapDiagnoseInput struct {
	Graph       gapdomain.Graph   `json:"graph"`
	Mastery     gapdomain.Mastery `json:"mastery"`
	LagModuleID string            `json:"lag_module_id"`
}

func readJSONFile(path string, out any) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("read %s: %w", path, err)
	}
	if err := json.Unmarshal(data, out); err != nil {
		return fmt.Errorf("decode %s: %w", path, err)
	}
	return nil
}

func writeCLIResult(cmd *cobra.Command, output string, value any, printTable func()) error {
	switch strings.ToLower(output) {
	case outputJSON:
		return json.NewEncoder(cmd.OutOrStdout()).Encode(value)
	case outputTable, "":
		printTable()
		return nil
	default:
		return fmt.Errorf("unsupported output format %q (expected table|json)", output)
	}
}
