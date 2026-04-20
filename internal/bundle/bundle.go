// Package bundle provides functionality for managing OpenCode configuration bundles.
package bundle

import (
	"crypto/sha256"
	_ "embed"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/qbicsoftware/occo/internal/source"
	"github.com/santhosh-tekuri/jsonschema/v5"
)

//go:embed 1.0.0.schema.json
var embeddedSchema string

var schemaCompiler *jsonschema.Compiler

var (
	githubHTTPClient       = http.DefaultClient
	githubAPIBaseURL       = "https://api.github.com"
	githubDownloadBaseURL  = "https://github.com"
	githubCacheDirOverride string
)

func init() {
	schemaCompiler = jsonschema.NewCompiler()
	if err := schemaCompiler.AddResource("1.0.0.schema.json", strings.NewReader(embeddedSchema)); err != nil {
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

// ValidationError represents an error found during bundle validation.
type ValidationError struct {
	Category string // "manifest", "schema", "entrypoint", "prompt_files", "preset"
	Message  string
	File     string // path to file with error
}

// ValidationResult represents the result of bundle validation.
type ValidationResult struct {
	Valid    bool
	Errors   []ValidationError
	Warnings []string
}

// ValidateBundle validates a bundle at the given root directory.
// Returns a ValidationResult with all errors collected (not fail-fast).
// System errors (IO, permissions) return a non-nil error.
// Validation errors return nil error with ValidationResult.Valid=false.
func ValidateBundle(bundleRoot string) (*ValidationResult, error) {
	result := &ValidationResult{
		Valid:    true,
		Errors:   []ValidationError{},
		Warnings: []string{},
	}

	// 1. Check manifest exists
	manifestPath := filepath.Join(bundleRoot, "opencode-bundle.manifest.json")
	manifestData, err := os.ReadFile(manifestPath)
	if err != nil {
		if os.IsNotExist(err) {
			result.Valid = false
			result.Errors = append(result.Errors, ValidationError{
				Category: "manifest",
				Message:  "manifest file not found",
				File:     manifestPath,
			})
			return result, nil
		}
		return nil, fmt.Errorf("failed to read manifest: %w", err)
	}

	// 2. Parse JSON and validate against schema
	var rawManifest interface{}
	if err := json.Unmarshal(manifestData, &rawManifest); err != nil {
		result.Valid = false
		result.Errors = append(result.Errors, ValidationError{
			Category: "manifest",
			Message:  fmt.Sprintf("invalid JSON: %v", err),
			File:     manifestPath,
		})
		return result, nil
	}

	// Validate against embedded schema
	schema, err := schemaCompiler.Compile("1.0.0.schema.json")
	if err != nil {
		return nil, fmt.Errorf("failed to compile schema: %w", err)
	}

	if err := schema.Validate(rawManifest); err != nil {
		result.Valid = false
		result.Errors = append(result.Errors, ValidationError{
			Category: "schema",
			Message:  fmt.Sprintf("schema validation failed: %v", err),
			File:     manifestPath,
		})
		return result, nil
	}

	// 3. Load manifest for version validation
	manifest, err := LoadManifest(manifestPath)
	if err != nil {
		result.Valid = false
		result.Errors = append(result.Errors, ValidationError{
			Category: "manifest",
			Message:  err.Error(),
			File:     manifestPath,
		})
		return result, nil
	}

	// 4. Validate each preset
	for _, preset := range manifest.Presets {
		// Check entrypoint exists
		entrypointPath := filepath.Join(bundleRoot, preset.Entrypoint)
		entrypointData, err := os.ReadFile(entrypointPath)
		if err != nil {
			result.Valid = false
			if os.IsNotExist(err) {
				result.Errors = append(result.Errors, ValidationError{
					Category: "entrypoint",
					Message:  "entrypoint file not found",
					File:     entrypointPath,
				})
			} else {
				result.Errors = append(result.Errors, ValidationError{
					Category: "entrypoint",
					Message:  fmt.Sprintf("failed to read entrypoint: %v", err),
					File:     entrypointPath,
				})
			}
			continue
		}

		// Check entrypoint is valid JSON
		var presetData interface{}
		if err := json.Unmarshal(entrypointData, &presetData); err != nil {
			result.Valid = false
			result.Errors = append(result.Errors, ValidationError{
				Category: "preset",
				Message:  fmt.Sprintf("invalid JSON in preset: %v", err),
				File:     entrypointPath,
			})
		}

		// Check prompt_files exist
		for _, promptFile := range preset.PromptFiles {
			promptPath := filepath.Join(bundleRoot, promptFile)
			if _, err := os.Stat(promptPath); err != nil {
				result.Valid = false
				if os.IsNotExist(err) {
					result.Errors = append(result.Errors, ValidationError{
						Category: "prompt_files",
						Message:  "prompt file not found",
						File:     promptPath,
					})
				} else {
					result.Errors = append(result.Errors, ValidationError{
						Category: "prompt_files",
						Message:  fmt.Sprintf("failed to read prompt file: %v", err),
						File:     promptPath,
					})
				}
			}
		}
	}

	return result, nil
}

// InstalledAsset represents a file installed from a bundle to the project.
type InstalledAsset struct {
	Source      string `json:"source"`      // Path in bundle (relative to bundle root)
	Destination string `json:"destination"` // Relative path in project .opencode/
}

// Provenance represents the bundle provenance file stored in the project.
type Provenance struct {
	SourceID        string           `json:"source_id"`
	SourceName      string           `json:"source_name"`
	SourceType      string           `json:"source_type"`
	BundleVersion   string           `json:"bundle_version"`
	PresetName      string           `json:"preset_name"`
	Entrypoint      string           `json:"entrypoint"`
	AppliedAt       string           `json:"applied_at"`
	InstalledAssets []InstalledAsset `json:"installed_assets,omitempty"`
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
			return nil, fmt.Errorf("no provenance file found (run 'bundle install' first)")
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
		tmpDir, err := os.MkdirTemp("", "opencode-bundle-install")
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
		bundleRoot, err := resolveGitHubReleaseToLocal(sourceLocation, versionTag)
		if err != nil {
			return "", nil, err
		}
		return bundleRoot, cleanup, nil

	default:
		return "", nil, fmt.Errorf("unknown source type: %s", sourceType)
	}
}

type githubReleaseResponse struct {
	TagName    string               `json:"tag_name"`
	Prerelease bool                 `json:"prerelease"`
	Assets     []githubReleaseAsset `json:"assets"`
}

type githubReleaseAsset struct {
	Name               string `json:"name"`
	BrowserDownloadURL string `json:"browser_download_url"`
}

// GitHubReleaseVersion describes a GitHub release version available for a bundle source.
type GitHubReleaseVersion struct {
	TagName    string
	Prerelease bool
}

func resolveGitHubReleaseToLocal(sourceLocation, versionTag string) (string, error) {
	ref, err := source.ParseGitHubLocation(sourceLocation)
	if err != nil {
		return "", err
	}

	tag := versionTag
	if tag == "" {
		tag = ref.Tag
	}
	if tag == "" {
		return "", fmt.Errorf("--version is required for github-release sources outside interactive mode (use --version latest or --version <tag>)")
	}

	release, err := fetchGitHubRelease(ref.Repo, tag)
	if err != nil {
		return "", err
	}

	archiveAsset, checksumsAsset := selectGitHubAssets(release)
	if archiveAsset == nil {
		return "", fmt.Errorf("no bundle asset (.tar.gz) found in release %s", release.TagName)
	}

	cacheDir, err := githubCacheDir()
	if err != nil {
		return "", err
	}
	if err := os.MkdirAll(cacheDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create cache directory: %w", err)
	}

	cachedBase := filepath.Join(cacheDir, strings.ReplaceAll(ref.Repo, "/", "-")+"-"+release.TagName)
	cachedTarball := cachedBase + ".tar.gz"
	cachedExtract := cachedBase

	if bundleRoot, err := cachedBundleRoot(cachedExtract); err == nil {
		return bundleRoot, nil
	}

	if err := downloadToFile(archiveAsset.BrowserDownloadURL, cachedTarball); err != nil {
		return "", err
	}

	if checksumsAsset != nil {
		if err := verifyGitHubChecksum(cachedTarball, archiveAsset.Name, checksumsAsset.BrowserDownloadURL); err != nil {
			_ = os.Remove(cachedTarball)
			_ = os.RemoveAll(cachedExtract)
			return "", err
		}
	}

	_ = os.RemoveAll(cachedExtract)
	if err := os.MkdirAll(cachedExtract, 0755); err != nil {
		return "", fmt.Errorf("failed to create cache extract directory: %w", err)
	}
	if err := extractTarball(cachedTarball, cachedExtract); err != nil {
		_ = os.Remove(cachedTarball)
		_ = os.RemoveAll(cachedExtract)
		return "", fmt.Errorf("failed to extract bundle: %w", err)
	}

	bundleRoot, err := findBundleRoot(cachedExtract)
	if err != nil {
		_ = os.Remove(cachedTarball)
		_ = os.RemoveAll(cachedExtract)
		return "", err
	}
	if _, err := os.Stat(filepath.Join(bundleRoot, "opencode-bundle.manifest.json")); err != nil {
		_ = os.Remove(cachedTarball)
		_ = os.RemoveAll(cachedExtract)
		return "", fmt.Errorf("bundle manifest not found in downloaded archive")
	}

	return bundleRoot, nil
}

func githubCacheDir() (string, error) {
	if githubCacheDirOverride != "" {
		return githubCacheDirOverride, nil
	}
	cacheRoot, err := os.UserCacheDir()
	if err != nil {
		return "", fmt.Errorf("failed to determine cache directory: %w", err)
	}
	return filepath.Join(cacheRoot, "opencode-helper", "github-releases"), nil
}

func cachedBundleRoot(cachedExtract string) (string, error) {
	if _, err := os.Stat(cachedExtract); err != nil {
		return "", err
	}
	return findBundleRoot(cachedExtract)
}

// ListGitHubReleases returns usable GitHub bundle releases for a source location.
func ListGitHubReleases(sourceLocation string) ([]GitHubReleaseVersion, error) {
	ref, err := source.ParseGitHubLocation(sourceLocation)
	if err != nil {
		return nil, err
	}

	releases, err := fetchGitHubReleases(ref.Repo)
	if err != nil {
		return nil, err
	}

	versions := make([]GitHubReleaseVersion, 0, len(releases))
	for _, release := range releases {
		if archiveAsset, _ := selectGitHubAssets(release); archiveAsset == nil {
			continue
		}
		versions = append(versions, GitHubReleaseVersion{TagName: release.TagName, Prerelease: release.Prerelease})
	}

	if len(versions) == 0 {
		return nil, fmt.Errorf("no usable bundle releases found for %s", ref.Repo)
	}

	return versions, nil
}

func fetchGitHubRelease(repo, tag string) (*githubReleaseResponse, error) {
	if tag == "latest" {
		releases, err := fetchGitHubReleases(repo)
		if err != nil {
			return nil, err
		}

		var sawPrerelease bool
		for _, release := range releases {
			if release.Prerelease {
				sawPrerelease = true
				continue
			}
			return fetchGitHubReleaseFromPath("/repos/" + repo + "/releases/tags/" + release.TagName)
		}

		if sawPrerelease {
			return nil, fmt.Errorf("no stable release found for %s; prereleases are available (use --version <tag>)", repo)
		}
		return nil, fmt.Errorf("no releases found for %s", repo)
	}

	path := "/repos/" + repo + "/releases/tags/" + tag
	return fetchGitHubReleaseFromPath(path)
}

func fetchGitHubReleases(repo string) ([]*githubReleaseResponse, error) {
	path := "/repos/" + repo + "/releases?per_page=100"
	body, err := fetchGitHubReleaseBody(path)
	if err != nil {
		return nil, err
	}

	var releases []*githubReleaseResponse
	if err := json.Unmarshal(body, &releases); err != nil {
		return nil, fmt.Errorf("failed to parse releases response: %w", err)
	}
	if len(releases) == 0 {
		return nil, fmt.Errorf("no releases found for %s", repo)
	}
	return releases, nil
}

func fetchGitHubReleaseFromPath(path string) (*githubReleaseResponse, error) {
	body, err := fetchGitHubReleaseBody(path)
	if err != nil {
		return nil, err
	}

	var release githubReleaseResponse
	if err := json.Unmarshal(body, &release); err != nil {
		return nil, fmt.Errorf("failed to parse release response: %w", err)
	}
	if release.TagName == "" {
		return nil, fmt.Errorf("failed to parse tag_name from release")
	}
	return &release, nil
}

func fetchGitHubReleaseBody(path string) ([]byte, error) {
	apiBase := githubAPIBaseURL
	if envBase := os.Getenv("OC_GITHUB_API_BASE_URL"); envBase != "" {
		apiBase = envBase
	}

	resp, err := githubHTTPClient.Get(strings.TrimRight(apiBase, "/") + path)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch release info from GitHub: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		msg := strings.TrimSpace(string(body))
		if msg == "" {
			msg = resp.Status
		}
		return nil, fmt.Errorf("GitHub API error: %s", msg)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read release response: %w", err)
	}
	return body, nil
}

func selectGitHubAssets(release *githubReleaseResponse) (*githubReleaseAsset, *githubReleaseAsset) {
	var archiveAsset *githubReleaseAsset
	var checksumsAsset *githubReleaseAsset
	for i := range release.Assets {
		asset := &release.Assets[i]
		if archiveAsset == nil && strings.HasSuffix(asset.Name, ".tar.gz") {
			archiveAsset = asset
		}
		if checksumsAsset == nil && strings.HasSuffix(asset.Name, "-checksums.txt") {
			checksumsAsset = asset
		}
	}
	return archiveAsset, checksumsAsset
}

func downloadToFile(url, dest string) error {
	resp, err := githubHTTPClient.Get(url)
	if err != nil {
		return fmt.Errorf("failed to download bundle: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to download bundle: %s", resp.Status)
	}

	file, err := os.Create(dest)
	if err != nil {
		return fmt.Errorf("failed to create cache file: %w", err)
	}
	defer file.Close()

	if _, err := io.Copy(file, resp.Body); err != nil {
		return fmt.Errorf("failed to write downloaded bundle: %w", err)
	}
	return nil
}

func verifyGitHubChecksum(archivePath, assetName, checksumsURL string) error {
	resp, err := githubHTTPClient.Get(checksumsURL)
	if err != nil {
		return fmt.Errorf("failed to download checksums file: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to download checksums file: %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read checksums file: %w", err)
	}

	var expected string
	for _, line := range strings.Split(string(body), "\n") {
		fields := strings.Fields(line)
		if len(fields) == 2 && fields[1] == assetName {
			expected = fields[0]
			break
		}
	}
	if expected == "" {
		return fmt.Errorf("checksums file missing entry for: %s", assetName)
	}

	data, err := os.ReadFile(archivePath)
	if err != nil {
		return fmt.Errorf("failed to read downloaded bundle: %w", err)
	}
	actual := sha256Hex(data)
	if actual != expected {
		return fmt.Errorf("bundle integrity check failed: SHA256 mismatch for %s", assetName)
	}
	return nil
}

func sha256Hex(data []byte) string {
	sum := sha256.Sum256(data)
	return fmt.Sprintf("%x", sum[:])
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
