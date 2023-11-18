package boilerplate

import (
	d "concertlist/internal/dataaccess"
	"errors"
)

/**
 * Should implement an interface
 */
type Extractor interface {
	Run() ([]d.Event, error)
}

type Resource struct {
	URL string
}

/**
 * Runs the OnHTML methods to collect links and create
 * []struct containing data object
 */
func (r Resource) Run() ([]d.Event, error) {
	var events []d.Event
	events = append(events, d.Event{
		Artist: "Greta Van Fleet",
		Venue:  "Forum",
	})

	if len(events) == 0 {
		return nil, errors.New("No Events found")
	}

	return events, nil
}
