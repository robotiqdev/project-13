package router

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

// errorResponse mirrors the JSON error envelope returned by all handlers.
type errorResponse struct {
	Error string `json:"error"`
}

// TestNotFoundRoute verifies that a GET request to an unregistered path returns 404 JSON.
func TestNotFoundRoute(t *testing.T) {
	handler := New()

	req := httptest.NewRequest(http.MethodGet, "/nonexistent", nil)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusNotFound {
		t.Errorf("expected status %d, got %d", http.StatusNotFound, rr.Code)
	}

	ct := rr.Header().Get("Content-Type")
	if ct != "application/json" {
		t.Errorf("expected Content-Type %q, got %q", "application/json", ct)
	}

	var resp errorResponse
	if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response body: %v", err)
	}

	if resp.Error != "not found" {
		t.Errorf("expected error %q, got %q", "not found", resp.Error)
	}
}

// TestNotFoundRoute_JsonBody verifies the exact JSON body for a 404 response.
func TestNotFoundRoute_JsonBody(t *testing.T) {
	handler := New()

	req := httptest.NewRequest(http.MethodGet, "/some/deep/path/that/does/not/exist", nil)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusNotFound {
		t.Errorf("expected status %d, got %d", http.StatusNotFound, rr.Code)
	}

	var raw map[string]interface{}
	if err := json.NewDecoder(rr.Body).Decode(&raw); err != nil {
		t.Fatalf("failed to decode response body as map: %v", err)
	}

	if len(raw) != 1 {
		t.Errorf("expected exactly 1 field in JSON body, got %d: %v", len(raw), raw)
	}

	errVal, ok := raw["error"]
	if !ok {
		t.Fatalf("expected JSON field \"error\" to exist, got keys: %v", raw)
	}

	if errVal != "not found" {
		t.Errorf("expected error value %q, got %q", "not found", errVal)
	}
}

// TestNotFoundRoute_ContentType verifies Content-Type is application/json for 404 responses.
func TestNotFoundRoute_ContentType(t *testing.T) {
	handler := New()

	paths := []string{"/unknown", "/foo/bar", "/api/v1/missing"}
	for _, path := range paths {
		t.Run(path, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, path, nil)
			rr := httptest.NewRecorder()

			handler.ServeHTTP(rr, req)

			ct := rr.Header().Get("Content-Type")
			if ct != "application/json" {
				t.Errorf("expected Content-Type %q, got %q", "application/json", ct)
			}
		})
	}
}

// TestMethodNotAllowed_Health verifies that POST /health returns 405 JSON.
func TestMethodNotAllowed_Health(t *testing.T) {
	handler := New()

	req := httptest.NewRequest(http.MethodPost, "/health", nil)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusMethodNotAllowed {
		t.Errorf("expected status %d, got %d", http.StatusMethodNotAllowed, rr.Code)
	}

	ct := rr.Header().Get("Content-Type")
	if ct != "application/json" {
		t.Errorf("expected Content-Type %q, got %q", "application/json", ct)
	}

	var resp errorResponse
	if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response body: %v", err)
	}

	if resp.Error != "method not allowed" {
		t.Errorf("expected error %q, got %q", "method not allowed", resp.Error)
	}
}

// TestMethodNotAllowed_Health_MultipleMethods verifies non-GET methods on /health return 405.
func TestMethodNotAllowed_Health_MultipleMethods(t *testing.T) {
	handler := New()

	methods := []string{http.MethodPost, http.MethodPut, http.MethodDelete, http.MethodPatch}
	for _, method := range methods {
		t.Run(method, func(t *testing.T) {
			req := httptest.NewRequest(method, "/health", nil)
			rr := httptest.NewRecorder()

			handler.ServeHTTP(rr, req)

			if rr.Code != http.StatusMethodNotAllowed {
				t.Errorf("expected status %d, got %d", http.StatusMethodNotAllowed, rr.Code)
			}

			ct := rr.Header().Get("Content-Type")
			if ct != "application/json" {
				t.Errorf("expected Content-Type %q, got %q", "application/json", ct)
			}

			var resp errorResponse
			if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
				t.Fatalf("failed to decode response body: %v", err)
			}

			if resp.Error != "method not allowed" {
				t.Errorf("expected error %q, got %q", "method not allowed", resp.Error)
			}
		})
	}
}

// TestMethodNotAllowed_ContentType verifies Content-Type is application/json for 405 responses.
func TestMethodNotAllowed_ContentType(t *testing.T) {
	handler := New()

	req := httptest.NewRequest(http.MethodPost, "/health", nil)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	ct := rr.Header().Get("Content-Type")
	if ct != "application/json" {
		t.Errorf("expected Content-Type %q, got %q", "application/json", ct)
	}
}

// TestHealthRoute_GET verifies that GET /health returns a successful response.
func TestHealthRoute_GET(t *testing.T) {
	handler := New()

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if rr.Code == http.StatusNotFound {
		t.Errorf("GET /health should not return 404; route should be registered")
	}
}

// TestRouterRegistersAllRoutes verifies that known routes don't fall through to the 404 handler.
func TestRouterRegistersAllRoutes(t *testing.T) {
	handler := New()

	knownRoutes := []string{"/health", "/version"}
	for _, path := range knownRoutes {
		t.Run(path, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, path, nil)
			rr := httptest.NewRecorder()

			handler.ServeHTTP(rr, req)

			if rr.Code == http.StatusNotFound {
				t.Errorf("route %s should be registered but returned 404", path)
			}
		})
	}
}

// TestNotFoundRoute_ValidJSON verifies that 404 error response body is valid JSON.
func TestNotFoundRoute_ValidJSON(t *testing.T) {
	handler := New()

	req := httptest.NewRequest(http.MethodGet, "/not-a-real-path", nil)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if !json.Valid(rr.Body.Bytes()) {
		t.Errorf("404 response body is not valid JSON: %s", rr.Body.String())
	}
}

// TestMethodNotAllowed_ValidJSON verifies that 405 error response body is valid JSON.
func TestMethodNotAllowed_ValidJSON(t *testing.T) {
	handler := New()

	req := httptest.NewRequest(http.MethodPost, "/health", nil)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if !json.Valid(rr.Body.Bytes()) {
		t.Errorf("405 response body is not valid JSON: %s", rr.Body.String())
	}
}

// TestHealthRoute_AccessibleWithoutAuthentication verifies that GET /health returns 200
// without any authentication headers present.
func TestHealthRoute_AccessibleWithoutAuthentication(t *testing.T) {
	r := New()

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
	r := New()

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
	r := New()

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rr := httptest.NewRecorder()

	r.ServeHTTP(rr, req)

	if rr.Code == http.StatusForbidden {
		t.Error("/health must be publicly accessible: got 403 Forbidden without auth header")
	}
}

// TestHealthRoute_ReturnsJSONBody verifies that /health returns a JSON response body.
func TestHealthRoute_ReturnsJSONBody(t *testing.T) {
	r := New()

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rr := httptest.NewRecorder()

	r.ServeHTTP(rr, req)

	if !json.Valid(rr.Body.Bytes()) {
		t.Errorf("/health response body is not valid JSON: %s", rr.Body.String())
	}
}

// TestHealthRoute_ReturnsStatusOKInBody verifies that /health returns {"status":"ok"} in the body.
func TestHealthRoute_ReturnsStatusOKInBody(t *testing.T) {
	r := New()

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
	r := New()

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
	r := New()

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
