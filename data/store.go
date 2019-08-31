//go:generate mockgen -destination=../mocks/mock_store.go -package=mocks github.com/lucactt/thelist/data Store

package data

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3" // SQLite driver
	"github.com/pkg/errors"
)

// Store represents a generic data store, which can be a database, a file, and so on.
// The Close method should always be called to close the store.
type Store interface {
	Close() error
	All() ([]*Item, error)
	Get(id uint) (*Item, error)
	Create(item *Item) error
	SetWatched(id uint, watched bool) error
	Delete(id uint) error
}

// DbStore wraps a database into a Store.
type DbStore struct {
	db *sql.DB
}

// NewDbStore creates a store that uses an SQLite db.
// If the tables do not already exist on the db, they will be created.
func NewDbStore(path string) (*DbStore, error) {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("open db at path %q failed", path))
	}

	err = db.Ping()
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("ping db at path %q failed", path))
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS items (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		show_id INTEGER NOT NULL,
		watched BOOLEAN NOT NULL CHECK (watched IN (0,1))
		)`)
	if err != nil {
		return nil, errors.Wrap(err, "query to create items table failed")
	}

	return &DbStore{db}, nil
}

// Close closes the store. Should be deferred.
// Closing multiple times a store has no effect.
func (s *DbStore) Close() error {
	err := s.db.Close()

	if err != nil {
		return errors.Wrap(err, "close db failed")
	}
	return nil
}

// All returns a slice containing all the items in the store.
// If there are no items, the slice will have length 0.
func (s *DbStore) All() ([]*Item, error) {
	rows, err := s.db.Query("SELECT * FROM items")
	if err != nil {
		return nil, errors.Wrap(err, "query to get all items failed")
	}
	defer rows.Close()

	var items []*Item
	for rows.Next() {
		var item Item
		err = rows.Scan(&item.ID, &item.ShowID, &item.Watched)
		if err != nil {
			return nil, errors.Wrap(err, "scanning item row failed")
		}
		items = append(items, &item)
	}

	if rows.Err() != nil {
		return nil, errors.Wrap(err, "item iteration error")
	}

	return items, nil
}

// Get returns the item found with the given id, or an error if there is no such item.
func (s *DbStore) Get(id uint) (*Item, error) {
	var item Item
	err := s.db.QueryRow("SELECT * FROM items WHERE id = ?", id).Scan(&item.ID, &item.ShowID, &item.Watched)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("query to get item with id %d failed", id))
	}

	return &item, nil
}

// Create adds the given item to the store.
func (s *DbStore) Create(item *Item) error {
	_, err := s.db.Exec(`INSERT INTO items (show_id, watched) VALUES (? ,?)`, item.ShowID, item.Watched)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("create item %+v failed", item))
	}

	return nil
}

// SetWatched sets the "watched" field of the item with the given id to the value passed as argument.
// If the item is not found, it will return an error.
func (s *DbStore) SetWatched(id uint, watched bool) error {
	value := 1
	if !watched {
		value = 0
	}

	r, err := s.db.Exec("UPDATE items SET watched = ? WHERE id = ?", value, id)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("set watched of item with id %d to %t failed", id, watched))
	}

	affected, err := r.RowsAffected()
	if err != nil {
		return errors.Wrap(err, "retrieve affected rows failed")
	}
	if affected == 0 {
		return fmt.Errorf("item with id %d not found", id)
	}

	return nil
}

// Delete deletes the item with the given id. If there is no such item, it
// will return an error.
func (s *DbStore) Delete(id uint) error {
	r, err := s.db.Exec("DELETE FROM items WHERE id = ?", id)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("delete item with id %d failed", id))
	}

	affected, err := r.RowsAffected()
	if err != nil {
		return errors.Wrap(err, "retrieving affected rows failed")
	}
	if affected == 0 {
		return fmt.Errorf("item with id %d not found", id)
	}

	return nil
}
