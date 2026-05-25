package queue

import (
	cfg "aggregator/internal/infra/config"
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

func NewQueueClient(ctx context.Context, cfg cfg.Settings) (*sqs.Client, error) {
	c, err := config.LoadDefaultConfig(
		ctx,
		config.WithRegion(cfg.Region),
		config.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider(
				cfg.KeysAwsAccessKeyId,
				cfg.KeysAwsSecretAccessKey,
				"",
			),
		),
	)

	if err != nil {
		return nil, err
	}

	client := sqs.NewFromConfig(c, func(o *sqs.Options) {
		o.BaseEndpoint = aws.String(cfg.AwsEndpointURL)
	})

	return client, nil
}
