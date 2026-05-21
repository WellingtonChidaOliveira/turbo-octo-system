package usecase

import (
	"context"
	"encoding/json"
	"time"

	"processor/internal/domain"
)

type ProcessedEventPublisher interface {
	PublishProcessedEvent(ctx context.Context, event domain.ProcessedEvent) error
}

type Clock interface {
	Now() time.Time
}

type Processor struct {
	publisher   ProcessedEventPublisher
	clock       Clock
	processorID string
}

func NewProcessor(publisher ProcessedEventPublisher, clock Clock, processorID string) *Processor {
	return &Processor{
		publisher:   publisher,
		clock:       clock,
		processorID: processorID,
	}
}

func (p *Processor) Handle(ctx context.Context, body []byte) error {
	event, err := decodeRawEvent(body)
	if err != nil {
		return err
	}

	if err := event.Validate(p.clock.Now()); err != nil {
		return err
	}

	processed := event.ToProcessed(p.processorID, p.clock.Now())
	return p.publisher.PublishProcessedEvent(ctx, processed)
}

func decodeRawEvent(body []byte) (domain.RawEvent, error) {
	var event domain.RawEvent
	err := json.Unmarshal(body, &event)
	return event, err
}

type SystemClock struct{}

func (SystemClock) Now() time.Time {
	return time.Now().UTC()
}
