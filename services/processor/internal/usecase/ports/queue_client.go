package ports

import (
	"context"
	"processor/internal/domain/entities"
)

type QueueConsumer interface {
	Receive(ctx context.Context) ([]entities.QueueMessage, error)
	Delete(ctx context.Context, receiptHandle string) error
}

type QueuePublisher interface {
	Send(ctx context.Context, message entities.QueueMessage) error
}
