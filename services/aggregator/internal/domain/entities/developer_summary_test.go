package entities

import "testing"

func TestNewSummaryDelta_Commits(t *testing.T) {
	event := ProcessedEvent{
		DeveloperID: "dev-123",
		MetricType:  MetricCommits,
		Value:       3,
		Timestamp:   "2026-04-15T10:30:00Z",
	}

	summary := NewSummaryDelta(event)

	if summary.DeveloperID != event.DeveloperID {
		t.Fatalf("expected developer id %s, got %s", event.DeveloperID, summary.DeveloperID)
	}
	if summary.TotalCommits != 3 {
		t.Fatalf("expected total commits 3, got %d", summary.TotalCommits)
	}
	if summary.EventsProcessed != 1 {
		t.Fatalf("expected events processed 1, got %d", summary.EventsProcessed)
	}
	if summary.LastActivity != event.Timestamp {
		t.Fatalf("expected last activity %s, got %s", event.Timestamp, summary.LastActivity)
	}
}

func TestNewSummaryDelta_ReviewTime(t *testing.T) {
	event := ProcessedEvent{
		DeveloperID: "dev-123",
		MetricType:  MetricReviewTimeMinutes,
		Value:       45,
		Timestamp:   "2026-04-15T10:30:00Z",
	}

	summary := NewSummaryDelta(event)

	if summary.TotalReviewTimeMinutes != 45 {
		t.Fatalf("expected review time 45, got %d", summary.TotalReviewTimeMinutes)
	}
	if summary.ReviewTimeEvents != 1 {
		t.Fatalf("expected review time events 1, got %d", summary.ReviewTimeEvents)
	}
}

func TestDeveloperSummary_AvgReviewTimeMinutes(t *testing.T) {
	summary := DeveloperSummary{
		TotalReviewTimeMinutes: 90,
		ReviewTimeEvents:       2,
	}

	if summary.AvgReviewTimeMinutes() != 45 {
		t.Fatalf("expected avg review time 45, got %.2f", summary.AvgReviewTimeMinutes())
	}
}

func TestDeveloperSummary_AvgReviewTimeMinutesWithoutReviews(t *testing.T) {
	summary := DeveloperSummary{}

	if summary.AvgReviewTimeMinutes() != 0 {
		t.Fatalf("expected avg review time 0, got %.2f", summary.AvgReviewTimeMinutes())
	}
}
