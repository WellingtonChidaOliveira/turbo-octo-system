package domain

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
