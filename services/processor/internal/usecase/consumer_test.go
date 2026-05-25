package usecase

import (
	"context"
	"errors"
	"processor/internal/dto"
	"processor/internal/usecase/retry"
	"testing"
	"time"
)

type scriptedConsumer struct {
	receiveCalls int
	batches      [][]dto.QueueMessage
	errs         []error
}

func (c *scriptedConsumer) Receive(ctx context.Context) ([]dto.QueueMessage, error) {
	call := c.receiveCalls
	c.receiveCalls++
	if call < len(c.errs) && c.errs[call] != nil {
		return nil, c.errs[call]
	}
	if call < len(c.batches) {
		return c.batches[call], nil
	}
	<-ctx.Done()
	return nil, ctx.Err()
}

func (c *scriptedConsumer) Delete(ctx context.Context, receiptHandle string) error {
	return nil
}

func TestRawEventConsumer_DispatchesMessagesAndClosesJobsOnCancel(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	consumer := &scriptedConsumer{
		batches: [][]dto.QueueMessage{
			{{ID: "message-1"}, {ID: "message-2"}},
		},
	}
	rawConsumer := NewRawEventsConsumer(consumer, retry.Policy{MaxAttempts: 1})
	jobs := make(chan dto.QueueMessage, 2)

	done := make(chan struct{})
	go func() {
		rawConsumer.Consumer(ctx, jobs)
		close(done)
	}()

	for len(jobs) < 2 {
		time.Sleep(time.Millisecond)
	}
	cancel()

	select {
	case <-done:
	case <-time.After(time.Second):
		t.Fatal("expected consumer to stop after cancellation")
	}

	count := 0
	for range jobs {
		count++
	}
	if count != 2 {
		t.Fatalf("expected 2 dispatched messages, got %d", count)
	}
}

func TestRawEventConsumer_RetriesReceiveErrorWithBackoff(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	consumer := &scriptedConsumer{
		errs: []error{errors.New("receive failed")},
		batches: [][]dto.QueueMessage{
			nil,
			{{ID: "message-1"}},
		},
	}
	rawConsumer := NewRawEventsConsumer(consumer, retry.Policy{
		InitialBackoff: 0,
		MaxBackoff:     0,
		MaxAttempts:    1,
	})
	jobs := make(chan dto.QueueMessage, 1)

	done := make(chan struct{})
	go func() {
		rawConsumer.Consumer(ctx, jobs)
		close(done)
	}()

	select {
	case msg := <-jobs:
		if msg.ID != "message-1" {
			t.Fatalf("expected message-1, got %s", msg.ID)
		}
	case <-time.After(time.Second):
		t.Fatal("expected message after receive retry")
	}

	cancel()
	select {
	case <-done:
	case <-time.After(time.Second):
		t.Fatal("expected consumer to stop after cancellation")
	}
}
