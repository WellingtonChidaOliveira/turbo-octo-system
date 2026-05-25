package dto

import "processor/internal/domain/entities"

type RawEvent struct {
	EventID     string `json:"event_id"`
	DeveloperID string `json:"developer_id"`
	MetricType  string `json:"metric_type"`
	Value       int    `json:"value"`
	Repository  string `json:"repository"`
	Timestamp   string `json:"timestamp"`
}

func (e RawEvent) ToEntity() entities.RawEvent {
	return entities.RawEvent{
		EventID:     e.EventID,
		DeveloperID: e.DeveloperID,
		MetricType:  e.MetricType,
		Value:       e.Value,
		Repository:  e.Repository,
		Timestamp:   e.Timestamp,
	}
}

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
