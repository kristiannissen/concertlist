# ConcertList ETL Pipeline

Dette projekt er bygget ved hjælp af **Hexagonal Arkitektur** (Ports & Adapters) for at sikre høj modularitet og adskillelse af forretningslogik fra ekstern infrastruktur.

## Arkitektur-oversigt

### 1. Domain (`internal/domain`)
Det inderste lag. Indeholder rene data-strukturer (f.eks. `MusicEvent`). Dette lag er uafhængigt af alt andet og kender ikke til databaser eller scrapere.

### 2. Ports (`internal/ports`)
Definerer de interfaces, som applikationen har brug for. Det er kontrakten for, hvordan man "trækker data ud" (Extractor) eller "gemmer data" (Queue).

### 3. UseCase (`internal/usecase`)
Det centrale orkestreringslag. Her placeres forretningslogikken.
- `orchestrator.go`: Modtager en `Extractor` og en `Queue`. Den kører flowet: *Hent data -> Loop -> Push til kø*. 
- Orkestratoren er generisk og kender ikke til de specifikke scraping-teknikker.

### 4. Adapters (`internal/adapters`)
De ydre lag, der implementerer portene:
- `extractor`: Specifikke implementeringer pr. hjemmeside (HTML/JSON/XML).
- `queue`: Implementering af kø-logik.
- `handler`: HTTP-endpoints der modtager eksterne kald (fra Cron/Web).

## Workflow

1. **Entrypoint** (`cmd/`): En main-fil initialiserer den ønskede adapter (f.eks. `SiteAExtractor`) og den generiske `Orchestrator`.
2. **Execution**: `Orchestrator` kalder `ExtractAll()` på adapteren.
3. **Processing**: Data returneres som `[]domain.MusicEvent` til orkestratoren, som pusher dem til køen.

## Hvordan tilføjes et nyt site?
1. Opret en ny fil i `internal/adapters/extractor/` (f.eks. `site_b.go`).
2. Implementér `ports.Extractor` interfacet (metoden `ExtractAll()`).
3. Opdater din `main.go` i `cmd/` til at bruge den nye adapter via dependency injection.

## Mappe struktur
```
├── cmd/
│   ├── api/              # Entrypoint til HTTP (Cron/API)
│   ├── cli/              # Entrypoint til CLI
│   └── cron/             # Entrypoint til Cron-jobs
├── internal/
│   ├── adapters/         # Eksterne implementeringer (infrastructure)
│   │   ├── extractor/    # Site-specifikke scrapere (adaptere)
│   │   │   ├── site_a.go # Implementerer ports.Extractor
│   │   │   └── site_b.go # Implementerer ports.Extractor
│   │   ├── queue/        # Kø-implementering
│   │   ├── storage/      # Blob/Database-implementering
│   │   ├── handler/      # HTTP handlers
│   │   └── router/       # Routing-konfiguration
│   ├── domain/           # Domæne-modeller (MusicEvent)
│   ├── ports/            # Interfaces (kontrakter)
│   │   ├── extractor.go
│   │   ├── orchestrator.go
│   │   └── queue.go
│   └── usecase/          # Forretningslogik
│       └── etl/
│           └── orchestrator.go
├── go.mod
└── go.sum
```

