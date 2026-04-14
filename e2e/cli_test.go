package e2e

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
	"testing"
	"time"

	"github.com/creack/pty"
)

type commandResult struct {
	stdout string
	stderr string
	err    error
}

type provenance struct {
	SourceID      string `json:"source_id"`
	SourceName    string `json:"source_name"`
	SourceType    string `json:"source_type"`
	BundleVersion string `json:"bundle_version"`
	PresetName    string `json:"preset_name"`
	Entrypoint    string `json:"entrypoint"`
	AppliedAt     string `json:"applied_at"`
}

func TestVersion(t *testing.T) {
	result := runOC(t, testEnv(t), "version")
	requireSuccess(t, result)
	if !strings.HasPrefix(result.stdout, "oc ") {
		t.Fatalf("expected version output, got stdout=%q stderr=%q", result.stdout, result.stderr)
	}
}

func TestLocalDirectoryFlow(t *testing.T) {
	env := testEnv(t)
	bundleDir := copyFixtureBundle(t)
	projectRoot := t.TempDir()

	addResult := runOC(t, env, "source", "add", bundleDir, "--name", "fixture-dir")
	requireSuccess(t, addResult)
	sourceID := extractSourceID(t, addResult.stdout)

	listResult := runOC(t, env, "source", "list")
	requireSuccess(t, listResult)
	requireContains(t, listResult.stdout, sourceID)
	requireContains(t, listResult.stdout, "fixture-dir")
	requireContains(t, listResult.stdout, "local-directory")

	applyResult := runOC(t, env, "bundle", "apply", sourceID, "--preset", "fixture", "--project-root", projectRoot)
	requireSuccess(t, applyResult)

	configPath := filepath.Join(projectRoot, "opencode.json")
	configData, err := os.ReadFile(configPath)
	if err != nil {
		t.Fatalf("failed to read applied config: %v", err)
	}
	requireContains(t, string(configData), `"mode": "fixture"`)

	provenancePath := filepath.Join(projectRoot, ".opencode", "bundle-provenance.json")
	prov := readProvenance(t, provenancePath)
	if prov.SourceID != sourceID {
		t.Fatalf("expected source id %q, got %q", sourceID, prov.SourceID)
	}
	if prov.SourceType != "local-directory" {
		t.Fatalf("expected source type local-directory, got %q", prov.SourceType)
	}
	if prov.PresetName != "fixture" {
		t.Fatalf("expected preset fixture, got %q", prov.PresetName)
	}

	statusResult := runOC(t, env, "bundle", "status", "--project-root", projectRoot)
	requireSuccess(t, statusResult)
	requireContains(t, statusResult.stdout, "Bundle Provenance:")
	requireContains(t, statusResult.stdout, sourceID)
	requireContains(t, statusResult.stdout, "fixture-dir")
	requireContains(t, statusResult.stdout, "fixture")

	overwriteResult := runOC(t, env, "bundle", "apply", sourceID, "--preset", "fixture", "--project-root", projectRoot)
	requireFailure(t, overwriteResult)
	requireContains(t, overwriteResult.stderr, "output file exists")

	forceResult := runOC(t, env, "bundle", "apply", sourceID, "--preset", "fixture", "--project-root", projectRoot, "--force")
	requireSuccess(t, forceResult)
	updatedProv := readProvenance(t, provenancePath)
	if updatedProv.SourceID != sourceID {
		t.Fatalf("expected forced apply to preserve source id %q, got %q", sourceID, updatedProv.SourceID)
	}
	if updatedProv.SourceName != "fixture-dir" {
		t.Fatalf("expected forced apply source name fixture-dir, got %q", updatedProv.SourceName)
	}
}

func TestLocalDirectoryApplyBySourceName(t *testing.T) {
	env := testEnv(t)
	bundleDir := copyFixtureBundle(t)
	projectRoot := t.TempDir()

	addResult := runOC(t, env, "source", "add", bundleDir, "--name", "fixture-dir")
	requireSuccess(t, addResult)

	applyResult := runOC(t, env, "bundle", "apply", "fixture-dir", "--preset", "fixture", "--project-root", projectRoot)
	requireSuccess(t, applyResult)

	prov := readProvenance(t, filepath.Join(projectRoot, ".opencode", "bundle-provenance.json"))
	if prov.SourceName != "fixture-dir" {
		t.Fatalf("expected source name fixture-dir, got %q", prov.SourceName)
	}
	if prov.PresetName != "fixture" {
		t.Fatalf("expected preset fixture, got %q", prov.PresetName)
	}
}

func TestLocalArchiveFlow(t *testing.T) {
	env := testEnv(t)
	bundleDir := copyFixtureBundle(t)
	archivePath := filepath.Join(t.TempDir(), "fixture-bundle.tar.gz")
	createTarGzFromDir(t, bundleDir, archivePath)
	projectRoot := t.TempDir()

	addResult := runOC(t, env, "source", "add", archivePath, "--name", "fixture-archive")
	requireSuccess(t, addResult)
	sourceID := extractSourceID(t, addResult.stdout)

	applyResult := runOC(t, env, "bundle", "apply", sourceID, "--preset", "fixture", "--project-root", projectRoot)
	requireSuccess(t, applyResult)

	provenancePath := filepath.Join(projectRoot, ".opencode", "bundle-provenance.json")
	prov := readProvenance(t, provenancePath)
	if prov.SourceType != "local-archive" {
		t.Fatalf("expected source type local-archive, got %q", prov.SourceType)
	}
	if prov.SourceName != "fixture-archive" {
		t.Fatalf("expected source name fixture-archive, got %q", prov.SourceName)
	}
}

func TestGitHubReleaseFlow(t *testing.T) {
	env := testEnv(t)
	projectRoot := t.TempDir()

	server := newGitHubReleaseE2EServer(t, githubReleaseE2EFixture{
		repo: "owner/repo",
		releases: []githubReleaseE2ERelease{
			newGitHubReleaseE2ERelease(t, "v1.2.3", false, "github-fixture"),
		},
	})
	defer server.Close()

	env = append(env, "OC_GITHUB_API_BASE_URL="+server.URL)

	addResult := runOC(t, env, "source", "add", "owner/repo", "--name", "fixture-github")
	requireSuccess(t, addResult)
	sourceID := extractSourceID(t, addResult.stdout)

	applyResult := runOC(t, env, "bundle", "apply", sourceID, "--version", "v1.2.3", "--preset", "fixture", "--project-root", projectRoot)
	requireSuccess(t, applyResult)

	configPath := filepath.Join(projectRoot, "opencode.json")
	configData, err := os.ReadFile(configPath)
	if err != nil {
		t.Fatalf("failed to read applied config: %v", err)
	}
	requireContains(t, string(configData), `"mode":"github-fixture"`)

	prov := readProvenance(t, filepath.Join(projectRoot, ".opencode", "bundle-provenance.json"))
	if prov.SourceType != "github-release" {
		t.Fatalf("expected source type github-release, got %q", prov.SourceType)
	}
	if prov.BundleVersion != "v1.2.3" {
		t.Fatalf("expected bundle version v1.2.3, got %q", prov.BundleVersion)
	}
	if prov.SourceName != "fixture-github" {
		t.Fatalf("expected source name fixture-github, got %q", prov.SourceName)
	}
}

func TestGitHubReleaseInteractiveVersionSelectionFlow(t *testing.T) {
	env := testEnv(t)
	projectRoot := t.TempDir()
	server := newGitHubReleaseE2EServer(t, githubReleaseE2EFixture{
		repo: "owner/repo",
		releases: []githubReleaseE2ERelease{
			newGitHubReleaseE2ERelease(t, "v1.3.0-alpha.1", true, "github-prerelease"),
			newGitHubReleaseE2ERelease(t, "v1.2.3", false, "github-stable"),
		},
	})
	defer server.Close()

	env = append(env, "OC_GITHUB_API_BASE_URL="+server.URL)

	addResult := runOC(t, env, "source", "add", "owner/repo", "--name", "fixture-github")
	requireSuccess(t, addResult)
	sourceID := extractSourceID(t, addResult.stdout)

	applyResult := runOCInPTY(t, env, "1\n", "bundle", "apply", sourceID, "--preset", "fixture", "--project-root", projectRoot)
	requireSuccess(t, applyResult)
	requireContains(t, applyResult.stdout, "Available versions for owner/repo:")
	requireContains(t, applyResult.stdout, "v1.3.0-alpha.1 (prerelease)")

	configData, err := os.ReadFile(filepath.Join(projectRoot, "opencode.json"))
	if err != nil {
		t.Fatalf("failed to read applied config: %v", err)
	}
	requireContains(t, string(configData), `"mode":"github-prerelease"`)

	prov := readProvenance(t, filepath.Join(projectRoot, ".opencode", "bundle-provenance.json"))
	if prov.BundleVersion != "v1.3.0-alpha.1" {
		t.Fatalf("expected bundle version v1.3.0-alpha.1, got %q", prov.BundleVersion)
	}
	if prov.PresetName != "fixture" {
		t.Fatalf("expected preset fixture, got %q", prov.PresetName)
	}
}

func TestGitHubReleaseApplyWithoutVersionNonInteractiveFails(t *testing.T) {
	env := testEnv(t)
	projectRoot := t.TempDir()
	server := newGitHubReleaseE2EServer(t, githubReleaseE2EFixture{
		repo: "owner/repo",
		releases: []githubReleaseE2ERelease{
			newGitHubReleaseE2ERelease(t, "v1.3.0-alpha.1", true, "github-prerelease"),
			newGitHubReleaseE2ERelease(t, "v1.2.3", false, "github-stable"),
		},
	})
	defer server.Close()

	env = append(env, "OC_GITHUB_API_BASE_URL="+server.URL)

	addResult := runOC(t, env, "source", "add", "owner/repo", "--name", "fixture-github")
	requireSuccess(t, addResult)
	sourceID := extractSourceID(t, addResult.stdout)

	applyResult := runOCWithStdin(t, env, strings.NewReader(""), "bundle", "apply", sourceID, "--preset", "fixture", "--project-root", projectRoot)
	requireFailure(t, applyResult)
	requireContains(t, applyResult.stderr, "--version is required for github-release sources outside interactive mode")
	if strings.Contains(applyResult.stderr, "Select a version") || strings.Contains(applyResult.stdout, "Select a version") {
		t.Fatalf("unexpected interactive prompt in non-interactive flow: stdout=%q stderr=%q", applyResult.stdout, applyResult.stderr)
	}
}

func TestBundleApplyFailsForUnknownSource(t *testing.T) {
	projectRoot := t.TempDir()
	result := runOC(t, testEnv(t), "bundle", "apply", "missing-id", "--preset", "fixture", "--project-root", projectRoot)
	requireFailure(t, result)
	requireContains(t, result.stderr, "source not found")
}

func TestBundleApplyRequiresPresetOutsideTTY(t *testing.T) {
	env := testEnv(t)
	bundleDir := copyFixtureBundle(t)
	projectRoot := t.TempDir()

	addResult := runOC(t, env, "source", "add", bundleDir, "--name", "fixture-dir")
	requireSuccess(t, addResult)

	applyResult := runOCWithStdin(t, env, strings.NewReader(""), "bundle", "apply", "fixture-dir", "--project-root", projectRoot)
	requireFailure(t, applyResult)
	requireContains(t, applyResult.stderr, "--preset is required outside interactive mode")
}

func TestBundleApplyNoArgsFailsInNonTTY(t *testing.T) {
	// When run without arguments in non-TTY (e2e tests), should fail with helpful message
	// Use empty stdin to ensure no TTY is detected
	projectRoot := t.TempDir()
	result := runOCWithStdin(t, testEnv(t), strings.NewReader(""), "bundle", "apply", "--project-root", projectRoot)
	requireFailure(t, result)
	// When there are no sources, it should fail with "no sources registered"
	// OR when there are sources but no TTY, fail with "source-ref is required"
	requireContains(t, result.stderr, "source-ref is required")
	requireContains(t, result.stderr, "non-interactive mode")
}

func TestBundleApplyAutoFlagRequiresSourceRef(t *testing.T) {
	// --auto flag should require source-ref argument regardless of TTY
	env := testEnv(t)
	bundleDir := copyFixtureBundle(t)

	addResult := runOC(t, env, "source", "add", bundleDir, "--name", "fixture-dir")
	requireSuccess(t, addResult)

	projectRoot := t.TempDir()
	// Using --auto without source-ref should fail
	result := runOCWithStdin(t, env, strings.NewReader(""), "bundle", "apply", "--auto", "--project-root", projectRoot)
	requireFailure(t, result)
	requireContains(t, result.stderr, "source-ref is required")
	requireContains(t, result.stderr, "--auto")
}

func TestBundleApplyAutoFlagWithPresetRequiresSource(t *testing.T) {
	// --auto with --preset but no source-ref should fail
	env := testEnv(t)
	bundleDir := copyFixtureBundle(t)

	addResult := runOC(t, env, "source", "add", bundleDir, "--name", "fixture-dir")
	requireSuccess(t, addResult)

	projectRoot := t.TempDir()
	result := runOCWithStdin(t, env, strings.NewReader(""), "bundle", "apply", "--auto", "--preset", "fixture", "--project-root", projectRoot)
	requireFailure(t, result)
	requireContains(t, result.stderr, "source-ref is required")
}

func TestSourceAddFailsWithoutManifest(t *testing.T) {
	bundleDir := t.TempDir()
	result := runOC(t, testEnv(t), "source", "add", bundleDir)
	requireFailure(t, result)
	requireContains(t, result.stderr, "bundle manifest not found")
}

func TestInvalidTarballFailsOnApply(t *testing.T) {
	env := testEnv(t)
	archivePath := filepath.Join(t.TempDir(), "invalid.tar.gz")
	if err := os.WriteFile(archivePath, []byte("not a tarball"), 0o644); err != nil {
		t.Fatalf("failed to write invalid archive: %v", err)
	}

	addResult := runOC(t, env, "source", "add", archivePath, "--name", "broken-archive")
	requireSuccess(t, addResult)
	sourceID := extractSourceID(t, addResult.stdout)

	projectRoot := t.TempDir()
	applyResult := runOC(t, env, "bundle", "apply", sourceID, "--preset", "fixture", "--project-root", projectRoot)
	requireFailure(t, applyResult)
	requireContains(t, applyResult.stderr, "failed to resolve source")
	if runtime.GOOS == "darwin" {
		requireContains(t, applyResult.stderr, "failed to extract tarball")
		return
	}
	requireContains(t, applyResult.stderr, "failed to extract tarball")
}

func testEnv(t *testing.T) []string {
	t.Helper()
	homeDir := t.TempDir()
	configHome := filepath.Join(homeDir, ".config")
	if err := os.MkdirAll(configHome, 0o755); err != nil {
		t.Fatalf("failed to create config home: %v", err)
	}

	pathValue := os.Getenv("PATH")
	if pathValue == "" {
		t.Fatal("PATH is required for subprocess execution")
	}

	return []string{
		"HOME=" + homeDir,
		"XDG_CONFIG_HOME=" + configHome,
		"PATH=" + pathValue,
	}
}

func runOC(t *testing.T, env []string, args ...string) commandResult {
	t.Helper()
	return runOCWithStdin(t, env, nil, args...)
}

func runOCInPTY(t *testing.T, env []string, input string, args ...string) commandResult {
	t.Helper()
	binaryPath := os.Getenv("OC_E2E_BINARY")
	if binaryPath == "" {
		t.Skip("OC_E2E_BINARY not set; skipping black-box CLI E2E tests")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, binaryPath, args...)
	cmd.Env = env
	ptmx, err := pty.Start(cmd)
	if err != nil {
		t.Fatalf("failed to start PTY command: %v", err)
	}
	defer ptmx.Close()

	var output bytes.Buffer
	readDone := make(chan error, 1)
	go func() {
		_, copyErr := io.Copy(&output, ptmx)
		readDone <- copyErr
	}()

	if input != "" {
		if _, err := ptmx.Write([]byte(input)); err != nil {
			t.Fatalf("failed to write PTY input: %v", err)
		}
	}
	err = cmd.Wait()
	_ = ptmx.Close()
	<-readDone

	return commandResult{stdout: output.String(), stderr: output.String(), err: err}
}

func runOCWithStdin(t *testing.T, env []string, stdin *strings.Reader, args ...string) commandResult {
	t.Helper()
	binaryPath := os.Getenv("OC_E2E_BINARY")
	if binaryPath == "" {
		t.Skip("OC_E2E_BINARY not set; skipping black-box CLI E2E tests")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, binaryPath, args...)
	cmd.Env = env
	if stdin != nil {
		cmd.Stdin = stdin
	}
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()

	return commandResult{
		stdout: stdout.String(),
		stderr: stderr.String(),
		err:    err,
	}
}

func requireSuccess(t *testing.T, result commandResult) {
	t.Helper()
	if result.err != nil {
		t.Fatalf("expected success, got err=%v stdout=%q stderr=%q", result.err, result.stdout, result.stderr)
	}
}

func requireFailure(t *testing.T, result commandResult) {
	t.Helper()
	if result.err == nil {
		t.Fatalf("expected failure, got stdout=%q stderr=%q", result.stdout, result.stderr)
	}
}

func requireContains(t *testing.T, haystack, needle string) {
	t.Helper()
	if !strings.Contains(haystack, needle) {
		t.Fatalf("expected %q to contain %q", haystack, needle)
	}
}

func extractSourceID(t *testing.T, stdout string) string {
	t.Helper()
	for _, line := range strings.Split(stdout, "\n") {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "ID:") {
			return strings.TrimSpace(strings.TrimPrefix(trimmed, "ID:"))
		}
	}
	t.Fatalf("failed to extract source id from stdout=%q", stdout)
	return ""
}

func readProvenance(t *testing.T, path string) provenance {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("failed to read provenance file: %v", err)
	}

	var prov provenance
	if err := json.Unmarshal(data, &prov); err != nil {
		t.Fatalf("failed to parse provenance file: %v", err)
	}
	return prov
}

func copyFixtureBundle(t *testing.T) string {
	t.Helper()
	sourceRoot := filepath.Join("testdata", "fixture-bundle")
	destRoot := filepath.Join(t.TempDir(), "fixture-bundle")
	copyDir(t, sourceRoot, destRoot)
	return destRoot
}

func copyDir(t *testing.T, srcDir, destDir string) {
	t.Helper()
	entries, err := os.ReadDir(srcDir)
	if err != nil {
		t.Fatalf("failed to read fixture dir %s: %v", srcDir, err)
	}
	if err := os.MkdirAll(destDir, 0o755); err != nil {
		t.Fatalf("failed to create fixture dir %s: %v", destDir, err)
	}

	for _, entry := range entries {
		srcPath := filepath.Join(srcDir, entry.Name())
		destPath := filepath.Join(destDir, entry.Name())
		if entry.IsDir() {
			copyDir(t, srcPath, destPath)
			continue
		}

		data, err := os.ReadFile(srcPath)
		if err != nil {
			t.Fatalf("failed to read fixture file %s: %v", srcPath, err)
		}
		if err := os.WriteFile(destPath, data, 0o644); err != nil {
			t.Fatalf("failed to write fixture file %s: %v", destPath, err)
		}
	}
}

func createTarGzFromDir(t *testing.T, sourceDir, archivePath string) {
	t.Helper()
	archiveFile, err := os.Create(archivePath)
	if err != nil {
		t.Fatalf("failed to create archive %s: %v", archivePath, err)
	}
	defer archiveFile.Close()

	gzipWriter := gzip.NewWriter(archiveFile)
	defer gzipWriter.Close()

	tarWriter := tar.NewWriter(gzipWriter)
	defer tarWriter.Close()

	if err := filepath.Walk(sourceDir, func(path string, info os.FileInfo, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if info.IsDir() {
			return nil
		}

		relPath, err := filepath.Rel(sourceDir, path)
		if err != nil {
			return err
		}
		content, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		header := &tar.Header{
			Name: filepath.ToSlash(filepath.Join(filepath.Base(sourceDir), relPath)),
			Mode: 0o644,
			Size: int64(len(content)),
		}
		if err := tarWriter.WriteHeader(header); err != nil {
			return err
		}
		if _, err := tarWriter.Write(content); err != nil {
			return err
		}
		return nil
	}); err != nil {
		t.Fatalf("failed to build tar.gz fixture: %v", err)
	}
	if err := tarWriter.Close(); err != nil {
		t.Fatalf("failed to finalize tar writer: %v", err)
	}
	if err := gzipWriter.Close(); err != nil {
		t.Fatalf("failed to finalize gzip writer: %v", err)
	}
	if err := archiveFile.Close(); err != nil {
		t.Fatalf("failed to close archive file: %v", err)
	}
}

type githubReleaseE2EFixture struct {
	repo     string
	releases []githubReleaseE2ERelease
}

type githubReleaseE2ERelease struct {
	tag          string
	prerelease   bool
	archiveName  string
	archiveBytes []byte
	checksums    string
}

func newGitHubReleaseE2EServer(t *testing.T, fixture githubReleaseE2EFixture) *httptest.Server {
	t.Helper()
	mux := http.NewServeMux()
	mux.HandleFunc("/repos/"+fixture.repo+"/releases", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(buildReleaseListResponse(r.Host, fixture.releases))
	})
	mux.HandleFunc("/repos/"+fixture.repo+"/releases/latest", func(w http.ResponseWriter, r *http.Request) {
		for _, release := range fixture.releases {
			if release.prerelease {
				continue
			}
			writeReleaseResponse(w, r, fixture.repo, release)
			return
		}
		http.NotFound(w, r)
	})
	for _, release := range fixture.releases {
		release := release
		mux.HandleFunc("/repos/"+fixture.repo+"/releases/tags/"+release.tag, func(w http.ResponseWriter, r *http.Request) {
			writeReleaseResponse(w, r, fixture.repo, release)
		})
		mux.HandleFunc("/downloads/"+fixture.repo+"/releases/download/"+release.tag+"/"+release.archiveName, func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/gzip")
			_, _ = w.Write(release.archiveBytes)
		})
		mux.HandleFunc("/downloads/"+fixture.repo+"/releases/download/"+release.tag+"/opencode-config-bundle-"+release.tag+"-checksums.txt", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "text/plain")
			_, _ = w.Write([]byte(release.checksums))
		})
	}
	return httptest.NewServer(mux)
}

func writeReleaseResponse(w http.ResponseWriter, r *http.Request, repo string, release githubReleaseE2ERelease) {
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]any{
		"tag_name":   release.tag,
		"prerelease": release.prerelease,
		"assets": []map[string]string{
			{
				"name":                 release.archiveName,
				"browser_download_url": fmt.Sprintf("http://%s/downloads/%s/releases/download/%s/%s", r.Host, repo, release.tag, release.archiveName),
			},
			{
				"name":                 "opencode-config-bundle-" + release.tag + "-checksums.txt",
				"browser_download_url": fmt.Sprintf("http://%s/downloads/%s/releases/download/%s/%s", r.Host, repo, release.tag, "opencode-config-bundle-"+release.tag+"-checksums.txt"),
			},
		},
	})
}

func buildReleaseListResponse(host string, releases []githubReleaseE2ERelease) []map[string]any {
	responses := make([]map[string]any, 0, len(releases))
	for _, release := range releases {
		responses = append(responses, map[string]any{
			"tag_name":   release.tag,
			"prerelease": release.prerelease,
			"assets": []map[string]string{{
				"name":                 release.archiveName,
				"browser_download_url": fmt.Sprintf("http://%s/downloads/unused/releases/download/%s/%s", host, release.tag, release.archiveName),
			}},
		})
	}
	sort.SliceStable(responses, func(i, j int) bool { return i < j })
	return responses
}

func newGitHubReleaseE2ERelease(t *testing.T, tag string, prerelease bool, mode string) githubReleaseE2ERelease {
	t.Helper()
	bundleDir := filepath.Join(t.TempDir(), "opencode-config-bundle-"+tag)
	if err := os.MkdirAll(bundleDir, 0o755); err != nil {
		t.Fatalf("failed to create bundle dir: %v", err)
	}
	if err := os.WriteFile(filepath.Join(bundleDir, "opencode-bundle.manifest.json"), []byte(fmt.Sprintf(`{
		"manifest_version": "1.0.0",
		"bundle_name": "fixture-github",
		"bundle_version": %q,
		"presets": [
			{"name": "fixture", "entrypoint": "opencode.json"}
		]
	}`, tag)), 0o644); err != nil {
		t.Fatalf("failed to write manifest: %v", err)
	}
	if err := os.WriteFile(filepath.Join(bundleDir, "opencode.json"), []byte(fmt.Sprintf(`{"mode":"%s"}`, mode)), 0o644); err != nil {
		t.Fatalf("failed to write preset: %v", err)
	}

	archiveName := "opencode-config-bundle-" + tag + ".tar.gz"
	archivePath := filepath.Join(t.TempDir(), archiveName)
	createTarGzFromDir(t, bundleDir, archivePath)
	archiveData, err := os.ReadFile(archivePath)
	if err != nil {
		t.Fatalf("failed to read archive: %v", err)
	}
	archiveSHA := fmt.Sprintf("%x", sha256.Sum256(archiveData))

	return githubReleaseE2ERelease{
		tag:          tag,
		prerelease:   prerelease,
		archiveName:  archiveName,
		archiveBytes: archiveData,
		checksums:    fmt.Sprintf("%s  %s\n", archiveSHA, archiveName),
	}
}

func TestPresetListSourcesGitHubWithStableRelease(t *testing.T) {
	env := testEnv(t)
	server := newGitHubReleaseE2EServer(t, githubReleaseE2EFixture{
		repo: "owner/repo",
		releases: []githubReleaseE2ERelease{
			newGitHubReleaseE2ERelease(t, "v1.3.0-alpha.1", true, "github-prerelease"),
			newGitHubReleaseE2ERelease(t, "v1.2.3", false, "github-stable"),
		},
	})
	defer server.Close()

	env = append(env, "OC_GITHUB_API_BASE_URL="+server.URL)

	addResult := runOC(t, env, "source", "add", "owner/repo", "--name", "test-stable")
	requireSuccess(t, addResult)

	listResult := runOC(t, env, "source", "list", "--with-presets")
	requireSuccess(t, listResult)
	requireContains(t, listResult.stdout, "v1.2.3")
	requireContains(t, listResult.stdout, "fixture")
}

func TestPresetListSourcesGitHubPrereleaseOnlyWarns(t *testing.T) {
	env := testEnv(t)
	server := newGitHubReleaseE2EServer(t, githubReleaseE2EFixture{
		repo: "owner/repo",
		releases: []githubReleaseE2ERelease{
			newGitHubReleaseE2ERelease(t, "v2.0.0-alpha.1", true, "github-prerelease"),
		},
	})
	defer server.Close()

	env = append(env, "OC_GITHUB_API_BASE_URL="+server.URL)

	addResult := runOC(t, env, "source", "add", "owner/repo", "--name", "test-prerelease")
	requireSuccess(t, addResult)

	listResult := runOCWithStdin(t, env, strings.NewReader(""), "source", "list", "--with-presets")
	requireFailure(t, listResult)
	requireContains(t, listResult.stderr, "no stable release found")
}

// Bundle Init E2E Tests for US-053
func TestBundleInitNonInteractiveCreatesValidBundle(t *testing.T) {
	outputDir := t.TempDir()
	bundleDir := filepath.Join(outputDir, "test-bundle")

	result := runOCWithStdin(t, testEnv(t), strings.NewReader(""), "bundle", "init", "--name", "my-test-bundle", "--output", bundleDir)
	requireSuccess(t, result)
	requireContains(t, result.stdout, "done: bundle initialized")

	// Verify manifest was created
	manifestPath := filepath.Join(bundleDir, "opencode-bundle.manifest.json")
	manifestData, err := os.ReadFile(manifestPath)
	if err != nil {
		t.Fatalf("failed to read manifest: %v", err)
	}
	requireContains(t, string(manifestData), `"manifest_version": "1.0.0"`)
	requireContains(t, string(manifestData), `"bundle_name": "my-test-bundle"`)
	requireContains(t, string(manifestData), `"bundle_version": "0.0.1"`)
	requireContains(t, string(manifestData), `"name": "default"`)

	// Verify preset was created at root level (per bundle contract)
	presetPath := filepath.Join(bundleDir, "default.json")
	_, err = os.ReadFile(presetPath)
	if err != nil {
		t.Fatalf("failed to read preset: %v", err)
	}

	// Verify README was created
	readmePath := filepath.Join(bundleDir, "README.md")
	readmeData, err := os.ReadFile(readmePath)
	if err != nil {
		t.Fatalf("failed to read README: %v", err)
	}
	requireContains(t, string(readmeData), "# my-test-bundle")
}

func TestBundleInitWithCustomVersion(t *testing.T) {
	outputDir := t.TempDir()
	bundleDir := filepath.Join(outputDir, "test-bundle")

	result := runOCWithStdin(t, testEnv(t), strings.NewReader(""), "bundle", "init", "--name", "custom-bundle", "--version", "1.2.3", "--output", bundleDir)
	requireSuccess(t, result)

	// Verify manifest has custom version
	manifestPath := filepath.Join(bundleDir, "opencode-bundle.manifest.json")
	manifestData, err := os.ReadFile(manifestPath)
	if err != nil {
		t.Fatalf("failed to read manifest: %v", err)
	}
	requireContains(t, string(manifestData), `"bundle_version": "1.2.3"`)
}

func TestBundleInitFailsWithoutNameInNonTTY(t *testing.T) {
	outputDir := t.TempDir()
	bundleDir := filepath.Join(outputDir, "test-bundle")

	result := runOCWithStdin(t, testEnv(t), strings.NewReader(""), "bundle", "init", "--output", bundleDir)
	requireFailure(t, result)
	requireContains(t, result.stderr, "--name is required in non-interactive mode")
}

func TestBundleInitFailsIfDirectoryExists(t *testing.T) {
	outputDir := t.TempDir()
	bundleDir := filepath.Join(outputDir, "test-bundle")

	// Create the directory first
	if err := os.MkdirAll(bundleDir, 0755); err != nil {
		t.Fatalf("failed to create directory: %v", err)
	}

	result := runOCWithStdin(t, testEnv(t), strings.NewReader(""), "bundle", "init", "--name", "my-bundle", "--output", bundleDir)
	requireFailure(t, result)
	requireContains(t, result.stderr, "already exists")
}

func TestBundleInitForceOverwritesExistingDirectory(t *testing.T) {
	outputDir := t.TempDir()
	bundleDir := filepath.Join(outputDir, "test-bundle")

	// Create the directory with a conflicting manifest first
	if err := os.MkdirAll(bundleDir, 0755); err != nil {
		t.Fatalf("failed to create directory: %v", err)
	}
	if err := os.WriteFile(filepath.Join(bundleDir, "opencode-bundle.manifest.json"), []byte("old manifest"), 0644); err != nil {
		t.Fatalf("failed to create old manifest: %v", err)
	}

	result := runOCWithStdin(t, testEnv(t), strings.NewReader(""), "bundle", "init", "--name", "force-bundle", "--output", bundleDir, "--force")
	requireSuccess(t, result)
	requireContains(t, result.stdout, "done: bundle initialized")

	// Verify manifest was overwritten
	manifestPath := filepath.Join(bundleDir, "opencode-bundle.manifest.json")
	manifestData, err := os.ReadFile(manifestPath)
	if err != nil {
		t.Fatalf("manifest should exist after force overwrite: %v", err)
	}
	requireContains(t, string(manifestData), `"bundle_name": "force-bundle"`)
	// Old content should be replaced
	if strings.Contains(string(manifestData), "old manifest") {
		t.Fatalf("old manifest content should have been overwritten")
	}
}

func TestBundleInitCreatesNestedDirectories(t *testing.T) {
	outputDir := t.TempDir()
	// Use a nested path that doesn't exist
	bundleDir := filepath.Join(outputDir, "nested", "deep", "bundle")

	result := runOCWithStdin(t, testEnv(t), strings.NewReader(""), "bundle", "init", "--name", "nested-bundle", "--output", bundleDir)
	requireSuccess(t, result)

	// Verify manifest was created in nested directory
	manifestPath := filepath.Join(bundleDir, "opencode-bundle.manifest.json")
	_, err := os.ReadFile(manifestPath)
	if err != nil {
		t.Fatalf("failed to read manifest in nested directory: %v", err)
	}
}

func TestBundleInitInvalidBundleName(t *testing.T) {
	outputDir := t.TempDir()
	bundleDir := filepath.Join(outputDir, "test-bundle")

	// Invalid name with special characters
	result := runOCWithStdin(t, testEnv(t), strings.NewReader(""), "bundle", "init", "--name", "my@invalid#bundle", "--output", bundleDir)
	requireFailure(t, result)
	requireContains(t, result.stderr, "bundle name must contain")
}

func TestBundleInitGeneratedBundleIsValid(t *testing.T) {
	env := testEnv(t)
	outputDir := t.TempDir()
	bundleDir := filepath.Join(outputDir, "valid-bundle")

	// Create the bundle
	initResult := runOCWithStdin(t, env, strings.NewReader(""), "bundle", "init", "--name", "valid-bundle", "--output", bundleDir)
	requireSuccess(t, initResult)

	// Add the bundle as a source
	addResult := runOCWithStdin(t, env, strings.NewReader(""), "source", "add", bundleDir, "--name", "test-valid")
	requireSuccess(t, addResult)

	// Try to apply the bundle
	projectRoot := t.TempDir()
	applyResult := runOCWithStdin(t, env, strings.NewReader(""), "bundle", "apply", "test-valid", "--preset", "default", "--project-root", projectRoot)
	requireSuccess(t, applyResult)

	// Verify config was applied
	configPath := filepath.Join(projectRoot, "opencode.json")
	_, err := os.ReadFile(configPath)
	if err != nil {
		t.Fatalf("config should have been applied: %v", err)
	}
}

func TestBundleInitInteractiveMode(t *testing.T) {
	outputDir := t.TempDir()
	bundleDir := filepath.Join(outputDir, "interactive-bundle")

	// Interactive mode: provide bundle name and version via PTY input
	// Input format: "my-interactive-bundle\n1.0.0\n" (name + newline + version + newline)
	result := runOCInPTY(t, testEnv(t), "my-interactive-bundle\n1.0.0\n", "bundle", "init", "--output", bundleDir)
	requireSuccess(t, result)
	requireContains(t, result.stdout, "done: bundle initialized")

	// Verify manifest was created with provided values
	manifestPath := filepath.Join(bundleDir, "opencode-bundle.manifest.json")
	manifestData, err := os.ReadFile(manifestPath)
	if err != nil {
		t.Fatalf("failed to read manifest: %v", err)
	}
	requireContains(t, string(manifestData), `"bundle_name": "my-interactive-bundle"`)
	requireContains(t, string(manifestData), `"bundle_version": "1.0.0"`)
}
