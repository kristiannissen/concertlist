// Package handler provides the Vercel Queue entry point.
package handler

import (
	"context"
	"log"
	"net/http"

	"github.com/kristiannissen/concertlist/internal/adapters/etl/extractors/richter_gladsaxe"
	"github.com/kristiannissen/concertlist/internal/adapters/storage/blob"
	"github.com/kristiannissen/concertlist/internal/adapters/queue"
	"github.com/kristiannissen/concertlist/internal/domain"
)

// QueueHandler processes queue jobs for Vercel Queues.
func QueueHandler(w http.ResponseWriter, r *http.Request) {
	// Initialize dependencies.
	blobStore, err := blob.NewBlobStore()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Initialize queue.
	q, err := queue.NewVercelQueue()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Initialize extractors (map venue name to extractor).
	extractors := map[string]domain.ExtractorPort{
		"richter_gladsaxe": richter_gladsaxe.NewExtractor(),
	}

	// Define queue job handler.
	jobHandler := func(ctx context.Context, job domain.ExtractionJob) error {
		extractor, ok := extractors[job.Venue]
		if !ok {
			log.Printf("Unknown venue: %s", job.Venue)
			return nil // Skip unknown venues.
		}

		// Run extraction for this venue.
		concerts, err := extractor.Extract(ctx)
		if err != nil {
			return err
		}

		// Save to storage.
		if err := blobStore.Save(ctx, concerts); err != nil {
			return err
		}

		log.Printf("Processed %d concerts for %s", len(concerts), job.Venue)
		return nil
	}

	// Process the queue.
	if err := q.Process(r.Context(), jobHandler); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
