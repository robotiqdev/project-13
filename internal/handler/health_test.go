package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/robotiqdev/project-13/internal/version"
)

// TestHealthHandler_GetReturns200 verifies that a GET request returns HTTP 200.
func TestHealthHandler_GetReturns200(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rr := httptest.NewRecorder()

	HealthHandler(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, rr.Code)
	}
}

// TestHealthHandler_GetReturnsJSONContentType verifies that a GET request sets Content-Type to application/json.
func TestHealthHandler_GetReturnsJSONContentType(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rr := httptest.NewRecorder()

	HealthHandler(rr, req)

	ct := rr.Header().Get("Content-Type")
	if ct != "application/json" {
		t.Errorf("expected Content-Type %q, got %q", "application/json", ct)
	}
}

// TestHealthHandler_GetReturnsStatusOKBody verifies that the JSON body contains {"status":"ok"}.
func TestHealthHandler_GetReturnsStatusOKBody(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rr := httptest.NewRecorder()

	HealthHandler(rr, req)

	var resp HealthResponse
	if err := json.Unmarshal(rr.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to unmarshal response body: %v", err)
	}

	if resp.Status != "ok" {
		t.Errorf("expected status field %q, got %q", "ok", resp.Status)
	}
}

// TestHealthHandler_GetBodyIsValidJSON verifies that the response body is valid JSON.
func TestHealthHandler_GetBodyIsValidJSON(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rr := httptest.NewRecorder()

	HealthHandler(rr, req)

	if !json.Valid(rr.Body.Bytes()) {
		t.Errorf("response body is not valid JSON: %s", rr.Body.String())
	}
}

// TestHealthHandler_GetBodyHasExactlyOneField verifies that the JSON body has exactly one field ("status").
func TestHealthHandler_GetBodyHasExactlyOneField(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rr := httptest.NewRecorder()

	HealthHandler(rr, req)

	var raw map[string]interface{}
	if err := json.Unmarshal(rr.Body.Bytes(), &raw); err != nil {
		t.Fatalf("failed to unmarshal response body as map: %v", err)
	}

	if len(raw) != 1 {
		t.Errorf("expected exactly 1 field in JSON body, got %d: %v", len(raw), raw)
	}

	if _, ok := raw["status"]; !ok {
		t.Errorf("expected JSON field \"status\" to exist, got keys: %v", raw)
	}
}

// TestHealthHandler_NonGetMethodsReturn405 verifies that POST, PUT, and DELETE return 405.
func TestHealthHandler_NonGetMethodsReturn405(t *testing.T) {
	methods := []string{http.MethodPost, http.MethodPut, http.MethodDelete}

	for _, method := range methods {
		t.Run(method, func(t *testing.T) {
			req := httptest.NewRequest(method, "/health", nil)
			rr := httptest.NewRecorder()

			HealthHandler(rr, req)

			if rr.Code != http.StatusMethodNotAllowed {
				t.Errorf("expected status %d, got %d", http.StatusMethodNotAllowed, rr.Code)
			}
		})
	}
}

// TestHealthHandler_NonGetMethodsReturnJSONErrorBody verifies that non-GET methods return a JSON error body.
func TestHealthHandler_NonGetMethodsReturnJSONErrorBody(t *testing.T) {
	methods := []string{http.MethodPost, http.MethodPut, http.MethodDelete}

	for _, method := range methods {
		t.Run(method, func(t *testing.T) {
			req := httptest.NewRequest(method, "/health", nil)
			rr := httptest.NewRecorder()

			HealthHandler(rr, req)

			var resp ErrorResponse
			if err := json.Unmarshal(rr.Body.Bytes(), &resp); err != nil {
				t.Fatalf("failed to unmarshal response body: %v", err)
			}

			if resp.Error == "" {
				t.Errorf("expected non-empty error field in response body")
			}
		})
	}
}

// TestHealthHandler_NonGetMethodsReturnJSONContentType verifies that non-GET methods set Content-Type to application/json.
func TestHealthHandler_NonGetMethodsReturnJSONContentType(t *testing.T) {
	methods := []string{http.MethodPost, http.MethodPut, http.MethodDelete}

	for _, method := range methods {
		t.Run(method, func(t *testing.T) {
			req := httptest.NewRequest(method, "/health", nil)
			rr := httptest.NewRecorder()

			HealthHandler(rr, req)

			ct := rr.Header().Get("Content-Type")
			if ct != "application/json" {
				t.Errorf("expected Content-Type %q, got %q", "application/json", ct)
			}
		})
	}
}

// TestHealthHandlerResponseIsMinimal verifies that /health does NOT include
// version fields. Load balancers expect a minimal body; extra fields may cause issues.
func TestHealthHandlerResponseIsMinimal(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	w := httptest.NewRecorder()

	HealthHandler(w, req)

	var got map[string]interface{}
	if err := json.NewDecoder(w.Body).Decode(&got); err != nil {
		t.Fatalf("GET /health body is not valid JSON: %v", err)
	}

	for _, forbidden := range []string{"version", "commit", "buildDate"} {
		if _, present := got[forbidden]; present {
			t.Errorf("GET /health response must not contain field %q (keep health minimal)", forbidden)
		}
	}
}

// --- /version endpoint tests ---

// TestVersionHandlerReturns200 verifies that GET /version responds with HTTP 200 OK.
func TestVersionHandlerReturns200(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/version", nil)
	w := httptest.NewRecorder()

	VersionHandler(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("GET /version status = %d; want %d", w.Code, http.StatusOK)
	}
}

// TestVersionHandlerReturnsValidJSON verifies that /version responds with valid JSON.
func TestVersionHandlerReturnsValidJSON(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/version", nil)
	w := httptest.NewRecorder()

	VersionHandler(w, req)

	var got map[string]interface{}
	if err := json.NewDecoder(w.Body).Decode(&got); err != nil {
		t.Fatalf("GET /version body is not valid JSON: %v", err)
	}
}

// TestVersionHandlerResponseContainsVersionField verifies the "version" field is present.
func TestVersionHandlerResponseContainsVersionField(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/version", nil)
	w := httptest.NewRecorder()

	VersionHandler(w, req)

	var got map[string]interface{}
	if err := json.NewDecoder(w.Body).Decode(&got); err != nil {
		t.Fatalf("GET /version body is not valid JSON: %v", err)
	}

	if _, ok := got["version"]; !ok {
		t.Error("GET /version response missing 'version' field")
	}
}

// TestVersionHandlerResponseContainsCommitField verifies the "commit" field is present.
func TestVersionHandlerResponseContainsCommitField(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/version", nil)
	w := httptest.NewRecorder()

	VersionHandler(w, req)

	var got map[string]interface{}
	if err := json.NewDecoder(w.Body).Decode(&got); err != nil {
		t.Fatalf("GET /version body is not valid JSON: %v", err)
	}

	if _, ok := got["commit"]; !ok {
		t.Error("GET /version response missing 'commit' field")
	}
}

// TestVersionHandlerResponseContainsBuildDateField verifies the "buildDate" field is present.
func TestVersionHandlerResponseContainsBuildDateField(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/version", nil)
	w := httptest.NewRecorder()

	VersionHandler(w, req)

	var got map[string]interface{}
	if err := json.NewDecoder(w.Body).Decode(&got); err != nil {
		t.Fatalf("GET /version body is not valid JSON: %v", err)
	}

	if _, ok := got["buildDate"]; !ok {
		t.Error("GET /version response missing 'buildDate' field")
	}
}

// TestVersionHandlerReflectsVersionPackageVars verifies that /version returns values
// that match the current state of the version package variables.
func TestVersionHandlerReflectsVersionPackageVars(t *testing.T) {
	origVersion := version.Version
	origCommit := version.Commit
	origBuildDate := version.BuildDate
	defer func() {
		version.Version = origVersion
		version.Commit = origCommit
		version.BuildDate = origBuildDate
	}()

	version.Version = "v1.2.3"
	version.Commit = "abc123def456"
	version.BuildDate = "2026-03-19T00:00:00Z"

	req := httptest.NewRequest(http.MethodGet, "/version", nil)
	w := httptest.NewRecorder()

	VersionHandler(w, req)

	var got map[string]string
	if err := json.NewDecoder(w.Body).Decode(&got); err != nil {
		t.Fatalf("GET /version body is not valid JSON: %v", err)
	}

	if got["version"] != version.Version {
		t.Errorf("GET /version 'version' = %q; want %q", got["version"], version.Version)
	}
	if got["commit"] != version.Commit {
		t.Errorf("GET /version 'commit' = %q; want %q", got["commit"], version.Commit)
	}
	if got["buildDate"] != version.BuildDate {
		t.Errorf("GET /version 'buildDate' = %q; want %q", got["buildDate"], version.BuildDate)
	}
}

// TestVersionHandlerSetsContentTypeJSON verifies that Content-Type is application/json.
func TestVersionHandlerSetsContentTypeJSON(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/version", nil)
	w := httptest.NewRecorder()

	VersionHandler(w, req)

	ct := w.Header().Get("Content-Type")
	if ct != "application/json" {
		t.Errorf("GET /version Content-Type = %q; want %q", ct, "application/json")
	}
}

// TestVersionHandlerEmptyVarsReturnEmptyStrings verifies that when no ldflags are
// injected, the /version response still contains all three fields with empty strings.
func TestVersionHandlerEmptyVarsReturnEmptyStrings(t *testing.T) {
	origVersion := version.Version
	origCommit := version.Commit
	origBuildDate := version.BuildDate
	defer func() {
		version.Version = origVersion
		version.Commit = origCommit
		version.BuildDate = origBuildDate
	}()

	version.Version = ""
	version.Commit = ""
	version.BuildDate = ""

	req := httptest.NewRequest(http.MethodGet, "/version", nil)
	w := httptest.NewRecorder()

	VersionHandler(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("GET /version with empty vars: status = %d; want %d", w.Code, http.StatusOK)
	}

	var got map[string]interface{}
	if err := json.NewDecoder(w.Body).Decode(&got); err != nil {
		t.Fatalf("GET /version body is not valid JSON: %v", err)
	}

	for _, field := range []string{"version", "commit", "buildDate"} {
		if _, ok := got[field]; !ok {
			t.Errorf("GET /version missing field %q even when empty", field)
		}
	}
}

// TestHealthAndVersionEndpointsAreIndependent verifies that calling /health does
// not affect the response of /version, and vice versa.
func TestHealthAndVersionEndpointsAreIndependent(t *testing.T) {
	reqH := httptest.NewRequest(http.MethodGet, "/health", nil)
	wH := httptest.NewRecorder()
	HealthHandler(wH, reqH)

	reqV := httptest.NewRequest(http.MethodGet, "/version", nil)
	wV := httptest.NewRecorder()
	VersionHandler(wV, reqV)

	var health map[string]interface{}
	if err := json.NewDecoder(wH.Body).Decode(&health); err != nil {
		t.Fatalf("/health body is not valid JSON: %v", err)
	}
	if health["status"] != "ok" {
		t.Errorf("/health status = %q; want %q", health["status"], "ok")
	}

	var ver map[string]interface{}
	if err := json.NewDecoder(wV.Body).Decode(&ver); err != nil {
		t.Fatalf("/version body is not valid JSON: %v", err)
	}
	for _, field := range []string{"version", "commit", "buildDate"} {
		if _, ok := ver[field]; !ok {
			t.Errorf("/version missing field %q", field)
		}
	}
}
