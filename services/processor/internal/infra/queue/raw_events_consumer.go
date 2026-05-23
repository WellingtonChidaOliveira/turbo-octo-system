package queue

import (
	"context"
	"processor/internal/domain/entities"
	"processor/internal/infra/config"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

type RawEventsConsumer struct {
	client   *sqs.Client
	queueUrl string
	settings config.QueueSettings
}

func NewRawEventsConsumer(client *sqs.Client, queueUrl string, settings config.QueueSettings) *RawEventsConsumer {
	return &RawEventsConsumer{
		client:   client,
		queueUrl: queueUrl,
		settings: settings,
	}
}

func (c *RawEventsConsumer) Receive(ctx context.Context) ([]entities.QueueMessage, error) {
	out, err := c.client.ReceiveMessage(ctx, &sqs.ReceiveMessageInput{
		QueueUrl:            aws.String(c.queueUrl),
		MaxNumberOfMessages: c.settings.MaxNumberOfMessages,
		WaitTimeSeconds:     c.settings.WaitTimeSeconds,
		VisibilityTimeout:   c.settings.VisibilityTimeout,
	})
	if err != nil {
		return nil, err
	}

	var queueMessages []entities.QueueMessage
	for _, msg := range out.Messages {
		queueMessages = append(queueMessages, entities.QueueMessage{
			ID:            aws.ToString(msg.MessageId),
			Body:          aws.ToString(msg.Body),
			ReceiptHandle: aws.ToString(msg.ReceiptHandle),
		})
	}

	return queueMessages, nil
}

func (c *RawEventsConsumer) Delete(ctx context.Context, receiptHandle string) error {
	_, err := c.client.DeleteMessage(ctx, &sqs.DeleteMessageInput{
		QueueUrl:      aws.String(c.queueUrl),
		ReceiptHandle: aws.String(receiptHandle),
	})
	return err
}
