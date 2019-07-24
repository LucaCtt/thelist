//go:generate mockgen -destination=../mocks/mock_store.go -package=mocks github.com/LucaCtt/thelist/data Store

package data

import (
	"fmt"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite" // SQLite driver
)

// Store represents a generic data store, which can be a database, a file, and so on.
// The Close method should always be called to close the store.
type Store interface {
	Close() error
	GetAllItems() ([]*Item, error)
	GetItem(id uint) (*Item, error)
	CreateItem(item *Item) error
	SetWatched(id uint, watched bool) error
	DeleteItem(id uint) error
	IsRecordNotFoundError(err error) bool
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
	connStr := fmt.Sprintf("file:%s",
		opt.Path)

	db, err := gorm.Open("sqlite3", connStr)
	if err != nil {
		return nil, err
	}

	db.AutoMigrate(&Item{})
	return &DbStore{db}, nil
}

// Close closes the store. Should be called with defer.
func (dbStore *DbStore) Close() error {
	return dbStore.db.Close()
}

// IsRecordNotFoundError returns true if err contains a RecordNotFound error.
func (dbStore *DbStore) IsRecordNotFoundError(err error) bool {
	return gorm.IsRecordNotFoundError(err)
}

// GetAllItems returns a slice containing all the items in the store.
// If there are no items, the slice will have length 0.
func (dbStore *DbStore) GetAllItems() ([]*Item, error) {
	var items []*Item
	err := dbStore.db.Find(&items).Error
	if err != nil {
		return nil, err
	}

	return items, nil
}

// GetItem returns the item with the given id, or error if there is no such item.
func (dbStore *DbStore) GetItem(id uint) (*Item, error) {
	var item Item
	err := dbStore.db.First(&item, id).Error
	if err != nil {
		return nil, err
	}

	return &item, nil
}

// CreateItem adds the given show to the store.
func (dbStore *DbStore) CreateItem(item *Item) error {
	err := dbStore.db.Create(item).Error
	return err
}

// SetWatched sets the "Watched" field of the item with the given id to the value passed as argument.
func (dbStore *DbStore) SetWatched(id uint, watched bool) error {
	item, err := dbStore.GetItem(id)
	if err != nil {
		return err
	}

	item.Watched = watched
	err = dbStore.db.Save(item).Error

	return err
}

// DeleteItem deletes the item with the given id. If there is no such item, it
// will return an error.
func (dbStore *DbStore) DeleteItem(id uint) error {
	var item Item
	err := dbStore.db.First(&item, id).Delete(&item).Error
	return err
}
