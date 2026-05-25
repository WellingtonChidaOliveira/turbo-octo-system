package queue

import (
	"aggregator/internal/dto"
	"aggregator/internal/infra/config"
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
)

type ProcessedEventConsumer struct {
	client   *sqs.Client
	queueURL string
	cfg      config.QueueSettings
}

func NewProcessedEventConsumer(
	client *sqs.Client,
	queueURL string,
	cfg config.QueueSettings) *ProcessedEventConsumer {
	return &ProcessedEventConsumer{
		client:   client,
		queueURL: queueURL,
		cfg:      cfg,
	}
}

func (c *ProcessedEventConsumer) Receive(ctx context.Context) ([]dto.QueueMessage, error) {
	out, err := c.client.ReceiveMessage(ctx, &sqs.ReceiveMessageInput{
		QueueUrl:            aws.String(c.queueURL),
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

func (c *ProcessedEventConsumer) Delete(ctx context.Context, receiptHandle string) error {
	_, err := c.client.DeleteMessage(ctx, &sqs.DeleteMessageInput{
		QueueUrl:      aws.String(c.queueURL),
		ReceiptHandle: aws.String(receiptHandle),
	})
	return err
}

func (c *ProcessedEventConsumer) CheckQueue(ctx context.Context) error {
	_, err := c.client.GetQueueAttributes(ctx, &sqs.GetQueueAttributesInput{
		QueueUrl: aws.String(c.queueURL),
		AttributeNames: []types.QueueAttributeName{
			types.QueueAttributeNameQueueArn,
		},
	})
	return err
}
