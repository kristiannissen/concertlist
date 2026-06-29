// Package main is the entry point for Vercel Cron Jobs.
package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/kristiannissen/concertlist/internal/adapters/etl/extractors/richter_gladsaxe"
	"github.com/kristiannissen/concertlist/internal/adapters/queue"
	"github.com/kristiannissen/concertlist/internal/domain"
)

func main() {
	// Initialize queue.
	q, err := queue.NewVercelQueue()
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

	// Set up context with timeout (e.g., 1 minute for enqueuing).
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()

	// Handle OS signals for graceful shutdown.
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigChan
		cancel()
	}()

	<-ctx.Done()
	log.Println("Cron job completed")
}
