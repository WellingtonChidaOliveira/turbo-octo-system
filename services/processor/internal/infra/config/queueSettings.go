package config

import (
	"os"
	"strconv"
	"strings"
)

const (
	defaultEndpoint    = "http://localhost:4566"
	defaultProcessorID = "processor-local-1"
	defaultRegion      = "us-east-1"
	defaultWorkerCount = 4
)

type Settings struct {
	Endpoint          string
	ProcessedQueueURL string
	ProcessorID       string
	RawQueueURL       string
	Region            string
	WorkerCount       int
}

func LoadSettings() Settings {
	endpoint := envOrDefault("AWS_ENDPOINT_URL", defaultEndpoint)

	return Settings{
		Endpoint:          endpoint,
		ProcessedQueueURL: envOrDefault("PROCESSED_QUEUE_URL", endpoint+"/000000000000/processed-events"),
		ProcessorID:       envOrDefault("PROCESSOR_ID", defaultProcessorID),
		RawQueueURL:       envOrDefault("RAW_QUEUE_URL", endpoint+"/000000000000/raw-events"),
		Region:            envOrDefault("AWS_REGION", defaultRegion),
		WorkerCount:       envIntOrDefault("WORKER_COUNT", defaultWorkerCount),
	}
}

func envOrDefault(key, fallback string) string {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return fallback
	}
	return value
}

func envIntOrDefault(key string, fallback int) int {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return fallback
	}

	parsed, err := strconv.Atoi(value)
	if err != nil || parsed <= 0 {
		return fallback
	}
	return parsed
}
