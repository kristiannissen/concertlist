// Package blob provides a ports.Blob implementation backed by Vercel Blob.
//
// Vercel doesn't publish a versioned public REST spec for this outside its
// own SDKs (there's no official Go client), so this wraps the same HTTP
// contract the @vercel/blob JS SDK uses internally: a control-plane API at
// https://vercel.com/api/blob, authenticated with a BLOB_READ_WRITE_TOKEN.
// Because that's an internal contract rather than a documented one, re-check
// https://github.com/vercel/storage/tree/main/packages/blob/src if requests
// start failing after a Vercel platform change.
package blob

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/kristiannissen/concertlist/internal/ports"
)

const (
	apiBase    = "https://vercel.com/api/blob"
	apiVersion = "12"
)

// VercelBlob is a ports.Blob implementation backed by Vercel Blob storage.
type VercelBlob struct {
	token   string
	storeID string
	client  *resty.Client
}

// New builds a VercelBlob from a BLOB_READ_WRITE_TOKEN. The store ID doesn't
// need to be supplied separately — it's embedded in the token itself
// (vercel_blob_rw_<storeId>_<random>).
func New(token string) (*VercelBlob, error) {
	parts := strings.Split(token, "_")
	if len(parts) < 4 {
		return nil, fmt.Errorf("blob: malformed BLOB_READ_WRITE_TOKEN")
	}

	return &VercelBlob{
		token:   token,
		storeID: parts[3],
		client:  resty.New(),
	}, nil
}

// request returns a resty request pre-populated with the headers every
// Vercel Blob API call needs.
func (b *VercelBlob) request(ctx context.Context) *resty.Request {
	return b.client.R().
		SetContext(ctx).
		SetAuthToken(b.token).
		SetHeader("x-api-version", apiVersion).
		SetHeader("x-vercel-blob-store-id", b.storeID)
}

// Put uploads body to pathname and returns the stored object's metadata.
// Configure the upload with the ports.With* PutOption functions.
func (b *VercelBlob) Put(ctx context.Context, pathname string, body []byte, opts ...ports.PutOption) (ports.Object, error) {
	cfg := ports.ResolvePutOptions(opts...)

	req := b.request(ctx).
		SetHeader("x-vercel-blob-access", string(cfg.Access)).
		SetQueryParam("pathname", pathname).
		SetBody(body)

	if cfg.ContentType != "" {
		req.SetHeader("x-content-type", cfg.ContentType)
	}
	if cfg.AddRandomSuffix {
		req.SetHeader("x-add-random-suffix", "1")
	}
	if cfg.AllowOverwrite {
		req.SetHeader("x-allow-overwrite", "1")
	}
	if cfg.CacheControlMaxAge > 0 {
		req.SetHeader("x-cache-control-max-age", strconv.Itoa(cfg.CacheControlMaxAge))
	}

	resp, err := req.Put(apiBase + "/")
	if err != nil {
		return ports.Object{}, fmt.Errorf("blob: put %q: %w", pathname, err)
	}
	if resp.IsError() {
		return ports.Object{}, fmt.Errorf("blob: put %q failed (%d): %s", pathname, resp.StatusCode(), resp.String())
	}

	var out struct {
		URL         string `json:"url"`
		DownloadURL string `json:"downloadUrl"`
		Pathname    string `json:"pathname"`
		ContentType string `json:"contentType"`
		ETag        string `json:"etag"`
	}
	if err := json.Unmarshal(resp.Body(), &out); err != nil {
		return ports.Object{}, fmt.Errorf("blob: put %q: decode response: %w", pathname, err)
	}

	return ports.Object{
		URL:         out.URL,
		DownloadURL: out.DownloadURL,
		Pathname:    out.Pathname,
		ContentType: out.ContentType,
		ETag:        out.ETag,
	}, nil
}

// Get downloads the contents previously stored at pathname. It resolves
// pathname to the object's canonical URL via the metadata endpoint first,
// since the control-plane API only ever returns metadata, never the body —
// the actual bytes have to be fetched from the CDN URL it hands back.
func (b *VercelBlob) Get(ctx context.Context, pathname string) ([]byte, error) {
	resp, err := b.request(ctx).
		SetQueryParam("url", pathname).
		Get(apiBase + "/")
	if err != nil {
		return nil, fmt.Errorf("blob: head %q: %w", pathname, err)
	}
	if resp.IsError() {
		return nil, fmt.Errorf("blob: head %q failed (%d): %s", pathname, resp.StatusCode(), resp.String())
	}

	var meta struct {
		URL string `json:"url"`
	}
	if err := json.Unmarshal(resp.Body(), &meta); err != nil {
		return nil, fmt.Errorf("blob: head %q: decode response: %w", pathname, err)
	}

	dl, err := b.client.R().SetContext(ctx).Get(meta.URL)
	if err != nil {
		return nil, fmt.Errorf("blob: download %q: %w", pathname, err)
	}
	if dl.IsError() {
		return nil, fmt.Errorf("blob: download %q failed (%d)", pathname, dl.StatusCode())
	}

	return dl.Body(), nil
}

// List returns one page of objects, optionally filtered with
// ports.WithPrefix / ports.WithLimit. Pass ports.WithCursor(previousPage.NextCursor)
// to page forward; keep calling List while the returned ListPage.HasMore is true.
func (b *VercelBlob) List(ctx context.Context, opts ...ports.ListOption) (ports.ListPage, error) {
	cfg := ports.ResolveListOptions(opts...)
	req := b.request(ctx)

	if cfg.Prefix != "" {
		req.SetQueryParam("prefix", cfg.Prefix)
	}
	if cfg.Limit > 0 {
		req.SetQueryParam("limit", strconv.Itoa(cfg.Limit))
	}
	if cfg.Cursor != "" {
		req.SetQueryParam("cursor", cfg.Cursor)
	}

	resp, err := req.Get(apiBase + "/")
	if err != nil {
		return ports.ListPage{}, fmt.Errorf("blob: list: %w", err)
	}
	if resp.IsError() {
		return ports.ListPage{}, fmt.Errorf("blob: list failed (%d): %s", resp.StatusCode(), resp.String())
	}

	var out struct {
		Blobs []struct {
			URL         string `json:"url"`
			DownloadURL string `json:"downloadUrl"`
			Pathname    string `json:"pathname"`
			Size        int64  `json:"size"`
			UploadedAt  string `json:"uploadedAt"`
			ETag        string `json:"etag"`
		} `json:"blobs"`
		Cursor  string `json:"cursor"`
		HasMore bool   `json:"hasMore"`
	}
	if err := json.Unmarshal(resp.Body(), &out); err != nil {
		return ports.ListPage{}, fmt.Errorf("blob: list: decode response: %w", err)
	}

	page := ports.ListPage{
		NextCursor: out.Cursor,
		HasMore:    out.HasMore,
		Objects:    make([]ports.Object, 0, len(out.Blobs)),
	}
	for _, o := range out.Blobs {
		// Vercel returns uploadedAt as an ISO 8601 string with millisecond
		// precision (e.g. "2024-04-05T19:12:36.679Z"). time.Parse accepts
		// the fractional-seconds component even though time.RFC3339 doesn't
		// spell it out in the layout.
		uploadedAt, _ := time.Parse(time.RFC3339, o.UploadedAt)
		page.Objects = append(page.Objects, ports.Object{
			URL:         o.URL,
			DownloadURL: o.DownloadURL,
			Pathname:    o.Pathname,
			Size:        o.Size,
			UploadedAt:  uploadedAt,
			ETag:        o.ETag,
		})
	}

	return page, nil
}

// Delete removes the object(s) at the given pathname(s) or URL(s).
func (b *VercelBlob) Delete(ctx context.Context, pathnames ...string) error {
	if len(pathnames) == 0 {
		return nil
	}

	body, err := json.Marshal(map[string][]string{"urls": pathnames})
	if err != nil {
		return fmt.Errorf("blob: delete: encode request: %w", err)
	}

	resp, err := b.request(ctx).
		SetHeader("Content-Type", "application/json").
		SetBody(body).
		Post(apiBase + "/delete")
	if err != nil {
		return fmt.Errorf("blob: delete: %w", err)
	}
	if resp.IsError() {
		return fmt.Errorf("blob: delete failed (%d): %s", resp.StatusCode(), resp.String())
	}

	return nil
}
