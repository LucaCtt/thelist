package store

import (
	"testing"

	"github.com/LucaCtt/thelist/common/testutils"
)

func assertItemsEqual(t *testing.T, got, want *Item) {
	t.Helper()

	if got == nil && want == nil {
		return
	}

	if got == nil ||
		want == nil ||
		got.ShowID != want.ShowID {
		t.Errorf("got %+v, want %+v", got, want)
	}
}

func assertItemsListsEqual(t *testing.T, got, want []*Item) {
	t.Helper()

	testutils.AssertLenEqual(t, len(got), len(want))
	for i, item := range got {
		assertItemsEqual(t, item, want[i])
	}
}

func genItems(len int) []*Item {
	result := make([]*Item, len)
	for i := 0; i < len; i++ {
		result[i] = &Item{ShowID: i + 1, Type: MovieType}
	}
	return result
}

func seedStore(t *testing.T, s *DbStore, items []*Item) {
	t.Helper()

	for _, i := range items {
		s.Create(i)
	}
}

func TestNewDbStore(t *testing.T) {
	tests := []struct {
		name    string
		path    string
		wantErr bool
	}{
		{"valid path", ":memory:", false},
		{"invalid path", "/a/b/c/d", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s, err := New(tt.path)
			testutils.AssertErr(t, err, tt.wantErr)
			if s != nil {
				s.Close()
			}
		})
	}
}

func TestDbStore_Close(t *testing.T) {
	t.Run("no error", func(t *testing.T) {
		s, _ := New(":memory:")
		err := s.Close()
		testutils.AssertErr(t, err, false)
	})

	t.Run("already closed", func(t *testing.T) {
		s, _ := New(":memory:")
		s.Close()
		err := s.Close()
		testutils.AssertErr(t, err, false)
	})
}

func TestDbStore_All(t *testing.T) {
	tests := []struct {
		name    string
		items   []*Item
		want    []*Item
		wantErr bool
	}{
		{"empty db", []*Item{}, []*Item{}, false},
		{"items in db", genItems(2), genItems(2), false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s, _ := New(":memory:")
			defer s.Close()
			seedStore(t, s, tt.items)

			got, err := s.All()
			testutils.AssertErr(t, err, tt.wantErr)
			assertItemsListsEqual(t, got, tt.want)
		})
	}
}

func TestDbStore_Get(t *testing.T) {
	tests := []struct {
		name    string
		items   []*Item
		id      uint
		want    *Item
		wantErr bool
	}{
		{"item doesn't exist", []*Item{}, 1, nil, true},
		{"item exists", genItems(1), 1, genItems(1)[0], false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s, _ := New(":memory:")
			defer s.Close()
			seedStore(t, s, tt.items)

			got, err := s.Get(tt.id)
			testutils.AssertErr(t, err, tt.wantErr)
			assertItemsEqual(t, got, tt.want)
		})
	}
}

func TestDbStore_Create(t *testing.T) {
	tests := []struct {
		name    string
		item    *Item
		wantErr bool
	}{
		{"valid item", &Item{ShowID: 1, Type: MovieType}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s, _ := New(":memory:")
			defer s.Close()

			err := s.Create(tt.item)
			testutils.AssertErr(t, err, tt.wantErr)

			got, _ := s.Get(1)
			assertItemsEqual(t, got, tt.item)
		})
	}
}

func TestDbStore_Delete(t *testing.T) {
	tests := []struct {
		name    string
		items   []*Item
		id      uint
		wantErr bool
	}{
		{"valid item", genItems(1), 1, false},
		{"item doesn't exist", []*Item{}, 1, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s, _ := New(":memory:")
			defer s.Close()
			seedStore(t, s, tt.items)

			err := s.Delete(tt.id)
			testutils.AssertErr(t, err, tt.wantErr)

			got, _ := s.Get(tt.id)
			assertItemsEqual(t, got, nil)
		})
	}
}
