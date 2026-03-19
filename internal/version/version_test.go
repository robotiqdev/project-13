package version_test

import (
	"testing"

	"github.com/robotiqdev/project-13/internal/version"
)

// TestDefaultVarsAreEmptyStrings verifies that the package-level vars
// have zero-value empty strings when no ldflags injection has occurred.
func TestDefaultVarsAreEmptyStrings(t *testing.T) {
	// Save and restore original values so this test does not pollute others.
	origVersion := version.Version
	origCommit := version.Commit
	origBuildDate := version.BuildDate
	defer func() {
		version.Version = origVersion
		version.Commit = origCommit
		version.BuildDate = origBuildDate
	}()

	// Reset to zero values to simulate the default (no ldflags) state.
	version.Version = ""
	version.Commit = ""
	version.BuildDate = ""

	if version.Version != "" {
		t.Errorf("Version default: got %q, want empty string", version.Version)
	}
	if version.Commit != "" {
		t.Errorf("Commit default: got %q, want empty string", version.Commit)
	}
	if version.BuildDate != "" {
		t.Errorf("BuildDate default: got %q, want empty string", version.BuildDate)
	}
}

// TestInfoReturnsStructWithVersionFromPackageVar verifies that Info().Version
// mirrors the package-level Version variable.
func TestInfoReturnsStructWithVersionFromPackageVar(t *testing.T) {
	orig := version.Version
	defer func() { version.Version = orig }()

	version.Version = "v1.2.3"
	got := version.Info()

	if got.Version != version.Version {
		t.Errorf("Info().Version = %q; want %q (version.Version)", got.Version, version.Version)
	}
}

// TestInfoReturnsStructWithCommitFromPackageVar verifies that Info().Commit
// mirrors the package-level Commit variable.
func TestInfoReturnsStructWithCommitFromPackageVar(t *testing.T) {
	orig := version.Commit
	defer func() { version.Commit = orig }()

	version.Commit = "abc1234def5678"
	got := version.Info()

	if got.Commit != version.Commit {
		t.Errorf("Info().Commit = %q; want %q (version.Commit)", got.Commit, version.Commit)
	}
}

// TestInfoReturnsStructWithBuildDateFromPackageVar verifies that Info().BuildDate
// mirrors the package-level BuildDate variable.
func TestInfoReturnsStructWithBuildDateFromPackageVar(t *testing.T) {
	orig := version.BuildDate
	defer func() { version.BuildDate = orig }()

	version.BuildDate = "2026-03-19T10:00:00Z"
	got := version.Info()

	if got.BuildDate != version.BuildDate {
		t.Errorf("Info().BuildDate = %q; want %q (version.BuildDate)", got.BuildDate, version.BuildDate)
	}
}

// TestInfoReturnsAllThreeFieldsPopulated verifies that a single call to Info()
// returns a BuildInfo with all three fields mirroring their respective vars.
func TestInfoReturnsAllThreeFieldsPopulated(t *testing.T) {
	origVersion := version.Version
	origCommit := version.Commit
	origBuildDate := version.BuildDate
	defer func() {
		version.Version = origVersion
		version.Commit = origCommit
		version.BuildDate = origBuildDate
	}()

	version.Version = "v2.0.0"
	version.Commit = "deadbeef"
	version.BuildDate = "2026-01-01T00:00:00Z"

	got := version.Info()

	if got.Version != version.Version {
		t.Errorf("Info().Version = %q; want %q", got.Version, version.Version)
	}
	if got.Commit != version.Commit {
		t.Errorf("Info().Commit = %q; want %q", got.Commit, version.Commit)
	}
	if got.BuildDate != version.BuildDate {
		t.Errorf("Info().BuildDate = %q; want %q", got.BuildDate, version.BuildDate)
	}
}

// TestInfoReturnsBuildInfoType verifies that Info() returns a value of type BuildInfo.
func TestInfoReturnsBuildInfoType(t *testing.T) {
	got := version.Info()
	// Compile-time check: assign to the concrete type.
	var _ version.BuildInfo = got
}

// TestInfoReflectsVarChanges verifies that Info() reads the current value of
// the vars at call time, not a snapshot taken at package init.
func TestInfoReflectsVarChanges(t *testing.T) {
	origVersion := version.Version
	defer func() { version.Version = origVersion }()

	version.Version = "v3.0.0-alpha"
	first := version.Info()

	version.Version = "v3.0.0"
	second := version.Info()

	if first.Version == second.Version {
		t.Errorf("Info() did not reflect updated Version: both calls returned %q", first.Version)
	}
	if second.Version != "v3.0.0" {
		t.Errorf("Info().Version after update = %q; want %q", second.Version, "v3.0.0")
	}
}
