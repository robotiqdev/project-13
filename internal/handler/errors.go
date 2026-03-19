package handler

import (
	"encoding/json"
	"net/http"
)

// ErrorResponse is the JSON error envelope returned by all handlers.
type ErrorResponse struct {
	Error string `json:"error"`
}

func writeError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(ErrorResponse{Error: message})
}

// MethodNotAllowedHandler returns a 405 Method Not Allowed JSON response.
func MethodNotAllowedHandler(w http.ResponseWriter, r *http.Request) {
	writeError(w, http.StatusMethodNotAllowed, "method not allowed")
}

// NotFoundHandler returns a 404 Not Found JSON response.
func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	writeError(w, http.StatusNotFound, "not found")
}

// HealthHandler handles GET /health requests.
func HealthHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"ok"}`))
}

// VersionHandler handles GET /version requests.
func VersionHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"version":"1.0.0"}`))
}
