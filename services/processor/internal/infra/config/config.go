package config

import (
	"os"
	"strconv"
)

type Settings struct {
	ProcessorID            string
	QueueUrl               string
	RawQueueURL            string
	ProcessedQueueURL      string
	Region                 string
	KeysAwsAccessKeyId     string
	KeysAwsSecretAccessKey string
	WorkerCount            int
}

func LoadSettings() Settings {
	var settings Settings
	settings.ProcessorID = getDefaultString("PROCESSOR_ID", "processor-1")
	settings.QueueUrl = getDefaultString("QUEUE_URL", "http://localhost:4566/")
	settings.RawQueueURL = getDefaultString("RAW_QUEUE_URL", "http://localhost:4566/000000000000/raw-events")
	settings.ProcessedQueueURL = getDefaultString("PROCESSED_QUEUE_URL", "http://localhost:4566/000000000000/processed-events")
	settings.Region = getDefaultString("REGION", "us-east-1")
	settings.KeysAwsAccessKeyId = getDefaultString("AWS_ACCESS_KEY_ID", "test")
	settings.KeysAwsSecretAccessKey = getDefaultString("AWS_SECRET_ACCESS_KEY", "test")
	settings.WorkerCount = normalizeWorkerCount(getDefaultInt("WORKER_COUNT", 5))

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
