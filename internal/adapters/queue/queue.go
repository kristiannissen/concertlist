url: https://raw.githubusercontent.com/kristiannissen/concertlist/main/internal/adapters/queue/queue.go

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

// DefaultRegion is the Vercel region used when none is specified in the config.
const DefaultRegion = "iad1"

// DefaultConsumer is the consumer name used when none is specified in the config.
const DefaultConsumer = "default"

// VercelQueue implements domain.QueuePort using Vercel Queues HTTP API.
type VercelQueue struct {
	config   QueueConfig
	client   *http.Client
	buildURL func(path string) string
}

// NewVercelQueue creates a new VercelQueue from the given configuration.
// Region defaults to DefaultRegion and Consumer defaults to DefaultConsumer
// when left empty.
func NewVercelQueue(config QueueConfig) (*VercelQueue, error) {
	if config.Topic == "" {
		return nil, fmt.Errorf("topic is required")
	}
	if config.Region == "" {
		config.Region = DefaultRegion
	}
	if config.Consumer == "" {
		config.Consumer = DefaultConsumer
	}

	q := &VercelQueue{
		config: config,
		client: http.DefaultClient,
	}
	q.buildURL = func(path string) string {
		return fmt.Sprintf("https://%s.vercel-queue.com/api/v3%s", q.config.Region, path)
	}

	return q, nil
}

// NewVercelQueueFromEnv creates a new VercelQueue using environment variables:
// QUEUE_NAME (required, used as Topic), QUEUE_REGION, QUEUE_CONSUMER, and
// VERCEL_OIDC_TOKEN (injected automatically by Vercel at runtime).
func NewVercelQueueFromEnv() (*VercelQueue, error) {
	topic := os.Getenv("QUEUE_NAME")
	if topic == "" {
		return nil, fmt.Errorf("QUEUE_NAME environment variable is required")
	}

	return NewVercelQueue(QueueConfig{
		Topic:     topic,
		Region:    os.Getenv("QUEUE_REGION"),
		Consumer:  os.Getenv("QUEUE_CONSUMER"),
		OIDCToken: os.Getenv("VERCEL_OIDC_TOKEN"),
	})
}

// authHeader sets the Authorization header on the request when an OIDC token is configured.
func (q *VercelQueue) aut
hHeader(req *http.Request) {
	if q.config.OIDCToken != "" {
		req.Header.Set("Authorization", "Bearer "+q.config.OIDCToken)
	}
}

// SendMessage sends a single message to the configured topic.
func (q *VercelQueue) SendMessage(ctx context.Context, body []byte, opts SendMessageOptions) error {
	url := q.buildURL(fmt.Sprintf("/topic/%s", q.config.Topic))

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return err
	}

	contentType := opts.ContentType
	if contentType == "" {
		contentType = "application/json"
	}
	req.Header.Set("Content-Type", contentType)
	q.authHeader(req)

	resp, err := q.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close() //nolint:errcheck

	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to send message: %s (status: %d)", string(respBody), resp.StatusCode)
	}

	return nil
}

// Enqueue sends a job to the Vercel Queue as a JSON-encoded message.
func (q *VercelQueue) Enqueue(ctx context.Context, job domain.ExtractionJob) error {
	data, err := json.Marshal(job)
	if err != nil {
		return err
	}

	return q.SendMessage(ctx, data, SendMessageOptions{ContentType: "application/json"})
}

// EnqueueConcert sends a concert to the Vercel Queue as a JSON-encoded message.
func (q *VercelQueue) EnqueueConcert(ctx context.Context, concert domain.Concert) error {
	data, err := json.Marshal(concert)
	if err != nil {
		return err
	}

	return q.SendMessage(ctx, data, SendMessageOptions{ContentType: "application/json"})
}


// ReceiveMessages retrieves messages from the configured topic/consumer.
func (q *VercelQueue) ReceiveMessages(ctx context.Context, opts ReceiveMessagesOptions) (*ReceiveMessagesResponse, error) {
	url := q.buildURL(fmt.Sprintf("/topic/%s/consumer/%s", q.config.Topic, q.config.Consumer))

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	q.authHeader(req)

	resp, err := q.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close() //nolint:errcheck

	if resp.StatusCode == http.StatusNoContent {
		return &ReceiveMessagesResponse{}, nil
	}

	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to receive m
essages: %s (status: %d)", string(respBody), resp.StatusCode)
	}

	var result ReceiveMessagesResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return &result, nil
}

// AcknowledgeMessage deletes a message from the queue by its receipt handle.
func (q *VercelQueue) AcknowledgeMessage(ctx context.Context, receiptHandle string) error {
	url := q.buildURL(fmt.Sprintf("/topic/%s/consumer/%s/message/%s", q.config.Topic, q.config.Consumer, receiptHandle))

	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, url, nil)
	if err != nil {
		return err
	}
	q.authHeader(req)

	resp, err := q.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close() //nolint:errcheck

	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to acknowledge message: %s (status: %d)", string(respBody), resp.StatusCode)
	}

	return nil
}

// ExtendLease extends the visibility timeout of a message currently being processed.
func (q *VercelQueue) ExtendLease(ctx context.Context, receiptHandle string, visibilityTimeoutSeconds int) error {
	url := q.buildURL(fmt.Sprintf("/topic/%s/consumer/%s/message/%s", q.config.Topic, q.config.Consumer, receiptHandle))

	body, err := json.Marshal(map[string]int{"visibilityTimeoutSeconds": visibilityTimeoutSeconds})
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPatch, url, bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	q.authHeader(req)

	resp, err := q.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close() //nolint:errcheck

	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to extend lease: %s (status: %d)", string(respBody), resp.StatusCode)
	}

	return nil
}

// Process is not needed for push-based consumers.
// Vercel automatically delivers messages to /api/queue.

func (q *VercelQueue) Process(ctx context.Context, handler domain.QueueHandler) error {
	return fmt.Errorf("push-based consumers do not require manual processing")
}
