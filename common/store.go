package common

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3" // SQLite driver
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
		return nil, fmt.Errorf("open db at path %q failed: %w", path, err)
	}

	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("ping db at path %q failed: %w", path, err)
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS items (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		show_id INTEGER NOT NULL,
		watched BOOLEAN NOT NULL CHECK (watched IN (0,1))
		)`)
	if err != nil {
		return nil, fmt.Errorf("query to create items table failed: %w", err)
	}

	return &DbStore{db}, nil
}

// Close closes the store. Should be deferred.
// Closing multiple times a store has no effect.
func (s *DbStore) Close() error {
	err := s.db.Close()

	if err != nil {
		return fmt.Errorf("close db failed: %w", err)
	}
	return nil
}

// All returns a slice containing all the items in the store.
// If there are no items, the slice will have length 0.
func (s *DbStore) All() ([]*Item, error) {
	rows, err := s.db.Query("SELECT * FROM items")
	if err != nil {
		return nil, fmt.Errorf("query to get all items failed: %w", err)
	}
	defer rows.Close()

	var items []*Item
	for rows.Next() {
		var item Item
		err = rows.Scan(&item.ID, &item.ShowID, &item.Watched)
		if err != nil {
			return nil, fmt.Errorf("scanning item row failed: %w", err)
		}
		items = append(items, &item)
	}

	if rows.Err() != nil {
		return nil, fmt.Errorf("item iteration error: %w", err)
	}

	return items, nil
}

// Get returns the item found with the given id, or an error if there is no such item.
func (s *DbStore) Get(id uint) (*Item, error) {
	var item Item
	err := s.db.QueryRow("SELECT * FROM items WHERE id = ?", id).Scan(&item.ID, &item.ShowID, &item.Watched)
	if err != nil {
		return nil, fmt.Errorf("query to get item with id %d failed: %w", id, err)
	}

	return &item, nil
}

// Create adds the given item to the store.
func (s *DbStore) Create(item *Item) error {
	_, err := s.db.Exec(`INSERT INTO items (show_id, watched) VALUES (? ,?)`, item.ShowID, item.Watched)
	if err != nil {
		return fmt.Errorf("create item %+v failed: %w", item, err)
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
		return fmt.Errorf("set watched of item with id %d to %t failed: %w", id, watched, err)
	}

	affected, err := r.RowsAffected()
	if err != nil {
		return fmt.Errorf("retrieve affected rows failed: %w", err)
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
		return fmt.Errorf("delete item with id %d failed:%w", id, err)
	}

	affected, err := r.RowsAffected()
	if err != nil {
		return fmt.Errorf("retrieving affected rows failed: %w", err)
	}
	if affected == 0 {
		return fmt.Errorf("item with id %d not found", id)
	}

	return nil
}
