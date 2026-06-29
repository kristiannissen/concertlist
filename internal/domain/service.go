// Package domain contains core business models and logic.
package domain

import "context"

// ETLServiceImpl implements ETLService.
type ETLServiceImpl struct {
	extractors []ExtractorPort
	storage    StoragePort
}

// NewETLService creates a new ETLService.
func NewETLService(extractors []ExtractorPort, storage StoragePort) *ETLServiceImpl {
	return &ETLServiceImpl{
		extractors: extractors,
		storage:    storage,
	}
}

// Run executes the ETL pipeline for all extractors.
func (s *ETLServiceImpl) Run(ctx context.Context) ([]Concert, error) {
	var allConcerts []Concert
	for _, extractor := range s.extractors {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
			concerts, err := extractor.Extract(ctx)
			if err != nil {
				return nil, err
			}
			allConcerts = append(allConcerts, concerts...)
		}
	}

	if err := s.storage.Save(ctx, allConcerts); err != nil {
		return nil, err
	}

	return allConcerts, nil
}
