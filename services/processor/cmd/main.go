package main

import (
	"context"
	"log/slog"
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
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})))

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	settings := config.LoadSettings()

	queueClient, err := queue.NewQueueClient(ctx, settings)
	if err != nil {
		slog.Error("failed to create queue client", "error", err)
		os.Exit(1)
	}

	var (
		eventPublisher = queue.NewProcessedEventsPublisher(queueClient, settings.ProcessedQueueURL)
		eventConsumer  = queue.NewRawEventsConsumer(queueClient, settings.RawQueueURL)
		consumer       = usecase.NewRawEventsConsumer(eventConsumer) // TODO: Remove this line after implementing the consumer interface
		processor      = usecase.NewEventProcessor(eventPublisher, eventConsumer, entities.SystemClock{}, settings.ProcessorID)
		pool           = worker.NewPool(processor, settings.WorkerCount)
		jobs           = make(chan entities.QueueMessage, settings.WorkerCount*2)
		workersDone    = make(chan struct{})
	)

	slog.Info("processor starting",
		"processor_id", settings.ProcessorID,
		"worker_count", settings.WorkerCount,
	)
	go func() {
		pool.Start(ctx, jobs)
		close(workersDone)
	}()

	consumer.Consumer(ctx, jobs)
	<-workersDone

	slog.Info("processor stopped", "processor_id", settings.ProcessorID)

}
