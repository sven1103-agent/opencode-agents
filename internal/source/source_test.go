package source

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestDetectSourceType_GitHubInputs(t *testing.T) {
	tests := []struct {
		name     string
		location string
	}{
		{
			name:     "bare owner repo",
			location: "qbicsoftware/opencode-config-bundle",
		},
		{
			name:     "github host path",
			location: "github.com/qbicsoftware/opencode-config-bundle",
		},
		{
			name:     "github release url",
			location: "https://github.com/qbicsoftware/opencode-config-bundle/releases/tag/1.0.0-alpha.1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := DetectSourceType(tt.location)
			if err != nil {
				t.Fatalf("DetectSourceType() error = %v", err)
			}
			if got != SourceTypeGitHubRelease {
				t.Fatalf("DetectSourceType() = %q, want %q", got, SourceTypeGitHubRelease)
			}
		})
	}
}

func TestDetectSourceType_LocalInputs(t *testing.T) {
	tmpDir := t.TempDir()

	detectedDir, err := DetectSourceType(tmpDir)
	if err != nil {
		t.Fatalf("DetectSourceType(directory) error = %v", err)
	}
	if detectedDir != SourceTypeLocalDirectory {
		t.Fatalf("DetectSourceType(directory) = %q, want %q", detectedDir, SourceTypeLocalDirectory)
	}

	archivePath := filepath.Join(t.TempDir(), "bundle.tar.gz")
	if got, err := DetectSourceType(archivePath); err != nil {
		t.Fatalf("DetectSourceType(archive) error = %v", err)
	} else if got != SourceTypeLocalArchive {
		t.Fatalf("DetectSourceType(archive) = %q, want %q", got, SourceTypeLocalArchive)
	}
}

func TestParseGitHubRef(t *testing.T) {
	tests := []struct {
		name        string
		location    string
		wantRepo    string
		wantTag     string
		wantErr     bool
		errContains string
	}{
		{
			name:     "bare owner repo",
			location: "qbicsoftware/opencode-config-bundle",
			wantRepo: "qbicsoftware/opencode-config-bundle",
		},
		{
			name:     "github host path",
			location: "github.com/qbicsoftware/opencode-config-bundle",
			wantRepo: "qbicsoftware/opencode-config-bundle",
		},
		{
			name:     "release url pins tag",
			location: "https://github.com/qbicsoftware/opencode-config-bundle/releases/tag/1.0.0-alpha.1",
			wantRepo: "qbicsoftware/opencode-config-bundle",
			wantTag:  "1.0.0-alpha.1",
		},
		{
			name:        "invalid github path",
			location:    "github.com/qbicsoftware",
			wantErr:     true,
			errContains: "invalid GitHub",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ref, err := parseGitHubRef(tt.location)
			if tt.wantErr {
				if err == nil {
					t.Fatalf("parseGitHubRef() error = nil, want error")
				}
				if tt.errContains != "" && !contains(err.Error(), tt.errContains) {
					t.Fatalf("parseGitHubRef() error = %v, want substring %q", err, tt.errContains)
				}
				return
			}

			if err != nil {
				t.Fatalf("parseGitHubRef() error = %v", err)
			}
			if ref.Repo != tt.wantRepo {
				t.Fatalf("parseGitHubRef().Repo = %q, want %q", ref.Repo, tt.wantRepo)
			}
			if ref.Tag != tt.wantTag {
				t.Fatalf("parseGitHubRef().Tag = %q, want %q", ref.Tag, tt.wantTag)
			}
		})
	}
}

func TestValidateSource_GitHubInput(t *testing.T) {
	err := ValidateSource("qbicsoftware/opencode-config-bundle", SourceTypeGitHubRelease)
	if err != nil {
		t.Fatalf("ValidateSource() error = %v", err)
	}
}

func TestResolveSourceRef(t *testing.T) {
	setRegistryForTest(t, []Source{
		{ID: "id-1", Name: "qbic", Type: SourceTypeLocalDirectory, Location: "/tmp/a"},
		{ID: "id-2", Name: "other", Type: SourceTypeLocalDirectory, Location: "/tmp/b"},
	})

	byID, err := ResolveSourceRef("id-1")
	if err != nil {
		t.Fatalf("ResolveSourceRef(id) error = %v", err)
	}
	if byID.Name != "qbic" {
		t.Fatalf("ResolveSourceRef(id).Name = %q", byID.Name)
	}

	byName, err := ResolveSourceRef("qbic")
	if err != nil {
		t.Fatalf("ResolveSourceRef(name) error = %v", err)
	}
	if byName.ID != "id-1" {
		t.Fatalf("ResolveSourceRef(name).ID = %q", byName.ID)
	}
}

func TestResolveSourceRefAmbiguousName(t *testing.T) {
	setRegistryForTest(t, []Source{
		{ID: "id-1", Name: "qbic", Type: SourceTypeLocalDirectory, Location: "/tmp/a"},
		{ID: "id-2", Name: "qbic", Type: SourceTypeLocalDirectory, Location: "/tmp/b"},
	})

	_, err := ResolveSourceRef("qbic")
	if err == nil {
		t.Fatal("ResolveSourceRef() error = nil, want ambiguity error")
	}
	if !strings.Contains(err.Error(), "ambiguous") || !strings.Contains(err.Error(), "id-1") || !strings.Contains(err.Error(), "id-2") {
		t.Fatalf("ResolveSourceRef() error = %v", err)
	}
}

func setRegistryForTest(t *testing.T, sources []Source) {
	t.Helper()
	configDir := t.TempDir()
	homeDir := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", configDir)
	t.Setenv("HOME", homeDir)

	registryPath, err := RegistryPath()
	if err != nil {
		t.Fatalf("failed to resolve registry path: %v", err)
	}
	if err := os.MkdirAll(filepath.Dir(registryPath), 0755); err != nil {
		t.Fatalf("failed to create config dir: %v", err)
	}

	data, err := json.Marshal(&Registry{Version: 1, Sources: sources})
	if err != nil {
		t.Fatalf("failed to marshal registry: %v", err)
	}
	if err := os.WriteFile(registryPath, data, 0644); err != nil {
		t.Fatalf("failed to write registry: %v", err)
	}
}

func contains(s, substr string) bool {
	return len(substr) == 0 || (len(s) >= len(substr) && indexOf(s, substr) >= 0)
}

func indexOf(s, substr string) int {
	for i := 0; i+len(substr) <= len(s); i++ {
		if s[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}
