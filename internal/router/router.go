package router

import (
	"net/http"

	"github.com/robotiqdev/project-13/internal/handler"
)

// NewRouter creates and returns a configured HTTP router.
// /health is intentionally public and requires no authentication for load balancer access.
func NewRouter() http.Handler {
	mux := http.NewServeMux()
	// /health is intentionally public and requires no authentication for load balancer access
	mux.HandleFunc("/health", handler.HealthHandler)
	return mux
}
