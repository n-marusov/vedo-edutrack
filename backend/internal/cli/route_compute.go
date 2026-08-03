package cli

import (
	"fmt"

	"github.com/spf13/cobra"
	"go.uber.org/zap"

	ontostub "vedo-edutrack/backend/internal/modules/ontologyport/adapters/stub"
	routestub "vedo-edutrack/backend/internal/modules/routeplanning/adapters/stub"
)

// newRouteComputeCmd builds the `route compute` subcommand — computes a route
// (--stub | from DB) as a dev/test tool.
func newRouteComputeCmd() *cobra.Command {
	parent := &cobra.Command{
		Use:   "route",
		Short: "Route operations (compute a route — dev/test tool)",
	}
	cmd := &cobra.Command{
		Use:   "compute",
		Short: "Compute a learning route (--stub | from DB)",
		RunE: func(cmd *cobra.Command, args []string) error {
			useStub, _ := cmd.Flags().GetBool("stub")
			if useStub {
				return runRouteComputeStub(args)
			}
			return stubNotImplemented("route compute")(cmd, args)
		},
	}
	cmd.Flags().Bool("stub", false, "use the in-memory stub graph")
	parent.AddCommand(cmd)
	return parent
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
