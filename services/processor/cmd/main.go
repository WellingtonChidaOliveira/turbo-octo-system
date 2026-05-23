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
	"processor/internal/usecase/retry"
	"syscall"
	"time"
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
		policy         = retry.NewPolicy(settings)
		eventPublisher = queue.NewProcessedEventsPublisher(queueClient, settings.ProcessedQueueURL)
		eventConsumer  = queue.NewRawEventsConsumer(queueClient, settings.RawQueueURL, settings.Queue)
		consumer       = usecase.NewRawEventsConsumer(eventConsumer, policy)
		processor      = usecase.NewEventProcessor(
			eventPublisher,
			eventConsumer,
			entities.SystemClock{},
			settings.ProcessorID,
			policy)
		pool        = worker.NewPool(processor, settings.WorkerCount)
		jobs        = make(chan entities.QueueMessage, settings.JobBufferSize)
		workersDone = make(chan struct{})
	)

	slog.Info("processor starting",
		"processor_id", settings.ProcessorID,
		"worker_count", settings.WorkerCount,
		"job_buffer_size", settings.JobBufferSize,
	)
	go func() {
		pool.Start(ctx, jobs)
		close(workersDone)
	}()

	consumer.Consumer(ctx, jobs)
	select {
	case <-workersDone:
	case <-time.After(settings.ShutdownTimeout):
		slog.Warn("processor shutdown timeout exceeded",
			"processor_id", settings.ProcessorID,
			"shutdown_timeout", settings.ShutdownTimeout.String(),
		)
	}

	slog.Info("processor stopped", "processor_id", settings.ProcessorID)

}
