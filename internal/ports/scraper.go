// This file is in internal/ports
package ports

import (
	"context"
	"sync"
)

type Scraper interface {
	Scrape(ctx context.Context, wg *sync.WaitGroup) error
	Extract(ctx context.Context, wg *sync.WaitGroup) error
}
