package router_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/robotiqdev/project-13/internal/router"
)

// TestHealthRoute_AccessibleWithoutAuthentication verifies that GET /health returns 200
// without any authentication headers present.
func TestHealthRoute_AccessibleWithoutAuthentication(t *testing.T) {
	r := router.NewRouter()

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	// No Authorization header, no token of any kind
	rr := httptest.NewRecorder()

	r.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected status %d on /health without auth headers, got %d", http.StatusOK, rr.Code)
	}
}

// TestHealthRoute_NoAuthorizationHeaderRequired verifies that /health does not return
// 401 Unauthorized when no Authorization header is provided.
func TestHealthRoute_NoAuthorizationHeaderRequired(t *testing.T) {
	r := router.NewRouter()

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rr := httptest.NewRecorder()

	r.ServeHTTP(rr, req)

	if rr.Code == http.StatusUnauthorized {
		t.Error("/health must be publicly accessible: got 401 Unauthorized without auth header")
	}
}

// TestHealthRoute_NotForbiddenWithoutAuth verifies that /health does not return
// 403 Forbidden when no authentication credentials are provided.
func TestHealthRoute_NotForbiddenWithoutAuth(t *testing.T) {
	r := router.NewRouter()

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rr := httptest.NewRecorder()

	r.ServeHTTP(rr, req)

	if rr.Code == http.StatusForbidden {
		t.Error("/health must be publicly accessible: got 403 Forbidden without auth header")
	}
}

// TestHealthRoute_ReturnsJSONBody verifies that /health returns a JSON response body.
func TestHealthRoute_ReturnsJSONBody(t *testing.T) {
	r := router.NewRouter()

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rr := httptest.NewRecorder()

	r.ServeHTTP(rr, req)

	if !json.Valid(rr.Body.Bytes()) {
		t.Errorf("/health response body is not valid JSON: %s", rr.Body.String())
	}
}

// TestHealthRoute_ReturnsStatusOKInBody verifies that /health returns {"status":"ok"} in the body.
func TestHealthRoute_ReturnsStatusOKInBody(t *testing.T) {
	r := router.NewRouter()

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rr := httptest.NewRecorder()

	r.ServeHTTP(rr, req)

	var body map[string]interface{}
	if err := json.Unmarshal(rr.Body.Bytes(), &body); err != nil {
		t.Fatalf("failed to unmarshal /health response body: %v", err)
	}

	status, ok := body["status"]
	if !ok {
		t.Fatal("/health response body missing \"status\" field")
	}
	if status != "ok" {
		t.Errorf("expected status field %q, got %q", "ok", status)
	}
}

// TestHealthRoute_ReturnsJSONContentType verifies that /health sets Content-Type to application/json.
func TestHealthRoute_ReturnsJSONContentType(t *testing.T) {
	r := router.NewRouter()

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rr := httptest.NewRecorder()

	r.ServeHTTP(rr, req)

	ct := rr.Header().Get("Content-Type")
	if ct != "application/json" {
		t.Errorf("expected Content-Type %q, got %q", "application/json", ct)
	}
}

// TestHealthRoute_RegistrationOutsideAuthMiddleware verifies that /health is accessible
// even when a request lacks auth credentials that a protected route would require.
// This confirms /health is on the base mux, not behind an auth middleware chain.
func TestHealthRoute_RegistrationOutsideAuthMiddleware(t *testing.T) {
	r := router.NewRouter()

	// Simulate a request with no credentials at all (no Bearer, no API key, nothing)
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	req.Header.Del("Authorization")
	req.Header.Del("X-Api-Key")
	rr := httptest.NewRecorder()

	r.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf(
			"/health must be registered on base mux (public), not behind auth middleware: expected 200, got %d",
			rr.Code,
		)
	}
}
