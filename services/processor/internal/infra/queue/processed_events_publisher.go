package queue

import (
	"context"
	"processor/internal/domain/abstractions"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

type ProcessedEventsPublisher struct {
	client   *sqs.Client
	queueUrl string
}

func NewProcessedEventsPublisher(client *sqs.Client, queueUrl string) *ProcessedEventsPublisher {
	return &ProcessedEventsPublisher{
		client:   client,
		queueUrl: queueUrl,
	}
}

func (p *ProcessedEventsPublisher) Send(ctx context.Context, message abstractions.QueueMessage) error {
	_, err := p.client.SendMessage(ctx, &sqs.SendMessageInput{
		QueueUrl:    aws.String(p.queueUrl),
		MessageBody: aws.String(message.Body),
	})
	return err
}
