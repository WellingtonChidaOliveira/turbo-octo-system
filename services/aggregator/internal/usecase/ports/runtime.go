package ports

import (
	"aggregator/internal/dto"
	"context"
)

type JobConsumer interface {
	Start(ctx context.Context, jobs chan<- dto.QueueMessage)
}

type WorkerPool interface {
	Start(ctx context.Context, jobs <-chan dto.QueueMessage)
}
