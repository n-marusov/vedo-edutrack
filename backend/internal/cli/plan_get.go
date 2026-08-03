package cli

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"go.uber.org/zap"

	planrepo "vedo-edutrack/backend/internal/modules/planmanagement/adapters/repository"
	"vedo-edutrack/backend/internal/modules/planmanagement/adapters/repository/sqlc"
	"vedo-edutrack/backend/internal/platform/config"
	platformpostgres "vedo-edutrack/backend/internal/platform/postgres"
)

// newPlanGetCmd builds the `plan get` subcommand — reads a fixed learning plan.
func newPlanGetCmd() *cobra.Command {
	parent := &cobra.Command{
		Use:   "plan",
		Short: "Plan operations (read plan / progress)",
	}
	cmd := &cobra.Command{
		Use:   "get",
		Short: "Read a fixed learning plan with timeline and checkpoints",
		RunE: func(cmd *cobra.Command, _ []string) error {
			return runPlanGet(cmd)
		},
	}
	cmd.Flags().String("plan", "", "learning plan id")
	cmd.Flags().String("learner", "", "optional learner id guard")
	parent.AddCommand(cmd)
	return parent
}

func runPlanGet(cmd *cobra.Command) error {
	planID, _ := cmd.Flags().GetString("plan")
	learnerID, _ := cmd.Flags().GetString("learner")
	output, _ := cmd.Root().PersistentFlags().GetString("output")
	if planID == "" {
		return fmt.Errorf("--plan is required")
	}

	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	pool, err := platformpostgres.Connect(ctx, cfg.DatabaseURL)
	if err != nil {
		return fmt.Errorf("connect to database: %w", err)
	}
	defer platformpostgres.Close(pool)

	repository := planrepo.NewPlanRepository(pool)
	plan, steps, checkpoints, err := repository.GetPlan(ctx, planID)
	if err != nil {
		zapLogger.Error("plan get failed", zap.String("plan_id", planID), zap.Error(err))
		return fmt.Errorf("get plan: %w", err)
	}
	if learnerID != "" && plan.LearnerID != learnerID {
		return fmt.Errorf("plan %s belongs to learner %s, not %s", planID, plan.LearnerID, learnerID)
	}
	zapLogger.Info("plan loaded", zap.String("plan_id", planID), zap.String("learner_id", plan.LearnerID), zap.Int("steps", len(steps)), zap.Int("checkpoints", len(checkpoints)))

	response := planGetResponse{Plan: plan, Steps: steps, Checkpoints: checkpoints}
	switch strings.ToLower(output) {
	case "json":
		return json.NewEncoder(cmd.OutOrStdout()).Encode(response)
	case "table", "":
		printPlanTable(cmd, response)
		return nil
	default:
		return fmt.Errorf("unsupported output format %q (expected table|json)", output)
	}
}

type planGetResponse struct {
	Plan        sqlc.LearningPlanRow     `json:"plan"`
	Steps       []sqlc.PlanStepRow       `json:"steps"`
	Checkpoints []sqlc.PlanCheckpointRow `json:"checkpoints"`
}

func printPlanTable(cmd *cobra.Command, response planGetResponse) {
	p := response.Plan
	fmt.Fprintf(cmd.OutOrStdout(), "plan: %s\n", p.ID)
	fmt.Fprintf(cmd.OutOrStdout(), "learner: %s\n", p.LearnerID)
	fmt.Fprintf(cmd.OutOrStdout(), "goal: %s\n", p.GoalModuleID)
	fmt.Fprintf(cmd.OutOrStdout(), "ontology_version: %s\n", p.OntologyVersion)
	fmt.Fprintf(cmd.OutOrStdout(), "timeline: %s → %s\n", p.TimelineStart.Format(time.RFC3339), p.TimelineEnd.Format(time.RFC3339))
	fmt.Fprintln(cmd.OutOrStdout(), "steps:")
	for _, step := range response.Steps {
		fmt.Fprintf(cmd.OutOrStdout(), "  %03d  %s  %s..%s  horizon=%s essential=%t\n", step.Position, step.ModuleID, step.PlannedStart.Format("2006-01-02"), step.PlannedEnd.Format("2006-01-02"), step.Horizon, step.IsEssential)
	}
	fmt.Fprintln(cmd.OutOrStdout(), "checkpoints:")
	for _, checkpoint := range response.Checkpoints {
		fmt.Fprintf(cmd.OutOrStdout(), "  %s  %s  target=%.1f%%\n", checkpoint.CheckpointDate.Format("2006-01-02"), checkpoint.Name, checkpoint.TargetCoveragePercent)
	}
}
