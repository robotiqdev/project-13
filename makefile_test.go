package integration_test

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

// repoRoot returns the absolute path to the repository root (where this file lives).
func repoRoot(t *testing.T) string {
	t.Helper()
	abs, err := filepath.Abs(".")
	if err != nil {
		t.Fatalf("failed to determine repo root: %v", err)
	}
	return abs
}

// readMakefile reads the Makefile from the repo root and returns its contents.
// Fails the test if the file does not exist.
func readMakefile(t *testing.T) string {
	t.Helper()
	root := repoRoot(t)
	path := filepath.Join(root, "Makefile")
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("Makefile not found at %s: %v", path, err)
	}
	return string(data)
}

// hasMakeTarget scans the Makefile text for a target line matching `<name>:`.
func hasMakeTarget(content, name string) bool {
	scanner := bufio.NewScanner(strings.NewReader(content))
	prefix := name + ":"
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, prefix) {
			return true
		}
	}
	return false
}

// buildTargetLines returns all lines belonging to the build target block
// (the target line itself plus any indented recipe lines that follow).
func buildTargetLines(content string) []string {
	var lines []string
	inBlock := false
	scanner := bufio.NewScanner(strings.NewReader(content))
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "build:") {
			inBlock = true
			lines = append(lines, line)
			continue
		}
		if inBlock {
			// Recipe lines are indented with a tab.
			if strings.HasPrefix(line, "\t") || strings.HasPrefix(line, " ") {
				lines = append(lines, line)
			} else {
				break
			}
		}
	}
	return lines
}

// buildTargetText returns the full text of the build target recipe.
func buildTargetText(content string) string {
	return strings.Join(buildTargetLines(content), "\n")
}

// TestMakefile_Exists verifies that a Makefile exists at the repository root.
func TestMakefile_Exists(t *testing.T) {
	root := repoRoot(t)
	path := filepath.Join(root, "Makefile")
	if _, err := os.Stat(path); os.IsNotExist(err) {
		t.Fatalf("Makefile does not exist at %s", path)
	}
}

// TestMakefile_HasBuildTarget verifies that the Makefile defines a build target.
func TestMakefile_HasBuildTarget(t *testing.T) {
	content := readMakefile(t)
	if !hasMakeTarget(content, "build") {
		t.Error("Makefile does not have a 'build:' target")
	}
}

// TestMakefile_BuildTarget_ContainsLDFlags verifies that the build target passes -ldflags to go build.
func TestMakefile_BuildTarget_ContainsLDFlags(t *testing.T) {
	content := readMakefile(t)
	recipe := buildTargetText(content)
	if !strings.Contains(recipe, "-ldflags") {
		t.Errorf("build target recipe does not contain -ldflags:\n%s", recipe)
	}
}

// TestMakefile_BuildTarget_InjectsVersion verifies the build target injects version.Version via ldflags.
func TestMakefile_BuildTarget_InjectsVersion(t *testing.T) {
	content := readMakefile(t)
	recipe := buildTargetText(content)
	if !strings.Contains(recipe, "version.Version") {
		t.Errorf("build target ldflags do not inject version.Version:\n%s", recipe)
	}
}

// TestMakefile_BuildTarget_InjectsCommit verifies the build target injects version.Commit via ldflags.
func TestMakefile_BuildTarget_InjectsCommit(t *testing.T) {
	content := readMakefile(t)
	recipe := buildTargetText(content)
	if !strings.Contains(recipe, "version.Commit") {
		t.Errorf("build target ldflags do not inject version.Commit:\n%s", recipe)
	}
}

// TestMakefile_BuildTarget_InjectsBuildDate verifies the build target injects version.BuildDate via ldflags.
func TestMakefile_BuildTarget_InjectsBuildDate(t *testing.T) {
	content := readMakefile(t)
	recipe := buildTargetText(content)
	if !strings.Contains(recipe, "version.BuildDate") {
		t.Errorf("build target ldflags do not inject version.BuildDate:\n%s", recipe)
	}
}

// TestMakefile_HasRunTarget verifies that the Makefile defines a run target for local development.
func TestMakefile_HasRunTarget(t *testing.T) {
	content := readMakefile(t)
	if !hasMakeTarget(content, "run") {
		t.Error("Makefile does not have a 'run:' target")
	}
}

// TestMakefile_HasTestTarget verifies that the Makefile defines a test target.
func TestMakefile_HasTestTarget(t *testing.T) {
	content := readMakefile(t)
	if !hasMakeTarget(content, "test") {
		t.Error("Makefile does not have a 'test:' target")
	}
}

// TestMakefile_RunTarget_ContainsLDFlags verifies the run target also passes -ldflags.
func TestMakefile_RunTarget_ContainsLDFlags(t *testing.T) {
	content := readMakefile(t)
	// Collect lines for the run target.
	var lines []string
	inBlock := false
	scanner := bufio.NewScanner(strings.NewReader(content))
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "run:") {
			inBlock = true
			lines = append(lines, line)
			continue
		}
		if inBlock {
			if strings.HasPrefix(line, "\t") || strings.HasPrefix(line, " ") {
				lines = append(lines, line)
			} else {
				break
			}
		}
	}
	recipe := strings.Join(lines, "\n")
	if !strings.Contains(recipe, "-ldflags") {
		t.Errorf("run target recipe does not contain -ldflags:\n%s", recipe)
	}
}

// buildInfo holds the JSON response from GET /version.
type buildInfo struct {
	Version   string `json:"version"`
	Commit    string `json:"commit"`
	BuildDate string `json:"build_date"`
}

// findFreePort returns a free TCP port on localhost.
func findFreePort(t *testing.T) int {
	t.Helper()
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("failed to find a free port: %v", err)
	}
	port := ln.Addr().(*net.TCPAddr).Port
	ln.Close()
	return port
}

// buildBinary compiles the cmd/server package with the given extra ldflags
// and writes the binary to outPath.  Returns the combined output on error.
func buildBinary(t *testing.T, outPath string, extraLDFlags string) {
	t.Helper()
	root := repoRoot(t)
	args := []string{"build"}
	if extraLDFlags != "" {
		args = append(args, "-ldflags", extraLDFlags)
	}
	args = append(args, "-o", outPath, "./cmd/server/")
	cmd := exec.Command("go", args...)
	cmd.Dir = root
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("go build failed: %v\n%s", err, out)
	}
}

// startServer starts the binary at binPath, waits until it is listening on
// addr (host:port), and registers a cleanup to stop it.
func startServer(t *testing.T, binPath string, addr string) {
	t.Helper()
	cmd := exec.Command(binPath)
	cmd.Env = append(os.Environ(), fmt.Sprintf("PORT=%s", strings.Split(addr, ":")[1]))
	cmd.Stdout = io.Discard
	cmd.Stderr = io.Discard
	if err := cmd.Start(); err != nil {
		t.Fatalf("failed to start server binary: %v", err)
	}
	t.Cleanup(func() { _ = cmd.Process.Kill() })

	// Wait for the server to accept connections.
	deadline := time.Now().Add(5 * time.Second)
	for time.Now().Before(deadline) {
		conn, err := net.DialTimeout("tcp", addr, 100*time.Millisecond)
		if err == nil {
			conn.Close()
			return
		}
		time.Sleep(50 * time.Millisecond)
	}
	t.Fatalf("server at %s did not become ready within 5 seconds", addr)
}

// TestMakeBuild_VersionEndpoint_AllFieldsNonEmpty runs 'make build', starts the
// resulting binary, sends GET /version, and asserts every field is non-empty —
// meaning all three ldflags were injected at build time.
func TestMakeBuild_VersionEndpoint_AllFieldsNonEmpty(t *testing.T) {
	root := repoRoot(t)

	// Run make build.
	makeCmd := exec.Command("make", "build")
	makeCmd.Dir = root
	out, err := makeCmd.CombinedOutput()
	if err != nil {
		t.Fatalf("'make build' failed: %v\n%s", err, out)
	}

	// Locate the produced binary.  Convention: bin/server or ./server.
	var binPath string
	candidates := []string{
		filepath.Join(root, "bin", "server"),
		filepath.Join(root, "server"),
	}
	for _, c := range candidates {
		if _, err := os.Stat(c); err == nil {
			binPath = c
			break
		}
	}
	if binPath == "" {
		t.Fatalf("'make build' succeeded but binary not found at any of: %v", candidates)
	}
	t.Cleanup(func() { os.Remove(binPath) })

	port := findFreePort(t)
	addr := fmt.Sprintf("127.0.0.1:%d", port)
	startServer(t, binPath, addr)

	resp, err := http.Get(fmt.Sprintf("http://%s/version", addr))
	if err != nil {
		t.Fatalf("GET /version failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("GET /version: expected 200, got %d", resp.StatusCode)
	}

	var info buildInfo
	if err := json.NewDecoder(resp.Body).Decode(&info); err != nil {
		t.Fatalf("failed to decode /version response: %v", err)
	}

	if info.Version == "" {
		t.Error("version field is empty — ldflags must inject a non-empty version string")
	}
	if info.Commit == "" {
		t.Error("commit field is empty — ldflags must inject a non-empty commit string")
	}
	if info.BuildDate == "" {
		t.Error("build_date field is empty — ldflags must inject a non-empty build date")
	}
}

// TestGoBuildWithoutLDFlags_DefaultVersionValues builds cmd/server without any
// ldflags and verifies that the /version endpoint returns the compiled-in
// defaults: version="dev", commit="unknown", build_date="unknown".
func TestGoBuildWithoutLDFlags_DefaultVersionValues(t *testing.T) {
	root := repoRoot(t)
	binPath := filepath.Join(root, "bin", "server-noflags-test")
	_ = os.MkdirAll(filepath.Dir(binPath), 0o755)

	buildBinary(t, binPath, "")
	t.Cleanup(func() { os.Remove(binPath) })

	port := findFreePort(t)
	addr := fmt.Sprintf("127.0.0.1:%d", port)
	startServer(t, binPath, addr)

	resp, err := http.Get(fmt.Sprintf("http://%s/version", addr))
	if err != nil {
		t.Fatalf("GET /version failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("GET /version: expected 200, got %d", resp.StatusCode)
	}

	var info buildInfo
	if err := json.NewDecoder(resp.Body).Decode(&info); err != nil {
		t.Fatalf("failed to decode /version response: %v", err)
	}

	if info.Version != "dev" {
		t.Errorf("default version: expected %q, got %q", "dev", info.Version)
	}
	if info.Commit != "unknown" {
		t.Errorf("default commit: expected %q, got %q", "unknown", info.Commit)
	}
	if info.BuildDate != "unknown" {
		t.Errorf("default build_date: expected %q, got %q", "unknown", info.BuildDate)
	}
}
