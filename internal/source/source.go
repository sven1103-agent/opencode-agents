// Package source provides functionality for managing OpenCode config sources.
package source

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"
)

// SourceType represents the type of a config source.
type SourceType string

const (
	// SourceTypeLocalDirectory represents a local directory source.
	SourceTypeLocalDirectory SourceType = "local-directory"
	// SourceTypeLocalArchive represents a local archive (.tar.gz) source.
	SourceTypeLocalArchive SourceType = "local-archive"
	// SourceTypeGitHubRelease represents a GitHub release source.
	SourceTypeGitHubRelease SourceType = "github-release"
)

// Source represents a registered config source.
type Source struct {
	ID        string     `json:"id"`
	Location  string     `json:"location"`
	Type      SourceType `json:"type"`
	Name      string     `json:"name"`
	CreatedAt string     `json:"created_at"`
}

// Registry represents the source registry file format.
type Registry struct {
	Version int      `json:"version"`
	Sources []Source `json:"sources"`
}

// GitHubRef represents a normalized GitHub repository reference.
type GitHubRef struct {
	Repo string
	Tag  string
}

// AmbiguousSourceRefError reports when a source name matches more than one source.
type AmbiguousSourceRefError struct {
	Ref     string
	Matches []Source
}

func (e *AmbiguousSourceRefError) Error() string {
	parts := make([]string, 0, len(e.Matches))
	for _, match := range e.Matches {
		parts = append(parts, fmt.Sprintf("%s (%s)", match.Name, match.ID))
	}
	return fmt.Sprintf("source name is ambiguous: %s (matches: %s)", e.Ref, strings.Join(parts, ", "))
}

var ownerRepoPattern = regexp.MustCompile(`^[A-Za-z0-9_.-]+/[A-Za-z0-9_.-]+$`)

// RegistryPath returns the path to the source registry file.
// It respects XDG_CONFIG_HOME for cross-platform config location.
// Defaults to ~/.config (XDG standard) if XDG_CONFIG_HOME is not set.
func RegistryPath() (string, error) {
	configDir := os.Getenv("XDG_CONFIG_HOME")
	if configDir == "" {
		home := os.Getenv("HOME")
		if home == "" {
			return "", fmt.Errorf("failed to get HOME directory")
		}
		configDir = filepath.Join(home, ".config")
	}
	return filepath.Join(configDir, "opencode-helper", "sources.json"), nil
}

// LegacyRegistryPath returns the path to the legacy config-sources.json file.
func LegacyRegistryPath() string {
	configDir := os.Getenv("XDG_CONFIG_HOME")
	if configDir == "" {
		home := os.Getenv("HOME")
		configDir = filepath.Join(home, ".config")
	}
	return filepath.Join(configDir, "opencode-helper", "config-sources.json")
}

// AppSupportRegistryPath returns the path to the macOS Application Support config.
func AppSupportRegistryPath() (string, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return "", fmt.Errorf("failed to get user config dir: %w", err)
	}
	return filepath.Join(configDir, "opencode-helper", "sources.json"), nil
}

// LoadRegistry loads the source registry from disk.
// If the new format doesn't exist, it attempts to migrate from legacy locations:
// 1. ~/.config/opencode-helper/config-sources.json (shell script legacy)
// 2. ~/Library/Application Support/opencode-helper/sources.json (pre-XDG Go CLI)
func LoadRegistry() (*Registry, error) {
	path, err := RegistryPath()
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			// Try legacy paths for migration
			return loadLegacyRegistry()
		}
		return nil, fmt.Errorf("failed to read registry: %w", err)
	}

	var registry Registry
	if err := json.Unmarshal(data, &registry); err != nil {
		return nil, fmt.Errorf("failed to parse registry: %w", err)
	}

	return &registry, nil
}

// LegacySource represents the legacy source format (without version field).
type LegacySource struct {
	ID        string     `json:"id"`
	Location  string     `json:"location"`
	Type      SourceType `json:"type"`
	Name      string     `json:"name"`
	CreatedAt string     `json:"added_at"`
}

// LegacyRegistry represents the legacy registry format.
type LegacyRegistry struct {
	Sources []LegacySource `json:"sources"`
}

// loadLegacyRegistry attempts to migrate from legacy config locations.
// It checks:
// 1. ~/.config/opencode-helper/config-sources.json (shell script legacy)
// 2. ~/Library/Application Support/opencode-helper/sources.json (pre-XDG Go CLI)
func loadLegacyRegistry() (*Registry, error) {
	// Try shell script legacy location first
	legacyPath := LegacyRegistryPath()
	if data, err := os.ReadFile(legacyPath); err == nil {
		if registry, err := migrateLegacyData(data, legacyPath); err == nil {
			return registry, nil
		}
	}

	// Try macOS Application Support location
	appSupportPath, err := AppSupportRegistryPath()
	if err != nil {
		return &Registry{Version: 1, Sources: []Source{}}, nil
	}
	if data, err := os.ReadFile(appSupportPath); err == nil {
		return migrateLegacyData(data, appSupportPath)
	}

	return &Registry{Version: 1, Sources: []Source{}}, nil
}

// migrateLegacyData migrates legacy JSON data (array or object format) to new Registry.
func migrateLegacyData(data []byte, legacyPath string) (*Registry, error) {
	// Try parsing as array first (legacy format)
	var legacySources []LegacySource
	if err := json.Unmarshal(data, &legacySources); err != nil {
		// Try parsing as object with "sources" field
		var legacy LegacyRegistry
		if err := json.Unmarshal(data, &legacy); err != nil {
			return nil, fmt.Errorf("failed to parse legacy registry: %w", err)
		}
		legacySources = legacy.Sources
	}

	// Migrate to new format
	registry := &Registry{Version: 1, Sources: make([]Source, len(legacySources))}
	for i, src := range legacySources {
		registry.Sources[i] = Source{
			ID:        src.ID,
			Location:  src.Location,
			Type:      src.Type,
			Name:      src.Name,
			CreatedAt: src.CreatedAt,
		}
	}

	// Migrate: save to new location
	if err := SaveRegistry(registry); err != nil {
		return nil, fmt.Errorf("failed to migrate registry: %w", err)
	}

	// Remove legacy file after successful migration
	if err := os.Remove(legacyPath); err != nil {
		// Log but don't fail - migration succeeded
		fmt.Printf("warning: failed to remove legacy config: %v\n", err)
	}

	return registry, nil
}

// SaveRegistry saves the source registry to disk.
func SaveRegistry(registry *Registry) error {
	path, err := RegistryPath()
	if err != nil {
		return err
	}

	// Ensure directory exists
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	data, err := json.MarshalIndent(registry, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal registry: %w", err)
	}

	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("failed to write registry: %w", err)
	}

	return nil
}

// DetectSourceType detects the source type from a location path.
func DetectSourceType(location string) (SourceType, error) {
	if isGitHubRef(location) {
		if _, err := parseGitHubRef(location); err != nil {
			return "", err
		}
		return SourceTypeGitHubRelease, nil
	}

	if isArchivePath(location) {
		return SourceTypeLocalArchive, nil
	}

	info, err := os.Stat(location)
	if err != nil {
		return "", fmt.Errorf("location does not exist: %s", location)
	}

	if info.IsDir() {
		return SourceTypeLocalDirectory, nil
	}

	return SourceTypeLocalArchive, nil
}

// isGitHubRef checks if a location appears to be a GitHub reference.
func isGitHubRef(location string) bool {
	if ownerRepoPattern.MatchString(location) && !strings.HasPrefix(location, "github.com/") {
		return true
	}

	githubPrefixes := []string{
		"https://github.com/",
		"http://github.com/",
		"github.com/",
		"git@github.com:",
	}
	for _, prefix := range githubPrefixes {
		if len(location) >= len(prefix) && location[:len(prefix)] == prefix {
			return true
		}
	}
	return false
}

func isArchivePath(location string) bool {
	return strings.HasSuffix(location, ".tar.gz") || strings.HasSuffix(location, ".tgz") || strings.HasSuffix(location, ".tar") || strings.HasSuffix(location, ".gz")
}

// ParseGitHubLocation returns a normalized repository and optional pinned tag.
func ParseGitHubLocation(location string) (GitHubRef, error) {
	return parseGitHubRef(location)
}

func parseGitHubRef(location string) (GitHubRef, error) {
	trimmed := strings.TrimSpace(location)
	if trimmed == "" {
		return GitHubRef{}, fmt.Errorf("invalid GitHub reference: location cannot be empty")
	}

	if ownerRepoPattern.MatchString(trimmed) && !strings.HasPrefix(trimmed, "github.com/") {
		return GitHubRef{Repo: trimmed}, nil
	}

	if strings.HasPrefix(trimmed, "git@github.com:") {
		trimmed = strings.TrimPrefix(trimmed, "git@github.com:")
		trimmed = strings.TrimSuffix(trimmed, ".git")
		if ownerRepoPattern.MatchString(trimmed) {
			return GitHubRef{Repo: trimmed}, nil
		}
		return GitHubRef{}, fmt.Errorf("invalid GitHub reference: %s", location)
	}

	if strings.HasPrefix(trimmed, "github.com/") {
		trimmed = "https://" + trimmed
	}

	parsed, err := url.Parse(trimmed)
	if err != nil {
		return GitHubRef{}, fmt.Errorf("invalid GitHub reference: %s", location)
	}
	if parsed.Host != "github.com" {
		return GitHubRef{}, fmt.Errorf("invalid GitHub reference: %s", location)
	}

	parts := strings.Split(strings.Trim(parsed.Path, "/"), "/")
	if len(parts) < 2 || parts[0] == "" || parts[1] == "" {
		return GitHubRef{}, fmt.Errorf("invalid GitHub reference: %s", location)
	}

	repo := parts[0] + "/" + strings.TrimSuffix(parts[1], ".git")
	if !ownerRepoPattern.MatchString(repo) {
		return GitHubRef{}, fmt.Errorf("invalid GitHub reference: %s", location)
	}

	ref := GitHubRef{Repo: repo}
	if len(parts) >= 5 && parts[2] == "releases" && parts[3] == "tag" && parts[4] != "" {
		ref.Tag = parts[4]
	}

	return ref, nil
}

// ValidateSource validates that a source location is accessible.
func ValidateSource(location string, sourceType SourceType) error {
	switch sourceType {
	case SourceTypeLocalDirectory:
		info, err := os.Stat(location)
		if err != nil {
			return fmt.Errorf("directory does not exist: %s", location)
		}
		if !info.IsDir() {
			return fmt.Errorf("location is not a directory: %s", location)
		}
		// Check for manifest file
		manifestPath := filepath.Join(location, "opencode-bundle.manifest.json")
		if _, err := os.Stat(manifestPath); os.IsNotExist(err) {
			return fmt.Errorf("bundle manifest not found: %s (run 'oc bundle init' to create one)", manifestPath)
		}
	case SourceTypeLocalArchive:
		info, err := os.Stat(location)
		if err != nil {
			return fmt.Errorf("archive does not exist: %s", location)
		}
		if info.IsDir() {
			return fmt.Errorf("archive must be a file, not a directory: %s", location)
		}
	case SourceTypeGitHubRelease:
		if _, err := parseGitHubRef(location); err != nil {
			return err
		}
	default:
		return fmt.Errorf("unknown source type: %s", sourceType)
	}
	return nil
}

func normalizeLocation(location string, sourceType SourceType) string {
	switch sourceType {
	case SourceTypeGitHubRelease:
		ref, err := parseGitHubRef(location)
		if err != nil {
			return location
		}
		return strings.ToLower(ref.Repo)
	case SourceTypeLocalDirectory:
		clean := filepath.Clean(location)
		if strings.HasPrefix(clean, "~/") {
			home, err := os.UserHomeDir()
			if err == nil {
				clean = filepath.Join(home, strings.TrimPrefix(clean, "~/"))
			}
		}
		abs, err := filepath.Abs(clean)
		if err != nil {
			return clean
		}
		return abs
	case SourceTypeLocalArchive:
		return filepath.Clean(location)
	default:
		return location
	}
}

func generateSourceID(normalizedLocation string) string {
	hash := sha256.Sum256([]byte(normalizedLocation))
	return hex.EncodeToString(hash[:])[:8]
}

// AddSource adds a new source to the registry.
func AddSource(location string, name string) (*Source, error) {
	// Detect source type
	sourceType, err := DetectSourceType(location)
	if err != nil {
		return nil, fmt.Errorf("failed to detect source type: %w", err)
	}

	// Validate the source
	if err := ValidateSource(location, sourceType); err != nil {
		return nil, err
	}

	// Load existing registry
	registry, err := LoadRegistry()
	if err != nil {
		return nil, err
	}

	// Normalize location for consistent ID generation and duplicate check
	normalizedLocation := normalizeLocation(location, sourceType)

	// Check for duplicate locations (compare normalized)
	for _, s := range registry.Sources {
		if normalizeLocation(s.Location, s.Type) == normalizedLocation {
			return nil, fmt.Errorf("source already registered at location: %s (id: %s)", location, s.ID)
		}
	}

	// Generate ID from normalized location
	id := generateSourceID(normalizedLocation)
	if name == "" {
		if sourceType == SourceTypeGitHubRelease {
			ref, err := parseGitHubRef(location)
			if err != nil {
				return nil, err
			}
			name = filepath.Base(ref.Repo)
		} else {
			name = filepath.Base(location)
		}
	}

	// Create new source
	source := Source{
		ID:        id,
		Location:  location,
		Type:      sourceType,
		Name:      name,
		CreatedAt: time.Now().UTC().Format(time.RFC3339),
	}

	registry.Sources = append(registry.Sources, source)

	if err := SaveRegistry(registry); err != nil {
		return nil, err
	}

	return &source, nil
}

// ListSources returns all registered sources.
func ListSources() ([]Source, error) {
	registry, err := LoadRegistry()
	if err != nil {
		return nil, err
	}
	return registry.Sources, nil
}

// RemoveSource removes a source from the registry by ID.
func RemoveSource(id string) error {
	registry, err := LoadRegistry()
	if err != nil {
		return err
	}

	found := false
	newSources := make([]Source, 0, len(registry.Sources))
	for _, s := range registry.Sources {
		if s.ID == id {
			found = true
		} else {
			newSources = append(newSources, s)
		}
	}

	if !found {
		return fmt.Errorf("source not found: %s", id)
	}

	registry.Sources = newSources
	return SaveRegistry(registry)
}

// GetSource returns a source by ID.
func GetSource(id string) (*Source, error) {
	registry, err := LoadRegistry()
	if err != nil {
		return nil, err
	}

	for _, s := range registry.Sources {
		if s.ID == id {
			return &s, nil
		}
	}

	return nil, fmt.Errorf("source not found: %s", id)
}

// ResolveSourceRef resolves a source reference by exact ID or unique name.
func ResolveSourceRef(ref string) (*Source, error) {
	registry, err := LoadRegistry()
	if err != nil {
		return nil, err
	}

	trimmed := strings.TrimSpace(ref)
	if trimmed == "" {
		return nil, fmt.Errorf("source reference cannot be empty")
	}

	for _, s := range registry.Sources {
		if s.ID == trimmed {
			return &s, nil
		}
	}

	var matches []Source
	for _, s := range registry.Sources {
		if s.Name == trimmed {
			matches = append(matches, s)
		}
	}

	switch len(matches) {
	case 0:
		return nil, fmt.Errorf("source not found: %s", ref)
	case 1:
		return &matches[0], nil
	default:
		sort.Slice(matches, func(i, j int) bool {
			return matches[i].ID < matches[j].ID
		})
		return nil, &AmbiguousSourceRefError{Ref: trimmed, Matches: matches}
	}
}
