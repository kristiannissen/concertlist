# Concert List

## Commands

```
go run ./internal/etl
```
Runs extract, transform and load

```
go run ./cmd/cli
```
Runs command-line


### Structure
.
├── cmd
│   └── cli
│       └── main.go
├── internal
│   ├── etl
│   │   ├── richtergladsaxe
│   │   │   └── richtergladsaxe.go
│   │   └── extract.go
│   └── random
│       └── number.go
├── LICENSE
├── README.md
└── go.mod

6 directories, 7 files
