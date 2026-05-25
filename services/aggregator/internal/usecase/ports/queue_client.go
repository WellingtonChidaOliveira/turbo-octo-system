package ports

import (
	"aggregator/internal/dto"
	"context"
)

type QueueConsumer interface {
	Receive(ctx context.Context) ([]dto.QueueMessage, error)
	QueueDeleter
}

type QueueDeleter interface {
	Delete(ctx context.Context, receiptHandle string) error
}
