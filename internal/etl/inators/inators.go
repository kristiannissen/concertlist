package inators

/**
 * Inators take care of all the extract and transform tasks
 * shared across all the custom modules using an interface
 * Each scraper should implement the interface and in the etl
 * the different scrapers will be executed
 * https://gobyexample.com/interfaces
 */

// Shared struct
type Resource struct {
	URL string
}
type Extractor interface {
	New()
}

func Run(e Extractor) {
	e.New()
}
