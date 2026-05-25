package config

import (
	"os"
	"strconv"
	"time"
)

type Settings struct {
	ProcessorID            string
	AwsEndpointURL         string
	RawQueueURL            string
	ProcessedQueueURL      string
	EventsTableName        string
	Region                 string
	KeysAwsAccessKeyId     string
	KeysAwsSecretAccessKey string
	WorkerCount            int
	JobBufferSize          int
	ShutdownTimeout        time.Duration
	ReceiveBackoff         RetrySettings
	PublishBackoff         RetrySettings
	Queue                  QueueSettings
}

type RetrySettings struct {
	InitialBackoff time.Duration
	MaxBackoff     time.Duration
	Jitter         time.Duration
	MaxAttempts    int
}

type QueueSettings struct {
	MaxNumberOfMessages int32
	WaitTimeSeconds     int32
	VisibilityTimeout   int32
}

func LoadSettings() Settings {
	var settings Settings
	settings.ProcessorID = getDefaultString("PROCESSOR_ID", "processor-1")
	settings.AwsEndpointURL = getDefaultString("AWS_ENDPOINT_URL", getDefaultString("QUEUE_URL", "http://localhost:4566/"))
	settings.RawQueueURL = getDefaultString("RAW_QUEUE_URL", "http://localhost:4566/000000000000/raw-events")
	settings.ProcessedQueueURL = getDefaultString("PROCESSED_QUEUE_URL", "http://localhost:4566/000000000000/processed-events")
	settings.EventsTableName = getDefaultString("EVENTS_TABLE_NAME", "events")
	settings.Region = getDefaultString("REGION", "us-east-1")
	settings.KeysAwsAccessKeyId = getDefaultString("AWS_ACCESS_KEY_ID", "test")
	settings.KeysAwsSecretAccessKey = getDefaultString("AWS_SECRET_ACCESS_KEY", "test")
	settings.WorkerCount = normalizeWorkerCount(getDefaultInt("WORKER_COUNT", 5))
	settings.JobBufferSize = normalizePositiveInt(getDefaultInt("JOB_BUFFER_SIZE", settings.WorkerCount*2), settings.WorkerCount*2)
	settings.ShutdownTimeout = getDuration("SHUTDOWN_TIMEOUT", 30*time.Second)
	settings.ReceiveBackoff = RetrySettings{
		InitialBackoff: getDuration("RECEIVE_BACKOFF_INITIAL", 1*time.Second),
		MaxBackoff:     getDuration("RECEIVE_BACKOFF_MAX", 30*time.Second),
		Jitter:         getDuration("RECEIVE_BACKOFF_JITTER", 250*time.Millisecond),
		MaxAttempts:    normalizePositiveInt(getDefaultInt("RECEIVE_MAX_ATTEMPTS", 0), 0),
	}
	settings.PublishBackoff = RetrySettings{
		InitialBackoff: getDuration("PUBLISH_BACKOFF_INITIAL", 100*time.Millisecond),
		MaxBackoff:     getDuration("PUBLISH_BACKOFF_MAX", 2*time.Second),
		Jitter:         getDuration("PUBLISH_BACKOFF_JITTER", 50*time.Millisecond),
		MaxAttempts:    normalizePositiveInt(getDefaultInt("PUBLISH_MAX_ATTEMPTS", 3), 3),
	}
	settings.Queue = QueueSettings{
		MaxNumberOfMessages: int32(normalizePositiveInt(getDefaultInt("SQS_MAX_NUMBER_OF_MESSAGES", 10), 10)),
		WaitTimeSeconds:     int32(normalizeNonNegativeInt(getDefaultInt("SQS_WAIT_TIME_SECONDS", 20), 20)),
		VisibilityTimeout:   int32(normalizePositiveInt(getDefaultInt("SQS_VISIBILITY_TIMEOUT_SECONDS", 30), 30)),
	}

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
