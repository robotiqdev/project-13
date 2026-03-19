package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
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
