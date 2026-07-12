// This file is in internal/ports
package ports

import (
	"context"
)

type Scraper interface {
	Scrape(ctx context.Context) error
	Extract(ctx context.Context) error
}
