package entities

import "time"

type RawEvent struct {
	EventID     string `json:"event_id"`
	DeveloperID string `json:"developer_id"`
	MetricType  string `json:"metric_type"`
	Value       int    `json:"value"`
	Repository  string `json:"repository"`
	Timestamp   string `json:"timestamp"`
}

func (e RawEvent) ToProcessed(processorID string, processedAt time.Time) ProcessedEvent {
	return ProcessedEvent{
		EventID:     e.EventID,
		DeveloperID: e.DeveloperID,
		MetricType:  e.MetricType,
		Value:       e.Value,
		Repository:  e.Repository,
		Timestamp:   e.Timestamp,
		ProcessedAt: processedAt.UTC().Format(time.RFC3339),
		ProcessorID: processorID,
	}
}
