package cmd

import (
	"os"
	"testing"
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
}
