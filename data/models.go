package data

import "time"

// Item represents an item of the show list.
type Item struct {
	ID        uint      `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	ShowID    *int      `json:"show_id" gorm:"not null"`
	Watched   bool      `json:"watched"`
}

// IsValid returns true if all the item's fields have valid values.
func (i *Item) IsValid() bool {
	return i.ShowID != nil
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
	Results      []*Show
	TotalResults int
}
