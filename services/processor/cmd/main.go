package main

import (
	"context"
	"processor/internal/infra/config"
	"processor/internal/infra/queue"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	endpoint := "http://localhost:4566"
	queueUrl := "/000000000000/raw-events"

	client := config.DefaultAWSConfigResolvers(ctx, endpoint)

	qe := queue.Queue{
		Client:   client,
		QueueUrl: queueUrl,
	}

	qe.Process(ctx)
}
