package entities

import (
	"testing"
	"time"
)

func mockRawEvent() RawEvent {
	return RawEvent{
		EventID:     "123e4567-e89b-12d3-a456-426614174000",
		DeveloperID: "dev123",
		MetricType:  "commits",
		Value:       5,
		Timestamp:   "2024-01-01T12:00:00Z",
	}
}

func TestRawEvent_Validate(t *testing.T) {
	now := time.Date(2026, 5, 22, 12, 0, 0, 0, time.UTC)

	tests := []struct {
		name    string
		event   RawEvent
		wantErr bool
	}{
		{
			name:    "valid event",
			event:   mockRawEvent(),
			wantErr: false,
		},
		{
			name: "missing event_id",
			event: RawEvent{
				EventID:     "",
				DeveloperID: "dev123",
				MetricType:  "commits",
				Value:       5,
				Timestamp:   "2024-01-01T12:00:00Z",
			},
			wantErr: true,
		},
		{
			name: "invalid event_id format",
			event: RawEvent{
				EventID:     "invalid-uuid",
				DeveloperID: "dev123",
				MetricType:  "commits",
				Value:       5,
				Timestamp:   "2024-01-01T12:00:00Z",
			},
			wantErr: true,
		},
		{
			name: "missing developer_id",
			event: RawEvent{
				EventID:     "123e4567-e89b-12d3-a456-426614174000",
				DeveloperID: "",
				MetricType:  "commits",
				Value:       5,
				Timestamp:   "2024-01-01T12:00:00Z",
			},
			wantErr: true,
		},
		{
			name: "blank developer_id",
			event: RawEvent{
				EventID:     "123e4567-e89b-12d3-a456-426614174000",
				DeveloperID: "   ",
				MetricType:  "commits",
				Value:       5,
				Timestamp:   "2024-01-01T12:00:00Z",
			},
			wantErr: true,
		},
		{
			name: "invalid metric_type",
			event: RawEvent{
				EventID:     "123e4567-e89b-12d3-a456-426614174000",
				DeveloperID: "dev123",
				MetricType:  "invalid_metric",
				Value:       5,
				Timestamp:   "2024-01-01T12:00:00Z",
			},
			wantErr: true,
		},
		{
			name: "negative value",
			event: RawEvent{
				EventID:     "123e4567-e89b-12d3-a456-426614174000",
				DeveloperID: "dev123",
				MetricType:  "commits",
				Value:       -1,
				Timestamp:   "2024-01-01T12:00:00Z",
			},
			wantErr: true,
		},
		{
			name: "review_time_minutes at max",
			event: RawEvent{
				EventID:     "123e4567-e89b-12d3-a456-426614174000",
				DeveloperID: "dev123",
				MetricType:  "review_time_minutes",
				Value:       1440,
				Timestamp:   "2024-01-01T12:00:00Z",
			},
			wantErr: false,
		},
		{
			name: "review_time_minutes exceeding max",
			event: RawEvent{
				EventID:     "123e4567-e89b-12d3-a456-426614174000",
				DeveloperID: "dev123",
				MetricType:  "review_time_minutes",
				Value:       1500,
				Timestamp:   "2024-01-01T12:00:00Z",
			},
			wantErr: true,
		},
		{
			name: "missing timestamp",
			event: RawEvent{
				EventID:     "123e4567-e89b-12d3-a456-426614174000",
				DeveloperID: "dev123",
				MetricType:  "commits",
				Value:       5,
				Timestamp:   "",
			},
			wantErr: true,
		},
		{name: "invalid timestamp format",
			event: RawEvent{
				EventID:     "123e4567-e89b-12d3-a456-426614174000",
				DeveloperID: "dev123",
				MetricType:  "commits",
				Value:       5,
				Timestamp:   "invalid-timestamp",
			},
			wantErr: true,
		},
		{
			name: "future timestamp",
			event: RawEvent{
				EventID:     "123e4567-e89b-12d3-a456-426614174000",
				DeveloperID: "dev123",
				MetricType:  "commits",
				Value:       5,
				Timestamp:   "2026-05-22T12:00:01Z",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.event.Validate(now)
			if (err != nil) != tt.wantErr {
				t.Errorf("RawEvent.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
