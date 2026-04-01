package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"
	"github.com/sven1103-agent/opencode-config-cli/internal/source"
)

var (
	sourceName string
)

// sourceCmd represents the source command
var sourceCmd = &cobra.Command{
	Use:   "source",
	Short: "Manage OpenCode config sources",
	Long: `Manage OpenCode configuration sources.

A config source is a location (local directory, archive, or GitHub release)
that contains OpenCode configuration bundles.

Examples:
  oc source add ./my-config-bundle
  oc source add ./release.tar.gz --name my-archive
  oc source add owner/repo --name my-config
  oc source add https://github.com/owner/repo/releases/tag/v1.2.3
  oc source list
  oc source remove abc12345`,
}

// sourceAddCmd adds a new config source
var sourceAddCmd = &cobra.Command{
	Use:   "add <location>",
	Short: "Register a new config source",
	Long: `Register a new config source.

The location can be:
  - A local directory containing a bundle
  - A local .tar.gz archive file
  - A GitHub repository or release URL

Examples:
  oc source add ./my-config-bundle
  oc source add ./release.tar.gz --name my-archive
  oc source add owner/repo
  oc source add github.com/owner/repo
  oc source add https://github.com/owner/repo/releases/tag/v1.2.3`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return runSourceAdd(args[0])
	},
}

// sourceListCmd lists all registered sources
var sourceListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all registered config sources",
	Long: `List all registered config sources.

Shows each source's ID, name, type, and location.

Example:
  oc source list`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return runSourceList()
	},
}

// sourceRemoveCmd removes a config source
var sourceRemoveCmd = &cobra.Command{
	Use:   "remove <id>",
	Short: "Remove a registered config source",
	Long: `Remove a registered config source by its ID.

Example:
  oc source remove abc12345`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return runSourceRemove(args[0])
	},
}

func init() {
	rootCmd.AddCommand(sourceCmd)

	// Add subcommands
	sourceCmd.AddCommand(sourceAddCmd)
	sourceCmd.AddCommand(sourceListCmd)
	sourceCmd.AddCommand(sourceRemoveCmd)

	// Flags for source add
	sourceAddCmd.Flags().StringVar(&sourceName, "name", "", "Friendly name for the source")
}

func runSourceAdd(location string) error {
	s, err := source.AddSource(location, sourceName)
	if err != nil {
		return fmt.Errorf("failed to add source: %w", err)
	}

	fmt.Printf("Source added successfully:\n")
	fmt.Printf("  ID:       %s\n", s.ID)
	fmt.Printf("  Name:     %s\n", s.Name)
	fmt.Printf("  Type:     %s\n", s.Type)
	fmt.Printf("  Location: %s\n", s.Location)
	fmt.Printf("  Created:  %s\n", s.CreatedAt)

	return nil
}

func runSourceList() error {
	sources, err := source.ListSources()
	if err != nil {
		return fmt.Errorf("failed to list sources: %w", err)
	}

	if len(sources) == 0 {
		fmt.Println("No sources registered.")
		fmt.Println("Use 'oc source add <location>' to register a source.")
		return nil
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintf(w, "ID\tNAME\tTYPE\tLOCATION\n")
	fmt.Fprintf(w, "%s\t%s\t%s\t%s\n", "───", "────", "────", "────────\n")

	for _, s := range sources {
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\n", s.ID, s.Name, s.Type, s.Location)
	}

	if err := w.Flush(); err != nil {
		return fmt.Errorf("failed to write output: %w", err)
	}

	return nil
}

func runSourceRemove(id string) error {
	// Check if source exists first
	_, err := source.GetSource(id)
	if err != nil {
		return fmt.Errorf("source not found: %s", id)
	}

	if err := source.RemoveSource(id); err != nil {
		return fmt.Errorf("failed to remove source: %w", err)
	}

	fmt.Printf("Source '%s' removed successfully.\n", id)
	return nil
}
