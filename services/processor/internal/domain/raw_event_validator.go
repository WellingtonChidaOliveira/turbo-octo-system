package domain

import (
	"errors"
	"strings"
	"time"
)

var validMetricTypes = map[string]struct{}{
	"commits":             {},
	"pull_requests":       {},
	"review_time_minutes": {},
}

func (e RawEvent) Validate(now time.Time) error {
	if err := e.validateIdentity(); err != nil {
		return err
	}
	if err := e.validateMetric(); err != nil {
		return err
	}
	return e.validateTimestamp(now)
}

func (e RawEvent) validateIdentity() error {
	if !isUUID(e.EventID) {
		return errors.New("event_id is required and must be a valid UUID")
	}
	if strings.TrimSpace(e.DeveloperID) == "" {
		return errors.New("developer_id is required")
	}
	return nil
}

func (e RawEvent) validateMetric() error {
	if !isValidMetricType(e.MetricType) {
		return errors.New("metric_type is invalid")
	}
	if e.Value < 0 {
		return errors.New("value must be greater than or equal to zero")
	}
	if e.MetricType == "review_time_minutes" && e.Value > 1440 {
		return errors.New("review_time_minutes cannot be greater than 1440")
	}
	return nil
}

func (e RawEvent) validateTimestamp(now time.Time) error {
	eventTime, err := time.Parse(time.RFC3339, e.Timestamp)
	if err != nil {
		return errors.New("timestamp is required")
	}
	if eventTime.After(now.UTC()) {
		return errors.New("timestamp cannot be in the future")
	}
	return nil
}

func isValidMetricType(metricType string) bool {
	_, ok := validMetricTypes[metricType]
	return ok
}

func isUUID(value string) bool {
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

	return true
}

func isUUIDSeparator(index int, char rune) bool {
	return (index == 8 || index == 13 || index == 18 || index == 23) && char == '-'
}

func isHex(char rune) bool {
	return (char >= '0' && char <= '9') ||
		(char >= 'a' && char <= 'f') ||
		(char >= 'A' && char <= 'F')
}
