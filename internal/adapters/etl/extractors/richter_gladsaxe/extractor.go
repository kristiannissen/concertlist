// Package richter_gladsaxe provides the extractor for Richter Gladsaxe venue.
package richter_gladsaxe

import (
	"context"

	"github.com/gocolly/colly/v2"
	"github.com/kristiannissen/concertlist/internal/domain"
)

// Extractor implements domain.ExtractorPort for Richter Gladsaxe.
type Extractor struct{}

// NewExtractor creates a new Extractor for Richter Gladsaxe.
func NewExtractor() *Extractor {
	return &Extractor{}
}

// Extract fetches concert data from Richter Gladsaxe's website.
func (e *Extractor) Extract(ctx context.Context) ([]domain.Concert, error) {
	var concerts []domain.Concert

	// Initialize gocolly collector.
	c := colly.NewCollector(
		colly.AllowedDomains("richter-gladsaxe.dk"),
	)

	// Handle HTML parsing.
	c.OnHTML("div.event", func(h *colly.HTMLElement) {
		concert := domain.Concert{
			Title: h.ChildText("h2"),
			Date:  h.ChildText("span.date"),
			URL:   h.Request.AbsoluteURL(h.Attr("href")),
			Venue: "Richter Gladsaxe",
		}
		concerts = append(concerts, concert)
	})

	// Handle context cancellation.
	go func() {
		<-ctx.Done()
		c.Stop()
	}()

	// Start scraping.
	if err := c.Visit("https://richter-gladsaxe.dk/"); err != nil {
		return nil, err
	}

	c.Wait()
	return concerts, nil
}
