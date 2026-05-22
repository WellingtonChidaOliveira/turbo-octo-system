package usecase

import (
	"context"
	"log"
	"processor/internal/domain/abstractions"
	"time"
)

// TODO: Receive messages from the queue, process them, and delete them from the queue after processing.
type EventProcessor struct {
	consumer   abstractions.QueueConsumerInterface
	publisher  abstractions.QueuePublisherInterface
	clock      time.Time
	processoID string
}

func NewEventProcessor(
	consumer abstractions.QueueConsumerInterface,
	publisher abstractions.QueuePublisherInterface,
	processorID string) *EventProcessor {
	return &EventProcessor{
		publisher:  publisher,
		consumer:   consumer,
		processoID: processorID,
	}
}

func (p *EventProcessor) Handle(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			log.Println("Shutting down event processor...")
			return nil
		default:
		}

		messages, err := p.consumer.Receive(ctx)
		if err != nil {
			log.Printf("Error receiving messages: %v", err)
			continue
		}

		for _, msg := range messages {
			log.Printf("Received message: %s", msg.Body)
			// Process the message (for demonstration, we just log it)

			// After processing, delete the message from the queue
			err := p.publisher.Send(ctx, abstractions.QueueMessage{Body: msg.Body})
			if err != nil {
				log.Printf("Error sending message: %v", err)
			}

			err = p.consumer.Delete(ctx, msg.ReceiptHandle)
			if err != nil {
				log.Printf("Error deleting message: %v", err)
			}
		}
	}
}
