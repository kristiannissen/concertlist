// Package queue provides Vercel Queues adapter using the HTTP API.
package queue

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"strings"

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
// QUEUE_NAME (required, used as Topic), QUEUE_REGION, QUEUE_CONSUMER,
// VERCEL_OIDC_TOKEN, and VERCEL_DEPLOYMENT_ID (all injected automatically by Vercel at runtime).
func NewVercelQueueFromEnv() (*VercelQueue, error) {
	topic := os.Getenv("QUEUE_NAME")
	if topic == "" {
		return nil, fmt.Errorf("QUEUE_NAME environment variable is required")
	}

	return NewVercelQueue(QueueConfig{
		Topic:        topic,
		Region:       os.Getenv("QUEUE_REGION"),
		Consumer:     os.Getenv("QUEUE_CONSUMER"),
		OIDCToken:    os.Getenv("VERCEL_OIDC_TOKEN"),
		DeploymentID: os.Getenv("VERCEL_DEPLOYMENT_ID"),
	})
}

// authHeader sets the Authorization header on the request when an OIDC token is configured.
func (q *VercelQueue) authHeader(req *http.Request) {
	if q.config.OIDCToken != "" {
		req.Header.Set("Authorization", "Bearer "+q.config.OIDCToken)
	}
}

// deploymentHeader sets the Vercel Deployment ID header when configured.
func (q *VercelQueue) deploymentHeader(req *http.Request) {
	if q.config.DeploymentID != "" {
		req.Header.Set("Vqs-Deployment-Id", q.config.DeploymentID)
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

	if opts.RetentionSeconds != 0 {
		req.Header.Set("Vqs-Retention-Seconds", fmt.Sprintf("%d", opts.RetentionSeconds))
	}
	if opts.DelaySeconds != 0 {
		req.Header.Set("Vqs-Delay-Seconds", fmt.Sprintf("%d", opts.DelaySeconds))
	}
	if opts.IdempotencyKey != "" {
		req.Header.Set("Vqs-Idempotency-Key", opts.IdempotencyKey)
	}

	q.authHeader(req)
	q.deploymentHeader(req)

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

// parseReceiveResponse parses the response body based on the Content-Type header.
// Vercel Queues returns either multipart/mixed or application/x-ndjson.
func parseReceiveResponse(resp *http.Response) (*ReceiveMessagesResponse, error) {
	contentType := resp.Header.Get("Content-Type")
	
	// Handle multipart/mixed response
	if strings.HasPrefix(contentType, "multipart/mixed") {
		return parseMultipartResponse(resp.Body)
	}
	
	// Handle application/x-ndjson response (newline-delimited JSON)
	if strings.HasPrefix(contentType, "application/x-ndjson") {
		return parseNDJSONResponse(resp.Body)
	}
	
	// Fallback: try to parse as JSON (for testing with mock servers)
	var result ReceiveMessagesResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}
	return &result, nil
}

// parseMultipartResponse parses a multipart/mixed response body.
func parseMultipartResponse(body io.Reader) (*ReceiveMessagesResponse, error) {
	var result ReceiveMessagesResponse
	
	// Parse the multipart response
	mr := multipart.NewReader(body, "")
	for {
		part, err := mr.NextPart()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, fmt.Errorf("failed to read multipart part: %w", err)
		}

		// Read the part content
		partData, err := io.ReadAll(part)
		if err != nil {
			return nil, fmt.Errorf("failed to read part data: %w", err)
		}

		// Parse the part as a Message
		var msg Message
		if err := json.Unmarshal(partData, &msg); err != nil {
			return nil, fmt.Errorf("failed to unmarshal message: %w", err)
		}
		
		// Extract receipt handle from headers if available
		if receiptHandle := part.Header.Get("Vqs-Receipt-Handle"); receiptHandle != "" {
			msg.ReceiptHandle = receiptHandle
		}
		
		result.Messages = append(result.Messages, msg)
		result.ReceiptHandles = append(result.ReceiptHandles, msg.ReceiptHandle)
	}
	
	return &result, nil
}

// parseNDJSONResponse parses a newline-delimited JSON response body.
func parseNDJSONResponse(body io.Reader) (*ReceiveMessagesResponse, error) {
	var result ReceiveMessagesResponse
	
	// Read the entire body
	data, err := io.ReadAll(body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}
	
	// Split by newlines
	lines := strings.Split(string(data), "
")
	
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		
		// Parse as JSON Message
		var msg Message
		if err := json.Unmarshal([]byte(line), &msg); err != nil {
			// Try to parse as a simple string message
			msg.Body = []byte(line)
			msg.ContentType = "text/plain"
		}
		
		result.Messages = append(result.Messages, msg)
		if msg.ReceiptHandle != "" {
			result.ReceiptHandles = append(result.ReceiptHandles, msg.ReceiptHandle)
		}
	}
	
	return &result, nil
}

// internalReceiveMessages is the internal implementation that returns the raw response.
func (q *VercelQueue) internalReceiveMessages(ctx context.Context, opts ReceiveMessagesOptions) (*ReceiveMessagesResponse, error) {
	url := q.buildURL(fmt.Sprintf("/topic/%s/consumer/%s", q.config.Topic, q.config.Consumer))

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, nil)
	if err != nil {
		return nil, err
	}

	// Set Accept header for response format
	accept := opts.Accept
	if accept == "" {
		accept = "application/x-ndjson"
	}
	req.Header.Set("Accept", accept)

	// Set optional headers
	if opts.MaxMessages != 0 {
		req.Header.Set("Vqs-Max-Messages", fmt.Sprintf("%d", opts.MaxMessages))
	}
	if opts.VisibilityTimeoutSeconds != 0 {
		req.Header.Set("Vqs-Visibility-Timeout-Seconds", fmt.Sprintf("%d", opts.VisibilityTimeoutSeconds))
	}
	if opts.MaxConcurrency != 0 {
		req.Header.Set("Vqs-Max-Concurrency", fmt.Sprintf("%d", opts.MaxConcurrency))
	}

	q.authHeader(req)
	q.deploymentHeader(req)

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
		return nil, fmt.Errorf("failed to receive messages: %s (status: %d)", string(respBody), resp.StatusCode)
	}

	// Parse response based on Content-Type
	return parseReceiveResponse(resp)
}

// ReceiveMessages retrieves messages from the configured topic/consumer and returns them as Concerts.
// Implements domain.QueuePort interface.
func (q *VercelQueue) ReceiveMessages(ctx context.Context) ([]domain.Concert, error) {
	resp, err := q.internalReceiveMessages(ctx, ReceiveMessagesOptions{
		Accept:                "application/x-ndjson",
		MaxMessages:          10,
		VisibilityTimeoutSeconds: 60,
	})
	if err != nil {
		return nil, err
	}

	var concerts []domain.Concert
	for _, msg := range resp.Messages {
		var concert domain.Concert
		if err := json.Unmarshal(msg.Body, &concert); err != nil {
			// If it's not a Concert, skip it
			continue
		}
		concerts = append(concerts, concert)
	}

	return concerts, nil
}

// internalReceiveMessageByID is the internal implementation that returns the raw response.
func (q *VercelQueue) internalReceiveMessageByID(ctx context.Context, messageID string, opts ReceiveMessagesOptions) (*ReceiveMessagesResponse, error) {
	// URL-encode the message ID as required by the API
	encodedMessageID := url.PathEscape(messageID)
	url := q.buildURL(fmt.Sprintf("/topic/%s/consumer/%s/id/%s", q.config.Topic, q.config.Consumer, encodedMessageID))

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, nil)
	if err != nil {
		return nil, err
	}

	// Set Accept header for response format
	accept := opts.Accept
	if accept == "" {
		accept = "application/x-ndjson"
	}
	req.Header.Set("Accept", accept)

	// Set optional headers
	if opts.VisibilityTimeoutSeconds != 0 {
		req.Header.Set("Vqs-Visibility-Timeout-Seconds", fmt.Sprintf("%d", opts.VisibilityTimeoutSeconds))
	}
	if opts.MaxConcurrency != 0 {
		req.Header.Set("Vqs-Max-Concurrency", fmt.Sprintf("%d", opts.MaxConcurrency))
	}

	q.authHeader(req)
	q.deploymentHeader(req)

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
		return nil, fmt.Errorf("failed to receive message by ID: %s (status: %d)", string(respBody), resp.StatusCode)
	}

	// Parse response based on Content-Type
	return parseReceiveResponse(resp)
}

// ReceiveMessageByID retrieves a specific message by its ID and returns it as a Concert.
// Implements domain.QueuePort interface.
func (q *VercelQueue) ReceiveMessageByID(ctx context.Context, messageID string) (domain.Concert, error) {
	resp, err := q.internalReceiveMessageByID(ctx, messageID, ReceiveMessagesOptions{
		Accept:                "application/x-ndjson",
		VisibilityTimeoutSeconds: 60,
	})
	if err != nil {
		return domain.Concert{}, err
	}

	if len(resp.Messages) == 0 {
		return domain.Concert{}, fmt.Errorf("no message found with ID %s", messageID)
	}

	var concert domain.Concert
	if err := json.Unmarshal(resp.Messages[0].Body, &concert); err != nil {
		return domain.Concert{}, fmt.Errorf("failed to parse message as Concert: %w", err)
	}

	return concert, nil
}

// AcknowledgeMessage deletes a message from the queue by its receipt handle.
// Uses the correct /lease/ endpoint path as per Vercel Queues HTTP API specification.
func (q *VercelQueue) AcknowledgeMessage(ctx context.Context, receiptHandle string) error {
	// URL-encode the receipt handle as required by the API
	encodedReceiptHandle := url.PathEscape(receiptHandle)
	url := q.buildURL(fmt.Sprintf("/topic/%s/consumer/%s/lease/%s", q.config.Topic, q.config.Consumer, encodedReceiptHandle))

	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, url, nil)
	if err != nil {
		return err
	}

	q.authHeader(req)
	q.deploymentHeader(req)

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
// Uses the correct /lease/ endpoint path as per Vercel Queues HTTP API specification.
func (q *VercelQueue) ExtendLease(ctx context.Context, receiptHandle string, visibilityTimeoutSeconds int) error {
	// URL-encode the receipt handle as required by the API
	encodedReceiptHandle := url.PathEscape(receiptHandle)
	url := q.buildURL(fmt.Sprintf("/topic/%s/consumer/%s/lease/%s", q.config.Topic, q.config.Consumer, encodedReceiptHandle))

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
	q.deploymentHeader(req)

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
