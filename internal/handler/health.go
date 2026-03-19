package handler

import (
	"net/http"
)

// HealthResponse is the JSON response returned by HealthHandler.
type HealthResponse struct {
	Status string `json:"status"`
}

// HealthHandler handles health check requests.
func HealthHandler(w http.ResponseWriter, r *http.Request) {
}
