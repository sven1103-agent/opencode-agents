// Package preset provides functionality for handling OpenCode presets.
package preset

import (
	"fmt"
	"os"
	"path/filepath"
)

// ValidPresets returns the list of valid preset names (for reference).
// Note: Preset selection will be handled via bundle/source commands (US-045).
func ValidPresets() []string {
	return []string{
		"mixed",
		"openai",
		"big-pickle",
		"minimax",
		"kimi",
	}
}

// GetDefaultConfig returns the default bundled config.
// The config is read from a bundled file that comes with the installation.
func GetDefaultConfig() (string, error) {
	// First try: read from bundled configs directory (relative to executable)
	execPath, err := os.Executable()
	if err == nil {
		execDir := filepath.Dir(execPath)
		bundledPath := filepath.Join(execDir, "presets", "default.json")
		if data, err := os.ReadFile(bundledPath); err == nil {
			return string(data), nil
		}
	}

	// Second try: development mode (check repo root)
	possiblePaths := []string{
		"opencode.mixed.json",
		"opencode.openai.json",
		"opencode.json",
	}
	for _, path := range possiblePaths {
		if data, err := os.ReadFile(path); err == nil {
			return string(data), nil
		}
	}

	return "", fmt.Errorf("no bundled config found - please ensure installation is complete")
}

// WriteConfig writes the config data to the destination path.
func WriteConfig(destPath string, data string, force bool) error {
	// Check if destination exists
	if !force {
		if _, err := os.Stat(destPath); err == nil {
			return fmt.Errorf("output file exists: %s (use --force to overwrite)", destPath)
		}
	}

	// Ensure destination directory exists
	destDir := filepath.Dir(destPath)
	if err := os.MkdirAll(destDir, 0755); err != nil {
		return fmt.Errorf("failed to create directory %s: %w", destDir, err)
	}

	// Write the file
	if err := os.WriteFile(destPath, []byte(data), 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}
