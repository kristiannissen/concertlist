// This file is in internal/ports
package ports

import (
	"context"

	"github.com/kristiannissen/concertlist/internal/domain"
)

type Scraper interface {
	Parse(ctx context.Context, data []byte, contentType string) (ScrapeResult, error)
}

type ScrapeResult struct {
	Events   []domain.MusicEvent
	NextURLs []string
}
