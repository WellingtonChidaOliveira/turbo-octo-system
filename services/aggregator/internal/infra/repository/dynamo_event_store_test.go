package repository

import (
	"aggregator/internal/domain/entities"
	"aggregator/internal/infra/repository/records"
	"aggregator/internal/usecase/apperrors"
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

func TestProcessedEventRecordMapping(t *testing.T) {
	event := entities.ProcessedEvent{
		EventID:     "event-1",
		DeveloperID: "dev-123",
		MetricType:  entities.MetricCommits,
		Value:       10,
		Repository:  "org/repo",
		Timestamp:   "2026-04-15T10:30:00Z",
		ProcessedAt: "2026-04-15T10:30:05Z",
		ProcessorID: "processor-1",
	}

	record := records.NewProcessedEvent(event)
	got := record.ToEntity()

	if got != event {
		t.Fatalf("expected mapped event %+v, got %+v", event, got)
	}
}

func TestMapDynamoWriteError_DuplicateEvent(t *testing.T) {
	err := &types.TransactionCanceledException{
		CancellationReasons: []types.CancellationReason{
			{Code: aws.String("ConditionalCheckFailed")},
		},
	}

	if !errors.Is(mapDynamoWriteError(err), apperrors.ErrEventAlreadyProcessed) {
		t.Fatal("expected duplicate event error")
	}
}

func TestMapDynamoWriteError_ReturnsOriginalError(t *testing.T) {
	err := errors.New("boom")

	if !errors.Is(mapDynamoWriteError(err), err) {
		t.Fatal("expected original error")
	}
}
