package config

import "testing"

func TestLoadSettings_NormalizesWorkerAndBuffer(t *testing.T) {
	t.Setenv("WORKER_COUNT", "0")
	t.Setenv("JOB_BUFFER_SIZE", "0")

	settings := LoadSettings()

	if settings.Worker.Count != 1 {
		t.Fatalf("expected worker count 1, got %d", settings.Worker.Count)
	}
	if settings.Worker.JobBufferSize != 2 {
		t.Fatalf("expected job buffer size 2, got %d", settings.Worker.JobBufferSize)
	}
}

func TestLoadSettings_LoadsQueueSettings(t *testing.T) {
	t.Setenv("SQS_MAX_NUMBER_OF_MESSAGES", "5")
	t.Setenv("SQS_WAIT_TIME_SECONDS", "10")
	t.Setenv("SQS_VISIBILITY_TIMEOUT_SECONDS", "45")

	settings := LoadSettings()

	if settings.Queue.MaxNumberOfMessages != 5 {
		t.Fatalf("expected max number of messages 5, got %d", settings.Queue.MaxNumberOfMessages)
	}
	if settings.Queue.WaitTimeSeconds != 10 {
		t.Fatalf("expected wait time seconds 10, got %d", settings.Queue.WaitTimeSeconds)
	}
	if settings.Queue.VisibilityTimeout != 45 {
		t.Fatalf("expected visibility timeout 45, got %d", settings.Queue.VisibilityTimeout)
	}
}

func TestLoadSettings_LoadsAPIPort(t *testing.T) {
	t.Setenv("API_PORT", "9090")

	settings := LoadSettings()

	if settings.API.Port != "9090" {
		t.Fatalf("expected api port 9090, got %s", settings.API.Port)
	}
}

func TestLoadSettings_PrefersAWSNamesAndKeepsLegacyAliases(t *testing.T) {
	t.Setenv("AWS_ENDPOINT_URL", "http://aws-endpoint:4566")
	t.Setenv("QUEUE_URL", "http://legacy-queue:4566")
	t.Setenv("AWS_REGION", "sa-east-1")
	t.Setenv("REGION", "us-east-1")

	settings := LoadSettings()

	if settings.AWS.EndpointURL != "http://aws-endpoint:4566" {
		t.Fatalf("expected preferred endpoint, got %s", settings.AWS.EndpointURL)
	}
	if settings.AWS.Region != "sa-east-1" {
		t.Fatalf("expected preferred region, got %s", settings.AWS.Region)
	}
}
