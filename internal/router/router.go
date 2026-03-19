package router

import "net/http"

// NewRouter creates and returns a configured HTTP router.
// /health is intentionally public and requires no authentication for load balancer access.
func NewRouter() http.Handler {
	return http.NewServeMux()
}
