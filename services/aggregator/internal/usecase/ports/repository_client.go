package ports

import (
	"aggregator/internal/domain/entities"
	"context"
)

type ProcessedEventStore interface {
	Save(ctx context.Context, event entities.ProcessedEvent) error
	FindByID(ctx context.Context, eventID string) (entities.ProcessedEvent, error)
	FindByDeveloperID(ctx context.Context, developerID string) ([]entities.ProcessedEvent, error)
}
