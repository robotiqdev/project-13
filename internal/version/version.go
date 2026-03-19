package version

// BuildInfo holds the build metadata for the service.
type BuildInfo struct {
	Version   string `json:"version"`
	Commit    string `json:"commit"`
	BuildDate string `json:"buildDate"`
}

// These variables are set at build time via -ldflags.
var (
	Version   = "dev"
	Commit    = "unknown"
	BuildDate = "unknown"
)

// Info returns the current build information.
func Info() BuildInfo {
	return BuildInfo{
		Version:   Version,
		Commit:    Commit,
		BuildDate: BuildDate,
	}
}
