// Package cipipeline_test validates the CI/CD pipeline configuration required by TASK-4773.
// Per the task spec, no code unit tests are written; instead these tests assert that
// the expected workflow file is present with the necessary structure so that the
// manual validation steps (checking ldflags in build output, verifying /version endpoint
// after deploy) can succeed.
package cipipeline_test

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

// workflowPath returns the path to the CI build workflow file.
func workflowPath(t *testing.T) string {
	t.Helper()
	return filepath.Join(repoRoot(t), ".github", "workflows", "build.yml")
}

// readWorkflow reads the workflow file contents, skipping the test if the file
// does not yet exist (allowing other tests to report the missing-file error first).
func readWorkflow(t *testing.T) string {
	t.Helper()
	path := workflowPath(t)
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		t.Skipf(".github/workflows/build.yml not found — skipping; create the file as described in TASK-4773")
	}
	if err != nil {
		t.Fatalf("reading .github/workflows/build.yml: %v", err)
	}
	return string(data)
}

// TestWorkflowFileExists verifies that .github/workflows/build.yml is present.
func TestWorkflowFileExists(t *testing.T) {
	path := workflowPath(t)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		t.Errorf(".github/workflows/build.yml does not exist at %s; create the file as described in TASK-4773", path)
	}
}

// TestWorkflowHasGoTestStep verifies that the workflow runs go test ./... before building.
func TestWorkflowHasGoTestStep(t *testing.T) {
	content := readWorkflow(t)
	if !strings.Contains(content, "go test") {
		t.Error(".github/workflows/build.yml does not contain a 'go test' step; add 'go test ./...' before the build step")
	}
}

// TestWorkflowHasBuildStep verifies that the workflow contains a build step
// (either via 'go build' directly or by invoking 'make build').
func TestWorkflowHasBuildStep(t *testing.T) {
	content := readWorkflow(t)
	hasMakeBuild := strings.Contains(content, "make build")
	hasGoBuild := strings.Contains(content, "go build")
	if !hasMakeBuild && !hasGoBuild {
		t.Error(".github/workflows/build.yml does not contain a build step; add 'make build' or 'go build' with ldflags")
	}
}

// TestWorkflowBuildUsesLdflags verifies that the workflow's build step injects
// ldflags OR delegates to make build / build.sh which contains the ldflags.
// Acceptable patterns:
//   - the workflow file itself contains "-ldflags"
//   - the workflow invokes "make build" (which delegates to scripts/build.sh with ldflags)
//   - the workflow invokes "scripts/build.sh" directly
func TestWorkflowBuildUsesLdflags(t *testing.T) {
	content := readWorkflow(t)
	hasLdflags := strings.Contains(content, "-ldflags") || strings.Contains(content, "ldflags")
	hasMakeBuild := strings.Contains(content, "make build")
	hasBuildScript := strings.Contains(content, "build.sh")
	if !hasLdflags && !hasMakeBuild && !hasBuildScript {
		t.Error(".github/workflows/build.yml build step must inject ldflags (via -ldflags, make build, or scripts/build.sh)")
	}
}

// TestWorkflowInjectsVersionLdflagOrDelegates verifies that when the workflow
// sets ldflags inline it references the Version variable, OR that it uses
// 'make build' / 'scripts/build.sh' (which already contain the ldflag).
func TestWorkflowInjectsVersionLdflagOrDelegates(t *testing.T) {
	content := readWorkflow(t)
	// Delegation to make build or build.sh is acceptable because those already
	// carry the correct ldflags.
	if strings.Contains(content, "make build") || strings.Contains(content, "build.sh") {
		return
	}
	want := "github.com/app/service/internal/version.Version"
	if !strings.Contains(content, want) {
		t.Errorf(".github/workflows/build.yml does not reference ldflag path %q and does not delegate to make build or build.sh", want)
	}
}

// TestWorkflowInjectsCommitLdflagOrDelegates verifies that when the workflow
// sets ldflags inline it references the Commit variable, OR delegates.
func TestWorkflowInjectsCommitLdflagOrDelegates(t *testing.T) {
	content := readWorkflow(t)
	if strings.Contains(content, "make build") || strings.Contains(content, "build.sh") {
		return
	}
	want := "github.com/app/service/internal/version.Commit"
	if !strings.Contains(content, want) {
		t.Errorf(".github/workflows/build.yml does not reference ldflag path %q and does not delegate to make build or build.sh", want)
	}
}

// TestWorkflowInjectsBuildDateLdflagOrDelegates verifies that when the workflow
// sets ldflags inline it references the BuildDate variable, OR delegates.
func TestWorkflowInjectsBuildDateLdflagOrDelegates(t *testing.T) {
	content := readWorkflow(t)
	if strings.Contains(content, "make build") || strings.Contains(content, "build.sh") {
		return
	}
	want := "github.com/app/service/internal/version.BuildDate"
	if !strings.Contains(content, want) {
		t.Errorf(".github/workflows/build.yml does not reference ldflag path %q and does not delegate to make build or build.sh", want)
	}
}

// TestWorkflowHasOnPushOrPullRequestTrigger verifies that the workflow is triggered
// on push or pull_request events, which is the standard CI trigger.
func TestWorkflowHasOnPushOrPullRequestTrigger(t *testing.T) {
	content := readWorkflow(t)
	hasPush := strings.Contains(content, "push")
	hasPullRequest := strings.Contains(content, "pull_request")
	if !hasPush && !hasPullRequest {
		t.Error(".github/workflows/build.yml does not contain a 'push' or 'pull_request' trigger")
	}
}

// TestWorkflowHasCheckoutStep verifies that the workflow checks out the repository
// (required so git commands in the build script can access commit info).
func TestWorkflowHasCheckoutStep(t *testing.T) {
	content := readWorkflow(t)
	if !strings.Contains(content, "checkout") {
		t.Error(".github/workflows/build.yml does not contain a checkout step; the build script needs git history to inject commit metadata")
	}
}

// TestWorkflowGoTestRunsBeforeBuild verifies that the 'go test' step appears
// before the build step in the workflow file (line-order check).
func TestWorkflowGoTestRunsBeforeBuild(t *testing.T) {
	content := readWorkflow(t)

	testIdx := strings.Index(content, "go test")
	if testIdx == -1 {
		t.Skip("no 'go test' step found — skipping order check")
	}

	buildIdx := -1
	for _, marker := range []string{"make build", "go build", "build.sh"} {
		idx := strings.Index(content, marker)
		if idx != -1 && (buildIdx == -1 || idx < buildIdx) {
			buildIdx = idx
		}
	}
	if buildIdx == -1 {
		t.Skip("no build step found — skipping order check")
	}

	if testIdx > buildIdx {
		t.Error(".github/workflows/build.yml: the 'go test' step must appear before the build step")
	}
}

// TestWorkflowTargetsBinServer verifies that when build flags are set inline
// (not via delegation) the output binary is bin/server, OR that the workflow
// delegates to make build / build.sh (which already target bin/server).
func TestWorkflowTargetsBinServer(t *testing.T) {
	content := readWorkflow(t)
	if strings.Contains(content, "make build") || strings.Contains(content, "build.sh") {
		return
	}
	if !strings.Contains(content, "bin/server") {
		t.Error(".github/workflows/build.yml build step does not target bin/server and does not delegate to make build or build.sh")
	}
}
