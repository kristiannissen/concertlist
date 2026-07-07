//
//
package ports

//
//
type Extractor interface {
    Extract(url string, selctors []string) error
}
