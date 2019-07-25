package data

import "testing"

var opt = &DbOptions{Path: ":memory:"}

// i returns a pointer to the given integer.
// Useful for passing ShowID values to Items.
func i(x int) *int {
	return &x
}

func TestNewDbStore_ValidDatabase_ConnectsSuccessfully(t *testing.T) {
	store, err := NewDbStore(opt)
	defer store.Close()

	if err != nil {
		t.Errorf("Cannot connect to database: %s", err)
	}
}

func TestNewDbStore_InvalidDatabase_ReturnsError(t *testing.T) {
	store, err := NewDbStore(&DbOptions{Path: "/bush/did/9/11"})

	if err == nil {
		t.Error("Connecting to invalid database did not return error")
		defer store.Close()
	}
}

func TestNewDbStore_InvalidDatabase_StoreNil(t *testing.T) {
	store, _ := NewDbStore(&DbOptions{Path: "/bush/did/9/11"})

	if store != nil {
		t.Error("Store should be nil if database connection is invalid")
		defer store.Close()
	}
}

func TestClose_ClosesSuccessfully(t *testing.T) {
	store, err := NewDbStore(opt)
	err = store.Close()

	if err != nil {
		t.Errorf("Error while closing store: %s", err)
	}
}

func TestGetAllItems_DatabaseEmpty_ReturnsNoItems(t *testing.T) {
	store, err := NewDbStore(opt)
	defer store.Close()

	items, err := store.GetAllItems()

	if err != nil {
		t.Errorf("Error while reading items %s", err)
	}

	if len(items) != 0 {
		t.Errorf("Items number is incorrect: expected %d, received %d", 0, len(items))
	}
}

func TestGetAllItems_ItemsInDatabase_ReturnsAllItems(t *testing.T) {
	store, err := NewDbStore(opt)
	defer store.Close()

	data := []*Item{&Item{ShowID: i(1)}, &Item{ShowID: i(2)}}
	for i := 0; i < len(data); i++ {
		store.CreateItem(data[i])
	}

	items, err := store.GetAllItems()
	if err != nil {
		t.Errorf("Error while reading items %s", err)
	}

	if len(items) != len(data) {
		t.Errorf("Items number is incorrect: expected %d, received %d", len(data), len(items))
	}

	for i := 0; i < len(items); i++ {
		if *items[i].ShowID != *data[i].ShowID {
			t.Errorf("Item ItemId is incorrect: expected %d, received %d", data[i].ShowID, items[i].ShowID)
		}
	}
}

func TestGetItem_ItemExists_ReturnsItem(t *testing.T) {
	store, err := NewDbStore(opt)
	defer store.Close()

	data := &Item{ShowID: i(1)}
	store.CreateItem(data)

	show, err := store.GetItem(data.ID)
	if err != nil {
		t.Errorf("Error while reading show %s", err)
	}

	if *show.ShowID != *data.ShowID {
		t.Errorf("Item ID is incorrect: expected %d, received %d", data.ShowID, show.ShowID)
	}
}

func TestGetItem_ItemDoesNotExists_ReturnsNil(t *testing.T) {
	store, err := NewDbStore(opt)
	defer store.Close()

	show, err := store.GetItem(1)
	if err == nil {
		t.Error("Reading show that does not exist did not return an error")
	}

	if show != nil {
		t.Error("Reading show that does not exist returned non-nil value")
	}
}

func TestCreateItem_CreatesItem(t *testing.T) {
	store, err := NewDbStore(opt)
	defer store.Close()

	data := &Item{ShowID: i(1)}
	err = store.CreateItem(data)
	if err != nil {
		t.Errorf("Error while creating show: %s", err)
	}
}

func TestDeleteItem_ItemExists_DeletesItem(t *testing.T) {
	store, err := NewDbStore(opt)
	defer store.Close()

	data := &Item{ShowID: i(1)}
	store.CreateItem(data)
	err = store.DeleteItem(data.ID)

	if err != nil {
		t.Errorf("Error while deleting show: %s", err)
	}
}
