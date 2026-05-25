package queue

import (
	"context"
	"processor/internal/dto"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
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

func (p *ProcessedEventsPublisher) Send(ctx context.Context, message dto.QueueMessage) error {
	_, err := p.client.SendMessage(ctx, &sqs.SendMessageInput{
		QueueUrl:    aws.String(p.queueURL),
		MessageBody: aws.String(message.Body),
	})
	return err
}
