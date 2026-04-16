package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"text/tabwriter"

	"github.com/qbicsoftware/occo/internal/bundle"
	"github.com/qbicsoftware/occo/internal/source"
	"github.com/qbicsoftware/occo/internal/styles"
	"github.com/spf13/cobra"
)

var (
	sourceName        string
	sourceWithPresets bool
)

// sourceCmd represents the source command
var sourceCmd = &cobra.Command{
	Use:   "source",
	Short: "Manage OpenCode config sources",
	Long: `Manage OpenCode configuration sources.

A config source is a location (local directory, archive, or GitHub release)
that contains OpenCode configuration bundles.

Examples:
  occo source add ./my-config-bundle
  occo source add ./release.tar.gz --name my-archive
  occo source add owner/repo --name my-config
  occo source add https://github.com/owner/repo/releases/tag/v1.2.3
  occo source list
  occo source remove abc12345`,
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
  occo source add ./my-config-bundle
  occo source add ./release.tar.gz --name my-archive
  occo source add owner/repo
  occo source add github.com/owner/repo
  occo source add https://github.com/owner/repo/releases/tag/v1.2.3`,
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
  occo source list`,
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
  occo source remove abc12345`,
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

	// Flags for source list
	sourceListCmd.Flags().BoolVar(&sourceWithPresets, "with-presets", false, "Show presets from all registered sources")
}

func runSourceAdd(location string) error {
	s, err := source.AddSource(location, sourceName)
	if err != nil {
		return fmt.Errorf("failed to add source: %w", err)
	}

	fmt.Println()
	fmt.Println(styles.Success("Source added successfully"))
	fmt.Println(styles.KeyValue("ID", s.ID))
	fmt.Println(styles.KeyValue("Name", s.Name))
	fmt.Println(styles.KeyValue("Type", string(s.Type)))
	fmt.Println(styles.KeyValue("Location", s.Location))
	fmt.Println(styles.KeyValueMuted("Created", s.CreatedAt))

	return nil
}

func runSourceList() error {
	if sourceWithPresets {
		return runSourceListWithPresets()
	}

	sources, err := source.ListSources()
	if err != nil {
		return fmt.Errorf("failed to list sources: %w", err)
	}

	if len(sources) == 0 {
		fmt.Println()
		fmt.Println(styles.Info("No sources registered."))
		fmt.Println(styles.Muted("Use 'occo source add <location>' to register a source."))
		return nil
	}

	// Sort sources by type, then by name
	for i := 0; i < len(sources)-1; i++ {
		for j := i + 1; j < len(sources); j++ {
			if sources[i].Type > sources[j].Type ||
				(sources[i].Type == sources[j].Type && sources[i].Name > sources[j].Name) {
				sources[i], sources[j] = sources[j], sources[i]
			}
		}
	}

	fmt.Println()
	fmt.Println(styles.SectionHeader("Registered Sources"))

	// Build table data
	headers := []string{"ID", "NAME", "TYPE", "LOCATION"}
	var rows [][]string
	for _, s := range sources {
		rows = append(rows, []string{s.ID, s.Name, string(s.Type), s.Location})
	}

	// Render modern table
	fmt.Println(styles.TableStyle(headers, rows))

	return nil
}

func runSourceListWithPresets() error {
	sources, err := source.ListSources()
	if err != nil {
		return fmt.Errorf("failed to list sources: %w", err)
	}
	if len(sources) == 0 {
		return fmt.Errorf("no sources registered. Use 'occo source add <location>' first")
	}

	fmt.Println()
	fmt.Println(styles.SectionHeader("Available Presets"))

	w := tabwriter.NewWriter(os.Stdout, 20, 0, 2, ' ', 0)
	fmt.Fprintln(w, "  SOURCE\t  ID\t  VERSION\t  PRESET\t  DESCRIPTION")
	fmt.Fprintln(w, "  ──────\t  ──\t  ───────\t  ──────\t  ───────────")

	foundPreset := false
	for _, src := range sources {
		versionTag := ""
		if string(src.Type) == "github-release" {
			versionTag, err = inspectGitHubBundleVersion(src.Location, "")
			if err != nil {
				fmt.Fprintln(os.Stderr, styles.Warning(fmt.Sprintf("failed to inspect source %s (%s): %v", src.Name, src.ID, err)))
				continue
			}
		}

		bundleRoot, cleanup, err := bundle.ResolveToLocal(string(src.Type), src.Location, versionTag)
		if err != nil {
			fmt.Fprintln(os.Stderr, styles.Warning(fmt.Sprintf("failed to inspect source %s (%s): %v", src.Name, src.ID, err)))
			continue
		}

		manifest, err := bundle.LoadManifest(filepath.Join(bundleRoot, "opencode-bundle.manifest.json"))
		cleanup()
		if err != nil {
			fmt.Fprintln(os.Stderr, styles.Warning(fmt.Sprintf("failed to inspect source %s (%s): %v", src.Name, src.ID, err)))
			continue
		}

		ref := src.ID
		if src.Name != "" {
			ref = src.Name
		}
		for _, preset := range manifest.Presets {
			foundPreset = true
			desc := preset.Description
			if desc == "" {
				desc = "-"
			}
			fmt.Fprintf(w, "  %s\t  %s\t  %s\t  %s\t  %s\n", ref, src.ID, manifest.BundleVersion, preset.Name, desc)
		}
	}

	if err := w.Flush(); err != nil {
		return fmt.Errorf("failed to write output: %w", err)
	}
	if !foundPreset {
		return fmt.Errorf("no inspectable source presets found")
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

	fmt.Println(styles.Success(fmt.Sprintf("Source '%s' removed successfully.", id)))
	return nil
}
