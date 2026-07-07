//
//
package extractor

import (
    "time"

    "github.com/kristiannissen/concertlist/internal/domain"
    "github.com/kristiannissen/concertlist/internal/ports"
    "github.com/gocolly/colly/v2"
)

//
type EventProcessor func(e *colly.HTMLElement) domain.MusicEvent

//
type CollyExtractor struct {
    collector *colly.Collector
    processor EventProcessor
    queue ports.Queue
}

//
type Option func(*CollyExtractor)

//
func WithQueue(q ports.Queue) Option {
    return func(e *CollyExtractor) {
        e.queue = q
    }
}

//
func NewCollyExtractor(opts ...Option) ports.Extractor {
    ce := &CollyExtractor{
        collector: colly.NewCollector(
            colly.Async(true),
        ),
    }

    //
    ce.collector.Limit(&colly.LimitRule{
        DomainGlob: "*",
        Parallelism: 2,
        Delay: 2 * time.Second,
    })

    for _, opt := range opts {
        opt(ce)
    }
    return ce
}

//
func (e *CollyExtractor) Extract(url string, selectors []string) error {
    for _, selector := range selectors {
        e.collector.OnHTML(selector, func(h *colly.HTMLElement) {
            if e.processor != nil {
                event := e.processor(h)
                _ = e.queue.Push(event) // TODO: handle _
            }
        })
    }

    err := e.collector.Visit(url)
    e.collector.Wait()

    return err
}
