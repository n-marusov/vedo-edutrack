package cli

import (
	"fmt"

	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

// newMigrateCmd builds the `migrate` subcommand — Atlas migrations.
// Subcommands: up, down, validate — stubs printing "not yet implemented".
func newMigrateCmd() *cobra.Command {
	migrateCmd := &cobra.Command{
		Use:   "migrate",
		Short: "Manage database migrations (Atlas)",
		RunE: func(cmd *cobra.Command, _ []string) error {
			zapLogger.Info("migrate subcommand invoked (no subcommand)")
			return cmd.Help()
		},
	}

	migrateCmd.AddCommand(
		&cobra.Command{
			Use:   "up",
			Short: "Apply all pending migrations",
			RunE:  stubNotImplemented("migrate up"),
		},
		&cobra.Command{
			Use:   "down",
			Short: "Revert the last migration",
			RunE:  stubNotImplemented("migrate down"),
		},
		&cobra.Command{
			Use:   "validate",
			Short: "Validate migration drift (drift = 0)",
			RunE:  stubNotImplemented("migrate validate"),
		},
	)

	return migrateCmd
}

// stubNotImplemented returns a RunE that logs and prints "not yet implemented".
func stubNotImplemented(cmdPath string) func(*cobra.Command, []string) error {
	return func(_ *cobra.Command, _ []string) error {
		zapLogger.Info("stub subcommand invoked", zap.String("command", cmdPath))
		fmt.Printf("%s: not yet implemented\n", cmdPath)
		return nil
	}
}
