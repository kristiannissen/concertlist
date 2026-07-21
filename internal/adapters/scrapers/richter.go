// internal/ports.Scraper for specific venues/sources.
package scrapers

import (
	"context"
	"sync"

	"go.uber.org/zap"
)

// Richter is a Scraper adapter for richter-gladsaxe.dk.
type Richter struct {
	URL string
	Log *zap.Logger

	// visited guards against colly's default revisit-prevention racing under
	// Async + Parallelism: two goroutines can both pass colly's internal
	// "not yet visited" check before either is marked visited, causing the
	// same URL to be fetched more than once. LoadOrStore is atomic, so this
	// closes that race regardless of parallelism.
	visited sync.Map
}

func (r *Richter) Scrape(ctx context.Context, wg *sync.WaitGroup) error {
	return nil
}

func (r *Richter) Extract(ctx context.Context, wg *sync.WaitGroup, URL string) error {
	return nil
}
