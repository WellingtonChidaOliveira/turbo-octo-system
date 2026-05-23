package retry

import (
	"testing"
	"time"
)

func TestPolicy_DelayUsesExponentialBackoffWithMax(t *testing.T) {
	policy := Policy{
		InitialBackoff: 100 * time.Millisecond,
		MaxBackoff:     250 * time.Millisecond,
		MaxAttempts:    5,
	}

	tests := []struct {
		attempt int
		want    time.Duration
	}{
		{attempt: 1, want: 100 * time.Millisecond},
		{attempt: 2, want: 200 * time.Millisecond},
		{attempt: 3, want: 250 * time.Millisecond},
		{attempt: 4, want: 250 * time.Millisecond},
	}

	for _, tt := range tests {
		if got := policy.Delay(tt.attempt); got != tt.want {
			t.Fatalf("attempt %d: expected %s, got %s", tt.attempt, tt.want, got)
		}
	}
}

func TestPolicy_Normalize(t *testing.T) {
	policy := Policy{
		InitialBackoff: -1,
		MaxBackoff:     -1,
		Jitter:         -1,
		MaxAttempts:    0,
	}.Normalize()

	if policy.MaxAttempts != 1 {
		t.Fatalf("expected max attempts 1, got %d", policy.MaxAttempts)
	}
	if policy.InitialBackoff != 0 {
		t.Fatalf("expected initial backoff 0, got %s", policy.InitialBackoff)
	}
	if policy.MaxBackoff != 0 {
		t.Fatalf("expected max backoff 0, got %s", policy.MaxBackoff)
	}
	if policy.Jitter != 0 {
		t.Fatalf("expected jitter 0, got %s", policy.Jitter)
	}
}
