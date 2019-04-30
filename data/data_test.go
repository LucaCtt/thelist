package data

import (
	"testing"
)

const dbHost string = "localhost"
const dbPort int = 5432
const dbUser string = "thelist_tester"
const dbName string = "thelist_test"
const dbPassword string = "password"

var opt = &DbOptions{Host: dbHost, Port: dbPort, User: dbUser, Name: dbName, Password: dbPassword}

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

func TestNewDbStoreIntegration_ValidDatabase(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	store, err := NewDbStore(opt)
	defer store.Close()

	if err != nil {
		t.Errorf("Cannot connect to database: %s", err)
	}
}

func TestNewDbStoreIntegration_InvalidDatabase(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	store, err := NewDbStore(&DbOptions{})

	if err == nil {
		t.Error("Connecting to invalid database did not return error")
	}

	if store != nil {
		t.Error("Store should be nil if database connection is invalid")
	}
}

func TestCloseIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	store := newCleanDbStore()
	err := store.Close()

	if err != nil {
		t.Errorf("Error while closing store: %s", err)
	}
}

func TestGetAllShowsIntegration_DatabaseEmpty(t *testing.T) {
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

func TestGetAllShowsIntegration_ShowsInDatabase(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	store := newCleanDbStore()
	defer store.Close()

	data := []*Show{&Show{ShowID: 1}, &Show{ShowID: 2}}
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
		if shows[i].ShowID != data[i].ShowID {
			t.Errorf("Show ShowId is incorrect: expected %d, received %d", data[i].ShowID, shows[i].ShowID)
		}
	}
}

func TestGetShowIntegration_ShowExists(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	store := newCleanDbStore()
	defer store.Close()

	data := &Show{ShowID: 1}
	store.CreateShow(data)

	show, err := store.GetShow(data.ID)

	if err != nil {
		t.Errorf("Error while reading show %s", err)
	}

	if show.ShowID != data.ShowID {
		t.Errorf("Show ID is incorrect: expected %d, received %d", data.ShowID, show.ShowID)
	}
}

func TestGetShowIntegration_ShowDoesNotExists(t *testing.T) {
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

func TestCreateShowIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	store := newCleanDbStore()
	defer store.Close()

	data := &Show{ShowID: 1}
	err := store.CreateShow(data)

	if err != nil {
		t.Errorf("Error while creating show: %s", err)
	}
}

func TestDeleteShowIntegration_ShowExists(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	store := newCleanDbStore()
	defer store.Close()

	data := &Show{ShowID: 1}
	store.CreateShow(data)
	err := store.DeleteShow(data.ID)

	if err != nil {
		t.Errorf("Error while deleting show: %s", err)
	}
}
