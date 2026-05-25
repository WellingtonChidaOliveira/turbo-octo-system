package queue

import (
	"context"
	"processor/internal/dto"
	"processor/internal/infra/config"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

type RawEventConsumer struct {
	client   *sqs.Client
	queueURL string
	settings config.QueueSettings
}

func NewRawEventConsumer(client *sqs.Client, queueURL string, settings config.QueueSettings) *RawEventConsumer {
	return &RawEventConsumer{
		client:   client,
		queueURL: queueURL,
		settings: settings,
	}
}

func (c *RawEventConsumer) Receive(ctx context.Context) ([]dto.QueueMessage, error) {
	out, err := c.client.ReceiveMessage(ctx, &sqs.ReceiveMessageInput{
		QueueUrl:            aws.String(c.queueURL),
		MaxNumberOfMessages: c.settings.MaxNumberOfMessages,
		WaitTimeSeconds:     c.settings.WaitTimeSeconds,
		VisibilityTimeout:   c.settings.VisibilityTimeout,
	})
	if err != nil {
		return nil, err
	}

	var queueMessages []dto.QueueMessage
	for _, msg := range out.Messages {
		queueMessages = append(queueMessages, dto.QueueMessage{
			ID:            aws.ToString(msg.MessageId),
			Body:          aws.ToString(msg.Body),
			ReceiptHandle: aws.ToString(msg.ReceiptHandle),
		})
	}

	return queueMessages, nil
}

func (c *RawEventConsumer) Delete(ctx context.Context, receiptHandle string) error {
	_, err := c.client.DeleteMessage(ctx, &sqs.DeleteMessageInput{
		QueueUrl:      aws.String(c.queueURL),
		ReceiptHandle: aws.String(receiptHandle),
	})
	return err
}
