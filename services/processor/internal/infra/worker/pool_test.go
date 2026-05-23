package worker

import (
	"context"
	"processor/internal/domain/entities"
	"sync"
	"testing"
	"time"
)

type fakeHandler struct {
	mu        sync.Mutex
	handled   []string
	handleErr error
}

func (h *fakeHandler) Handle(ctx context.Context, message entities.QueueMessage) error {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.handled = append(h.handled, message.ID)
	return h.handleErr
}

func (h *fakeHandler) handledCount() int {
	h.mu.Lock()
	defer h.mu.Unlock()
	return len(h.handled)
}

func TestPool_StartProcessesJobsUntilChannelCloses(t *testing.T) {
	handler := &fakeHandler{}
	pool := NewPool(handler, 2)
	jobs := make(chan entities.QueueMessage, 2)

	done := make(chan struct{})
	go func() {
		pool.Start(context.Background(), jobs)
		close(done)
	}()

	jobs <- entities.QueueMessage{ID: "message-1"}
	jobs <- entities.QueueMessage{ID: "message-2"}
	close(jobs)

	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatal("expected pool to stop after jobs channel closes")
	}

	if handler.handledCount() != 2 {
		t.Fatalf("expected 2 handled jobs, got %d", handler.handledCount())
	}
}

func TestPool_StartStopsWhenContextIsCanceled(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	handler := &fakeHandler{}
	pool := NewPool(handler, 2)
	jobs := make(chan entities.QueueMessage)

	done := make(chan struct{})
	go func() {
		pool.Start(ctx, jobs)
		close(done)
	}()

	cancel()

	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatal("expected pool to stop after context cancellation")
	}
}
