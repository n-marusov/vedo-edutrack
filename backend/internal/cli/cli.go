// Package cli implements the cobra command tree of the vedo-edutrack binary.
//
// The CLI is an input adapter over the Application layer — the same pattern
// as the HTTP handler and the MCP server (ADR-DES.API.cli-interface): every
// command calls module use cases through wire providers; there is no second
// path to data.
//
// Commands:
//
//	server         long-running process (HTTP API + SPARQL + webhooks + MCP SSE)
//	mcp            MCP server over stdio for AI agents (F6.6)
//	migrate        Atlas migrations up/down/validate (drift = 0)
//	seed           RBAC role catalog + demo data
//	ontology sync  copy subgraph from VEDO Hub (F0.2)
//	route compute  compute a route (--stub | from DB) — dev/test tool
//	plan get       read plan / progress
//	gap diagnose   root-cause diagnosis of learner lag
//	report         attestation / FGOS coverage reports to file (batch)
//
// Each command builds its own minimal wire graph (per-command lazy wire).
// Commands are scriptable (no interactive prompts) and emit structured
// audit logs via zap.
package cli

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
	"go.uber.org/zap"

	"vedo-edutrack/backend/internal/platform/config"
	platformlogger "vedo-edutrack/backend/internal/platform/logger"
)

// package-level logger set by the persistent pre-run of the root command.
var zapLogger = platformlogger.NewNop()

// Execute runs the root command. Returns an error on failure.
func Execute() error {
	root := NewRoot()
	return root.Execute()
}

// NewRoot builds the root command with all subcommands registered.
func NewRoot() *cobra.Command {
	root := &cobra.Command{
		Use:           "vedo-edutrack",
		Short:         "VEDO EduTrack — educational trajectory service",
		Long:          "VEDO EduTrack reads knowledge ontologies via the VEDO Hub API and adds educational mechanics: trajectory generation, learning plans, gap diagnosis, FGOS coverage, and knowledge-graph visualization.",
		SilenceUsage:  true,
		SilenceErrors: true,
		PersistentPreRunE: func(cmd *cobra.Command, _ []string) error {
			return initLogger()
		},
	}

	// Persistent flags shared by all subcommands.
	root.PersistentFlags().String("output", "table", "output format: table|json|csv")
	root.PersistentFlags().Bool("yes", false, "skip confirmation prompts")

	root.AddCommand(
		newServerCmd(),
		newMcpCmd(),
		newMigrateCmd(),
		newSeedCmd(),
		newOntologySyncCmd(),
		newRouteComputeCmd(),
		newPlanGetCmd(),
		newGapDiagnoseCmd(),
		newReportCmd(),
		newVersionCmd(),
		newHealthCmd(),
	)

	return root
}

// initLogger creates the zap logger (LOG_LEVEL from env) and bridges the
// standard log package to zap so third-party libraries emit structured logs.
func initLogger() error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	l, err := platformlogger.New(cfg.LogLevel)
	if err != nil {
		return fmt.Errorf("init logger: %w", err)
	}
	zapLogger = l

	// Bridge the std log package to zap: std log writes become INFO zap entries.
	// The std logger is replaced globally; zap remains the primary logger.
	_ = log.New(zapWriter{logger: zapLogger}, "", 0)
	log.SetOutput(&zapStdWriter{logger: zapLogger})

	zapLogger.Info("cli command invoked", zap.String("env", cfg.Environment))
	return nil
}

// zapStdWriter adapts the std log output to zap INFO entries.
type zapStdWriter struct{ logger *zap.Logger }

func (w zapStdWriter) Write(p []byte) (int, error) {
	w.logger.Info(string(p))
	return len(p), nil
}

// zapWriter is retained for callers that want a dedicated log.Logger bound to zap.
type zapWriter struct{ logger *zap.Logger }

func (w zapWriter) Write(p []byte) (int, error) {
	w.logger.Info(string(p))
	return len(p), nil
}
