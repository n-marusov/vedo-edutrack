package cli

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"go.uber.org/zap"

	ontostub "vedo-edutrack/backend/internal/modules/ontologyport/adapters/stub"
	ontologyapp "vedo-edutrack/backend/internal/modules/ontologyport/application"
	ontologydomain "vedo-edutrack/backend/internal/modules/ontologyport/domain"
	routestub "vedo-edutrack/backend/internal/modules/routeplanning/adapters/stub"
	routeapp "vedo-edutrack/backend/internal/modules/routeplanning/application"
	routedomain "vedo-edutrack/backend/internal/modules/routeplanning/domain"
)

// newRouteComputeCmd builds the `route compute` subcommand — computes a route
// (--stub | cached ontology graph) as a dev/test tool.
func newRouteComputeCmd() *cobra.Command {
	parent := &cobra.Command{
		Use:   "route",
		Short: "Route operations (compute a route — dev/test tool)",
	}
	cmd := &cobra.Command{
		Use:   "compute",
		Short: "Compute a learning route (--stub | cached ontology graph)",
		RunE: func(cmd *cobra.Command, args []string) error {
			useStub, _ := cmd.Flags().GetBool("stub")
			if useStub {
				return runRouteComputeStub(args)
			}
			return runRouteCompute(cmd)
		},
	}
	cmd.Flags().Bool("stub", false, "use the in-memory stub graph")
	cmd.Flags().String("learner", "", "learner id")
	cmd.Flags().String("position", "", "current module id")
	cmd.Flags().String("goal", "", "goal module id")
	cmd.Flags().String("pedagogy", "", "pedagogy concept id")
	cmd.Flags().String("ontology-version", "", "cached ontology version")
	parent.AddCommand(cmd)
	return parent
}

func runRouteCompute(cmd *cobra.Command) error {
	learnerID, _ := cmd.Flags().GetString("learner")
	positionID, _ := cmd.Flags().GetString("position")
	goalID, _ := cmd.Flags().GetString("goal")
	pedagogyID, _ := cmd.Flags().GetString("pedagogy")
	ontologyVersion, _ := cmd.Flags().GetString("ontology-version")
	output, _ := cmd.Root().PersistentFlags().GetString("output")
	if learnerID == "" || positionID == "" || goalID == "" {
		return fmt.Errorf("--learner, --position, and --goal are required")
	}

	service := routeapp.NewComputeService(cachedGraphProvider{cache: ontologyapp.DefaultCache()}, routedomain.NewPathfinder(routedomain.DefaultWeightProfile()), zapLogger)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	result, err := service.Compute(ctx, routeapp.ComputeRequest{LearnerID: learnerID, PositionID: positionID, GoalID: goalID, PedagogyConceptID: pedagogyID, OntologyVersion: ontologyVersion})
	if err != nil {
		return fmt.Errorf("route compute: %w", err)
	}
	zapLogger.Info("route computed", zap.String("learner_id", learnerID), zap.String("goal_topic_id", goalID), zap.Int("topics", len(result.Route.Steps)))

	switch strings.ToLower(output) {
	case "json":
		return json.NewEncoder(cmd.OutOrStdout()).Encode(result)
	case "table", "":
		for _, step := range result.Route.Steps {
			fmt.Fprintf(cmd.OutOrStdout(), "%2d  %s  (%s)\n", step.Order+1, step.ModuleID, step.Via)
		}
		return nil
	default:
		return fmt.Errorf("unsupported output format %q (expected table|json)", output)
	}
}

type cachedGraphProvider struct {
	cache *ontologyapp.VersionedCache
}

func (p cachedGraphProvider) GraphForVersion(_ context.Context, ontologyVersion string) (routedomain.OntologyGraph, error) {
	if p.cache == nil {
		return routedomain.OntologyGraph{}, fmt.Errorf("ontology cache is not initialized")
	}
	candidates := []string{"mvp"}
	for _, ontologyID := range candidates {
		subgraph, _, ok := p.cache.Get(ontologyID)
		if !ok {
			continue
		}
		if ontologyVersion != "" && subgraph.Version != ontologyVersion {
			continue
		}
		return routeGraphFromSubgraph(subgraph), nil
	}
	return routedomain.OntologyGraph{}, fmt.Errorf("no cached ontology graph found; run ontology sync before route compute")
}

func routeGraphFromSubgraph(subgraph ontologydomain.Subgraph) routedomain.OntologyGraph {
	modules := make([]routedomain.Module, 0, len(subgraph.Modules))
	links := make([]routedomain.Link, 0, len(subgraph.Links))
	for _, module := range subgraph.Modules {
		modules = append(modules, routedomain.Module{ID: module.ID, Title: module.Title})
	}
	for _, link := range subgraph.Links {
		links = append(links, routedomain.Link{SourceID: link.SourceModuleID, TargetID: link.TargetModuleID, Type: routedomain.LinkType(link.Type)})
	}
	return routedomain.OntologyGraph{Modules: modules, Links: links}
}

// runRouteComputeStub computes a route over the stub graph. Args: [goal_topic_id]
// (learner id defaults to "cli-user").
func runRouteComputeStub(args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("usage: vedo-edutrack route compute --stub <goal_topic_id>")
	}
	goalTopicID := args[0]
	learnerID := "cli-user"
	if len(args) > 1 {
		learnerID = args[1]
	}

	graph := ontostub.NewGraph()
	computer := routestub.NewComputer(graph)

	route, err := computer.ComputeRoute(learnerID, goalTopicID)
	if err != nil {
		return fmt.Errorf("route compute: %w", err)
	}

	zapLogger.Info("route computed",
		zap.String("learner_id", learnerID),
		zap.String("goal_topic_id", goalTopicID),
		zap.Int("topics", len(route)),
	)

	for _, t := range route {
		fmt.Printf("%2d  %s  (%s)\n", t.Order+1, t.TopicID, t.LinkType)
	}
	return nil
}
