package queue

import (
	"aggregator/internal/dto"
	"aggregator/internal/infra/config"
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

type ProcessedEventsConsumer struct {
	client   *sqs.Client
	queueUrl string
	cfg      config.QueueSettings
}

func NewProcessedEventsConsumer(
	client *sqs.Client,
	queueUrl string,
	cfg config.QueueSettings) *ProcessedEventsConsumer {
	return &ProcessedEventsConsumer{
		client:   client,
		queueUrl: queueUrl,
		cfg:      cfg,
	}
}

func (c *ProcessedEventsConsumer) Receive(ctx context.Context) ([]dto.QueueMessage, error) {
	out, err := c.client.ReceiveMessage(ctx, &sqs.ReceiveMessageInput{
		QueueUrl:            aws.String(c.queueUrl),
		MaxNumberOfMessages: c.cfg.MaxNumberOfMessages,
		WaitTimeSeconds:     c.cfg.WaitTimeSeconds,
		VisibilityTimeout:   c.cfg.VisibilityTimeout,
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

func (c *ProcessedEventsConsumer) Delete(ctx context.Context, receiptHandle string) error {
	_, err := c.client.DeleteMessage(ctx, &sqs.DeleteMessageInput{
		QueueUrl:      aws.String(c.queueUrl),
		ReceiptHandle: aws.String(receiptHandle),
	})
	return err
}
