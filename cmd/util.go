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

// convertToShowList wraps the results of the TMDb API library into the
// ShowSearchResult struct.
func convertToShowsList(movies []*common.Movie, tv []*common.TvShow) ([]*Show, error) {
	shows := make([]*Show, len(movies)+len(tv))

	for i, movie := range movies {
		shows[i] = &Show{
			ID:          movie.ID,
			Name:        movie.Title,
			ReleaseDate: parseDate(movie.ReleaseDate),
			VoteAverage: movie.VoteAverage,
		}
	}

	for i, tv := range tv {
		shows[i+len(movies)] = &Show{
			ID:          tv.ID,
			Name:        tv.Name,
			ReleaseDate: parseDate(tv.FirstAirDate),
			VoteAverage: tv.VoteAverage,
		}
	}

	return shows, nil
}

func getShows(c common.Client, name string) ([]*Show, error) {
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

	return convertToShowsList(<-moviesChan, <-tvChan)
}
