package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestWriteError_StatusCode verifies that writeError sets the correct HTTP status code.
func TestWriteError_StatusCode(t *testing.T) {
	rr := httptest.NewRecorder()
	writeError(rr, http.StatusNotFound, "not found")

	if rr.Code != http.StatusNotFound {
		t.Errorf("expected status %d, got %d", http.StatusNotFound, rr.Code)
	}
}

// TestWriteError_ContentTypeHeader verifies that writeError sets Content-Type to application/json.
func TestWriteError_ContentTypeHeader(t *testing.T) {
	rr := httptest.NewRecorder()
	writeError(rr, http.StatusNotFound, "not found")

	ct := rr.Header().Get("Content-Type")
	if ct != "application/json" {
		t.Errorf("expected Content-Type %q, got %q", "application/json", ct)
	}
}

// TestWriteError_BodyDecodesCorrectly verifies that the body decodes to the expected error message.
func TestWriteError_BodyDecodesCorrectly(t *testing.T) {
	rr := httptest.NewRecorder()
	writeError(rr, http.StatusNotFound, "not found")

	var resp ErrorResponse
	if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response body: %v", err)
	}

	if resp.Error != "not found" {
		t.Errorf("expected error message %q, got %q", "not found", resp.Error)
	}
}

// TestWriteError_NoExtraJSONFields verifies that the JSON body contains exactly the "error" field and nothing else.
func TestWriteError_NoExtraJSONFields(t *testing.T) {
	rr := httptest.NewRecorder()
	writeError(rr, http.StatusNotFound, "not found")

	var raw map[string]interface{}
	if err := json.NewDecoder(rr.Body).Decode(&raw); err != nil {
		t.Fatalf("failed to decode response body as map: %v", err)
	}

	if len(raw) != 1 {
		t.Errorf("expected exactly 1 field in JSON body, got %d: %v", len(raw), raw)
	}

	if _, ok := raw["error"]; !ok {
		t.Errorf("expected JSON field \"error\" to exist, got keys: %v", raw)
	}
}

// TestWriteError_ValidJSON verifies that the response body is valid JSON.
func TestWriteError_ValidJSON(t *testing.T) {
	rr := httptest.NewRecorder()
	writeError(rr, http.StatusBadRequest, "bad request")

	if !json.Valid(rr.Body.Bytes()) {
		t.Errorf("response body is not valid JSON: %s", rr.Body.String())
	}
}

// TestWriteError_DifferentStatusCodes verifies writeError works with various HTTP status codes.
func TestWriteError_DifferentStatusCodes(t *testing.T) {
	cases := []struct {
		status  int
		message string
	}{
		{http.StatusBadRequest, "bad request"},
		{http.StatusUnauthorized, "unauthorized"},
		{http.StatusForbidden, "forbidden"},
		{http.StatusInternalServerError, "internal server error"},
		{http.StatusMethodNotAllowed, "method not allowed"},
	}

	for _, tc := range cases {
		t.Run(http.StatusText(tc.status), func(t *testing.T) {
			rr := httptest.NewRecorder()
			writeError(rr, tc.status, tc.message)

			if rr.Code != tc.status {
				t.Errorf("expected status %d, got %d", tc.status, rr.Code)
			}

			var resp ErrorResponse
			if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
				t.Fatalf("failed to decode response body: %v", err)
			}

			if resp.Error != tc.message {
				t.Errorf("expected error %q, got %q", tc.message, resp.Error)
			}
		})
	}
}

// TestMethodNotAllowedHandler verifies that MethodNotAllowedHandler returns 405 with correct JSON body.
func TestMethodNotAllowedHandler(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rr := httptest.NewRecorder()
	MethodNotAllowedHandler(rr, req)

	if rr.Code != http.StatusMethodNotAllowed {
		t.Errorf("expected status %d, got %d", http.StatusMethodNotAllowed, rr.Code)
	}

	ct := rr.Header().Get("Content-Type")
	if ct != "application/json" {
		t.Errorf("expected Content-Type %q, got %q", "application/json", ct)
	}

	var resp ErrorResponse
	if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response body: %v", err)
	}

	if resp.Error != "method not allowed" {
		t.Errorf("expected error %q, got %q", "method not allowed", resp.Error)
	}
}

// TestMethodNotAllowedHandler_AnyMethod verifies MethodNotAllowedHandler works regardless of the request method.
func TestMethodNotAllowedHandler_AnyMethod(t *testing.T) {
	methods := []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete, http.MethodPatch}
	for _, method := range methods {
		t.Run(method, func(t *testing.T) {
			req := httptest.NewRequest(method, "/some/path", nil)
			rr := httptest.NewRecorder()
			MethodNotAllowedHandler(rr, req)

			if rr.Code != http.StatusMethodNotAllowed {
				t.Errorf("expected status %d, got %d", http.StatusMethodNotAllowed, rr.Code)
			}

			var resp ErrorResponse
			if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
				t.Fatalf("failed to decode response body: %v", err)
			}

			if resp.Error != "method not allowed" {
				t.Errorf("expected error %q, got %q", "method not allowed", resp.Error)
			}
		})
	}
}

// TestWriteError_EmptyMessage verifies writeError handles an empty message gracefully.
func TestWriteError_EmptyMessage(t *testing.T) {
	rr := httptest.NewRecorder()
	writeError(rr, http.StatusInternalServerError, "")

	if rr.Code != http.StatusInternalServerError {
		t.Errorf("expected status %d, got %d", http.StatusInternalServerError, rr.Code)
	}

	var resp ErrorResponse
	if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response body: %v", err)
	}

	if resp.Error != "" {
		t.Errorf("expected empty error message, got %q", resp.Error)
	}
}
