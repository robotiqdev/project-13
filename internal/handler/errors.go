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
