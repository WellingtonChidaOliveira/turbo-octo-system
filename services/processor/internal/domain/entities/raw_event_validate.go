package entities

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"
)

var uuidRegex = regexp.MustCompile(
	`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`,
)

var validMetricTypes = map[string]bool{
	"commits":             true,
	"pull_requests":       true,
	"review_time_minutes": true,
}

const maxReviewTimeMinutes = 1440

type ValidationErrors []string

func (ve ValidationErrors) Error() string {
	return strings.Join(ve, "; ")
}

func (e RawEvent) Validate() error {
	var errs ValidationErrors

	if e.EventID == "" {
		errs = append(errs, "event_id: obrigatório")
	} else if !uuidRegex.MatchString(strings.ToLower(e.EventID)) {
		errs = append(errs, "event_id: formato UUID inválido")
	}

	if strings.TrimSpace(e.DeveloperID) == "" {
		errs = append(errs, "developer_id: obrigatório")
	}

	if !validMetricTypes[e.MetricType] {
		errs = append(errs, fmt.Sprintf(
			"metric_type: valor inválido %q (esperado: commits, pull_requests, review_time_minutes)",
			e.MetricType,
		))
	}

	if e.Value < 0 {
		errs = append(errs, "value: deve ser >= 0")
	} else if e.MetricType == "review_time_minutes" && e.Value > maxReviewTimeMinutes {
		errs = append(errs, fmt.Sprintf(
			"value: para review_time_minutes o máximo é %d (24h)", maxReviewTimeMinutes,
		))
	}

	if e.Timestamp == "" {
		errs = append(errs, "timestamp: obrigatório")
	} else {
		t, err := time.Parse(time.RFC3339, e.Timestamp)
		if err != nil {
			errs = append(errs, "timestamp: formato inválido (esperado RFC3339, ex: 2006-01-02T15:04:05Z)")
		} else if t.After(time.Now()) {
			errs = append(errs, "timestamp: não pode ser uma data futura")
		}
	}

	if len(errs) == 0 {
		return nil
	}
	return errs
}

func IsValidationError(err error) bool {
	var ve ValidationErrors
	return errors.As(err, &ve)
}
