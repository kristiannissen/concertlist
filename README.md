# ConcertList: Queue-Driven Concert Scraper

> **Scalable, resilient concert data extraction** from music venue websites using queue-driven processing.

---

## 🎯 Key Design Idea

**Queue-Driven Crawling** – Each venue's scraping job is processed as an independent queue message. This enables:

- **Automatic parallel processing** via Vercel Queue workers
- **Resilience** with built-in retries for failed jobs
- **Scalability** – add more scrapers without changing core infrastructure
- **Decoupling** – scrapers publish to queue, consumers process asynchronously

---

## 🏗️ Architecture: Hexagonal (Ports & Adapters)

```
┌─────────────────────────────────────────────────────────────────────────┐
│                        Entry Points                                    │
│  ┌─────────────────┐  ┌─────────────────┐  ┌─────────────────┐          │
│  │  api/index.go    │  │  queue-consumer  │  │  cmd/api/main.go │          │
│  │  (Vercel HTTP)   │  │  (Vercel Queue)  │  │  (Local Server)  │          │
│  └────────┬────────┘  └────────┬────────┘  └────────┬────────┘          │
└───────────┼──────────────────────┼──────────────────────┼──────────────────┘
            │                      │                      │
            ▼                      ▼                      ▼
┌─────────────────────────────────────────────────────────────────────────┐
│                        Gateway Layer                                   │
│  ┌─────────────────────────────────────────────────────────────────────┐ │
│  │  gateway/gateway.go                                              │ │
│  │  - Bridges Vercel entry points to internal packages               │ │
│  │  - Works around Vercel's internal package import restrictions    │ │
│  └─────────────────────────────────────────────────────────────────────┘ │
└─────────────────────────────────────────────────────────────────────────┘
                            │
                            ▼
┌─────────────────────────────────────────────────────────────────────────┐
│                      Domain Layer (internal/domain/)                      │
│  ┌─────────────────────────────────────────────────────────────────────┐ │
│  │  model.go                                                         │ │
│  │  - MusicEvent (schema.org structured data)                        │ │
│  │  - Location, PostalAddress, Performer, Offer                       │ │
│  └─────────────────────────────────────────────────────────────────────┘ │
└─────────────────────────────────────────────────────────────────────────┘
                            │
                            ▼
┌─────────────────────────────────────────────────────────────────────────┐
│                      Ports (internal/ports/)                             │
│  ┌─────────────────────────────────────────────────────────────────────┐ │
│  │  scraper.go                                                       │ │
│  │  - Scraper interface: Scrape(ctx, wg) and Extract(ctx, wg)         │ │
│  └─────────────────────────────────────────────────────────────────────┘ │
└─────────────────────────────────────────────────────────────────────────┘
                            │
                            ▼
┌─────────────────────────────────────────────────────────────────────────┐
│                   Driven Adapters (internal/adapters/)                   │
│  ┌─────────────────────┐  ┌───────────────────────────────────────────┐ │
│  │  router.go           │  │  scrapers/                                      │ │
│  │  - HTTP routing      │  │    └── richter.go (Richter Gladsaxe)       │ │
│  │  - API endpoints     │  │    └── [future scrapers...]                │ │
│  │  - QueueConsumer     │  │                                                │ │
│  └─────────────────────┘  └───────────────────────────────────────────┘ │
└─────────────────────────────────────────────────────────────────────────┘
```

### Layer Responsibilities

| Layer | Path | Responsibility |
|-------|------|-----------------|
| **Entry Points** | `api/`, `cmd/` | Vercel HTTP, Queue, CLI entry points |
| **Gateway** | `gateway/` | Bridges Vercel to internal packages |
| **Domain** | `internal/domain/` | Core models (MusicEvent, Location, etc.) |
| **Ports** | `internal/ports/` | Interfaces (Scraper) |
| **Adapters** | `internal/adapters/` | Implementations (router, scrapers) |

---

## 📁 Folder Structure

```
.
├── api/
│   ├── index.go              # Vercel HTTP gateway (exports Handler)
│   └── queue-consumer/
│       └── index.go          # Vercel Queue consumer (exports Handler)
├── cmd/
│   ├── api/
│   │   └── main.go           # Local API server entry point
│   └── cli/
│       └── main.go           # CLI entry point for manual scraping
├── gateway/
│   └── gateway.go            # Bridges api/ to internal/ packages
├── internal/
│   ├── domain/
│   │   └── model.go          # Core data models (MusicEvent, etc.)
│   ├── ports/
│   │   └── scraper.go        # Scraper interface
│   └── adapters/
│       ├── router.go         # HTTP router with endpoints
│       └── scrapers/
│           └── richter.go     # Richter Gladsaxe venue scraper
│           # └── [add more venue scrapers here]
├── public/                   # Static files for Vercel output
├── go.mod
├── go.sum
└── vercel.json               # Vercel configuration
```

---

## 🔄 Data Flow

### Cron-Triggered Scraping

```
1. [Vercel Cron] → GET /api/musicevent/richter (daily at 00:00 UTC)
   
2. [Router] → Creates Richter scraper instance
   
3. [Richter Scraper] → Uses Colly to crawl richter-gladsaxe.dk
   - Follows links with MaxDepth(2)
   - Extracts concert data from .single-concert elements
   - Posts each MusicEvent to Vercel Queue topic "musicevent"
   
4. [Vercel Queue] → Buffers messages for async processing
   
5. [Queue Consumer] → api/queue-consumer/index.go processes messages
   - Currently logs receipt (to be extended)
```

### CLI Scraping

```bash
# Run Richter scraper manually
go run ./cmd/cli

# Output: Logs scraping progress and events found
```

---

## 🛠️ Core Components

### Scraper Interface (`internal/ports/scraper.go`)

```go
type Scraper interface {
    Scrape(ctx context.Context, wg *sync.WaitGroup) error
    Extract(ctx context.Context, wg *sync.WaitGroup) error
}
```

### Richter Scraper Implementation

The Richter scraper (`internal/adapters/scrapers/richter.go`) uses:

- **Colly** for HTML crawling with:
  - `AllowedDomains` to restrict to richter-gladsaxe.dk
  - `MaxDepth(2)` for controlled crawling
  - `Async(true)` with `Parallelism: 2` for concurrent requests
  - `RandomDelay: 5s` to avoid overwhelming the site
  - Custom `sync.Map` for thread-safe visited URL tracking

- **Resty** for posting scraped events to Vercel Queue

- **Zap** for structured logging

### Domain Model (`internal/domain/model.go`)

```go
type MusicEvent struct {
    Context   string    `json:"@context"`
    Type      string    `json:"@type"`
    Name      string    `json:"name"`
    StartDate string    `json:"startDate"`
    Location  Location  `json:"location,omitempty"`
    Performer Performer `json:"performer,omitempty"`
    Offer     Offer     `json:"offers,omitempty"`
}

type Location struct {
    Type    string        `json:"@type"`
    Name    string        `json:"name"`
    Address PostalAddress `json:"address"`
}

// Plus: PostalAddress, Performer, Offer
```

---

## 📡 Queue Infrastructure

### Vercel Queue Configuration

The queue is configured in `vercel.json`:

```json
{
  "functions": {
    "api/queue-consumer/index.go": {
      "experimentalTriggers": [
        {
          "type": "queue/v2beta",
          "topic": "musicevent",
          "retryAfterSeconds": 60,
          "initialDelaySeconds": 0
        }
      ]
    }
  }
}
```

### Queue Topic: `musicevent`

- **Producer**: Richter scraper posts MusicEvent JSON to queue
- **Consumer**: `api/queue-consumer/index.go` processes messages
- **Retry**: Failed messages retry after 60 seconds
- **Scaling**: Vercel automatically scales queue consumers

---

## ⏰ Cron Jobs

Configured in `vercel.json`:

```json
{
  "crons": [
    {
      "path": "/api/musicevent/richter",
      "schedule": "0 0 * * *"
    }
  ]
}
```

This triggers the Richter scraper daily at midnight UTC.

---

## 🚀 Adding New Scrapers

To add a new venue scraper:

### 1. Create Scraper File

Add a new file in `internal/adapters/scrapers/`:

```go
// internal/adapters/scrapers/newvenue.go
package scrapers

import (
    "context"
    "sync"
    
    "github.com/gocolly/colly"
    "github.com/kristiannissen/concertlist/internal/domain"
    "go.uber.org/zap"
)

type NewVenue struct {
    URL string
    Log *zap.Logger
    visited sync.Map
}

func (n *NewVenue) Scrape(ctx context.Context, wg *sync.WaitGroup) error {
    // Implement crawling logic using Colly
    // Extract concert data and post to queue
    return nil
}

func (n *NewVenue) Extract(ctx context.Context, wg *sync.WaitGroup) error {
    // Optional: Implement data extraction logic
    return nil
}
```

### 2. Add API Endpoint

Add a new endpoint in `internal/adapters/router.go`:

```go
mux.HandleFunc("GET /api/musicevent/newvenue", func(w http.ResponseWriter, r *http.Request) {
    var scraper ports.Scraper = &scrapers.NewVenue{
        URL: "https://newvenue.com/events",
        Log: logger,
    }
    wg := &sync.WaitGroup{}
    err := scraper.Scrape(r.Context(), wg)
    wg.Wait()
    
    w.Header().Set("Content-Type", "application/json")
    if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        json.NewEncoder(w).Encode(map[string]string{"status": err.Error()})
        return
    }
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
})
```

### 3. Add Cron Job (Optional)

Add to `vercel.json`:

```json
{
  "crons": [
    {
      "path": "/api/musicevent/richter",
      "schedule": "0 0 * * *"
    },
    {
      "path": "/api/musicevent/newvenue",
      "schedule": "0 6 * * *"
    }
  ]
}
```

---

## 🎭 Supported Formats

Scrapers can handle various data formats:

- **HTML** (via [Colly](https://github.com/gocolly/colly)) – Primary format for venue websites
- **JSON** (via `encoding/json`) – For APIs returning JSON
- **XML** (via `encoding/xml`) – For XML-based feeds
- **CSV** (via `encoding/csv`) – For CSV exports

Each scraper implementation chooses the appropriate parsing method for its target.

---

## 🛠️ Technologies

| Component | Technology | Purpose |
|-----------|------------|---------|
| **Web Framework** | `net/http` | HTTP routing and serving |
| **Crawling** | [Colly](https://github.com/gocolly/colly) | HTML scraping with concurrency |
| **HTTP Client** | [Resty](https://github.com/go-resty/resty) | Posting to Vercel Queue |
| **Logging** | [Zap](https://github.com/uber-go/zap) | Structured logging |
| **Concurrency** | `sync.WaitGroup`, `context.Context` | Goroutine coordination |
| **Queue** | Vercel Queue v2beta | Async message processing |
| **Deployment** | Vercel | Serverless hosting |

---

## 🚀 Deployment

### Vercel Configuration (`vercel.json`)

```json
{
   "$schema":"https://openapi.vercel.sh/vercel.json",
   "framework":null,
   "outputDirectory":"public",
   "regions":["arn1"],
   "rewrites":[
      {
         "source":"/api/(.*)",
         "destination":"/api"
      }
   ],
   "build":{
      "env":{
         "GO_BUILD_FLAGS":"-ldflags '-s -w'"
      }
   },
   "crons":[
      {
         "path":"/api/musicevent/richter",
         "schedule":"0 0 * * *"
      }
   ],
   "functions":{
      "api/queue-consumer/index.go":{
         "experimentalTriggers":[
            {
               "type":"queue/v2beta",
               "topic":"musicevent",
               "retryAfterSeconds":60,
               "initialDelaySeconds":0
            }
         ]
      }
   }
}
```

### Local Development

```bash
# Install dependencies
go mod tidy

# Run local API server
go run ./cmd/api

# Run CLI scraper manually
go run ./cmd/cli

# Test API endpoints
curl http://localhost:3000/api/health
curl http://localhost:3000/api/musicevent/richter
```

The API server listens on port 3000 by default, or on the `PORT` environment variable if set.

---

## 📌 Design Principles

1. **Hexagonal Architecture** – Domain isolated from infrastructure, easy to test and swap implementations
2. **Queue-Driven** – Scalable, resilient, decoupled processing
3. **Idiomatic Go** – `context.Context`, `error` returns, `sync.WaitGroup` for concurrency
4. **Vercel-Native** – Uses Vercel Cron, Queue, and serverless functions
5. **Keep It Simple** – No DI frameworks, minimal abstractions
6. **Extensible** – Add new scrapers by implementing the Scraper interface

---

## 🔧 Environment Variables

| Variable | Purpose | Required |
|----------|---------|----------|
| `VERCEL_OIDC_TOKEN` | Authentication for Vercel Queue API | Yes (for queue posting) |
| `PORT` | Local server port | No (defaults to 3000) |
| `HOST` | Local server host | No (defaults to localhost) |

---

## 📊 Current Status

- ✅ **Richter Gladsaxe** scraper implemented and deployed
- ✅ **Vercel Queue** integration for async processing
- ✅ **Cron job** for daily scraping
- ✅ **CLI** for manual scraping
- 🚧 **Queue consumer** – Currently logs messages (processing logic to be added)
- 🚧 **Additional scrapers** – Ready to be added to `internal/adapters/scrapers/`

---

## 🎯 Roadmap

1. **Enhance Queue Consumer** – Implement message processing and storage
2. **Add More Scrapers** – Extend to additional venues
3. **Add Storage Adapter** – Persist scraped events (Vercel Blob, database, etc.)
4. **Add API Endpoints** – Query scraped concert data
5. **Improve Error Handling** – Better retry logic and dead-letter queue
6. **Add Monitoring** – Metrics and alerts for scraping jobs
