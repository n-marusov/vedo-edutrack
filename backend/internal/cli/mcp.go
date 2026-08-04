package cli

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"
	"go.uber.org/zap"

	execapp "vedo-edutrack/backend/internal/modules/executionprogress/application"
	execdomain "vedo-edutrack/backend/internal/modules/executionprogress/domain"
	gapapp "vedo-edutrack/backend/internal/modules/gapcoverage/application"
	gapdomain "vedo-edutrack/backend/internal/modules/gapcoverage/domain"
	"vedo-edutrack/backend/internal/modules/integrations/adapters/mcp"
	ontostub "vedo-edutrack/backend/internal/modules/ontologyport/adapters/stub"
	resourceapp "vedo-edutrack/backend/internal/modules/resources/application"
	resourcedomain "vedo-edutrack/backend/internal/modules/resources/domain"
	routestub "vedo-edutrack/backend/internal/modules/routeplanning/adapters/stub"
)

// newMcpCmd builds the `mcp` subcommand — MCP server over stdio (F6.6).
// The server speaks JSON-RPC 2.0 on stdin/stdout; all logs go to stderr.
func newMcpCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "mcp",
		Short: "MCP server over stdio for AI agents (F6.6)",
		RunE: func(_ *cobra.Command, _ []string) error {
			return runMcpServer()
		},
	}
}

// mcpProgressRepo is an in-memory ProgressRepository for MCP forecast tools.
type mcpProgressRepo struct{}

func (mcpProgressRepo) GetProgress(_ context.Context, _ string) ([]execdomain.ModuleProgress, error) {
	return []execdomain.ModuleProgress{
		{ModuleID: "percent", Status: execdomain.StatusMastered, MasteredAt: mcpTimePtr(time.Date(2026, 9, 5, 0, 0, 0, 0, time.UTC))},
		{ModuleID: "solutions", Status: execdomain.StatusInProgress},
	}, nil
}
func (mcpProgressRepo) GetPlannedModuleCount(_ context.Context, _, _ string) (int, error) {
	return 10, nil
}
func (mcpProgressRepo) GetRemainingModules(_ context.Context, _, _ string) (int, error) {
	return 5, nil
}

// mcpTools builds the read-oriented MCP tool registry over the stub
// application services (same fixtures as the HTTP API handler).
func mcpTools() []mcp.Tool {
	graph := ontostub.NewGraph()
	computer := routestub.NewComputer(graph)
	progress := execapp.NewProgressService(zapLogger)
	forecast := execapp.NewForecastService(mcpProgressRepo{}, zapLogger)
	gaps := gapapp.NewGapService(zapLogger)
	coverage := gapapp.NewCoverageService(zapLogger)
	catalog, _ := resourcedomain.NewCatalog([]resourcedomain.Resource{
		{ID: "res-1", Title: "Percent video", Type: resourcedomain.ResourceTypeContent, Format: "video", Source: "school", Difficulty: "basic", DurationMinutes: 10, URI: "https://example.test/res-1"},
		{ID: "res-2", Title: "Percent text", Type: resourcedomain.ResourceTypeContent, Format: "text", Source: "school", Difficulty: "basic", DurationMinutes: 30, URI: "https://example.test/res-2"},
	})
	_ = catalog.BindToModule(resourcedomain.ResourceBinding{ResourceID: "res-1", ModuleID: "math-5-11", LinkType: resourcedomain.LinkAppliesTo})
	resources := resourceapp.NewCatalogService(catalog, zapLogger)

	return []mcp.Tool{
		{
			Name:        "get_route",
			Description: "Returns the current learning route for a learner (horizons, steps).",
			InputSchema: map[string]any{
				"type": "object",
				"properties": map[string]any{
					"learner_id":    map[string]any{"type": "string"},
					"goal_topic_id": map[string]any{"type": "string"},
				},
				"required": []string{"learner_id", "goal_topic_id"},
			},
			Handler: func(_ context.Context, args map[string]any) (any, error) {
				learnerID, _ := args["learner_id"].(string)
				goal, _ := args["goal_topic_id"].(string)
				if learnerID == "" || goal == "" {
					return nil, errors.New("learner_id and goal_topic_id are required")
				}
				route, err := computer.ComputeRoute(learnerID, goal)
				if err != nil {
					return nil, err
				}
				zapLogger.Info("mcp get_route", zap.String("learner_id", learnerID), zap.String("goal", goal), zap.Int("steps", len(route)))
				return map[string]any{"learner_id": learnerID, "goal_topic_id": goal, "steps": route}, nil
			},
		},
		{
			Name:        "get_progress",
			Description: "Returns plan-vs-actual progress and forecast for a learner.",
			InputSchema: map[string]any{
				"type": "object",
				"properties": map[string]any{
					"learner_id": map[string]any{"type": "string"},
				},
				"required": []string{"learner_id"},
			},
			Handler: func(ctx context.Context, args map[string]any) (any, error) {
				learnerID, _ := args["learner_id"].(string)
				if learnerID == "" {
					return nil, errors.New("learner_id is required")
				}
				plan := execdomain.FixedPlan{ID: "plan-1", LearnerID: learnerID, Modules: []execdomain.PlannedModule{
					{ModuleID: "percent", PlannedStart: mcpDay(0), PlannedEnd: mcpDay(10)},
				}}
				actual := []execdomain.ModuleProgress{{ModuleID: "percent", Status: execdomain.StatusMastered, MasteredAt: mcpTimePtr(mcpDay(5))}}
				comparison := progress.Compare(plan, actual)
				fc, err := forecast.ForecastReadiness(ctx, learnerID, "plan-1", 5, 6)
				if err != nil {
					zapLogger.Warn("forecast failed", zap.Error(err))
				}
				return map[string]any{
					"learner_id":  learnerID,
					"deviations":  comparison.Deviations,
					"divergences": comparison.Divergences,
					"forecast":    fc,
				}, nil
			},
		},
		{
			Name:        "get_coverage",
			Description: "Returns FGOS/profstandard coverage for a learner.",
			InputSchema: map[string]any{
				"type": "object",
				"properties": map[string]any{
					"learner_id": map[string]any{"type": "string"},
				},
				"required": []string{"learner_id"},
			},
			Handler: func(_ context.Context, args map[string]any) (any, error) {
				learnerID, _ := args["learner_id"].(string)
				if learnerID == "" {
					return nil, errors.New("learner_id is required")
				}
				report := coverage.Coverage(
					[]gapdomain.FgosBinding{{ModuleID: "percent", RequirementID: "fgos-math-5"}, {ModuleID: "solutions", RequirementID: "fgos-chem-8"}},
					gapdomain.Mastery{Modules: map[string]float64{"percent": 1.0}},
				)
				return map[string]any{"learner_id": learnerID, "coverage": report}, nil
			},
		},
		{
			Name:        "get_gaps",
			Description: "Returns diagnosed gaps with root causes for a learner lag module.",
			InputSchema: map[string]any{
				"type": "object",
				"properties": map[string]any{
					"learner_id":    map[string]any{"type": "string"},
					"lag_module_id": map[string]any{"type": "string"},
				},
				"required": []string{"learner_id", "lag_module_id"},
			},
			Handler: func(_ context.Context, args map[string]any) (any, error) {
				learnerID, _ := args["learner_id"].(string)
				lag, _ := args["lag_module_id"].(string)
				if learnerID == "" || lag == "" {
					return nil, errors.New("learner_id and lag_module_id are required")
				}
				gapGraph := gapdomain.Graph{Modules: []gapdomain.Module{
					{ID: "percent", Title: "Проценты"}, {ID: "solutions", Title: "Растворы"}, {ID: "chemistry", Title: "Химия"},
				}, Links: []gapdomain.Link{
					{SourceID: "percent", TargetID: "solutions", Type: gapdomain.LinkStrictPrerequisite},
					{SourceID: "solutions", TargetID: "chemistry", Type: gapdomain.LinkStrictPrerequisite},
				}}
				mastery := gapdomain.Mastery{Modules: map[string]float64{"percent": 0.7, "solutions": 0.2}}
				result := gaps.Diagnose(gapGraph, mastery, lag)
				return map[string]any{"learner_id": learnerID, "lag_module_id": lag, "diagnosis": result}, nil
			},
		},
		{
			Name:        "get_resources",
			Description: "Returns resources for a module.",
			InputSchema: map[string]any{
				"type": "object",
				"properties": map[string]any{
					"module_id": map[string]any{"type": "string"},
				},
				"required": []string{"module_id"},
			},
			Handler: func(_ context.Context, args map[string]any) (any, error) {
				moduleID, _ := args["module_id"].(string)
				if moduleID == "" {
					return nil, errors.New("module_id is required")
				}
				result := resources.Query(resourceapp.CatalogQuery{ModuleID: moduleID})
				return map[string]any{"module_id": moduleID, "items": result.Items, "total": result.Total}, nil
			},
		},
	}
}

// runMcpServer runs the MCP server over stdio until stdin closes.
func runMcpServer() error {
	zapLogger.Info("mcp server starting (stdio)")

	server := mcp.New(mcp.Config{
		APIKey: os.Getenv("VEDO_MCP_API_KEY"),
		Tools:  mcpTools(),
		OnShutdown: func() {
			zapLogger.Info("mcp server shutting down")
		},
	}, os.Stdin, os.Stdout, zapLogger)

	if err := server.Serve(context.Background()); err != nil {
		return fmt.Errorf("mcp server: %w", err)
	}
	return nil
}

// mcpDay returns a UTC timestamp n days from the epoch.
func mcpDay(offset int) time.Time {
	return time.Date(2026, 9, 1+offset, 0, 0, 0, 0, time.UTC)
}

// mcpTimePtr returns a pointer to a time.
func mcpTimePtr(t time.Time) *time.Time { return &t }
