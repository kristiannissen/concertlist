// Package domain contains pure business logic and JSON core models.
// This package must never import external frameworks or transport layers.
package domain

// Concert represents a concert event in the system.
type Concert struct {
	ID      string `json:"id"`
	Title   string `json:"title"`
	Artist  string `json:"artist"`
	Venue   string `json:"venue"`
	Date    string `json:"date"`
	Price   float64 `json:"price"`
	Tickets int     `json:"tickets_available"`
}

// CreateConcertRequest represents the request payload for creating a new concert.
type CreateConcertRequest struct {
	Title   string  `json:"title"`
	Artist  string  `json:"artist"`
	Venue   string  `json:"venue"`
	Date    string  `json:"date"`
	Price   float64 `json:"price"`
	Tickets int     `json:"tickets_available"`
}

// ErrorResponse represents a standardized error response.
type ErrorResponse struct {
	Error string `json:"error"`
}
