package router

import (
	"net/http"

	"github.com/robotiqdev/project-13/internal/handler"
)

// New returns an http.Handler with all application routes registered.
func New() http.Handler {
	mux := http.NewServeMux()

	// /version is public for deployment tracking and incident correlation
	mux.HandleFunc("/version", handler.VersionHandler)

	return mux
}
