package usecase

import (
	"aggregator/internal/domain/entities"
	"aggregator/internal/usecase/ports"
	"context"
)

type GetProcessedEventByID struct {
	store ports.ProcessedEventReader
}

func NewGetProcessedEventByID(store ports.ProcessedEventReader) *GetProcessedEventByID {
	return &GetProcessedEventByID{store: store}
}

func (uc *GetProcessedEventByID) Handle(ctx context.Context, eventID string) (entities.ProcessedEvent, error) {
	return uc.store.FindByID(ctx, eventID)
}
