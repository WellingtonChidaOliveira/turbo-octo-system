package abstractions

import (
	"context"
)

type QueueMessage struct {
	Body          string
	ReceiptHandle string
}

type QueueConsumerInterface interface {
	Receive(ctx context.Context) ([]QueueMessage, error)
	Delete(ctx context.Context, receiptHandle string) error
}

type QueuePublisherInterface interface {
	Send(ctx context.Context, message QueueMessage) error
}
