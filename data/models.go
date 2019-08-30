package data

import (
	"time"
)

// Item represents an item of the show list.
type Item struct {
	ID        uint      `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	ShowID    int       `json:"show_id"`
	Watched   bool      `json:"watched"`
}

// Show represents a movie or a tv series.
type Show struct {
	ID          int
	Name        string
	ReleaseDate string
	Popularity  float32
	VoteAverage float32
}

// ShowList represents a list of shows.
type ShowList struct {
	Results      []Show
	TotalResults int
}
