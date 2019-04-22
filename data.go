package main

import (
	"encoding/json"
	"errors"
	"io/ioutil"
)

type Store interface {
	GetAllItems() ([]*Item, error)
	GetItem(id int) (*Item, error)
}

type fileStore struct {
	filePath string
}

type Item struct {
	ID int `json:"id"`
}

func (fileStore *fileStore) GetAllItems() ([]*Item, error) {
	items, err := ioutil.ReadFile(fileStore.filePath)
	if err != nil {
		return nil, errors.New("Cannot read file")
	}

	itemsJSON := []*Item{}
	err = json.Unmarshal([]byte(items), &itemsJSON)
	if err != nil {
		return nil, errors.New("Cannot convert items to json")
	}

	return itemsJSON, nil
}

func (fileStore *fileStore) GetItem(id int) (*Item, error) {
	return nil, nil
}

func NewFileStore(filePath string) *fileStore {
	return &fileStore{filePath}
}
