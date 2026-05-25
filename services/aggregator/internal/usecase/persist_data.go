package usecase

import (
	"aggregator/internal/dto"
	"aggregator/internal/usecase/apperrors"
	"aggregator/internal/usecase/ports"
	"context"
	"encoding/json"
	"errors"
	"log/slog"
)

type PersistProcessedEvent struct {
	store   ports.ProcessedEventWriter
	deleter ports.QueueDeleter
}

func NewPersistProcessedEventHandler(store ports.ProcessedEventWriter, deleter ports.QueueDeleter) *PersistProcessedEvent {
	return &PersistProcessedEvent{
		store:   store,
		deleter: deleter,
	}
}

func (h *PersistProcessedEvent) Handle(ctx context.Context, msg dto.QueueMessage) error {
	var eventDTO dto.ProcessedEvent
	if err := json.Unmarshal([]byte(msg.Body), &eventDTO); err != nil {
		slog.Error("failed to decode processed event",
			"message_id", msg.ID,
			"stage", "decode",
			"error", err,
		)
		return err
	}

	event := eventDTO.ToEntity()
	if err := event.Validate(); err != nil {
		slog.Error("invalid processed event",
			"event_id", event.EventID,
			"correlation_id", event.EventID,
			"message_id", msg.ID,
			"stage", "validate",
			"error", err,
		)
		return err
	}

	err := h.store.SaveEventAndIncrementSummary(ctx, event)
	if err != nil && !errors.Is(err, apperrors.ErrEventAlreadyProcessed) {
		slog.Error("failed to persist processed event",
			"event_id", event.EventID,
			"correlation_id", event.EventID,
			"message_id", msg.ID,
			"stage", "persist",
			"error", err,
		)
		return err
	}
	if errors.Is(err, apperrors.ErrEventAlreadyProcessed) {
		slog.Info("processed event already persisted",
			"event_id", event.EventID,
			"correlation_id", event.EventID,
			"message_id", msg.ID,
			"stage", "persist",
		)
	}
	if err == nil || errors.Is(err, apperrors.ErrEventAlreadyProcessed) {
		if err := h.store.UpdateLastActivityIfNewer(ctx, event.DeveloperID, event.Timestamp); err != nil {
			slog.Error("failed to update last activity",
				"event_id", event.EventID,
				"correlation_id", event.EventID,
				"message_id", msg.ID,
				"stage", "last_activity",
				"error", err,
			)
			return err
		}
	}

	if err := h.deleter.Delete(ctx, msg.ReceiptHandle); err != nil {
		slog.Error("failed to delete processed event message",
			"event_id", event.EventID,
			"correlation_id", event.EventID,
			"message_id", msg.ID,
			"stage", "delete",
			"error", err,
		)
		return err
	}

	slog.Info("processed event persisted",
		"event_id", event.EventID,
		"correlation_id", event.EventID,
		"message_id", msg.ID,
		"stage", "persist",
	)

	return nil
}
