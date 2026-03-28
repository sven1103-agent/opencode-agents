package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/sven1103-agent/opencode-helper/internal/version"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print version information",
	Long:  `Print the version information for oc.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("oc %s\n", version.Version)
	},
}
