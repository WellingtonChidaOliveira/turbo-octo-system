package config

import (
	"os"
	"strconv"
	"time"
)

type Settings struct {
	AWS             AWSSettings
	Queue           QueueSettings
	Dynamo          DynamoSettings
	API             APISettings
	Worker          WorkerSettings
	ReceiveBackoff  RetrySettings
	ShutdownTimeout time.Duration
}

type AWSSettings struct {
	EndpointURL     string
	Region          string
	AccessKeyID     string
	SecretAccessKey string
}

type QueueSettings struct {
	ProcessedQueueURL   string
	MaxNumberOfMessages int32
	WaitTimeSeconds     int32
	VisibilityTimeout   int32
}

type DynamoSettings struct {
	EventsTableName           string
	DeveloperSummaryTableName string
}

type APISettings struct {
	Port string
}

type WorkerSettings struct {
	Count         int
	JobBufferSize int
}

type RetrySettings struct {
	InitialBackoff time.Duration
	MaxBackoff     time.Duration
	Jitter         time.Duration
	MaxAttempts    int
}

func LoadSettings() Settings {
	var settings Settings
	settings.AWS = AWSSettings{
		EndpointURL:     getDefaultString("AWS_ENDPOINT_URL", getDefaultString("QUEUE_URL", "http://localhost:4566/")),
		Region:          getDefaultString("AWS_REGION", getDefaultString("REGION", "us-east-1")),
		AccessKeyID:     getDefaultString("AWS_ACCESS_KEY_ID", "test"),
		SecretAccessKey: getDefaultString("AWS_SECRET_ACCESS_KEY", "test"),
	}
	settings.Queue = QueueSettings{
		ProcessedQueueURL:   getDefaultString("PROCESSED_QUEUE_URL", "http://localhost:4566/000000000000/processed-events"),
		MaxNumberOfMessages: int32(normalizePositiveInt(getDefaultInt("SQS_MAX_NUMBER_OF_MESSAGES", 10), 10)),
		WaitTimeSeconds:     int32(normalizeNonNegativeInt(getDefaultInt("SQS_WAIT_TIME_SECONDS", 20), 20)),
		VisibilityTimeout:   int32(normalizePositiveInt(getDefaultInt("SQS_VISIBILITY_TIMEOUT_SECONDS", 30), 30)),
	}
	settings.Dynamo = DynamoSettings{
		EventsTableName:           getDefaultString("EVENTS_TABLE_NAME", "events"),
		DeveloperSummaryTableName: getDefaultString("DEVELOPER_SUMMARY_TABLE_NAME", "developer_summary"),
	}
	settings.API = APISettings{
		Port: getDefaultString("API_PORT", "8080"),
	}
	workerCount := normalizeWorkerCount(getDefaultInt("WORKER_COUNT", 5))
	settings.Worker = WorkerSettings{
		Count:         workerCount,
		JobBufferSize: normalizePositiveInt(getDefaultInt("JOB_BUFFER_SIZE", workerCount*2), workerCount*2),
	}
	settings.ReceiveBackoff = RetrySettings{
		InitialBackoff: getDuration("RECEIVE_BACKOFF_INITIAL", 1*time.Second),
		MaxBackoff:     getDuration("RECEIVE_BACKOFF_MAX", 30*time.Second),
		Jitter:         getDuration("RECEIVE_BACKOFF_JITTER", 250*time.Millisecond),
		MaxAttempts:    normalizePositiveInt(getDefaultInt("RECEIVE_MAX_ATTEMPTS", 0), 0),
	}
	settings.ShutdownTimeout = getDuration("SHUTDOWN_TIMEOUT", 30*time.Second)

	return settings
}

func getDefaultString(key string, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		value = defaultValue
	}
	return value
}

func getDefaultInt(key string, defaultValue int) int {
	valueStr := os.Getenv(key)
	if valueStr == "" {
		return defaultValue
	}
	value, err := strconv.Atoi(valueStr)
	if err != nil {
		return defaultValue
	}
	return value
}

func normalizeWorkerCount(workerCount int) int {
	if workerCount <= 0 {
		return 1
	}
	return workerCount
}

func normalizePositiveInt(value int, fallback int) int {
	if value <= 0 {
		return fallback
	}
	return value
}

func normalizeNonNegativeInt(value int, fallback int) int {
	if value < 0 {
		return fallback
	}
	return value
}

func getDuration(key string, defaultValue time.Duration) time.Duration {
	valueStr := os.Getenv(key)
	if valueStr == "" {
		return defaultValue
	}
	value, err := time.ParseDuration(valueStr)
	if err != nil {
		return defaultValue
	}
	return value
}
