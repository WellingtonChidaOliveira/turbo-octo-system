package worker

import (
	"context"
	"log/slog"
	"processor/internal/dto"
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

	for workerId := 1; workerId <= p.size; workerId++ {
		wg.Add(1)
		go p.runWorker(ctx, workerId, jobs, &wg)
	}

	wg.Wait()
}

func (p *Pool) runWorker(ctx context.Context, workerId int, jobs <-chan dto.QueueMessage, wg *sync.WaitGroup) {
	defer wg.Done()

	for job := range jobs {
		p.handleJob(ctx, workerId, job)
	}
}

func (p *Pool) handleJob(ctx context.Context, workerId int, job dto.QueueMessage) {
	if err := p.handler.Handle(ctx, job); err != nil {
		slog.Error("worker failed to process message",
			"worker_id", workerId,
			"message_id", job.ID,
			"error", err,
		)
		return
	}

	slog.Info("worker processed message",
		"worker_id", workerId,
		"message_id", job.ID,
	)
}
