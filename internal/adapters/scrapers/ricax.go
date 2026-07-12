// internal/ports.Scraper for specific venues/sources.
package scrapers

import (
	"context"
	"net/url"
	"sync"
	"time"

	"github.com/gocolly/colly"
	"github.com/kristiannissen/concertlist/internal/domain"
	"go.uber.org/zap"
)

// RicAx is a Scraper adapter (currently a stub).
type RicAx struct {
	URL string
	Log *zap.Logger

	// visited guards against colly's default revisit-prevention racing under
	// Async + Parallelism: two goroutines can both pass colly's internal
	// "not yet visited" check before either is marked visited, causing the
	// same URL to be fetched more than once. LoadOrStore is atomic, so this
	// closes that race regardless of parallelism.
	visited sync.Map
}

func (r *RicAx) Scrape(ctx context.Context, wg *sync.WaitGroup) error {
	// AllowedDomains needs a bare hostname, not the full URL. Passing the
	// full URL (scheme + path) means it never matches the request's actual
	// host, so colly silently rejects every request as forbidden-domain.
	parsed, err := url.Parse(r.URL)
	if err != nil {
		return err
	}

	c := colly.NewCollector(
		colly.AllowedDomains(parsed.Hostname()),
		colly.MaxDepth(2),
		colly.Async(true),
	)
	c.Limit(&colly.LimitRule{DomainGlob: "*", Parallelism: 2, RandomDelay: 5 * time.Second})

	// Custom scraper job
	c.OnHTML("a.card-img-top", func(e *colly.HTMLElement) {
		l := e.Request.AbsoluteURL(e.Attr("href"))

		if _, seen := r.visited.LoadOrStore(l, true); seen {
			return
		}

		r.Log.Info("Visiting", zap.String("URL", l))
		c.Visit(l)
	})
	// Error handling
	c.OnError(func(res *colly.Response, err error) {
		r.Log.Info("Error", zap.String("msg", err.Error()))
	})
	// Scrape data
	c.OnHTML(".single-concert", func(e *colly.HTMLElement) {
		//
		wg.Add(1)
		go func() {
			defer wg.Done()

			m := domain.MusicEvent{
				Name:      e.ChildText("#concertTitle"),
				StartDate: e.ChildText("#concertDate"),
			}
			r.Log.Info("Event", zap.String("title", m.Name))
		}()
	})

	r.visited.Store(r.URL, true)

	if err := c.Visit(r.URL); err != nil {
		r.Log.Error("visit failed", zap.String("URL", r.URL), zap.Error(err))
		return err
	}
	c.Wait()

	return nil
}

func (r *RicAx) Extract(ctx context.Context, wg *sync.WaitGroup) error {

	return nil
}
