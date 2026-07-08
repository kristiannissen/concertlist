// Package queue provides types for Vercel Queues API.
package queue

import (
	"testing"
)

func TestMessage(t *testing.T) {
	t.Parallel()

	msg := Message{
		ID:            "msg-123",
		Body:          []byte("test message"),
		ContentType:   "application/json",
		ReceiptHandle: "receipt-123",
		Attributes:    map[string]string{"key": "value"},
	}

	if msg.ID != "msg-123" {
		t.Errorf("Expected ID to be 'msg-123', got %s", msg.ID)
	}
	if string(msg.Body) != "test message" {
		t.Errorf("Expected Body to be 'test message', got %s", string(msg.Body))
	}
	if msg.ContentType != "application/json" {
		t.Errorf("Expected ContentType to be 'application/json', got %s", msg.ContentType)
	}
	if msg.ReceiptHandle != "receipt-123" {
		t.Errorf("Expected ReceiptHandle to be 'receipt-123', got %s", msg.ReceiptHandle)
	}
	if msg.Attributes["key"] != "value" {
		t.Errorf("Expected Attributes['key'] to be 'value', got %s", msg.Attributes["key"])
	}
}

func TestSendMessageOptions(t *testing.T) {
	t.Parallel()

	opts := SendMessageOptions{
		ContentType:      "application/json",
		RetentionSeconds: 86400,
		DelaySeconds:     60,
		IdempotencyKey:   "unique-key-123",
	}

	if opts.ContentType != "application/json" {
		t.Errorf("Expected ContentType to be 'application/json', got %s", opts.ContentType)
	}
	if opts.RetentionSeconds != 86400 {
		t.Errorf("Expected RetentionSeconds to be 86400, got %d", opts.RetentionSeconds)
	}
	if opts.DelaySeconds != 60 {
		t.Errorf("Expected DelaySeconds to be 60, got %d", opts.DelaySeconds)
	}
	if opts.IdempotencyKey != "unique-key-123" {
		t.Errorf("Expected IdempotencyKey to be 'unique-key-123', got %s", opts.IdempotencyKey)
	}
}

func TestReceiveMessagesOptions(t *testing.T) {
	t.Parallel()

	opts := ReceiveMessagesOptions{
		MaxMessages:              5,
		VisibilityTimeoutSeconds: 120,
		MaxConcurrency:           10,
		Accept:                   "application/x-ndjson",
	}

	if opts.MaxMessages != 5 {
		t.Errorf("Expected MaxMessages to be 5, got %d", opts.MaxMessages)
	}
	if opts.VisibilityTimeoutSeconds != 120 {
		t.Errorf("Expected VisibilityTimeoutSeconds to be 120, got %d", opts.VisibilityTimeoutSeconds)
	}
	if opts.MaxConcurrency != 10 {
		t.Errorf("Expected MaxConcurrency to be 10, got %d", opts.MaxConcurrency)
	}
	if opts.Accept != "application/x-ndjson" {
		t.Errorf("Expected Accept to be 'application/x-ndjson', got %s", opts.Accept)
	}
}

func TestReceiveMessagesResponse(t *testing.T) {
	t.Parallel()

	msg1 := Message{ID: "msg-1", Body: []byte("message 1")}
	msg2 := Message{ID: "msg-2", Body: []byte("message 2")}

	resp := ReceiveMessagesResponse{
		Messages:       []Message{msg1, msg2},
		ReceiptHandles: []string{"receipt-1", "receipt-2"},
	}

	if len(resp.Messages) != 2 {
		t.Errorf("Expected 2 messages, got %d", len(resp.Messages))
	}
	if len(resp.ReceiptHandles) != 2 {
		t.Errorf("Expected 2 receipt handles, got %d", len(resp.ReceiptHandles))
	}
}

func TestAcknowledgeMessageOptions(t *testing.T) {
	t.Parallel()

	opts := AcknowledgeMessageOptions{
		ReceiptHandle: "receipt-123",
	}

	if opts.ReceiptHandle != "receipt-123" {
		t.Errorf("Expected ReceiptHandle to be 'receipt-123', got %s", opts.ReceiptHandle)
	}
}

func TestExtendLeaseOptions(t *testing.T) {
	t.Parallel()

	opts := ExtendLeaseOptions{
		ReceiptHandle:            "receipt-123",
		VisibilityTimeoutSeconds: 120,
	}

	if opts.ReceiptHandle != "receipt-123" {
		t.Errorf("Expected ReceiptHandle to be 'receipt-123', got %s", opts.ReceiptHandle)
	}
	if opts.VisibilityTimeoutSeconds != 120 {
		t.Errorf("Expected VisibilityTimeoutSeconds to be 120, got %d", opts.VisibilityTimeoutSeconds)
	}
}

func TestQueueConfig(t *testing.T) {
	t.Parallel()

	config := QueueConfig{
		Region:       "fra1",
		Topic:        "my-topic",
		Consumer:     "my-consumer",
		OIDCToken:    "token-123",
		DeploymentID: "deployment-123",
	}

	if config.Region != "fra1" {
		t.Errorf("Expected Region to be 'fra1', got %s", config.Region)
	}
	if config.Topic != "my-topic" {
		t.Errorf("Expected Topic to be 'my-topic', got %s", config.Topic)
	}
	if config.Consumer != "my-consumer" {
		t.Errorf("Expected Consumer to be 'my-consumer', got %s", config.Consumer)
	}
	if config.OIDCToken != "token-123" {
		t.Errorf("Expected OIDCToken to be 'token-123', got %s", config.OIDCToken)
	}
	if config.DeploymentID != "deployment-123" {
		t.Errorf("Expected DeploymentID to be 'deployment-123', got %s", config.DeploymentID)
	}
}

func TestErrorResponse(t *testing.T) {
	t.Parallel()

	errResp := ErrorResponse{
		Error:   "NotFound",
		Message: "Topic not found",
	}

	if errResp.Error != "NotFound" {
		t.Errorf("Expected Error to be 'NotFound', got %s", errResp.Error)
	}
	if errResp.Message != "Topic not found" {
		t.Errorf("Expected Message to be 'Topic not found', got %s", errResp.Message)
	}
}
