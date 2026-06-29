// Package domain contains pure business logic and JSON core models.
package domain

import (
	"encoding/json"
	"testing"
)

func TestConcert_JSONMarshal(t *testing.T) {
	t.Parallel()

	concert := Concert{
		ID:      "123",
		Title:   "Test Concert",
		Artist:  "Test Artist",
		Venue:   "Test Venue",
		Date:    "2024-01-01",
		Price:   25.50,
		Tickets: 100,
	}

	got, err := json.Marshal(concert)
	if err != nil {
		t.Fatalf("json.Marshal failed: %v", err)
	}

	want := `{"id":"123","title":"Test Concert","artist":"Test Artist","venue":"Test Venue","date":"2024-01-01","price":25.5,"tickets_available":100}`
	if string(got) != want {
		t.Errorf("json.Marshal() = %s, want %s", string(got), want)
	}
}

func TestConcert_JSONUnmarshal(t *testing.T) {
	t.Parallel()

	jsonStr := `{"id":"123","title":"Test Concert","artist":"Test Artist","venue":"Test Venue","date":"2024-01-01","price":25.5,"tickets_available":100}`
	var concert Concert
	if err := json.Unmarshal([]byte(jsonStr), &concert); err != nil {
		t.Fatalf("json.Unmarshal failed: %v", err)
	}

	if concert.ID != "123" {
		t.Errorf("concert.ID = %s, want 123", concert.ID)
	}
	if concert.Title != "Test Concert" {
		t.Errorf("concert.Title = %s, want Test Concert", concert.Title)
	}
	if concert.Artist != "Test Artist" {
		t.Errorf("concert.Artist = %s, want Test Artist", concert.Artist)
	}
	if concert.Venue != "Test Venue" {
		t.Errorf("concert.Venue = %s, want Test Venue", concert.Venue)
	}
	if concert.Date != "2024-01-01" {
		t.Errorf("concert.Date = %s, want 2024-01-01", concert.Date)
	}
	if concert.Price != 25.50 {
		t.Errorf("concert.Price = %f, want 25.50", concert.Price)
	}
	if concert.Tickets != 100 {
		t.Errorf("concert.Tickets = %d, want 100", concert.Tickets)
	}
}

func TestCreateConcertRequest_JSONMarshal(t *testing.T) {
	t.Parallel()

	req := CreateConcertRequest{
		Title:   "New Concert",
		Artist:  "New Artist",
		Venue:   "New Venue",
		Date:    "2024-12-31",
		Price:   30.00,
		Tickets: 200,
	}

	got, err := json.Marshal(req)
	if err != nil {
		t.Fatalf("json.Marshal failed: %v", err)
	}

	want := `{"title":"New Concert","artist":"New Artist","venue":"New Venue","date":"2024-12-31","price":30,"tickets_available":200}`
	if string(got) != want {
		t.Errorf("json.Marshal() = %s, want %s", string(got), want)
	}
}

func TestCreateConcertRequest_JSONUnmarshal(t *testing.T) {
	t.Parallel()

	jsonStr := `{"title":"New Concert","artist":"New Artist","venue":"New Venue","date":"2024-12-31","price":30,"tickets_available":200}`
	var req CreateConcertRequest
	if err := json.Unmarshal([]byte(jsonStr), &req); err != nil {
		t.Fatalf("json.Unmarshal failed: %v", err)
	}

	if req.Title != "New Concert" {
		t.Errorf("req.Title = %s, want New Concert", req.Title)
	}
	if req.Artist != "New Artist" {
		t.Errorf("req.Artist = %s, want New Artist", req.Artist)
	}
	if req.Venue != "New Venue" {
		t.Errorf("req.Venue = %s, want New Venue", req.Venue)
	}
	if req.Date != "2024-12-31" {
		t.Errorf("req.Date = %s, want 2024-12-31", req.Date)
	}
	if req.Price != 30.00 {
		t.Errorf("req.Price = %f, want 30.00", req.Price)
	}
	if req.Tickets != 200 {
		t.Errorf("req.Tickets = %d, want 200", req.Tickets)
	}
}

func TestErrorResponse_JSONMarshal(t *testing.T) {
	t.Parallel()

	errResp := ErrorResponse{Error: "not found"}

	got, err := json.Marshal(errResp)
	if err != nil {
		t.Fatalf("json.Marshal failed: %v", err)
	}

	want := `{"error":"not found"}`
	if string(got) != want {
		t.Errorf("json.Marshal() = %s, want %s", string(got), want)
	}
}

func TestErrorResponse_JSONUnmarshal(t *testing.T) {
	t.Parallel()

	jsonStr := `{"error":"not found"}`
	var errResp ErrorResponse
	if err := json.Unmarshal([]byte(jsonStr), &errResp); err != nil {
		t.Fatalf("json.Unmarshal failed: %v", err)
	}

	if errResp.Error != "not found" {
		t.Errorf("errResp.Error = %s, want not found", errResp.Error)
	}
}
