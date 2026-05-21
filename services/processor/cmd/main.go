package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"processor/internal/infra/config"
	"processor/internal/infra/queue"
	"processor/internal/infra/worker"
	"processor/internal/usecase"
)

func main() {

	// Set up context to handle graceful shutdown on interrupt signals
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	settings := config.LoadSettings()

	// Initialize AWS SQS client with custom settings
	client, err := queue.NewSQSClient(ctx, settings)
	if err != nil {
		log.Fatalf("unable to load AWS SDK config: %v", err)
	}

	var (
		publisher   = queue.NewProcessedEventsPublisher(client, settings.ProcessedQueueURL)
		processor   = usecase.NewProcessor(publisher, usecase.SystemClock{}, settings.ProcessorID)
		consumer    = queue.NewRawEventsConsumer(client, settings.RawQueueURL)
		pool        = worker.NewPool(settings.WorkerCount, processor)
		jobs        = make(chan worker.Job, settings.WorkerCount*2)
		workersDone = make(chan struct{})
	)

	log.Printf("processor started with %d workers", settings.WorkerCount)
	go func() {
		pool.Run(ctx, jobs)
		close(workersDone)
	}()

	consumer.Consume(ctx, jobs)
	close(jobs)
	<-workersDone
	log.Println("processor stopped")
}
