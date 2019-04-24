package main

import (
	"fmt"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

// Show represents a generic show, like a movie or a TV series.
type Show struct {
	gorm.Model
	ShowID int
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
	Dbname   string
	Password string
}

// NewDbStore opens a connection to the specified db, updates its schema
// and returns it wrapped into a Store.
func NewDbStore(opt *DbOptions) (*DbStore, error) {
	connString := fmt.Sprintf("host=%s port=%d user=%s dbname=%s password=%s",
		opt.Host,
		opt.Port,
		opt.User,
		opt.Dbname,
		opt.Password)

	db, err := gorm.Open("postgres", connString)
	if err != nil {
		return nil, err
	}

	db.AutoMigrate(&Show{})
	return &DbStore{db}, nil
}

func (dbStore *DbStore) Close() error {
	return dbStore.db.Close()
}

func (dbStore *DbStore) GetAllShows() ([]*Show, error) {
	return nil, nil
}

func (dbStore *DbStore) GetShow(id uint) (*Show, error) {
	return nil, nil
}

func (dbStore *DbStore) CreateShow(show *Show) error {
	return nil
}

func (dbStore *DbStore) UpdateShow(id uint, show *Show) error {
	return nil
}

func (dbStore *DbStore) DeleteShow(id uint) error {
	return nil
}
