package cmd

import (
	"errors"
	"fmt"
	"testing"

	"github.com/LucaCtt/thelist/common/client"
	"github.com/LucaCtt/thelist/mocks"
	"github.com/stretchr/testify/mock"
)

func assertErr(t *testing.T, got error, wantErr bool) {
	t.Helper()

	if (got != nil) != wantErr {
		t.Errorf("got %v, wantErr %v", got, wantErr)
	}
}

func genMovies(len int) []*client.Movie {
	res := make([]*client.Movie, len)

	for i := 0; i < len; i++ {
		res[i] = &client.Movie{
			ID:          i,
			Title:       fmt.Sprintf("test%d", i),
			ReleaseDate: "2001-01-01",
		}
	}

	return res
}

func genTvShows(len int) []*client.TvShow {
	res := make([]*client.TvShow, len)

	for i := 0; i < len; i++ {
		res[i] = &client.TvShow{
			ID:           i,
			Name:         fmt.Sprintf("test%d", i),
			FirstAirDate: "2001-01-01",
		}
	}

	return res
}

func Test_add(t *testing.T) {
	t.Run("name in args", func(t *testing.T) {
		show := "test"
		p := &mocks.Prompt{}
		c := &mocks.Client{}
		s := &mocks.Store{}

		s.On("Create", mock.Anything).Return(nil)
		c.On("SearchMovie", show).Return(genMovies(1), nil)
		c.On("SearchTvShow", show).Return([]*client.TvShow{}, nil)

		err := add([]string{show}, p, c, s)

		assertErr(t, err, false)
		mock.AssertExpectationsForObjects(t, p, c, s)
	})

	t.Run("name not in args", func(t *testing.T) {
		show := "test"
		p := &mocks.Prompt{}
		c := &mocks.Client{}
		s := &mocks.Store{}

		p.On("Input", mock.Anything).Return(show, nil)
		s.On("Create", mock.Anything).Return(nil)
		c.On("SearchMovie", show).Return(genMovies(1), nil)
		c.On("SearchTvShow", show).Return([]*client.TvShow{}, nil)

		err := add([]string{}, p, c, s)

		assertErr(t, err, false)
		mock.AssertExpectationsForObjects(t, p, c, s)
	})

	t.Run("name empty", func(t *testing.T) {
		p := &mocks.Prompt{}
		c := &mocks.Client{}
		s := &mocks.Store{}

		p.On("Input", mock.Anything).Return("", nil)

		err := add([]string{}, p, c, s)

		assertErr(t, err, true)
		mock.AssertExpectationsForObjects(t, p, c, s)
	})

	t.Run("no shows found", func(t *testing.T) {
		show := "test"
		p := &mocks.Prompt{}
		c := &mocks.Client{}
		s := &mocks.Store{}

		c.On("SearchMovie", show).Return([]*client.Movie{}, nil)
		c.On("SearchTvShow", show).Return([]*client.TvShow{}, nil)

		err := add([]string{show}, p, c, s)

		assertErr(t, err, true)
		mock.AssertExpectationsForObjects(t, p, c, s)
	})

	t.Run("multiple shows found", func(t *testing.T) {
		show := "test"
		p := &mocks.Prompt{}
		c := &mocks.Client{}
		s := &mocks.Store{}

		p.On("Select", mock.Anything, mock.Anything).Return(0, nil)
		s.On("Create", mock.Anything).Return(nil)
		c.On("SearchMovie", show).Return(genMovies(2), nil)
		c.On("SearchTvShow", show).Return(genTvShows(2), nil)

		err := add([]string{show}, p, c, s)

		assertErr(t, err, false)
		mock.AssertExpectationsForObjects(t, p, c, s)
	})

	t.Run("prompter input err", func(t *testing.T) {
		p := &mocks.Prompt{}
		c := &mocks.Client{}
		s := &mocks.Store{}

		p.On("Input", mock.Anything, mock.Anything).Return("", errors.New("test"))

		err := add([]string{}, p, c, s)

		assertErr(t, err, true)
		mock.AssertExpectationsForObjects(t, p, c, s)
	})

	t.Run("prompter select err", func(t *testing.T) {
		show := "test"
		p := &mocks.Prompt{}
		c := &mocks.Client{}
		s := &mocks.Store{}

		p.On("Select", mock.Anything, mock.Anything).Return(0, errors.New("test"))
		c.On("SearchMovie", show).Return(genMovies(1), nil)
		c.On("SearchTvShow", show).Return(genTvShows(1), nil)

		err := add([]string{show}, p, c, s)

		assertErr(t, err, true)
		mock.AssertExpectationsForObjects(t, p, c, s)
	})

	t.Run("client err", func(t *testing.T) {
		show := "test"
		p := &mocks.Prompt{}
		c := &mocks.Client{}
		s := &mocks.Store{}

		c.On("SearchMovie", show).Return(nil, errors.New("test"))
		c.On("SearchTvShow", show).Return([]*client.TvShow{}, nil)

		err := add([]string{show}, p, c, s)

		assertErr(t, err, true)
		mock.AssertExpectationsForObjects(t, p, c, s)
	})

	t.Run("store err", func(t *testing.T) {
		show := "test"
		p := &mocks.Prompt{}
		c := &mocks.Client{}
		s := &mocks.Store{}

		s.On("Create", mock.Anything).Return(errors.New("test"))
		c.On("SearchMovie", show).Return(genMovies(1), nil)
		c.On("SearchTvShow", show).Return([]*client.TvShow{}, nil)

		err := add([]string{show}, p, c, s)

		assertErr(t, err, true)
		mock.AssertExpectationsForObjects(t, p, c, s)
	})
}
