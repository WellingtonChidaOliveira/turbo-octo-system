package usecase

import (
	"context"
	"log/slog"
	"processor/internal/domain/entities"
	"processor/internal/usecase/ports"
	"time"
)

const (
	receiveInitialBackoff = time.Second
	receiveMaxBackoff     = 30 * time.Second
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
	defer close(jobs)

	backoff := receiveInitialBackoff
	for ctx.Err() == nil {
		messages, err := c.consumer.Receive(ctx)
		if err != nil {
			backoff = c.handleReceiveError(ctx, err, backoff)
			continue
		}

		backoff = receiveInitialBackoff
		c.dispatchMessages(ctx, messages, jobs)
	}

	slog.Info("raw event consumer stopped")
}

func (c *RawEventConsumer) dispatchMessages(ctx context.Context, messages []entities.QueueMessage, jobs chan<- entities.QueueMessage) {
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

func (c *RawEventConsumer) handleReceiveError(ctx context.Context, err error, backoff time.Duration) time.Duration {
	if ctx.Err() != nil {
		return backoff
	}

	slog.Warn("failed to receive raw events",
		"stage", "receive",
		"backoff_ms", backoff.Milliseconds(),
		"error", err,
	)

	if waitErr := waitBackoff(ctx, backoff); waitErr != nil {
		return backoff
	}
	return nextBackoff(backoff, receiveMaxBackoff)
}
