// Package richter_gladsaxe provides the extractor for Richter Gladsaxe venue.
package richter_gladsaxe

import (
	"context"
	"log"
	"strings"

	"github.com/gocolly/colly/v2"
	"github.com/kristiannissen/concertlist/pkg/domain"
)

// Extractor implements domain.ExtractorPort for Richter Gladsaxe.
type Extractor struct{}

// NewExtractor creates a new Extractor for Richter Gladsaxe.
func NewExtractor() *Extractor {
	return &Extractor{}
}

// Extract fetches concert data from Richter Gladsaxe's index page.
func (e *Extractor) Extract(ctx context.Context) ([]domain.Concert, error) {
	var concerts []domain.Concert

	// Initialize gocolly collector.
	c := colly.NewCollector(
		colly.AllowedDomains("richter-gladsaxe.dk"),
		colly.UserAgent("Mozilla/5.0 (compatible; ConcertList/1.0)"),
	)

	// Handle context cancellation. colly.Collector has no Stop method, so
	// abort any request that fires after the context has been cancelled.
	c.OnRequest(func(r *colly.Request) {
		select {
		case <-ctx.Done():
			r.Abort()
		default:
		}
	})

	// On every .card .text-overlay element, extract date and title.
	c.OnHTML(".card .text-overlay", func(h *colly.HTMLElement) {
		// Extract date and title from the text overlay.
		date := h.ChildText(".date")
		title := h.ChildText(".title")

		// Clean up whitespace.
		date = strings.TrimSpace(date)
		title = strings.TrimSpace(title)

		// Skip if no title or date.
		if title == "" || date == "" {
			return
		}

		// Create a new Concert.
		concert := domain.Concert{
			ID:    generateID(title, date), // Simple ID generation.
			Title: title,
			Date:  date,
			Venue: "Richter Gladsaxe",
			URL:   h.Request.AbsoluteURL(h.Attr("href")),
		}

		concerts = append(concerts, concert)
	})

	// Log errors during scraping.
	c.OnError(func(r *colly.Response, err error) {
		log.Printf("Failed to scrape %s: %v", r.Request.URL, err)
	})

	// Start scraping the index page.
	if err := c.Visit("https://richter-gladsaxe.dk/"); err != nil {
		return nil, err
	}

	// Wait for all requests to finish.
	c.Wait()

	return concerts, nil
}

// generateID creates a simple ID from the title and date.
func generateID(title, date string) string {
	return strings.ToLower(strings.ReplaceAll(title+"-"+date, " ", "-"))
}
