package api

import (
	"aggregator/internal/domain/entities"
	"aggregator/internal/infra/api/responses"
	"aggregator/internal/usecase/apperrors"
	"context"
	"encoding/json"
	"errors"
	"net/http/httptest"
	"testing"
)

func TestHealth_ReturnsOK(t *testing.T) {
	server := newTestServer(&fakeHealth{}, &fakeHealth{}, &fakeEventsGetter{}, &fakeSummaryGetter{})

	resp, err := server.Test(httptest.NewRequest("GET", "/health", nil))
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if resp.StatusCode != 200 {
		t.Fatalf("expected status 200, got %d", resp.StatusCode)
	}
}

func TestHealth_ReturnsUnavailableWhenQueueFails(t *testing.T) {
	server := newTestServer(
		&fakeHealth{queueErr: errors.New("queue unavailable")},
		&fakeHealth{},
		&fakeEventsGetter{},
		&fakeSummaryGetter{},
	)

	resp, err := server.Test(httptest.NewRequest("GET", "/health", nil))
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if resp.StatusCode != 503 {
		t.Fatalf("expected status 503, got %d", resp.StatusCode)
	}
}

func TestHealth_ReturnsUnavailableWhenStorageFails(t *testing.T) {
	server := newTestServer(
		&fakeHealth{},
		&fakeHealth{storageErr: errors.New("storage unavailable")},
		&fakeEventsGetter{},
		&fakeSummaryGetter{},
	)

	resp, err := server.Test(httptest.NewRequest("GET", "/health", nil))
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if resp.StatusCode != 503 {
		t.Fatalf("expected status 503, got %d", resp.StatusCode)
	}
}

func TestGetMetrics_ReturnsEvents(t *testing.T) {
	server := newTestServer(
		&fakeHealth{},
		&fakeHealth{},
		&fakeEventsGetter{events: []entities.ProcessedEvent{{
			EventID:     "event-1",
			DeveloperID: "dev-001",
			MetricType:  entities.MetricCommits,
			Value:       12,
			Repository:  "org/api",
			Timestamp:   "2026-04-15T10:30:00Z",
			ProcessedAt: "2026-04-15T10:30:05Z",
			ProcessorID: "processor-1",
		}}},
		&fakeSummaryGetter{},
	)

	resp, err := server.Test(httptest.NewRequest("GET", "/metrics/dev-001", nil))
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if resp.StatusCode != 200 {
		t.Fatalf("expected status 200, got %d", resp.StatusCode)
	}

	var body []responses.ProcessedEvent
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		t.Fatalf("expected valid response body, got %v", err)
	}
	if len(body) != 1 || body[0].DeveloperID != "dev-001" {
		t.Fatalf("expected dev-001 event, got %+v", body)
	}
}

func TestGetSummary_ReturnsSummary(t *testing.T) {
	server := newTestServer(
		&fakeHealth{},
		&fakeHealth{},
		&fakeEventsGetter{},
		&fakeSummaryGetter{summary: entities.DeveloperSummary{
			DeveloperID:            "dev-001",
			TotalCommits:           12,
			TotalPullRequests:      3,
			TotalReviewTimeMinutes: 90,
			ReviewTimeEvents:       2,
			EventsProcessed:        5,
			LastActivity:           "2026-04-15T10:30:00Z",
		}},
	)

	resp, err := server.Test(httptest.NewRequest("GET", "/metrics/dev-001/summary", nil))
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if resp.StatusCode != 200 {
		t.Fatalf("expected status 200, got %d", resp.StatusCode)
	}

	var body responses.DeveloperSummary
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		t.Fatalf("expected valid response body, got %v", err)
	}
	if body.DeveloperID != "dev-001" || body.AvgReviewTimeMinutes != 45 {
		t.Fatalf("expected dev-001 summary with avg 45, got %+v", body)
	}
}

func TestGetSummary_ReturnsNotFound(t *testing.T) {
	server := newTestServer(
		&fakeHealth{},
		&fakeHealth{},
		&fakeEventsGetter{},
		&fakeSummaryGetter{err: apperrors.ErrNotFound},
	)

	resp, err := server.Test(httptest.NewRequest("GET", "/metrics/dev-001/summary", nil))
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if resp.StatusCode != 404 {
		t.Fatalf("expected status 404, got %d", resp.StatusCode)
	}
}

func newTestServer(
	queueHealth *fakeHealth,
	storageHealth *fakeHealth,
	eventsGetter *fakeEventsGetter,
	summaryGetter *fakeSummaryGetter,
) *Server {
	return NewServer(
		":0",
		NewHealthHandler(queueHealth, storageHealth),
		NewMetricsHandler(eventsGetter),
		NewSummaryHandler(summaryGetter),
	)
}

type fakeHealth struct {
	queueErr   error
	storageErr error
}

func (f *fakeHealth) CheckQueue(ctx context.Context) error {
	return f.queueErr
}

func (f *fakeHealth) CheckStorage(ctx context.Context) error {
	return f.storageErr
}

type fakeEventsGetter struct {
	events []entities.ProcessedEvent
	err    error
}

func (f *fakeEventsGetter) Handle(ctx context.Context, developerID string) ([]entities.ProcessedEvent, error) {
	return f.events, f.err
}

type fakeSummaryGetter struct {
	summary entities.DeveloperSummary
	err     error
}

func (f *fakeSummaryGetter) Handle(ctx context.Context, developerID string) (entities.DeveloperSummary, error) {
	return f.summary, f.err
}
