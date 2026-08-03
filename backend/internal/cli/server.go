package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

// newServerCmd builds the `server` subcommand — the long-running HTTP process
// (API + metrics + health + embedded SPA). See ADR-DES.API.cli-interface.
func newServerCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "server",
		Short: "Start the HTTP server (API + metrics + health + embedded SPA)",
		RunE: func(cmd *cobra.Command, _ []string) error {
			zapLogger.Info("server subcommand invoked")
			if err := runServer(cmd); err != nil {
				return fmt.Errorf("server: %w", err)
			}
			return nil
		},
	}
}

// runServer starts the production HTTP server (chi router with graceful
// shutdown). Full implementation lands in T4.
func runServer(cmd *cobra.Command) error {
	return serveHTTP(cmd)
}
