package common

import "time"

const (
	// MovieType is used to identify a show as a movie
	MovieType = "Movie"

	// TvShowType is used to identify a show as a tv show.
	TvShowType = "Tv Show"
)

// Item represents an item of the show list.
type Item struct {
	ID      uint
	Type    string
	ShowID  int
	Watched bool
}

// Show represents a movie or a tv series.
type Show struct {
	ID          int
	Type        string
	Name        string
	ReleaseDate time.Time
	Popularity  float32
	VoteAverage float32
}
