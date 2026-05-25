package repository

import (
	"aggregator/internal/domain/entities"
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type DynamoEventStore struct {
	client    *dynamodb.Client
	tableName string
}

func NewDynamoEventStore(client *dynamodb.Client, tableName string) *DynamoEventStore {
	return &DynamoEventStore{
		client:    client,
		tableName: tableName,
	}
}

func (s *DynamoEventStore) Save(ctx context.Context, event entities.ProcessedEvent) error {
	item, err := attributevalue.MarshalMap(processedEventRecordFromEntity(event))
	if err != nil {
		return err
	}

	_, err = s.client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName:           aws.String(s.tableName),
		Item:                item,
		ConditionExpression: aws.String("attribute_not_exists(event_id)"),
	})
	return err
}

func (s *DynamoEventStore) FindByID(ctx context.Context, eventID string) (entities.ProcessedEvent, error) {
	out, err := s.client.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: aws.String(s.tableName),
		Key: map[string]types.AttributeValue{
			"event_id": &types.AttributeValueMemberS{Value: eventID},
		},
	})
	if err != nil {
		return entities.ProcessedEvent{}, err
	}
	if len(out.Item) == 0 {
		return entities.ProcessedEvent{}, nil
	}

	var record processedEventRecord
	if err := attributevalue.UnmarshalMap(out.Item, &record); err != nil {
		return entities.ProcessedEvent{}, err
	}
	return record.toEntity(), nil
}

func (s *DynamoEventStore) FindByDeveloperID(ctx context.Context, developerID string) ([]entities.ProcessedEvent, error) {
	out, err := s.client.Scan(ctx, &dynamodb.ScanInput{
		TableName:        aws.String(s.tableName),
		FilterExpression: aws.String("developer_id = :developer_id"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":developer_id": &types.AttributeValueMemberS{Value: developerID},
		},
	})
	if err != nil {
		return nil, err
	}

	var records []processedEventRecord
	if err := attributevalue.UnmarshalListOfMaps(out.Items, &records); err != nil {
		return nil, err
	}

	events := make([]entities.ProcessedEvent, 0, len(records))
	for _, record := range records {
		events = append(events, record.toEntity())
	}
	return events, nil
}

type processedEventRecord struct {
	EventID     string `dynamodbav:"event_id"`
	DeveloperID string `dynamodbav:"developer_id"`
	MetricType  string `dynamodbav:"metric_type"`
	Value       int    `dynamodbav:"value"`
	Repository  string `dynamodbav:"repository"`
	Timestamp   string `dynamodbav:"timestamp"`
	ProcessedAt string `dynamodbav:"processed_at"`
	ProcessorID string `dynamodbav:"processor_id"`
}

func processedEventRecordFromEntity(event entities.ProcessedEvent) processedEventRecord {
	return processedEventRecord{
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

func (r processedEventRecord) toEntity() entities.ProcessedEvent {
	return entities.ProcessedEvent{
		EventID:     r.EventID,
		DeveloperID: r.DeveloperID,
		MetricType:  r.MetricType,
		Value:       r.Value,
		Repository:  r.Repository,
		Timestamp:   r.Timestamp,
		ProcessedAt: r.ProcessedAt,
		ProcessorID: r.ProcessorID,
	}
}
