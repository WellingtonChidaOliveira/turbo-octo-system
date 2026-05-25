package usecase

import (
	"aggregator/internal/domain/entities"
	"aggregator/internal/dto"
	"aggregator/internal/usecase/apperrors"
	"context"
	"errors"
	"strings"
	"testing"
)

func TestPersistProcessedEvent_Handle_PersistsUpdatesLastActivityAndDeletes(t *testing.T) {
	store := &fakeProcessedEventStore{}
	deleter := &fakeQueueDeleter{}
	handler := NewPersistProcessedEventHandler(store, deleter)

	err := handler.Handle(context.Background(), validQueueMessage())
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if store.saved.EventID != "2f4db003-dbc5-4b1b-bf7b-1f847b69ce1f" {
		t.Fatalf("expected event to be saved, got %+v", store.saved)
	}
	if store.lastActivityDeveloperID != "dev-123" || store.lastActivityTimestamp != "2026-04-15T10:30:00Z" {
		t.Fatalf("expected last activity update, got developer=%s timestamp=%s", store.lastActivityDeveloperID, store.lastActivityTimestamp)
	}
	if deleter.deletedReceiptHandle != "receipt-1" {
		t.Fatalf("expected message to be deleted, got %s", deleter.deletedReceiptHandle)
	}
}

func TestPersistProcessedEvent_Handle_DuplicateEventDeletesMessage(t *testing.T) {
	store := &fakeProcessedEventStore{saveErr: apperrors.ErrEventAlreadyProcessed}
	deleter := &fakeQueueDeleter{}
	handler := NewPersistProcessedEventHandler(store, deleter)

	err := handler.Handle(context.Background(), validQueueMessage())
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if store.lastActivityDeveloperID != "dev-123" || store.lastActivityTimestamp != "2026-04-15T10:30:00Z" {
		t.Fatalf("expected duplicate to update last activity, got developer=%s timestamp=%s", store.lastActivityDeveloperID, store.lastActivityTimestamp)
	}
	if deleter.deletedReceiptHandle != "receipt-1" {
		t.Fatalf("expected duplicate message to be deleted, got %s", deleter.deletedReceiptHandle)
	}
}

func TestPersistProcessedEvent_Handle_StoreErrorDoesNotDeleteMessage(t *testing.T) {
	storeErr := errors.New("dynamodb unavailable")
	store := &fakeProcessedEventStore{saveErr: storeErr}
	deleter := &fakeQueueDeleter{}
	handler := NewPersistProcessedEventHandler(store, deleter)

	err := handler.Handle(context.Background(), validQueueMessage())
	if !errors.Is(err, storeErr) {
		t.Fatalf("expected store error, got %v", err)
	}
	if deleter.deletedReceiptHandle != "" {
		t.Fatalf("expected message to remain in queue, got deleted receipt %s", deleter.deletedReceiptHandle)
	}
}

func TestPersistProcessedEvent_Handle_LastActivityErrorDoesNotDeleteMessage(t *testing.T) {
	updateErr := errors.New("conditional update unavailable")
	store := &fakeProcessedEventStore{lastActivityErr: updateErr}
	deleter := &fakeQueueDeleter{}
	handler := NewPersistProcessedEventHandler(store, deleter)

	err := handler.Handle(context.Background(), validQueueMessage())
	if !errors.Is(err, updateErr) {
		t.Fatalf("expected last activity error, got %v", err)
	}
	if deleter.deletedReceiptHandle != "" {
		t.Fatalf("expected message to remain in queue, got deleted receipt %s", deleter.deletedReceiptHandle)
	}
}

func TestPersistProcessedEvent_Handle_DeleteErrorReturnsError(t *testing.T) {
	deleteErr := errors.New("delete failed")
	store := &fakeProcessedEventStore{}
	deleter := &fakeQueueDeleter{deleteErr: deleteErr}
	handler := NewPersistProcessedEventHandler(store, deleter)

	err := handler.Handle(context.Background(), validQueueMessage())
	if !errors.Is(err, deleteErr) {
		t.Fatalf("expected delete error, got %v", err)
	}
}

func TestPersistProcessedEvent_Handle_InvalidJSONDoesNotDeleteMessage(t *testing.T) {
	store := &fakeProcessedEventStore{}
	deleter := &fakeQueueDeleter{}
	handler := NewPersistProcessedEventHandler(store, deleter)

	err := handler.Handle(context.Background(), dto.QueueMessage{
		ID:            "message-1",
		Body:          "invalid",
		ReceiptHandle: "receipt-1",
	})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if deleter.deletedReceiptHandle != "" {
		t.Fatalf("expected message to remain in queue, got deleted receipt %s", deleter.deletedReceiptHandle)
	}
}

func TestPersistProcessedEvent_Handle_InvalidContractDoesNotPersistOrDelete(t *testing.T) {
	tests := []struct {
		name        string
		replacement string
	}{
		{name: "missing developer_id", replacement: `"developer_id": ""`},
		{name: "invalid metric_type", replacement: `"metric_type": "deploys"`},
		{name: "negative value", replacement: `"value": -1`},
		{name: "review_time_minutes above max", replacement: `"metric_type": "review_time_minutes", "value": 1441`},
		{name: "missing processed_at", replacement: `"processed_at": ""`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := &fakeProcessedEventStore{}
			deleter := &fakeQueueDeleter{}
			handler := NewPersistProcessedEventHandler(store, deleter)
			msg := validQueueMessage()
			msg.Body = strings.Replace(msg.Body, contractSnippet(tt.replacement), tt.replacement, 1)

			err := handler.Handle(context.Background(), msg)
			if err == nil {
				t.Fatal("expected validation error, got nil")
			}
			if store.saved.EventID != "" {
				t.Fatalf("expected event not to be saved, got %+v", store.saved)
			}
			if deleter.deletedReceiptHandle != "" {
				t.Fatalf("expected message to remain in queue, got deleted receipt %s", deleter.deletedReceiptHandle)
			}
		})
	}
}

func contractSnippet(replacement string) string {
	switch {
	case strings.Contains(replacement, "developer_id"):
		return `"developer_id": "dev-123"`
	case strings.Contains(replacement, "metric_type") && strings.Contains(replacement, "value"):
		return `"metric_type": "commits",
				"value": 10`
	case strings.Contains(replacement, "metric_type"):
		return `"metric_type": "commits"`
	case strings.Contains(replacement, "value"):
		return `"value": 10`
	case strings.Contains(replacement, "processed_at"):
		return `"processed_at": "2026-04-15T10:30:05Z"`
	default:
		return replacement
	}
}

func validQueueMessage() dto.QueueMessage {
	return dto.QueueMessage{
		ID:            "message-1",
		ReceiptHandle: "receipt-1",
		Body: `{
				"event_id": "2f4db003-dbc5-4b1b-bf7b-1f847b69ce1f",
				"developer_id": "dev-123",
				"metric_type": "commits",
				"value": 10,
				"repository": "org/repo",
				"timestamp": "2026-04-15T10:30:00Z",
				"processed_at": "2026-04-15T10:30:05Z",
				"processor_id": "processor-1"
			}`,
	}
}

type fakeProcessedEventStore struct {
	saved                   entities.ProcessedEvent
	saveErr                 error
	lastActivityDeveloperID string
	lastActivityTimestamp   string
	lastActivityErr         error
}

func (s *fakeProcessedEventStore) SaveEventAndIncrementSummary(ctx context.Context, event entities.ProcessedEvent) error {
	s.saved = event
	return s.saveErr
}

func (s *fakeProcessedEventStore) UpdateLastActivityIfNewer(ctx context.Context, developerID string, timestamp string) error {
	s.lastActivityDeveloperID = developerID
	s.lastActivityTimestamp = timestamp
	return s.lastActivityErr
}

type fakeQueueDeleter struct {
	deletedReceiptHandle string
	deleteErr            error
}

func (d *fakeQueueDeleter) Delete(ctx context.Context, receiptHandle string) error {
	d.deletedReceiptHandle = receiptHandle
	return d.deleteErr
}
