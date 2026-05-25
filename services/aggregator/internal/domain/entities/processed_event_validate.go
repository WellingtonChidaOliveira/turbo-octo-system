package entities

import (
	"errors"
	"strings"
	"time"
)

var validProcessedMetricTypes = map[string]struct{}{
	MetricCommits:           {},
	MetricPullRequests:      {},
	MetricReviewTimeMinutes: {},
}

func (e ProcessedEvent) Validate() error {
	if err := e.validateIdentity(); err != nil {
		return err
	}
	if err := e.validateMetric(); err != nil {
		return err
	}
	return e.validateTimestamps()
}

func (e ProcessedEvent) validateIdentity() error {
	if !isUUIDV4(e.EventID) {
		return errors.New("event_id is required and must be a valid UUID v4")
	}
	if strings.TrimSpace(e.DeveloperID) == "" {
		return errors.New("developer_id is required")
	}
	if strings.TrimSpace(e.Repository) == "" {
		return errors.New("repository is required")
	}
	if strings.TrimSpace(e.ProcessorID) == "" {
		return errors.New("processor_id is required")
	}
	return nil
}

func (e ProcessedEvent) validateMetric() error {
	if _, ok := validProcessedMetricTypes[e.MetricType]; !ok {
		return errors.New("metric_type is invalid")
	}
	if e.Value < 0 {
		return errors.New("value must be greater than or equal to zero")
	}
	if e.MetricType == MetricReviewTimeMinutes && e.Value > 1440 {
		return errors.New("review_time_minutes cannot be greater than 1440")
	}
	return nil
}

func (e ProcessedEvent) validateTimestamps() error {
	if _, err := time.Parse(time.RFC3339, e.Timestamp); err != nil {
		return errors.New("timestamp is required and must be RFC3339")
	}
	if _, err := time.Parse(time.RFC3339, e.ProcessedAt); err != nil {
		return errors.New("processed_at is required and must be RFC3339")
	}
	return nil
}

func isUUIDV4(value string) bool {
	if len(value) != 36 {
		return false
	}

	for index, char := range value {
		if isUUIDSeparator(index, char) {
			continue
		}
		if index == 8 || index == 13 || index == 18 || index == 23 {
			return false
		}
		if !isHex(char) {
			return false
		}
	}

	return value[14] == '4' && isRFC4122Variant(value[19])
}

func isUUIDSeparator(index int, char rune) bool {
	return (index == 8 || index == 13 || index == 18 || index == 23) && char == '-'
}

func isHex(char rune) bool {
	return (char >= '0' && char <= '9') ||
		(char >= 'a' && char <= 'f') ||
		(char >= 'A' && char <= 'F')
}

func isRFC4122Variant(char byte) bool {
	return char == '8' || char == '9' || char == 'a' || char == 'A' || char == 'b' || char == 'B'
}
