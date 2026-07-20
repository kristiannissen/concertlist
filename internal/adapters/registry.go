// internal/adapters
package adapters

import (
	"os"

	blobadapter "github.com/kristiannissen/concertlist/internal/adapters/blob"
	"github.com/kristiannissen/concertlist/internal/adapters/scrapers"
	"github.com/kristiannissen/concertlist/internal/ports"
	"go.uber.org/zap"
)

func NewScraperRegistry(log *zap.Logger) map[string]ports.Scraper {
	// A missing/malformed BLOB_READ_WRITE_TOKEN shouldn't stop scraping —
	// Vega.Extract logs and skips the upload step when Blob is nil.
	var blob ports.Blob
	if b, err := blobadapter.New(os.Getenv("BLOB_READ_WRITE_TOKEN")); err != nil {
		log.Warn("blob client not configured", zap.Error(err))
	} else {
		blob = b
	}

	return map[string]ports.Scraper{
		"richter": &scrapers.Richter{URL: "https://richter-gladsaxe.dk", Log: log},
		"vega":    &scrapers.Vega{URL: "https://vega.dk", Log: log, Blob: blob},
	}
}
