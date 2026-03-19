package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/robotiqdev/project-13/internal/version"
)

// TestVersionHandler_GET_Returns200 verifies that a GET request returns HTTP 200.
func TestVersionHandler_GET_Returns200(t *testing.T) {
	version.Version = "1.2.3"
	version.Commit = "abc123def456"
	version.BuildDate = "2026-03-19T00:00:00Z"

	req := httptest.NewRequest(http.MethodGet, "/version", nil)
	rr := httptest.NewRecorder()

	VersionHandler(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, rr.Code)
	}
}

// TestVersionHandler_GET_ContentTypeJSON verifies that a GET request sets Content-Type to application/json.
func TestVersionHandler_GET_ContentTypeJSON(t *testing.T) {
	version.Version = "1.2.3"
	version.Commit = "abc123def456"
	version.BuildDate = "2026-03-19T00:00:00Z"

	req := httptest.NewRequest(http.MethodGet, "/version", nil)
	rr := httptest.NewRecorder()

	VersionHandler(rr, req)

	ct := rr.Header().Get("Content-Type")
	if ct != "application/json" {
		t.Errorf("expected Content-Type %q, got %q", "application/json", ct)
	}
}

// TestVersionHandler_GET_BodyDecodesIntoBuildInfo verifies the response body decodes into a BuildInfo struct.
func TestVersionHandler_GET_BodyDecodesIntoBuildInfo(t *testing.T) {
	version.Version = "1.2.3"
	version.Commit = "abc123def456"
	version.BuildDate = "2026-03-19T00:00:00Z"

	req := httptest.NewRequest(http.MethodGet, "/version", nil)
	rr := httptest.NewRecorder()

	VersionHandler(rr, req)

	var info version.BuildInfo
	if err := json.NewDecoder(rr.Body).Decode(&info); err != nil {
		t.Fatalf("failed to decode response body into BuildInfo: %v", err)
	}

	if info.Version != "1.2.3" {
		t.Errorf("expected Version %q, got %q", "1.2.3", info.Version)
	}
	if info.Commit != "abc123def456" {
		t.Errorf("expected Commit %q, got %q", "abc123def456", info.Commit)
	}
	if info.BuildDate != "2026-03-19T00:00:00Z" {
		t.Errorf("expected BuildDate %q, got %q", "2026-03-19T00:00:00Z", info.BuildDate)
	}
}

// TestVersionHandler_GET_BodyReflectsKnownValues verifies that the response body contains the exact values set on version variables.
func TestVersionHandler_GET_BodyReflectsKnownValues(t *testing.T) {
	version.Version = "0.0.1-test"
	version.Commit = "deadbeef"
	version.BuildDate = "1970-01-01T00:00:00Z"

	req := httptest.NewRequest(http.MethodGet, "/version", nil)
	rr := httptest.NewRecorder()

	VersionHandler(rr, req)

	var info version.BuildInfo
	if err := json.NewDecoder(rr.Body).Decode(&info); err != nil {
		t.Fatalf("failed to decode response body: %v", err)
	}

	if info.Version != "0.0.1-test" {
		t.Errorf("expected Version %q, got %q", "0.0.1-test", info.Version)
	}
	if info.Commit != "deadbeef" {
		t.Errorf("expected Commit %q, got %q", "deadbeef", info.Commit)
	}
	if info.BuildDate != "1970-01-01T00:00:00Z" {
		t.Errorf("expected BuildDate %q, got %q", "1970-01-01T00:00:00Z", info.BuildDate)
	}
}

// TestVersionHandler_GET_ValidJSON verifies that the response body is valid JSON.
func TestVersionHandler_GET_ValidJSON(t *testing.T) {
	version.Version = "1.0.0"
	version.Commit = "cafebabe"
	version.BuildDate = "2026-01-01T00:00:00Z"

	req := httptest.NewRequest(http.MethodGet, "/version", nil)
	rr := httptest.NewRecorder()

	VersionHandler(rr, req)

	if !json.Valid(rr.Body.Bytes()) {
		t.Errorf("response body is not valid JSON: %s", rr.Body.String())
	}
}

// TestVersionHandler_POST_Returns405 verifies that a POST request returns HTTP 405.
func TestVersionHandler_POST_Returns405(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/version", nil)
	rr := httptest.NewRecorder()

	VersionHandler(rr, req)

	if rr.Code != http.StatusMethodNotAllowed {
		t.Errorf("expected status %d, got %d", http.StatusMethodNotAllowed, rr.Code)
	}
}

// TestVersionHandler_POST_Returns405WithJSONError verifies that a POST request returns a JSON error body.
func TestVersionHandler_POST_Returns405WithJSONError(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/version", nil)
	rr := httptest.NewRecorder()

	VersionHandler(rr, req)

	ct := rr.Header().Get("Content-Type")
	if ct != "application/json" {
		t.Errorf("expected Content-Type %q for 405 response, got %q", "application/json", ct)
	}

	var resp ErrorResponse
	if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode 405 response body: %v", err)
	}

	if resp.Error == "" {
		t.Error("expected non-empty error message in 405 response body")
	}
}

// TestVersionHandler_PUT_Returns405 verifies that a PUT request returns HTTP 405.
func TestVersionHandler_PUT_Returns405(t *testing.T) {
	req := httptest.NewRequest(http.MethodPut, "/version", nil)
	rr := httptest.NewRecorder()

	VersionHandler(rr, req)

	if rr.Code != http.StatusMethodNotAllowed {
		t.Errorf("expected status %d, got %d", http.StatusMethodNotAllowed, rr.Code)
	}
}

// TestVersionHandler_PUT_Returns405WithJSONError verifies that a PUT request returns a JSON error body.
func TestVersionHandler_PUT_Returns405WithJSONError(t *testing.T) {
	req := httptest.NewRequest(http.MethodPut, "/version", nil)
	rr := httptest.NewRecorder()

	VersionHandler(rr, req)

	ct := rr.Header().Get("Content-Type")
	if ct != "application/json" {
		t.Errorf("expected Content-Type %q for 405 response, got %q", "application/json", ct)
	}

	var resp ErrorResponse
	if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode 405 response body: %v", err)
	}

	if resp.Error == "" {
		t.Error("expected non-empty error message in 405 response body")
	}
}

// TestVersionHandler_DELETE_Returns405 verifies that a DELETE request returns HTTP 405.
func TestVersionHandler_DELETE_Returns405(t *testing.T) {
	req := httptest.NewRequest(http.MethodDelete, "/version", nil)
	rr := httptest.NewRecorder()

	VersionHandler(rr, req)

	if rr.Code != http.StatusMethodNotAllowed {
		t.Errorf("expected status %d, got %d", http.StatusMethodNotAllowed, rr.Code)
	}
}

// TestVersionHandler_DELETE_Returns405WithJSONError verifies that a DELETE request returns a JSON error body.
func TestVersionHandler_DELETE_Returns405WithJSONError(t *testing.T) {
	req := httptest.NewRequest(http.MethodDelete, "/version", nil)
	rr := httptest.NewRecorder()

	VersionHandler(rr, req)

	ct := rr.Header().Get("Content-Type")
	if ct != "application/json" {
		t.Errorf("expected Content-Type %q for 405 response, got %q", "application/json", ct)
	}

	var resp ErrorResponse
	if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode 405 response body: %v", err)
	}

	if resp.Error == "" {
		t.Error("expected non-empty error message in 405 response body")
	}
}

// TestVersionHandler_NonGETMethods_AllReturn405 is a table-driven test verifying all non-GET methods return 405.
func TestVersionHandler_NonGETMethods_AllReturn405(t *testing.T) {
	methods := []string{
		http.MethodPost,
		http.MethodPut,
		http.MethodDelete,
		http.MethodPatch,
		http.MethodHead,
		http.MethodOptions,
	}

	for _, method := range methods {
		t.Run(method, func(t *testing.T) {
			req := httptest.NewRequest(method, "/version", nil)
			rr := httptest.NewRecorder()

			VersionHandler(rr, req)

			if rr.Code != http.StatusMethodNotAllowed {
				t.Errorf("method %s: expected status %d, got %d", method, http.StatusMethodNotAllowed, rr.Code)
			}

			ct := rr.Header().Get("Content-Type")
			if ct != "application/json" {
				t.Errorf("method %s: expected Content-Type %q, got %q", method, "application/json", ct)
			}

			var resp ErrorResponse
			if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
				t.Fatalf("method %s: failed to decode 405 response body: %v", method, err)
			}

			if resp.Error == "" {
				t.Errorf("method %s: expected non-empty error message in 405 response", method)
			}
		})
	}
}
