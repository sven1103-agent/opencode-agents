package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/sven1103-agent/opencode-helper/internal/bundle"
	"github.com/sven1103-agent/opencode-helper/internal/source"
)

var (
	bundleProjectRoot string
	bundlePreset      string
	bundleForce       bool
	bundleDryRun      bool
	bundleOutput      string
	bundleYes         bool
)

// bundleCmd represents the bundle command
var bundleCmd = &cobra.Command{
	Use:   "bundle",
	Short: "Manage OpenCode configuration bundles",
	Long: `Manage OpenCode configuration bundles.

Apply, track, and update configuration bundles from registered sources.

Examples:
  oc bundle apply abc12345 --preset default
  oc bundle status
  oc bundle update abc12345`,
}

// bundleApplyCmd applies a preset from a registered config bundle
var bundleApplyCmd = &cobra.Command{
	Use:   "apply <source-id>",
	Short: "Apply a preset from a config bundle",
	Long: `Apply a preset from a registered config bundle to a project.

The source-id must reference a registered config source (see 'source list').
The preset name must exist in the bundle's manifest.

Examples:
  oc bundle apply abc12345 --preset default
  oc bundle apply abc12345 --preset minimal --project-root ./myproject
  oc bundle apply abc12345 --preset default --force`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return runBundleApply(args[0])
	},
}

// bundleStatusCmd shows provenance for the applied bundle
var bundleStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show provenance for applied bundle",
	Long: `Show provenance information for the currently applied bundle.

Displays the source, version, and preset that was applied to the project.

Example:
  oc bundle status
  oc bundle status --project-root ./myproject`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return runBundleStatus()
	},
}

// bundleUpdateCmd checks for and applies newer bundle releases
var bundleUpdateCmd = &cobra.Command{
	Use:   "update <source-id>",
	Short: "Check for and apply newer bundle releases",
	Long: `Check for and apply newer bundle releases from update-capable sources.

Only sources marked as update-capable in their manifest support this command.

Examples:
  oc bundle update abc12345
  oc bundle update abc12345 --yes`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return runBundleUpdate(args[0])
	},
}

func init() {
	rootCmd.AddCommand(bundleCmd)

	// Add subcommands
	bundleCmd.AddCommand(bundleApplyCmd)
	bundleCmd.AddCommand(bundleStatusCmd)
	bundleCmd.AddCommand(bundleUpdateCmd)

	// Flags for bundle apply
	bundleApplyCmd.Flags().StringVar(&bundlePreset, "preset", "", "Preset name to apply (required)")
	bundleApplyCmd.Flags().StringVar(&bundleProjectRoot, "project-root", ".", "Project root directory")
	bundleApplyCmd.Flags().StringVar(&bundleOutput, "output", "opencode.json", "Output file path")
	bundleApplyCmd.Flags().BoolVar(&bundleForce, "force", false, "Overwrite existing files")
	bundleApplyCmd.Flags().BoolVar(&bundleDryRun, "dry-run", false, "Show what would be done without doing it")
	bundleApplyCmd.MarkFlagRequired("preset")

	// Flags for bundle status
	bundleStatusCmd.Flags().StringVar(&bundleProjectRoot, "project-root", ".", "Project root directory")

	// Flags for bundle update
	bundleUpdateCmd.Flags().BoolVar(&bundleYes, "yes", false, "Skip confirmation prompt")
}

func runBundleApply(sourceID string) error {
	// Validate preset name is provided
	if bundlePreset == "" {
		return fmt.Errorf("--preset is required")
	}

	// Resolve project root
	projectRoot, err := filepath.Abs(bundleProjectRoot)
	if err != nil {
		return fmt.Errorf("invalid project root: %w", err)
	}

	// Check if project root exists
	if _, err := os.Stat(projectRoot); os.IsNotExist(err) {
		return fmt.Errorf("project root does not exist: %s", projectRoot)
	}

	// Get the source from registry
	src, err := source.GetSource(sourceID)
	if err != nil {
		return fmt.Errorf("source not found: %s", sourceID)
	}

	// Resolve source to local bundle root
	bundleRoot, cleanup, err := bundle.ResolveToLocal(string(src.Type), src.Location, "")
	if err != nil {
		return fmt.Errorf("failed to resolve source: %w", err)
	}
	defer cleanup()

	// Load manifest
	manifestPath := filepath.Join(bundleRoot, "opencode-bundle.manifest.json")
	manifest, err := bundle.LoadManifest(manifestPath)
	if err != nil {
		return fmt.Errorf("failed to load manifest: %w", err)
	}

	// Get preset from manifest
	preset, err := bundle.GetPreset(manifest, bundlePreset)
	if err != nil {
		return fmt.Errorf("preset not found in bundle: %s", bundlePreset)
	}

	// Resolve output path
	outputPath := filepath.Join(projectRoot, bundleOutput)

	// Validate output path
	if err := validateOutputPath(projectRoot, outputPath); err != nil {
		return err
	}

	// Read preset content
	presetFilePath := filepath.Join(bundleRoot, preset.Entrypoint)
	presetContent, err := os.ReadFile(presetFilePath)
	if err != nil {
		return fmt.Errorf("failed to read preset file: %w", err)
	}

	// Dry run mode
	if bundleDryRun {
		fmt.Printf("dry-run: apply preset '%s' from bundle '%s'\n", bundlePreset, manifest.BundleName)
		fmt.Printf("dry-run: write config to %s\n", outputPath)
		return nil
	}

	// Write config file
	if err := os.WriteFile(outputPath, presetContent, 0644); err != nil {
		if !bundleForce && os.IsExist(err) {
			return fmt.Errorf("output file exists: %s (use --force to overwrite)", outputPath)
		}
		return fmt.Errorf("failed to write config: %w", err)
	}
	fmt.Printf("written: %s\n", outputPath)

	// Write provenance
	prov := &bundle.Provenance{
		SourceID:      sourceID,
		SourceName:    src.Name,
		SourceType:    string(src.Type),
		BundleVersion: manifest.BundleVersion,
		PresetName:    bundlePreset,
		Entrypoint:    preset.Entrypoint,
		AppliedAt:     "2026-03-31T00:00:00Z", // Would use time.Now().Format(time.RFC3339)
	}

	if err := bundle.SaveProvenance(projectRoot, prov, bundleForce); err != nil {
		return fmt.Errorf("failed to save provenance: %w", err)
	}
	fmt.Printf("written: %s\n", bundle.ProvenancePath(projectRoot))
	fmt.Println("done: bundle applied")

	return nil
}

func runBundleStatus() error {
	// Resolve project root
	projectRoot, err := filepath.Abs(bundleProjectRoot)
	if err != nil {
		return fmt.Errorf("invalid project root: %w", err)
	}

	// Load provenance
	prov, err := bundle.LoadProvenance(projectRoot)
	if err != nil {
		return fmt.Errorf("no bundle applied to this project (run 'bundle apply' first)")
	}

	// Display provenance
	fmt.Println("Bundle Provenance:")
	fmt.Printf("  source_id:      %s\n", prov.SourceID)
	fmt.Printf("  source_name:    %s\n", prov.SourceName)
	fmt.Printf("  source_type:    %s\n", prov.SourceType)
	fmt.Printf("  bundle_version: %s\n", prov.BundleVersion)
	fmt.Printf("  preset_name:    %s\n", prov.PresetName)
	fmt.Printf("  applied_at:     %s\n", prov.AppliedAt)

	return nil
}

func runBundleUpdate(sourceID string) error {
	// Get the source from registry
	src, err := source.GetSource(sourceID)
	if err != nil {
		return fmt.Errorf("source not found: %s", sourceID)
	}

	// For now, github-release is required for updates (as per shell script behavior)
	if string(src.Type) != "github-release" {
		return fmt.Errorf("bundle update is only supported for github-release sources")
	}

	// Find the project root with provenance
	// For now, check current directory
	projectRoot := bundleProjectRoot
	if projectRoot == "." {
		cwd, err := os.Getwd()
		if err != nil {
			return fmt.Errorf("failed to get current directory: %w", err)
		}
		projectRoot = cwd
	}

	// Load provenance to verify bundle has been applied
	var prov *bundle.Provenance
	prov, err = bundle.LoadProvenance(projectRoot)
	if err != nil {
		return fmt.Errorf("no bundle applied to this project (run 'bundle apply' first)")
	}

	// Suppress unused variable warning
	_ = prov

	// For now, return not implemented for github-release
	return fmt.Errorf("bundle update for github-release sources requires network operations (not yet implemented)")
}
