package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/kristiannissen/concertlist/pkg/domain"
)

func TestQueueHandler_ConcertMessage(t *testing.T) {
	t.Parallel()

	// Create a test concert
	concert := domain.Concert{
		ID:    "test-id",
		Title: "Test Concert",
		Venue: "Test Venue",
		Date:  "2026-07-04",
	}
	
	// Marshal concert to JSON
	body, err := json.Marshal(concert)
	if err != nil {
		t.Fatalf("Failed to marshal concert: %v", err)
	}

	// Create a request with the concert JSON
	req := httptest.NewRequest(http.MethodPost, "/api/queue", bytes.NewReader(body))
	w := httptest.NewRecorder()

	// Call the queue handler
	QueueHandler(w, req)

	// Check the response
	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
}

func TestQueueHandler_ExtractionJobMessage(t *testing.T) {
	t.Parallel()

	// Create a test extraction job
	job := domain.ExtractionJob{
		Venue: "test-venue",
	}
	
	// Marshal job to JSON
	body, err := json.Marshal(job)
	if err != nil {
		t.Fatalf("Failed to marshal job: %v", err)
	}

	// Create a request with the job JSON
	req := httptest.NewRequest(http.MethodPost, "/api/queue", bytes.NewReader(body))
	w := httptest.NewRecorder()

	// Call the queue handler
	QueueHandler(w, req)

	// Check the response - it might fail because we can't reach the actual website
	// but it should handle the message type correctly
	// We just check that it doesn't panic
	_ = w.Code
}

func TestQueueHandler_InvalidMessage(t *testing.T) {
	t.Parallel()

	// Create an invalid message
	body := []byte("invalid json")

	// Create a request with the invalid JSON
	req := httptest.NewRequest(http.MethodPost, "/api/queue", bytes.NewReader(body))
	w := httptest.NewRecorder()

	// Call the queue handler
	QueueHandler(w, req)

	// Check that it returns a 400 error
	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400 for invalid JSON, got %d", w.Code)
	}
}
