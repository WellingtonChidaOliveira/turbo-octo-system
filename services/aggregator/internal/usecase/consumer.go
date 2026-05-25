package usecase

import (
	"aggregator/internal/dto"
	"aggregator/internal/usecase/ports"
	"context"
	"log/slog"
)

type ProcessedEventConsumer struct {
	consumer ports.QueueConsumer
}

func NewProcessedEventConsumer(consumer ports.QueueConsumer) *ProcessedEventConsumer {
	return &ProcessedEventConsumer{
		consumer: consumer,
	}
}

func (c *ProcessedEventConsumer) Consume(ctx context.Context, jobs chan<- dto.QueueMessage) {
	defer close(jobs)

	for ctx.Err() == nil {
		messages, err := c.consumer.Receive(ctx)
		if err != nil {
			slog.Error("failed to receive messages", "error", err)
			continue
		}

		for _, msg := range messages {
			select {
			case <-ctx.Done():
				return
			case jobs <- msg:
				slog.Info("message queued",
					"message_id", msg.ID,
					"stage", "dispatch",
				)
			}
		}
	}

	slog.Info("processed event consumer stopped")
}
