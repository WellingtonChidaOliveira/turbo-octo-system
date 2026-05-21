package worker

import (
	"context"
	"log"
	"sync"
)

type Handler interface {
	Handle(ctx context.Context, body []byte) error
}

type AckFunc func(ctx context.Context) error

type Job struct {
	ID   string
	Body []byte
	// Ack confirms the message only after the handler succeeds, preserving SQS redrive behavior.
	Ack AckFunc
}

type Pool struct {
	handler Handler
	size    int
}

func NewPool(size int, handler Handler) *Pool {
	return &Pool{
		handler: handler,
		size:    size,
	}
}

func (p *Pool) Run(ctx context.Context, jobs <-chan Job) {
	var wg sync.WaitGroup

	for workerID := 1; workerID <= p.size; workerID++ {
		wg.Add(1)
		go p.runWorker(ctx, workerID, jobs, &wg)
	}

	wg.Wait()
}

func (p *Pool) runWorker(ctx context.Context, workerID int, jobs <-chan Job, wg *sync.WaitGroup) {
	defer wg.Done()

	for {
		select {
		case <-ctx.Done():
			return
		case job, ok := <-jobs:
			if !ok {
				return
			}

			p.handleJob(ctx, workerID, job)
		}
	}
}

func (p *Pool) handleJob(ctx context.Context, workerID int, job Job) {
	if err := p.handler.Handle(ctx, job.Body); err != nil {
		log.Printf("worker=%d message_id=%s failed: %v", workerID, job.ID, err)
		return
	}

	if err := job.Ack(ctx); err != nil {
		log.Printf("worker=%d message_id=%s ack failed: %v", workerID, job.ID, err)
		return
	}

	log.Printf("worker=%d message_id=%s processed", workerID, job.ID)
}
