// Package blob provides HTTP client implementation for Vercel Blob storage.
package blob

import (
	"context"
	"io"

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
	// TODO: Implement using Vercel Blob HTTP API PUT endpoint
	return nil
}

// Load retrieves concerts from Vercel Blob storage.
func (c *HTTPClient) Load(ctx context.Context) ([]domain.Concert, error) {
	// TODO: Implement using Vercel Blob HTTP API GET endpoint
	return nil, nil
}

// uploadBlob uploads data to Vercel Blob.
func (c *HTTPClient) uploadBlob(ctx context.Context, pathname string, data io.Reader, contentType string, access string) error {
	// TODO: Implement upload using PUT /?pathname={pathname}
	return nil
}

// downloadBlob downloads data from Vercel Blob.
func (c *HTTPClient) downloadBlob(ctx context.Context, pathname string) (io.ReadCloser, error) {
	// TODO: Implement download using GET /?url={url}
	return nil, nil
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