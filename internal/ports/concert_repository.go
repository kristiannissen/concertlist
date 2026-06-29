// Package ports defines interfaces for the application's ports (inbound and outbound).
// This package contains the contract that adapters must implement.
package ports

import (
	"context"

	"github.com/kristiannissen/concertlist/internal/domain"
)

// ConcertRepository defines the interface for concert data operations.
// This is an outbound (driven) port that adapters will implement.
type ConcertRepository interface {
	// GetAll returns all concerts from the repository.
	GetAll(ctx context.Context) ([]domain.Concert, error)

	// GetByID returns a single concert by its ID.
	GetByID(ctx context.Context, id string) (*domain.Concert, error)

	// Create adds a new concert to the repository.
	Create(ctx context.Context, concert domain.CreateConcertRequest) (*domain.Concert, error)

	// Update modifies an existing concert.
	Update(ctx context.Context, id string, concert domain.CreateConcertRequest) (*domain.Concert, error)

	// Delete removes a concert from the repository.
	Delete(ctx context.Context, id string) error
}
