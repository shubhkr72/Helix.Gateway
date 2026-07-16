package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHealth(t *testing.T) {
	// Create a request
	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)

	// Record the response
	rr := httptest.NewRecorder()

	// Call the handler
	Health(rr, req)

	// Check status code
	if rr.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rr.Code)
	}

	// Check Content-Type
	expectedContentType := "application/json"
	if rr.Header().Get("Content-Type") != expectedContentType {
		t.Fatalf(
			"expected Content-Type %q, got %q",
			expectedContentType,
			rr.Header().Get("Content-Type"),
		)
	}

	// Decode and verify JSON response
	var response map[string]string
	if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if response["status"] != "ok" {
		t.Fatalf("expected status %q, got %q", "ok", response["status"])
	}
}
