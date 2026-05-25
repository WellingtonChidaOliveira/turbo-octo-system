package repository

import (
	"aggregator/internal/domain/entities"
	"aggregator/internal/infra/repository/records"
	"aggregator/internal/usecase/apperrors"
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type DynamoEventStore struct {
	client           *dynamodb.Client
	eventsTableName  string
	summaryTableName string
}

func NewDynamoEventStore(client *dynamodb.Client, eventsTableName string, summaryTableName string) *DynamoEventStore {
	return &DynamoEventStore{
		client:           client,
		eventsTableName:  eventsTableName,
		summaryTableName: summaryTableName,
	}
}

func (s *DynamoEventStore) SaveEventAndIncrementSummary(ctx context.Context, event entities.ProcessedEvent) error {
	input, err := buildSaveEventAndSummaryInput(s.eventsTableName, s.summaryTableName, event)
	if err != nil {
		return err
	}

	_, err = s.client.TransactWriteItems(ctx, input)
	return mapDynamoWriteError(err)
}

func (s *DynamoEventStore) FindByID(ctx context.Context, eventID string) (entities.ProcessedEvent, error) {
	out, err := s.client.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: aws.String(s.eventsTableName),
		Key: map[string]types.AttributeValue{
			"event_id": &types.AttributeValueMemberS{Value: eventID},
		},
	})
	if err != nil {
		return entities.ProcessedEvent{}, err
	}
	if len(out.Item) == 0 {
		return entities.ProcessedEvent{}, apperrors.ErrNotFound
	}

	var record records.ProcessedEvent
	if err := attributevalue.UnmarshalMap(out.Item, &record); err != nil {
		return entities.ProcessedEvent{}, err
	}
	return record.ToEntity(), nil
}

func (s *DynamoEventStore) FindByDeveloperID(ctx context.Context, developerID string) ([]entities.ProcessedEvent, error) {
	paginator := dynamodb.NewQueryPaginator(s.client, &dynamodb.QueryInput{
		TableName:              aws.String(s.eventsTableName),
		IndexName:              aws.String("developer_id-index"),
		KeyConditionExpression: aws.String("developer_id = :developer_id"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":developer_id": &types.AttributeValueMemberS{Value: developerID},
		},
	})

	events := make([]entities.ProcessedEvent, 0)
	for paginator.HasMorePages() {
		out, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		var eventRecords []records.ProcessedEvent
		if err := attributevalue.UnmarshalListOfMaps(out.Items, &eventRecords); err != nil {
			return nil, err
		}

		for _, record := range eventRecords {
			events = append(events, record.ToEntity())
		}
	}

	return events, nil
}

func (s *DynamoEventStore) FindSummaryByDeveloperID(ctx context.Context, developerID string) (entities.DeveloperSummary, error) {
	out, err := s.client.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: aws.String(s.summaryTableName),
		Key: map[string]types.AttributeValue{
			"developer_id": &types.AttributeValueMemberS{Value: developerID},
		},
	})
	if err != nil {
		return entities.DeveloperSummary{}, err
	}
	if len(out.Item) == 0 {
		return entities.DeveloperSummary{}, apperrors.ErrNotFound
	}

	var record records.DeveloperSummary
	if err := attributevalue.UnmarshalMap(out.Item, &record); err != nil {
		return entities.DeveloperSummary{}, err
	}
	return record.ToEntity(), nil
}
