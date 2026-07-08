// Package domain contains core business models and logic.
package domain

import "time"

// MusicEvent represents a music event with schema.org context.
type MusicEvent struct {
	Context    string      `json:"@context"`
	Type       string      `json:"@type"`
	Name       string      `json:"name"`
	StartDate  string      `json:"startDate"`
	Location   Location    `json:"location"`
	Performers []Performer `json:"performer"`
}

// Location represents the venue where the event takes place.
type Location struct {
	Type    string        `json:"@type"`
	Name    string        `json:"name"`
	Address PostalAddress `json:"address"`
}

// PostalAddress represents the physical address of a location.
type PostalAddress struct {
	Type       string `json:"@type"`
	Street     string `json:"streetAddress"`
	Locality   string `json:"addressLocality"`
	PostalCode string `json:"postalCode"`
	Country    string `json:"addressCountry"`
}

// Performer represents a music group or artist performing at an event.
type Performer struct {
	Type string `json:"@type"`
	Name string `json:"name"`
}

// Concert represents a simplified concert event for internal use.
type Concert struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Venue       string `json:"venue"`
	Date        string `json:"date"`
	Description string `json:"description,omitempty"`
	URL         string `json:"url,omitempty"`
}

// ExtractionJob represents a job to extract data from a venue.
type ExtractionJob struct {
	Venue     string    `json:"venue"`
	Timestamp time.Time `json:"timestamp"`
}
