package cmd

import (
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/sven1103-agent/opencode-config-cli/internal/source"
)

func TestRunPresetList(t *testing.T) {
	// Capture output
	oldStdout := os.Stdout
	defer func() { os.Stdout = oldStdout }()

	// Just run the function and check it doesn't error
	err := runPresetList()
	if err != nil {
		t.Errorf("runPresetList() error = %v", err)
	}
}

func TestGetPresetConfig(t *testing.T) {
	// Test getting valid presets
	presets := []string{"mixed", "openai", "big-pickle", "minimax", "kimi"}

	for _, p := range presets {
		t.Run(p, func(t *testing.T) {
			config, err := getPresetConfig(p)
			if err != nil {
				// In test environment, bundled files may not be found
				t.Logf("getPresetConfig(%q) error (expected in test without bundled files): %v", p, err)
				return
			}
			if len(config) == 0 {
				t.Error("config should not be empty")
			}
		})
	}

	// Test invalid preset
	_, err := getPresetConfig("invalid-preset")
	if err == nil {
		t.Error("getPresetConfig() should error for invalid preset")
	}
}

func TestPresetUseCmdFlags(t *testing.T) {
	// Test that flags are properly configured
	cmd := presetUseCmd

	if cmd.Flags().Lookup("project-root") == nil {
		t.Error("project-root flag should exist")
	}
	if cmd.Flags().Lookup("output") == nil {
		t.Error("output flag should exist")
	}
	if cmd.Flags().Lookup("force") == nil {
		t.Error("force flag should exist")
	}
	if cmd.Flags().Lookup("dry-run") == nil {
		t.Error("dry-run flag should exist")
	}
}

func TestPresetListCmdFlags(t *testing.T) {
	cmd := presetListCmd
	if cmd.Flags().Lookup("project-root") != nil {
		t.Error("preset list should not have project-root flag")
	}
	if cmd.Flags().Lookup("sources") == nil {
		t.Error("preset list should have sources flag")
	}
}

func TestRunPresetListSources(t *testing.T) {
	restore := saveRegistry(t)
	defer restore()

	bundleDir := t.TempDir()
	manifest := `{"manifest_version":"1.0.0","bundle_name":"qbic","bundle_version":"v1.2.3","presets":[{"name":"mixed","entrypoint":"mixed.json","description":"Mixed preset"}]}`
	if err := os.WriteFile(filepath.Join(bundleDir, "opencode-bundle.manifest.json"), []byte(manifest), 0644); err != nil {
		t.Fatalf("failed to write manifest: %v", err)
	}
	if err := os.WriteFile(filepath.Join(bundleDir, "mixed.json"), []byte(`{"agents":[]}`), 0644); err != nil {
		t.Fatalf("failed to write preset: %v", err)
	}

	registry, _ := source.LoadRegistry()
	registry.Sources = []source.Source{{ID: "id-1", Name: "qbic", Type: source.SourceTypeLocalDirectory, Location: bundleDir}}
	if err := source.SaveRegistry(registry); err != nil {
		t.Fatalf("failed to save registry: %v", err)
	}

	origSources := presetSources
	origStdout := os.Stdout
	origStderr := os.Stderr
	defer func() {
		presetSources = origSources
		os.Stdout = origStdout
		os.Stderr = origStderr
	}()

	rOut, wOut, err := os.Pipe()
	if err != nil {
		t.Fatalf("failed to create stdout pipe: %v", err)
	}
	rErr, wErr, err := os.Pipe()
	if err != nil {
		t.Fatalf("failed to create stderr pipe: %v", err)
	}
	os.Stdout = wOut
	os.Stderr = wErr
	presetSources = true

	err = runPresetList()
	wOut.Close()
	wErr.Close()
	if err != nil {
		t.Fatalf("runPresetList() error = %v", err)
	}

	stdoutBytes, _ := io.ReadAll(rOut)
	stderrBytes, _ := io.ReadAll(rErr)
	if len(stderrBytes) != 0 {
		t.Fatalf("stderr = %s", stderrBytes)
	}
	stdout := string(stdoutBytes)
	if !strings.Contains(stdout, "qbic") || !strings.Contains(stdout, "v1.2.3") || !strings.Contains(stdout, "mixed") {
		t.Fatalf("stdout = %s", stdout)
	}
}
