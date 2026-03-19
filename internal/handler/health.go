package handler

import (
	"encoding/json"
	"net/http"

	"github.com/app/service/internal/version"
)

// HealthHandler handles GET /health requests.
// It returns a minimal {"status":"ok"} response suitable for load-balancer health checks.
func HealthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
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
