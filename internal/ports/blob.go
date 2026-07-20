// internal/ports
package ports

import (
	"context"
	"time"
)

// BlobAccess controls whether a stored object is publicly reachable via its
// URL or requires an authenticated request.
type BlobAccess string

const (
	BlobAccessPublic  BlobAccess = "public"
	BlobAccessPrivate BlobAccess = "private"
)

// PutOptions is the resolved configuration for a Put call, built by applying
// PutOption values on top of ResolvePutOptions' defaults. Callers configure
// Put via the With* functions below rather than constructing PutOptions
// directly.
type PutOptions struct {
	Access             BlobAccess
	ContentType        string
	AddRandomSuffix    bool
	AllowOverwrite     bool
	CacheControlMaxAge int
}

// PutOption configures a single aspect of a Put call.
type PutOption func(*PutOptions)

// WithAccess overrides the default (BlobAccessPublic) visibility of the
// stored object.
func WithAccess(access BlobAccess) PutOption {
	return func(o *PutOptions) { o.Access = access }
}

// WithContentType sets the Content-Type reported when the object is later
// downloaded. Providers infer this from the pathname when omitted.
func WithContentType(contentType string) PutOption {
	return func(o *PutOptions) { o.ContentType = contentType }
}

// WithAddRandomSuffix appends a random suffix to the pathname, so repeated
// uploads under the same name don't collide.
func WithAddRandomSuffix() PutOption {
	return func(o *PutOptions) { o.AddRandomSuffix = true }
}

// WithAllowOverwrite permits Put to replace an existing object at the same
// pathname. Providers reject overwrites by default.
func WithAllowOverwrite() PutOption {
	return func(o *PutOptions) { o.AllowOverwrite = true }
}

// WithCacheControlMaxAge sets how long (in seconds) caches may serve the
// object before revalidating.
func WithCacheControlMaxAge(seconds int) PutOption {
	return func(o *PutOptions) { o.CacheControlMaxAge = seconds }
}

// ResolvePutOptions applies opts on top of sane defaults and returns the
// resolved configuration. Blob implementations call this once at the top of
// Put instead of each re-implementing the same defaulting/apply loop.
func ResolvePutOptions(opts ...PutOption) PutOptions {
	o := PutOptions{Access: BlobAccessPublic}
	for _, opt := range opts {
		opt(&o)
	}
	return o
}

// Object describes a stored blob, as returned by Put and List.
type Object struct {
	URL         string
	DownloadURL string
	Pathname    string
	ContentType string
	Size        int64
	UploadedAt  time.Time
	ETag        string
}

// ListOptions is the resolved configuration for a List call, built by
// applying ListOption values via ResolveListOptions.
type ListOptions struct {
	// Prefix restricts results to pathnames starting with this value.
	Prefix string
	// Limit caps the number of results in this page. Zero uses the
	// provider's default page size.
	Limit int
	// Cursor resumes listing from a previous ListPage.NextCursor. Leave
	// empty to fetch the first page.
	Cursor string
}

// ListOption configures a single aspect of a List call.
type ListOption func(*ListOptions)

// WithPrefix restricts List to pathnames starting with prefix.
func WithPrefix(prefix string) ListOption {
	return func(o *ListOptions) { o.Prefix = prefix }
}

// WithLimit caps the number of results returned in a single List page.
func WithLimit(limit int) ListOption {
	return func(o *ListOptions) { o.Limit = limit }
}

// WithCursor resumes a List call from a previous ListPage.NextCursor.
func WithCursor(cursor string) ListOption {
	return func(o *ListOptions) { o.Cursor = cursor }
}

// ResolveListOptions applies opts on top of the zero value and returns the
// resolved configuration.
func ResolveListOptions(opts ...ListOption) ListOptions {
	var o ListOptions
	for _, opt := range opts {
		opt(&o)
	}
	return o
}

// ListPage is one page of List results.
type ListPage struct {
	Objects []Object
	// NextCursor is passed back via WithCursor to fetch the next page. Only
	// meaningful when HasMore is true.
	NextCursor string
	HasMore    bool
}

// Blob is a reusable, provider-agnostic interface for object storage.
// Scrapers and other callers should depend on this interface rather than a
// concrete client, so they aren't coupled to a specific storage provider.
type Blob interface {
	// Put uploads body to pathname and returns the stored object's
	// metadata. Configure the upload with the With* PutOption functions,
	// e.g. Put(ctx, "foo.json", body, ports.WithContentType("application/json")).
	Put(ctx context.Context, pathname string, body []byte, opts ...PutOption) (Object, error)

	// Get downloads the contents previously stored at pathname.
	Get(ctx context.Context, pathname string) ([]byte, error)

	// List returns one page of objects, optionally filtered with WithPrefix
	// / WithLimit. Pass the returned ListPage.NextCursor back via
	// WithCursor to fetch subsequent pages; keep paging while
	// ListPage.HasMore is true.
	List(ctx context.Context, opts ...ListOption) (ListPage, error)

	// Delete removes the object(s) at the given pathname(s) or URL(s).
	Delete(ctx context.Context, pathnames ...string) error
}
