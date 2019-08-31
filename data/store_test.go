//go:generate mockgen -destination=../mocks/mock_store.go -package=mocks github.com/lucactt/thelist/data Store

package data

import (
	"testing"

	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

func assertErr(t *testing.T, got error, want bool) {
	t.Helper()
	if (got != nil) != want {
		t.Fatalf("got %q, wantErr %t", got, want)
	}
}

func assertEquals(t *testing.T, got, want *Item) {
	t.Helper()
	if got == nil && want == nil {
		return
	}

	if got == nil ||
		want == nil ||
		got.ShowID != want.ShowID ||
		got.Watched != want.Watched {
		t.Errorf("got %+v, want %+v", got, want)
	}
}

func assertEqualsAll(t *testing.T, got, want []Item) {
	t.Helper()
	for i, item := range got {
		assertEquals(t, &item, &want[i])
	}
}

func makeItems(ids ...int) []Item {
	result := make([]Item, len(ids))
	for i, id := range ids {
		result[i] = Item{ShowID: id}
	}
	return result
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
			_, err := NewDbStore(tt.path)
			assertErr(t, err, tt.wantErr)
		})
	}
}

func TestDbStore_Close(t *testing.T) {
	t.Run("no error", func(t *testing.T) {
		s, err := NewDbStore(":memory:")
		err = s.Close()
		assertErr(t, err, false)
	})

	t.Run("already closed", func(t *testing.T) {
		s, err := NewDbStore(":memory:")
		err = s.Close()
		err = s.Close()
		assertErr(t, err, false)
	})
}

func TestDbStore_All(t *testing.T) {
	tests := []struct {
		name    string
		items   []Item
		want    []Item
		wantErr bool
	}{
		{"empty db", []Item{}, []Item{}, false},
		{"items in db", makeItems(1, 2), makeItems(1, 2), false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s, err := NewDbStore(":memory:")
			defer s.Close()

			for _, i := range tt.items {
				s.Create(&i)
			}

			got, err := s.All()
			assertErr(t, err, tt.wantErr)
			assertEqualsAll(t, got, tt.want)
		})
	}
}

func TestDbStore_Get(t *testing.T) {
	tests := []struct {
		name    string
		items   []Item
		id      uint
		want    *Item
		wantErr bool
	}{
		{"item doesn't exist", []Item{}, 1, nil, true},
		{"item exists", makeItems(1), 1, &Item{ID: 1, ShowID: 1, Watched: false}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s, err := NewDbStore(":memory:")
			defer s.Close()

			for _, i := range tt.items {
				s.Create(&i)
			}

			got, err := s.Get(tt.id)
			assertErr(t, err, tt.wantErr)
			assertEquals(t, got, tt.want)
		})
	}
}

func TestDbStore_Create(t *testing.T) {
	tests := []struct {
		name    string
		item    *Item
		wantErr bool
	}{
		{"valid item", &Item{ShowID: 1}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s, err := NewDbStore(":memory:")
			defer s.Close()

			err = s.Create(tt.item)
			got, err := s.Get(1)
			assertErr(t, err, tt.wantErr)
			assertEquals(t, got, tt.item)
		})
	}
}

func TestDbStore_SetWatched(t *testing.T) {
	tests := []struct {
		name    string
		items   []Item
		id      uint
		watched bool
		want    *Item
		wantErr bool
	}{
		{"item doesn't exist", []Item{}, 1, false, nil, true},
		{"item exists", makeItems(1), 1, true, &Item{ID: 1, ShowID: 1, Watched: true}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s, err := NewDbStore(":memory:")
			defer s.Close()

			for _, i := range tt.items {
				s.Create(&i)
			}

			err = s.SetWatched(tt.id, tt.watched)
			got, _ := s.Get(tt.id)

			assertErr(t, err, tt.wantErr)
			assertEquals(t, got, tt.want)
		})
	}
}

func TestDbStore_Delete(t *testing.T) {
	tests := []struct {
		name    string
		items   []Item
		id      uint
		wantErr bool
	}{
		{"valid item", makeItems(1), 1, false},
		{"item doesn't exist", []Item{}, 1, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s, err := NewDbStore(":memory:")
			defer s.Close()

			for _, i := range tt.items {
				s.Create(&i)
			}
			err = s.Delete(tt.id)

			got, _ := s.Get(tt.id)
			assertErr(t, err, tt.wantErr)
			assertEquals(t, got, nil)
		})
	}
}
