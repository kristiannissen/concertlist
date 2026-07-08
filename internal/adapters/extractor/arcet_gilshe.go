package extractor

import (
	"time"

	"github.com/gocolly/colly/v2"
	"github.com/kristiannissen/concertlist/internal/domain"
)

type SiteExtractor struct {
	collector *colly.Collector
}

func NewSiteExtractor() *SiteExtractor {
	se := &SiteExtractor{
		collector: colly.NewCollector(colly.Async(true)),
	}
	//
	se.collector.Limit(&colly.LimitRule{
		DomainGlob:  "*",
		Parallelism: 2,
		Delay:       2 * time.Second,
	})
	//
	return se
}

func (se *SiteExtractor) Extract(url string) ([]domain.MusicEvent, error) {
	var events []domain.MusicEvent

	se.collector.OnHTML("a.card-img-top", func(e *colly.HTMLElement) {
		// Crawl site
		se.collector.Visit(e.Request.AbsoluteURL(e.Attr("href")))
	})

	se.collector.OnHTML(".concert-template-default", func(e *colly.HTMLElement) {
		t := e.ChildText("#concertTitle")

		//
		event := domain.MusicEvent{
			Context: "https://schema.org",
			Type:    "MusicEvent",
			Name:    t,
			Location: domain.Location{
				Type: "MucisVenue",
				Name: "Richter",
				Address: domain.PostalAddress{
					Type:       "PostalAddress",
					Street:     "",
					Locality:   "",
					PostalCode: "",
					Country:    "",
				},
			},
		}
		//
		events = append(events, event)
	})

	//
	err := se.collector.Visit(url)
	se.collector.Wait()

	return events, err
}
