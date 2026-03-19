package handler_test

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/robotiqdev/project-13/internal/router"
	"github.com/robotiqdev/project-13/internal/version"
)

// newTestServer creates a real httptest.Server backed by the full application router.
func newTestServer(t *testing.T) *httptest.Server {
	t.Helper()
	srv := httptest.NewServer(router.New())
	t.Cleanup(srv.Close)
	return srv
}

// TestVersionEndpoint_Integration_GET_Returns200 verifies that a real HTTP GET to /version returns 200 OK.
func TestVersionEndpoint_Integration_GET_Returns200(t *testing.T) {
	version.Version = "test-v"
	version.Commit = "test-commit"
	version.BuildDate = "test-date"

	srv := newTestServer(t)

	resp, err := http.Get(srv.URL + "/version")
	if err != nil {
		t.Fatalf("GET /version failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, resp.StatusCode)
	}
}

// TestVersionEndpoint_Integration_GET_ContentTypeJSON verifies that a real HTTP GET to /version returns Content-Type application/json.
func TestVersionEndpoint_Integration_GET_ContentTypeJSON(t *testing.T) {
	version.Version = "test-v"
	version.Commit = "test-commit"
	version.BuildDate = "test-date"

	srv := newTestServer(t)

	resp, err := http.Get(srv.URL + "/version")
	if err != nil {
		t.Fatalf("GET /version failed: %v", err)
	}
	defer resp.Body.Close()

	ct := resp.Header.Get("Content-Type")
	if ct != "application/json" {
		t.Errorf("expected Content-Type %q, got %q", "application/json", ct)
	}
}

// TestVersionEndpoint_Integration_GET_BodyIsValidJSON verifies that the response body is valid JSON.
func TestVersionEndpoint_Integration_GET_BodyIsValidJSON(t *testing.T) {
	version.Version = "test-v"
	version.Commit = "test-commit"
	version.BuildDate = "test-date"

	srv := newTestServer(t)

	resp, err := http.Get(srv.URL + "/version")
	if err != nil {
		t.Fatalf("GET /version failed: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("failed to read response body: %v", err)
	}

	if !json.Valid(body) {
		t.Errorf("response body is not valid JSON: %s", string(body))
	}
}

// TestVersionEndpoint_Integration_GET_AllFieldsPresent verifies that the response body contains
// the three required fields: version, commit, and build_date as non-empty strings.
func TestVersionEndpoint_Integration_GET_AllFieldsPresent(t *testing.T) {
	version.Version = "test-v"
	version.Commit = "test-commit"
	version.BuildDate = "test-date"

	srv := newTestServer(t)

	resp, err := http.Get(srv.URL + "/version")
	if err != nil {
		t.Fatalf("GET /version failed: %v", err)
	}
	defer resp.Body.Close()

	var info version.BuildInfo
	if err := json.NewDecoder(resp.Body).Decode(&info); err != nil {
		t.Fatalf("failed to decode response body: %v", err)
	}

	if info.Version == "" {
		t.Error("expected non-empty 'version' field in response body")
	}
	if info.Commit == "" {
		t.Error("expected non-empty 'commit' field in response body")
	}
	if info.BuildDate == "" {
		t.Error("expected non-empty 'build_date' field in response body")
	}
}

// TestVersionEndpoint_Integration_GET_BodyReflectsKnownValues verifies that the response body
// reflects the exact values set on the version package variables.
func TestVersionEndpoint_Integration_GET_BodyReflectsKnownValues(t *testing.T) {
	version.Version = "test-v"
	version.Commit = "abc123"
	version.BuildDate = "2026-03-19T00:00:00Z"

	srv := newTestServer(t)

	resp, err := http.Get(srv.URL + "/version")
	if err != nil {
		t.Fatalf("GET /version failed: %v", err)
	}
	defer resp.Body.Close()

	var info version.BuildInfo
	if err := json.NewDecoder(resp.Body).Decode(&info); err != nil {
		t.Fatalf("failed to decode response body: %v", err)
	}

	if info.Version != "test-v" {
		t.Errorf("expected Version %q, got %q", "test-v", info.Version)
	}
	if info.Commit != "abc123" {
		t.Errorf("expected Commit %q, got %q", "abc123", info.Commit)
	}
	if info.BuildDate != "2026-03-19T00:00:00Z" {
		t.Errorf("expected BuildDate %q, got %q", "2026-03-19T00:00:00Z", info.BuildDate)
	}
}

// TestVersionEndpoint_Integration_GET_BodyAsMap verifies all three JSON keys are present
// using a raw map decode (not a typed struct), ensuring key names match the JSON spec.
func TestVersionEndpoint_Integration_GET_BodyAsMap(t *testing.T) {
	version.Version = "test-v"
	version.Commit = "test-commit"
	version.BuildDate = "test-date"

	srv := newTestServer(t)

	resp, err := http.Get(srv.URL + "/version")
	if err != nil {
		t.Fatalf("GET /version failed: %v", err)
	}
	defer resp.Body.Close()

	var raw map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&raw); err != nil {
		t.Fatalf("failed to decode response body into map: %v", err)
	}

	requiredKeys := []string{"version", "commit", "build_date"}
	for _, key := range requiredKeys {
		val, ok := raw[key]
		if !ok {
			t.Errorf("expected JSON key %q to be present in response body", key)
			continue
		}
		str, ok := val.(string)
		if !ok {
			t.Errorf("expected JSON key %q to be a string, got %T", key, val)
			continue
		}
		if str == "" {
			t.Errorf("expected JSON key %q to be a non-empty string", key)
		}
	}
}

// TestVersionEndpoint_Integration_POST_Returns405 verifies that a real HTTP POST to /version returns 405.
func TestVersionEndpoint_Integration_POST_Returns405(t *testing.T) {
	srv := newTestServer(t)

	resp, err := http.Post(srv.URL+"/version", "application/json", nil)
	if err != nil {
		t.Fatalf("POST /version failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusMethodNotAllowed {
		t.Errorf("expected status %d for POST, got %d", http.StatusMethodNotAllowed, resp.StatusCode)
	}
}

// TestVersionEndpoint_Integration_NonGETMethods_AllReturn405 is a table-driven integration test
// verifying that all non-GET methods return 405 through the full router stack.
func TestVersionEndpoint_Integration_NonGETMethods_AllReturn405(t *testing.T) {
	srv := newTestServer(t)

	methods := []string{
		http.MethodPost,
		http.MethodPut,
		http.MethodDelete,
		http.MethodPatch,
	}

	for _, method := range methods {
		t.Run(method, func(t *testing.T) {
			req, err := http.NewRequest(method, srv.URL+"/version", nil)
			if err != nil {
				t.Fatalf("failed to create %s request: %v", method, err)
			}

			client := &http.Client{}
			resp, err := client.Do(req)
			if err != nil {
				t.Fatalf("%s /version failed: %v", method, err)
			}
			defer resp.Body.Close()

			if resp.StatusCode != http.StatusMethodNotAllowed {
				t.Errorf("method %s: expected status %d, got %d", method, http.StatusMethodNotAllowed, resp.StatusCode)
			}

			ct := resp.Header.Get("Content-Type")
			if ct != "application/json" {
				t.Errorf("method %s: expected Content-Type %q, got %q", method, "application/json", ct)
			}
		})
	}
}
