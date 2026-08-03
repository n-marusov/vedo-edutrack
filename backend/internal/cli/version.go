package cli

import (
	"fmt"

	"github.com/spf13/cobra"

	"vedo-edutrack/backend/internal/platform/config"
)

// newVersionCmd builds the `version` subcommand — prints the build version
// (injected via ldflags, ADR-DES.INFRA.dynamic-config-injection).
func newVersionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print the build version",
		RunE: func(_ *cobra.Command, _ []string) error {
			cfg, err := config.Load()
			if err != nil {
				return err
			}
			fmt.Printf("vedo-edutrack %s (env %s)\n", config.Version, cfg.Environment)
			return nil
		},
	}
}
