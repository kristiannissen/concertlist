// Package queue provides types for Vercel Queues API.
package queue

// Message represents a message received from Vercel Queues.
type Message struct {
	ID            string            "json:"id""
	Body          []byte            "json:"body""
	ContentType   string            "json:"contentType""
	ReceiptHandle string            "json:"receiptHandle""
	Attributes    map[string]string "json:"attributes,omitempty""
}

// SendMessageOptions contains options for sending a message.
type SendMessageOptions struct {
	ContentType     string
	RetentionSeconds int
	DelaySeconds     int
	IdempotencyKey   string
}

// ReceiveMessagesOptions contains options for receiving messages.
type ReceiveMessagesOptions struct {
	MaxMessages            int
	VisibilityTimeoutSeconds int
	MaxConcurrency         int
	Accept                 string
}

// ReceiveMessagesResponse contains the response from receiving messages.
type ReceiveMessagesResponse struct {
	Messages      []Message
	ReceiptHandles []string
}

// AcknowledgeMessageOptions contains options for acknowledging a message.
type AcknowledgeMessageOptions struct {
	ReceiptHandle string
}

// ExtendLeaseOptions contains options for extending a message lease.
type ExtendLeaseOptions struct {
	ReceiptHandle            string
	VisibilityTimeoutSeconds int
}

// QueueConfig contains configuration for the Vercel Queue client.
type QueueConfig struct {
	Region       string
	Topic        string
	Consumer     string
	OIDCToken    string
	DeploymentID string
}

// ErrorResponse contains error information from the API.
type ErrorResponse struct {
	Error   string
	Message string
}
