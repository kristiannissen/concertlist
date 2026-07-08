//
// 
package ports

import (
    "github.com/kristiannissen/concertlist/internal/domain"
)

//
type Orchestrator interface {
    RunCLI(url string) ([]domain.MusicEvent, error)
    RunHTTP(url string) ([]domain.MusicEvent, error)
}
