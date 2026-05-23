package usecase

import (
	"context"
	"encoding/json"
	"log"
	"processor/internal/domain/entities"
	"processor/internal/usecase/ports"
)

type EventProcessor struct {
	publisher   ports.QueuePublisher
	deleter     ports.QueueDeleter
	clock       ports.Clock
	processorID string
}

func NewEventProcessor(
	publisher ports.QueuePublisher,
	deleter ports.QueueDeleter,
	clock ports.Clock,
	processorID string) *EventProcessor {
	return &EventProcessor{
		publisher:   publisher,
		deleter:     deleter,
		processorID: processorID,
		clock:       clock,
	}
}

func (p *EventProcessor) Handle(ctx context.Context, message entities.QueueMessage) error {
	return p.ProcessMessage(ctx, message)
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

	err = p.deleter.Delete(ctx, msg.ReceiptHandle)
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
