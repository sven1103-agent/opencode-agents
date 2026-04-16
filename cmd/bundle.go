package cmd

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/qbicsoftware/occo/internal/bundle"
	configpreset "github.com/qbicsoftware/occo/internal/preset"
	"github.com/qbicsoftware/occo/internal/source"
	"github.com/qbicsoftware/occo/internal/styles"
	"github.com/spf13/cobra"
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
	bundleInstallAssets                = true
	bundleResolveToLocal               = bundle.ResolveToLocal
	bundleListGitHubReleases           = bundle.ListGitHubReleases
	bundlePromptIn           io.Reader = os.Stdin
	bundlePromptOut          io.Writer = os.Stdout
	bundleInputIsTTY                   = isInteractiveTTY

	// bundle init flags
	bundleInitName    string
	bundleInitVersion string
	bundleInitOutput  string
	bundleInitForce   bool
)

// Bundle name validation: alphanumeric, hyphens, and underscores only
// Must start with alphanumeric, hyphens/underscores allowed after first character
var bundleNameRegex = regexp.MustCompile(`^[a-zA-Z0-9][a-zA-Z0-9_-]*$`)

// Default version for new bundles
const defaultBundleVersion = "0.0.1"

// bundleCmd represents the bundle command
var bundleCmd = &cobra.Command{
	Use:   "bundle",
	Short: "Manage OpenCode configuration bundles",
	Long: `Manage OpenCode configuration bundles.

Install, track, and update configuration bundles from registered sources.

Examples:
	  occo bundle install qbic --preset default
	  occo bundle install qbic
	  occo bundle status
	  occo bundle update abc12345`,
}

// bundleInstallCmd installs a preset from a registered config bundle
var bundleInstallCmd = &cobra.Command{
	Use:   "install [source-ref]",
	Short: "Install a preset from a config bundle",
	Long: `Install a preset from a registered config bundle to a project.

The source-ref may be either a registered source ID or a unique source name.
When omitted in interactive mode, the command prompts for source and preset selection.

Examples:
	  occo bundle install qbic --preset default
	  occo bundle install qbic
	  occo bundle install abc12345 --version v1.2.3 --preset default
	  occo bundle install qbic --preset minimal --project-root ./myproject
	  occo bundle install qbic --auto --preset default --force
	  occo bundle install (interactive mode)`,
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
  occo bundle status
  occo bundle status --project-root ./myproject`,
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
  occo bundle update abc12345
  occo bundle update abc12345 --yes`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return runBundleUpdate(args[0])
	},
}

// bundleInitCmd initializes a new bundle directory with a proper structure
var bundleInitCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize a new configuration bundle",
	Long: `Initialize a new configuration bundle with the proper directory structure.

Creates the following files in the output directory:
  - opencode-bundle.manifest.json: Bundle manifest with metadata
  - presets/default.json: Default preset placeholder
  - README.md: Bundle documentation

Examples:
  occo bundle init
  occo bundle init --name mybundle
  occo bundle init --name mybundle --version v1.0.0
  occo bundle init --name mybundle --output ./my-bundle`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return runBundleInit()
	},
}

func init() {
	rootCmd.AddCommand(bundleCmd)

	// Add subcommands
	bundleCmd.AddCommand(bundleInstallCmd)
	bundleCmd.AddCommand(bundleStatusCmd)
	bundleCmd.AddCommand(bundleUpdateCmd)
	bundleCmd.AddCommand(bundleInitCmd)

	// Flags for bundle apply
	bundleInstallCmd.Flags().StringVar(&bundlePreset, "preset", "", "Preset name to apply")
	bundleInstallCmd.Flags().StringVar(&bundleVersion, "version", "", "Bundle version/tag to apply for github-release sources")
	bundleInstallCmd.Flags().StringVar(&bundleProjectRoot, "project-root", ".", "Project root directory")
	bundleInstallCmd.Flags().StringVar(&bundleOutput, "output", "opencode.json", "Output file path")
	bundleInstallCmd.Flags().BoolVar(&bundleAuto, "auto", false, "Run in non-interactive mode (requires source-ref and --preset)")
	bundleInstallCmd.Flags().BoolVar(&bundleForce, "force", false, "Overwrite existing files")
	bundleInstallCmd.Flags().BoolVar(&bundleDryRun, "dry-run", false, "Show what would be done without doing it")
	bundleInstallCmd.Flags().BoolVar(&bundleInstallAssets, "assets", true, "Install prompt files to .opencode/ directory")
	bundleInstallCmd.ValidArgsFunction = completeSourceRefs
	_ = bundleInstallCmd.RegisterFlagCompletionFunc("preset", completeBundlePresetNames)
	bundleUpdateCmd.ValidArgsFunction = completeSourceRefs

	// Flags for bundle status
	bundleStatusCmd.Flags().StringVar(&bundleProjectRoot, "project-root", ".", "Project root directory")

	// Flags for bundle update
	bundleUpdateCmd.Flags().BoolVar(&bundleYes, "yes", false, "Skip confirmation prompt")

	// Flags for bundle init
	bundleInitCmd.Flags().StringVar(&bundleInitName, "name", "", "Bundle name (required in non-interactive mode)")
	bundleInitCmd.Flags().StringVar(&bundleInitVersion, "version", "", "Bundle version (defaults to 0.0.1)")
	bundleInitCmd.Flags().StringVar(&bundleInitOutput, "output", ".", "Output directory for the bundle")
	bundleInitCmd.Flags().BoolVar(&bundleInitForce, "force", false, "Overwrite existing directory contents")
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
		fmt.Println()
		fmt.Println(styles.SectionHeader("Dry Run"))
		fmt.Println(styles.KeyValue("Preset", selectedPreset))
		fmt.Println(styles.KeyValue("Bundle", manifest.BundleName))
		fmt.Println(styles.KeyValue("Output", outputPath))
		fmt.Println()
		fmt.Println(styles.Done("dry-run complete"))
		return nil
	}

	// Check for existing files and prompt for override if needed
	if !bundleForce {
		if filesToOverwrite, err := checkExistingFiles(projectRoot, outputPath, bundlePresetEntry.PromptFiles); err != nil {
			return err
		} else if len(filesToOverwrite) > 0 {
			proceed, err := promptForOverrideConfirmation(filesToOverwrite)
			if err != nil {
				return err
			}
			if !proceed {
				return fmt.Errorf("installation cancelled by user")
			}
			bundleForce = true
		}
	}

	// Adjust prompt file paths in preset content from ./prompts/ to .opencode/prompts/
	// The CLI normalizes paths to the target installation location (local project)
	if bundleInstallAssets && len(bundlePresetEntry.PromptFiles) > 0 {
		adjustedContent := strings.ReplaceAll(string(presetContent), "./prompts/", ".opencode/prompts/")
		presetContent = []byte(adjustedContent)
	}

	// Reuse the shared write semantics so bundle apply matches init/preset overwrite behavior.
	if err := configpreset.WriteConfig(outputPath, string(presetContent), bundleForce); err != nil {
		return fmt.Errorf("failed to write config: %w", err)
	}
	fmt.Println(styles.Written(outputPath))

	// Install prompt files if enabled and present
	var installedAssets []bundle.InstalledAsset
	if bundleInstallAssets && len(bundlePresetEntry.PromptFiles) > 0 {
		installed, err := installPromptFiles(bundleRoot, projectRoot, bundlePresetEntry.PromptFiles, bundleForce)
		if err != nil {
			return fmt.Errorf("failed to install prompt files: %w", err)
		}
		installedAssets = installed

		for _, a := range installedAssets {
			fmt.Printf("written: %s\n", a.Destination)
		}
	}

	// Write provenance
	prov := &bundle.Provenance{
		SourceID:        src.ID,
		SourceName:      src.Name,
		SourceType:      string(src.Type),
		BundleVersion:   manifest.BundleVersion,
		PresetName:      selectedPreset,
		Entrypoint:      bundlePresetEntry.Entrypoint,
		AppliedAt:       "2026-03-31T00:00:00Z", // Would use time.Now().Format(time.RFC3339)
		InstalledAssets: installedAssets,
	}

	if err := bundle.SaveProvenance(projectRoot, prov, bundleForce); err != nil {
		return fmt.Errorf("failed to save provenance: %w", err)
	}
	fmt.Println(styles.Written(bundle.ProvenancePath(projectRoot)))
	fmt.Println(styles.Done("bundle applied"))

	return nil
}

// installPromptFiles copies prompt files from the bundle to the project's .opencode/prompts/ directory.
func installPromptFiles(bundleRoot, projectRoot string, promptFiles []string, force bool) ([]bundle.InstalledAsset, error) {
	if len(promptFiles) == 0 {
		return nil, nil
	}

	var installed []bundle.InstalledAsset
	promptsDir := filepath.Join(projectRoot, ".opencode", "prompts")

	// Create prompts directory if it doesn't exist
	if err := os.MkdirAll(promptsDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create prompts directory: %w", err)
	}

	for _, pf := range promptFiles {
		// Normalize path: strip leading "prompts/" prefix to avoid duplication
		normalizedPath := strings.TrimPrefix(pf, "prompts/")
		normalizedPath = strings.TrimPrefix(normalizedPath, "/")

		// Validate normalized path isn't empty
		if normalizedPath == "" {
			return nil, fmt.Errorf("invalid prompt path in config: %q - path cannot be empty after normalization", pf)
		}

		sourcePath := filepath.Join(bundleRoot, pf)

		// Verify source file exists
		if _, err := os.Stat(sourcePath); err != nil {
			return nil, fmt.Errorf("prompt file not found in bundle: %s", pf)
		}

		// Determine destination (preserve relative structure, but strip leading "prompts/" prefix)
		destPath := filepath.Join(promptsDir, filepath.Base(normalizedPath))

		// Create subdirectory if needed (for paths like subdir/file.md)
		if strings.Contains(normalizedPath, string(filepath.Separator)) {
			subdir := filepath.Dir(normalizedPath)
			subdirPath := filepath.Join(promptsDir, subdir)
			if err := os.MkdirAll(subdirPath, 0755); err != nil {
				return nil, fmt.Errorf("failed to create subdirectory: %w", err)
			}
			destPath = filepath.Join(promptsDir, normalizedPath)
		}

		// Check if destination exists (unless force)
		if !force {
			if _, err := os.Stat(destPath); err == nil {
				return nil, fmt.Errorf("prompt file already exists: %s (use --force to overwrite)", destPath)
			}
		}

		// Copy the file
		sourceData, err := os.ReadFile(sourcePath)
		if err != nil {
			return nil, fmt.Errorf("failed to read prompt file: %w", err)
		}

		if err := os.WriteFile(destPath, sourceData, 0644); err != nil {
			return nil, fmt.Errorf("failed to write prompt file: %w", err)
		}

		installed = append(installed, bundle.InstalledAsset{
			Source:      pf,
			Destination: destPath,
		})
	}

	return installed, nil
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

	fmt.Print(styles.Loading("Fetching available versions from GitHub..."))
	releases, err := bundleListGitHubReleases(sourceLocation)
	fmt.Print("\r")
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
	var latestPrerelease string
	for _, release := range releases {
		if !release.Prerelease {
			return release.TagName, nil
		}
		if latestPrerelease == "" || release.TagName > latestPrerelease {
			latestPrerelease = release.TagName
		}
	}

	if latestPrerelease != "" {
		return latestPrerelease, nil
	}

	return "", fmt.Errorf("no releases found for %s", ref.Repo)
}

func checkExistingFiles(projectRoot, outputPath string, promptFiles []string) ([]string, error) {
	var existing []string

	if _, err := os.Stat(outputPath); err == nil {
		existing = append(existing, outputPath)
	}

	provPath := bundle.ProvenancePath(projectRoot)
	if _, err := os.Stat(provPath); err == nil {
		existing = append(existing, provPath)
	}

	if len(promptFiles) > 0 {
		promptsDir := filepath.Join(projectRoot, ".opencode", "prompts")
		for _, pf := range promptFiles {
			destPath := filepath.Join(promptsDir, filepath.Base(pf))
			if _, err := os.Stat(destPath); err == nil {
				existing = append(existing, destPath)
			}
		}
	}

	return existing, nil
}

func promptForOverrideConfirmation(files []string) (bool, error) {
	fmt.Fprintln(bundlePromptOut)
	fmt.Fprintln(bundlePromptOut, styles.Warning("The following files will be overwritten:"))
	for _, f := range files {
		fmt.Fprintf(bundlePromptOut, "  - %s\n", styles.Muted(f))
	}
	fmt.Fprintln(bundlePromptOut)
	fmt.Fprint(bundlePromptOut, styles.YesNoPrompt("Do you want to overwrite", "n"))

	reader := bufio.NewReader(bundlePromptIn)
	selection, err := reader.ReadString('\n')
	if err != nil {
		if err == io.EOF {
			return false, nil
		}
		return false, fmt.Errorf("failed to read confirmation: %w", err)
	}

	selection = strings.TrimSpace(strings.ToLower(selection))
	if selection == "y" || selection == "yes" {
		return true, nil
	}
	return false, nil
}

func promptForPresetSelection(manifest *bundle.Manifest) (string, error) {
	if len(manifest.Presets) == 0 {
		return "", fmt.Errorf("bundle has no presets to select")
	}

	defaultIdx := 0
	for i, p := range manifest.Presets {
		if p.Name == "default" {
			defaultIdx = i
			break
		}
	}

	reader := bufio.NewReader(bundlePromptIn)
	for {
		fmt.Fprintln(bundlePromptOut)
		fmt.Fprint(bundlePromptOut, styles.SectionHeader("Select Preset for "+manifest.BundleName))
		for i, preset := range manifest.Presets {
			assetIndicator := styles.AssetIndicator(len(preset.PromptFiles) > 0, len(preset.PromptFiles))
			if i == defaultIdx {
				if preset.Description != "" {
					fmt.Fprintf(bundlePromptOut, "▸ %d) %s  %s  %s\n",
						i+1,
						preset.Name,
						styles.Muted("- "+preset.Description),
						assetIndicator)
					continue
				}
				fmt.Fprintf(bundlePromptOut, "▸ %d) %s  %s\n", i+1, preset.Name, assetIndicator)
			} else {
				if preset.Description != "" {
					fmt.Fprintf(bundlePromptOut, "  %d) %s  %s  %s\n",
						i+1,
						preset.Name,
						styles.Muted("- "+preset.Description),
						assetIndicator)
					continue
				}
				fmt.Fprintf(bundlePromptOut, "  %d) %s  %s\n", i+1, preset.Name, assetIndicator)
			}
		}
		fmt.Fprintln(bundlePromptOut)
		fmt.Fprintf(bundlePromptOut, "%s(%s)%s ", styles.Prompt("Enter selection (default: "), styles.Highlight(fmt.Sprint(defaultIdx+1)), styles.Prompt("): "))

		selection, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				return "", fmt.Errorf("interactive preset selection cancelled")
			}
			return "", fmt.Errorf("failed to read preset selection: %w", err)
		}

		selection = strings.TrimSpace(selection)
		if selection == "" {
			return manifest.Presets[defaultIdx].Name, nil
		}
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

		fmt.Fprintln(bundlePromptOut, styles.Invalid("Please enter a preset number or exact name, or press Enter for default."))
	}
}

func promptForSourceSelection() (string, error) {
	sources, err := source.ListSources()
	if err != nil {
		return "", fmt.Errorf("failed to list sources: %w", err)
	}

	if len(sources) == 0 {
		return "", fmt.Errorf("no sources registered (run 'occo source add <location>' to register a source)")
	}

	reader := bufio.NewReader(bundlePromptIn)
	for {
		fmt.Fprintln(bundlePromptOut)
		fmt.Fprintln(bundlePromptOut, styles.SectionHeader("Select Source"))
		for i, src := range sources {
			icon := styles.SourceTypeIcon(string(src.Type))
			typeLabel := styles.SourceTypeLabel(string(src.Type))
			if i == 0 {
				fmt.Fprintf(bundlePromptOut, "▸ %d) %s  %s  %s\n",
					i+1,
					icon,
					styles.Highlight(src.Name),
					styles.Muted("("+typeLabel+")"))
			} else {
				fmt.Fprintf(bundlePromptOut, "  %d) %s  %s  %s\n",
					i+1,
					icon,
					styles.Highlight(src.Name),
					styles.Muted("("+typeLabel+")"))
			}
		}
		fmt.Fprintln(bundlePromptOut)
		fmt.Fprint(bundlePromptOut, styles.Prompt("Enter selection (default: 1): "))

		selection, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				return "", fmt.Errorf("interactive source selection cancelled")
			}
			return "", fmt.Errorf("failed to read source selection: %w", err)
		}

		selection = strings.TrimSpace(selection)

		if selection == "" {
			return sources[0].ID, nil
		}

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

		fmt.Fprintln(bundlePromptOut, styles.Invalid("Please enter a source number, name, or press Enter for default."))
	}
}

func promptForGitHubReleaseSelection(sourceLocation string, releases []bundle.GitHubReleaseVersion) (string, error) {
	if len(releases) == 0 {
		return "", fmt.Errorf("github-release source has no versions to select")
	}

	// GitHub API returns releases sorted by date (newest first), so index 0 is latest
	// Prefer stable over prerelease
	latestIdx := 0
	for i, r := range releases {
		if !r.Prerelease {
			// First stable release is the latest stable since sorted by date
			latestIdx = i
			break
		}
	}

	reader := bufio.NewReader(bundlePromptIn)
	for {
		fmt.Fprintln(bundlePromptOut)
		fmt.Fprint(bundlePromptOut, styles.SectionHeader("Select Version for "+sourceLocation))
		for i, release := range releases {
			label := release.TagName
			if release.Prerelease {
				label += " " + styles.Muted("(prerelease)")
			}
			if i == latestIdx {
				label += " " + styles.Muted("(recommended)")
			}
			if i == latestIdx {
				fmt.Fprintf(bundlePromptOut, "▸ %d) %s\n", i+1, label)
			} else {
				fmt.Fprintf(bundlePromptOut, "  %d) %s\n", i+1, label)
			}
		}
		fmt.Fprintln(bundlePromptOut)
		fmt.Fprintf(bundlePromptOut, "%s(%s)%s ", styles.Prompt("Enter selection (default: "), styles.Highlight(fmt.Sprint(latestIdx+1)), styles.Prompt("): "))

		selection, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				return "", fmt.Errorf("interactive version selection cancelled")
			}
			return "", fmt.Errorf("failed to read version selection: %w", err)
		}

		selection = strings.TrimSpace(selection)
		if selection == "" {
			return releases[latestIdx].TagName, nil
		}

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

		fmt.Fprintln(bundlePromptOut, styles.Invalid("Please enter a version number, exact tag, or press Enter for default."))
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
		return fmt.Errorf("no bundle applied to this project (run 'bundle apply' first)")
	}

	// Display provenance
	fmt.Println()
	fmt.Println(styles.SectionHeader("Bundle Provenance"))
	fmt.Println(styles.KeyValue("Source ID", prov.SourceID))
	fmt.Println(styles.KeyValue("Source Name", prov.SourceName))
	fmt.Println(styles.KeyValue("Source Type", string(prov.SourceType)))
	fmt.Println(styles.KeyValue("Bundle Version", prov.BundleVersion))
	fmt.Println(styles.KeyValue("Preset", prov.PresetName))
	fmt.Println(styles.KeyValueMuted("Applied At", prov.AppliedAt))

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

func runBundleInit() error {
	// Resolve output directory
	outputDir, err := filepath.Abs(bundleInitOutput)
	if err != nil {
		return fmt.Errorf("invalid output directory: %w", err)
	}

	// Determine if we're in interactive mode
	isInteractive := bundleInputIsTTY()

	// Get bundle name (interactive or from flag)
	bundleName := bundleInitName
	if bundleName == "" {
		if !isInteractive {
			return fmt.Errorf("--name is required in non-interactive mode")
		}
		bundleName, err = promptForBundleName()
		if err != nil {
			return err
		}
	}

	// Validate bundle name
	if err := validateBundleName(bundleName); err != nil {
		return err
	}

	// Get bundle version (interactive or from flag, default if neither)
	version := bundleInitVersion
	if version == "" {
		if isInteractive {
			version, err = promptForBundleVersion()
			if err != nil {
				return err
			}
		} else {
			version = defaultBundleVersion
		}
	}

	// Check if output directory exists
	info, err := os.Stat(outputDir)
	if err == nil && info.IsDir() {
		// Directory exists
		if !bundleInitForce {
			return fmt.Errorf("output directory already exists: %s (use --force to overwrite)", outputDir)
		}
	} else if os.IsNotExist(err) {
		// Directory does not exist, create it
		if err := os.MkdirAll(outputDir, 0755); err != nil {
			return fmt.Errorf("failed to create output directory: %w", err)
		}
	} else if err != nil {
		return fmt.Errorf("failed to check output directory: %w", err)
	}

	// Generate manifest
	manifest := bundle.Manifest{
		ManifestVersion: "1.0.0",
		BundleName:      bundleName,
		BundleVersion:   version,
		Presets: []bundle.Preset{
			{
				Name:        "default",
				Entrypoint:  "default.json",
				Description: "Default preset",
				PromptFiles: []string{},
			},
		},
	}

	manifestData, err := json.MarshalIndent(manifest, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal manifest: %w", err)
	}

	manifestPath := filepath.Join(outputDir, "opencode-bundle.manifest.json")
	if err := os.WriteFile(manifestPath, manifestData, 0644); err != nil {
		return fmt.Errorf("failed to write manifest: %w", err)
	}
	fmt.Println(styles.Written(manifestPath))

	// Generate default preset placeholder
	defaultPresetContent := `{
  "name": "default",
  "description": "Default preset for ` + bundleName + `",
  "agents": []
}
`
	presetPath := filepath.Join(outputDir, "default.json")
	if err := os.WriteFile(presetPath, []byte(defaultPresetContent), 0644); err != nil {
		return fmt.Errorf("failed to write preset: %w", err)
	}
	fmt.Println(styles.Written(presetPath))

	// Generate README
	readmeContent := "# " + bundleName + "\n\n" +
		"OpenCode configuration bundle.\n\n" +
		"## Bundle Info\n\n" +
		"- **Name**: " + bundleName + "\n" +
		"- **Version**: " + version + "\n\n" +
		"## Presets\n\n" +
		"### default\n\n" +
		"Default preset with basic configuration.\n\n" +
		"## Usage\n\n" +
		"```bash\n" +
		"occo bundle install <source-ref> --preset default\n" +
		"```\n\n" +
		"## Structure\n\n" +
		"- `opencode-bundle.manifest.json` - Bundle manifest\n" +
		"- `default.json` - Default preset definition\n\n" +
		"## Publishing\n\n" +
		"This bundle can be distributed via GitHub releases. See the [bundle contract](https://github.com/qbicsoftware/opencode-config-cli/blob/main/docs/specs/bundle-contract.md) for details.\n"
	readmePath := filepath.Join(outputDir, "README.md")
	if err := os.WriteFile(readmePath, []byte(readmeContent), 0644); err != nil {
		return fmt.Errorf("failed to write README: %w", err)
	}
	fmt.Println(styles.Written(readmePath))

	fmt.Println(styles.Done("bundle initialized"))
	return nil
}

// promptForBundleName prompts the user for a bundle name in interactive mode
func promptForBundleName() (string, error) {
	reader := bufio.NewReader(bundlePromptIn)
	for {
		fmt.Fprint(bundlePromptOut, styles.Prompt("Enter bundle name: "))
		name, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				return "", fmt.Errorf("interactive input cancelled")
			}
			return "", fmt.Errorf("failed to read bundle name: %w", err)
		}

		name = strings.TrimSpace(name)
		if name == "" {
			fmt.Fprintln(bundlePromptOut, styles.Error("Bundle name cannot be empty. Please enter a valid name."))
			continue
		}

		if err := validateBundleName(name); err != nil {
			fmt.Fprintf(bundlePromptOut, "%s: %s\n", styles.Error("Invalid bundle name"), err)
			continue
		}

		return name, nil
	}
}

// promptForBundleVersion prompts the user for a bundle version in interactive mode
func promptForBundleVersion() (string, error) {
	reader := bufio.NewReader(bundlePromptIn)
	fmt.Fprintf(bundlePromptOut, styles.Prompt("Enter bundle version (press Enter for default '%s'): "), defaultBundleVersion)
	version, err := reader.ReadString('\n')
	if err != nil {
		if err == io.EOF {
			return "", fmt.Errorf("interactive input cancelled")
		}
		return "", fmt.Errorf("failed to read bundle version: %w", err)
	}

	version = strings.TrimSpace(version)
	if version == "" {
		return defaultBundleVersion, nil
	}

	return version, nil
}

// validateBundleName validates the bundle name according to bundle contract rules
func validateBundleName(name string) error {
	if name == "" {
		return fmt.Errorf("bundle name cannot be empty")
	}
	if len(name) < 1 || len(name) > 64 {
		return fmt.Errorf("bundle name must be 1-64 characters")
	}
	// First character must be alphanumeric
	if !bundleNameRegex.MatchString(name[:1]) {
		return fmt.Errorf("bundle name must start with a letter or number")
	}
	// Rest of the name can be alphanumeric, hyphens, or underscores
	if len(name) > 1 && !bundleNameRegex.MatchString(name) {
		return fmt.Errorf("bundle name must contain only alphanumeric characters, hyphens, and underscores")
	}
	return nil
}
