package data

import (
	"errors"
	"testing"
)

var opt = &DbOptions{Path: ":memory:"}

func TestIsRecordNotFoundError(t *testing.T) {
	type args struct {
		err error
	}
	tests := []struct {
		name string
		args *args
		want bool
	}{
		{"Empty error", &args{errors.New("")}, false},
		{"RecordNotFoundError", &args{ErrRecordNotFound}, true},
		{"Not RecordNotFoundError", &args{errors.New("abcd")}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsRecordNotFoundError(tt.args.err); got != tt.want {
				t.Errorf("DbStore.IsRecordNotFoundError() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewDbStore(t *testing.T) {
	type args struct {
		opt *DbOptions
	}
	tests := []struct {
		name    string
		args    *args
		wantDb  bool
		wantErr bool
	}{
		{"Valid db", &args{opt}, true, false},
		{"Invalid db", &args{&DbOptions{Path: "/bush/did/9/11"}}, false, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewDbStore(tt.args.opt)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewDbStore() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if (got != nil) != tt.wantDb {
				t.Errorf("NewDbStore() got = %v, wantDb %v", err, tt.wantDb)
				return
			}
		})
	}
}

func TestDbStore_Close(t *testing.T) {
	valid, _ := NewDbStore(opt)
	closed, _ := NewDbStore(opt)
	closed.Close()
	tests := []struct {
		name    string
		dbStore *DbStore
		wantErr bool
	}{
		{"Valid store", valid, false},
		{"Store already closed", closed, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.dbStore.Close(); (err != nil) != tt.wantErr {
				t.Errorf("DbStore.Close() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDbStore_GetAllItems(t *testing.T) {
	emptyStore, _ := NewDbStore(opt)
	notEmptyStore, _ := NewDbStore(opt)
	data := []*Item{&Item{ShowID: p(1)}, &Item{ShowID: p(2)}}
	for i := 0; i < len(data); i++ {
		notEmptyStore.CreateItem(data[i])
	}

	tests := []struct {
		name    string
		dbStore *DbStore
		want    []*Item
		wantErr bool
	}{
		{"Empty store", emptyStore, []*Item{}, false},
		{"Items in store", notEmptyStore, data, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.dbStore.GetAllItems()
			if (err != nil) != tt.wantErr {
				t.Errorf("DbStore.GetAllItems() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if len(got) != len(tt.want) {
				t.Errorf("DbStore.GetAllItems() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDbStore_GetItem(t *testing.T) {
	emptyStore, _ := NewDbStore(opt)
	notEmptyStore, _ := NewDbStore(opt)
	notEmptyStore.CreateItem(&Item{ShowID: p(1)})

	type args struct {
		id uint
	}
	tests := []struct {
		name     string
		dbStore  *DbStore
		args     *args
		wantData bool
		wantErr  bool
	}{
		{"Item exists", notEmptyStore, &args{1}, true, false},
		{"Item doesn't exist", emptyStore, &args{1}, false, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.dbStore.GetItem(tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("DbStore.GetItem() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if (got != nil) != tt.wantData {
				t.Errorf("DbStore.GetItem() item = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestDbStore_CreateItem(t *testing.T) {
	store, _ := NewDbStore(opt)
	type args struct {
		item *Item
	}
	tests := []struct {
		name    string
		dbStore *DbStore
		args    *args
		wantErr bool
	}{
		{"Valid item", store, &args{&Item{ShowID: p(1)}}, false},
		{"Invalid item", store, &args{&Item{}}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.dbStore.CreateItem(tt.args.item); (err != nil) != tt.wantErr {
				t.Errorf("DbStore.CreateItem() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDbStore_SetWatched(t *testing.T) {
	emptyStore, _ := NewDbStore(opt)
	notEmptyStore, _ := NewDbStore(opt)
	item := &Item{ShowID: p(1)}
	notEmptyStore.CreateItem(item)
	type args struct {
		id      uint
		watched bool
	}
	tests := []struct {
		name    string
		dbStore *DbStore
		args    *args
		wantErr bool
	}{
		{"Item exists", notEmptyStore, &args{item.ID, true}, false},
		{"Item doesn't exists", emptyStore, &args{item.ID, true}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.dbStore.SetWatched(tt.args.id, tt.args.watched); (err != nil) != tt.wantErr {
				t.Errorf("DbStore.SetWatched() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDbStore_DeleteItem(t *testing.T) {
	emptyStore, _ := NewDbStore(opt)
	notEmptyStore, _ := NewDbStore(opt)
	item := &Item{ShowID: p(1)}
	notEmptyStore.CreateItem(item)
	type args struct {
		id uint
	}
	tests := []struct {
		name    string
		dbStore *DbStore
		args    *args
		wantErr bool
	}{
		{"Item exists", notEmptyStore, &args{item.ID}, false},
		{"Item doesn't exist", emptyStore, &args{item.ID}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.dbStore.DeleteItem(tt.args.id); (err != nil) != tt.wantErr {
				t.Errorf("DbStore.DeleteItem() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
