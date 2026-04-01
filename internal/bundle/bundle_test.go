package bundle

import (
	"os"
	"path/filepath"
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
