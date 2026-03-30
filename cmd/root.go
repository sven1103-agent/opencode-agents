package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/sven1103-agent/opencode-helper/internal/version"
)

var versionFlag bool

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "oc",
	Short: "oc - OpenCode configuration manager",
	Long: `oc is the OpenCode configuration manager CLI.

Manage OpenCode configurations, including presets, sources, and bundle operations.`,
	Run: func(cmd *cobra.Command, args []string) {
		if versionFlag {
			fmt.Printf("oc %s\n", version.Version)
			os.Exit(0)
		}
		_ = cmd.Help()
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() error {
	rootCmd.Flags().BoolVarP(&versionFlag, "version", "v", false, "Print version information")
	return rootCmd.Execute()
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
