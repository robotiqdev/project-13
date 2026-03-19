package handler

import "net/http"

// ErrorResponse is the JSON error envelope returned by all handlers.
type ErrorResponse struct {
	Error string `json:"error"`
}

// writeError is a stub — implementation will be provided separately.
func writeError(w http.ResponseWriter, status int, message string) {
	panic("not implemented")
}
