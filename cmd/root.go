package cmd

import (
	"fmt"
	"os"

	"github.com/qbicsoftware/occo/internal/styles"
	"github.com/qbicsoftware/occo/internal/version"
	"github.com/spf13/cobra"
)

var versionFlag bool
var colorMode string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "occo",
	Short: "occo - OpenCode configuration manager",
	Long: `occo is the OpenCode configuration manager CLI.

Manage OpenCode configurations, including presets, sources, and bundle operations.`,
	Run: func(cmd *cobra.Command, args []string) {
		if versionFlag {
			fmt.Printf("occo %s\n", version.Version)
			os.Exit(0)
		}
		_ = cmd.Help()
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() error {
	// Set color mode from flag
	styles.SetColorMode(styles.ColorMode(colorMode))

	rootCmd.PersistentFlags().BoolVarP(&versionFlag, "version", "v", false, "Print version information")
	rootCmd.PersistentFlags().StringVar(&colorMode, "color", "always", "Color output mode: auto, always, never")
	return rootCmd.Execute()
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
