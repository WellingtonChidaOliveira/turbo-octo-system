package usecase

import (
	"context"
	"log"
	"processor/internal/domain/entities"
	"processor/internal/usecase/ports"
	"time"
)

type RawEventConsumer struct {
	consumer ports.QueueConsumer
}

func NewRawEventsConsumer(client ports.QueueConsumer) *RawEventConsumer {
	return &RawEventConsumer{
		consumer: client,
	}
}

func (c *RawEventConsumer) Consumer(ctx context.Context, jobs chan<- entities.QueueMessage) {
	for ctx.Err() == nil {
		messages, err := c.consumer.Receive(ctx)
		if err != nil {
			c.handleReceiveError(ctx, err)
			continue
		}

		c.dispatchMessages(ctx, messages, jobs)
	}
}

func (c *RawEventConsumer) dispatchMessages(ctx context.Context, messages []entities.QueueMessage, jobs chan<- entities.QueueMessage) {
	for _, msg := range messages {
		select {
		case <-ctx.Done():
			return
		case jobs <- msg:
		}
	}
}

func (c *RawEventConsumer) handleReceiveError(ctx context.Context, err error) {
	if ctx.Err() != nil {
		return
	}

	log.Printf("failed to receive raw events message %v", err)
	time.Sleep(time.Second)
}
