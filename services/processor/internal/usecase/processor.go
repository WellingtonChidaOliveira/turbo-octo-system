package usecase

import (
	"context"
	"encoding/json"
	"log"
	"processor/internal/domain/entities"
	"processor/internal/usecase/ports"
)

// TODO: Receive messages from the queue, process them, and delete them from the queue after processing.
type EventProcessor struct {
	consumer    ports.QueueConsumer
	publisher   ports.QueuePublisher
	clock       ports.Clock
	processorID string
}

func NewEventProcessor(
	consumer ports.QueueConsumer,
	publisher ports.QueuePublisher,
	clock ports.Clock,
	processorID string) *EventProcessor {
	return &EventProcessor{
		publisher:   publisher,
		consumer:    consumer,
		processorID: processorID,
		clock:       clock,
	}
}

func (p *EventProcessor) Handle(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			log.Println("Shutting down event processor...")
			return nil
		default:
			messages, err := p.consumer.Receive(ctx)
			if err != nil {
				log.Printf("Error receiving messages: %v", err)
				continue
			}
			for _, msg := range messages {
				if err := p.ProcessMessage(ctx, msg); err != nil {
					log.Printf("Error processing message: %v", err)
					continue
				}
			}
		}

	}
}

func (p *EventProcessor) ProcessMessage(ctx context.Context, msg entities.QueueMessage) error {
	processedMsg, err := p.trateMessage(msg)
	if err != nil {
		log.Printf("Error treating message: %v", err)
		return err
	}
	log.Printf("Received message: %s", msg.Body)

	err = p.publisher.Send(ctx, entities.QueueMessage{Body: string(processedMsg)})
	if err != nil {
		log.Printf("Error sending message: %v", err)
		return err
	}

	err = p.consumer.Delete(ctx, msg.ReceiptHandle)
	if err != nil {
		log.Printf("Error deleting message: %v", err)
		return err
	}

	return nil
}

func (p *EventProcessor) trateMessage(msg entities.QueueMessage) (string, error) {
	var event entities.RawEvent
	if err := json.Unmarshal([]byte(msg.Body), &event); err != nil {
		log.Printf("Error unmarshaling message: %v", err)
		return "", err
	}

	if err := event.Validate(p.clock.Now()); err != nil {
		log.Printf("Invalid event: %v", err)
		return "", err
	}
	//TODO: system clock
	processEvent := event.ToProcessed(p.processorID, p.clock.Now())
	processedMsg, err := json.Marshal(processEvent)
	if err != nil {
		log.Printf("Error marshaling to processed event: %v", err)
		return "", err
	}

	log.Printf("validated message: %s", msg.Body)

	return string(processedMsg), nil
}
