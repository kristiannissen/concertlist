// Package ports defines interfaces for the application's ports (inbound and outbound).
// This package contains the contract that adapters must implement.
package ports

import (
	"context"

	"github.com/kristiannissen/concertlist/internal/domain"
)

// ConcertService defines the interface for concert business logic operations.
// This is an inbound (driving) port that the domain layer will use.
type ConcertService interface {
	// GetAllConcerts returns all concerts.
	GetAllConcerts(ctx context.Context) ([]domain.Concert, error)

	// GetConcertByID returns a single concert by its ID.
	GetConcertByID(ctx context.Context, id string) (*domain.Concert, error)

	// CreateConcert creates a new concert.
	CreateConcert(ctx context.Context, req domain.CreateConcertRequest) (*domain.Concert, error)

	// UpdateConcert updates an existing concert.
	UpdateConcert(ctx context.Context, id string, req domain.CreateConcertRequest) (*domain.Concert, error)

	// DeleteConcert deletes a concert.
	DeleteConcert(ctx context.Context, id string) error
}
