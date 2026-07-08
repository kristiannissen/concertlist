# ConcertList: ETL Concert Scraper

> **Queue-Driven Crawling** for scalable, resilient concert data extraction from music venue websites.

---

## 🎯 Key Design Idea

**Queue-Driven Crawling** – Each URL (index page, sub-pages like `/event-1`, `/event-2`) is processed as an independent queue job. This enables:

- **Automatic depth handling** (index → sub-pages → sub-sub-pages)
- **Parallel processing** via Vercel Queue workers
- **Resilience** with built-in retries for failed jobs
- **Statelessness** – no crawler state to manage

---

## 🏗️ Architecture: Hexagonal (Ports &amp; Adapters)

```
┌─────────────────────────────────────────────────────────────┐
│                        API / CLI                               │
│                  (Driving Adapters)                           │
└─────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────┐
│                      Domain Layer                              │
│  ┌─────────────┐    ┌─────────────┐    ┌───────────────────┐ │
│  │  Concert    │    │   Scraper   │    │    ScrapeResult   │ │
│  │   Model    │    │  Interface  │    │  (Concerts+Links) │ │
│  └─────────────┘    └─────────────┘    └───────────────────┘ │
│       ▲                  ▲                  ▲                │
└───────┼──────────────────┼──────────────────┼────────────────┘
        │                  │                  │
        ▼                  ▼                  ▼
┌─────────────────────────────────────────────────────────────┐
│                   Driven Adapters                              │
│  ┌───────────────────┐  ┌───────────────────┐  ┌─────────────┐ │
│  │  Vercel Blob       │  │  Vercel Queue      │  │   HTTP      │ │
│  │  (Storage)         │  │  (Processing)      │  │  (Fetching) │ │
│  └───────────────────┘  └───────────────────┘  └─────────────┘ │
└─────────────────────────────────────────────────────────────┘
```

- **Domain Layer** (`internal/domain`): Pure business logic, no external dependencies
- **Ports** (`internal/ports`): Interfaces defining what the app needs (Storage, Queue)
- **Adapters** (`internal/adapters`): Implementations for Vercel Blob, Vercel Queue, HTTP

---

## 📁 Folder Structure

```
.
├── api/
│   └── index.go              # Vercel gateway (exports `Handler`)
├── cmd/
│   ├── api/
│   │   └── main.go           # Local API server
│   └── cli/
│       └── main.go           # CLI entry point
├── internal/
│   ├── domain/
│   │   ├── concert.go        # Core models (Concert, ScrapeJob, ScrapeResult)
│   │   ├── fetcher.go        # HTTP client
│   │   └── scrapers/
│   │       ├── registry.go   # Scraper registry (maps URLs → scrapers)
│   │       └── venue_a.go    # Site-specific scrapers (HTML/JSON/CSV/XML)
│   ├── ports/
│   │   ├── storage.go        # Storage interface
│   │   └── queue.go          # Queue interface
│   └── adapters/
│       ├── http/
│       │   └── handler.go     # API routes (/scrape, /queue)
│       ├── vercelblob/
│       │   └── storage.go     # Vercel Blob storage adapter
│       └── vercelqueue/
│           └── queue.go      # Vercel Queue adapter
├── go.mod
├── go.sum
└── vercel.json               # Vercel config (rewrites, crons, queues)
```

---

## 🔄 ETL Flow

```
1. [Cron/CLI] → Enqueue index page URL
   Queue: [ "/events" ]

2. [Worker] → Scrape "/events"
   → Extracts: "/event-1", "/event-2", "/event-3"
   → Saves: [] (index page has no concerts)
   → Enqueues: "/event-1", "/event-2", "/event-3"
   Queue: [ "/event-1", "/event-2", "/event-3" ]

3. [Worker] → Scrape "/event-1"
   → Extracts: [] (no further links)
   → Saves: [ConcertA, ConcertB]
   Queue: [ "/event-2", "/event-3" ]

4. [Worker] → Scrape "/event-2"
   → Saves: [ConcertC]
   Queue: [ "/event-3" ]

5. [Worker] → Scrape "/event-3"
   → Saves: [ConcertD, ConcertE]
   Queue: [ ]
```

---

## 🛠️ Core Components

### Scraper Interface

```go
type Scraper interface {
    Match(url string) bool
    Parse(ctx context.Context, data []byte, contentType string) (ScrapeResult, error)
}

type ScrapeResult struct {
    Concerts []Concert   // Extracted concerts
    NextURLs []string    // Sub-pages to scrape next
}
```

### Queue Processing

- **Vercel Queue** handles job distribution and retries
- Each URL = one job
- Workers process jobs concurrently
- Failed jobs auto-retry (configurable in `vercel.json`)

### Storage

- **Vercel Blob** for persistent concert data
- Simple key-value storage
- No database needed

---

## 🚀 Deployment

### Vercel Configuration (`vercel.json`)

```json
{
  "$schema": "https://openapi.vercel.sh/vercel.json",
  "cleanUrls": true,
  "rewrites": [{"source": "/api/(.*)", "destination": "api/index.go"}],
  "crons": [{"path": "/api/scrape", "schedule": "0 0 * * *"}],
  "queues": {
    "scrape-queue": {"memory": 300, "timeout": 60, "maxRetries": 3}
  }
}
```

### Local Development

```bash
# Run API server
go run ./cmd/api

# Run CLI scraper
go run ./cmd/cli --url https://venue.com/events
```

---

## 🎭 Supported Formats

Each scraper handles its own format:

- **HTML** (via Gocolly)
- **JSON** (via `encoding/json`)
- **XML** (via `encoding/xml`)
- **CSV** (via `encoding/csv`)

Add new venues by implementing the `Scraper` interface and registering in `registry.go`.

---

## 📌 Design Principles

1. **Keep It Simple** – No DI frameworks, no enterprise patterns
2. **Idiomatic Go** – `context.Context`, `error` returns, channels for concurrency
3. **Hexagonal Architecture** – Domain isolated from infrastructure
4. **Queue-Driven** – Scalable, resilient, stateless crawling
5. **Vercel-Native** – Uses Vercel Cron, Queue, Blob for serverless deployment
