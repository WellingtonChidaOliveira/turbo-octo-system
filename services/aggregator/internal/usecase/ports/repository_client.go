package ports

import (
	"aggregator/internal/domain/entities"
	"context"
)

type ProcessedEventWriter interface {
	SaveEventAndIncrementSummary(ctx context.Context, event entities.ProcessedEvent) error
	UpdateLastActivityIfNewer(ctx context.Context, developerID string, timestamp string) error
}

type ProcessedEventReader interface {
	FindByID(ctx context.Context, eventID string) (entities.ProcessedEvent, error)
	FindByDeveloperID(ctx context.Context, developerID string) ([]entities.ProcessedEvent, error)
}

type DeveloperSummaryStore interface {
	FindSummaryByDeveloperID(ctx context.Context, developerID string) (entities.DeveloperSummary, error)
}
