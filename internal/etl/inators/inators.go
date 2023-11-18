package inators

/**
 * Inators take care of all the extract and transform tasks
 * shared across all the custom modules using an interface
 * Each scraper should implement the interface and in the etl
 * the different scrapers will be executed
 * https://gobyexample.com/interfaces
 */

type Extractor interface {
	Run()
}

// Shared struct
// TODO: Should return struct or error
func Runner(e interface{ Extractor }) {
	// fmt.Println("URL ", r.URL)
	e.Run()
}
