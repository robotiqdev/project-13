package handler

import "net/http"

// HealthHandler handles GET /health requests.
// It returns a minimal {"status":"ok"} response suitable for load-balancer health checks.
func HealthHandler(w http.ResponseWriter, r *http.Request) {
	// stub — implementation pending
}

// VersionHandler handles GET /version requests.
// It returns full build metadata (version, commit, buildDate) as JSON.
func VersionHandler(w http.ResponseWriter, r *http.Request) {
	// stub — implementation pending
}
