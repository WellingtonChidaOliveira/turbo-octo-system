package responses

import "aggregator/internal/domain/entities"

type ProcessedEvent struct {
	EventID     string `json:"event_id"`
	DeveloperID string `json:"developer_id"`
	MetricType  string `json:"metric_type"`
	Value       int    `json:"value"`
	Repository  string `json:"repository"`
	Timestamp   string `json:"timestamp"`
	ProcessedAt string `json:"processed_at"`
	ProcessorID string `json:"processor_id"`
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

func NewProcessedEvents(events []entities.ProcessedEvent) []ProcessedEvent {
	response := make([]ProcessedEvent, 0, len(events))
	for _, event := range events {
		response = append(response, NewProcessedEvent(event))
	}
	return response
}
