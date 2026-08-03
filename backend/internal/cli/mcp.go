package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

// newMcpCmd builds the `mcp` subcommand — MCP server over stdio (F6.6).
// Not yet implemented in M0.3.
func newMcpCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "mcp",
		Short: "MCP server over stdio for AI agents (F6.6)",
		RunE: func(_ *cobra.Command, _ []string) error {
			zapLogger.Info("mcp subcommand invoked (stub)")
			fmt.Println("MCP server not yet implemented")
			return nil
		},
	}
}
