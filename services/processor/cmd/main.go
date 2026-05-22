package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"processor/internal/domain/entities"
	"processor/internal/infra/config"
	"processor/internal/infra/queue"
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
		publisher = queue.NewProcessedEventsPublisher(queueClient, settings.ProcessedQueueURL)
		consumer  = queue.NewRawEventsConsumer(queueClient, settings.RawQueueURL)
		processor = usecase.NewEventProcessor(consumer, publisher, entities.SystemClock{}, settings.ProcessorID)
	)

	if err := processor.Handle(ctx); err != nil {
		log.Fatal(err)
	}

}
