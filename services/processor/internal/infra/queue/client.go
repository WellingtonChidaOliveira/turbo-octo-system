package queue

import (
	"context"
	conf "processor/internal/infra/config"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

type SQSQueueClient struct {
	client   *sqs.Client
	queueUrl string
}

func NewQueueClient(ctx context.Context, setting conf.Settings) (*sqs.Client, error) {
	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithRegion(setting.Region),
		config.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider(setting.KeysAwsAccessKeyId, setting.KeysAwsSecretAccessKey, ""),
		),
	)

	if err != nil {
		return nil, err
	}

	client := sqs.NewFromConfig(cfg, func(o *sqs.Options) {
		o.BaseEndpoint = aws.String(setting.QueueUrl)
	})
	return client, nil
}
