// Package main is the entry point for Vercel Cron Jobs.
package main

import (
	"context"
	"log"
	"time"

	"github.com/kristiannissen/concertlist/pkg/adapters/queue"
	"github.com/kristiannissen/concertlist/pkg/domain"
)

func main() {
	// Initialize queue.
	q, err := queue.NewVercelQueueFromEnv()
	if err != nil {
		log.Fatalf("Failed to initialize queue: %v", err)
	}

	// Enqueue extraction jobs for all venues.
	venues := []string{"richter_gladsaxe"} // Add more venues later.
	for _, venue := range venues {
		job := domain.ExtractionJob{
			Venue:     venue,
			Timestamp: time.Now().UTC(),
		}
		if err := q.Enqueue(context.Background(), job); err != nil {
			log.Printf("Failed to enqueue job for %s: %v", venue, err)
		} else {
			log.Printf("Enqueued job for %s", venue)
		}
	}
}
