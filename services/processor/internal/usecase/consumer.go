package usecase

import (
	"context"
	"log/slog"
	"processor/internal/dto"
	"processor/internal/usecase/ports"
	"processor/internal/usecase/retry"
)

type RawEventConsumer struct {
	consumer ports.QueueConsumer
	retry    retry.Policy
}

func NewRawEventsConsumer(client ports.QueueConsumer, retryPolicy retry.Policy) *RawEventConsumer {
	return &RawEventConsumer{
		consumer: client,
		retry:    retryPolicy.Normalize(),
	}
}

func (c *RawEventConsumer) Consumer(ctx context.Context, jobs chan<- dto.QueueMessage) {
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

	slog.Info("raw event consumer stopped")
}

func (c *RawEventConsumer) dispatchMessages(ctx context.Context, messages []dto.QueueMessage, jobs chan<- dto.QueueMessage) {
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

func (c *RawEventConsumer) handleReceiveError(ctx context.Context, err error, attempt int) {
	if ctx.Err() != nil {
		return
	}

	backoff := c.retry.Delay(attempt)
	slog.Warn("failed to receive raw events",
		"stage", "receive",
		"attempt", attempt,
		"backoff_ms", backoff.Milliseconds(),
		"error", err,
	)

	_ = c.retry.Wait(ctx, attempt)
}
