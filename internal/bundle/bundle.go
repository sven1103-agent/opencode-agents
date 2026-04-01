// Package bundle provides functionality for managing OpenCode configuration bundles.
package bundle

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/santhosh-tekuri/jsonschema/v5"
)

//go:embed schema.json
var embeddedSchema string

var schemaCompiler *jsonschema.Compiler

func init() {
	schemaCompiler = jsonschema.NewCompiler()
	if err := schemaCompiler.AddResource("schema.json", strings.NewReader(embeddedSchema)); err != nil {
		panic(fmt.Sprintf("failed to load embedded schema: %v", err))
	}
}

func isVersionSupported(version string) bool {
	for _, v := range supportedVersions {
		if version == v {
			return true
		}
	}
	return false
}

// Manifest represents the bundle manifest file format.
type Manifest struct {
	ManifestVersion string   `json:"manifest_version"`
	BundleName      string   `json:"bundle_name"`
	BundleVersion   string   `json:"bundle_version"`
	BundleRoot      string   `json:"bundle_root"`
	Presets         []Preset `json:"presets"`
	UpdateCapable   bool     `json:"update_capable,omitempty"`
	UpdateCheckURL  string   `json:"update_check_url,omitempty"`
}

var supportedVersions = []string{"1.0.0"}

// Preset represents a preset entry in the bundle manifest.
type Preset struct {
	Name        string   `json:"name"`
	Entrypoint  string   `json:"entrypoint"`
	PromptFiles []string `json:"prompt_files,omitempty"`
	Description string   `json:"description,omitempty"`
}

// Provenance represents the bundle provenance file stored in the project.
type Provenance struct {
	SourceID      string `json:"source_id"`
	SourceName    string `json:"source_name"`
	SourceType    string `json:"source_type"`
	BundleVersion string `json:"bundle_version"`
	PresetName    string `json:"preset_name"`
	Entrypoint    string `json:"entrypoint"`
	AppliedAt     string `json:"applied_at"`
}

// LoadManifest loads a bundle manifest from the given path.
func LoadManifest(manifestPath string) (*Manifest, error) {
	data, err := os.ReadFile(manifestPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read manifest: %w", err)
	}

	var manifest Manifest
	if err := json.Unmarshal(data, &manifest); err != nil {
		return nil, fmt.Errorf("failed to parse manifest: %w", err)
	}

	if !isVersionSupported(manifest.ManifestVersion) {
		return nil, fmt.Errorf("unsupported manifest version: %s (supported: %v)", manifest.ManifestVersion, supportedVersions)
	}

	return &manifest, nil
}

// GetPreset returns a preset by name from the manifest.
func GetPreset(manifest *Manifest, name string) (*Preset, error) {
	for _, p := range manifest.Presets {
		if p.Name == name {
			return &p, nil
		}
	}
	return nil, fmt.Errorf("preset not found: %s", name)
}

// ProvenancePath returns the path to the bundle provenance file in a project.
func ProvenancePath(projectRoot string) string {
	return filepath.Join(projectRoot, ".opencode", "bundle-provenance.json")
}

// LoadProvenance loads the bundle provenance from a project.
func LoadProvenance(projectRoot string) (*Provenance, error) {
	provPath := ProvenancePath(projectRoot)
	data, err := os.ReadFile(provPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("no provenance file found (run 'bundle apply' first)")
		}
		return nil, fmt.Errorf("failed to read provenance: %w", err)
	}

	var prov Provenance
	if err := json.Unmarshal(data, &prov); err != nil {
		return nil, fmt.Errorf("failed to parse provenance: %w", err)
	}

	return &prov, nil
}

// SaveProvenance saves the bundle provenance to a project.
func SaveProvenance(projectRoot string, prov *Provenance, force bool) error {
	opencodeDir := filepath.Join(projectRoot, ".opencode")
	if err := os.MkdirAll(opencodeDir, 0755); err != nil {
		return fmt.Errorf("failed to create .opencode directory: %w", err)
	}

	provPath := ProvenancePath(projectRoot)

	// Check if provenance exists (unless force)
	if !force {
		if _, err := os.Stat(provPath); err == nil {
			return fmt.Errorf("provenance already exists: %s (use --force to overwrite)", provPath)
		}
	}

	data, err := json.MarshalIndent(prov, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal provenance: %w", err)
	}

	if err := os.WriteFile(provPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write provenance: %w", err)
	}

	return nil
}

// ResolveToLocal resolves a source to a local bundle root directory.
// For local directories, returns the path as-is.
// For archives, extracts to a temp directory.
// For GitHub releases, downloads and extracts.
// Returns the local bundle root path and a cleanup function.
func ResolveToLocal(sourceType, sourceLocation, versionTag string) (string, func(), error) {
	cleanup := func() {}

	switch sourceType {
	case "local-directory":
		if _, err := os.Stat(sourceLocation); err != nil {
			return "", nil, fmt.Errorf("source directory not found: %s", sourceLocation)
		}
		return sourceLocation, cleanup, nil

	case "local-archive":
		if _, err := os.Stat(sourceLocation); err != nil {
			return "", nil, fmt.Errorf("source archive not found: %s", sourceLocation)
		}

		// Extract to temp directory
		tmpDir, err := os.MkdirTemp("", "opencode-bundle-apply")
		if err != nil {
			return "", nil, fmt.Errorf("failed to create temp directory: %w", err)
		}

		if err := extractTarball(sourceLocation, tmpDir); err != nil {
			os.RemoveAll(tmpDir)
			return "", nil, fmt.Errorf("failed to extract tarball: %w", err)
		}

		// Determine bundle root based on archive structure
		bundleRoot, err := findBundleRoot(tmpDir)
		if err != nil {
			os.RemoveAll(tmpDir)
			return "", nil, err
		}

		cleanup = func() { os.RemoveAll(tmpDir) }
		return bundleRoot, cleanup, nil

	case "github-release":
		// For now, return an error - GitHub support would require network operations
		return "", nil, fmt.Errorf("github-release sources require network operations (not yet implemented)")

	default:
		return "", nil, fmt.Errorf("unknown source type: %s", sourceType)
	}
}

// extractTarball extracts a .tar.gz archive to the destination directory.
func extractTarball(archivePath, destDir string) error {
	cmd := exec.Command("tar", "-xzf", archivePath, "-C", destDir)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to extract tarball: %w", err)
	}
	return nil
}

// findBundleRoot determines the bundle root from an extracted archive.
func findBundleRoot(extractDir string) (string, error) {
	// Check if manifest exists directly at root (Pattern 2)
	manifestAtRoot := filepath.Join(extractDir, "opencode-bundle.manifest.json")
	if _, err := os.Stat(manifestAtRoot); err == nil {
		return extractDir, nil
	}

	// Pattern 1: Single top-level directory
	entries, err := os.ReadDir(extractDir)
	if err != nil {
		return "", fmt.Errorf("failed to read extract directory: %w", err)
	}

	var dirs []os.DirEntry
	for _, e := range entries {
		if e.IsDir() {
			dirs = append(dirs, e)
		}
	}

	if len(dirs) == 1 {
		return filepath.Join(extractDir, dirs[0].Name()), nil
	}

	if len(dirs) == 0 {
		return "", fmt.Errorf("archive has no content")
	}

	return "", fmt.Errorf("archive has multiple top-level items (expected single directory)")
}
