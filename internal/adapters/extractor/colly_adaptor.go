package extractor

import (
	"time"

	"github.com/gocolly/colly/v2"
	"github.com/kristiannissen/concertlist/internal/domain"
	"github.com/kristiannissen/concertlist/internal/ports"
)

type CollyExtractor struct {
	collector *colly.Collector
}

type Option func(*CollyExtractor)

func NewCollyExtractor(opts ...Option) ports.Extractor {
	ce := &CollyExtractor{
		collector: colly.NewCollector(
			colly.Async(true),
		),
	}

	//
	ce.collector.Limit(&colly.LimitRule{
		DomainGlob:  "*",
		Parallelism: 2,
		Delay:       2 * time.Second,
	})

	for _, opt := range opts {
		opt(ce)
	}
	return ce
}

func (e *CollyExtractor) Extract(url string) ([]domain.MusicEvent, error) {
	var events []domain.MusicEvent

	err := e.collector.Visit(url)
	e.collector.Wait()

	return events, err
}
