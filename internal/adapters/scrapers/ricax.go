// Package scrapers holds concrete Scraper adapters that implement
// internal/ports.Scraper for specific venues/sources.
package scrapers

import (
	"context"
	"time"

	"github.com/gocolly/colly"
	"github.com/kristiannissen/concertlist/internal/ports"
	"go.uber.org/zap"
)

// RicAx is a Scraper adapter (currently a stub).
type RicAx struct {
	URL string
}

func (r *RicAx) Parse(ctx context.Context, data []byte, contentType string) (ports.ScrapeResult, error) {
	logger, _ := zap.NewDevelopment()
	defer logger.Sync()
	sugar := logger.Sugar()

	c := colly.NewCollector(
		colly.AllowedDomains(r.URL),
		colly.MaxDepth(2),
		colly.Async(true),
	)
	c.Limit(&colly.LimitRule{DomainGlob: "*", Parallelism: 2, RandomDelay: 5 * time.Second})

	// Custom scraper job
	c.OnHTML("a.card-img-top", func(e *colly.HTMLElement) {
		l := e.Attr("href")
		c.Visit(e.Request.AbsoluteURL(l))
	})
	// Error handling
	c.OnError(func(r *colly.Response, err error) {
		sugar.Errorf("Error %s", err.Error())
	})

	c.Visit(r.URL)
	c.Wait()

	return ports.ScrapeResult{}, nil
}
