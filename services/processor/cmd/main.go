package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"processor/internal/domain/entities"
	"processor/internal/infra/config"
	"processor/internal/infra/queue"
	"processor/internal/infra/worker"
	"processor/internal/usecase"
	"syscall"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	settings := config.LoadSettings()

	queueClient, err := queue.NewQueueClient(ctx, settings)
	if err != nil {
		panic(err)
	}

	var (
		eventPublisher = queue.NewProcessedEventsPublisher(queueClient, settings.ProcessedQueueURL)
		eventConsumer  = queue.NewRawEventsConsumer(queueClient, settings.RawQueueURL)
		consumer       = usecase.NewRawEventsConsumer(eventConsumer) // TODO: Remove this line after implementing the consumer interface
		processor      = usecase.NewEventProcessor(eventPublisher, eventConsumer, entities.SystemClock{}, settings.ProcessorID)
		pool           = worker.NewPool(processor, settings.WorkerCount) // TODO: Replace nil with the actual handler implementation
		jobs           = make(chan entities.QueueMessage, settings.WorkerCount*2)
		workersDone    = make(chan struct{})
	)

	log.Printf("Starting processor with ID %s", settings.ProcessorID)
	go func() {
		pool.Start(ctx, jobs)
		close(workersDone)
	}()

	consumer.Consumer(ctx, jobs)
	close(jobs)
	<-workersDone

	log.Printf("Processor with ID %s is shutting down", settings.ProcessorID)

}
