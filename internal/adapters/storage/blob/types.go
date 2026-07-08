// Package blob provides types for Vercel Blob API.
package blob

// BlobAccess defines the access level for blobs.
type BlobAccess string

const (
	AccessPublic  BlobAccess = "public"
	AccessPrivate BlobAccess = "private"
)

// UploadOptions contains options for uploading a blob.
type UploadOptions struct {
	ContentType        string
	Access             BlobAccess
	AddRandomSuffix    bool
	AllowOverwrite     bool
	CacheControlMaxAge int
}

// BlobMetadata contains metadata about a blob.
type BlobMetadata struct {
	Pathname      string
	ContentType   string
	ContentLength int64
	UploadedAt    string
	Access        string
	URL           string
}

// ListBlobsResponse contains the response from listing blobs.
type ListBlobsResponse struct {
	Blobs      []BlobMetadata
	Cursor     string
	HasMore    bool
	TotalCount int
}

// DeleteBlobsRequest contains the request for deleting multiple blobs.
type DeleteBlobsRequest struct {
	URLs []string
}

// DeleteBlobsResponse contains the response from deleting blobs.
type DeleteBlobsResponse struct {
	Deleted []string
}

// CopyBlobRequest contains the request for copying a blob.
type CopyBlobRequest struct {
	FromURL string
}

// MultipartUploadCreateResponse contains the response from creating a multipart upload.
type MultipartUploadCreateResponse struct {
	UploadID string
	Key      string
}

// MultipartUploadPart contains information about a part in a multipart upload.
type MultipartUploadPart struct {
	PartNumber int
	ETag       string
}

// MultipartUploadCompleteRequest contains the request for completing a multipart upload.
type MultipartUploadCompleteRequest struct {
	Parts []MultipartUploadPart
}

// SignedTokenRequest contains the request for generating a signed token.
type SignedTokenRequest struct {
	Pathname            string
	AllowedContentTypes []string
	MaximumSizeInBytes  int
	ValidUntil          string
	AddRandomSuffix     bool
	Access              string
}

// SignedTokenResponse contains the response for a signed token.
type SignedTokenResponse struct {
	Token   string
	URL     string
	Expires string
}

// ErrorResponse contains error information from the API.
type ErrorResponse struct {
	Error   string
	Message string
}
