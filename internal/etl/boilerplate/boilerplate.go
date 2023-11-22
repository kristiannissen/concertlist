package boilerplate

import (
	d "concertlist/internal/dataaccess"
	"errors"
	"log"
	"net/url"
	"time"

	"github.com/gocolly/colly"
)

var coll *colly.Collector

/**
 * Should implement an interface
 */
type Extractor interface {
	Run() ([]d.Event, error)
	New(c *colly.Collector) error
}

type Resource struct {
	URL   string
	Delay time.Duration
}

// New takes the collector and adds custom rules
// and options to fit the target site
func (r Resource) New(c *colly.Collector) error {
	u, _ := url.Parse(r.URL)
	log.Println("url ", r.URL, u)

	c.Limit(&colly.LimitRule{
		Delay:       r.Delay,
		RandomDelay: r.Delay,
		Parallelism: 2,
	})
	c.AllowedDomains = []string{u.Host}

	coll = c.Clone()

	return errors.New("Shit!")
}

/**
 * Runs the OnHTML methods to collect links and create
 * []struct containing data object
 */
func (r Resource) Run() ([]d.Event, error) {
	var events []d.Event

	coll.OnHTML("a[href]", func(h *colly.HTMLElement) {
		coll.Visit(h.Request.AbsoluteURL(h.Attr("href")))
	})

	coll.OnHTML(".product_main", func(h *colly.HTMLElement) {
		events = append(events, d.Event{Artist: h.ChildText("h1"), Venue: "Boilerplate"})
	})

	coll.Visit(r.URL)
	coll.Wait()

	return events, nil
}
