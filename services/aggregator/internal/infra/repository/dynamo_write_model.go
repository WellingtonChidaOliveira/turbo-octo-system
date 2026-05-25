package repository

import (
	"aggregator/internal/domain/entities"
	"aggregator/internal/infra/repository/records"
	"strconv"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

func buildSaveEventAndSummaryInput(eventsTable string, summaryTable string, event entities.ProcessedEvent) (*dynamodb.TransactWriteItemsInput, error) {
	eventItem, err := attributevalue.MarshalMap(records.NewProcessedEvent(event))
	if err != nil {
		return nil, err
	}

	delta := entities.NewSummaryDelta(event)
	return &dynamodb.TransactWriteItemsInput{
		TransactItems: []types.TransactWriteItem{
			buildPutEventItem(eventsTable, eventItem),
			buildIncrementSummaryItem(summaryTable, delta),
		},
	}, nil
}

func buildPutEventItem(table string, item map[string]types.AttributeValue) types.TransactWriteItem {
	return types.TransactWriteItem{
		Put: &types.Put{
			TableName:           aws.String(table),
			Item:                item,
			ConditionExpression: aws.String("attribute_not_exists(event_id)"),
		},
	}
}

func buildIncrementSummaryItem(table string, delta entities.DeveloperSummary) types.TransactWriteItem {
	return types.TransactWriteItem{
		Update: &types.Update{
			TableName: aws.String(table),
			Key: map[string]types.AttributeValue{
				"developer_id": &types.AttributeValueMemberS{Value: delta.DeveloperID},
			},
			UpdateExpression: aws.String(
				"SET last_activity = :last_activity " +
					"ADD total_commits :total_commits, total_pull_requests :total_pull_requests, " +
					"total_review_time_minutes :total_review_time_minutes, review_time_events :review_time_events, " +
					"events_processed :events_processed",
			),
			ExpressionAttributeValues: map[string]types.AttributeValue{
				":last_activity":             &types.AttributeValueMemberS{Value: delta.LastActivity},
				":total_commits":             &types.AttributeValueMemberN{Value: intToString(delta.TotalCommits)},
				":total_pull_requests":       &types.AttributeValueMemberN{Value: intToString(delta.TotalPullRequests)},
				":total_review_time_minutes": &types.AttributeValueMemberN{Value: intToString(delta.TotalReviewTimeMinutes)},
				":review_time_events":        &types.AttributeValueMemberN{Value: intToString(delta.ReviewTimeEvents)},
				":events_processed":          &types.AttributeValueMemberN{Value: intToString(delta.EventsProcessed)},
			},
		},
	}
}

func intToString(value int) string {
	return strconv.Itoa(value)
}
