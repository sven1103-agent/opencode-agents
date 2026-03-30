package cmd

import (
	"os"
	"path/filepath"
	"testing"
)

// TestValidateOutputPath tests the path traversal protection
func TestValidateOutputPath(t *testing.T) {
	tests := []struct {
		name        string
		projectRoot string
		outputPath  string
		wantErr     bool
	}{
		{
			name:        "valid path within project root",
			projectRoot: "/tmp/project",
			outputPath:  "/tmp/project/opencode.json",
			wantErr:     false,
		},
		{
			name:        "valid nested path",
			projectRoot: "/tmp/project",
			outputPath:  "/tmp/project/subdir/opencode.json",
			wantErr:     false,
		},
		{
			name:        "path traversal attempt",
			projectRoot: "/tmp/project",
			outputPath:  "/tmp/project/../../../etc/passwd",
			wantErr:     true,
		},
		{
			name:        "absolute path traversal",
			projectRoot: "/tmp/project",
			outputPath:  "/etc/passwd",
			wantErr:     true,
		},
		{
			name:        "sibling directory traversal",
			projectRoot: "/tmp/project",
			outputPath:  "/tmp/other-file",
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateOutputPath(tt.projectRoot, tt.outputPath)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateOutputPath() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// TestInitAcceptanceCriteria tests the full init command against acceptance criteria
func TestInitAcceptanceCriteria(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "oc-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Test: initCmd is registered
	if initCmd == nil {
		t.Error("initCmd should be registered")
	}

	// Test: Output path validation
	validOutput := filepath.Join(tmpDir, "opencode.json")
	if err := validateOutputPath(tmpDir, validOutput); err != nil {
		t.Errorf("valid output path should not error: %v", err)
	}

	// Test: Path traversal should be rejected
	traversalOutput := filepath.Join(tmpDir, "..", "..", "etc", "passwd")
	if err := validateOutputPath(tmpDir, traversalOutput); err == nil {
		t.Error("path traversal should be rejected")
	}
}
