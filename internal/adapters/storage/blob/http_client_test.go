// Package blob provides HTTP client implementation for Vercel Blob storage.
package blob

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/kristiannissen/concertlist/internal/domain"
)

func TestNewHTTPClient(t *testing.T) {
	t.Parallel()

	config := HTTPClientConfig{
		StoreID:       "test-store",
		AccessToken:   "test-token",
		APIVersion:    "12",
		BaseURL:       "https://vercel.com/api/blob",
		StorageBaseURL: "https://test-store.public.blob.vercel-storage.com",
	}

	client := NewHTTPClient(config)

	if client.storeID != config.StoreID {
		t.Errorf("Expected storeID to be %s, got %s", config.StoreID, client.storeID)
	}
	if client.accessToken != config.AccessToken {
		t.Errorf("Expected accessToken to be %s, got %s", config.AccessToken, client.accessToken)
	}
	if client.apiVersion != config.APIVersion {
		t.Errorf("Expected apiVersion to be %s, got %s", config.APIVersion, client.apiVersion)
	}
}

func TestHTTPClient_Save(t *testing.T) {
	t.Parallel()

	// Create test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify request method
		if r.Method != http.MethodPut {
			t.Errorf("Expected PUT method, got %s", r.Method)
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		// Verify pathname query parameter
		if r.URL.Query().Get("pathname") != "concerts.json" {
			t.Errorf("Expected pathname=concerts.json, got %s", r.URL.Query().Get("pathname"))
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		// Verify headers
		if r.Header.Get("Authorization") != "Bearer test-token" {
			t.Errorf("Expected Authorization header to be Bearer test-token, got %s", r.Header.Get("Authorization"))
		}
		if r.Header.Get("x-vercel-blob-store-id") != "test-store" {
			t.Errorf("Expected x-vercel-blob-store-id header to be test-store, got %s", r.Header.Get("x-vercel-blob-store-id"))
		}
		if r.Header.Get("x-api-version") != "12" {
			t.Errorf("Expected x-api-version header to be 12, got %s", r.Header.Get("x-api-version"))
		}
		if r.Header.Get("x-vercel-blob-access") != "public" {
			t.Errorf("Expected x-vercel-blob-access header to be public, got %s", r.Header.Get("x-vercel-blob-access"))
		}
		if r.Header.Get("x-content-type") != "application/json" {
			t.Errorf("Expected x-content-type header to be application/json, got %s", r.Header.Get("x-content-type"))
		}
		if r.Header.Get("x-allow-overwrite") != "1" {
			t.Errorf("Expected x-allow-overwrite header to be 1, got %s", r.Header.Get("x-allow-overwrite"))
		}

		// Read and verify body
		body, err := io.ReadAll(r.Body)
		if err != nil {
			t.Errorf("Failed to read request body: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		var concerts []domain.Concert
		if err := json.Unmarshal(body, &concerts); err != nil {
			t.Errorf("Failed to unmarshal request body: %v", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		// Verify we received the expected concert
		if len(concerts) != 1 {
			t.Errorf("Expected 1 concert, got %d", len(concerts))
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if concerts[0].ID != "1" {
			t.Errorf("Expected concert ID to be 1, got %s", concerts[0].ID)
		}

		// Return success
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	config := HTTPClientConfig{
		StoreID:       "test-store",
		AccessToken:   "test-token",
		APIVersion:    "12",
		BaseURL:       server.URL,
		StorageBaseURL: "https://test-store.public.blob.vercel-storage.com",
	}

	client := NewHTTPClient(config)
	ctx := context.Background()
	concerts := []domain.Concert{
		{ID: "1", Title: "Test Concert", Venue: "Test Venue", Date: "2026-07-02"},
	}

	err := client.Save(ctx, concerts)
	if err != nil {
		t.Errorf("Save failed: %v", err)
	}
}

func TestHTTPClient_Save_EmptyConcerts(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify we can handle empty concerts array
		body, _ := io.ReadAll(r.Body)
		var concerts []domain.Concert
		json.Unmarshal(body, &concerts)
		
		if len(concerts) != 0 {
			t.Errorf("Expected 0 concerts, got %d", len(concerts))
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	config := HTTPClientConfig{
		StoreID:       "test-store",
		AccessToken:   "test-token",
		APIVersion:    "12",
		BaseURL:       server.URL,
		StorageBaseURL: "https://test-store.public.blob.vercel-storage.com",
	}

	client := NewHTTPClient(config)
	ctx := context.Background()
	concerts := []domain.Concert{}

	err := client.Save(ctx, concerts)
	if err != nil {
		t.Errorf("Save failed with empty concerts: %v", err)
	}
}

func TestHTTPClient_Save_ServerError(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal Server Error"))
	}))
	defer server.Close()

	config := HTTPClientConfig{
		StoreID:       "test-store",
		AccessToken:   "test-token",
		APIVersion:    "12",
		BaseURL:       server.URL,
		StorageBaseURL: "https://test-store.public.blob.vercel-storage.com",
	}

	client := NewHTTPClient(config)
	ctx := context.Background()
	concerts := []domain.Concert{
		{ID: "1", Title: "Test Concert", Venue: "Test Venue", Date: "2026-07-02"},
	}

	err := client.Save(ctx, concerts)
	if err == nil {
		t.Error("Expected Save to fail with server error")
	}
	if !strings.Contains(err.Error(), "500") {
		t.Errorf("Expected error to contain 500, got: %s", err.Error())
	}
}

func TestHTTPClient_Load(t *testing.T) {
	t.Parallel()

	// Create test server that returns sample concerts
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify request method
		if r.Method != http.MethodGet {
			t.Errorf("Expected GET method, got %s", r.Method)
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		// Verify Authorization header
		if r.Header.Get("Authorization") != "Bearer test-token" {
			t.Errorf("Expected Authorization header to be Bearer test-token, got %s", r.Header.Get("Authorization"))
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		// Return sample concerts JSON
		sampleConcerts := []domain.Concert{
			{ID: "1", Title: "Test Concert 1", Venue: "Test Venue 1", Date: "2026-07-01"},
			{ID: "2", Title: "Test Concert 2", Venue: "Test Venue 2", Date: "2026-07-02"},
		}
		
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(sampleConcerts)
	}))
	defer server.Close()

	config := HTTPClientConfig{
		StoreID:       "test-store",
		AccessToken:   "test-token",
		APIVersion:    "12",
		BaseURL:       "https://vercel.com/api/blob",
		StorageBaseURL: server.URL,
	}

	client := NewHTTPClient(config)
	ctx := context.Background()

	concerts, err := client.Load(ctx)
	if err != nil {
		t.Errorf("Load failed: %v", err)
	}

	if len(concerts) != 2 {
		t.Errorf("Expected 2 concerts, got %d", len(concerts))
	}

	if concerts[0].ID != "1" {
		t.Errorf("Expected first concert ID to be 1, got %s", concerts[0].ID)
	}

	if concerts[1].ID != "2" {
		t.Errorf("Expected second concert ID to be 2, got %s", concerts[1].ID)
	}
}

func TestHTTPClient_Load_NotFound(t *testing.T) {
	t.Parallel()

	// Create test server that returns 404
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	config := HTTPClientConfig{
		StoreID:       "test-store",
		AccessToken:   "test-token",
		APIVersion:    "12",
		BaseURL:       "https://vercel.com/api/blob",
		StorageBaseURL: server.URL,
	}

	client := NewHTTPClient(config)
	ctx := context.Background()

	concerts, err := client.Load(ctx)
	if err != nil {
		t.Errorf("Load should not return error for 404, got: %v", err)
	}

	if len(concerts) != 0 {
		t.Errorf("Expected empty concerts for 404, got %d", len(concerts))
	}
}

func TestHTTPClient_Load_InvalidJSON(t *testing.T) {
	t.Parallel()

	// Create test server that returns invalid JSON
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("invalid json"))
	}))
	defer server.Close()

	config := HTTPClientConfig{
		StoreID:       "test-store",
		AccessToken:   "test-token",
		APIVersion:    "12",
		BaseURL:       "https://vercel.com/api/blob",
		StorageBaseURL: server.URL,
	}

	client := NewHTTPClient(config)
	ctx := context.Background()

	_, err := client.Load(ctx)
	if err == nil {
		t.Error("Expected Load to fail with invalid JSON")
	}
	if !strings.Contains(err.Error(), "failed to unmarshal concerts") {
		t.Errorf("Expected unmarshal error, got: %s", err.Error())
	}
}

func TestHTTPClient_Load_ServerError(t *testing.T) {
	t.Parallel()

	// Create test server that returns server error
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal Server Error"))
	}))
	defer server.Close()

	config := HTTPClientConfig{
		StoreID:       "test-store",
		AccessToken:   "test-token",
		APIVersion:    "12",
		BaseURL:       "https://vercel.com/api/blob",
		StorageBaseURL: server.URL,
	}

	client := NewHTTPClient(config)
	ctx := context.Background()

	_, err := client.Load(ctx)
	if err == nil {
		t.Error("Expected Load to fail with server error")
	}
	if !strings.Contains(err.Error(), "500") {
		t.Errorf("Expected error to contain 500, got: %s", err.Error())
	}
}

func TestHTTPClient_uploadBlob(t *testing.T) {
	t.Parallel()

	config := HTTPClientConfig{
		StoreID:       "test-store",
		AccessToken:   "test-token",
		APIVersion:    "12",
		BaseURL:       "https://vercel.com/api/blob",
		StorageBaseURL: "https://test-store.public.blob.vercel-storage.com",
	}

	client := NewHTTPClient(config)
	ctx := context.Background()

	// TODO: Implement test with mock HTTP server
	err := client.uploadBlob(ctx, "test.txt", nil, "text/plain", "public")
	if err != nil {
		t.Logf("uploadBlob returned error (expected for skeleton): %v", err)
	}
}

func TestHTTPClient_downloadBlob(t *testing.T) {
	t.Parallel()

	config := HTTPClientConfig{
		StoreID:       "test-store",
		AccessToken:   "test-token",
		APIVersion:    "12",
		BaseURL:       "https://vercel.com/api/blob",
		StorageBaseURL: "https://test-store.public.blob.vercel-storage.com",
	}

	client := NewHTTPClient(config)
	ctx := context.Background()

	// TODO: Implement test with mock HTTP server
	reader, err := client.downloadBlob(ctx, "test.txt")
	if err != nil {
		t.Logf("downloadBlob returned error (expected for skeleton): %v", err)
	}
	if reader != nil {
		defer reader.Close()
	}
}

func TestHTTPClient_listBlobs(t *testing.T) {
	t.Parallel()

	config := HTTPClientConfig{
		StoreID:       "test-store",
		AccessToken:   "test-token",
		APIVersion:    "12",
		BaseURL:       "https://vercel.com/api/blob",
		StorageBaseURL: "https://test-store.public.blob.vercel-storage.com",
	}

	client := NewHTTPClient(config)
	ctx := context.Background()

	// TODO: Implement test with mock HTTP server
	concerts, err := client.listBlobs(ctx, "")
	if err != nil {
		t.Logf("listBlobs returned error (expected for skeleton): %v", err)
	}
	if concerts != nil {
		t.Logf("Listed %d concerts", len(concerts))
	}
}

func TestHTTPClient_deleteBlob(t *testing.T) {
	t.Parallel()

	config := HTTPClientConfig{
		StoreID:       "test-store",
		AccessToken:   "test-token",
		APIVersion:    "12",
		BaseURL:       "https://vercel.com/api/blob",
		StorageBaseURL: "https://test-store.public.blob.vercel-storage.com",
	}

	client := NewHTTPClient(config)
	ctx := context.Background()

	// TODO: Implement test with mock HTTP server
	err := client.deleteBlob(ctx, "test.txt")
	if err != nil {
		t.Logf("deleteBlob returned error (expected for skeleton): %v", err)
	}
}

func TestHTTPClient_deleteBlobs(t *testing.T) {
	t.Parallel()

	config := HTTPClientConfig{
		StoreID:       "test-store",
		AccessToken:   "test-token",
		APIVersion:    "12",
		BaseURL:       "https://vercel.com/api/blob",
		StorageBaseURL: "https://test-store.public.blob.vercel-storage.com",
	}

	client := NewHTTPClient(config)
	ctx := context.Background()
	urls := []string{"https://test-store.public.blob.vercel-storage.com/test1.txt"}

	// TODO: Implement test with mock HTTP server
	err := client.deleteBlobs(ctx, urls)
	if err != nil {
		t.Logf("deleteBlobs returned error (expected for skeleton): %v", err)
	}
}

func TestHTTPClient_getBlobMetadata(t *testing.T) {
	t.Parallel()

	config := HTTPClientConfig{
		StoreID:       "test-store",
		AccessToken:   "test-token",
		APIVersion:    "12",
		BaseURL:       "https://vercel.com/api/blob",
		StorageBaseURL: "https://test-store.public.blob.vercel-storage.com",
	}

	client := NewHTTPClient(config)
	ctx := context.Background()
	url := "https://test-store.public.blob.vercel-storage.com/test.txt"

	// TODO: Implement test with mock HTTP server
	metadata, err := client.getBlobMetadata(ctx, url)
	if err != nil {
		t.Logf("getBlobMetadata returned error (expected for skeleton): %v", err)
	}
	if metadata != nil {
		t.Logf("Got metadata: %v", metadata)
	}
}

func TestHTTPClient_copyBlob(t *testing.T) {
	t.Parallel()

	config := HTTPClientConfig{
		StoreID:       "test-store",
		AccessToken:   "test-token",
		APIVersion:    "12",
		BaseURL:       "https://vercel.com/api/blob",
		StorageBaseURL: "https://test-store.public.blob.vercel-storage.com",
	}

	client := NewHTTPClient(config)
	ctx := context.Background()
	fromURL := "https://test-store.public.blob.vercel-storage.com/source.txt"
	toPathname := "destination.txt"

	// TODO: Implement test with mock HTTP server
	err := client.copyBlob(ctx, toPathname, fromURL)
	if err != nil {
		t.Logf("copyBlob returned error (expected for skeleton): %v", err)
	}
}