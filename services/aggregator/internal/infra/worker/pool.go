package worker

import (
	"aggregator/internal/dto"
	"context"
	"log/slog"
	"sync"
)

type Handler interface {
	Handle(ctx context.Context, message dto.QueueMessage) error
}

type Pool struct {
	handler Handler
	size    int
}

func NewPool(handler Handler, size int) *Pool {
	if size <= 0 {
		size = 1
	}

	return &Pool{
		handler: handler,
		size:    size,
	}
}

func (p *Pool) Start(ctx context.Context, jobs <-chan dto.QueueMessage) {
	var wg sync.WaitGroup

	for workerID := 1; workerID <= p.size; workerID++ {
		wg.Add(1)
		go p.runWorker(ctx, workerID, jobs, &wg)
	}

	wg.Wait()
}

func (p *Pool) runWorker(ctx context.Context, workerID int, jobs <-chan dto.QueueMessage, wg *sync.WaitGroup) {
	defer wg.Done()

	for job := range jobs {
		p.handleJob(ctx, workerID, job)
	}
}

func (p *Pool) handleJob(ctx context.Context, workerID int, job dto.QueueMessage) {
	if err := p.handler.Handle(ctx, job); err != nil {
		slog.Error("worker failed to process message",
			"worker_id", workerID,
			"message_id", job.ID,
			"error", err,
		)
		return
	}

	slog.Info("worker processed message",
		"worker_id", workerID,
		"message_id", job.ID,
	)
}
