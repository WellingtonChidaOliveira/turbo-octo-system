package usecase

import (
	"aggregator/internal/dto"
	"aggregator/internal/usecase/ports"
	"aggregator/internal/usecase/retry"
	"context"
	"log/slog"
)

type ProcessedEventConsumer struct {
	consumer ports.QueueConsumer
	retry    retry.Policy
}

func NewProcessedEventConsumer(consumer ports.QueueConsumer, retryPolicy retry.Policy) *ProcessedEventConsumer {
	return &ProcessedEventConsumer{
		consumer: consumer,
		retry:    retryPolicy.Normalize(),
	}
}

func (c *ProcessedEventConsumer) Start(ctx context.Context, jobs chan<- dto.QueueMessage) {
	defer close(jobs)

	attempt := 0
	for ctx.Err() == nil {
		messages, err := c.consumer.Receive(ctx)
		if err != nil {
			attempt++
			c.handleReceiveError(ctx, err, attempt)
			continue
		}

		attempt = 0
		c.dispatchMessages(ctx, messages, jobs)
	}

	slog.Info("processed event consumer stopped")
}

func (c *ProcessedEventConsumer) dispatchMessages(ctx context.Context, messages []dto.QueueMessage, jobs chan<- dto.QueueMessage) {
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

func (c *ProcessedEventConsumer) handleReceiveError(ctx context.Context, err error, attempt int) {
	if ctx.Err() != nil {
		return
	}

	backoff := c.retry.Delay(attempt)
	slog.Warn("failed to receive processed events",
		"stage", "receive",
		"attempt", attempt,
		"backoff_ms", backoff.Milliseconds(),
		"error", err,
	)

	_ = c.retry.Wait(ctx, attempt)
}
