package records

import "aggregator/internal/domain/entities"

type DeveloperSummary struct {
	DeveloperID            string `dynamodbav:"developer_id"`
	TotalCommits           int    `dynamodbav:"total_commits"`
	TotalPullRequests      int    `dynamodbav:"total_pull_requests"`
	TotalReviewTimeMinutes int    `dynamodbav:"total_review_time_minutes"`
	ReviewTimeEvents       int    `dynamodbav:"review_time_events"`
	EventsProcessed        int    `dynamodbav:"events_processed"`
	LastActivity           string `dynamodbav:"last_activity"`
}

func (r DeveloperSummary) ToEntity() entities.DeveloperSummary {
	return entities.DeveloperSummary{
		DeveloperID:            r.DeveloperID,
		TotalCommits:           r.TotalCommits,
		TotalPullRequests:      r.TotalPullRequests,
		TotalReviewTimeMinutes: r.TotalReviewTimeMinutes,
		ReviewTimeEvents:       r.ReviewTimeEvents,
		EventsProcessed:        r.EventsProcessed,
		LastActivity:           r.LastActivity,
	}
}
