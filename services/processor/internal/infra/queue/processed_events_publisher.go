package queue

import (
	"context"
	"encoding/json"
	"time"

	"processor/internal/domain"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
)

type ProcessedEventsPublisher struct {
	client   *sqs.Client
	queueURL string
}

func NewProcessedEventsPublisher(client *sqs.Client, queueURL string) *ProcessedEventsPublisher {
	return &ProcessedEventsPublisher{
		client:   client,
		queueURL: queueURL,
	}
}

func (p *ProcessedEventsPublisher) PublishProcessedEvent(ctx context.Context, event domain.ProcessedEvent) error {
	body, err := json.Marshal(event)
	if err != nil {
		return err
	}

	return p.sendWithBackoff(ctx, body, event.EventID)
}

func (p *ProcessedEventsPublisher) sendWithBackoff(ctx context.Context, body []byte, eventID string) error {
	var lastErr error

	for attempt := 0; attempt < 3; attempt++ {
		if err := p.send(ctx, body, eventID); err == nil {
			return nil
		} else {
			lastErr = err
		}

		if err := waitBackoff(ctx, attempt); err != nil {
			return err
		}
	}

	return lastErr
}

func (p *ProcessedEventsPublisher) send(ctx context.Context, body []byte, eventID string) error {
	_, err := p.client.SendMessage(ctx, &sqs.SendMessageInput{
		QueueUrl:    aws.String(p.queueURL),
		MessageBody: aws.String(string(body)),
		MessageAttributes: map[string]types.MessageAttributeValue{
			"event_id": {
				DataType:    aws.String("String"),
				StringValue: aws.String(eventID),
			},
		},
	})
	return err
}

func waitBackoff(ctx context.Context, attempt int) error {
	delay := time.Duration(1<<attempt) * time.Second

	select {
	case <-time.After(delay):
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}
