package schema

import (
	"testing"
)

// TestSchemas tests that all expected schemas are defined
func TestSchemas(t *testing.T) {
	schemas := Schemas()

	if len(schemas) != 2 {
		t.Errorf("expected 2 schemas, got %d", len(schemas))
	}

	expected := map[string]string{
		"handoff": "handoff.schema.json",
		"result":  "result.schema.json",
	}

	for _, s := range schemas {
		if expected[s.Name] != s.Filename {
			t.Errorf("expected schema %s to have filename %s, got %s", s.Name, expected[s.Name], s.Filename)
		}
	}
}

// TestSchemaDir tests schema directory name
func TestSchemaDir(t *testing.T) {
	if got := SchemaDir(); got != "schemas" {
		t.Errorf("SchemaDir() = %v, want 'schemas'", got)
	}
}

// TestValidatePath tests path traversal protection
func TestValidatePath(t *testing.T) {
	tests := []struct {
		name      string
		baseDir   string
		targetDir string
		wantErr   bool
	}{
		{
			name:      "valid path within base",
			baseDir:   "/tmp/project",
			targetDir: "/tmp/project/subdir",
			wantErr:   false,
		},
		{
			name:      "same directory",
			baseDir:   "/tmp/project",
			targetDir: "/tmp/project",
			wantErr:   false,
		},
		{
			name:      "path traversal attempt",
			baseDir:   "/tmp/project",
			targetDir: "/tmp/project/../../../etc",
			wantErr:   true,
		},
		{
			name:      "sibling directory",
			baseDir:   "/tmp/project",
			targetDir: "/tmp/other",
			wantErr:   true,
		},
		{
			name:      "absolute path outside",
			baseDir:   "/tmp/project",
			targetDir: "/etc/passwd",
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidatePath(tt.baseDir, tt.targetDir)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidatePath() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
