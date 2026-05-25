package usecase

import (
	"aggregator/internal/dto"
	"aggregator/internal/usecase/retry"
	"context"
	"testing"
)

func TestProcessedEventConsumer_Start_ClosesJobsWhenContextCanceled(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	jobs := make(chan dto.QueueMessage)
	consumer := NewProcessedEventConsumer(&fakeQueueConsumer{}, retry.Policy{MaxAttempts: 1})

	consumer.Start(ctx, jobs)

	if _, ok := <-jobs; ok {
		t.Fatal("expected jobs channel to be closed")
	}
}

type fakeQueueConsumer struct{}

func (f *fakeQueueConsumer) Receive(ctx context.Context) ([]dto.QueueMessage, error) {
	return nil, ctx.Err()
}

func (f *fakeQueueConsumer) Delete(ctx context.Context, receiptHandle string) error {
	return nil
}
