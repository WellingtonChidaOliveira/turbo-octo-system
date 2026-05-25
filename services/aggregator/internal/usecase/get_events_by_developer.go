package usecase

import (
	"aggregator/internal/domain/entities"
	"aggregator/internal/usecase/ports"
	"context"
)

type GetProcessedEventsByDeveloper struct {
	store ports.ProcessedEventStore
}

func NewGetProcessedEventsByDeveloper(store ports.ProcessedEventStore) *GetProcessedEventsByDeveloper {
	return &GetProcessedEventsByDeveloper{store: store}
}

func (uc *GetProcessedEventsByDeveloper) Handle(ctx context.Context, developerID string) ([]entities.ProcessedEvent, error) {
	return uc.store.FindByDeveloperID(ctx, developerID)
}
