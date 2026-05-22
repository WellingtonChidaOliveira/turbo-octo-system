package entities

import (
	"testing"
	"time"
)

func TestRawEvent_ToProcessed(t *testing.T) {
	raw := RawEvent{
		EventID:     "123e4567-e89b-12d3-a456-426614174000",
		DeveloperID: "dev123",
		MetricType:  "commits",
		Value:       5,
		Repository:  "org/repo",
		Timestamp:   "2026-04-15T10:30:00Z",
	}
	processedAt := time.Date(2026, 4, 15, 10, 30, 5, 0, time.FixedZone("BRT", -3*60*60))

	processed := raw.ToProcessed("processor-1", processedAt)

	if processed.EventID != raw.EventID {
		t.Fatalf("expected event_id %s, got %s", raw.EventID, processed.EventID)
	}
	if processed.DeveloperID != raw.DeveloperID {
		t.Fatalf("expected developer_id %s, got %s", raw.DeveloperID, processed.DeveloperID)
	}
	if processed.MetricType != raw.MetricType {
		t.Fatalf("expected metric_type %s, got %s", raw.MetricType, processed.MetricType)
	}
	if processed.Value != raw.Value {
		t.Fatalf("expected value %d, got %d", raw.Value, processed.Value)
	}
	if processed.Repository != raw.Repository {
		t.Fatalf("expected repository %s, got %s", raw.Repository, processed.Repository)
	}
	if processed.Timestamp != raw.Timestamp {
		t.Fatalf("expected timestamp %s, got %s", raw.Timestamp, processed.Timestamp)
	}
	if processed.ProcessorID != "processor-1" {
		t.Fatalf("expected processor_id processor-1, got %s", processed.ProcessorID)
	}
	if processed.ProcessedAt != "2026-04-15T13:30:05Z" {
		t.Fatalf("expected processed_at in UTC, got %s", processed.ProcessedAt)
	}
}
