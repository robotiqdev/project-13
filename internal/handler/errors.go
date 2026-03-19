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

// MethodNotAllowedHandler is a stub — implementation pending TASK-4778.
func MethodNotAllowedHandler(w http.ResponseWriter, r *http.Request) {
	panic("not implemented")
}
