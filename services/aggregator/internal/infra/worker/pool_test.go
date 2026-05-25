package worker

import (
	"aggregator/internal/dto"
	"context"
	"sync"
	"testing"
	"time"
)

func TestPool_Start_ReturnsWhenJobsChannelCloses(t *testing.T) {
	pool := NewPool(&fakeHandler{}, 1)
	jobs := make(chan dto.QueueMessage)
	done := make(chan struct{})

	go func() {
		pool.Start(context.Background(), jobs)
		close(done)
	}()
	close(jobs)

	select {
	case <-done:
	case <-time.After(time.Second):
		t.Fatal("expected pool to stop after jobs channel closes")
	}
}

type fakeHandler struct {
	mu      sync.Mutex
	handled int
}

func (h *fakeHandler) Handle(ctx context.Context, message dto.QueueMessage) error {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.handled++
	return nil
}
