// Package handler provides the Vercel Queue entry point.
package handler

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/kristiannissen/concertlist/internal/adapters/etl/extractors/richter_gladsaxe"
	"github.com/kristiannissen/concertlist/internal/adapters/storage/blob"
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

	// Parse the ExtractionJob from the message.
	var job domain.ExtractionJob
	if err := json.Unmarshal(body, &job); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Initialize dependencies.
	blobStore, err := blob.NewBlobStore()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Initialize extractors (map venue name to extractor).
	extractors := map[string]domain.ExtractorPort{
		"richter_gladsaxe": richter_gladsaxe.NewExtractor(),
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

	// Save to storage.
	if err := blobStore.Save(r.Context(), concerts); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Printf("Processed %d concerts for %s", len(concerts), job.Venue)
	w.WriteHeader(http.StatusOK) // Acknowledge successful processing.
}
