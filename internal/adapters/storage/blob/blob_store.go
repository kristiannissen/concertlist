// Package blob provides Vercel Blob storage adapter.
package blob

import (
	"context"
	"encoding/json"
	"os"

	"github.com/kristiannissen/concertlist/internal/domain"
	"github.com/vercel/blob"
)

// BlobStore implements domain.StoragePort using Vercel Blob.
type BlobStore struct {
	client *blob.Client
}

// NewBlobStore creates a new BlobStore.
func NewBlobStore() (*BlobStore, error) {
	client := blob.NewClient()
	return &BlobStore{client: client}, nil
}

// Save stores concerts in Vercel Blob.
func (s *BlobStore) Save(ctx context.Context, concerts []domain.Concert) error {
	data, err := json.Marshal(concerts)
	if err != nil {
		return err
	}

	storeID := os.Getenv("BLOB_STORE_ID")
	if storeID == "" {
		return os.ErrNotExist
	}

	// Upload to Vercel Blob.
	_, err = s.client.Put(ctx, storeID+"/concerts.json", data, nil)
	return err
}

// Load retrieves concerts from Vercel Blob.
func (s *BlobStore) Load(ctx context.Context) ([]domain.Concert, error) {
	storeID := os.Getenv("BLOB_STORE_ID")
	if storeID == "" {
		return nil, os.ErrNotExist
	}

	// Download from Vercel Blob.
	resp, err := s.client.Get(ctx, storeID+"/concerts.json")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close() //nolint:errcheck

	var concerts []domain.Concert
	if err := json.NewDecoder(resp.Body).Decode(&concerts); err != nil {
		return nil, err
	}

	return concerts, nil
}
