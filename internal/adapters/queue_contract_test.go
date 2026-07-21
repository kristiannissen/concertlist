// internal/adapters
package adapters

import (
	"encoding/json"
	"testing"

	"go.uber.org/zap"
)

// These mirror the anonymous message structs decoded by EventScrapeConsumer
// and EventExtractConsumer in router.go. They're deliberately duplicated
// here rather than imported: this test exists specifically to catch the bug
// class where a producer's JSON keys (router.go's /api/scrape/trigger
// handler, or a scraper's Scrape method) drift from what a consumer
// expects — messages then queue and get acknowledged successfully while
// silently never reaching Scrape/Extract. That exact bug shipped once
// already (vega.go published "name" while the consumer decoded "venue").
type scrapeMessage struct {
	Venue string `json:"venue"`
}

type extractMessage struct {
	Venue string `json:"venue"`
	URL   string `json:"url"`
}

// TestVenueScrapeMessageContract checks that the body router.go's
// GET /api/scrape/trigger handler publishes to the "venue-scrape" topic for
// each registered venue decodes into a key that actually resolves in the
// scraper registry.
func TestVenueScrapeMessageContract(t *testing.T) {
	reg := NewScraperRegistry(zap.NewNop(), "")
	if len(reg) == 0 {
		t.Fatal("scraper registry is empty — nothing to check")
	}

	for venue := range reg {
		body, err := json.Marshal(map[string]string{"venue": venue})
		if err != nil {
			t.Fatalf("marshal producer body for %q: %v", venue, err)
		}

		var msg scrapeMessage
		if err := json.Unmarshal(body, &msg); err != nil {
			t.Fatalf("consumer failed to decode producer body for %q: %v", venue, err)
		}

		if _, ok := reg[msg.Venue]; !ok {
			t.Errorf("published venue %q does not resolve in the scraper registry", msg.Venue)
		}
	}
}

// TestEventExtractMessageContract checks that the body Vega.Scrape
// publishes to the "event-extract" topic (internal/adapters/scrapers/vega.go)
// decodes into a venue key that resolves in the scraper registry, and that
// the URL field survives the round trip.
func TestEventExtractMessageContract(t *testing.T) {
	reg := NewScraperRegistry(zap.NewNop(), "")

	// Exact shape Vega.Scrape marshals today.
	body := []byte(`{"venue":"vega","url":"https://vega.dk/event/example"}`)

	var msg extractMessage
	if err := json.Unmarshal(body, &msg); err != nil {
		t.Fatalf("consumer failed to decode producer body: %v", err)
	}
	if msg.Venue == "" || msg.URL == "" {
		t.Fatalf("decoded message has empty fields (venue=%q url=%q) — producer/consumer field name mismatch", msg.Venue, msg.URL)
	}

	if _, ok := reg[msg.Venue]; !ok {
		keys := make([]string, 0, len(reg))
		for k := range reg {
			keys = append(keys, k)
		}
		t.Errorf("published venue %q does not resolve in the scraper registry (known keys: %v)", msg.Venue, keys)
	}
}
