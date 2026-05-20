package domain

type ProcessedEvent struct {
	EventId     string `json:"event_id"`
	DeveloperId string `json:"developer_id"`
	MetricType  string `json:"metric_type"`
	Value       int    `json:"value"`
	Repository  string `json:"repository"`
	Timestamp   string `json:"timestamp"`
	ProcessedAt string `json:"processed_at"`
	ProcessorId string `json:"processor_id"`
}
