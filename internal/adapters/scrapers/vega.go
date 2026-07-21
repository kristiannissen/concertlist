// internal/ports.Scraper for specific venues/sources.
package scrapers

import (
	"context"
	"encoding/json"
	"net/url"
	"os"
	"sync"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/gocolly/colly"
	"github.com/kristiannissen/concertlist/internal/domain"
	"github.com/kristiannissen/concertlist/internal/ports"
	"go.uber.org/zap"
)

// __NEXT_DATA__
type NextData struct {
	Props struct {
		PageProps struct {
			Data EventData `json:"data"`
		} `json:"pageProps"`
	} `json:"props"`
}

type EventData struct {
	Title     string    `json:"title"`
	Slug      string    `json:"slug"`
	Price     int       `json:"price"`
	FirstDate time.Time `json:"firstDate"`
	LastDate  time.Time `json:"lastDate"`
}

type Vega struct {
	URL  string
	Log  *zap.Logger
	Blob ports.Blob

	visited sync.Map
}

func (r *Vega) Scrape(ctx context.Context, wg *sync.WaitGroup) error {
	parsed, err := url.Parse(r.URL)
	if err != nil {
		return err
	}

	c := colly.NewCollector(
		colly.AllowedDomains(parsed.Hostname()),
		colly.MaxDepth(2),
		colly.Async(true),
		colly.UserAgent("Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/150.0.0.0 Safari/537.36"),
	)
	c.Limit(&colly.LimitRule{DomainGlob: "*", Parallelism: 2, RandomDelay: 5 * time.Second})
	// Resty client
	oidcToken, _ := ctx.Value("x-vercel-oidc-token").(string)
	client := resty.New()
	client.SetAuthToken(oidcToken)
	// Custom scraper job that looks for all relevant event URLs
	// and sends them off to a new queue for Extract()
	c.OnHTML("[data-theme='secondary']", func(e *colly.HTMLElement) {
		l := e.Request.AbsoluteURL(e.Attr("href"))

		if _, seen := r.visited.LoadOrStore(l, true); seen {
			return
		}
		//
		wg.Add(1)
		go func() {
			defer wg.Done()

			v := map[string]string{
				"venue": "vega",
				"url":   l,
			}

			_, err := client.R().
				SetHeader("Vqs-Deployment-Id", os.Getenv("VERCEL_DEPLOYMENT_ID")).
				SetBody(v).
				Post("https://arn1.vercel-queue.com/api/v3/topic/event-extract")
			if err != nil {
				r.Log.Error("post failed", zap.String("URL", l), zap.Error(err))
			}
		}()

		r.Log.Info("Visiting", zap.String("URL", l))
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

func (r *Vega) Extract(ctx context.Context, wg *sync.WaitGroup, URL string) error {
	parsed, err := url.Parse(URL)
	if err != nil {
		return err
	}

	c := colly.NewCollector(
		colly.AllowedDomains(parsed.Hostname()),
		colly.MaxDepth(2),
		colly.UserAgent("Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/150.0.0.0 Safari/537.36"),
	)
	c.Limit(&colly.LimitRule{DomainGlob: "*", Parallelism: 2, RandomDelay: 5 * time.Second})

	c.OnHTML("script#__NEXT_DATA__", func(e *colly.HTMLElement) {
		var next NextData
		if err := json.Unmarshal([]byte(e.Text), &next); err != nil {
			r.Log.Error("failed to parse", zap.Error(err))
			return
		}

		event := next.Props.PageProps.Data

		m := domain.MusicEvent{
			Context:   "https://schema.org",
			Type:      "MusicEvent",
			Name:      event.Title,
			StartDate: event.FirstDate.Format(time.RFC3339),
			Performer: domain.Performer{
				Type: "MusicGroup",
				Name: event.Title,
			},
			Offer: domain.Offer{
				Type:          "Offer",
				Price:         event.Price,
				PriceCurrency: "DKK",
				URL:           URL,
			},
		}

		if r.Blob == nil {
			r.Log.Warn("no blob client configured, skipping upload")
			return
		}

		body, err := json.Marshal(m)
		if err != nil {
			r.Log.Error("failed to marshal event", zap.Error(err))
			return
		}

		obj, err := r.Blob.Put(ctx, event.Slug+".json", body, ports.WithContentType("application/json"))
		if err != nil {
			r.Log.Error("failed to put blob", zap.Error(err))
			return
		}

		r.Log.Info("stored event", zap.String("url", obj.URL))
	})
	//
	c.Visit(URL)
	c.Wait()

	return nil
}
