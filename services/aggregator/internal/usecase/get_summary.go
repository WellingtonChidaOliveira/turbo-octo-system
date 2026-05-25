package usecase

import (
	"aggregator/internal/domain/entities"
	"aggregator/internal/usecase/ports"
	"context"
)

type GetDeveloperSummary struct {
	store ports.DeveloperSummaryStore
}

func NewGetDeveloperSummary(store ports.DeveloperSummaryStore) *GetDeveloperSummary {
	return &GetDeveloperSummary{store: store}
}

func (uc *GetDeveloperSummary) Handle(ctx context.Context, developerID string) (entities.DeveloperSummary, error) {
	return uc.store.FindSummaryByDeveloperID(ctx, developerID)
}
