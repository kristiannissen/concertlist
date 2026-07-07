//
//
package ports

import "github.com/kristiannissen/concertlist/internal/domain"

//
//
type Extractor interface {
    Extract() ([]domain.MusicEvent, error)  
}
