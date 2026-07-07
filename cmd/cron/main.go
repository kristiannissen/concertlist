// Package main is the entry point for Vercel Cron Jobs.
package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/kristiannissen/concertlist/internal/adapters/queue"
	"github.com/kristiannissen/concertlist/internal/domain"
)

func main() {
	// Check if we should run in producer mode or consumer mode
	mode := os.Getenv("QUEUE_MODE")
	if mode == "" {
		mode = "producer" // Default to producer for backward compatibility
	}

	// Initialize queue.
	q, err := queue.NewVercelQueueFromEnv()
	if err != nil {
		log.Fatalf("Failed to initialize queue: %v", err)
	}

	switch mode {
	case "producer":
		runProducer(q)
	case "consumer":
		runConsumer(q)
	default:
		log.Fatalf("Unknown QUEUE_MODE: %s (must be 'producer' or 'consumer')", mode)
	}
}

// runProducer enqueues extraction jobs for all venues.
func runProducer(q *queue.VercelQueue) {
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

// runConsumer starts an async consumer to process messages from the queue.
func runConsumer(q *queue.VercelQueue) {
	// Create a handler for processing concerts
	handler := func(ctx context.Context, concert domain.Concert) error {
		log.Printf("Processing concert: %s at %s", concert.Title, concert.Date)
		
		// Simulate processing work
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(100 * time.Millisecond):
			return nil
		}
	}

	// Create async consumer with options
	consumer := queue.NewAsyncConsumer(
		q,
		handler,
		queue.WithConcurrency(10),
		queue.WithVisibilityTimeout(60*time.Second),
		queue.WithProcessTimeout(30*time.Second),
	)

	// Create context for graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Handle signals for graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	
	go func() {
		<-sigChan
		log.Println("Shutting down...")
		consumer.Stop()
		cancel()
	}()

	// Start the consumer
	log.Println("Starting async queue consumer...")
	if err := consumer.Start(ctx); err != nil {
		log.Fatalf("Consumer failed: %v", err)
	}
	
	log.Println("Consumer stopped")
}
