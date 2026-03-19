package router

import (
	"net/http"

	"github.com/robotiqdev/project-13/internal/handler"
)

// New returns an http.Handler with all application routes registered.
func New() http.Handler {
	mux := http.NewServeMux()
	// /health is intentionally public and requires no authentication for load balancer access
	mux.HandleFunc("/health", handler.HealthHandler)
	// /version is public for deployment tracking and incident correlation
	mux.HandleFunc("/version", handler.VersionHandler)
	mux.HandleFunc("/", handler.NotFoundHandler)
	return mux
}

// NewRouter creates and returns a configured HTTP router.
// /health is intentionally public and requires no authentication for load balancer access.
func NewRouter() http.Handler {
	return New()
}
