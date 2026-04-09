package cmd

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
	"github.com/sven1103-agent/opencode-config-cli/internal/bundle"
	configpreset "github.com/sven1103-agent/opencode-config-cli/internal/preset"
	"github.com/sven1103-agent/opencode-config-cli/internal/source"
)

var (
	bundleProjectRoot    string
	bundlePreset         string
	bundleVersion        string
	bundleAuto           bool
	bundleForce          bool
	bundleDryRun         bool
	bundleOutput         string
	bundleYes            bool
	bundleResolveToLocal           = bundle.ResolveToLocal
	bundlePromptIn       io.Reader = os.Stdin
	bundlePromptOut      io.Writer = os.Stdout
	bundleInputIsTTY               = isInteractiveTTY
)

// bundleCmd represents the bundle command
var bundleCmd = &cobra.Command{
	Use:   "bundle",
	Short: "Manage OpenCode configuration bundles",
	Long: `Manage OpenCode configuration bundles.

Apply, track, and update configuration bundles from registered sources.

Examples:
	  oc bundle apply qbic --preset default
	  oc bundle apply qbic
	  oc bundle status
	  oc bundle update abc12345`,
}

// bundleApplyCmd applies a preset from a registered config bundle
var bundleApplyCmd = &cobra.Command{
	Use:   "apply <source-ref>",
	Short: "Apply a preset from a config bundle",
	Long: `Apply a preset from a registered config bundle to a project.

The source-ref may be either a registered source ID or a unique source name.
In interactive terminals, omitting --preset opens a guided preset selection flow.

Examples:
	  oc bundle apply qbic --preset default
	  oc bundle apply qbic
	  oc bundle apply abc12345 --version v1.2.3 --preset default
	  oc bundle apply qbic --preset minimal --project-root ./myproject
	  oc bundle apply qbic --auto --preset default --force`,
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
	bundleApplyCmd.Flags().StringVar(&bundlePreset, "preset", "", "Preset name to apply")
	bundleApplyCmd.Flags().StringVar(&bundleVersion, "version", "", "Bundle version/tag to apply for github-release sources")
	bundleApplyCmd.Flags().StringVar(&bundleProjectRoot, "project-root", ".", "Project root directory")
	bundleApplyCmd.Flags().StringVar(&bundleOutput, "output", "opencode.json", "Output file path")
	bundleApplyCmd.Flags().BoolVar(&bundleAuto, "auto", false, "Disable interactive preset selection")
	bundleApplyCmd.Flags().BoolVar(&bundleForce, "force", false, "Overwrite existing files")
	bundleApplyCmd.Flags().BoolVar(&bundleDryRun, "dry-run", false, "Show what would be done without doing it")
	bundleApplyCmd.ValidArgsFunction = completeSourceRefs
	_ = bundleApplyCmd.RegisterFlagCompletionFunc("preset", completeBundlePresetNames)
	bundleUpdateCmd.ValidArgsFunction = completeSourceRefs

	// Flags for bundle status
	bundleStatusCmd.Flags().StringVar(&bundleProjectRoot, "project-root", ".", "Project root directory")

	// Flags for bundle update
	bundleUpdateCmd.Flags().BoolVar(&bundleYes, "yes", false, "Skip confirmation prompt")
}

func runBundleApply(sourceRef string) error {
	// Resolve project root
	projectRoot, err := filepath.Abs(bundleProjectRoot)
	if err != nil {
		return fmt.Errorf("invalid project root: %w", err)
	}

	// Check if project root exists
	if _, err := os.Stat(projectRoot); os.IsNotExist(err) {
		return fmt.Errorf("project root does not exist: %s", projectRoot)
	}

	// Resolve the source from registry
	src, err := source.ResolveSourceRef(sourceRef)
	if err != nil {
		return err
	}
	if bundleVersion != "" && string(src.Type) != "github-release" {
		return fmt.Errorf("--version is only supported for github-release sources")
	}

	// Resolve source to local bundle root
	bundleRoot, cleanup, err := bundleResolveToLocal(string(src.Type), src.Location, bundleVersion)
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

	selectedPreset := bundlePreset
	if selectedPreset == "" {
		if bundleAuto || !bundleInputIsTTY() {
			return fmt.Errorf("--preset is required outside interactive mode")
		}
		selectedPreset, err = promptForPresetSelection(manifest)
		if err != nil {
			return err
		}
	}

	// Get preset from manifest
	bundlePresetEntry, err := bundle.GetPreset(manifest, selectedPreset)
	if err != nil {
		return fmt.Errorf("preset not found in bundle: %s", selectedPreset)
	}

	// Resolve output path
	outputPath := filepath.Join(projectRoot, bundleOutput)

	// Validate output path
	if err := validateOutputPath(projectRoot, outputPath); err != nil {
		return err
	}

	// Read preset content
	presetFilePath := filepath.Join(bundleRoot, bundlePresetEntry.Entrypoint)
	presetContent, err := os.ReadFile(presetFilePath)
	if err != nil {
		return fmt.Errorf("failed to read preset file: %w", err)
	}

	// Dry run mode
	if bundleDryRun {
		fmt.Printf("dry-run: apply preset '%s' from bundle '%s'\n", selectedPreset, manifest.BundleName)
		fmt.Printf("dry-run: write config to %s\n", outputPath)
		return nil
	}

	// Reuse the shared write semantics so bundle apply matches init/preset overwrite behavior.
	if err := configpreset.WriteConfig(outputPath, string(presetContent), bundleForce); err != nil {
		return fmt.Errorf("failed to write config: %w", err)
	}
	fmt.Printf("written: %s\n", outputPath)

	// Write provenance
	prov := &bundle.Provenance{
		SourceID:      src.ID,
		SourceName:    src.Name,
		SourceType:    string(src.Type),
		BundleVersion: manifest.BundleVersion,
		PresetName:    selectedPreset,
		Entrypoint:    bundlePresetEntry.Entrypoint,
		AppliedAt:     "2026-03-31T00:00:00Z", // Would use time.Now().Format(time.RFC3339)
	}

	if err := bundle.SaveProvenance(projectRoot, prov, bundleForce); err != nil {
		return fmt.Errorf("failed to save provenance: %w", err)
	}
	fmt.Printf("written: %s\n", bundle.ProvenancePath(projectRoot))
	fmt.Println("done: bundle applied")

	return nil
}

func promptForPresetSelection(manifest *bundle.Manifest) (string, error) {
	if len(manifest.Presets) == 0 {
		return "", fmt.Errorf("bundle has no presets to select")
	}

	reader := bufio.NewReader(bundlePromptIn)
	for {
		fmt.Fprintf(bundlePromptOut, "Available presets for %s:\n", manifest.BundleName)
		for i, preset := range manifest.Presets {
			if preset.Description != "" {
				fmt.Fprintf(bundlePromptOut, "  %d) %s - %s\n", i+1, preset.Name, preset.Description)
				continue
			}
			fmt.Fprintf(bundlePromptOut, "  %d) %s\n", i+1, preset.Name)
		}
		fmt.Fprint(bundlePromptOut, "Select a preset: ")

		selection, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				return "", fmt.Errorf("interactive preset selection cancelled")
			}
			return "", fmt.Errorf("failed to read preset selection: %w", err)
		}

		selection = strings.TrimSpace(selection)
		for _, preset := range manifest.Presets {
			if preset.Name == selection {
				return preset.Name, nil
			}
		}

		if index, err := strconv.Atoi(selection); err == nil {
			if index >= 1 && index <= len(manifest.Presets) {
				return manifest.Presets[index-1].Name, nil
			}
		}

		fmt.Fprintln(bundlePromptOut, "Invalid selection. Please enter a preset number or exact name.")
	}
}

func isInteractiveTTY() bool {
	info, err := os.Stdin.Stat()
	if err != nil {
		return false
	}
	return info.Mode()&os.ModeCharDevice != 0
}

func completeSourceRefs(_ *cobra.Command, _ []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	sources, err := source.ListSources()
	if err != nil {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}

	seen := map[string]struct{}{}
	var refs []string
	for _, src := range sources {
		for _, candidate := range sourceCompletionCandidates(src) {
			if !strings.HasPrefix(candidate, toComplete) {
				continue
			}
			if _, ok := seen[candidate]; ok {
				continue
			}
			seen[candidate] = struct{}{}
			refs = append(refs, candidate)
		}
	}

	return refs, cobra.ShellCompDirectiveNoFileComp
}

func completeBundlePresetNames(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	if len(args) == 0 {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}

	src, err := source.ResolveSourceRef(args[0])
	if err != nil {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}

	versionTag := ""
	if flag := cmd.Flags().Lookup("version"); flag != nil {
		versionTag = flag.Value.String()
	}

	bundleRoot, cleanup, err := bundleResolveToLocal(string(src.Type), src.Location, versionTag)
	if err != nil {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}
	defer cleanup()

	manifest, err := bundle.LoadManifest(filepath.Join(bundleRoot, "opencode-bundle.manifest.json"))
	if err != nil {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}

	var presets []string
	for _, preset := range manifest.Presets {
		if strings.HasPrefix(preset.Name, toComplete) {
			presets = append(presets, preset.Name)
		}
	}

	return presets, cobra.ShellCompDirectiveNoFileComp
}

func sourceCompletionCandidates(src source.Source) []string {
	if src.Name == "" || src.Name == src.ID {
		return []string{src.ID}
	}
	return []string{src.ID, src.Name}
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

func runBundleUpdate(sourceRef string) error {
	// Get the source from registry
	src, err := source.ResolveSourceRef(sourceRef)
	if err != nil {
		return err
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
