# Concert List

## Commands

```
go run ./cmd/api
```
Runs the HTTP server locally (same entry point Vercel's Go runtime uses via `api/index.go`).

```
go run ./cmd/cron
```
Runs the queue producer/consumer job. Controlled by `QUEUE_MODE=producer|consumer` (defaults to `producer`).

Note: no venue extractors are implemented yet, so the producer currently has no venues to enqueue jobs for.


### Structure
```
.
├── api
│   └── index.go              # Vercel serverless entry point, delegates to internal/router
├── cmd
│   ├── api
│   │   └── main.go            # local dev server entry point
│   └── cron
│       └── main.go            # queue producer/consumer entry point
├── internal
│   ├── adapters
│   │   ├── queue
│   │   │   ├── consumer.go
│   │   │   ├── queue.go        # Vercel Queue client
│   │   │   └── types.go
│   │   └── storage
│   │       └── blob            # Vercel Blob storage client
│   ├── domain                  # core models, ports, ETL service
│   ├── handler                 # HTTP handlers
│   └── router                  # HTTP routing
├── LICENSE
├── README.md
├── go.mod
└── vercel.json

11 directories, 13 files
```
