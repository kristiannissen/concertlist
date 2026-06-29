// Package queue provides Vercel Queues adapter using the HTTP API.
package queue

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/kristiannissen/concertlist/internal/domain"
)

// VercelQueue implements domain.QueuePort using Vercel Queues HTTP API.
type VercelQueue struct {
	queueName string
	region    string // e.g., "iad1" (default Vercel region)
}

// NewVercelQueue creates a new VercelQueue.
func NewVercelQueue() (*VercelQueue, error) {
	queueName := os.Getenv("QUEUE_NAME")
	if queueName == "" {
		return nil, fmt.Errorf("QUEUE_NAME environment variable is required")
	}
	return &VercelQueue{
		queueName: queueName,
		region:    "iad1", // Default to Vercel's US East region.
	}, nil
}

// Enqueue sends a message to the Vercel Queue.
func (q *VercelQueue) Enqueue(ctx context.Context, job domain.ExtractionJob) error {
	// Marshal the job to JSON.
	data, err := json.Marshal(job)
	if err != nil {
		return err
	}

	// Create the request body.
	body := map[string]interface{}{
		"body": string(data),
	}
	bodyBytes, err := json.Marshal(body)
	if err != nil {
		return err
	}

	// Send the request to Vercel Queues HTTP API.
	url := fmt.Sprintf("https://%s.vercel-queue.com/api/v3/topic/%s", q.region, q.queueName)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(bodyBytes))
	if err != nil {
		return err
	}

	// Vercel automatically injects the OIDC token for authentication.
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close() //nolint:errcheck

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to enqueue job: %s (status: %d)", string(body), resp.StatusCode)
	}

	return nil
}

// Process is not needed for push-based consumers.
// Vercel automatically delivers messages to /api/queue.
func (q *VercelQueue) Process(ctx context.Context, handler domain.QueueHandler) error {
	return fmt.Errorf("push-based consumers do not require manual processing")
}
