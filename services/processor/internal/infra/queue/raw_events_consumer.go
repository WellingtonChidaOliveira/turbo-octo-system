package queue

import (
	"context"
	"processor/internal/domain/abstractions"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

type RawEventsConsumer struct {
	client   *sqs.Client
	queueUrl string
}

func NewRawEventsConsumer(client *sqs.Client, queueUrl string) *RawEventsConsumer {
	return &RawEventsConsumer{
		client:   client,
		queueUrl: queueUrl,
	}
}

func (c *RawEventsConsumer) Receive(ctx context.Context) ([]abstractions.QueueMessage, error) {
	out, err := c.client.ReceiveMessage(ctx, &sqs.ReceiveMessageInput{
		QueueUrl:            aws.String(c.queueUrl),
		MaxNumberOfMessages: 10,
		WaitTimeSeconds:     20,
		VisibilityTimeout:   30,
	})
	if err != nil {
		return nil, err
	}

	var queueMessages []abstractions.QueueMessage
	for _, msg := range out.Messages {
		queueMessages = append(queueMessages, abstractions.QueueMessage{
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
