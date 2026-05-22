package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"processor/internal/domain/usecase"
	"processor/internal/infra/config"
	"processor/internal/infra/queue"
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
		publisher = queue.NewProcessedEventsPublisher(queueClient, settings.ProcessedQueueURL)
		consumer  = queue.NewRawEventsConsumer(queueClient, settings.RawQueueURL)
		processor = usecase.NewEventProcessor(consumer, publisher, settings.ProcessorID)
	)

	if err := processor.Handle(ctx); err != nil {
		log.Fatal(err)
	}

}
