// Package domain contains core business models and logic.
package domain

import "context"

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

// ExtractorPort defines the interface for extracting concert data from a venue.
type ExtractorPort interface {
	Extract(ctx context.Context) ([]Concert, error)
}

// StoragePort defines the interface for storing concerts.
type StoragePort interface {
	Save(ctx context.Context, concerts []Concert) error
	Load(ctx context.Context) ([]Concert, error)
}

// QueuePort defines the interface for queue operations.
type QueuePort interface {
	Enqueue(ctx context.Context, job ExtractionJob) error
	EnqueueConcert(ctx context.Context, concert Concert) error
	Process(ctx context.Context, handler QueueHandler) error
}

// QueueHandler is a function that processes a queue job.
type QueueHandler func(ctx context.Context, job ExtractionJob) error
