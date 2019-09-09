package cmd

import (
	"fmt"
	"time"

	"github.com/LucaCtt/thelist/common"
)

// Show is a common representation for movies and tv series.
type Show struct {
	ID          int
	Type        string
	Name        string
	ReleaseDate time.Time
	VoteAverage float32
}

func parseDate(date string) time.Time {
	parsed, err := time.Parse(common.DateFormat, date)
	if err != nil {
		panic(fmt.Errorf("parse movie release date failed: %w", err))
	}
	return parsed
}

func movieToShow(movie *common.Movie) *Show {
	return &Show{
		ID:          movie.ID,
		Type:        common.MovieType,
		Name:        movie.Title,
		ReleaseDate: parseDate(movie.ReleaseDate),
		VoteAverage: movie.VoteAverage,
	}
}

func tvToShow(tv *common.TvShow) *Show {
	return &Show{
		ID:          tv.ID,
		Type:        common.TvShowType,
		Name:        tv.Name,
		ReleaseDate: parseDate(tv.FirstAirDate),
		VoteAverage: tv.VoteAverage,
	}
}

func convertToShowsList(movies []*common.Movie, tv []*common.TvShow) []*Show {
	shows := make([]*Show, len(movies)+len(tv))

	for i, movie := range movies {
		shows[i] = movieToShow(movie)
	}

	for i, tv := range tv {
		shows[i+len(movies)] = tvToShow(tv)
	}

	return shows
}

func getShow(c common.Client, id int, t string) (*Show, error) {
	var show *Show

	switch t {
	case common.MovieType:
		movie, err := c.GetMovie(id)
		if err != nil {
			return nil, fmt.Errorf("get show with id %d and type %q failed: %w", id, t, err)
		}
		show = movieToShow(movie)
	case common.TvShowType:
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

func searchShow(c common.Client, name string) ([]*Show, error) {
	moviesChan := make(chan []*common.Movie)
	tvChan := make(chan []*common.TvShow)
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
