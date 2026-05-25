package entities

type ProcessedEvent struct {
	EventID     string
	DeveloperID string
	MetricType  string
	Value       int
	Repository  string
	Timestamp   string
	ProcessedAt string
	ProcessorID string
}
