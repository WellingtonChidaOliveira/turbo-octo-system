package dto

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

func (e ProcessedEvent) ToEntity() entities.ProcessedEvent {
	return entities.ProcessedEvent{
		EventID:     e.EventID,
		DeveloperID: e.DeveloperID,
		MetricType:  e.MetricType,
		Value:       e.Value,
		Repository:  e.Repository,
		Timestamp:   e.Timestamp,
		ProcessedAt: e.ProcessedAt,
		ProcessorID: e.ProcessorID,
	}
}
