package errors

import (
	"fmt"
	"net/http"
)

type ErrorCode int
type ErrorSeverity int

// Error codes
const (
	CodeUnexpected = http.StatusInternalServerError
	CodeNotFound   = http.StatusNotFound
)

// Severity levels
const (
	SeverityErr = iota
	SeverityWarn
)

type Error struct {
	Message  string
	Err      error
	Code     ErrorCode
	Severity ErrorSeverity
}

func (e *Error) Error() string {
	return fmt.Sprintf("%d - %q", e.Code, e.Message)
}

func (e *Error) Unwrap() error {
	return e.Err
}
