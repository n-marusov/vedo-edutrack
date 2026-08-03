package cli

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"go.uber.org/zap"

	"vedo-edutrack/backend/internal/modules/ontologyport/adapters/hub"
	ontologyapp "vedo-edutrack/backend/internal/modules/ontologyport/application"
	"vedo-edutrack/backend/internal/platform/config"
)

var mvpSubjectConceptIDs = []string{
	"subject:math",
	"subject:biology",
	"subject:physics",
	"subject:chemistry",
	"subject:history",
	"subject:literature",
	"subject:geography",
	"subject:computer-science",
	"subject:social-studies",
}

// newOntologySyncCmd builds the `ontology sync` subcommand — copies a
// subgraph from VEDO Hub (F0.2).
func newOntologySyncCmd() *cobra.Command {
	parent := &cobra.Command{
		Use:   "ontology",
		Short: "Ontology operations (sync subgraph from VEDO Hub)",
	}

	syncCmd := &cobra.Command{
		Use:   "sync",
		Short: "Copy ontology subgraph from VEDO Hub (F0.2)",
		RunE: func(cmd *cobra.Command, _ []string) error {
			return runOntologySync(cmd)
		},
	}
	syncCmd.Flags().String("ontology", "mvp", "Hub ontology id")
	syncCmd.Flags().StringSlice("subject", nil, "subject/concept id to sync; repeat or comma-separate (default: all MVP subjects)")
	syncCmd.Flags().Int("depth", 4, "GraphQL graphNeighborhood traversal depth")
	syncCmd.Flags().Bool("incremental", false, "diff against cached ontology version and report changed modules")
	parent.AddCommand(syncCmd)
	return parent
}

func runOntologySync(cmd *cobra.Command) error {
	ontologyID, _ := cmd.Flags().GetString("ontology")
	subjects, _ := cmd.Flags().GetStringSlice("subject")
	depth, _ := cmd.Flags().GetInt("depth")
	incremental, _ := cmd.Flags().GetBool("incremental")
	output, _ := cmd.Root().PersistentFlags().GetString("output")
	if len(subjects) == 0 {
		subjects = mvpSubjectConceptIDs
	}

	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}
	hubClient, err := hub.NewClientFromConfig(cfg, zapLogger)
	if err != nil {
		return fmt.Errorf("init Hub GraphQL client: %w", err)
	}
	service := ontologyapp.NewSyncService(hubClient, ontologyapp.DefaultCache(), zapLogger)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	result, err := service.Sync(ctx, ontologyapp.SyncOptions{
		OntologyID:  ontologyID,
		ConceptIDs:  subjects,
		Depth:       depth,
		Incremental: incremental,
	})
	if err != nil {
		zapLogger.Error("ontology sync failed", zap.Error(err))
		return fmt.Errorf("ontology sync: %w", err)
	}
	zapLogger.Info("ontology sync completed", zap.String("ontology_id", result.OntologyID), zap.Int("modules", result.ModuleCount), zap.Int("links", result.LinkCount), zap.Duration("duration", result.Duration))

	switch strings.ToLower(output) {
	case outputJSON:
		return json.NewEncoder(cmd.OutOrStdout()).Encode(result)
	case outputTable, "":
		fmt.Fprintf(cmd.OutOrStdout(), "ontology: %s\n", result.OntologyID)
		fmt.Fprintf(cmd.OutOrStdout(), "version: %s\n", result.Version)
		fmt.Fprintf(cmd.OutOrStdout(), "modules: %d\n", result.ModuleCount)
		fmt.Fprintf(cmd.OutOrStdout(), "links: %d\n", result.LinkCount)
		fmt.Fprintf(cmd.OutOrStdout(), "changed_modules: %d\n", result.ChangedModules)
		fmt.Fprintf(cmd.OutOrStdout(), "cached_at: %s\n", result.CachedAt.Format(time.RFC3339))
		fmt.Fprintf(cmd.OutOrStdout(), "duration_ms: %d\n", result.Duration.Milliseconds())
		return nil
	default:
		return fmt.Errorf("unsupported output format %q (expected table|json)", output)
	}
}
