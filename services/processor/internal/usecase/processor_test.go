package usecase

import (
	"context"
	"encoding/json"
	"errors"
	"processor/internal/domain/entities"
	"sync"
	"testing"
	"time"
)

type fakeConsumer struct {
	mu        sync.Mutex
	deleted   []string
	deleteErr error
}

func (f *fakeConsumer) Receive(ctx context.Context) ([]entities.QueueMessage, error) {
	return nil, nil
}

func (f *fakeConsumer) Delete(ctx context.Context, msg string) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.deleted = append(f.deleted, msg)
	return f.deleteErr
}

func (f *fakeConsumer) deletedCount() int {
	f.mu.Lock()
	defer f.mu.Unlock()
	return len(f.deleted)
}

type fakePublisher struct {
	mu      sync.Mutex
	sent    []string
	sendErr error
}

func (f *fakePublisher) Send(ctx context.Context, msg entities.QueueMessage) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.sent = append(f.sent, msg.Body)
	return f.sendErr
}

func (f *fakePublisher) sentCount() int {
	f.mu.Lock()
	defer f.mu.Unlock()
	return len(f.sent)
}

type fakeClock struct {
	now time.Time
}

func (f *fakeClock) Now() time.Time {
	return f.now
}

func newTestProcessor(consumer *fakeConsumer, publisher *fakePublisher) *EventProcessor {
	clock := &fakeClock{now: time.Date(2026, 4, 15, 10, 30, 5, 0, time.UTC)}
	return NewEventProcessor(publisher, consumer, clock, "processor-1")
}

func validQueueMessage() entities.QueueMessage {
	return validQueueMessageWithReceipt("receipt-1")
}

func validQueueMessageWithReceipt(receiptHandle string) entities.QueueMessage {
	return entities.QueueMessage{
		Body: `{
	              "event_id": "550e8400-e29b-41d4-a716-446655440000",
              "developer_id": "dev-123",
              "metric_type": "commits",
              "value": 10,
              "repository": "org/repo",
	              "timestamp": "2026-04-15T10:30:00Z"
	          }`,
		ReceiptHandle: receiptHandle,
	}
}

func TestProcessMessage_ValidMessage_PublishesAndDeletes(t *testing.T) {
	consumer := &fakeConsumer{}
	publisher := &fakePublisher{}
	processor := newTestProcessor(consumer, publisher)
	msg := validQueueMessage()

	err := processor.ProcessMessage(context.Background(), msg)
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	if len(publisher.sent) != 1 {
		t.Fatalf("expected 1 published message, got %d", len(publisher.sent))
	}

	if len(consumer.deleted) != 1 {
		t.Fatalf("expected 1 deleted message, got %d", len(consumer.deleted))
	}

	if consumer.deleted[0] != "receipt-1" {
		t.Fatalf("expected receipt-1 deleted, got %s", consumer.deleted[0])
	}

	var processed entities.ProcessedEvent
	err = json.Unmarshal([]byte(publisher.sent[0]), &processed)
	if err != nil {
		t.Fatal(err)
	}

	if processed.ProcessorID != "processor-1" {
		t.Fatalf("expected processor_id processor-1, got %s", processed.ProcessorID)
	}
	if processed.ProcessedAt != "2026-04-15T10:30:05Z" {
		t.Fatalf("expected processed_at 2026-04-15T10:30:05Z, got %s", processed.ProcessedAt)
	}
	if processed.EventID != "550e8400-e29b-41d4-a716-446655440000" {
		t.Fatalf("expected original event_id to be preserved, got %s", processed.EventID)
	}

}

func TestProcessMessage_InvalidMessage_ReturnsError(t *testing.T) {
	consumer := &fakeConsumer{}
	publisher := &fakePublisher{}
	processor := newTestProcessor(consumer, publisher)

	msg := entities.QueueMessage{
		Body:          `invalid json`,
		ReceiptHandle: "receipt-1",
	}

	err := processor.ProcessMessage(context.Background(), msg)
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if len(publisher.sent) != 0 {
		t.Fatalf("expected 0 published messages, got %d", len(publisher.sent))
	}

	if len(consumer.deleted) != 0 {
		t.Fatalf("expected 0 deleted messages, got %d", len(consumer.deleted))
	}
}

func TestProcessMessage_InvalidEvent_DoesNotPublishOrDelete(t *testing.T) {
	consumer := &fakeConsumer{}
	publisher := &fakePublisher{}
	processor := newTestProcessor(consumer, publisher)

	msg := validQueueMessage()
	msg.Body = `{
        "event_id": "550e8400-e29b-41d4-a716-446655440000",
        "developer_id": "",
        "metric_type": "commits",
        "value": 10,
        "repository": "org/repo",
        "timestamp": "2026-04-15T10:30:00Z"
    }`

	err := processor.ProcessMessage(context.Background(), msg)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if len(publisher.sent) != 0 {
		t.Fatalf("expected 0 published messages, got %d", len(publisher.sent))
	}
	if len(consumer.deleted) != 0 {
		t.Fatalf("expected 0 deleted messages, got %d", len(consumer.deleted))
	}
}

func TestProcessMessage_PublishError_DoesNotDelete(t *testing.T) {
	consumer := &fakeConsumer{}
	publisher := &fakePublisher{sendErr: errors.New("send failed")}
	processor := newTestProcessor(consumer, publisher)

	err := processor.ProcessMessage(context.Background(), validQueueMessage())
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if len(publisher.sent) != 1 {
		t.Fatalf("expected 1 publish attempt, got %d", len(publisher.sent))
	}
	if len(consumer.deleted) != 0 {
		t.Fatalf("expected 0 deleted messages, got %d", len(consumer.deleted))
	}
}

func TestProcessMessage_DeleteError_ReturnsError(t *testing.T) {
	consumer := &fakeConsumer{deleteErr: errors.New("delete failed")}
	publisher := &fakePublisher{}
	processor := newTestProcessor(consumer, publisher)

	err := processor.ProcessMessage(context.Background(), validQueueMessage())
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if len(publisher.sent) != 1 {
		t.Fatalf("expected 1 published message, got %d", len(publisher.sent))
	}
	if len(consumer.deleted) != 1 {
		t.Fatalf("expected 1 delete attempt, got %d", len(consumer.deleted))
	}
}
