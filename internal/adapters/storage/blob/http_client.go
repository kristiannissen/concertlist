// Package blob provides HTTP client implementation for Vercel Blob storage.
package blob

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/kristiannissen/concertlist/internal/domain"
)

// HTTPClient implements domain.StoragePort using Vercel Blob HTTP API.
type HTTPClient struct {
	storeID       string
	accessToken   string
	apiVersion    string
	baseURL       string
	storageBaseURL string
}

// HTTPClientConfig holds configuration for the HTTP client.
type HTTPClientConfig struct {
	StoreID       string
	AccessToken   string
	APIVersion    string
	BaseURL       string
	StorageBaseURL string
}

// NewHTTPClient creates a new HTTPClient for Vercel Blob storage.
func NewHTTPClient(config HTTPClientConfig) *HTTPClient {
	return &HTTPClient{
		storeID:       config.StoreID,
		accessToken:   config.AccessToken,
		apiVersion:    config.APIVersion,
		baseURL:       config.BaseURL,
		storageBaseURL: config.StorageBaseURL,
	}
}

// Save stores concerts in Vercel Blob storage.
func (c *HTTPClient) Save(ctx context.Context, concerts []domain.Concert) error {
	const defaultPathname = "concerts.json"
	
	// Marshal concerts to JSON
	data, err := json.MarshalIndent(concerts, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal concerts: %w", err)
	}

	// Create request body
	body := bytes.NewReader(data)

	// Create request URL with pathname query parameter
	url := fmt.Sprintf("%s/?pathname=%s", c.baseURL, defaultPathname)

	// Create HTTP request
	req, err := http.NewRequestWithContext(ctx, http.MethodPut, url, body)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	// Set required headers
	req.Header.Set("Authorization", "Bearer "+c.accessToken)
	req.Header.Set("x-vercel-blob-store-id", c.storeID)
	req.Header.Set("x-api-version", c.apiVersion)
	req.Header.Set("x-vercel-blob-access", string(AccessPublic))
	req.Header.Set("x-content-type", "application/json")
	req.Header.Set("x-allow-overwrite", "1")

	// Create HTTP client and execute request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	// Check response status
	if resp.StatusCode >= http.StatusBadRequest {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("request failed with status %d: %s", resp.StatusCode, string(body))
	}

	return nil
}

// Load retrieves concerts from Vercel Blob storage.
func (c *HTTPClient) Load(ctx context.Context) ([]domain.Concert, error) {
	const defaultPathname = "concerts.json"
	
	// Create URL for the stored concerts file
	url := fmt.Sprintf("%s/%s", c.storageBaseURL, defaultPathname)
	
	// Create HTTP request with context
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	
	// Set Authorization header for private blobs
	req.Header.Set("Authorization", "Bearer "+c.accessToken)
	
	// Create HTTP client and execute request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close() //nolint:errcheck

	// Check if file doesn't exist yet (empty store). This must be checked
	// before the generic error check below, since 404 also satisfies
	// resp.StatusCode >= http.StatusBadRequest.
	if resp.StatusCode == http.StatusNotFound {
		return []domain.Concert{}, nil
	}

	// Check response status
	if resp.StatusCode >= http.StatusBadRequest {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("request failed with status %d: %s", resp.StatusCode, string(body))
	}

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// Unmarshal JSON into concerts
	var concerts []domain.Concert
	if err := json.Unmarshal(body, &concerts); err != nil {
		return nil, fmt.Errorf("failed to unmarshal concerts: %w", err)
	}

	return concerts, nil
}

// uploadBlob uploads data to Vercel Blob.
func (c *HTTPClient) uploadBlob(ctx context.Context, pathname string, data io.Reader, contentType string, access string) error {
	// Create request URL with pathname query parameter
	url := fmt.Sprintf("%s/?pathname=%s", c.baseURL, pathname)
	
	// Create HTTP request with context and data body
	req, err := http.NewRequestWithContext(ctx, http.MethodPut, url, data)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	// Set required headers
	req.Header.Set("Authorization", "Bearer "+c.accessToken)
	req.Header.Set("x-vercel-blob-store-id", c.storeID)
	req.Header.Set("x-api-version", c.apiVersion)
	req.Header.Set("x-vercel-blob-access", access)
	req.Header.Set("x-content-type", contentType)
	req.Header.Set("x-allow-overwrite", "1")
	
	// Create HTTP client and execute request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close() //nolint:errcheck

	// Check response status
	if resp.StatusCode >= http.StatusBadRequest {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("request failed with status %d: %s", resp.StatusCode, string(body))
	}

	return nil
}

// downloadBlob downloads data from Vercel Blob.
func (c *HTTPClient) downloadBlob(ctx context.Context, pathname string) (io.ReadCloser, error) {
	// Create URL for the blob file
	url := fmt.Sprintf("%s/%s", c.storageBaseURL, pathname)
	
	// Create HTTP request with context
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	
	// Set Authorization header for private blobs
	req.Header.Set("Authorization", "Bearer "+c.accessToken)
	
	// Create HTTP client and execute request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	
	// Check response status
	if resp.StatusCode >= http.StatusBadRequest {
		body, _ := io.ReadAll(resp.Body)
		defer resp.Body.Close() //nolint:errcheck
		return nil, fmt.Errorf("request failed with status %d: %s", resp.StatusCode, string(body))
	}

	// Return the response body for the caller to read and close
	return resp.Body, nil
}

// listBlobs lists blobs in the store.
func (c *HTTPClient) listBlobs(ctx context.Context, prefix string) ([]domain.Concert, error) {
	// TODO: Implement using GET /?prefix={prefix}
	return nil, nil
}

// deleteBlob deletes a blob from the store.
func (c *HTTPClient) deleteBlob(ctx context.Context, pathname string) error {
	// TODO: Implement using DELETE /?pathname={pathname}
	return nil
}

// deleteBlobs deletes multiple blobs from the store.
func (c *HTTPClient) deleteBlobs(ctx context.Context, urls []string) error {
	// TODO: Implement using POST /delete
	return nil
}

// getBlobMetadata retrieves metadata for a blob.
func (c *HTTPClient) getBlobMetadata(ctx context.Context, url string) (map[string]string, error) {
	// TODO: Implement using GET /?url={url}
	return nil, nil
}

// copyBlob copies a blob to a new pathname.
func (c *HTTPClient) copyBlob(ctx context.Context, toPathname string, fromURL string) error {
	// TODO: Implement using PUT /?pathname={to_pathname}&fromUrl={from_url}
	return nil
}