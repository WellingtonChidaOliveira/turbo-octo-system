package usecase

import "fmt"

type ErrorType string

const (
	ErrorTypeDecode     ErrorType = "decode"
	ErrorTypeValidation ErrorType = "validation"
	ErrorTypeEncode     ErrorType = "encode"
	ErrorTypePublish    ErrorType = "publish"
	ErrorTypeDelete     ErrorType = "delete"
)

type ProcessingError struct {
	Type    ErrorType
	EventID string
	Err     error
}

func (e ProcessingError) Error() string {
	if e.EventID == "" {
		return fmt.Sprintf("%s error: %v", e.Type, e.Err)
	}
	return fmt.Sprintf("%s error for event %s: %v", e.Type, e.EventID, e.Err)
}

func (e ProcessingError) Unwrap() error {
	return e.Err
}
