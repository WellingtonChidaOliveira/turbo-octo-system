package main

import (
	"aggregator/internal/infra/api"
	"aggregator/internal/infra/config"
	"aggregator/internal/infra/queue"
	"aggregator/internal/infra/repository"
	"aggregator/internal/infra/runtime"
	"aggregator/internal/infra/worker"
	"aggregator/internal/usecase"
	"aggregator/internal/usecase/retry"
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})))

	cfg := config.LoadSettings()
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

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
		processedEventsQueue = queue.NewProcessedEventConsumer(queueClient, cfg.Queue.ProcessedQueueURL, cfg.Queue)
		eventStore           = repository.NewDynamoEventStore(dynamoClient, cfg.Dynamo.EventsTableName, cfg.Dynamo.DeveloperSummaryTableName)

		receivePolicy          = retryPolicy(cfg.ReceiveBackoff)
		processedEventConsumer = usecase.NewProcessedEventConsumer(processedEventsQueue, receivePolicy)
		persistData            = usecase.NewPersistProcessedEventHandler(eventStore, processedEventsQueue)
		eventsGetter           = usecase.NewGetProcessedEventsByDeveloper(eventStore)
		summaryGetter          = usecase.NewGetDeveloperSummary(eventStore)

		apiServer = api.NewServer(
			":"+cfg.API.Port,
			api.NewHealthHandler(processedEventsQueue, eventStore),
			api.NewMetricsHandler(eventsGetter),
			api.NewSummaryHandler(summaryGetter),
		)
		workerPool = worker.NewPool(persistData, cfg.Worker.Count)
		service    = runtime.NewService(processedEventConsumer, workerPool, apiServer, cfg)
	)

	service.Run(ctx, stop)
}

func retryPolicy(settings config.RetrySettings) retry.Policy {
	return retry.Policy{
		InitialBackoff: settings.InitialBackoff,
		MaxBackoff:     settings.MaxBackoff,
		Jitter:         settings.Jitter,
		MaxAttempts:    settings.MaxAttempts,
	}
}
