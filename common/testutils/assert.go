package testutils

import (
	"reflect"
	"testing"
)

func AssertErr(t *testing.T, got error, wantErr bool) {
	t.Helper()

	if (got != nil) != wantErr {
		t.Errorf("got %v, wantErr %v", got, wantErr)
	}
}

func AssertEqual(t *testing.T, got, want interface{}) {
	t.Helper()

	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %+v, want %+v", got, want)
	}
}

func AssertLenEqual(t *testing.T, got, want int) {
	t.Helper()

	if got != want {
		t.Errorf("got has len %d, want %d", got, want)
		return
	}
}
