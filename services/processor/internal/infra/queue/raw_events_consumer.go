package queue

import (
	"context"
	"log"
	"time"

	"processor/internal/infra/worker"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
)

type RawEventsConsumer struct {
	client   *sqs.Client
	queueURL string
}

func NewRawEventsConsumer(client *sqs.Client, queueURL string) *RawEventsConsumer {
	return &RawEventsConsumer{
		client:   client,
		queueURL: queueURL,
	}
}

func (c *RawEventsConsumer) Consume(ctx context.Context, jobs chan<- worker.Job) {
	for ctx.Err() == nil {
		messages, err := c.receive(ctx)
		if err != nil {
			c.handleReceiveError(ctx, err)
			continue
		}

		c.dispatch(ctx, jobs, messages)
	}
}

func (c *RawEventsConsumer) receive(ctx context.Context) ([]types.Message, error) {
	output, err := c.client.ReceiveMessage(ctx, &sqs.ReceiveMessageInput{
		QueueUrl:            aws.String(c.queueURL),
		MaxNumberOfMessages: 10,
		WaitTimeSeconds:     20,
		VisibilityTimeout:   30,
	})
	if err != nil {
		return nil, err
	}

	return output.Messages, nil
}

func (c *RawEventsConsumer) dispatch(ctx context.Context, jobs chan<- worker.Job, messages []types.Message) {
	for _, message := range messages {
		select {
		case jobs <- c.toJob(message):
		case <-ctx.Done():
			return
		}
	}
}

func (c *RawEventsConsumer) toJob(message types.Message) worker.Job {
	return worker.Job{
		ID:   aws.ToString(message.MessageId),
		Body: []byte(aws.ToString(message.Body)),
		Ack: func(ctx context.Context) error {
			return c.delete(ctx, message)
		},
	}
}

func (c *RawEventsConsumer) delete(ctx context.Context, message types.Message) error {
	_, err := c.client.DeleteMessage(ctx, &sqs.DeleteMessageInput{
		QueueUrl:      aws.String(c.queueURL),
		ReceiptHandle: message.ReceiptHandle,
	})
	return err
}

func (c *RawEventsConsumer) handleReceiveError(ctx context.Context, err error) {
	if ctx.Err() != nil {
		return
	}

	log.Printf("failed to receive raw events: %v", err)
	time.Sleep(time.Second)
}
