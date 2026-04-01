package e2e

import (
	"archive/zip"
	"bytes"
	"context"
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestGoInstallTaggedVersionReportsTag(t *testing.T) {
	const modulePath = "github.com/sven1103-agent/opencode-config-cli"
	const moduleVersion = "v1.0.0-alpha.3"

	proxyRoot := t.TempDir()
	createModuleProxy(t, proxyRoot, modulePath, moduleVersion)

	binDir := filepath.Join(t.TempDir(), "bin")
	if err := os.MkdirAll(binDir, 0o755); err != nil {
		t.Fatalf("failed to create bin dir: %v", err)
	}
	goPath := filepath.Join(t.TempDir(), "gopath")
	if err := os.MkdirAll(goPath, 0o755); err != nil {
		t.Fatalf("failed to create GOPATH dir: %v", err)
	}

	install := runCommand(t, commandSpec{
		name: "go",
		args: []string{"install", modulePath + "@" + moduleVersion},
		env: []string{
			"GOBIN=" + binDir,
			"GOPATH=" + goPath,
			"GOMODCACHE=" + filepath.Join(goPath, "pkg", "mod"),
			"GOPROXY=file://" + proxyRoot + ",https://proxy.golang.org,direct",
			"GOSUMDB=off",
			"GONOSUMDB=*",
			"GOFLAGS=-modcacherw",
		},
	})
	requireSuccess(t, install)

	binaryPath := filepath.Join(binDir, "opencode-config-cli")
	metadata := runCommand(t, commandSpec{name: "go", args: []string{"version", "-m", binaryPath}})
	requireSuccess(t, metadata)
	requireContains(t, metadata.stdout, modulePath+"\t"+moduleVersion)

	version := runCommand(t, commandSpec{name: binaryPath, args: []string{"version"}})
	requireSuccess(t, version)
	requireContains(t, version.stdout, "oc "+moduleVersion)
}

type commandSpec struct {
	name string
	args []string
	env  []string
}

func runCommand(t *testing.T, spec commandSpec) commandResult {
	t.Helper()

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	cmd := exec.CommandContext(ctx, spec.name, spec.args...)
	cmd.Env = mergeEnv(os.Environ(), spec.env)

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()

	return commandResult{stdout: stdout.String(), stderr: stderr.String(), err: err}
}

func mergeEnv(base, overrides []string) []string {
	values := map[string]string{}
	order := []string{}

	for _, entry := range append(base, overrides...) {
		key, value, ok := strings.Cut(entry, "=")
		if !ok {
			continue
		}
		if _, seen := values[key]; !seen {
			order = append(order, key)
		}
		values[key] = value
	}

	merged := make([]string, 0, len(order))
	for _, key := range order {
		merged = append(merged, key+"="+values[key])
	}

	return merged
}

func createModuleProxy(t *testing.T, proxyRoot, modulePath, moduleVersion string) {
	t.Helper()

	workDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get working directory: %v", err)
	}
	repoRoot := filepath.Dir(workDir)

	moduleDir := filepath.Join(proxyRoot, filepath.FromSlash(modulePath), "@v")
	if err := os.MkdirAll(moduleDir, 0o755); err != nil {
		t.Fatalf("failed to create module proxy dir: %v", err)
	}

	goModPath := filepath.Join(repoRoot, "go.mod")
	goModData, err := os.ReadFile(goModPath)
	if err != nil {
		t.Fatalf("failed to read go.mod: %v", err)
	}
	if err := os.WriteFile(filepath.Join(moduleDir, moduleVersion+".mod"), goModData, 0o644); err != nil {
		t.Fatalf("failed to write module file: %v", err)
	}

	infoData, err := json.Marshal(struct {
		Version string
		Time    time.Time
	}{
		Version: moduleVersion,
		Time:    time.Date(2026, time.April, 1, 0, 0, 0, 0, time.UTC),
	})
	if err != nil {
		t.Fatalf("failed to marshal module info: %v", err)
	}
	if err := os.WriteFile(filepath.Join(moduleDir, moduleVersion+".info"), append(infoData, '\n'), 0o644); err != nil {
		t.Fatalf("failed to write module info: %v", err)
	}
	if err := os.WriteFile(filepath.Join(moduleDir, "list"), []byte(moduleVersion+"\n"), 0o644); err != nil {
		t.Fatalf("failed to write module version list: %v", err)
	}

	zipPath := filepath.Join(moduleDir, moduleVersion+".zip")
	createModuleZip(t, repoRoot, zipPath, modulePath+"@"+moduleVersion, moduleVersion)
}

func createModuleZip(t *testing.T, repoRoot, zipPath, zipRoot, moduleVersion string) {
	t.Helper()

	zipFile, err := os.Create(zipPath)
	if err != nil {
		t.Fatalf("failed to create module zip: %v", err)
	}
	defer zipFile.Close()

	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	if err := filepath.WalkDir(repoRoot, func(path string, entry os.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}

		relPath, err := filepath.Rel(repoRoot, path)
		if err != nil {
			return err
		}
		if relPath == "." {
			return nil
		}

		name := entry.Name()
		if entry.IsDir() {
			if name == ".git" || name == ".worktrees" || name == ".tmp" {
				return filepath.SkipDir
			}
			return nil
		}
		if strings.HasPrefix(relPath, ".git/") || strings.HasPrefix(relPath, ".worktrees/") || strings.HasPrefix(relPath, ".tmp/") {
			return nil
		}

		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		if filepath.ToSlash(relPath) == "internal/version/version.go" {
			data = []byte(strings.ReplaceAll(string(data), "var Version = \"0.0.0-dev\"", "var Version = \""+moduleVersion+"\""))
		}

		info, err := entry.Info()
		if err != nil {
			return err
		}

		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}
		header.Name = filepath.ToSlash(filepath.Join(zipRoot, relPath))
		header.Method = zip.Deflate

		writer, err := zipWriter.CreateHeader(header)
		if err != nil {
			return err
		}
		if _, err := writer.Write(data); err != nil {
			return err
		}
		return nil
	}); err != nil {
		t.Fatalf("failed to build module zip: %v", err)
	}

	if err := zipWriter.Close(); err != nil {
		t.Fatalf("failed to finalize module zip: %v", err)
	}
	if err := zipFile.Close(); err != nil {
		t.Fatalf("failed to close module zip: %v", err)
	}
}
