package common

import (
	"fmt"
	"time"

	"github.com/LucaCtt/thelist/common/client"
	"github.com/LucaCtt/thelist/common/store"
)

// Show is a common representation for movies and tv series.
type Show struct {
	ID          int
	Type        string
	Name        string
	Year        int
	VoteAverage float32
}

func parseYear(date string) int {
	parsed, err := time.Parse(client.DateFormat, date)
	if err != nil {
		return 0
	}
	return parsed.Year()
}

func movieToShow(movie *client.Movie) *Show {
	return &Show{
		ID:          movie.ID,
		Type:        store.MovieType,
		Name:        movie.Title,
		Year:        parseYear(movie.ReleaseDate),
		VoteAverage: movie.VoteAverage,
	}
}

func tvToShow(tv *client.TvShow) *Show {
	return &Show{
		ID:          tv.ID,
		Type:        store.TvShowType,
		Name:        tv.Name,
		Year:        parseYear(tv.FirstAirDate),
		VoteAverage: tv.VoteAverage,
	}
}

func convertToShowsList(movies []*client.Movie, tv []*client.TvShow) []*Show {
	shows := make([]*Show, len(movies)+len(tv))

	for i, movie := range movies {
		shows[i] = movieToShow(movie)
	}

	for i, tv := range tv {
		shows[i+len(movies)] = tvToShow(tv)
	}

	return shows
}

// GetShow retrieves info on the movie or tv show (specified by the passed type) with the given id
// and returns it in the form of a Show.
func GetShow(c client.Client, id int, t string) (*Show, error) {
	var show *Show

	switch t {
	case store.MovieType:
		movie, err := c.GetMovie(id)
		if err != nil {
			return nil, fmt.Errorf("get show with id %d and type %q failed: %w", id, t, err)
		}
		show = movieToShow(movie)
	case store.TvShowType:
		tv, err := c.GetTvShow(id)
		if err != nil {
			return nil, fmt.Errorf("get show with id %d and type %q failed: %w", id, t, err)
		}
		show = tvToShow(tv)
	default:
		panic("invalid show type")
	}

	return show, nil
}

// SearchShow searches for movies and tv shows using the given client and
// returns the results in the form a list of shows.
func SearchShow(c client.Client, name string) ([]*Show, error) {
	moviesChan := make(chan []*client.Movie)
	tvChan := make(chan []*client.TvShow)
	errChan := make(chan error)

	go func() {
		m, err := c.SearchMovie(name)
		errChan <- err
		moviesChan <- m
	}()

	go func() {
		t, err := c.SearchTvShow(name)
		errChan <- err
		tvChan <- t
	}()

	for i := 0; i < 2; i++ {
		err := <-errChan
		if err != nil {
			return nil, fmt.Errorf("get shows failed: %w", err)
		}
	}

	return convertToShowsList(<-moviesChan, <-tvChan), nil
}
