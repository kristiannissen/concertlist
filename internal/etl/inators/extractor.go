// Package inators takes offers a suite of methods
// to run the custom adaptors in a comon way
package inators

import (
	d "concertlist/internal/dataaccess"
	"log"

	"github.com/gocolly/colly"
)

// An Extractor needs to be implemented by the
// custom adaptors
type Extractor interface {
	Run() ([]d.Event, error)
	New(c *colly.Collector) error
}

// Runner takes an interface extractor and runs
// the methods
func Runner(e interface{ Extractor }) {
	// fmt.Println("URL ", r.URL)
	var events []d.Event

	events, _ = e.Run()

	for _, evnt := range events {
		_, err := evnt.UpSert()
		if err != nil {
			log.Fatal(err)
		}
	}
}
