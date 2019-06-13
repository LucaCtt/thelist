package data

import (
	"fmt"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres" // Postgresql driver
)

// Show represents a generic show, like a movie or a TV series.
type Show struct {
	ID        uint      `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	ShowID    *int      `json:"show_id" gorm:"not null"`
}

// IsValid returns true if all the show's fields have valid values.
func (show *Show) IsValid() bool {
	if show.ShowID == nil {
		return false
	}

	return true
}

// Store represents a generic data store, which can be a database, a file, and so on.
// The Close method should always be called to close the store.
type Store interface {
	Close() error
	GetAllShows() ([]*Show, error)
	GetShow(id uint) (*Show, error)
	CreateShow(show *Show) error
	UpdateShow(id uint, show *Show) error
	DeleteShow(id uint) error
	IsRecordNotFoundError(err error) bool
}

// DbStore wraps a database into a Store.
type DbStore struct {
	db *gorm.DB
}

// DbOptions contains the values required to connect to a database.
type DbOptions struct {
	Host     string
	Port     int
	User     string
	Name     string
	Password string
}

// NewDbStore opens a connection to the specified postgresql db, updates its schema
// and returns it wrapped into a Store.
func NewDbStore(opt *DbOptions) (*DbStore, error) {
	connStr := fmt.Sprintf("host=%s port=%d user=%s dbname=%s password=%s sslmode=disable",
		opt.Host,
		opt.Port,
		opt.User,
		opt.Name,
		opt.Password)

	db, err := gorm.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	db.AutoMigrate(&Show{})
	return &DbStore{db}, nil
}

// Close closes the store. Should always be called with defer.
func (dbStore *DbStore) Close() error {
	return dbStore.db.Close()
}

// IsRecordNotFoundError returns true if error contains a RecordNotFound error.
func (dbStore *DbStore) IsRecordNotFoundError(err error) bool {
	return gorm.IsRecordNotFoundError(err)
}

// GetAllShows returns a slice containing all the shows in the store.
// If there are no shows, the slice will have length 0.
func (dbStore *DbStore) GetAllShows() ([]*Show, error) {
	var shows []*Show
	err := dbStore.db.Find(&shows).Error

	if err != nil {
		return nil, err
	}

	return shows, nil
}

// GetShow returns the show with the given id, or error if there is no such show.
func (dbStore *DbStore) GetShow(id uint) (*Show, error) {
	var show Show
	err := dbStore.db.First(&show, id).Error

	if err != nil {
		return nil, err
	}

	return &show, nil
}

// CreateShow adds the given show to the store.
func (dbStore *DbStore) CreateShow(show *Show) error {
	err := dbStore.db.Create(show).Error
	return err
}

// UpdateShow updates the show with the given id to match the given show. Returns an error
// if the id doesn't match any show.
func (dbStore *DbStore) UpdateShow(id uint, show *Show) error {
	showToUpdate, err := dbStore.GetShow(id)

	if err != nil {
		return err
	}

	showToUpdate.ShowID = show.ShowID
	err = dbStore.db.Save(showToUpdate).Error

	return err
}

// DeleteShow deletes the show with the given id. If there is no such show,
// will return error.
func (dbStore *DbStore) DeleteShow(id uint) error {
	var show Show
	err := dbStore.db.First(&show, id).Delete(&show).Error
	return err
}
