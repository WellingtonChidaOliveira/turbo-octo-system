package usecase

import (
	"aggregator/internal/domain/entities"
	"aggregator/internal/usecase/ports"
	"context"
)

type GetProcessedEventByID struct {
	store ports.ProcessedEventStore
}

func NewGetProcessedEventByID(store ports.ProcessedEventStore) *GetProcessedEventByID {
	return &GetProcessedEventByID{store: store}
}

func (uc *GetProcessedEventByID) Execute(ctx context.Context, eventID string) (entities.ProcessedEvent, error) {
	return uc.store.FindByID(ctx, eventID)
}
