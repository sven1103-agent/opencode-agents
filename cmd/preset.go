package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"text/tabwriter"

	"github.com/spf13/cobra"
	"github.com/sven1103-agent/opencode-config-cli/internal/bundle"
	"github.com/sven1103-agent/opencode-config-cli/internal/preset"
	"github.com/sven1103-agent/opencode-config-cli/internal/source"
)

var (
	presetProjectRoot    string
	presetOutput         string
	presetForce          bool
	presetDryRun         bool
	presetSources        bool
	presetResolveToLocal = bundle.ResolveToLocal
)

// presetCmd represents the preset command
var presetCmd = &cobra.Command{
	Use:   "preset",
	Short: "Manage OpenCode presets",
	Long: `Manage OpenCode configuration presets.

Use 'preset list' for built-in presets or 'preset list --sources' to browse registered sources.`,
}

// presetListCmd lists all available presets
var presetListCmd = &cobra.Command{
	Use:   "list",
	Short: "List available presets",
	Long:  "List built-in presets or browse presets across registered sources with --sources",
	RunE: func(cmd *cobra.Command, args []string) error {
		return runPresetList()
	},
}

// presetUseCmd applies a preset to the project
var presetUseCmd = &cobra.Command{
	Use:   "use [preset-name]",
	Short: "Apply a preset to the project",
	Long: `Apply a preset configuration to the project.

The preset name is one of: mixed, openai, big-pickle, minimax, kimi.

Examples:
  oc preset use openai
  oc preset use minimax --project-root /path/to/project
  oc preset use kimi --output custom-config.json --force`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return runPresetUse(args[0])
	},
}

func init() {
	// Add preset command to root
	rootCmd.AddCommand(presetCmd)

	// Add list subcommand
	presetCmd.AddCommand(presetListCmd)

	// Add use subcommand
	presetCmd.AddCommand(presetUseCmd)

	// Flags for preset use
	presetUseCmd.Flags().StringVar(&presetProjectRoot, "project-root", ".", "Project root directory")
	presetUseCmd.Flags().StringVar(&presetOutput, "output", "opencode.json", "Output file path")
	presetUseCmd.Flags().BoolVar(&presetForce, "force", false, "Overwrite existing files")
	presetUseCmd.Flags().BoolVar(&presetDryRun, "dry-run", false, "Show what would be done without doing it")
	presetListCmd.Flags().BoolVar(&presetSources, "sources", false, "List presets across registered sources")
}

func runPresetList() error {
	if presetSources {
		return runSourcePresetList()
	}

	presets := preset.ValidPresets()
	fmt.Println("Available presets:")
	for _, p := range presets {
		fmt.Printf("  - %s\n", p)
	}
	return nil
}

func runSourcePresetList() error {
	sources, err := source.ListSources()
	if err != nil {
		return fmt.Errorf("failed to list sources: %w", err)
	}
	if len(sources) == 0 {
		return fmt.Errorf("no sources registered. Use 'oc source add <location>' first")
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintf(w, "SOURCE_REF\tSOURCE_ID\tBUNDLE_VERSION\tPRESET\tDESCRIPTION\n")
	fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\n", "──────────", "─────────", "──────────────", "──────", "───────────")

	foundPreset := false
	for _, src := range sources {
		bundleRoot, cleanup, err := presetResolveToLocal(string(src.Type), src.Location, "")
		if err != nil {
			fmt.Fprintf(os.Stderr, "warning: failed to inspect source %s (%s): %v\n", src.Name, src.ID, err)
			continue
		}

		manifest, err := bundle.LoadManifest(filepath.Join(bundleRoot, "opencode-bundle.manifest.json"))
		cleanup()
		if err != nil {
			fmt.Fprintf(os.Stderr, "warning: failed to inspect source %s (%s): %v\n", src.Name, src.ID, err)
			continue
		}

		ref := src.ID
		if src.Name != "" {
			ref = src.Name
		}
		for _, preset := range manifest.Presets {
			foundPreset = true
			fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\n", ref, src.ID, manifest.BundleVersion, preset.Name, preset.Description)
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

func runPresetUse(presetName string) error {
	// Validate preset name
	validPresets := preset.ValidPresets()
	valid := false
	for _, p := range validPresets {
		if p == presetName {
			valid = true
			break
		}
	}
	if !valid {
		return fmt.Errorf("invalid preset: %s (valid presets: %v)", presetName, validPresets)
	}

	// Resolve project root
	projectRoot, err := filepath.Abs(presetProjectRoot)
	if err != nil {
		return fmt.Errorf("invalid project root: %w", err)
	}

	// Check if project root exists
	if _, err := os.Stat(projectRoot); os.IsNotExist(err) {
		return fmt.Errorf("project root does not exist: %s", projectRoot)
	}

	// Resolve output path
	outputPath := filepath.Join(projectRoot, presetOutput)

	// Validate output path to prevent path traversal
	if err := validateOutputPath(projectRoot, outputPath); err != nil {
		return err
	}

	// Get preset config (from repo bundled files)
	configData, err := getPresetConfig(presetName)
	if err != nil {
		return fmt.Errorf("failed to get preset config: %w", err)
	}

	// Dry run mode
	if presetDryRun {
		fmt.Printf("dry-run: apply preset '%s' to %s\n", presetName, outputPath)
		return nil
	}

	// Write config file
	if err := preset.WriteConfig(outputPath, configData, presetForce); err != nil {
		return fmt.Errorf("failed to write config: %w", err)
	}
	fmt.Printf("written: %s\n", outputPath)
	fmt.Println("done: preset applied")

	return nil
}

// getPresetConfig returns the config for a given preset name
func getPresetConfig(presetName string) (string, error) {
	// Map preset names to config files
	presetFiles := map[string]string{
		"mixed":      "opencode.mixed.json",
		"openai":     "opencode.openai.json",
		"big-pickle": "opencode.big-pickle.json",
		"minimax":    "opencode.minimax.json",
		"kimi":       "opencode.kimi.json",
	}

	filename, ok := presetFiles[presetName]
	if !ok {
		return "", fmt.Errorf("unknown preset: %s", presetName)
	}

	// First try: read from bundled configs directory (relative to executable)
	execPath, err := os.Executable()
	if err == nil {
		execDir := filepath.Dir(execPath)
		bundledPath := filepath.Join(execDir, "presets", filename)
		if data, err := os.ReadFile(bundledPath); err == nil {
			return string(data), nil
		}
	}

	// Second try: development mode (check worktree root)
	if data, err := os.ReadFile(filename); err == nil {
		return string(data), nil
	}

	return "", fmt.Errorf("preset config not found: %s", filename)
}
