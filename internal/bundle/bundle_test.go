package bundle

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestLoadManifest_VersionValidation(t *testing.T) {
	tmpDir := t.TempDir()

	tests := []struct {
		name        string
		manifest    string
		wantErr     bool
		errContains string
	}{
		{
			name: "valid version 1.0.0",
			manifest: `{
				"manifest_version": "1.0.0",
				"bundle_name": "test",
				"bundle_version": "v1.0.0",
				"presets": [{"name": "p", "entrypoint": "p.json"}]
			}`,
			wantErr: false,
		},
		{
			name: "unsupported version",
			manifest: `{
				"manifest_version": "2.0.0",
				"bundle_name": "test",
				"bundle_version": "v1.0.0",
				"presets": [{"name": "p", "entrypoint": "p.json"}]
			}`,
			wantErr:     true,
			errContains: "unsupported manifest version",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			manifestPath := filepath.Join(tmpDir, "opencode-bundle.manifest.json")
			if err := os.WriteFile(manifestPath, []byte(tt.manifest), 0644); err != nil {
				t.Fatal(err)
			}

			_, err := LoadManifest(manifestPath)
			if (err != nil) != tt.wantErr {
				t.Errorf("LoadManifest() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && tt.errContains != "" {
				if err == nil || !contains(err.Error(), tt.errContains) {
					t.Errorf("LoadManifest() error = %v, want contains %v", err, tt.errContains)
				}
			}
		})
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && containsAt(s, substr))
}

func containsAt(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

// TestValidateBundle_ValidBundle tests validation of a valid bundle
func TestValidateBundle_ValidBundle(t *testing.T) {
	tmpDir := t.TempDir()

	// Create valid manifest
	manifest := `{
		"manifest_version": "1.0.0",
		"bundle_name": "test-bundle",
		"bundle_version": "v1.0.0",
		"presets": [
			{
				"name": "test",
				"description": "Test preset",
				"entrypoint": "test.json",
				"prompt_files": []
			}
		]
	}`
	manifestPath := filepath.Join(tmpDir, "opencode-bundle.manifest.json")
	if err := os.WriteFile(manifestPath, []byte(manifest), 0644); err != nil {
		t.Fatal(err)
	}

	// Create valid entrypoint
	entrypoint := `{"model": "gpt-4", "temperature": 0.7}`
	if err := os.WriteFile(filepath.Join(tmpDir, "test.json"), []byte(entrypoint), 0644); err != nil {
		t.Fatal(err)
	}

	result, err := ValidateBundle(tmpDir)
	if err != nil {
		t.Errorf("ValidateBundle() unexpected error = %v", err)
		return
	}
	if !result.Valid {
		t.Errorf("ValidateBundle() expected Valid=true, got false, errors: %v", result.Errors)
	}
	if len(result.Errors) != 0 {
		t.Errorf("ValidateBundle() expected 0 errors, got %d: %v", len(result.Errors), result.Errors)
	}
}

// TestValidateBundle_MissingManifest tests validation with missing manifest
func TestValidateBundle_MissingManifest(t *testing.T) {
	tmpDir := t.TempDir()

	result, err := ValidateBundle(tmpDir)
	if err != nil {
		t.Errorf("ValidateBundle() unexpected error = %v", err)
		return
	}
	if result.Valid {
		t.Errorf("ValidateBundle() expected Valid=false for missing manifest")
	}
	// Should have a manifest error
	found := false
	for _, e := range result.Errors {
		if e.Category == "manifest" && strings.Contains(e.Message, "manifest") {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("ValidateBundle() expected manifest error, got: %v", result.Errors)
	}
}

// TestValidateBundle_InvalidJSON tests validation with malformed JSON
func TestValidateBundle_InvalidJSON(t *testing.T) {
	tmpDir := t.TempDir()

	// Create invalid JSON manifest
	manifest := `{invalid json}`
	manifestPath := filepath.Join(tmpDir, "opencode-bundle.manifest.json")
	if err := os.WriteFile(manifestPath, []byte(manifest), 0644); err != nil {
		t.Fatal(err)
	}

	result, err := ValidateBundle(tmpDir)
	if err != nil {
		t.Errorf("ValidateBundle() unexpected error = %v", err)
		return
	}
	if result.Valid {
		t.Errorf("ValidateBundle() expected Valid=false for invalid JSON")
	}
	// Should have a schema or manifest parse error
	found := false
	for _, e := range result.Errors {
		if e.Category == "schema" || e.Category == "manifest" {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("ValidateBundle() expected schema/manifest error for invalid JSON, got: %v", result.Errors)
	}
}

// TestValidateBundle_InvalidManifestVersion tests validation with unsupported version
func TestValidateBundle_InvalidManifestVersion(t *testing.T) {
	tmpDir := t.TempDir()

	// Create manifest with unsupported version
	manifest := `{
		"manifest_version": "999.0.0",
		"bundle_name": "test-bundle",
		"bundle_version": "v1.0.0",
		"presets": [{"name": "test", "description": "Test", "entrypoint": "test.json", "prompt_files": []}]
	}`
	manifestPath := filepath.Join(tmpDir, "opencode-bundle.manifest.json")
	if err := os.WriteFile(manifestPath, []byte(manifest), 0644); err != nil {
		t.Fatal(err)
	}

	// Create dummy entrypoint to avoid entrypoint errors masking version error
	if err := os.WriteFile(filepath.Join(tmpDir, "test.json"), []byte("{}"), 0644); err != nil {
		t.Fatal(err)
	}

	result, err := ValidateBundle(tmpDir)
	if err != nil {
		t.Errorf("ValidateBundle() unexpected error = %v", err)
		return
	}
	if result.Valid {
		t.Errorf("ValidateBundle() expected Valid=false for unsupported version")
	}
	// Should have a version error
	found := false
	for _, e := range result.Errors {
		if strings.Contains(e.Message, "version") || strings.Contains(e.Message, "supported") {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("ValidateBundle() expected version error, got: %v", result.Errors)
	}
}

// TestValidateBundle_MissingRequiredFields tests validation with missing required fields
func TestValidateBundle_MissingRequiredFields(t *testing.T) {
	tmpDir := t.TempDir()

	// Create manifest missing bundle_name, bundle_version, presets
	manifest := `{
		"manifest_version": "1.0.0"
	}`
	manifestPath := filepath.Join(tmpDir, "opencode-bundle.manifest.json")
	if err := os.WriteFile(manifestPath, []byte(manifest), 0644); err != nil {
		t.Fatal(err)
	}

	result, err := ValidateBundle(tmpDir)
	if err != nil {
		t.Errorf("ValidateBundle() unexpected error = %v", err)
		return
	}
	if result.Valid {
		t.Errorf("ValidateBundle() expected Valid=false for missing required fields")
	}
	// Should have schema validation errors
	found := false
	for _, e := range result.Errors {
		if e.Category == "schema" && (strings.Contains(e.Message, "bundle_name") ||
			strings.Contains(e.Message, "bundle_version") ||
			strings.Contains(e.Message, "presets")) {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("ValidateBundle() expected schema errors for missing fields, got: %v", result.Errors)
	}
}

// TestValidateBundle_MissingEntrypoint tests validation with missing entrypoint file
func TestValidateBundle_MissingEntrypoint(t *testing.T) {
	tmpDir := t.TempDir()

	// Create valid manifest referencing non-existent entrypoint
	manifest := `{
		"manifest_version": "1.0.0",
		"bundle_name": "test-bundle",
		"bundle_version": "v1.0.0",
		"presets": [{"name": "test", "description": "Test", "entrypoint": "nonexistent.json", "prompt_files": []}]
	}`
	manifestPath := filepath.Join(tmpDir, "opencode-bundle.manifest.json")
	if err := os.WriteFile(manifestPath, []byte(manifest), 0644); err != nil {
		t.Fatal(err)
	}

	result, err := ValidateBundle(tmpDir)
	if err != nil {
		t.Errorf("ValidateBundle() unexpected error = %v", err)
		return
	}
	if result.Valid {
		t.Errorf("ValidateBundle() expected Valid=false for missing entrypoint")
	}
	// Should have an entrypoint error
	found := false
	for _, e := range result.Errors {
		if e.Category == "entrypoint" {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("ValidateBundle() expected entrypoint error, got: %v", result.Errors)
	}
}

// TestValidateBundle_InvalidPresetJSON tests validation with invalid preset JSON
func TestValidateBundle_InvalidPresetJSON(t *testing.T) {
	tmpDir := t.TempDir()

	// Create valid manifest
	manifest := `{
		"manifest_version": "1.0.0",
		"bundle_name": "test-bundle",
		"bundle_version": "v1.0.0",
		"presets": [{"name": "test", "description": "Test", "entrypoint": "test.json", "prompt_files": []}]
	}`
	manifestPath := filepath.Join(tmpDir, "opencode-bundle.manifest.json")
	if err := os.WriteFile(manifestPath, []byte(manifest), 0644); err != nil {
		t.Fatal(err)
	}

	// Create invalid JSON entrypoint
	if err := os.WriteFile(filepath.Join(tmpDir, "test.json"), []byte("not valid json"), 0644); err != nil {
		t.Fatal(err)
	}

	result, err := ValidateBundle(tmpDir)
	if err != nil {
		t.Errorf("ValidateBundle() unexpected error = %v", err)
		return
	}
	if result.Valid {
		t.Errorf("ValidateBundle() expected Valid=false for invalid preset JSON")
	}
	// Should have a preset error
	found := false
	for _, e := range result.Errors {
		if e.Category == "preset" {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("ValidateBundle() expected preset error, got: %v", result.Errors)
	}
}

// TestValidateBundle_MissingPromptFile tests validation with missing prompt file
func TestValidateBundle_MissingPromptFile(t *testing.T) {
	tmpDir := t.TempDir()

	// Create valid manifest referencing non-existent prompt file
	manifest := `{
		"manifest_version": "1.0.0",
		"bundle_name": "test-bundle",
		"bundle_version": "v1.0.0",
		"presets": [{"name": "test", "description": "Test", "entrypoint": "test.json", "prompt_files": ["missing.txt"]}]
	}`
	manifestPath := filepath.Join(tmpDir, "opencode-bundle.manifest.json")
	if err := os.WriteFile(manifestPath, []byte(manifest), 0644); err != nil {
		t.Fatal(err)
	}

	// Create valid entrypoint
	if err := os.WriteFile(filepath.Join(tmpDir, "test.json"), []byte("{}"), 0644); err != nil {
		t.Fatal(err)
	}

	result, err := ValidateBundle(tmpDir)
	if err != nil {
		t.Errorf("ValidateBundle() unexpected error = %v", err)
		return
	}
	if result.Valid {
		t.Errorf("ValidateBundle() expected Valid=false for missing prompt file")
	}
	// Should have a prompt_files error
	found := false
	for _, e := range result.Errors {
		if e.Category == "prompt_files" {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("ValidateBundle() expected prompt_files error, got: %v", result.Errors)
	}
}

// TestValidateBundle_EmptyPresetsArray tests validation with empty presets array (should fail per schema)
func TestValidateBundle_EmptyPresetsArray(t *testing.T) {
	tmpDir := t.TempDir()

	// Create manifest with empty presets array - schema requires minItems: 1
	manifest := `{
		"manifest_version": "1.0.0",
		"bundle_name": "test-bundle",
		"bundle_version": "v1.0.0",
		"presets": []
	}`
	manifestPath := filepath.Join(tmpDir, "opencode-bundle.manifest.json")
	if err := os.WriteFile(manifestPath, []byte(manifest), 0644); err != nil {
		t.Fatal(err)
	}

	result, err := ValidateBundle(tmpDir)
	if err != nil {
		t.Errorf("ValidateBundle() unexpected error = %v", err)
		return
	}
	// Empty presets array should fail validation (schema minItems: 1)
	if result.Valid {
		t.Errorf("ValidateBundle() expected Valid=false for empty presets array")
	}
	// Should have a schema error about presets
	found := false
	for _, e := range result.Errors {
		if e.Category == "schema" && (strings.Contains(e.Message, "presets") || strings.Contains(e.Message, "minItems")) {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("ValidateBundle() expected schema error for empty presets, got: %v", result.Errors)
	}
}

// TestValidateBundle_DeeplyNestedBundle tests validation with nested directory structure
func TestValidateBundle_DeeplyNestedBundle(t *testing.T) {
	tmpDir := t.TempDir()

	// Create nested directory structure with manifest at root
	presetsDir := filepath.Join(tmpDir, "presets")
	subDir := filepath.Join(presetsDir, "v1")
	if err := os.MkdirAll(subDir, 0755); err != nil {
		t.Fatal(err)
	}

	// Create manifest with entrypoints in subdirectories
	manifest := `{
		"manifest_version": "1.0.0",
		"bundle_name": "nested-bundle",
		"bundle_version": "v1.0.0",
		"presets": [
			{
				"name": "default",
				"description": "Default preset",
				"entrypoint": "presets/v1/default.json",
				"prompt_files": ["prompts/readme.md"]
			}
		]
	}`
	manifestPath := filepath.Join(tmpDir, "opencode-bundle.manifest.json")
	if err := os.WriteFile(manifestPath, []byte(manifest), 0644); err != nil {
		t.Fatal(err)
	}

	// Create preset file in subdirectory
	presetContent := `{"agents":[]}`
	if err := os.WriteFile(filepath.Join(subDir, "default.json"), []byte(presetContent), 0644); err != nil {
		t.Fatal(err)
	}

	// Create prompts directory and file
	promptsDir := filepath.Join(tmpDir, "prompts")
	if err := os.MkdirAll(promptsDir, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(promptsDir, "readme.md"), []byte("# Readme"), 0644); err != nil {
		t.Fatal(err)
	}

	result, err := ValidateBundle(tmpDir)
	if err != nil {
		t.Errorf("ValidateBundle() unexpected error = %v", err)
		return
	}
	if !result.Valid {
		t.Errorf("ValidateBundle() expected Valid=true for nested bundle, got errors: %v", result.Errors)
	}
}

// TestValidateBundle_MultipleErrors tests that all validation errors are collected (not fail-fast)
func TestValidateBundle_MultipleErrors(t *testing.T) {
	tmpDir := t.TempDir()

	// Create manifest with multiple problems:
	// - Preset 1: missing entrypoint
	// - Preset 2: missing prompt_file
	// - Preset 3: invalid preset JSON
	manifest := `{
		"manifest_version": "1.0.0",
		"bundle_name": "multi-error-bundle",
		"bundle_version": "v1.0.0",
		"presets": [
			{
				"name": "missing-entrypoint",
				"description": "Preset with missing entrypoint",
				"entrypoint": "nonexistent1.json"
			},
			{
				"name": "missing-prompt",
				"description": "Preset with missing prompt file",
				"entrypoint": "valid.json",
				"prompt_files": ["missing-prompt.md"]
			},
			{
				"name": "invalid-json",
				"description": "Preset with invalid JSON",
				"entrypoint": "invalid.json"
			}
		]
	}`
	manifestPath := filepath.Join(tmpDir, "opencode-bundle.manifest.json")
	if err := os.WriteFile(manifestPath, []byte(manifest), 0644); err != nil {
		t.Fatal(err)
	}

	// Create valid.json (exists)
	if err := os.WriteFile(filepath.Join(tmpDir, "valid.json"), []byte("{}"), 0644); err != nil {
		t.Fatal(err)
	}

	// Create invalid.json (invalid JSON)
	if err := os.WriteFile(filepath.Join(tmpDir, "invalid.json"), []byte("not valid json {"), 0644); err != nil {
		t.Fatal(err)
	}

	result, err := ValidateBundle(tmpDir)
	if err != nil {
		t.Errorf("ValidateBundle() unexpected error = %v", err)
		return
	}
	if result.Valid {
		t.Errorf("ValidateBundle() expected Valid=false for multiple errors")
	}

	// Should collect at least 3 errors (missing entrypoint, missing prompt, invalid json)
	// We expect: 1 entrypoint error + 1 prompt_files error + 1 preset error = 3 errors minimum
	if len(result.Errors) < 3 {
		t.Errorf("ValidateBundle() expected at least 3 errors collected, got %d: %v", len(result.Errors), result.Errors)
	}

	// Verify we have each category of error
	categories := make(map[string]bool)
	for _, e := range result.Errors {
		categories[e.Category] = true
	}
	if !categories["entrypoint"] {
		t.Errorf("ValidateBundle() expected entrypoint error category")
	}
	if !categories["prompt_files"] {
		t.Errorf("ValidateBundle() expected prompt_files error category")
	}
	if !categories["preset"] {
		t.Errorf("ValidateBundle() expected preset error category")
	}
}
