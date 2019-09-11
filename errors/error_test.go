package errors

import (
	"fmt"
	"testing"
)

func TestError_Error(t *testing.T) {
	e := &Error{
		Message:  "test",
		Err:      fmt.Errorf("error"),
		Code:     CodeNotFound,
		Severity: SeverityErr,
	}
	want := fmt.Sprintf("[%d] %d - %q", e.Severity, e.Code, e.Message)

	if got := e.Error(); got != want {
		t.Errorf("got %q, want %q", got, want)
	}

}

func TestError_Unwrap(t *testing.T) {
	e := &Error{
		Message:  "test",
		Err:      fmt.Errorf("error"),
		Code:     CodeNotFound,
		Severity: SeverityErr,
	}
	want := e.Err

	if got := e.Unwrap(); got != want {
		t.Errorf("got %v, want %v", got, want)
	}

}
