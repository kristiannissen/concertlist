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
	// Enqueue sends an extraction job to the queue.
	Enqueue(ctx context.Context, job ExtractionJob) error
	// EnqueueConcert sends a concert directly to the queue.
	EnqueueConcert(ctx context.Context, concert Concert) error
	// Process processes queue messages (for push-based consumers).
	Process(ctx context.Context, handler QueueHandler) error
	// ReceiveMessages retrieves messages from the queue for polling-based consumers.
	ReceiveMessages(ctx context.Context) ([]byte, error)
	// ReceiveMessageByID retrieves a specific message by its ID.
	ReceiveMessageByID(ctx context.Context, messageID string) ([]byte, error)
	// AcknowledgeMessage acknowledges successful processing of a message.
	AcknowledgeMessage(ctx context.Context, receiptHandle string) error
	// ExtendLease extends the visibility timeout of a message.
	ExtendLease(ctx context.Context, receiptHandle string, visibilityTimeoutSeconds int) error
}

// QueueHandler is a function that processes a queue job.
type QueueHandler func(ctx context.Context, job ExtractionJob) error
