//go:generate mockgen -destination=../mocks/mock_store.go -package=mocks github.com/lucactt/thelist/data Store

package data

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite" // SQLite driver
	"github.com/pkg/errors"
)

// Store represents a generic data store, which can be a database, a file, and so on.
// The Close method should always be called to close the store.
type Store interface {
	Close() error
	All() ([]Item, error)
	Get(id uint) (*Item, error)
	Create(item *Item) error
	SetWatched(id uint, watched bool) error
	Delete(id uint) error
}

// DbStore wraps a database into a Store.
type DbStore struct {
	db *gorm.DB
}

// DbOptions contains the values required to connect to a database.
type DbOptions struct {
	Path string
}

// NewDbStore opens a connection to the specified postgresql db, updates its schema
// and returns it wrapped into a Store.
func NewDbStore(opt *DbOptions) (*DbStore, error) {
	db, err := gorm.Open("sqlite3", opt.Path)
	if err != nil {
		return nil, errors.Wrap(err, "create db store failed")
	}

	db.AutoMigrate(&Item{})
	return &DbStore{db}, nil
}

// Close closes the store. Should be called with defer.
func (s *DbStore) Close() error {
	err := s.db.Close()

	if err != nil {
		return errors.Wrap(err, "close db store failed")
	}
	return nil
}

// All returns a slice containing all the items in the store.
// If there are no items, the slice will have length 0.
func (s *DbStore) All() ([]Item, error) {
	var items []Item
	err := s.db.Find(&items).Error
	if err != nil {
		return nil, errors.Wrap(err, "get all items failed")
	}

	return items, nil
}

// First returns the item with the given id, or error if there is no such item.
func (s *DbStore) First(id uint) (*Item, error) {
	var item Item
	err := s.db.First(&item, id).Error
	if err != nil {
		return nil, errors.Wrap(err, "find item failed")
	}

	return &item, nil
}

// Create adds the given show to the store.
func (s *DbStore) Create(item *Item) error {
	err := s.db.Create(item).Error
	if err != nil {
		return errors.Wrap(err, "create item failed")
	}

	return nil
}

// SetWatched sets the "Watched" field of the item with the given id to the value passed as argument.
func (s *DbStore) SetWatched(id uint, watched bool) error {
	item, err := s.First(id)
	if err != nil {
		return errors.Wrap(err, "find item to set watched failed")
	}

	item.Watched = watched
	err = s.db.Save(item).Error
	if err != nil {
		return errors.Wrap(err, "set watched failed")
	}

	return err
}

// Delete deletes the item with the given id. If there is no such item, it
// will return an error.
func (s *DbStore) Delete(id uint) error {
	item, err := s.First(id)
	if err != nil {
		return errors.Wrap(err, "find item to delete failed")
	}

	err = s.db.Delete(item).Error
	if err != nil {
		return errors.Wrap(err, "delete item failed")
	}

	return err
}
