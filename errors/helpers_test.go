package errors

import (
	"fmt"
	"reflect"
	"testing"
)

func Test_E(t *testing.T) {
	t.Run("valid args", func(t *testing.T) {
		got := E("test", fmt.Errorf("test"), CodeNotFound, SeverityErr)
		want := &Error{
			Message:  "test",
			Err:      fmt.Errorf("test"),
			Code:     CodeNotFound,
			Severity: SeverityErr,
		}

		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %+v, want %+v", got, want)
		}
	})

	t.Run("invalid arg", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("did not panic")
			}
		}()
		E(69.420)
	})
}

func Test_Code(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want ErrorCode
	}{
		{"error type is not *Error", fmt.Errorf("test"), CodeUnexpected},
		{"error code is not 0", &Error{Code: CodeNotFound, Err: fmt.Errorf("test")}, CodeNotFound},
		{"wrapped error is *Error", &Error{Err: &Error{Code: CodeNotFound}}, CodeNotFound},
		{"wrapped error is not *Error", &Error{Err: fmt.Errorf("test")}, CodeUnexpected},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Code(tt.err)

			if got != tt.want {
				t.Errorf("got %d, want %d", got, tt.want)
			}
		})
	}
}
