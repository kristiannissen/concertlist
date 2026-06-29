// Package queue provides Vercel Queues adapter.
package queue

import (
	"context"
	"encoding/json"
	"os"

	"github.com/kristiannissen/concertlist/internal/domain"
	"github.com/vercel/queues-go"
)

// VercelQueue implements domain.QueuePort.
type VercelQueue struct {
	client *queues.Client
}

// NewVercelQueue creates a new VercelQueue.
func NewVercelQueue() (*VercelQueue, error) {
	queueName := os.Getenv("QUEUE_NAME")
	queueToken := os.Getenv("QUEUE_TOKEN")
	if queueName == "" || queueToken == "" {
		return nil, os.ErrNotExist
	}
	client, err := queues.NewClient(queueName, queueToken)
	if err != nil {
		return nil, err
	}
	return &VercelQueue{client: client}, nil
}

// Enqueue adds a job to the queue.
func (q *VercelQueue) Enqueue(ctx context.Context, job domain.ExtractionJob) error {
	data, err := json.Marshal(job)
	if err != nil {
		return err
	}
	return q.client.Send(ctx, data)
}

// Process registers a handler for queue jobs.
func (q *VercelQueue) Process(ctx context.Context, handler domain.QueueHandler) error {
	return q.client.Process(ctx, func(ctx context.Context, job *queues.Job) error {
		var extractionJob domain.ExtractionJob
		if err := json.Unmarshal(job.Payload, &extractionJob); err != nil {
			return err
		}
		return handler(ctx, extractionJob)
	})
}
