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
	bundleProjectRoot        string
	bundlePreset             string
	bundleVersion            string
	bundleAuto               bool
	bundleForce              bool
	bundleDryRun             bool
	bundleOutput             string
	bundleYes                bool
	bundleResolveToLocal               = bundle.ResolveToLocal
	bundleListGitHubReleases           = bundle.ListGitHubReleases
	bundlePromptIn           io.Reader = os.Stdin
	bundlePromptOut          io.Writer = os.Stdout
	bundleInputIsTTY                   = isInteractiveTTY
)

// bundleCmd represents the bundle command
var bundleCmd = &cobra.Command{
	Use:   "bundle",
	Short: "Manage OpenCode configuration bundles",
	Long: `Manage OpenCode configuration bundles.

Install, track, and update configuration bundles from registered sources.

Examples:
	  oc bundle install qbic --preset default
	  oc bundle install qbic
	  oc bundle status
	  oc bundle update abc12345`,
}

// bundleInstallCmd applies a preset from a registered config bundle
var bundleInstallCmd = &cobra.Command{
	Use:   "install [source-ref]",
	Short: "Install a preset from a config bundle",
	Long: `Install a preset from a registered config bundle to a project.

The source-ref may be either a registered source ID or a unique source name.
When omitted in interactive mode, the command prompts for source and preset selection.

Examples:
	  oc bundle install qbic --preset default
	  oc bundle install qbic
	  oc bundle install abc12345 --version v1.2.3 --preset default
	  oc bundle install qbic --preset minimal --project-root ./myproject
	  oc bundle install qbic --auto --preset default --force
	  oc bundle install (interactive mode)`,
	Args: cobra.RangeArgs(0, 1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) > 0 {
			return runBundleInstall(args[0], false)
		}
		// No arguments provided - enter interactive mode
		if !bundleInputIsTTY() || bundleAuto {
			return fmt.Errorf("source-ref is required in non-interactive mode or when using --auto flag")
		}
		// Interactive source and preset selection
		sourceRef, err := promptForSourceSelection()
		if err != nil {
			return err
		}
		return runBundleInstall(sourceRef, true)
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
	bundleCmd.AddCommand(bundleInstallCmd)
	bundleCmd.AddCommand(bundleStatusCmd)
	bundleCmd.AddCommand(bundleUpdateCmd)

	// Flags for bundle install
	bundleInstallCmd.Flags().StringVar(&bundlePreset, "preset", "", "Preset name to install")
	bundleInstallCmd.Flags().StringVar(&bundleVersion, "version", "", "Bundle version/tag to install for github-release sources")
	bundleInstallCmd.Flags().StringVar(&bundleProjectRoot, "project-root", ".", "Project root directory")
	bundleInstallCmd.Flags().StringVar(&bundleOutput, "output", "opencode.json", "Output file path")
	bundleInstallCmd.Flags().BoolVar(&bundleAuto, "auto", false, "Run in non-interactive mode (requires source-ref and --preset)")
	bundleInstallCmd.Flags().BoolVar(&bundleForce, "force", false, "Overwrite existing files")
	bundleInstallCmd.Flags().BoolVar(&bundleDryRun, "dry-run", false, "Show what would be done without doing it")
	bundleInstallCmd.ValidArgsFunction = completeSourceRefs
	_ = bundleInstallCmd.RegisterFlagCompletionFunc("preset", completeBundlePresetNames)
	bundleUpdateCmd.ValidArgsFunction = completeSourceRefs

	// Flags for bundle status
	bundleStatusCmd.Flags().StringVar(&bundleProjectRoot, "project-root", ".", "Project root directory")

	// Flags for bundle update
	bundleUpdateCmd.Flags().BoolVar(&bundleYes, "yes", false, "Skip confirmation prompt")
}

func runBundleInstall(sourceRef string, interactivePreset bool) error {
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
	selectedVersion := bundleVersion
	if string(src.Type) == "github-release" {
		selectedVersion, err = resolveGitHubBundleVersion(src.Location, bundleVersion, true)
		if err != nil {
			return err
		}
	}
	if bundleVersion != "" && string(src.Type) != "github-release" {
		return fmt.Errorf("--version is only supported for github-release sources")
	}

	// Resolve source to local bundle root
	bundleRoot, cleanup, err := bundleResolveToLocal(string(src.Type), src.Location, selectedVersion)
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
		if bundleAuto || (!interactivePreset && !bundleInputIsTTY()) {
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
		fmt.Printf("dry-run: install preset '%s' from bundle '%s'\n", selectedPreset, manifest.BundleName)
		fmt.Printf("dry-run: write config to %s\n", outputPath)
		return nil
	}

	// Reuse the shared write semantics so bundle install matches init/preset overwrite behavior.
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
	fmt.Println("done: bundle installed")

	return nil
}

func resolveGitHubBundleVersion(sourceLocation, requestedVersion string, allowPrompt bool) (string, error) {
	if requestedVersion != "" {
		return requestedVersion, nil
	}

	ref, err := source.ParseGitHubLocation(sourceLocation)
	if err != nil {
		return "", err
	}
	if ref.Tag != "" {
		return ref.Tag, nil
	}

	releases, err := bundleListGitHubReleases(sourceLocation)
	if err != nil {
		return "", err
	}
	if !allowPrompt || bundleAuto || !bundleInputIsTTY() {
		if hasStableGitHubRelease(releases) {
			return "", fmt.Errorf("--version is required for github-release sources outside interactive mode (use --version latest or --version <tag>)")
		}
		return "", fmt.Errorf("--version is required for github-release sources outside interactive mode; only prereleases are available (use --version <tag>)")
	}

	if len(releases) == 1 {
		return releases[0].TagName, nil
	}

	return promptForGitHubReleaseSelection(sourceLocation, releases)
}

func inspectGitHubBundleVersion(sourceLocation, requestedVersion string) (string, error) {
	if requestedVersion != "" {
		return requestedVersion, nil
	}

	ref, err := source.ParseGitHubLocation(sourceLocation)
	if err != nil {
		return "", err
	}
	if ref.Tag != "" {
		return ref.Tag, nil
	}

	releases, err := bundleListGitHubReleases(sourceLocation)
	if err != nil {
		return "", err
	}
	for _, release := range releases {
		if !release.Prerelease {
			return release.TagName, nil
		}
	}

	return "", fmt.Errorf("no stable release found for %s; prereleases are available", ref.Repo)
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

func promptForSourceSelection() (string, error) {
	sources, err := source.ListSources()
	if err != nil {
		return "", fmt.Errorf("failed to list sources: %w", err)
	}

	if len(sources) == 0 {
		return "", fmt.Errorf("no sources registered (run 'oc bundle source add <location>' to register a source)")
	}

	reader := bufio.NewReader(bundlePromptIn)
	for {
		fmt.Fprintln(bundlePromptOut, "Available sources:")
		for i, src := range sources {
			fmt.Fprintf(bundlePromptOut, "  %d) %s (%s) - %s\n", i+1, src.ID, src.Type, src.Name)
		}
		fmt.Fprint(bundlePromptOut, "Select a source: ")

		selection, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				return "", fmt.Errorf("interactive source selection cancelled")
			}
			return "", fmt.Errorf("failed to read source selection: %w", err)
		}

		selection = strings.TrimSpace(selection)

		// First check for exact name match
		for _, src := range sources {
			if src.Name == selection {
				return src.ID, nil
			}
		}

		// Then check for exact ID match
		for _, src := range sources {
			if src.ID == selection {
				return src.ID, nil
			}
		}

		// Finally check for index
		if index, err := strconv.Atoi(selection); err == nil {
			if index >= 1 && index <= len(sources) {
				return sources[index-1].ID, nil
			}
		}

		fmt.Fprintln(bundlePromptOut, "Invalid selection. Please enter a source number, ID, or name.")
	}
}

func promptForGitHubReleaseSelection(sourceLocation string, releases []bundle.GitHubReleaseVersion) (string, error) {
	if len(releases) == 0 {
		return "", fmt.Errorf("github-release source has no versions to select")
	}

	reader := bufio.NewReader(bundlePromptIn)
	for {
		fmt.Fprintf(bundlePromptOut, "Available versions for %s:\n", sourceLocation)
		for i, release := range releases {
			label := release.TagName
			if release.Prerelease {
				label += " (prerelease)"
			}
			fmt.Fprintf(bundlePromptOut, "  %d) %s\n", i+1, label)
		}
		fmt.Fprint(bundlePromptOut, "Select a version: ")

		selection, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				return "", fmt.Errorf("interactive version selection cancelled")
			}
			return "", fmt.Errorf("failed to read version selection: %w", err)
		}

		selection = strings.TrimSpace(selection)
		for _, release := range releases {
			if release.TagName == selection {
				return release.TagName, nil
			}
		}

		if index, err := strconv.Atoi(selection); err == nil {
			if index >= 1 && index <= len(releases) {
				return releases[index-1].TagName, nil
			}
		}

		fmt.Fprintln(bundlePromptOut, "Invalid selection. Please enter a version number or exact tag.")
	}
}

func hasStableGitHubRelease(releases []bundle.GitHubReleaseVersion) bool {
	for _, release := range releases {
		if !release.Prerelease {
			return true
		}
	}
	return false
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
	if string(src.Type) == "github-release" {
		versionTag, err = inspectGitHubBundleVersion(src.Location, versionTag)
		if err != nil {
			return nil, cobra.ShellCompDirectiveNoFileComp
		}
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
		return fmt.Errorf("no bundle installed in this project (run 'bundle install' first)")
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
		return fmt.Errorf("no bundle installed in this project (run 'bundle install' first)")
	}

	// Suppress unused variable warning
	_ = prov

	// For now, return not implemented for github-release
	return fmt.Errorf("bundle update for github-release sources requires network operations (not yet implemented)")
}
