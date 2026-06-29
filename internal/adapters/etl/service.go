// Package etl provides the ETL (Extract, Transform, Load) workflow implementation.
package etl

import (
	"context"

	"github.com/kristiannissen/concertlist/internal/domain"
)

// Service implements the domain.ETLService interface.
type Service struct {
	extractors []domain.Extractor
	storage    domain.BlobStorage
}

// NewService creates a new ETL service with the given extractors and storage.
func NewService(extractors []domain.Extractor, storage domain.BlobStorage) *Service {
	return &Service{
		extractors: extractors,
		storage:    storage,
	}
}

// Run executes the full ETL workflow.
func (s *Service) Run() ([]MusicEvent, error) {
	ctx := context.Background()

	// Extract: Collect events from all sources
	var allEvents []domain.MusicEvent
	for _, extractor := range s.extractors {
		events, err := extractor.Extract()
		if err != nil {
			return nil, err
		}
		allEvents = append(allEvents, events...)
	}

	// Load: Store the extracted events
	if err := s.storage.Store(allEvents); err != nil {
		return nil, err
	}

	return allEvents, nil
}

// RunWithContext executes the full ETL workflow with a provided context.
func (s *Service) RunWithContext(ctx context.Context) ([]domain.MusicEvent, error) {
	// Extract: Collect events from all sources
	var allEvents []domain.MusicEvent
	for _, extractor := range s.extractors {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
			events, err := extractor.Extract()
			if err != nil {
				return nil, err
			}
			allEvents = append(allEvents, events...)
		}
	}

	// Load: Store the extracted events
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		if err := s.storage.Store(allEvents); err != nil {
			return nil, err
		}
	}

	return allEvents, nil
}
