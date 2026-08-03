package cli

import (
	"github.com/spf13/cobra"
)

// newOntologySyncCmd builds the `ontology sync` subcommand — copies a
// subgraph from VEDO Hub (F0.2). Stub in M0.3.
func newOntologySyncCmd() *cobra.Command {
	parent := &cobra.Command{
		Use:   "ontology",
		Short: "Ontology operations (sync subgraph from VEDO Hub)",
	}
	parent.AddCommand(&cobra.Command{
		Use:   "sync",
		Short: "Copy ontology subgraph from VEDO Hub (F0.2)",
		RunE:  stubNotImplemented("ontology sync"),
	})
	return parent
}
