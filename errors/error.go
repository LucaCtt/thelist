package errors

import (
	"fmt"
	"net/http"
)

// ErrorSeverity is used to classify the severity of an error.
// The zero value is SeverityErr.
type ErrorSeverity int
type ErrorCode int

// HTTP status codes used to identify errors.
const (
	CodeUnexpected ErrorCode = http.StatusInternalServerError
	CodeNotFound   ErrorCode = http.StatusNotFound
	CodeBadValue   ErrorCode = http.StatusBadRequest
)

// Severity levels
const (
	SeverityErr ErrorSeverity = iota
	SeverityWarn
)

// Error is a wrapper for an error value with added context.
type Error struct {
	Message  string
	Err      error
	Code     ErrorCode
	Severity ErrorSeverity
}

func (e *Error) Error() string {
	return fmt.Sprintf("[%d] %d - %q", e.Severity, e.Code, e.Message)
}

// Unwrap returns the wrapped error.
func (e *Error) Unwrap() error {
	return e.Err
}
