//go:generate mockgen -destination=../mocks/mock_store.go -package=mocks github.com/lucactt/thelist/data Store

package data

import (
	"reflect"
	"testing"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

func assertErr(t *testing.T, got error, want bool) {
	if (got != nil) != want {
		t.Errorf("error = %v, wantErr %v", got, want)
	}
}

func assertEquals(t *testing.T, got, want *Item) {
	t.Helper()
	if got == nil && want == nil {
		return
	}

	if got.ShowID != want.ShowID || got.Watched != want.Watched {
		t.Errorf("got %v, want %v", got, want)
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
	type args struct {
		opt *DbOptions
	}
	tests := []struct {
		name    string
		args    args
		want    *DbStore
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewDbStore(tt.args.opt)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewDbStore() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewDbStore() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDbStore_Close(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		{"no error", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s, err := NewDbStore(&DbOptions{":memory:"})
			err = s.Close()
			assertErr(t, err, tt.wantErr)
		})
	}
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
			s, err := NewDbStore(&DbOptions{":memory:"})
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

func TestDbStore_First(t *testing.T) {
	tests := []struct {
		name    string
		items   []Item
		id      uint
		want    *Item
		wantErr bool
	}{
		{"item doesn't exist", []Item{}, 1, nil, true},
		{"item exists", makeItems(1), 1, &Item{ShowID: 1, Watched: false}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s, err := NewDbStore(&DbOptions{":memory:"})
			defer s.Close()

			for _, i := range tt.items {
				s.Create(&i)
			}

			got, err := s.First(tt.id)
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
			s, err := NewDbStore(&DbOptions{":memory:"})
			defer s.Close()

			err = s.Create(tt.item)
			assertErr(t, err, tt.wantErr)
		})
	}
}

func TestDbStore_SetWatched(t *testing.T) {
	type fields struct {
		db *gorm.DB
	}
	type args struct {
		id      uint
		watched bool
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &DbStore{
				db: tt.fields.db,
			}
			if err := s.SetWatched(tt.args.id, tt.args.watched); (err != nil) != tt.wantErr {
				t.Errorf("DbStore.SetWatched() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDbStore_Delete(t *testing.T) {
	type fields struct {
		db *gorm.DB
	}
	type args struct {
		id uint
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &DbStore{
				db: tt.fields.db,
			}
			if err := s.Delete(tt.args.id); (err != nil) != tt.wantErr {
				t.Errorf("DbStore.Delete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
