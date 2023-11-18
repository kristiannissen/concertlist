package boilerplate

import "fmt"

/**
 * Should implement an interface
 */
type Extractor interface {
	Run() string
}

type Resource struct {
	URL string
}

func (r Resource) New(e interface{ Extractor }) {
	fmt.Println("URL ", r.URL)
	e.Run()
}
