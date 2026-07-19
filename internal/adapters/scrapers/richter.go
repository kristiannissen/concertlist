// internal/ports.Scraper for specific venues/sources.
package scrapers

import (
	"context"
	"encoding/json"
	"net/url"
	"sync"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/gocolly/colly"
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
		// Identify ourselves so the venue can see who's crawling and reach
		// out (instead of just blocking) if this is ever a problem.
		colly.UserAgent("Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/150.0.0.0 Safari/537.36"),
	)
	c.Limit(&colly.LimitRule{DomainGlob: "*", Parallelism: 2, RandomDelay: 5 * time.Second})

	// Send queue request using resty
	client := resty.New()
	// Auth token from incoming request
	client.SetAuthToken("")
	// Custom scraper job it finds the relevant URL and passes it to the queue for processing
	c.OnHTML("a.card-img-top", func(e *colly.HTMLElement) {
		l := e.Request.AbsoluteURL(e.Attr("href"))

		if _, seen := r.visited.LoadOrStore(l, true); seen {
			// Exists the callback
			return
		}

		r.Log.Info("Visiting", zap.String("URL", l))
		// Payload
		body, _ := json.Marshal(map[string]string{"venue": "", "url": l})
		//
		resp, err := client.R().SetBody(body).Post("")
		if err != nil || resp.IsError() {
			r.Log.Error("Failed", zap.String("venue-url", l), zap.Error(err))
		}
		r.Log.Info("enqueue value url", zap.String("url", l))

		c.Visit(l)
	})
	// Error handling
	c.OnError(func(res *colly.Response, err error) {
		r.Log.Info("Error", zap.String("msg", err.Error()))
	})

	r.visited.Store(r.URL, true)

	if err := c.Visit(r.URL); err != nil {
		r.Log.Error("visit failed", zap.String("URL", r.URL), zap.Error(err))
		return err
	}
	c.Wait()

	return nil
}

func (r *Richter) Extract(ctx context.Context, wg *sync.WaitGroup, URL string) error {
	return nil
}
