package queue

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

type Queue struct {
	Client   *sqs.Client
	QueueUrl string
}

type QueueMessage interface {
	Process(ctx context.Context)
}

func (q *Queue) Process(ctx context.Context) {
	for {
		out, err := q.Client.ReceiveMessage(ctx, &sqs.ReceiveMessageInput{
			QueueUrl:            aws.String(q.QueueUrl),
			MaxNumberOfMessages: 10,
			WaitTimeSeconds:     20,
			VisibilityTimeout:   30,
		})
		if err != nil {
			log.Fatal(err)
		}

		for _, msg := range out.Messages {
			log.Printf("Received message: %s", aws.ToString(msg.Body))

			_, err := q.Client.DeleteMessage(ctx, &sqs.DeleteMessageInput{
				QueueUrl:      aws.String(q.QueueUrl),
				ReceiptHandle: msg.ReceiptHandle,
			})
			if err != nil {
				log.Fatal(err)
			}
		}
	}
}
