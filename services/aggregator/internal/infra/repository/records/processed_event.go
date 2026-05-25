package records

import "aggregator/internal/domain/entities"

type ProcessedEvent struct {
	EventID     string `dynamodbav:"event_id"`
	DeveloperID string `dynamodbav:"developer_id"`
	MetricType  string `dynamodbav:"metric_type"`
	Value       int    `dynamodbav:"value"`
	Repository  string `dynamodbav:"repository"`
	Timestamp   string `dynamodbav:"timestamp"`
	ProcessedAt string `dynamodbav:"processed_at"`
	ProcessorID string `dynamodbav:"processor_id"`
}

func NewProcessedEvent(event entities.ProcessedEvent) ProcessedEvent {
	return ProcessedEvent{
		EventID:     event.EventID,
		DeveloperID: event.DeveloperID,
		MetricType:  event.MetricType,
		Value:       event.Value,
		Repository:  event.Repository,
		Timestamp:   event.Timestamp,
		ProcessedAt: event.ProcessedAt,
		ProcessorID: event.ProcessorID,
	}
}

func (r ProcessedEvent) ToEntity() entities.ProcessedEvent {
	return entities.ProcessedEvent{
		EventID:     r.EventID,
		DeveloperID: r.DeveloperID,
		MetricType:  r.MetricType,
		Value:       r.Value,
		Repository:  r.Repository,
		Timestamp:   r.Timestamp,
		ProcessedAt: r.ProcessedAt,
		ProcessorID: r.ProcessorID,
	}
}
