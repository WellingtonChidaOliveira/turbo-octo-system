package usecase

import (
	"context"
	"encoding/json"
	"log/slog"
	"processor/internal/domain/entities"
	"processor/internal/usecase/ports"
	"processor/internal/usecase/retry"
)

type EventProcessor struct {
	publisher   ports.QueuePublisher
	deleter     ports.QueueDeleter
	clock       ports.Clock
	processorID string
	retry       retry.Policy
}

func NewEventProcessor(
	publisher ports.QueuePublisher,
	deleter ports.QueueDeleter,
	clock ports.Clock,
	processorID string,
	retryPolicy retry.Policy) *EventProcessor {
	return &EventProcessor{
		publisher:   publisher,
		deleter:     deleter,
		processorID: processorID,
		clock:       clock,
		retry:       retryPolicy.Normalize(),
	}
}

func (p *EventProcessor) Handle(ctx context.Context, message entities.QueueMessage) error {
	return p.ProcessMessage(ctx, message)
}

func (p *EventProcessor) ProcessMessage(ctx context.Context, msg entities.QueueMessage) error {
	processedMsg, eventID, errorType, err := p.buildProcessedMessage(msg)
	if err != nil {
		processingErr := ProcessingError{Type: errorType, EventID: eventID, Err: err}
		slog.Error("event processing failed",
			"event_id", eventID,
			"correlation_id", correlationID(eventID, msg.ID),
			"message_id", msg.ID,
			"stage", "process",
			"error_type", processingErr.Type,
			"error", err,
		)
		return processingErr
	}
	slog.Info("event validated",
		"event_id", eventID,
		"correlation_id", correlationID(eventID, msg.ID),
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
			"correlation_id", correlationID(eventID, msg.ID),
			"message_id", msg.ID,
			"stage", "delete",
			"error_type", ErrorTypeDelete,
			"error", err,
		)
		return ProcessingError{Type: ErrorTypeDelete, EventID: eventID, Err: err}
	}
	slog.Info("event deleted",
		"event_id", eventID,
		"correlation_id", correlationID(eventID, msg.ID),
		"message_id", msg.ID,
		"stage", "delete",
	)

	return nil
}

func (p *EventProcessor) publishWithBackoff(ctx context.Context, eventID string, messageID string, msg entities.QueueMessage) error {
	var err error
	for attempt := 1; attempt <= p.retry.MaxAttempts; attempt++ {
		err = p.publisher.Send(ctx, msg)
		if err == nil {
			slog.Info("event published",
				"event_id", eventID,
				"correlation_id", correlationID(eventID, messageID),
				"message_id", messageID,
				"stage", "publish",
				"attempt", attempt,
			)
			return nil
		}

		backoff := p.retry.Delay(attempt)
		slog.Warn("event publish attempt failed",
			"event_id", eventID,
			"correlation_id", correlationID(eventID, messageID),
			"message_id", messageID,
			"stage", "publish",
			"attempt", attempt,
			"max_attempts", p.retry.MaxAttempts,
			"backoff_ms", backoff.Milliseconds(),
			"error_type", ErrorTypePublish,
			"error", err,
		)

		if attempt == p.retry.MaxAttempts {
			break
		}
		if waitErr := p.retry.Wait(ctx, attempt); waitErr != nil {
			return ProcessingError{Type: ErrorTypePublish, EventID: eventID, Err: waitErr}
		}
	}

	slog.Error("event publish failed",
		"event_id", eventID,
		"correlation_id", correlationID(eventID, messageID),
		"message_id", messageID,
		"stage", "publish",
		"attempts", p.retry.MaxAttempts,
		"error_type", ErrorTypePublish,
		"error", err,
	)
	return ProcessingError{Type: ErrorTypePublish, EventID: eventID, Err: err}
}

func (p *EventProcessor) buildProcessedMessage(msg entities.QueueMessage) (string, string, ErrorType, error) {
	var event entities.RawEvent
	if err := json.Unmarshal([]byte(msg.Body), &event); err != nil {
		return "", "", ErrorTypeDecode, err
	}

	if err := event.Validate(p.clock.Now()); err != nil {
		return "", event.EventID, ErrorTypeValidation, err
	}

	processEvent := event.ToProcessed(p.processorID, p.clock.Now())
	processedMsg, err := json.Marshal(processEvent)
	if err != nil {
		return "", event.EventID, ErrorTypeEncode, err
	}

	return string(processedMsg), event.EventID, "", nil
}

func correlationID(eventID string, messageID string) string {
	if eventID != "" {
		return eventID
	}
	return messageID
}
