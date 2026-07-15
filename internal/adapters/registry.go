// internal/adapters
package adapters

import (
	"github.com/kristiannissen/concertlist/internal/adapters/scrapers"
	"github.com/kristiannissen/concertlist/internal/ports"
	"go.uber.org/zap"
)

func NewScraperRegistry(log *zap.Logger) map[string]ports.Scraper {
	return map[string]ports.Scraper{
		"ricter": &scrapers.Richter{URL: "https://richter-gladsaxe.dk", Log: log},
		"vega":   &scrapers.Vega{URL: "https://vega.dk", Log: log},
	}
}
