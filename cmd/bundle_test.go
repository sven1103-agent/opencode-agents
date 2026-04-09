package cmd

import (
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/sven1103-agent/opencode-config-cli/internal/bundle"
	"github.com/sven1103-agent/opencode-config-cli/internal/source"
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

	err := runBundleApply("nonexistent-id")
	if err == nil {
		t.Error("runBundleApply() expected error for nonexistent source")
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

	err := runBundleApply("abc12345")
	if err == nil {
		t.Error("runBundleApply() expected error when preset is missing")
	}
	if !strings.Contains(err.Error(), "--preset is required outside interactive mode") {
		t.Fatalf("runBundleApply() error = %v", err)
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

	if err := runBundleApply("github1"); err != nil {
		t.Fatalf("runBundleApply() error = %v", err)
	}
	if _, err := os.Stat(filepath.Join(projectRoot, "opencode.json")); err != nil {
		t.Fatalf("expected opencode.json to be written: %v", err)
	}
}

// TestBundleApplyFlags tests that bundle apply flags are properly configured
func TestBundleApplyFlags(t *testing.T) {
	if bundleApplyCmd.Flags().Lookup("preset") == nil {
		t.Error("preset flag should exist on bundle apply command")
	}
	if bundleApplyCmd.Flags().Lookup("auto") == nil {
		t.Error("auto flag should exist on bundle apply command")
	}
	if bundleApplyCmd.Flags().Lookup("project-root") == nil {
		t.Error("project-root flag should exist on bundle apply command")
	}
	if bundleApplyCmd.Flags().Lookup("force") == nil {
		t.Error("force flag should exist on bundle apply command")
	}
	if bundleApplyCmd.Flags().Lookup("dry-run") == nil {
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
	if bundleApplyCmd.Flags().Lookup("version") == nil {
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

	err := runBundleApply("local1")
	if err == nil {
		t.Fatal("runBundleApply() error = nil, want error")
	}
	if !strings.Contains(err.Error(), "--version is only supported for github-release sources") {
		t.Fatalf("runBundleApply() error = %v", err)
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

	if err := runBundleApply("qbic"); err != nil {
		t.Fatalf("runBundleApply() error = %v", err)
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

	err := runBundleApply("qbic")
	if err == nil {
		t.Fatal("runBundleApply() error = nil, want ambiguous source error")
	}
	if !strings.Contains(err.Error(), "ambiguous") || !strings.Contains(err.Error(), "id-1") || !strings.Contains(err.Error(), "id-2") {
		t.Fatalf("runBundleApply() error = %v", err)
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

	if err := runBundleApply("qbic"); err != nil {
		t.Fatalf("runBundleApply() error = %v", err)
	}

	content, err := os.ReadFile(filepath.Join(projectRoot, "opencode.json"))
	if err != nil {
		t.Fatalf("failed to read written config: %v", err)
	}
	if string(content) != `{"name":"second"}` {
		t.Fatalf("written config = %s", content)
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
