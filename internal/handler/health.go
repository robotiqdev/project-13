package handler

import (
	"encoding/json"
	"net/http"
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
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(healthResponseBytes)
}
