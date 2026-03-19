package router

import (
	"net/http"
)

// New returns an http.Handler with all application routes registered.
func New() http.Handler {
	mux := http.NewServeMux()
	return mux
}
