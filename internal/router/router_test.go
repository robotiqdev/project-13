package router_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/robotiqdev/project-13/internal/router"
	"github.com/robotiqdev/project-13/internal/version"
)

// TestVersionRoute_GET_Returns200 verifies that GET /version returns HTTP 200.
func TestVersionRoute_GET_Returns200(t *testing.T) {
	version.Version = "1.2.3"
	version.Commit = "abc123"
	version.BuildDate = "2026-03-19T00:00:00Z"

	handler := router.New()
	req := httptest.NewRequest(http.MethodGet, "/version", nil)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, rr.Code)
	}
}

// TestVersionRoute_GET_ContentTypeJSON verifies that GET /version returns Content-Type application/json.
func TestVersionRoute_GET_ContentTypeJSON(t *testing.T) {
	version.Version = "1.2.3"
	version.Commit = "abc123"
	version.BuildDate = "2026-03-19T00:00:00Z"

	handler := router.New()
	req := httptest.NewRequest(http.MethodGet, "/version", nil)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	ct := rr.Header().Get("Content-Type")
	if ct != "application/json" {
		t.Errorf("expected Content-Type %q, got %q", "application/json", ct)
	}
}

// TestVersionRoute_GET_ReturnsVersionJSON verifies that GET /version returns valid version JSON.
func TestVersionRoute_GET_ReturnsVersionJSON(t *testing.T) {
	version.Version = "2.0.0"
	version.Commit = "deadbeef"
	version.BuildDate = "2026-01-01T00:00:00Z"

	handler := router.New()
	req := httptest.NewRequest(http.MethodGet, "/version", nil)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	var info version.BuildInfo
	if err := json.NewDecoder(rr.Body).Decode(&info); err != nil {
		t.Fatalf("failed to decode response body into BuildInfo: %v", err)
	}

	if info.Version != "2.0.0" {
		t.Errorf("expected Version %q, got %q", "2.0.0", info.Version)
	}
	if info.Commit != "deadbeef" {
		t.Errorf("expected Commit %q, got %q", "deadbeef", info.Commit)
	}
	if info.BuildDate != "2026-01-01T00:00:00Z" {
		t.Errorf("expected BuildDate %q, got %q", "2026-01-01T00:00:00Z", info.BuildDate)
	}
}

// TestVersionRoute_GET_NoAuthHeaderRequired verifies that GET /version does not require any
// authorization headers — the route must be public.
func TestVersionRoute_GET_NoAuthHeaderRequired(t *testing.T) {
	version.Version = "1.0.0"
	version.Commit = "cafe"
	version.BuildDate = "2026-03-19T00:00:00Z"

	handler := router.New()

	// Explicitly send no Authorization header.
	req := httptest.NewRequest(http.MethodGet, "/version", nil)
	// Ensure no Authorization header is present.
	req.Header.Del("Authorization")
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("GET /version without auth header: expected status %d, got %d — route must be public", http.StatusOK, rr.Code)
	}
}

// TestVersionRoute_GET_NoAPIKeyRequired verifies that GET /version does not require an
// X-API-Key header or any other custom auth header.
func TestVersionRoute_GET_NoAPIKeyRequired(t *testing.T) {
	version.Version = "1.0.0"
	version.Commit = "cafe"
	version.BuildDate = "2026-03-19T00:00:00Z"

	handler := router.New()

	req := httptest.NewRequest(http.MethodGet, "/version", nil)
	req.Header.Del("X-API-Key")
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if rr.Code == http.StatusUnauthorized || rr.Code == http.StatusForbidden {
		t.Errorf("GET /version returned auth error %d — route must be public with no auth required", rr.Code)
	}
	if rr.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, rr.Code)
	}
}

// TestVersionRoute_POST_Returns405 verifies that POST /version returns HTTP 405.
func TestVersionRoute_POST_Returns405(t *testing.T) {
	handler := router.New()
	req := httptest.NewRequest(http.MethodPost, "/version", nil)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusMethodNotAllowed {
		t.Errorf("expected status %d for POST /version, got %d", http.StatusMethodNotAllowed, rr.Code)
	}
}

// TestVersionRoute_POST_Returns405WithJSONBody verifies that POST /version returns a JSON error body.
func TestVersionRoute_POST_Returns405WithJSONBody(t *testing.T) {
	handler := router.New()
	req := httptest.NewRequest(http.MethodPost, "/version", nil)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusMethodNotAllowed {
		t.Errorf("expected status %d for POST /version, got %d", http.StatusMethodNotAllowed, rr.Code)
	}

	ct := rr.Header().Get("Content-Type")
	if ct != "application/json" {
		t.Errorf("expected Content-Type %q for 405 response, got %q", "application/json", ct)
	}

	if !json.Valid(rr.Body.Bytes()) {
		t.Errorf("POST /version 405 response body is not valid JSON: %s", rr.Body.String())
	}
}

// TestVersionRoute_GET_ResponseBodyIsValidJSON verifies that the response body is valid JSON.
func TestVersionRoute_GET_ResponseBodyIsValidJSON(t *testing.T) {
	version.Version = "1.0.0"
	version.Commit = "abc"
	version.BuildDate = "2026-03-19T00:00:00Z"

	handler := router.New()
	req := httptest.NewRequest(http.MethodGet, "/version", nil)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if !json.Valid(rr.Body.Bytes()) {
		t.Errorf("GET /version response body is not valid JSON: %s", rr.Body.String())
	}
}

// TestVersionRoute_NonGETMethods_AllReturn405 is a table-driven test verifying non-GET methods return 405.
func TestVersionRoute_NonGETMethods_AllReturn405(t *testing.T) {
	methods := []string{
		http.MethodPost,
		http.MethodPut,
		http.MethodDelete,
		http.MethodPatch,
	}

	h := router.New()

	for _, method := range methods {
		t.Run(method, func(t *testing.T) {
			req := httptest.NewRequest(method, "/version", nil)
			rr := httptest.NewRecorder()

			h.ServeHTTP(rr, req)

			if rr.Code != http.StatusMethodNotAllowed {
				t.Errorf("method %s /version: expected status %d, got %d", method, http.StatusMethodNotAllowed, rr.Code)
			}
		})
	}
}
