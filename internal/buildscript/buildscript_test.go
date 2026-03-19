// Package buildscript_test validates the build infrastructure required by TASK-4772.
// Per the task spec, no unit tests are written for the shell script logic itself;
// instead, these tests assert that the expected files are present with the
// necessary structure so that the manual validation steps can succeed.
package buildscript_test

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

// repoRoot returns the absolute path to the repository root by walking up from
// this test file's location until go.mod is found.
func repoRoot(t *testing.T) string {
	t.Helper()
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("runtime.Caller failed")
	}
	dir := filepath.Dir(file)
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			t.Fatal("could not find repository root (go.mod not found)")
		}
		dir = parent
	}
}

// TestBuildScriptExists verifies that scripts/build.sh is present in the repository.
func TestBuildScriptExists(t *testing.T) {
	root := repoRoot(t)
	path := filepath.Join(root, "scripts", "build.sh")
	if _, err := os.Stat(path); os.IsNotExist(err) {
		t.Errorf("scripts/build.sh does not exist at %s; create the file as described in TASK-4772", path)
	}
}

// TestBuildScriptIsExecutable verifies that scripts/build.sh has the executable bit set.
func TestBuildScriptIsExecutable(t *testing.T) {
	root := repoRoot(t)
	path := filepath.Join(root, "scripts", "build.sh")
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		t.Skipf("scripts/build.sh not found — skipping executable check")
	}
	if err != nil {
		t.Fatalf("os.Stat(%s): %v", path, err)
	}
	if info.Mode()&0111 == 0 {
		t.Errorf("scripts/build.sh is not executable (mode %s); run chmod +x scripts/build.sh", info.Mode())
	}
}

// TestBuildScriptHasShebang verifies that scripts/build.sh starts with #!/bin/bash.
func TestBuildScriptHasShebang(t *testing.T) {
	root := repoRoot(t)
	path := filepath.Join(root, "scripts", "build.sh")
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		t.Skipf("scripts/build.sh not found — skipping shebang check")
	}
	if err != nil {
		t.Fatalf("reading scripts/build.sh: %v", err)
	}
	content := string(data)
	if !strings.HasPrefix(content, "#!/bin/bash") {
		t.Errorf("scripts/build.sh does not start with #!/bin/bash; got first line: %q",
			strings.SplitN(content, "\n", 2)[0])
	}
}

// TestBuildScriptHasSetEUO verifies that scripts/build.sh contains "set -euo pipefail"
// so the script exits non-zero on build failure.
func TestBuildScriptHasSetEUO(t *testing.T) {
	root := repoRoot(t)
	path := filepath.Join(root, "scripts", "build.sh")
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		t.Skipf("scripts/build.sh not found — skipping pipefail check")
	}
	if err != nil {
		t.Fatalf("reading scripts/build.sh: %v", err)
	}
	if !strings.Contains(string(data), "set -euo pipefail") {
		t.Error("scripts/build.sh does not contain 'set -euo pipefail'; the script must exit non-zero on failure")
	}
}

// TestBuildScriptInjectsVersionLdflag verifies that scripts/build.sh passes the
// Version ldflag targeting the version package.
func TestBuildScriptInjectsVersionLdflag(t *testing.T) {
	root := repoRoot(t)
	path := filepath.Join(root, "scripts", "build.sh")
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		t.Skipf("scripts/build.sh not found — skipping ldflag check")
	}
	if err != nil {
		t.Fatalf("reading scripts/build.sh: %v", err)
	}
	content := string(data)
	want := "github.com/robotiqdev/project-13/internal/version.Version"
	if !strings.Contains(content, want) {
		t.Errorf("scripts/build.sh does not reference ldflag path %q", want)
	}
}

// TestBuildScriptInjectsCommitLdflag verifies that scripts/build.sh passes the
// Commit ldflag targeting the version package.
func TestBuildScriptInjectsCommitLdflag(t *testing.T) {
	root := repoRoot(t)
	path := filepath.Join(root, "scripts", "build.sh")
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		t.Skipf("scripts/build.sh not found — skipping ldflag check")
	}
	if err != nil {
		t.Fatalf("reading scripts/build.sh: %v", err)
	}
	content := string(data)
	want := "github.com/robotiqdev/project-13/internal/version.Commit"
	if !strings.Contains(content, want) {
		t.Errorf("scripts/build.sh does not reference ldflag path %q", want)
	}
}

// TestBuildScriptInjectsBuildDateLdflag verifies that scripts/build.sh passes the
// BuildDate ldflag targeting the version package.
func TestBuildScriptInjectsBuildDateLdflag(t *testing.T) {
	root := repoRoot(t)
	path := filepath.Join(root, "scripts", "build.sh")
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		t.Skipf("scripts/build.sh not found — skipping ldflag check")
	}
	if err != nil {
		t.Fatalf("reading scripts/build.sh: %v", err)
	}
	content := string(data)
	want := "github.com/robotiqdev/project-13/internal/version.BuildDate"
	if !strings.Contains(content, want) {
		t.Errorf("scripts/build.sh does not reference ldflag path %q", want)
	}
}

// TestBuildScriptCapturesGitDescribeVersion verifies that scripts/build.sh uses
// git describe to capture the version string.
func TestBuildScriptCapturesGitDescribeVersion(t *testing.T) {
	root := repoRoot(t)
	path := filepath.Join(root, "scripts", "build.sh")
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		t.Skipf("scripts/build.sh not found — skipping git describe check")
	}
	if err != nil {
		t.Fatalf("reading scripts/build.sh: %v", err)
	}
	if !strings.Contains(string(data), "git describe") {
		t.Error("scripts/build.sh does not use 'git describe' to capture version")
	}
}

// TestBuildScriptCapturesGitRevParseCommit verifies that scripts/build.sh uses
// git rev-parse to capture the short commit hash.
func TestBuildScriptCapturesGitRevParseCommit(t *testing.T) {
	root := repoRoot(t)
	path := filepath.Join(root, "scripts", "build.sh")
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		t.Skipf("scripts/build.sh not found — skipping git rev-parse check")
	}
	if err != nil {
		t.Fatalf("reading scripts/build.sh: %v", err)
	}
	if !strings.Contains(string(data), "git rev-parse") {
		t.Error("scripts/build.sh does not use 'git rev-parse' to capture commit hash")
	}
}

// TestBuildScriptHasFallbackForVersion verifies that scripts/build.sh provides a
// fallback value (e.g. 'dev') when git describe is unavailable (shallow CI clone).
func TestBuildScriptHasFallbackForVersion(t *testing.T) {
	root := repoRoot(t)
	path := filepath.Join(root, "scripts", "build.sh")
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		t.Skipf("scripts/build.sh not found — skipping fallback check")
	}
	if err != nil {
		t.Fatalf("reading scripts/build.sh: %v", err)
	}
	content := string(data)
	// The fallback can be implemented as "|| echo 'dev'" or similar.
	if !strings.Contains(content, "dev") {
		t.Error("scripts/build.sh does not contain a 'dev' fallback for version when git describe fails")
	}
}

// TestBuildScriptTargetsBinServer verifies that scripts/build.sh outputs to bin/server.
func TestBuildScriptTargetsBinServer(t *testing.T) {
	root := repoRoot(t)
	path := filepath.Join(root, "scripts", "build.sh")
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		t.Skipf("scripts/build.sh not found — skipping output path check")
	}
	if err != nil {
		t.Fatalf("reading scripts/build.sh: %v", err)
	}
	if !strings.Contains(string(data), "bin/server") {
		t.Error("scripts/build.sh does not output to bin/server")
	}
}

// TestMakefileExists verifies that a Makefile is present in the repository root.
func TestMakefileExists(t *testing.T) {
	root := repoRoot(t)
	path := filepath.Join(root, "Makefile")
	if _, err := os.Stat(path); os.IsNotExist(err) {
		t.Errorf("Makefile does not exist at %s; create the file as described in TASK-4772", path)
	}
}

// TestMakefileHasBuildTarget verifies that the Makefile defines a 'build' target.
func TestMakefileHasBuildTarget(t *testing.T) {
	root := repoRoot(t)
	path := filepath.Join(root, "Makefile")
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		t.Skipf("Makefile not found — skipping build target check")
	}
	if err != nil {
		t.Fatalf("reading Makefile: %v", err)
	}
	content := string(data)
	// A Makefile target starts at column 0 followed by a colon.
	if !strings.Contains(content, "build:") {
		t.Error("Makefile does not define a 'build:' target")
	}
}

// TestMakefileBuildTargetCallsBuildScript verifies that the Makefile 'build' target
// delegates to scripts/build.sh.
func TestMakefileBuildTargetCallsBuildScript(t *testing.T) {
	root := repoRoot(t)
	path := filepath.Join(root, "Makefile")
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		t.Skipf("Makefile not found — skipping build script delegation check")
	}
	if err != nil {
		t.Fatalf("reading Makefile: %v", err)
	}
	if !strings.Contains(string(data), "scripts/build.sh") {
		t.Error("Makefile 'build' target does not call scripts/build.sh")
	}
}
