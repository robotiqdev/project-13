package router

import (
	"net/http"

	"github.com/robotiqdev/project-13/internal/handler"
)

// New returns an http.Handler with all routes registered.
func New() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/health", handler.HealthHandler)
	mux.HandleFunc("/version", handler.VersionHandler)
	mux.HandleFunc("/", handler.NotFoundHandler)
	return mux
}
