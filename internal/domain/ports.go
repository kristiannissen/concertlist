// Package domain contains the core business models and interfaces.
package domain

// Extractor defines the interface for site-specific data extractors.
type Extractor interface {
	// Extract retrieves events from a specific source and returns them as MusicEvents.
	Extract() ([]MusicEvent, error)
}

// BlobStorage defines the interface for storing events in blob storage.
type BlobStorage interface {
	// Store saves the given events to blob storage.
	Store(events []MusicEvent) error
	// GetAll retrieves all stored events.
	GetAll() ([]MusicEvent, error)
}

// ETLService defines the interface for the ETL workflow.
type ETLService interface {
	// Run executes the full ETL workflow: extract from all sources, transform, and load.
	Run() ([]MusicEvent, error)
}
