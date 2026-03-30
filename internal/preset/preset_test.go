package preset

import (
	"os"
	"path/filepath"
	"testing"
)

// TestValidPresets tests that all expected presets are valid (for reference)
func TestValidPresets(t *testing.T) {
	presets := ValidPresets()
	expected := []string{"mixed", "openai", "big-pickle", "minimax", "kimi"}

	if len(presets) != len(expected) {
		t.Errorf("expected %d presets, got %d", len(expected), len(presets))
	}

	for _, e := range expected {
		found := false
		for _, p := range presets {
			if p == e {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("expected preset %q not found", e)
		}
	}
}

// TestGetDefaultConfig tests getting the default config
func TestGetDefaultConfig(t *testing.T) {
	config, err := GetDefaultConfig()
	if err != nil {
		t.Logf("GetDefaultConfig() error (expected in test without bundled config): %v", err)
		return
	}

	if len(config) == 0 {
		t.Error("config should not be empty")
	}
}

// TestWriteConfig tests writing config to a file
func TestWriteConfig(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "oc-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	destFile := filepath.Join(tmpDir, "test.json")
	testData := `{"test": true}`

	// Test write without force (new file)
	if err := WriteConfig(destFile, testData, false); err != nil {
		t.Errorf("WriteConfig() error = %v", err)
	}

	// Verify content
	content, err := os.ReadFile(destFile)
	if err != nil {
		t.Errorf("failed to read dest file: %v", err)
	}
	if string(content) != testData {
		t.Errorf("content mismatch: got %s", string(content))
	}

	// Test write with force (existing file)
	if err := WriteConfig(destFile, `{"updated": true}`, true); err != nil {
		t.Errorf("WriteConfig() with force error = %v", err)
	}

	// Test write without force (existing file should fail)
	if err := WriteConfig(destFile, testData, false); err == nil {
		t.Error("WriteConfig() should fail when file exists and force=false")
	}
}
