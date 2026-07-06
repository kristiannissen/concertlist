// Package queue provides Vercel Queues adapter using the HTTP API.
package queue

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/kristiannissen/concertlist/pkg/domain"
)

func TestNewVercelQueue(t *testing.T) {
	t.Parallel()

	config := QueueConfig{
		Region:    "fra1",
		Topic:     "test-topic",
		Consumer:  "test-consumer",
		OIDCToken: "test-token",
	}

	q, err := NewVercelQueue(config)
	if err != nil {
		t.Fatalf("Failed to create VercelQueue: %v", err)
	}

	if q.config.Topic != "test-topic" {
		t.Errorf("Expected Topic to be test-topic, got %s", q.config.Topic)
	}
}

func TestNewVercelQueue_EmptyTopic(t *testing.T) {
	t.Parallel()

	config := QueueConfig{
		Region: "fra1",
	}

	_, err := NewVercelQueue(config)
	if err == nil {
		t.Error("Expected error for empty topic, got nil")
	}
}

func TestNewVercelQueue_DefaultValues(t *testing.T) {
	t.Parallel()

	config := QueueConfig{
		Topic:     "test-topic",
		OIDCToken: "test-token",
	}

	q, err := NewVercelQueue(config)
	if err != nil {
		t.Fatalf("Failed to create VercelQueue: %v", err)
	}

	if q.config.Region != DefaultRegion {
		t.Errorf("Expected Region to be %s, got %s", DefaultRegion, q.config.Region)
	}
	if q.config.Consumer != "default" {
		t.Errorf("Expected Consumer to be default, got %s", q.config.Consumer)
	}
}

func TestBuildURL(t *testing.T) {
	t.Parallel()

	config := QueueConfig{
		Region: "fra1",
		Topic:  "test-topic",
	}
	q, _ := NewVercelQueue(config)

	url := q.buildURL("/topic/test")
	expected := "https://fra1.vercel-queue.com/api/v3/topic/test"
	if url != expected {
		t.Errorf("Expected URL to be %s, got %s", expected, url)
	}
}

func TestSendMessage(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("Expected POST method, got %s", r.Method)
		}
		if r.Header.Get("Authorization") != "Bearer test-token" {
			t.Errorf("Expected Authorization header")
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	config := QueueConfig{
		Region:    "test-region",
		Topic:     "test-topic",
		Consumer:  "test-consumer",
		OIDCToken: "test-token",
	}
	q := &VercelQueue{
		config: config,
		client: server.Client(),
	}
	q.buildURL = func(path string) string {
		return server.URL + path
	}

	ctx := context.Background()
	messageBody := []byte("test message")
	opts := SendMessageOptions{
		ContentType:     "application/json",
		RetentionSeconds: 86400,
	}

	err := q.SendMessage(ctx, messageBody, opts)
	if err != nil {
		t.Fatalf("Failed to send message: %v", err)
	}
}

func TestEnqueue(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	config := QueueConfig{
		Region:    "test-region",
		Topic:     "test-topic",
		Consumer:  "test-consumer",
		OIDCToken: "test-token",
	}
	q := &VercelQueue{
		config: config,
		client: server.Client(),
	}
	q.buildURL = func(path string) string {
		return server.URL + path
	}

	ctx := context.Background()
	job := domain.ExtractionJob{
		Venue: "test-venue",
	}

	err := q.Enqueue(ctx, job)
	if err != nil {
		t.Fatalf("Failed to enqueue job: %v", err)
	}
}

func TestReceiveMessages_NoContent(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	config := QueueConfig{
		Region:    "test-region",
		Topic:     "test-topic",
		Consumer:  "test-consumer",
		OIDCToken: "test-token",
	}
	q := &VercelQueue{
		config: config,
		client: server.Client(),
	}
	q.buildURL = func(path string) string {
		return server.URL + path
	}

	ctx := context.Background()
	opts := ReceiveMessagesOptions{}

	resp, err := q.ReceiveMessages(ctx, opts)
	if err != nil {
		t.Fatalf("Failed to receive messages: %v", err)
	}

	if len(resp.Messages) != 0 {
		t.Errorf("Expected 0 messages, got %d", len(resp.Messages))
	}
}

func TestAcknowledgeMessage(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("Expected DELETE method, got %s", r.Method)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	config := QueueConfig{
		Region:    "test-region",
		Topic:     "test-topic",
		Consumer:  "test-consumer",
		OIDCToken: "test-token",
	}
	q := &VercelQueue{
		config: config,
		client: server.Client(),
	}
	q.buildURL = func(path string) string {
		return server.URL + path
	}

	ctx := context.Background()
	err := q.AcknowledgeMessage(ctx, "test-receipt")
	if err != nil {
		t.Fatalf("Failed to acknowledge message: %v", err)
	}
}

func TestExtendLease(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPatch {
			t.Errorf("Expected PATCH method, got %s", r.Method)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	config := QueueConfig{
		Region:    "test-region",
		Topic:     "test-topic",
		Consumer:  "test-consumer",
		OIDCToken: "test-token",
	}
	q := &VercelQueue{
		config: config,
		client: server.Client(),
	}
	q.buildURL = func(path string) string {
		return server.URL + path
	}

	ctx := context.Background()
	err := q.ExtendLease(ctx, "test-receipt", 120)
	if err != nil {
		t.Fatalf("Failed to extend lease: %v", err)
	}
}

func TestEnqueueConcert(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("Expected POST method, got %s", r.Method)
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		if r.Header.Get("Authorization") != "Bearer test-token" {
			t.Errorf("Expected Authorization header")
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		if r.Header.Get("Content-Type") != "application/json" {
			t.Errorf("Expected Content-Type to be application/json, got %s", r.Header.Get("Content-Type"))
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	config := QueueConfig{
		Region:    "test-region",
		Topic:     "test-topic",
		Consumer:  "test-consumer",
		OIDCToken: "test-token",
	}
	q := &VercelQueue{
		config: config,
		client: server.Client(),
	}
	q.buildURL = func(path string) string {
		return server.URL + path
	}

	ctx := context.Background()
	concert := domain.Concert{
		ID:    "test-id",
		Title: "Test Concert",
		Venue: "Test Venue",
		Date:  "2026-07-04",
	}

	err := q.EnqueueConcert(ctx, concert)
	if err != nil {
		t.Fatalf("EnqueueConcert failed: %v", err)
	}
}

func TestEnqueueConcert_Error(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal Server Error"))
	}))
	defer server.Close()

	config := QueueConfig{
		Region:    "test-region",
		Topic:     "test-topic",
		Consumer:  "test-consumer",
		OIDCToken: "test-token",
	}
	q := &VercelQueue{
		config: config,
		client: server.Client(),
	}
	q.buildURL = func(path string) string {
		return server.URL + path
	}

	ctx := context.Background()
	concert := domain.Concert{
		ID:    "test-id",
		Title: "Test Concert",
		Venue: "Test Venue",
		Date:  "2026-07-04",
	}

	err := q.EnqueueConcert(ctx, concert)
	if err == nil {
		t.Error("Expected EnqueueConcert to return error")
	}
}
