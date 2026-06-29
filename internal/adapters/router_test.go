// Package adapters contains HTTP handlers and routing implementations.
package adapters

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewRouter_HealthCheck(t *testing.T) {
	t.Parallel()

	// Create the router
	router := NewRouter()

	// Create a request to the health check endpoint
	req, err := http.NewRequest(http.MethodGet, "/api/health", nil)
	if err != nil {
		t.Fatalf("http.NewRequest failed: %v", err)
	}

	// Create a response recorder
	rr := httptest.NewRecorder()

	// Serve the request
	router.ServeHTTP(rr, req)

	// Check the status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Check the response body
	expected := `{"status": "ok"}`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}

	// Check the content type
	contentType := rr.Header().Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf("handler returned wrong content type: got %v want %v", contentType, "application/json")
	}
}

func TestNewRouter_NotFound(t *testing.T) {
	t.Parallel()

	// Create the router
	router := NewRouter()

	// Create a request to a non-existent endpoint
	req, err := http.NewRequest(http.MethodGet, "/api/nonexistent", nil)
	if err != nil {
		t.Fatalf("http.NewRequest failed: %v", err)
	}

	// Create a response recorder
	rr := httptest.NewRecorder()

	// Serve the request
	router.ServeHTTP(rr, req)

	// Check the status code (should be 404 for non-existent routes)
	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusNotFound)
	}
}
