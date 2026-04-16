package cmd

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/qbicsoftware/occo/internal/bundle"
	"github.com/qbicsoftware/occo/internal/source"
	"github.com/spf13/cobra"
)

// setupTestProject creates a temporary project directory
func setupTestProject(t *testing.T) string {
	t.Helper()
	tempDir, err := os.MkdirTemp("", "opencode-test-project-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}

	// Create .opencode directory
	opencodeDir := filepath.Join(tempDir, ".opencode")
	if err := os.MkdirAll(opencodeDir, 0755); err != nil {
		os.RemoveAll(tempDir)
		t.Fatalf("failed to create .opencode directory: %v", err)
	}

	return tempDir
}

// TestBundleApplyNoSource tests applying bundle without a source
func TestBundleApplyNoSource(t *testing.T) {
	// Save original flag values
	origPreset := bundlePreset
	origProjectRoot := bundleProjectRoot
	origForce := bundleForce
	origDryRun := bundleDryRun
	origOutput := bundleOutput
	defer func() {
		bundlePreset = origPreset
		bundleProjectRoot = origProjectRoot
		bundleForce = origForce
		bundleDryRun = origDryRun
		bundleOutput = origOutput
		bundleAuto = false
	}()

	// Test with nonexistent source
	bundlePreset = "test"
	bundleProjectRoot = "."
	bundleDryRun = false

	err := runBundleInstall("nonexistent-id", false)
	if err == nil {
		t.Error("runBundleInstall() expected error for nonexistent source")
	}
}

// TestBundleApplyMissingPreset tests applying with missing preset flag
func TestBundleApplyMissingPreset(t *testing.T) {
	restore := saveRegistry(t)
	defer restore()

	bundleDir := t.TempDir()
	manifest := `{"manifest_version":"1.0.0","bundle_name":"local","bundle_version":"v1.0.0","presets":[{"name":"test","entrypoint":"test.json","description":"Test preset"}]}`
	if err := os.WriteFile(filepath.Join(bundleDir, "opencode-bundle.manifest.json"), []byte(manifest), 0644); err != nil {
		t.Fatalf("failed to write manifest: %v", err)
	}
	if err := os.WriteFile(filepath.Join(bundleDir, "test.json"), []byte(`{"agents":[]}`), 0644); err != nil {
		t.Fatalf("failed to write preset: %v", err)
	}

	registry, _ := source.LoadRegistry()
	registry.Sources = []source.Source{{ID: "abc12345", Name: "qbic", Type: source.SourceTypeLocalDirectory, Location: bundleDir}}
	if err := source.SaveRegistry(registry); err != nil {
		t.Fatalf("failed to save registry: %v", err)
	}

	origPreset := bundlePreset
	origAuto := bundleAuto
	origTTY := bundleInputIsTTY
	origProjectRoot := bundleProjectRoot
	defer func() { bundlePreset = origPreset }()
	defer func() {
		bundleAuto = origAuto
		bundleInputIsTTY = origTTY
		bundleProjectRoot = origProjectRoot
	}()

	bundlePreset = ""
	bundleAuto = false
	bundleInputIsTTY = func() bool { return false }
	bundleProjectRoot = t.TempDir()

	err := runBundleInstall("abc12345", false)
	if err == nil {
		t.Error("runBundleInstall() expected error when preset is missing")
	}
	if !strings.Contains(err.Error(), "--preset is required outside interactive mode") {
		t.Fatalf("runBundleInstall() error = %v", err)
	}
}

// TestBundleApplyAdjustsPromptFilePaths tests that CLI adjusts prompt paths
func TestBundleApplyAdjustsPromptFilePaths(t *testing.T) {
	restore := saveRegistry(t)
	defer restore()

	bundleDir := t.TempDir()

	manifest := `{"manifest_version":"1.0.0","bundle_name":"test","bundle_version":"v1.0.0","presets":[{"name":"test","entrypoint":"test.json","description":"Test","prompt_files":["prompts/coder.md"]}]}`
	if err := os.WriteFile(filepath.Join(bundleDir, "opencode-bundle.manifest.json"), []byte(manifest), 0644); err != nil {
		t.Fatalf("failed to write manifest: %v", err)
	}

	preset := `{"agents":{"coder":{"prompt":"{file:./prompts/coder.md}"}}}`
	if err := os.WriteFile(filepath.Join(bundleDir, "test.json"), []byte(preset), 0644); err != nil {
		t.Fatalf("failed to write preset: %v", err)
	}

	if err := os.MkdirAll(filepath.Join(bundleDir, "prompts"), 0755); err != nil {
		t.Fatalf("failed to create prompts dir: %v", err)
	}
	if err := os.WriteFile(filepath.Join(bundleDir, "prompts", "coder.md"), []byte("# coder"), 0644); err != nil {
		t.Fatalf("failed to write coder.md: %v", err)
	}

	reg, err := source.LoadRegistry()
	if err != nil {
		t.Fatalf("failed to load registry: %v", err)
	}
	reg.Sources = []source.Source{{ID: "tid", Name: "t", Type: source.SourceTypeLocalDirectory, Location: bundleDir}}
	if err := source.SaveRegistry(reg); err != nil {
		t.Fatalf("failed to save registry: %v", err)
	}

	origPreset := bundlePreset
	origAuto := bundleAuto
	origTTY := bundleInputIsTTY
	origRoot := bundleProjectRoot
	origForce := bundleForce
	origAssets := bundleInstallAssets
	defer func() {
		bundlePreset = origPreset
		bundleAuto = origAuto
		bundleInputIsTTY = origTTY
		bundleProjectRoot = origRoot
		bundleForce = origForce
		bundleInstallAssets = origAssets
	}()

	bundlePreset = "test"
	bundleAuto = true
	bundleInputIsTTY = func() bool { return false }
	bundleProjectRoot = t.TempDir()
	bundleForce = true
	bundleInstallAssets = true

	err = runBundleInstall("tid", false)
	if err != nil {
		t.Fatalf("runBundleInstall() error = %v", err)
	}

	cfg, err := os.ReadFile(filepath.Join(bundleProjectRoot, "opencode.json"))
	if err != nil {
		t.Fatalf("failed to read config: %v", err)
	}
	if strings.Contains(string(cfg), "./prompts/") {
		t.Errorf("config should not contain ./prompts/, got: %s", cfg)
	}
	if !strings.Contains(string(cfg), ".opencode/prompts/") {
		t.Errorf("config should contain .opencode/prompts/, got: %s", cfg)
	}
	if _, err := os.Stat(filepath.Join(bundleProjectRoot, ".opencode", "prompts", "coder.md")); err != nil {
		t.Errorf("prompt should be installed to .opencode/prompts/: %v", err)
	}
}

// TestBundleStatusNoProvenance tests status command with no provenance
func TestBundleStatusNoProvenance(t *testing.T) {
	origProjectRoot := bundleProjectRoot
	defer func() { bundleProjectRoot = origProjectRoot }()

	// Use a temp directory with no provenance
	tempDir, err := os.MkdirTemp("", "opencode-test-noprovenance-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	bundleProjectRoot = tempDir

	err = runBundleStatus()
	if err == nil {
		t.Error("runBundleStatus() expected error when no provenance exists")
	}
}

// TestBundleStatusWithProvenance tests status command with provenance
func TestBundleStatusWithProvenance(t *testing.T) {
	origProjectRoot := bundleProjectRoot
	defer func() { bundleProjectRoot = origProjectRoot }()

	// Create temp project with provenance
	tempDir := setupTestProject(t)
	defer os.RemoveAll(tempDir)

	prov := &bundle.Provenance{
		SourceID:      "test-id",
		SourceName:    "test-source",
		SourceType:    "local-directory",
		BundleVersion: "v1.0.0",
		PresetName:    "test",
		Entrypoint:    "test.json",
		AppliedAt:     "2026-03-31T00:00:00Z",
	}
	if err := bundle.SaveProvenance(tempDir, prov, false); err != nil {
		t.Fatalf("failed to save provenance: %v", err)
	}

	bundleProjectRoot = tempDir

	err := runBundleStatus()
	if err != nil {
		t.Errorf("runBundleStatus() error = %v", err)
	}
}

// TestBundleUpdateNonGitHub tests update command with non-github source
func TestBundleUpdateNonGitHub(t *testing.T) {
	// This test requires a source in the registry
	// For now just verify it returns error for non-github source

	err := runBundleUpdate("nonexistent")
	if err == nil {
		t.Error("runBundleUpdate() expected error for nonexistent source")
	}
}

func TestBundleApplyPassesVersionForGitHubSources(t *testing.T) {
	restore := saveRegistry(t)
	defer restore()

	registry, _ := source.LoadRegistry()
	registry.Sources = []source.Source{{
		ID:       "github1",
		Location: "qbicsoftware/opencode-config-bundle",
		Type:     source.SourceTypeGitHubRelease,
		Name:     "qbic",
	}}
	if err := source.SaveRegistry(registry); err != nil {
		t.Fatalf("failed to save registry: %v", err)
	}

	projectRoot := setupTestProject(t)
	defer os.RemoveAll(projectRoot)

	origPreset := bundlePreset
	origProjectRoot := bundleProjectRoot
	origVersion := bundleVersion
	origResolver := bundleResolveToLocal
	defer func() {
		bundlePreset = origPreset
		bundleProjectRoot = origProjectRoot
		bundleVersion = origVersion
		bundleResolveToLocal = origResolver
	}()

	bundlePreset = "test"
	bundleProjectRoot = projectRoot
	bundleVersion = "v1.2.3"
	bundleResolveToLocal = func(sourceType, sourceLocation, versionTag string) (string, func(), error) {
		if sourceType != "github-release" {
			t.Fatalf("sourceType = %q, want github-release", sourceType)
		}
		if sourceLocation != "qbicsoftware/opencode-config-bundle" {
			t.Fatalf("sourceLocation = %q", sourceLocation)
		}
		if versionTag != "v1.2.3" {
			t.Fatalf("versionTag = %q, want v1.2.3", versionTag)
		}

		bundleRoot := t.TempDir()
		manifest := `{"manifest_version":"1.0.0","bundle_name":"qbic","bundle_version":"v1.2.3","presets":[{"name":"test","entrypoint":"test.json"}]}`
		if err := os.WriteFile(filepath.Join(bundleRoot, "opencode-bundle.manifest.json"), []byte(manifest), 0644); err != nil {
			return "", nil, err
		}
		if err := os.WriteFile(filepath.Join(bundleRoot, "test.json"), []byte(`{"agents":[]}`), 0644); err != nil {
			return "", nil, err
		}
		return bundleRoot, func() {}, nil
	}

	if err := runBundleInstall("github1", false); err != nil {
		t.Fatalf("runBundleInstall() error = %v", err)
	}
	if _, err := os.Stat(filepath.Join(projectRoot, "opencode.json")); err != nil {
		t.Fatalf("expected opencode.json to be written: %v", err)
	}
}

func TestBundleApplyInteractiveSelectsGitHubReleaseVersion(t *testing.T) {
	restore := saveRegistry(t)
	defer restore()

	registry, _ := source.LoadRegistry()
	registry.Sources = []source.Source{{
		ID:       "github1",
		Location: "qbicsoftware/opencode-config-bundle",
		Type:     source.SourceTypeGitHubRelease,
		Name:     "qbic",
	}}
	if err := source.SaveRegistry(registry); err != nil {
		t.Fatalf("failed to save registry: %v", err)
	}

	projectRoot := setupTestProject(t)
	defer os.RemoveAll(projectRoot)

	origPreset := bundlePreset
	origProjectRoot := bundleProjectRoot
	origVersion := bundleVersion
	origResolver := bundleResolveToLocal
	origListReleases := bundleListGitHubReleases
	origTTY := bundleInputIsTTY
	origPromptIn := bundlePromptIn
	origPromptOut := bundlePromptOut
	defer func() {
		bundlePreset = origPreset
		bundleProjectRoot = origProjectRoot
		bundleVersion = origVersion
		bundleResolveToLocal = origResolver
		bundleListGitHubReleases = origListReleases
		bundleInputIsTTY = origTTY
		bundlePromptIn = origPromptIn
		bundlePromptOut = origPromptOut
	}()

	bundlePreset = "test"
	bundleProjectRoot = projectRoot
	bundleVersion = ""
	bundleInputIsTTY = func() bool { return true }
	bundlePromptIn = strings.NewReader("2\n")
	bundlePromptOut = io.Discard
	bundleListGitHubReleases = func(location string) ([]bundle.GitHubReleaseVersion, error) {
		if location != "qbicsoftware/opencode-config-bundle" {
			t.Fatalf("location = %q", location)
		}
		return []bundle.GitHubReleaseVersion{{TagName: "v1.3.0", Prerelease: false}, {TagName: "v1.4.0-alpha.1", Prerelease: true}}, nil
	}
	bundleResolveToLocal = func(sourceType, sourceLocation, versionTag string) (string, func(), error) {
		if versionTag != "v1.4.0-alpha.1" {
			t.Fatalf("versionTag = %q, want prerelease selection", versionTag)
		}
		bundleRoot := t.TempDir()
		manifest := `{"manifest_version":"1.0.0","bundle_name":"qbic","bundle_version":"v1.4.0-alpha.1","presets":[{"name":"test","entrypoint":"test.json"}]}`
		if err := os.WriteFile(filepath.Join(bundleRoot, "opencode-bundle.manifest.json"), []byte(manifest), 0644); err != nil {
			return "", nil, err
		}
		if err := os.WriteFile(filepath.Join(bundleRoot, "test.json"), []byte(`{"agents":[]}`), 0644); err != nil {
			return "", nil, err
		}
		return bundleRoot, func() {}, nil
	}

	if err := runBundleInstall("github1", false); err != nil {
		t.Fatalf("runBundleInstall() error = %v", err)
	}
}

func TestBundleApplyGitHubSourceRequiresVersionOutsideInteractiveMode(t *testing.T) {
	restore := saveRegistry(t)
	defer restore()

	registry, _ := source.LoadRegistry()
	registry.Sources = []source.Source{{
		ID:       "github1",
		Location: "qbicsoftware/opencode-config-bundle",
		Type:     source.SourceTypeGitHubRelease,
		Name:     "qbic",
	}}
	if err := source.SaveRegistry(registry); err != nil {
		t.Fatalf("failed to save registry: %v", err)
	}

	origPreset := bundlePreset
	origVersion := bundleVersion
	origListReleases := bundleListGitHubReleases
	origTTY := bundleInputIsTTY
	defer func() {
		bundlePreset = origPreset
		bundleVersion = origVersion
		bundleListGitHubReleases = origListReleases
		bundleInputIsTTY = origTTY
	}()

	bundlePreset = "test"
	bundleVersion = ""
	bundleInputIsTTY = func() bool { return false }
	bundleListGitHubReleases = func(string) ([]bundle.GitHubReleaseVersion, error) {
		return []bundle.GitHubReleaseVersion{{TagName: "v1.3.0", Prerelease: false}, {TagName: "v1.4.0-alpha.1", Prerelease: true}}, nil
	}

	err := runBundleInstall("github1", false)
	if err == nil {
		t.Fatal("runBundleInstall() error = nil, want version-selection error")
	}
	if !strings.Contains(err.Error(), "--version is required for github-release sources outside interactive mode") {
		t.Fatalf("runBundleInstall() error = %v", err)
	}
}

func TestBundleApplyGitHubSourceReportsPrereleaseOnlyOutsideInteractiveMode(t *testing.T) {
	restore := saveRegistry(t)
	defer restore()

	registry, _ := source.LoadRegistry()
	registry.Sources = []source.Source{{
		ID:       "github1",
		Location: "qbicsoftware/opencode-config-bundle",
		Type:     source.SourceTypeGitHubRelease,
		Name:     "qbic",
	}}
	if err := source.SaveRegistry(registry); err != nil {
		t.Fatalf("failed to save registry: %v", err)
	}

	origPreset := bundlePreset
	origVersion := bundleVersion
	origListReleases := bundleListGitHubReleases
	origTTY := bundleInputIsTTY
	defer func() {
		bundlePreset = origPreset
		bundleVersion = origVersion
		bundleListGitHubReleases = origListReleases
		bundleInputIsTTY = origTTY
	}()

	bundlePreset = "test"
	bundleVersion = ""
	bundleInputIsTTY = func() bool { return false }
	bundleListGitHubReleases = func(string) ([]bundle.GitHubReleaseVersion, error) {
		return []bundle.GitHubReleaseVersion{{TagName: "v1.4.0-alpha.1", Prerelease: true}, {TagName: "v1.3.0-alpha.2", Prerelease: true}}, nil
	}

	err := runBundleInstall("github1", false)
	if err == nil {
		t.Fatal("runBundleInstall() error = nil, want prerelease-only version-selection error")
	}
	if !strings.Contains(err.Error(), "only prereleases are available") {
		t.Fatalf("runBundleInstall() error = %v", err)
	}
}

func TestBundleApplyGitHubSourceSinglePrereleaseStillRequiresVersionOutsideInteractiveMode(t *testing.T) {
	restore := saveRegistry(t)
	defer restore()

	registry, _ := source.LoadRegistry()
	registry.Sources = []source.Source{{
		ID:       "github1",
		Location: "qbicsoftware/opencode-config-bundle",
		Type:     source.SourceTypeGitHubRelease,
		Name:     "qbic",
	}}
	if err := source.SaveRegistry(registry); err != nil {
		t.Fatalf("failed to save registry: %v", err)
	}

	origPreset := bundlePreset
	origVersion := bundleVersion
	origListReleases := bundleListGitHubReleases
	origTTY := bundleInputIsTTY
	defer func() {
		bundlePreset = origPreset
		bundleVersion = origVersion
		bundleListGitHubReleases = origListReleases
		bundleInputIsTTY = origTTY
	}()

	bundlePreset = "test"
	bundleVersion = ""
	bundleInputIsTTY = func() bool { return false }
	bundleListGitHubReleases = func(string) ([]bundle.GitHubReleaseVersion, error) {
		return []bundle.GitHubReleaseVersion{{TagName: "v1.4.0-alpha.1", Prerelease: true}}, nil
	}

	err := runBundleInstall("github1", false)
	if err == nil {
		t.Fatal("runBundleInstall() error = nil, want prerelease-only version-selection error")
	}
	if !strings.Contains(err.Error(), "only prereleases are available") {
		t.Fatalf("runBundleInstall() error = %v", err)
	}
}

func TestCompleteBundlePresetNamesGitHubSourceUsesNonInteractiveInspection(t *testing.T) {
	restore := saveRegistry(t)
	defer restore()

	registry, _ := source.LoadRegistry()
	registry.Sources = []source.Source{{ID: "github1", Name: "qbic", Type: source.SourceTypeGitHubRelease, Location: "qbicsoftware/opencode-config-bundle"}}
	if err := source.SaveRegistry(registry); err != nil {
		t.Fatalf("failed to save registry: %v", err)
	}

	origResolver := bundleResolveToLocal
	origListReleases := bundleListGitHubReleases
	origTTY := bundleInputIsTTY
	origPromptOut := bundlePromptOut
	defer func() {
		bundleResolveToLocal = origResolver
		bundleListGitHubReleases = origListReleases
		bundleInputIsTTY = origTTY
		bundlePromptOut = origPromptOut
	}()

	bundleInputIsTTY = func() bool { return true }
	var promptOutput bytes.Buffer
	bundlePromptOut = &promptOutput
	bundleListGitHubReleases = func(string) ([]bundle.GitHubReleaseVersion, error) {
		return []bundle.GitHubReleaseVersion{{TagName: "v2.0.0-alpha.1", Prerelease: true}, {TagName: "v1.9.0", Prerelease: false}}, nil
	}
	bundleResolveToLocal = func(sourceType, sourceLocation, versionTag string) (string, func(), error) {
		if versionTag != "v1.9.0" {
			t.Fatalf("versionTag = %q, want latest stable", versionTag)
		}
		bundleRoot := t.TempDir()
		manifest := `{"manifest_version":"1.0.0","bundle_name":"qbic","bundle_version":"v1.9.0","presets":[{"name":"test","entrypoint":"test.json"}]}`
		if err := os.WriteFile(filepath.Join(bundleRoot, "opencode-bundle.manifest.json"), []byte(manifest), 0644); err != nil {
			return "", nil, err
		}
		if err := os.WriteFile(filepath.Join(bundleRoot, "test.json"), []byte(`{"agents":[]}`), 0644); err != nil {
			return "", nil, err
		}
		return bundleRoot, func() {}, nil
	}

	cmd := &cobra.Command{}
	cmd.Flags().String("version", "", "")
	completions, directive := completeBundlePresetNames(cmd, []string{"qbic"}, "t")
	if directive != cobra.ShellCompDirectiveNoFileComp {
		t.Fatalf("directive = %v", directive)
	}
	if promptOutput.Len() != 0 {
		t.Fatalf("completion should not prompt, got %q", promptOutput.String())
	}
	if len(completions) != 1 || completions[0] != "test" {
		t.Fatalf("completions = %v", completions)
	}
}

// TestBundleInstallFlags tests that bundle install flags are properly configured
func TestBundleInstallFlags(t *testing.T) {
	if bundleInstallCmd.Flags().Lookup("preset") == nil {
		t.Error("preset flag should exist on bundle apply command")
	}
	if bundleInstallCmd.Flags().Lookup("auto") == nil {
		t.Error("auto flag should exist on bundle apply command")
	}
	if bundleInstallCmd.Flags().Lookup("project-root") == nil {
		t.Error("project-root flag should exist on bundle apply command")
	}
	if bundleInstallCmd.Flags().Lookup("force") == nil {
		t.Error("force flag should exist on bundle apply command")
	}
	if bundleInstallCmd.Flags().Lookup("dry-run") == nil {
		t.Error("dry-run flag should exist on bundle apply command")
	}
}

// TestBundleStatusFlags tests that bundle status flags are properly configured
func TestBundleStatusFlags(t *testing.T) {
	if bundleStatusCmd.Flags().Lookup("project-root") == nil {
		t.Error("project-root flag should exist on bundle status command")
	}
}

// TestBundleUpdateFlags tests that bundle update flags are properly configured
func TestBundleUpdateFlags(t *testing.T) {
	if bundleUpdateCmd.Flags().Lookup("yes") == nil {
		t.Error("yes flag should exist on bundle update command")
	}
}

func TestBundleApplyVersionFlagExists(t *testing.T) {
	if bundleInstallCmd.Flags().Lookup("version") == nil {
		t.Fatal("version flag should exist on bundle apply command")
	}
}

func TestBundleApplyRejectsVersionForLocalSources(t *testing.T) {
	restore := saveRegistry(t)
	defer restore()

	bundleDir := t.TempDir()
	manifest := `{"manifest_version":"1.0.0","bundle_name":"local","bundle_version":"v1.0.0","presets":[{"name":"test","entrypoint":"test.json"}]}`
	if err := os.WriteFile(filepath.Join(bundleDir, "opencode-bundle.manifest.json"), []byte(manifest), 0644); err != nil {
		t.Fatalf("failed to write manifest: %v", err)
	}
	if err := os.WriteFile(filepath.Join(bundleDir, "test.json"), []byte(`{"agents":[]}`), 0644); err != nil {
		t.Fatalf("failed to write preset: %v", err)
	}

	registry, _ := source.LoadRegistry()
	registry.Sources = []source.Source{{
		ID:       "local1",
		Location: bundleDir,
		Type:     source.SourceTypeLocalDirectory,
		Name:     "local",
	}}
	if err := source.SaveRegistry(registry); err != nil {
		t.Fatalf("failed to save registry: %v", err)
	}

	projectRoot := setupTestProject(t)
	defer os.RemoveAll(projectRoot)

	origPreset := bundlePreset
	origProjectRoot := bundleProjectRoot
	origVersion := bundleVersion
	defer func() {
		bundlePreset = origPreset
		bundleProjectRoot = origProjectRoot
		bundleVersion = origVersion
	}()

	bundlePreset = "test"
	bundleProjectRoot = projectRoot
	bundleVersion = "v1.2.3"

	err := runBundleInstall("local1", false)
	if err == nil {
		t.Fatal("runBundleInstall() error = nil, want error")
	}
	if !strings.Contains(err.Error(), "--version is only supported for github-release sources") {
		t.Fatalf("runBundleInstall() error = %v", err)
	}
}

func TestBundleApplyResolvesSourceByName(t *testing.T) {
	restore := saveRegistry(t)
	defer restore()

	bundleDir := t.TempDir()
	manifest := `{"manifest_version":"1.0.0","bundle_name":"local","bundle_version":"v1.0.0","presets":[{"name":"test","entrypoint":"test.json","description":"Test preset"}]}`
	if err := os.WriteFile(filepath.Join(bundleDir, "opencode-bundle.manifest.json"), []byte(manifest), 0644); err != nil {
		t.Fatalf("failed to write manifest: %v", err)
	}
	if err := os.WriteFile(filepath.Join(bundleDir, "test.json"), []byte(`{"agents":[]}`), 0644); err != nil {
		t.Fatalf("failed to write preset: %v", err)
	}

	registry, _ := source.LoadRegistry()
	registry.Sources = []source.Source{{
		ID:       "local1",
		Location: bundleDir,
		Type:     source.SourceTypeLocalDirectory,
		Name:     "qbic",
	}}
	if err := source.SaveRegistry(registry); err != nil {
		t.Fatalf("failed to save registry: %v", err)
	}

	projectRoot := setupTestProject(t)
	defer os.RemoveAll(projectRoot)

	origPreset := bundlePreset
	origProjectRoot := bundleProjectRoot
	defer func() {
		bundlePreset = origPreset
		bundleProjectRoot = origProjectRoot
	}()

	bundlePreset = "test"
	bundleProjectRoot = projectRoot

	if err := runBundleInstall("qbic", false); err != nil {
		t.Fatalf("runBundleInstall() error = %v", err)
	}

	prov, err := bundle.LoadProvenance(projectRoot)
	if err != nil {
		t.Fatalf("LoadProvenance() error = %v", err)
	}
	if prov.SourceID != "local1" {
		t.Fatalf("provenance SourceID = %q, want local1", prov.SourceID)
	}
}

func TestBundleApplyRejectsAmbiguousSourceName(t *testing.T) {
	restore := saveRegistry(t)
	defer restore()

	registry, _ := source.LoadRegistry()
	registry.Sources = []source.Source{{ID: "id-1", Name: "qbic", Type: source.SourceTypeLocalDirectory, Location: "/tmp/a"}, {ID: "id-2", Name: "qbic", Type: source.SourceTypeLocalDirectory, Location: "/tmp/b"}}
	if err := source.SaveRegistry(registry); err != nil {
		t.Fatalf("failed to save registry: %v", err)
	}

	origPreset := bundlePreset
	bundlePreset = "test"
	defer func() { bundlePreset = origPreset }()

	err := runBundleInstall("qbic", false)
	if err == nil {
		t.Fatal("runBundleInstall() error = nil, want ambiguous source error")
	}
	if !strings.Contains(err.Error(), "ambiguous") || !strings.Contains(err.Error(), "id-1") || !strings.Contains(err.Error(), "id-2") {
		t.Fatalf("runBundleInstall() error = %v", err)
	}
}

func TestBundleApplyInteractiveSelectsPreset(t *testing.T) {
	restore := saveRegistry(t)
	defer restore()

	bundleDir := t.TempDir()
	manifest := `{"manifest_version":"1.0.0","bundle_name":"local","bundle_version":"v1.0.0","presets":[{"name":"first","entrypoint":"first.json","description":"First preset"},{"name":"second","entrypoint":"second.json","description":"Second preset"}]}`
	if err := os.WriteFile(filepath.Join(bundleDir, "opencode-bundle.manifest.json"), []byte(manifest), 0644); err != nil {
		t.Fatalf("failed to write manifest: %v", err)
	}
	if err := os.WriteFile(filepath.Join(bundleDir, "first.json"), []byte(`{"name":"first"}`), 0644); err != nil {
		t.Fatalf("failed to write first preset: %v", err)
	}
	if err := os.WriteFile(filepath.Join(bundleDir, "second.json"), []byte(`{"name":"second"}`), 0644); err != nil {
		t.Fatalf("failed to write second preset: %v", err)
	}

	registry, _ := source.LoadRegistry()
	registry.Sources = []source.Source{{ID: "local1", Name: "qbic", Type: source.SourceTypeLocalDirectory, Location: bundleDir}}
	if err := source.SaveRegistry(registry); err != nil {
		t.Fatalf("failed to save registry: %v", err)
	}

	projectRoot := setupTestProject(t)
	defer os.RemoveAll(projectRoot)

	origPreset := bundlePreset
	origAuto := bundleAuto
	origTTY := bundleInputIsTTY
	origPromptIn := bundlePromptIn
	origPromptOut := bundlePromptOut
	defer func() {
		bundlePreset = origPreset
		bundleAuto = origAuto
		bundleInputIsTTY = origTTY
		bundlePromptIn = origPromptIn
		bundlePromptOut = origPromptOut
	}()

	bundlePreset = ""
	bundleAuto = false
	bundleProjectRoot = projectRoot
	bundleInputIsTTY = func() bool { return true }
	bundlePromptIn = strings.NewReader("2\n")
	bundlePromptOut = io.Discard

	if err := runBundleInstall("qbic", false); err != nil {
		t.Fatalf("runBundleInstall() error = %v", err)
	}

	content, err := os.ReadFile(filepath.Join(projectRoot, "opencode.json"))
	if err != nil {
		t.Fatalf("failed to read written config: %v", err)
	}
	if string(content) != `{"name":"second"}` {
		t.Fatalf("written config = %s", content)
	}
}

func TestBundleApplyInteractiveAcceptsNumericLikePresetName(t *testing.T) {
	manifest := &bundle.Manifest{
		BundleName: "numeric-fixture",
		Presets: []bundle.Preset{
			{Name: "first", Description: "First preset"},
			{Name: "2", Description: "Numeric-like preset"},
		},
	}

	origPromptIn := bundlePromptIn
	origPromptOut := bundlePromptOut
	defer func() {
		bundlePromptIn = origPromptIn
		bundlePromptOut = origPromptOut
	}()

	bundlePromptIn = strings.NewReader("2\n")
	bundlePromptOut = io.Discard

	selected, err := promptForPresetSelection(manifest)
	if err != nil {
		t.Fatalf("promptForPresetSelection() error = %v", err)
	}
	if selected != "2" {
		t.Fatalf("selected preset = %q, want %q", selected, "2")
	}
}

func TestCompleteSourceRefs(t *testing.T) {
	restore := saveRegistry(t)
	defer restore()

	registry, _ := source.LoadRegistry()
	registry.Sources = []source.Source{{ID: "id-1", Name: "qbic", Type: source.SourceTypeLocalDirectory, Location: "/tmp/a"}}
	if err := source.SaveRegistry(registry); err != nil {
		t.Fatalf("failed to save registry: %v", err)
	}

	completions, directive := completeSourceRefs(nil, nil, "q")
	if directive != 4 { // cobra.ShellCompDirectiveNoFileComp
		t.Fatalf("directive = %v", directive)
	}
	if len(completions) != 1 || completions[0] != "qbic" {
		t.Fatalf("completions = %v", completions)
	}
}

// ============================================================================
// Bundle Init Tests
// ============================================================================

func TestBundleInitFlags(t *testing.T) {
	if bundleInitCmd.Flags().Lookup("name") == nil {
		t.Error("name flag should exist on bundle init command")
	}
	if bundleInitCmd.Flags().Lookup("version") == nil {
		t.Error("version flag should exist on bundle init command")
	}
	if bundleInitCmd.Flags().Lookup("output") == nil {
		t.Error("output flag should exist on bundle init command")
	}
	if bundleInitCmd.Flags().Lookup("force") == nil {
		t.Error("force flag should exist on bundle init command")
	}
}

func TestBundleInitWithNameAndVersion(t *testing.T) {
	// Save and restore original values
	origName := bundleInitName
	origVersion := bundleInitVersion
	origOutput := bundleInitOutput
	origForce := bundleInitForce
	origTTY := bundleInputIsTTY
	defer func() {
		bundleInitName = origName
		bundleInitVersion = origVersion
		bundleInitOutput = origOutput
		bundleInitForce = origForce
		bundleInputIsTTY = origTTY
	}()

	// Create temp output directory
	tempDir := t.TempDir()
	outputDir := filepath.Join(tempDir, "my-bundle")

	bundleInitName = "my-bundle"
	bundleInitVersion = "v1.0.0"
	bundleInitOutput = outputDir
	bundleInitForce = false
	bundleInputIsTTY = func() bool { return false }

	err := runBundleInit()
	if err != nil {
		t.Fatalf("runBundleInit() error = %v", err)
	}

	// Verify manifest was created
	manifestPath := filepath.Join(outputDir, "opencode-bundle.manifest.json")
	if _, err := os.Stat(manifestPath); err != nil {
		t.Fatalf("expected manifest to be created: %v", err)
	}

	// Verify manifest content
	manifest, err := bundle.LoadManifest(manifestPath)
	if err != nil {
		t.Fatalf("failed to load manifest: %v", err)
	}
	if manifest.BundleName != "my-bundle" {
		t.Fatalf("manifest.BundleName = %q, want %q", manifest.BundleName, "my-bundle")
	}
	if manifest.BundleVersion != "v1.0.0" {
		t.Fatalf("manifest.BundleVersion = %q, want %q", manifest.BundleVersion, "v1.0.0")
	}
	if manifest.ManifestVersion != "1.0.0" {
		t.Fatalf("manifest.ManifestVersion = %q, want %q", manifest.ManifestVersion, "1.0.0")
	}
	if len(manifest.Presets) != 1 {
		t.Fatalf("len(manifest.Presets) = %d, want 1", len(manifest.Presets))
	}
	if manifest.Presets[0].Name != "default" {
		t.Fatalf("manifest.Presets[0].Name = %q, want %q", manifest.Presets[0].Name, "default")
	}
	if manifest.Presets[0].Entrypoint != "default.json" {
		t.Fatalf("manifest.Presets[0].Entrypoint = %q, want %q", manifest.Presets[0].Entrypoint, "default.json")
	}

	// Verify preset was created
	presetPath := filepath.Join(outputDir, "default.json")
	if _, err := os.Stat(presetPath); err != nil {
		t.Fatalf("expected preset to be created: %v", err)
	}

	// Verify README was created
	readmePath := filepath.Join(outputDir, "README.md")
	if _, err := os.Stat(readmePath); err != nil {
		t.Fatalf("expected README to be created: %v", err)
	}

	readmeData, err := os.ReadFile(readmePath)
	if err != nil {
		t.Fatalf("failed to read README: %v", err)
	}
	if !strings.Contains(string(readmeData), "my-bundle") {
		t.Fatalf("README should contain bundle name")
	}
}

func TestBundleInitWithNameOnlyDefaultsVersion(t *testing.T) {
	// Save and restore original values
	origName := bundleInitName
	origVersion := bundleInitVersion
	origOutput := bundleInitOutput
	origTTY := bundleInputIsTTY
	defer func() {
		bundleInitName = origName
		bundleInitVersion = origVersion
		bundleInitOutput = origOutput
		bundleInputIsTTY = origTTY
	}()

	// Create temp output directory
	tempDir := t.TempDir()
	outputDir := filepath.Join(tempDir, "test-bundle")

	bundleInitName = "test-bundle"
	bundleInitVersion = ""
	bundleInitOutput = outputDir
	bundleInputIsTTY = func() bool { return false }

	err := runBundleInit()
	if err != nil {
		t.Fatalf("runBundleInit() error = %v", err)
	}

	// Verify version defaults to 0.0.1
	manifestPath := filepath.Join(outputDir, "opencode-bundle.manifest.json")
	manifest, err := bundle.LoadManifest(manifestPath)
	if err != nil {
		t.Fatalf("failed to load manifest: %v", err)
	}
	if manifest.BundleVersion != defaultBundleVersion {
		t.Fatalf("manifest.BundleVersion = %q, want default %q", manifest.BundleVersion, defaultBundleVersion)
	}
}

func TestBundleInitFailsWithoutNameInNonInteractiveMode(t *testing.T) {
	// Save and restore original values
	origName := bundleInitName
	origOutput := bundleInitOutput
	origTTY := bundleInputIsTTY
	defer func() {
		bundleInitName = origName
		bundleInitOutput = origOutput
		bundleInputIsTTY = origTTY
	}()

	tempDir := t.TempDir()
	bundleInitName = ""
	bundleInitOutput = filepath.Join(tempDir, "output")
	bundleInputIsTTY = func() bool { return false }

	err := runBundleInit()
	if err == nil {
		t.Fatal("runBundleInit() expected error when name is empty in non-interactive mode")
	}
	if !strings.Contains(err.Error(), "--name is required") {
		t.Fatalf("runBundleInit() error = %v, want error containing '--name is required'", err)
	}
}

func TestBundleInitInteractivePromptsForName(t *testing.T) {
	// Save and restore original values
	origName := bundleInitName
	origVersion := bundleInitVersion
	origOutput := bundleInitOutput
	origTTY := bundleInputIsTTY
	origPromptIn := bundlePromptIn
	origPromptOut := bundlePromptOut
	defer func() {
		bundleInitName = origName
		bundleInitVersion = origVersion
		bundleInitOutput = origOutput
		bundleInputIsTTY = origTTY
		bundlePromptIn = origPromptIn
		bundlePromptOut = origPromptOut
	}()

	tempDir := t.TempDir()
	outputDir := filepath.Join(tempDir, "interactive-bundle")
	bundleInitName = ""
	bundleInitVersion = "v2.0.0"
	bundleInitOutput = outputDir
	bundleInputIsTTY = func() bool { return true }
	bundlePromptIn = strings.NewReader("interactive-name\n")
	bundlePromptOut = io.Discard

	err := runBundleInit()
	if err != nil {
		t.Fatalf("runBundleInit() error = %v", err)
	}

	// Verify manifest was created with the interactive name
	manifestPath := filepath.Join(outputDir, "opencode-bundle.manifest.json")
	manifest, err := bundle.LoadManifest(manifestPath)
	if err != nil {
		t.Fatalf("failed to load manifest: %v", err)
	}
	if manifest.BundleName != "interactive-name" {
		t.Fatalf("manifest.BundleName = %q, want %q", manifest.BundleName, "interactive-name")
	}
	if manifest.BundleVersion != "v2.0.0" {
		t.Fatalf("manifest.BundleVersion = %q, want %q", manifest.BundleVersion, "v2.0.0")
	}
}

func TestBundleInitInteractivePromptsForVersion(t *testing.T) {
	// Save and restore original values
	origName := bundleInitName
	origVersion := bundleInitVersion
	origOutput := bundleInitOutput
	origTTY := bundleInputIsTTY
	origPromptIn := bundlePromptIn
	origPromptOut := bundlePromptOut
	defer func() {
		bundleInitName = origName
		bundleInitVersion = origVersion
		bundleInitOutput = origOutput
		bundleInputIsTTY = origTTY
		bundlePromptIn = origPromptIn
		bundlePromptOut = origPromptOut
	}()

	tempDir := t.TempDir()
	outputDir := filepath.Join(tempDir, "version-prompt-bundle")
	bundleInitName = "test-bundle"
	bundleInitVersion = ""
	bundleInitOutput = outputDir
	bundleInputIsTTY = func() bool { return true }
	bundlePromptIn = strings.NewReader("\n") // Empty version, should use default
	bundlePromptOut = io.Discard

	err := runBundleInit()
	if err != nil {
		t.Fatalf("runBundleInit() error = %v", err)
	}

	// Verify manifest was created with default version
	manifestPath := filepath.Join(outputDir, "opencode-bundle.manifest.json")
	manifest, err := bundle.LoadManifest(manifestPath)
	if err != nil {
		t.Fatalf("failed to load manifest: %v", err)
	}
	if manifest.BundleVersion != defaultBundleVersion {
		t.Fatalf("manifest.BundleVersion = %q, want default %q", manifest.BundleVersion, defaultBundleVersion)
	}
}

func TestBundleInitFailsWhenDirectoryExistsWithoutForce(t *testing.T) {
	// Save and restore original values
	origName := bundleInitName
	origOutput := bundleInitOutput
	origTTY := bundleInputIsTTY
	defer func() {
		bundleInitName = origName
		bundleInitOutput = origOutput
		bundleInputIsTTY = origTTY
	}()

	// Create temp directory with existing content
	tempDir := t.TempDir()
	outputDir := filepath.Join(tempDir, "existing-dir")
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		t.Fatalf("failed to create directory: %v", err)
	}
	// Write something in the directory
	if err := os.WriteFile(filepath.Join(outputDir, "existing.txt"), []byte("existing"), 0644); err != nil {
		t.Fatalf("failed to write file: %v", err)
	}

	bundleInitName = "new-bundle"
	bundleInitOutput = outputDir
	bundleInitForce = false
	bundleInputIsTTY = func() bool { return false }

	err := runBundleInit()
	if err == nil {
		t.Fatal("runBundleInit() expected error when directory exists without --force")
	}
	if !strings.Contains(err.Error(), "already exists") {
		t.Fatalf("runBundleInit() error = %v, want error containing 'already exists'", err)
	}
}

func TestBundleInitSucceedsWithForceOnExistingDirectory(t *testing.T) {
	// Save and restore original values
	origName := bundleInitName
	origOutput := bundleInitOutput
	origTTY := bundleInputIsTTY
	defer func() {
		bundleInitName = origName
		bundleInitOutput = origOutput
		bundleInputIsTTY = origTTY
	}()

	// Create temp directory with existing content
	tempDir := t.TempDir()
	outputDir := filepath.Join(tempDir, "existing-dir")
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		t.Fatalf("failed to create directory: %v", err)
	}
	// Write something in the directory
	if err := os.WriteFile(filepath.Join(outputDir, "existing.txt"), []byte("existing"), 0644); err != nil {
		t.Fatalf("failed to write file: %v", err)
	}

	bundleInitName = "new-bundle"
	bundleInitOutput = outputDir
	bundleInitForce = true
	bundleInputIsTTY = func() bool { return false }

	err := runBundleInit()
	if err != nil {
		t.Fatalf("runBundleInit() error = %v", err)
	}

	// Verify manifest was created
	manifestPath := filepath.Join(outputDir, "opencode-bundle.manifest.json")
	if _, err := os.Stat(manifestPath); err != nil {
		t.Fatalf("expected manifest to be created: %v", err)
	}
}

func TestBundleInitCreatesPresetsSubdirectory(t *testing.T) {
	// Save and restore original values
	origName := bundleInitName
	origOutput := bundleInitOutput
	origTTY := bundleInputIsTTY
	defer func() {
		bundleInitName = origName
		bundleInitOutput = origOutput
		bundleInputIsTTY = origTTY
	}()

	tempDir := t.TempDir()
	outputDir := filepath.Join(tempDir, "preset-test-bundle")

	bundleInitName = "preset-test"
	bundleInitOutput = outputDir
	bundleInputIsTTY = func() bool { return false }

	err := runBundleInit()
	if err != nil {
		t.Fatalf("runBundleInit() error = %v", err)
	}

	// Verify preset file was created at root level (per bundle contract)
	presetPath := filepath.Join(outputDir, "default.json")
	if _, err := os.Stat(presetPath); err != nil {
		t.Fatalf("preset file not created at root level: %v", err)
	}
}

func TestValidateBundleName(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		{"valid-bundle", false},
		{"valid_bundle", false},
		{"validBundle123", false},
		{"a", false},
		{"123", false},
		{"-invalid", true},     // can't start with hyphen
		{"_invalid", true},     // can't start with underscore
		{"", true},             // empty
		{"invalid name", true}, // contains space
		{"invalid.name", true}, // contains dot
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateBundleName(tt.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateBundleName(%q) error = %v, wantErr %v", tt.name, err, tt.wantErr)
			}
		})
	}
}

func TestBundleInitRejectsInvalidBundleName(t *testing.T) {
	// Save and restore original values
	origName := bundleInitName
	origOutput := bundleInitOutput
	origTTY := bundleInputIsTTY
	defer func() {
		bundleInitName = origName
		bundleInitOutput = origOutput
		bundleInputIsTTY = origTTY
	}()

	tempDir := t.TempDir()
	bundleInitName = "-invalid" // starts with hyphen (invalid)
	bundleInitOutput = filepath.Join(tempDir, "output")
	bundleInputIsTTY = func() bool { return false }

	err := runBundleInit()
	if err == nil {
		t.Fatal("runBundleInit() expected error for invalid bundle name")
	}
	if !strings.Contains(err.Error(), "must start with a letter or number") {
		t.Fatalf("runBundleInit() error = %v, want error containing 'must start with a letter or number'", err)
	}
}

func TestBundleInitInteractiveCancelsOnEOF(t *testing.T) {
	// Save and restore original values
	origName := bundleInitName
	origOutput := bundleInitOutput
	origTTY := bundleInputIsTTY
	origPromptIn := bundlePromptIn
	defer func() {
		bundleInitName = origName
		bundleInitOutput = origOutput
		bundleInputIsTTY = origTTY
		bundlePromptIn = origPromptIn
	}()

	tempDir := t.TempDir()
	bundleInitName = ""
	bundleInitOutput = filepath.Join(tempDir, "output")
	bundleInputIsTTY = func() bool { return true }
	bundlePromptIn = strings.NewReader("") // EOF

	err := runBundleInit()
	if err == nil {
		t.Fatal("runBundleInit() expected error on EOF")
	}
	if !strings.Contains(err.Error(), "cancelled") {
		t.Fatalf("runBundleInit() error = %v, want error containing 'cancelled'", err)
	}
}

func TestBundleInitCreatesOutputDirectory(t *testing.T) {
	// Save and restore original values
	origName := bundleInitName
	origOutput := bundleInitOutput
	origTTY := bundleInputIsTTY
	defer func() {
		bundleInitName = origName
		bundleInitOutput = origOutput
		bundleInputIsTTY = origTTY
	}()

	tempDir := t.TempDir()
	outputDir := filepath.Join(tempDir, "new", "nested", "path")

	bundleInitName = "nested-bundle"
	bundleInitOutput = outputDir
	bundleInputIsTTY = func() bool { return false }

	err := runBundleInit()
	if err != nil {
		t.Fatalf("runBundleInit() error = %v", err)
	}

	// Verify directory was created
	if info, err := os.Stat(outputDir); err != nil {
		t.Fatalf("output directory not created: %v", err)
	} else if !info.IsDir() {
		t.Fatal("output is not a directory")
	}
}

func TestBundleInitManifestIsLoadable(t *testing.T) {
	// Save and restore original values
	origName := bundleInitName
	origVersion := bundleInitVersion
	origOutput := bundleInitOutput
	origTTY := bundleInputIsTTY
	defer func() {
		bundleInitName = origName
		bundleInitVersion = origVersion
		bundleInitOutput = origOutput
		bundleInputIsTTY = origTTY
	}()

	tempDir := t.TempDir()
	outputDir := filepath.Join(tempDir, "loadable-bundle")

	bundleInitName = "loadable-bundle"
	bundleInitVersion = "v3.0.0"
	bundleInitOutput = outputDir
	bundleInputIsTTY = func() bool { return false }

	err := runBundleInit()
	if err != nil {
		t.Fatalf("runBundleInit() error = %v", err)
	}

	// Verify manifest can be loaded by LoadManifest
	manifestPath := filepath.Join(outputDir, "opencode-bundle.manifest.json")
	manifest, err := bundle.LoadManifest(manifestPath)
	if err != nil {
		t.Fatalf("LoadManifest() failed: %v", err)
	}

	// Verify manifest has all required fields
	if manifest.ManifestVersion != "1.0.0" {
		t.Fatalf("manifest.ManifestVersion = %q", manifest.ManifestVersion)
	}
	if manifest.BundleName != "loadable-bundle" {
		t.Fatalf("manifest.BundleName = %q", manifest.BundleName)
	}
	if manifest.BundleVersion != "v3.0.0" {
		t.Fatalf("manifest.BundleVersion = %q", manifest.BundleVersion)
	}
	if len(manifest.Presets) != 1 {
		t.Fatalf("len(manifest.Presets) = %d", len(manifest.Presets))
	}
}

func TestBundleInitPresetFilePathMatchesEntrypoint(t *testing.T) {
	// Save and restore original values
	origName := bundleInitName
	origOutput := bundleInitOutput
	origTTY := bundleInputIsTTY
	defer func() {
		bundleInitName = origName
		bundleInitOutput = origOutput
		bundleInputIsTTY = origTTY
	}()

	tempDir := t.TempDir()
	outputDir := filepath.Join(tempDir, "entrypoint-test")

	bundleInitName = "entrypoint-test"
	bundleInitOutput = outputDir
	bundleInputIsTTY = func() bool { return false }

	err := runBundleInit()
	if err != nil {
		t.Fatalf("runBundleInit() error = %v", err)
	}

	// Load manifest and verify preset entrypoint
	manifestPath := filepath.Join(outputDir, "opencode-bundle.manifest.json")
	manifest, err := bundle.LoadManifest(manifestPath)
	if err != nil {
		t.Fatalf("LoadManifest() failed: %v", err)
	}

	preset := manifest.Presets[0]
	expectedPresetPath := filepath.Join(outputDir, preset.Entrypoint)

	// Verify the preset file exists at the entrypoint path
	if _, err := os.Stat(expectedPresetPath); err != nil {
		t.Fatalf("preset file does not exist at entrypoint path %s: %v", expectedPresetPath, err)
	}
}
