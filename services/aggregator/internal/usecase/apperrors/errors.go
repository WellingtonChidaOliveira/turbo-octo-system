package apperrors

import "errors"

var (
	ErrEventAlreadyProcessed = errors.New("event already processed")
	ErrNotFound              = errors.New("record not found")
)
