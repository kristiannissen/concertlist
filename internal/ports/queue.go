package ports

import "github.com/kristiannissen/concertlist/internal/domain"

type Queue interface {
	Push(event domain.MusicEvent) error
}
