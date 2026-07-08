//
//
package etl

import (
    "github.com/kristiannissen/concertlist/internal/ports"
    "fmt"
)

// TODO: add queue
type Orchestrator struct {
    extractor ports.Extractor
}

//
func NewOrchestrator(e ports.Extractor) *Orchestrator {
    return &Orchestrator{
        extractor: e,
    }
}

//
func (o *Orchestrator) RunCLI(url string) {
    //
    events, _ := o.extractor.Extract(url)
    for event := range events {
        fmt.Println(event)
    }
}

//
func (o *Orchestrator) RunHTTP(url string) {
}
