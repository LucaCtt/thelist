package cmd

import (
	"errors"
	"testing"

	"github.com/LucaCtt/thelist/common"
	"github.com/LucaCtt/thelist/mocks"
	"github.com/stretchr/testify/mock"
)

func Test_list(t *testing.T) {
	t.Run("items in store", func(t *testing.T) {
		p := &mocks.Prompter{}
		c := &mocks.Client{}
		s := &mocks.Store{}

		p.On("MultiSelect", mock.Anything, mock.Anything).Return([]int{0}, nil)
		c.On("GetMovie", 1).Return(genMovies(1)[0], nil)
		s.On("All").Return([]*common.Item{&common.Item{ID: 1, ShowID: 1, Type: common.MovieType}}, nil)
		s.On("Delete", uint(1)).Return(nil)

		err := list(p, c, s)
		assertErr(t, err, false)
		mock.AssertExpectationsForObjects(t, p, c, s)
	})

	t.Run("store empty", func(t *testing.T) {
		p := &mocks.Prompter{}
		c := &mocks.Client{}
		s := &mocks.Store{}

		s.On("All").Return([]*common.Item{}, nil)

		err := list(p, c, s)
		assertErr(t, err, true)
		mock.AssertExpectationsForObjects(t, p, c, s)
	})

	t.Run("movie not found", func(t *testing.T) {
		p := &mocks.Prompter{}
		c := &mocks.Client{}
		s := &mocks.Store{}

		c.On("GetMovie", 1).Return(nil, errors.New("test"))
		s.On("All").Return([]*common.Item{&common.Item{ID: 1, ShowID: 1, Type: common.MovieType}}, nil)

		err := list(p, c, s)
		assertErr(t, err, true)
		mock.AssertExpectationsForObjects(t, p, c, s)
	})

	t.Run("store get all error", func(t *testing.T) {
		p := &mocks.Prompter{}
		c := &mocks.Client{}
		s := &mocks.Store{}

		s.On("All").Return([]*common.Item{}, errors.New("test"))

		err := list(p, c, s)
		assertErr(t, err, true)
		mock.AssertExpectationsForObjects(t, p, c, s)
	})

	t.Run("store delete error", func(t *testing.T) {
		p := &mocks.Prompter{}
		c := &mocks.Client{}
		s := &mocks.Store{}

		p.On("MultiSelect", mock.Anything, mock.Anything).Return([]int{0}, nil)
		c.On("GetMovie", 1).Return(genMovies(1)[0], nil)
		s.On("All").Return([]*common.Item{&common.Item{ID: 1, ShowID: 1, Type: common.MovieType}}, nil)
		s.On("Delete", uint(1)).Return(errors.New("test"))

		err := list(p, c, s)
		assertErr(t, err, true)
		mock.AssertExpectationsForObjects(t, p, c, s)
	})

	t.Run("prompter error", func(t *testing.T) {
		p := &mocks.Prompter{}
		c := &mocks.Client{}
		s := &mocks.Store{}

		p.On("MultiSelect", mock.Anything, mock.Anything).Return([]int{}, errors.New("test"))
		c.On("GetMovie", 1).Return(genMovies(1)[0], nil)
		s.On("All").Return([]*common.Item{&common.Item{ID: 1, ShowID: 1, Type: common.MovieType}}, nil)

		err := list(p, c, s)
		assertErr(t, err, true)
		mock.AssertExpectationsForObjects(t, p, c, s)
	})
}
