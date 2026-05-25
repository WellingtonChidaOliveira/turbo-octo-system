package main

import (
	"aggregator/internal/dto"
	"aggregator/internal/infra/config"
	"aggregator/internal/infra/queue"
	"aggregator/internal/infra/repository"
	"aggregator/internal/infra/worker"
	"aggregator/internal/usecase"
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	cfg := config.LoadSettings()

	queueClient, err := queue.NewQueueClient(ctx, cfg)
	if err != nil {
		slog.Error("failed to create queue client", "error", err)
		os.Exit(1)
	}

	dynamoClient, err := repository.NewDynamoClient(ctx, cfg)
	if err != nil {
		slog.Error("failed to create dynamodb client", "error", err)
		os.Exit(1)
	}

	var (
		eventConsumer = queue.NewProcessedEventsConsumer(queueClient, cfg.ProcessedQueueURL, cfg.Queue)
		consumer      = usecase.NewProcessedEventConsumer(eventConsumer)
		eventStore    = repository.NewDynamoEventStore(dynamoClient, cfg.EventsTableName, cfg.DeveloperSummaryTableName)
		persistData   = usecase.NewPersistDataHandler(eventStore, eventConsumer)
		workerPool    = worker.NewPool(persistData, cfg.WorkerCount)
		jobs          = make(chan dto.QueueMessage, cfg.JobBufferSize)
	)

	go consumer.Consume(ctx, jobs)
	workerPool.Start(ctx, jobs)
}
