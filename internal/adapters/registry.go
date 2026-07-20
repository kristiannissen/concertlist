// internal/adapters
package adapters

import (
	"os"

	blobadapter "github.com/kristiannissen/concertlist/internal/adapters/blob"
	"github.com/kristiannissen/concertlist/internal/adapters/scrapers"
	"github.com/kristiannissen/concertlist/internal/ports"
	"go.uber.org/zap"
)

// NewScraperRegistry builds the venue scrapers, wiring in a blob client when
// credentials are available. oidcToken should be the caller's
// x-vercel-oidc-token request header (see router.go) — this project's Blob
// store is OIDC-authenticated, not a static BLOB_READ_WRITE_TOKEN, so the
// token has to come from the current request rather than an env var. Pass
// "" if there's no request in scope (e.g. tests); Blob is left nil and
// Vega.Extract logs and skips the upload step rather than failing.
func NewScraperRegistry(log *zap.Logger, oidcToken string) map[string]ports.Scraper {
	blob := resolveBlob(log, oidcToken)

	return map[string]ports.Scraper{
		"richter": &scrapers.Richter{URL: "https://richter-gladsaxe.dk", Log: log},
		"vega":    &scrapers.Vega{URL: "https://vega.dk", Log: log, Blob: blob},
	}
}

func resolveBlob(log *zap.Logger, oidcToken string) ports.Blob {
	// A static BLOB_READ_WRITE_TOKEN takes priority if one's ever added
	// (e.g. local override), matching the priority order the official SDK
	// uses (explicit token > OIDC+storeId > env token).
	if token := os.Getenv("BLOB_READ_WRITE_TOKEN"); token != "" {
		if b, err := blobadapter.New(token); err != nil {
			log.Warn("blob client not configured", zap.Error(err))
		} else {
			return b
		}
	}

	if storeID := os.Getenv("BLOB_STORE_ID"); storeID != "" {
		if b, err := blobadapter.NewFromOIDC(oidcToken, storeID); err != nil {
			log.Warn("blob client not configured", zap.Error(err))
		} else {
			return b
		}
	}

	return nil
}
