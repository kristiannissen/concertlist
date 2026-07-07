//
//
package ports

//
//
type Extractor interface {
    Extract(url, selctor string) error
}
