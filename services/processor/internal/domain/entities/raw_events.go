package entities

import "time"

type RawEvent struct {
	EventID     string
	DeveloperID string
	MetricType  string
	Value       int
	Repository  string
	Timestamp   string
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
