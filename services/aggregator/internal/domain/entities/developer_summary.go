package entities

type DeveloperSummary struct {
	DeveloperID            string
	TotalCommits           int
	TotalPullRequests      int
	TotalReviewTimeMinutes int
	ReviewTimeEvents       int
	EventsProcessed        int
	LastActivity           string
}

func (s DeveloperSummary) AvgReviewTimeMinutes() float64 {
	if s.ReviewTimeEvents == 0 {
		return 0
	}
	return float64(s.TotalReviewTimeMinutes) / float64(s.ReviewTimeEvents)
}

func NewSummaryDelta(event ProcessedEvent) DeveloperSummary {
	summary := DeveloperSummary{
		DeveloperID:     event.DeveloperID,
		EventsProcessed: 1,
		LastActivity:    event.Timestamp,
	}

	switch event.MetricType {
	case MetricCommits:
		summary.TotalCommits = event.Value
	case MetricPullRequests:
		summary.TotalPullRequests = event.Value
	case MetricReviewTimeMinutes:
		summary.TotalReviewTimeMinutes = event.Value
		summary.ReviewTimeEvents = 1
	}

	return summary
}
