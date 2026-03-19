package version

// BuildInfo holds the build metadata for the service.
type BuildInfo struct {
	Version   string `json:"version"`
	Commit    string `json:"commit"`
	BuildDate string `json:"buildDate"`
}

// Package-level vars targeted by ldflags injection.
var Version string
var Commit string
var BuildDate string

// Info returns a BuildInfo populated from the package-level vars.
// Stub: intentionally broken — does not use vars so TDD tests fail.
func Info() BuildInfo {
	return BuildInfo{}
}
