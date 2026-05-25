package usecase

import (
	"aggregator/internal/domain/entities"
	"aggregator/internal/dto"
	"aggregator/internal/usecase/apperrors"
	"context"
	"errors"
	"testing"
)

func TestPersistData_Execute_PersistsAndDeletes(t *testing.T) {
	store := &fakeProcessedEventStore{}
	deleter := &fakeQueueDeleter{}
	handler := NewPersistDataHandler(store, deleter)

	err := handler.Execute(context.Background(), validQueueMessage())
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if store.saved.EventID != "2f4db003-dbc5-4b1b-bf7b-1f847b69ce1f" {
		t.Fatalf("expected event to be saved, got %+v", store.saved)
	}
	if deleter.deletedReceiptHandle != "receipt-1" {
		t.Fatalf("expected message to be deleted, got %s", deleter.deletedReceiptHandle)
	}
}

func TestPersistData_Execute_DuplicateEventDeletesMessage(t *testing.T) {
	store := &fakeProcessedEventStore{saveErr: apperrors.ErrEventAlreadyProcessed}
	deleter := &fakeQueueDeleter{}
	handler := NewPersistDataHandler(store, deleter)

	err := handler.Execute(context.Background(), validQueueMessage())
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if deleter.deletedReceiptHandle != "receipt-1" {
		t.Fatalf("expected duplicate message to be deleted, got %s", deleter.deletedReceiptHandle)
	}
}

func TestPersistData_Execute_StoreErrorDoesNotDeleteMessage(t *testing.T) {
	storeErr := errors.New("dynamodb unavailable")
	store := &fakeProcessedEventStore{saveErr: storeErr}
	deleter := &fakeQueueDeleter{}
	handler := NewPersistDataHandler(store, deleter)

	err := handler.Execute(context.Background(), validQueueMessage())
	if !errors.Is(err, storeErr) {
		t.Fatalf("expected store error, got %v", err)
	}
	if deleter.deletedReceiptHandle != "" {
		t.Fatalf("expected message to remain in queue, got deleted receipt %s", deleter.deletedReceiptHandle)
	}
}

func TestPersistData_Execute_DeleteErrorReturnsError(t *testing.T) {
	deleteErr := errors.New("delete failed")
	store := &fakeProcessedEventStore{}
	deleter := &fakeQueueDeleter{deleteErr: deleteErr}
	handler := NewPersistDataHandler(store, deleter)

	err := handler.Execute(context.Background(), validQueueMessage())
	if !errors.Is(err, deleteErr) {
		t.Fatalf("expected delete error, got %v", err)
	}
}

func TestPersistData_Execute_InvalidJSONDoesNotDeleteMessage(t *testing.T) {
	store := &fakeProcessedEventStore{}
	deleter := &fakeQueueDeleter{}
	handler := NewPersistDataHandler(store, deleter)

	err := handler.Execute(context.Background(), dto.QueueMessage{
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
	saved   entities.ProcessedEvent
	saveErr error
}

func (s *fakeProcessedEventStore) SaveEventAndIncrementSummary(ctx context.Context, event entities.ProcessedEvent) error {
	s.saved = event
	return s.saveErr
}

func (s *fakeProcessedEventStore) FindByID(ctx context.Context, eventID string) (entities.ProcessedEvent, error) {
	return entities.ProcessedEvent{}, nil
}

func (s *fakeProcessedEventStore) FindByDeveloperID(ctx context.Context, developerID string) ([]entities.ProcessedEvent, error) {
	return nil, nil
}

type fakeQueueDeleter struct {
	deletedReceiptHandle string
	deleteErr            error
}

func (d *fakeQueueDeleter) Delete(ctx context.Context, receiptHandle string) error {
	d.deletedReceiptHandle = receiptHandle
	return d.deleteErr
}
