package main

import (
	"aggregator/internal/dto"
	"aggregator/internal/infra/config"
	"aggregator/internal/infra/queue"
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

	var (
		eventConsumer = queue.NewProcessedEventsConsumer(queueClient, cfg.ProcessedQueueURL, cfg.Queue)
		consumer      = usecase.NewProcessedEventConsumer(eventConsumer)
		jobs          = make(chan dto.QueueMessage, cfg.JobBufferSize)
	)

	consumer.Consume(ctx, jobs)
}
