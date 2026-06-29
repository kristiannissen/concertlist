// Package extractors provides site-specific data extraction implementations.
package extractors

import (
	"github.com/gocolly/colly/v2"
)

// BaseExtractor provides common functionality for all site extractors.
type BaseExtractor struct {
	collector *colly.Collector
}

// NewBaseExtractor creates a new base extractor with default colly settings.
func NewBaseExtractor() *BaseExtractor {
	return &BaseExtractor{
		collector: colly.NewCollector(
			colly.AllowedDomains(),
			colly.MaxDepth(2),
			colly.Async(true),
		),
	}
}

// Collector returns the colly collector instance.
func (e *BaseExtractor) Collector() *colly.Collector {
	return e.collector
}

// Clone creates a new colly collector with the same settings as the base.
// This allows each site extractor to have its own isolated collector.
func (e *BaseExtractor) Clone() *colly.Collector {
	return e.collector.Clone()
}
