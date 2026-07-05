// Package handler provides the Vercel Queue entry point.
package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/kristiannissen/concertlist/internal/adapters/etl/extractors/richter_gladsaxe"
	"github.com/kristiannissen/concertlist/internal/adapters/queue"
	"github.com/kristiannissen/concertlist/internal/domain"
)

// QueueHandler processes queue messages from Vercel Queues (push-based).
func QueueHandler(w http.ResponseWriter, r *http.Request) {
	// Only allow POST requests.
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Read the message body.
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close() //nolint:errcheck

	// Try to parse as ExtractionJob first, then as Concert
	var job domain.ExtractionJob
	if err := json.Unmarshal(body, &job); err != nil {
		// Try parsing as Concert
		var concert domain.Concert
		if err2 := json.Unmarshal(body, &concert); err2 != nil {
			http.Error(w, fmt.Sprintf("failed to parse message as ExtractionJob or Concert: %v", err), http.StatusBadRequest)
			return
		}
		// Handle concert message - save to storage
		log.Printf("Processing concert: %s", concert.Title)
		// TODO: Save concert to storage
		w.WriteHeader(http.StatusOK)
		return
	}

	// Initialize queue for extractor
	queueClient, err := queue.NewVercelQueueFromEnv()
	if err != nil {
		log.Printf("Failed to create queue client: %v", err)
		// Continue without queue - concerts will still be collected
	}

	// Initialize extractors (map venue name to extractor).
	extractors := map[string]domain.ExtractorPort{
		"richter_gladsaxe": richter_gladsaxe.NewExtractor(queueClient),
	}

	// Get the extractor for this venue.
	extractor, ok := extractors[job.Venue]
	if !ok {
		log.Printf("Unknown venue: %s", job.Venue)
		w.WriteHeader(http.StatusOK) // Acknowledge the message even if venue is unknown.
		return
	}

	// Run extraction for this venue.
	concerts, err := extractor.Extract(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// TODO: Save to storage once blob storage (internal/adapters/storage/blob)
	// is fully implemented. Wiring it in now would call unfinished stubs.

	log.Printf("Processed %d concerts for %s", len(concerts), job.Venue)
	w.WriteHeader(http.StatusOK) // Acknowledge successful processing.
}
