// Package queue provides Vercel Queues adapter using the HTTP API.
package queue

import (
	"context"
	"encoding/json"
	"log"
	"sync"
	"time"

	"github.com/kristiannissen/concertlist/internal/domain"
)

// AsyncConsumer provides asynchronous message processing for Vercel Queues.
// It implements a worker pool pattern for concurrent message processing.
type AsyncConsumer struct {
	queue             *VercelQueue
	concurrency       int
	visibilityTimeout time.Duration
	processTimeout    time.Duration
	handler           func(ctx context.Context, concert domain.Concert) error
	stopChan          chan struct{}
	wg                sync.WaitGroup
}

// AsyncConsumerOption is a function that configures an AsyncConsumer.
type AsyncConsumerOption func(*AsyncConsumer)

// WithConcurrency sets the number of concurrent workers.
func WithConcurrency(n int) AsyncConsumerOption {
	return func(c *AsyncConsumer) {
		if n > 0 {
			c.concurrency = n
		}
	}
}

// WithVisibilityTimeout sets the visibility timeout for received messages.
func WithVisibilityTimeout(d time.Duration) AsyncConsumerOption {
	return func(c *AsyncConsumer) {
		c.visibilityTimeout = d
	}
}

// WithProcessTimeout sets the timeout for processing a single message.
func WithProcessTimeout(d time.Duration) AsyncConsumerOption {
	return func(c *AsyncConsumer) {
		c.processTimeout = d
	}
}

// NewAsyncConsumer creates a new async consumer for the given queue.
func NewAsyncConsumer(queue *VercelQueue, handler func(ctx context.Context, concert domain.Concert) error, opts ...AsyncConsumerOption) *AsyncConsumer {
	c := &AsyncConsumer{
		queue:             queue,
		concurrency:       10,
		visibilityTimeout: 60 * time.Second,
		processTimeout:    30 * time.Second,
		handler:           handler,
		stopChan:          make(chan struct{}),
	}

	for _, opt := range opts {
		opt(c)
	}

	return c
}

// Start begins processing messages asynchronously.
// It creates a pool of workers that continuously poll for messages.
func (c *AsyncConsumer) Start(ctx context.Context) error {
	// Create a channel for message processing
	messageChan := make(chan domain.Concert, c.concurrency)

	// Start workers
	for i := 0; i < c.concurrency; i++ {
		c.wg.Add(1)
		go c.worker(ctx, messageChan)
	}

	// Start poller
	c.wg.Add(1)
	go c.poll(ctx, messageChan)

	// Wait for stop signal
	<-c.stopChan

	// Close channel and wait for workers to finish
	close(messageChan)
	c.wg.Wait()

	return nil
}

// Stop signals the consumer to stop processing messages.
func (c *AsyncConsumer) Stop() {
	close(c.stopChan)
}

// poll continuously retrieves messages from the queue and sends them to workers.
func (c *AsyncConsumer) poll(ctx context.Context, messageChan chan<- domain.Concert) {
	defer c.wg.Done()

	for {
		select {
		case <-ctx.Done():
			return
		case <-c.stopChan:
			return
		default:
			// Receive messages with visibility timeout
			opts := ReceiveMessagesOptions{
				MaxMessages:              c.concurrency,
				VisibilityTimeoutSeconds: int(c.visibilityTimeout.Seconds()),
				Accept:                   "application/x-ndjson",
			}

			resp, err := c.queue.internalReceiveMessages(ctx, opts)
			if err != nil {
				log.Printf("Error receiving messages: %v", err)
				// Backoff on error
				time.Sleep(1 * time.Second)
				continue
			}

			// Send messages to workers
			for _, msg := range resp.Messages {
				var concert domain.Concert
				// Try to unmarshal as Concert
				if err := json.Unmarshal(msg.Body, &concert); err != nil {
					log.Printf("Error unmarshaling message: %v", err)
					// Try to acknowledge the message anyway to avoid being stuck
					if msg.ReceiptHandle != "" {
						_ = c.queue.AcknowledgeMessage(ctx, msg.ReceiptHandle)
					}
					continue
				}

				select {
				case messageChan <- concert:
					// Message sent to worker
				case <-ctx.Done():
					return
				case <-c.stopChan:
					return
				}
			}

			// If no messages, add a small delay to avoid busy waiting
			if len(resp.Messages) == 0 {
				time.Sleep(100 * time.Millisecond)
			}
		}
	}
}

// worker processes messages from the channel.
func (c *AsyncConsumer) worker(ctx context.Context, messageChan <-chan domain.Concert) {
	defer c.wg.Done()

	for concert := range messageChan {
		// Create a context with timeout for processing
		processCtx, cancel := context.WithTimeout(ctx, c.processTimeout)

		// Process the message
		err := c.handler(processCtx, concert)
		cancel()

		if err != nil {
			log.Printf("Error processing concert %s: %v", concert.ID, err)
			// Message will become visible again after visibility timeout expires
			continue
		}

		log.Printf("Successfully processed concert: %s", concert.ID)
	}
}
