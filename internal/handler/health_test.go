package handler_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/app/service/internal/handler"
	"github.com/app/service/internal/version"
)

// --- /health endpoint tests ---

// TestHealthHandlerReturns200 verifies that GET /health responds with HTTP 200 OK.
func TestHealthHandlerReturns200(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	w := httptest.NewRecorder()

	handler.HealthHandler(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("GET /health status = %d; want %d", w.Code, http.StatusOK)
	}
}

// TestHealthHandlerReturnsStatusOK verifies that the health response body
// contains exactly {"status":"ok"}.
func TestHealthHandlerReturnsStatusOK(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	w := httptest.NewRecorder()

	handler.HealthHandler(w, req)

	var got map[string]string
	if err := json.NewDecoder(w.Body).Decode(&got); err != nil {
		t.Fatalf("GET /health body is not valid JSON: %v", err)
	}

	status, ok := got["status"]
	if !ok {
		t.Fatal("GET /health response missing 'status' field")
	}
	if status != "ok" {
		t.Errorf("GET /health status = %q; want %q", status, "ok")
	}
}

// TestHealthHandlerResponseIsMinimal verifies that /health does NOT include
// version fields. Load balancers expect a minimal body; extra fields may cause issues.
func TestHealthHandlerResponseIsMinimal(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	w := httptest.NewRecorder()

	handler.HealthHandler(w, req)

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

// TestHealthHandlerSetsContentTypeJSON verifies that Content-Type is application/json.
func TestHealthHandlerSetsContentTypeJSON(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	w := httptest.NewRecorder()

	handler.HealthHandler(w, req)

	ct := w.Header().Get("Content-Type")
	if ct != "application/json" {
		t.Errorf("GET /health Content-Type = %q; want %q", ct, "application/json")
	}
}

// --- /version endpoint tests ---

// TestVersionHandlerReturns200 verifies that GET /version responds with HTTP 200 OK.
func TestVersionHandlerReturns200(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/version", nil)
	w := httptest.NewRecorder()

	handler.VersionHandler(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("GET /version status = %d; want %d", w.Code, http.StatusOK)
	}
}

// TestVersionHandlerReturnsValidJSON verifies that /version responds with valid JSON.
func TestVersionHandlerReturnsValidJSON(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/version", nil)
	w := httptest.NewRecorder()

	handler.VersionHandler(w, req)

	var got map[string]interface{}
	if err := json.NewDecoder(w.Body).Decode(&got); err != nil {
		t.Fatalf("GET /version body is not valid JSON: %v", err)
	}
}

// TestVersionHandlerResponseContainsVersionField verifies the "version" field is present.
func TestVersionHandlerResponseContainsVersionField(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/version", nil)
	w := httptest.NewRecorder()

	handler.VersionHandler(w, req)

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

	handler.VersionHandler(w, req)

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

	handler.VersionHandler(w, req)

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
	// Set known values into the version package vars.
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

	handler.VersionHandler(w, req)

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

	handler.VersionHandler(w, req)

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

	handler.VersionHandler(w, req)

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
// not affect the response of /version, and vice versa — the two endpoints are fully
// independent (no shared mutable state side-effects).
func TestHealthAndVersionEndpointsAreIndependent(t *testing.T) {
	// Call /health then /version.
	reqH := httptest.NewRequest(http.MethodGet, "/health", nil)
	wH := httptest.NewRecorder()
	handler.HealthHandler(wH, reqH)

	reqV := httptest.NewRequest(http.MethodGet, "/version", nil)
	wV := httptest.NewRecorder()
	handler.VersionHandler(wV, reqV)

	// Health must be minimal.
	var health map[string]interface{}
	if err := json.NewDecoder(wH.Body).Decode(&health); err != nil {
		t.Fatalf("/health body is not valid JSON: %v", err)
	}
	if health["status"] != "ok" {
		t.Errorf("/health status = %q; want %q", health["status"], "ok")
	}

	// Version must have build info fields.
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
