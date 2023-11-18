package boilerplate

import (
	"fmt"
)

/**
 * Should implement an interface
 */
type Extractor interface {
	Run() string
}

type Resource struct {
	URL string
}

// TODO: should return resource or error
func (r Resource) Run() {
	fmt.Println("Boilerplate Run: ", r.URL)
}
