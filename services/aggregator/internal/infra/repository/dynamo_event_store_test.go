package repository

import (
	"aggregator/internal/domain/entities"
	"aggregator/internal/infra/repository/records"
	"aggregator/internal/usecase/apperrors"
	"errors"
	"strings"
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

func TestDeveloperSummaryRecordMapping(t *testing.T) {
	summary := entities.DeveloperSummary{
		DeveloperID:            "dev-001",
		TotalCommits:           12,
		TotalPullRequests:      3,
		TotalReviewTimeMinutes: 90,
		ReviewTimeEvents:       2,
		EventsProcessed:        5,
		LastActivity:           "2026-04-15T10:30:00Z",
	}

	record := records.DeveloperSummary{
		DeveloperID:            summary.DeveloperID,
		TotalCommits:           summary.TotalCommits,
		TotalPullRequests:      summary.TotalPullRequests,
		TotalReviewTimeMinutes: summary.TotalReviewTimeMinutes,
		ReviewTimeEvents:       summary.ReviewTimeEvents,
		EventsProcessed:        summary.EventsProcessed,
		LastActivity:           summary.LastActivity,
	}

	if got := record.ToEntity(); got != summary {
		t.Fatalf("expected mapped summary %+v, got %+v", summary, got)
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

func TestBuildSaveEventAndSummaryInput_DoesNotOverwriteLastActivity(t *testing.T) {
	input, err := buildSaveEventAndSummaryInput("events", "developer_summary", entities.ProcessedEvent{
		EventID:     "2f4db003-dbc5-4b1b-bf7b-1f847b69ce1f",
		DeveloperID: "dev-123",
		MetricType:  entities.MetricCommits,
		Value:       10,
		Repository:  "org/repo",
		Timestamp:   "2026-04-15T10:30:00Z",
		ProcessedAt: "2026-04-15T10:30:05Z",
		ProcessorID: "processor-1",
	})
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	update := input.TransactItems[1].Update
	if strings.Contains(aws.ToString(update.UpdateExpression), "last_activity") {
		t.Fatalf("expected summary increment not to overwrite last_activity, got %s", aws.ToString(update.UpdateExpression))
	}
}

func TestMapDynamoConditionalUpdateError_IgnoresOlderLastActivity(t *testing.T) {
	err := &types.ConditionalCheckFailedException{}

	if got := mapDynamoConditionalUpdateError(err); got != nil {
		t.Fatalf("expected nil error for older last_activity, got %v", got)
	}
}
