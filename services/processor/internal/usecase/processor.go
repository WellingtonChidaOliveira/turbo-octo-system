package usecase

import (
	"context"
	"encoding/json"
	"log/slog"
	"processor/internal/domain/entities"
	"processor/internal/usecase/ports"
	"time"
)

const (
	publishMaxAttempts    = 3
	publishInitialBackoff = 100 * time.Millisecond
	publishMaxBackoff     = 2 * time.Second
)

type EventProcessor struct {
	publisher   ports.QueuePublisher
	deleter     ports.QueueDeleter
	clock       ports.Clock
	processorID string
}

func NewEventProcessor(
	publisher ports.QueuePublisher,
	deleter ports.QueueDeleter,
	clock ports.Clock,
	processorID string) *EventProcessor {
	return &EventProcessor{
		publisher:   publisher,
		deleter:     deleter,
		processorID: processorID,
		clock:       clock,
	}
}

func (p *EventProcessor) Handle(ctx context.Context, message entities.QueueMessage) error {
	return p.ProcessMessage(ctx, message)
}

func (p *EventProcessor) ProcessMessage(ctx context.Context, msg entities.QueueMessage) error {
	processedMsg, eventID, err := p.trateMessage(msg)
	if err != nil {
		slog.Error("event processing failed",
			"event_id", eventID,
			"message_id", msg.ID,
			"stage", "process",
			"error", err,
		)
		return err
	}
	slog.Info("event validated",
		"event_id", eventID,
		"message_id", msg.ID,
		"stage", "process",
	)

	if err = p.publishWithBackoff(ctx, eventID, msg.ID, entities.QueueMessage{Body: string(processedMsg)}); err != nil {
		return err
	}

	err = p.deleter.Delete(ctx, msg.ReceiptHandle)
	if err != nil {
		slog.Error("event delete failed",
			"event_id", eventID,
			"message_id", msg.ID,
			"stage", "delete",
			"error", err,
		)
		return err
	}
	slog.Info("event deleted",
		"event_id", eventID,
		"message_id", msg.ID,
		"stage", "delete",
	)

	return nil
}

func (p *EventProcessor) publishWithBackoff(ctx context.Context, eventID string, messageID string, msg entities.QueueMessage) error {
	backoff := publishInitialBackoff

	var err error
	for attempt := 1; attempt <= publishMaxAttempts; attempt++ {
		err = p.publisher.Send(ctx, msg)
		if err == nil {
			slog.Info("event published",
				"event_id", eventID,
				"message_id", messageID,
				"stage", "publish",
				"attempt", attempt,
			)
			return nil
		}

		slog.Warn("event publish attempt failed",
			"event_id", eventID,
			"message_id", messageID,
			"stage", "publish",
			"attempt", attempt,
			"max_attempts", publishMaxAttempts,
			"backoff_ms", backoff.Milliseconds(),
			"error", err,
		)

		if attempt == publishMaxAttempts {
			break
		}
		if waitErr := waitBackoff(ctx, backoff); waitErr != nil {
			return waitErr
		}
		backoff = nextBackoff(backoff, publishMaxBackoff)
	}

	slog.Error("event publish failed",
		"event_id", eventID,
		"message_id", messageID,
		"stage", "publish",
		"attempts", publishMaxAttempts,
		"error", err,
	)
	return err
}

func waitBackoff(ctx context.Context, delay time.Duration) error {
	timer := time.NewTimer(delay)
	defer timer.Stop()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-timer.C:
		return nil
	}
}

func nextBackoff(current time.Duration, max time.Duration) time.Duration {
	next := current * 2
	if next > max {
		return max
	}
	return next
}

func (p *EventProcessor) trateMessage(msg entities.QueueMessage) (string, string, error) {
	var event entities.RawEvent
	if err := json.Unmarshal([]byte(msg.Body), &event); err != nil {
		return "", "", err
	}

	if err := event.Validate(p.clock.Now()); err != nil {
		return "", event.EventID, err
	}

	processEvent := event.ToProcessed(p.processorID, p.clock.Now())
	processedMsg, err := json.Marshal(processEvent)
	if err != nil {
		return "", event.EventID, err
	}

	return string(processedMsg), event.EventID, nil
}
