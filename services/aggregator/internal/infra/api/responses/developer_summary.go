package responses

import "aggregator/internal/domain/entities"

type DeveloperSummary struct {
	DeveloperID          string  `json:"developer_id"`
	TotalCommits         int     `json:"total_commits"`
	TotalPullRequests    int     `json:"total_pull_requests"`
	AvgReviewTimeMinutes float64 `json:"avg_review_time_minutes"`
	EventsProcessed      int     `json:"events_processed"`
	LastActivity         string  `json:"last_activity"`
}

func NewDeveloperSummary(summary entities.DeveloperSummary) DeveloperSummary {
	return DeveloperSummary{
		DeveloperID:          summary.DeveloperID,
		TotalCommits:         summary.TotalCommits,
		TotalPullRequests:    summary.TotalPullRequests,
		AvgReviewTimeMinutes: summary.AvgReviewTimeMinutes(),
		EventsProcessed:      summary.EventsProcessed,
		LastActivity:         summary.LastActivity,
	}
}
