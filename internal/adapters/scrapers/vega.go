// internal/ports.Scraper for specific venues/sources.
package scrapers

import (
	"context"
	"net/url"
	"os"
	"sync"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/gocolly/colly"
	"github.com/kristiannissen/concertlist/internal/domain"
	"go.uber.org/zap"
)

// Vega is a Scraper adapter for Vega-gladsaxe.dk.
type Vega struct {
	URL string
	Log *zap.Logger

	// visited guards against colly's default revisit-prevention racing under
	// Async + Parallelism: two goroutines can both pass colly's internal
	// "not yet visited" check before either is marked visited, causing the
	// same URL to be fetched more than once. LoadOrStore is atomic, so this
	// closes that race regardless of parallelism.
	visited sync.Map
}

func (r *Vega) Scrape(ctx context.Context, wg *sync.WaitGroup) error {
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

	// Custom scraper job
	c.OnHTML("[data-theme='secondary']", func(e *colly.HTMLElement) {
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
	//
	client := resty.New()
	client.SetAuthToken(os.Getenv("VERCEL_OIDC_TOKEN"))
	// Scrape data
	c.OnHTML(".single-concert", func(e *colly.HTMLElement) {
		//
		title := e.ChildText("#concertTitle")
		if title == "" {
			// .single-concert matches the outer block, but #concertTitle
			// isn't guaranteed to exist inside every element carrying that
			// class (e.g. listing/teaser cards reuse it for styling without
			// the detail-page markup). Skip rather than log/send a blank
			// event.
			r.Log.Warn("skipping single-concert match with no title",
				zap.String("URL", e.Request.URL.String()))
			return
		}

		wg.Add(1)
		go func() {
			defer wg.Done()

			m := domain.MusicEvent{
				Context:   "https://schema.org",
				Type:      "MusicEvent",
				Name:      title,
				StartDate: e.ChildText("#concertDate"),
				Location: domain.Location{
					Type: "MusicVenue",
					Name: "Vega",
					Address: domain.PostalAddress{
						Type:       "",
						Street:     "Telefonvej 16",
						Locality:   "Søborg",
						PostalCode: "2860",
						Country:    "Danmark",
					},
				},
				Performer: domain.Performer{
					Type: "MusicGroup",
					Name: title,
				},
				Offer: domain.Offer{
					Type: "Offer",
					URL:  e.Request.AbsoluteURL(e.Request.URL.String()),
				},
			}
			//
			client.R().SetBody(m).Post("https://arn1.vercel-queue.com/api/v3/topic/musicevent")
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

func (r *Vega) Extract(ctx context.Context, wg *sync.WaitGroup) error {

	return nil
}
