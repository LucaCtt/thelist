package data

import (
	"testing"
)

const dbPath string = "./thelisttest.db"

var opt = &DbOptions{Path: dbPath}

// i returns a pointer to the given integer.
// Useful for passing ShowID values to Shows.
func i(x int) *int {
	return &x
}

func clearDbStore(store *DbStore) {
	shows, _ := store.GetAllShows()

	for i := 0; i < len(shows); i++ {
		store.DeleteShow(shows[i].ID)
	}
}

func newCleanDbStore() *DbStore {
	store, _ := NewDbStore(opt)
	clearDbStore(store)

	return store
}

func TestNewDbStoreIntegration_ValidDatabase_ConnectsSuccessfully(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	store, err := NewDbStore(opt)
	defer store.Close()

	if err != nil {
		t.Errorf("Cannot connect to database: %s", err)
	}
}

func TestNewDbStoreIntegration_InvalidDatabase_ReturnsError(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	_, err := NewDbStore(&DbOptions{Path: "/bush/did/9/11"})

	if err == nil {
		t.Error("Connecting to invalid database did not return error")
	}
}

func TestNewDbStoreIntegration_InvalidDatabase_StoreNil(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	store, _ := NewDbStore(&DbOptions{Path: "/bush/did/9/11"})

	if store != nil {
		t.Error("Store should be nil if database connection is invalid")
	}
}

func TestCloseIntegration_ClosesSuccessfully(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	store := newCleanDbStore()
	err := store.Close()

	if err != nil {
		t.Errorf("Error while closing store: %s", err)
	}
}

func TestGetAllShowsIntegration_DatabaseEmpty_ReturnsNoShows(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	store := newCleanDbStore()
	defer store.Close()

	shows, err := store.GetAllShows()

	if err != nil {
		t.Errorf("Error while reading shows %s", err)
	}

	if len(shows) != 0 {
		t.Errorf("Shows number is incorrect: expected %d, received %d", 0, len(shows))
	}
}

func TestGetAllShowsIntegration_ShowsInDatabase_ReturnsAllShows(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	store := newCleanDbStore()
	defer store.Close()

	data := []*Show{&Show{ShowID: i(1)}, &Show{ShowID: i(2)}}
	for i := 0; i < len(data); i++ {
		store.CreateShow(data[i])
	}

	shows, err := store.GetAllShows()

	if err != nil {
		t.Errorf("Error while reading shows %s", err)
	}

	if len(shows) != len(data) {
		t.Errorf("Shows number is incorrect: expected %d, received %d", len(data), len(shows))
	}

	for i := 0; i < len(shows); i++ {
		if *shows[i].ShowID != *data[i].ShowID {
			t.Errorf("Show ShowId is incorrect: expected %d, received %d", data[i].ShowID, shows[i].ShowID)
		}
	}
}

func TestGetShowIntegration_ShowExists_ReturnsShow(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	store := newCleanDbStore()
	defer store.Close()

	data := &Show{ShowID: i(1)}
	store.CreateShow(data)

	show, err := store.GetShow(data.ID)

	if err != nil {
		t.Errorf("Error while reading show %s", err)
	}

	if *show.ShowID != *data.ShowID {
		t.Errorf("Show ID is incorrect: expected %d, received %d", data.ShowID, show.ShowID)
	}
}

func TestGetShowIntegration_ShowDoesNotExists_ReturnsNil(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	store := newCleanDbStore()
	defer store.Close()

	show, err := store.GetShow(1)

	if err == nil {
		t.Error("Reading show that does not exist did not return an error")
	}

	if show != nil {
		t.Error("Reading show that does not exist returned non-nil value")
	}
}

func TestCreateShowIntegration_CreatesShow(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	store := newCleanDbStore()
	defer store.Close()

	data := &Show{ShowID: i(1)}
	err := store.CreateShow(data)

	if err != nil {
		t.Errorf("Error while creating show: %s", err)
	}
}

func TestDeleteShowIntegration_ShowExists_DeletesShow(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	store := newCleanDbStore()
	defer store.Close()

	data := &Show{ShowID: i(1)}
	store.CreateShow(data)
	err := store.DeleteShow(data.ID)

	if err != nil {
		t.Errorf("Error while deleting show: %s", err)
	}
}
