package handler

import (
	"encoding/json"
	"net/http"

	"github.com/robotiqdev/project-13/internal/version"
)

// VersionHandler serves build version information.
func VersionHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(version.Info())
}
