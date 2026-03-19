package handler

import (
	"encoding/json"
	"net/http"

	"github.com/robotiqdev/project-13/internal/version"
)

// HealthResponse is the JSON response returned by HealthHandler.
type HealthResponse struct {
	Status string `json:"status"`
}

var healthResponseBytes []byte

func init() {
	var err error
	healthResponseBytes, err = json.Marshal(HealthResponse{Status: "ok"})
	if err != nil {
		panic(err)
	}
}

// HealthHandler handles health check requests.
func HealthHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", "GET")
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(healthResponseBytes)
}

// VersionHandler handles GET /version requests.
// It returns full build metadata (version, commit, buildDate) as JSON.
func VersionHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	info := version.Info()
	json.NewEncoder(w).Encode(map[string]string{
		"version":   info.Version,
		"commit":    info.Commit,
		"buildDate": info.BuildDate,
	})
}
