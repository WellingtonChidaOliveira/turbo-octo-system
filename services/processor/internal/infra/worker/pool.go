package worker

import (
	"context"
	"log"
	"processor/internal/domain/entities"
	"sync"
)

type Handler interface {
	Handle(ctx context.Context, message entities.QueueMessage) error
}

type Pool struct {
	handler Handler
	size    int
}

func NewPool(handler Handler, size int) *Pool {
	return &Pool{
		handler: handler,
		size:    size,
	}
}

func (p *Pool) Start(ctx context.Context, jobs <-chan entities.QueueMessage) {
	var wg sync.WaitGroup

	for workerId := 1; workerId <= p.size; workerId++ {
		wg.Add(1)
		go p.runWorker(ctx, workerId, jobs, &wg)
	}

	wg.Wait()
}

func (p *Pool) runWorker(ctx context.Context, workerId int, jobs <-chan entities.QueueMessage, wg *sync.WaitGroup) {
	defer wg.Done()

	for {
		select {
		case <-ctx.Done():
			return
		case job, ok := <-jobs:
			if !ok {
				return
			}

			p.handleJob(ctx, workerId, job)
		}
	}
}

func (p *Pool) handleJob(ctx context.Context, workerId int, job entities.QueueMessage) {
	if err := p.handler.Handle(ctx, job); err != nil {
		log.Printf("worker=%d message_id=%s err=%v", workerId, job.ID, err)
		return
	}

	log.Printf("worker=%d message_id=%s processed successfully", workerId, job.ID)
}
