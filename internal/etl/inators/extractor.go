package inators

import (
	d "concertlist/internal/dataaccess"
	"log"
)

/**
 * Inators take care of all the extract and transform tasks
 * shared across all the custom modules using an interface
 * Each scraper should implement the interface and in the etl
 * the different scrapers will be executed
 * https://gobyexample.com/interfaces
 *
 */

type Extractor interface {
	Run() ([]d.Event, error)
}

func Runner(e interface{ Extractor }) {
	// fmt.Println("URL ", r.URL)
	var events []d.Event
	events, _ = e.Run()

	for k, v := range events {
		log.Println("event: ", k, v)
	}
}
