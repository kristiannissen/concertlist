// Package blob provides types for Vercel Blob API.
package blob

import (
	"testing"
)

func TestBlobAccessConstants(t *testing.T) {
	t.Parallel()

	if AccessPublic != "public" {
		t.Errorf("Expected AccessPublic to be 'public', got %s", AccessPublic)
	}
	if AccessPrivate != "private" {
		t.Errorf("Expected AccessPrivate to be 'private', got %s", AccessPrivate)
	}
}

func TestUploadOptions(t *testing.T) {
	t.Parallel()

	opts := UploadOptions{
		ContentType:     "application/json",
		Access:          AccessPublic,
		AddRandomSuffix: true,
		AllowOverwrite:  true,
		CacheControlMaxAge: 3600,
	}

	if opts.ContentType != "application/json" {
		t.Errorf("Expected ContentType to be 'application/json', got %s", opts.ContentType)
	}
	if opts.Access != AccessPublic {
		t.Errorf("Expected Access to be 'public', got %s", opts.Access)
	}
	if !opts.AddRandomSuffix {
		t.Error("Expected AddRandomSuffix to be true")
	}
	if !opts.AllowOverwrite {
		t.Error("Expected AllowOverwrite to be true")
	}
	if opts.CacheControlMaxAge != 3600 {
		t.Errorf("Expected CacheControlMaxAge to be 3600, got %d", opts.CacheControlMaxAge)
	}
}

func TestBlobMetadata(t *testing.T) {
	t.Parallel()

	metadata := BlobMetadata{
		Pathname:      "test/file.json",
		ContentType:   "application/json",
		ContentLength: 1024,
		UploadedAt:    "2026-07-02T00:00:00Z",
		Access:        "public",
		URL:           "https://store.public.blob.vercel-storage.com/test/file.json",
	}

	if metadata.Pathname != "test/file.json" {
		t.Errorf("Expected Pathname to be 'test/file.json', got %s", metadata.Pathname)
	}
	if metadata.ContentLength != 1024 {
		t.Errorf("Expected ContentLength to be 1024, got %d", metadata.ContentLength)
	}
}

func TestListBlobsResponse(t *testing.T) {
	t.Parallel()

	blob1 := BlobMetadata{
		Pathname: "file1.json",
	}
	blob2 := BlobMetadata{
		Pathname: "file2.json",
	}

	response := ListBlobsResponse{
		Blobs:      []BlobMetadata{blob1, blob2},
		Cursor:     "next-page",
		HasMore:    true,
		TotalCount: 2,
	}

	if len(response.Blobs) != 2 {
		t.Errorf("Expected 2 blobs, got %d", len(response.Blobs))
	}
	if response.Cursor != "next-page" {
		t.Errorf("Expected Cursor to be 'next-page', got %s", response.Cursor)
	}
	if !response.HasMore {
		t.Error("Expected HasMore to be true")
	}
	if response.TotalCount != 2 {
		t.Errorf("Expected TotalCount to be 2, got %d", response.TotalCount)
	}
}

func TestDeleteBlobsRequest(t *testing.T) {
	t.Parallel()

	request := DeleteBlobsRequest{
		URLs: []string{
			"https://store.public.blob.vercel-storage.com/file1.json",
			"https://store.public.blob.vercel-storage.com/file2.json",
		},
	}

	if len(request.URLs) != 2 {
		t.Errorf("Expected 2 URLs, got %d", len(request.URLs))
	}
}

func TestDeleteBlobsResponse(t *testing.T) {
	t.Parallel()

	response := DeleteBlobsResponse{
		Deleted: []string{
			"https://store.public.blob.vercel-storage.com/file1.json",
		},
	}

	if len(response.Deleted) != 1 {
		t.Errorf("Expected 1 deleted URL, got %d", len(response.Deleted))
	}
}

func TestCopyBlobRequest(t *testing.T) {
	t.Parallel()

	request := CopyBlobRequest{
		FromURL: "https://store.public.blob.vercel-storage.com/source.json",
	}

	if request.FromURL != "https://store.public.blob.vercel-storage.com/source.json" {
		t.Errorf("Expected FromURL to match, got %s", request.FromURL)
	}
}

func TestMultipartUploadCreateResponse(t *testing.T) {
	t.Parallel()

	response := MultipartUploadCreateResponse{
		UploadID: "upload-123",
		Key:      "large-file.zip",
	}

	if response.UploadID != "upload-123" {
		t.Errorf("Expected UploadID to be 'upload-123', got %s", response.UploadID)
	}
	if response.Key != "large-file.zip" {
		t.Errorf("Expected Key to be 'large-file.zip', got %s", response.Key)
	}
}

func TestMultipartUploadPart(t *testing.T) {
	t.Parallel()

	part := MultipartUploadPart{
		PartNumber: 1,
		ETag:      "etag-123",
	}

	if part.PartNumber != 1 {
		t.Errorf("Expected PartNumber to be 1, got %d", part.PartNumber)
	}
	if part.ETag != "etag-123" {
		t.Errorf("Expected ETag to be 'etag-123', got %s", part.ETag)
	}
}

func TestMultipartUploadCompleteRequest(t *testing.T) {
	t.Parallel()

	part1 := MultipartUploadPart{PartNumber: 1, ETag: "etag-1"}
	part2 := MultipartUploadPart{PartNumber: 2, ETag: "etag-2"}

	request := MultipartUploadCompleteRequest{
		Parts: []MultipartUploadPart{part1, part2},
	}

	if len(request.Parts) != 2 {
		t.Errorf("Expected 2 parts, got %d", len(request.Parts))
	}
}

func TestSignedTokenRequest(t *testing.T) {
	t.Parallel()

	request := SignedTokenRequest{
		Pathname:           "upload.txt",
		AllowedContentTypes: []string{"text/plain"},
		MaximumSizeInBytes: 1048576,
		ValidUntil:         "2026-07-04T00:00:00Z",
		AddRandomSuffix:    true,
		Access:            "public",
	}

	if request.Pathname != "upload.txt" {
		t.Errorf("Expected Pathname to be 'upload.txt', got %s", request.Pathname)
	}
	if len(request.AllowedContentTypes) != 1 {
		t.Errorf("Expected 1 allowed content type, got %d", len(request.AllowedContentTypes))
	}
	if request.MaximumSizeInBytes != 1048576 {
		t.Errorf("Expected MaximumSizeInBytes to be 1048576, got %d", request.MaximumSizeInBytes)
	}
}

func TestSignedTokenResponse(t *testing.T) {
	t.Parallel()

	response := SignedTokenResponse{
		Token:   "signed-token-123",
		URL:     "https://store.public.blob.vercel-storage.com/upload.txt",
		Expires: "2026-07-04T00:00:00Z",
	}

	if response.Token != "signed-token-123" {
		t.Errorf("Expected Token to be 'signed-token-123', got %s", response.Token)
	}
	if response.URL != "https://store.public.blob.vercel-storage.com/upload.txt" {
		t.Errorf("Expected URL to match, got %s", response.URL)
	}
}

func TestErrorResponse(t *testing.T) {
	t.Parallel()

	errResp := ErrorResponse{
		Error:   "NotFound",
		Message: "Blob not found",
	}

	if errResp.Error != "NotFound" {
		t.Errorf("Expected Error to be 'NotFound', got %s", errResp.Error)
	}
	if errResp.Message != "Blob not found" {
		t.Errorf("Expected Message to be 'Blob not found', got %s", errResp.Message)
	}
}