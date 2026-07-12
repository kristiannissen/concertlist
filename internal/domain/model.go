// Package domain contains core business models and logic.
package domain

// MusicEvent represents a music event with schema.org context.
type MusicEvent struct {
	Context   string    `json:"@context"`
	Type      string    `json:"@type"`
	Name      string    `json:"name"`
	StartDate string    `json:"startDate"`
	Location  Location  `json:"location,omitempty"`
	Performer Performer `json:"performer,omitempty"`
	Offer     Offer     `json:"offers,omitempty"`
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

type Offer struct {
	Type          string `json:"@type"`
	URL           string `json:"url"`
	Price         int    `json:"price"`
	PriceCurrency string `json:"priceCurrency"`
}
