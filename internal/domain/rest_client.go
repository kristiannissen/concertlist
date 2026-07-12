// This file is in internal/domain
package domain

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"time"
)

// RESTClient is an async HTTP client for making non-blocking requests.
// It implements a driven port for external API calls.
type RESTClient struct {
	client      *http.Client
	baseURL     string
	timeout     time.Duration
	retryPolicy RetryPolicy
}

// RetryPolicy defines retry behavior for failed requests.
type RetryPolicy struct {
	MaxRetries  int
	RetryDelay  time.Duration
	RetryCodes  []int // HTTP status codes to retry (e.g., 429, 500)
	ShouldRetry func(resp *http.Response, err error) bool
}

// NewRESTClient initializes a new async REST client.
func NewRESTClient(baseURL string, opts ...func(*RESTClient)) *RESTClient {
	client := &RESTClient{
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
		baseURL: baseURL,
		timeout: 10 * time.Second,
		retryPolicy: RetryPolicy{
			MaxRetries: 3,
			RetryDelay: 1 * time.Second,
			RetryCodes: []int{429, 500, 502, 503, 504},
			ShouldRetry: func(resp *http.Response, err error) bool {
				if err != nil {
					return true // Retry on network errors
				}
				for _, code := range []int{429, 500, 502, 503, 504} {
					if resp.StatusCode == code {
						return true
					}
				}
				return false
			},
		},
	}

	for _, opt := range opts {
		opt(client)
	}

	return client
}

// WithHTTPClient allows injecting a custom *http.Client.
func WithHTTPClient(client *http.Client) func(*RESTClient) {
	return func(c *RESTClient) {
		c.client = client
	}
}

// WithTimeout sets the request timeout.
func WithTimeout(timeout time.Duration) func(*RESTClient) {
	return func(c *RESTClient) {
		c.timeout = timeout
	}
}

// WithRetryPolicy sets the retry policy.
func WithRetryPolicy(policy RetryPolicy) func(*RESTClient) {
	return func(c *RESTClient) {
		c.retryPolicy = policy
	}
}

// Request represents an HTTP request to be made asynchronously.
type Request struct {
	Method  string
	Path    string
	Headers map[string]string
	Body    []byte
}

// Response represents the HTTP response from an async request.
type Response struct {
	StatusCode int
	Body       []byte
	Headers    http.Header
	Error      error
}

// DoAsync sends an HTTP request asynchronously and returns a channel to receive the response.
// The caller must read from the channel to avoid goroutine leaks.
func (c *RESTClient) DoAsync(ctx context.Context, req Request) <-chan Response {
	responseChan := make(chan Response, 1)

	go func() {
		defer close(responseChan)

		// Apply context timeout
		ctx, cancel := context.WithTimeout(ctx, c.timeout)
		defer cancel()

		// Retry loop
		var lastErr error
		var lastResp *http.Response
		for attempt := 0; attempt <= c.retryPolicy.MaxRetries; attempt++ {
			select {
			case <-ctx.Done():
				responseChan <- Response{Error: ctx.Err()}
				return
			default:
				// Build the request
				url := c.baseURL + req.Path
				httpReq, err := http.NewRequestWithContext(ctx, req.Method, url, bytes.NewReader(req.Body))
				if err != nil {
					lastErr = fmt.Errorf("failed to create request: %w", err)
					if !c.retryPolicy.ShouldRetry(nil, lastErr) {
						break
					}
					time.Sleep(c.retryPolicy.RetryDelay)
					continue
				}

				// Set headers
				for key, value := range req.Headers {
					httpReq.Header.Set(key, value)
				}

				// Send the request
				resp, err := c.client.Do(httpReq)
				if err != nil {
					lastErr = fmt.Errorf("request failed: %w", err)
					if !c.retryPolicy.ShouldRetry(nil, lastErr) {
						break
					}
					time.Sleep(c.retryPolicy.RetryDelay)
					continue
				}
				defer resp.Body.Close() //nolint:errcheck

				// Read the response body
				body, err := io.ReadAll(resp.Body)
				if err != nil {
					lastErr = fmt.Errorf("failed to read response body: %w", err)
					if !c.retryPolicy.ShouldRetry(resp, lastErr) {
						break
					}
					time.Sleep(c.retryPolicy.RetryDelay)
					continue
				}

				// Check if we should retry
				if c.retryPolicy.ShouldRetry(resp, nil) {
					lastResp = resp
					time.Sleep(c.retryPolicy.RetryDelay)
					continue
				}

				// Success
				responseChan <- Response{
					StatusCode: resp.StatusCode,
					Body:       body,
					Headers:    resp.Header,
					Error:      nil,
				}
				return
			}
		}

		// All retries exhausted
		if lastResp != nil {
			responseChan <- Response{
				StatusCode: lastResp.StatusCode,
				Error:      fmt.Errorf("request failed after %d retries", c.retryPolicy.MaxRetries),
			}
		} else {
			responseChan <- Response{Error: lastErr}
		}
	}()

	return responseChan
}

// GetAsync sends a GET request asynchronously.
func (c *RESTClient) GetAsync(ctx context.Context, path string, headers map[string]string) <-chan Response {
	return c.DoAsync(ctx, Request{
		Method:  "GET",
		Path:    path,
		Headers: headers,
	})
}

// PostAsync sends a POST request asynchronously.
func (c *RESTClient) PostAsync(ctx context.Context, path string, body []byte, headers map[string]string) <-chan Response {
	return c.DoAsync(ctx, Request{
		Method:  "POST",
		Path:    path,
		Body:    body,
		Headers: headers,
	})
}
