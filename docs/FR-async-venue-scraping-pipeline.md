# Functional Requirement: Async, Queue-Driven Venue Scraping Pipeline

**Status:** Draft
**Author:** Kristian Nissen
**Date:** 2026-07-15
**Related files:** `internal/adapters/router.go`, `internal/adapters/scrapers/richter.go`, `internal/adapters/scrapers/vega.go`, `internal/ports/scraper.go`

## Problem Statement

Today, each venue is scraped by a dedicated cron-triggered HTTP route (e.g. `GET /api/musicevent/richter`) that synchronously crawls the venue's listing page, follows every event link, extracts event data, and posts each event to a queue — all inside one request/response cycle, blocked on with `wg.Wait()`.

This does not scale past a couple of venues: as more venues are added (6-8+), either each one needs its own cron route (unmanageable), or a single handler loops over all of them sequentially, which risks exceeding the serverless function's execution time limit. Long-running work inside an HTTP handler is also fragile on Vercel specifically — the runtime can freeze or tear down a function shortly after it responds, so any "fire a goroutine and return early" approach is not safe here; work that must survive past the response needs to be handed to something durable (a queue), not left in memory.

Separately, `Extract()` is already defined on `ports.Scraper` but unused — today's extraction logic lives inline inside `Scrape()`'s HTML callbacks instead.

## Goals

1. No single HTTP request performs more than one venue's worth of crawling.
2. Adding a venue does not require adding a new cron route or touching unrelated code.
3. A failure while processing one venue or one event does not require re-running the whole pipeline — only the failed unit is retried.
4. Extraction logic is decoupled from HTML/JSON/XML specifics at the point where it's dispatched, so the pipeline plumbing doesn't need per-venue special-casing.

## Non-Goals

- Per-venue queue topics or a dead-letter queue — the queue's default redelivery/retry is sufficient at this scale; do not build custom backoff logic.
- A persistent "already processed" datastore — Blob storage overwritten by a deterministic key already makes reprocessing idempotent.
- Migrating `Richter`/`Vega` off colly to a lighter HTML library — reasonable future cleanup, not required for this requirement.
- Onboarding venues beyond Richter and Vega — this requirement covers migrating the existing two to the new pipeline; adding further venues should require only a registry entry once this is done.

## Architecture Overview

Three stages, each a small, independently-retryable unit of work:

```
[cron, e.g. daily]
      |
      v
1. Trigger handler (HTTP, GET)
   - for each venue in the registry: publish {venue} to "venue-scrape" topic
   - returns as soon as all publishes succeed (no crawling here)
      |
      v
2. "venue-scrape" queue consumer  (one invocation per venue, run concurrently)
   - looks up the venue's ports.Scraper in the registry
   - calls Scrape(ctx, ...): crawls the venue's listing page(s) only,
     discovers event detail URLs (no field extraction here)
   - for each discovered URL: publish {venue, url} to "event-extract" topic
      |
      v
3. "event-extract" queue consumer  (one invocation per event URL)
   - looks up the venue's ports.Scraper in the registry
   - calls Extract(ctx, url) -> (domain.MusicEvent, error)
   - on success: writes the event to Blob, keyed by a slug derived from the
     event URL (idempotent — reprocessing overwrites, never duplicates)
   - on any error (extract or blob write): returns non-2xx so the queue
     redelivers the message
```

## Functional Requirements

### FR-1: Trigger handler becomes enqueue-only

The existing per-venue cron routes are replaced by a single handler that fans out to the `venue-scrape` topic instead of scraping inline.

Acceptance criteria:
- [ ] One cron schedule triggers one HTTP endpoint (e.g. `GET /api/scrape/trigger`)
- [ ] The endpoint publishes exactly one `{venue}` message per entry in the scraper registry
- [ ] The endpoint returns success once all publishes succeed, without waiting on any scraping or extraction
- [ ] No HTML/JSON/XML fetching happens inside this handler

### FR-2: `ports.Scraper` interface split

Replace the current single-purpose extraction-inside-`Scrape` design with two distinct responsibilities matching the pipeline stages.

- `Scrape(ctx context.Context, listingURL string) ([]string, error)` — crawls listing page(s) for one venue and returns discovered event URLs. Performs no field parsing and no queue/Blob I/O.
- `Extract(ctx context.Context, eventURL string) (domain.MusicEvent, error)` — fetches one event URL and returns the parsed event. Performs no queue/Blob I/O. Drops the current `*sync.WaitGroup` parameter — no longer needed since each call handles exactly one item.

Acceptance criteria:
- [ ] `Richter` and `Vega` both implement the updated interface
- [ ] Existing `.single-concert` field-parsing logic is moved from `Scrape`'s `OnHTML` callback into `Extract`, not duplicated
- [ ] `Scrape` no longer visits/parses individual event detail pages — it only follows listing-page pagination/link discovery needed to collect URLs

### FR-3: Venue registry

A single registry maps venue key to concrete `ports.Scraper`, so no consumer or handler hardcodes a concrete scraper type.

- Location: `internal/adapters/registry.go`
- Shape: `func NewScraperRegistry(log *zap.Logger) map[string]ports.Scraper`

Acceptance criteria:
- [ ] Both the trigger handler (FR-1) and both queue consumers (FR-4, FR-5) resolve their target scraper via this registry, keyed by the `venue` field on the incoming message
- [ ] Adding a new venue requires only a new entry in this map (plus its adapter implementation) — no changes to router or consumer logic

### FR-4: `venue-scrape` queue consumer

Acceptance criteria:
- [ ] Consumes messages of shape `{"venue": "richter"}`
- [ ] Looks up the venue in the registry; unknown venue is logged and the message is not retried (permanent failure, not transient)
- [ ] Calls `Scrape` for that venue only
- [ ] Publishes one `{"venue": ..., "url": ...}` message per discovered event URL to the `event-extract` topic
- [ ] Returns non-2xx on any error from `Scrape` itself, so the queue retries that venue only — not the others

### FR-5: `event-extract` queue consumer

Acceptance criteria:
- [ ] Consumes messages of shape `{"venue": "richter", "url": "https://richter-gladsaxe.dk/event/..."}`
- [ ] Looks up the venue in the registry and calls `Extract(ctx, url)`
- [ ] On success, serializes the returned `domain.MusicEvent` and writes it to Blob under a key derived deterministically from `url` (e.g. `events/{venue}/{slug}.json`)
- [ ] On extraction error or Blob write error, returns non-2xx so the queue redelivers this single message — re-processing does not require re-running venue discovery (FR-4)

### FR-6: Preserve existing scraping etiquette

Carried forward from the current implementation — must not regress during this refactor.

Acceptance criteria:
- [ ] All outbound requests set a descriptive `User-Agent` identifying the bot and a contact address (already added to `richter.go`/`vega.go`)
- [ ] `colly.IgnoreRobotsTxt()` is never called — robots.txt restrictions continue to apply
- [ ] Per-domain rate limiting (`colly.Limit` — parallelism/random delay) is preserved in `Scrape`

### FR-7: No silently discarded errors

Acceptance criteria:
- [ ] Every queue publish (venue-scrape, event-extract) checks its returned error and propagates/logs it — fixes the current bug where the event-post error is discarded, making failed posts indistinguishable from successful ones
- [ ] Every Blob write checks its returned error before the consumer responds 2xx

## Open Questions

- Final topic names — reuse the existing `musicevent` topic name for one of the two new stages, or introduce two new topic names (`venue-scrape`, `event-extract`)? *(Engineering)*
- Exact Blob key/slug convention — derive from URL path, or hash the URL? *(Engineering)*
- Queue concurrency/parallelism settings appropriate for the target request volume per venue. *(Engineering/Infra)*
- Should an unknown/misconfigured venue key in a message be surfaced anywhere (alert, dashboard) beyond a log line? *(Engineering)*
