package config

import "testing"

func TestLoadSettings_NormalizesWorkerAndBuffer(t *testing.T) {
	t.Setenv("WORKER_COUNT", "0")
	t.Setenv("JOB_BUFFER_SIZE", "0")

	settings := LoadSettings()

	if settings.WorkerCount != 1 {
		t.Fatalf("expected worker count 1, got %d", settings.WorkerCount)
	}
	if settings.JobBufferSize != 2 {
		t.Fatalf("expected job buffer size 2, got %d", settings.JobBufferSize)
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
