package cmd

import (
	"errors"
	"testing"

	"github.com/LucaCtt/thelist/common"
	"github.com/LucaCtt/thelist/mocks"
	"github.com/stretchr/testify/mock"
)

func assertErr(t *testing.T, got error, wantErr bool) {
	t.Helper()

	if (got != nil) != wantErr {
		t.Errorf("got error %v, wantErr %v", got, wantErr)
	}
}

func Test_add(t *testing.T) {
	t.Run("name in args", func(t *testing.T) {
		show := "test"
		p := &mocks.Prompter{}
		c := &mocks.Client{}
		s := &mocks.Store{}

		s.On("Create", mock.Anything).Return(nil)
		c.On("Search", show).Return([]*common.Show{&common.Show{ID: 1, Name: "test1"}}, nil)

		err := add([]string{show}, p, c, s)

		assertErr(t, err, false)
		mock.AssertExpectationsForObjects(t, p, c, s)
	})

	t.Run("name not in args", func(t *testing.T) {
		show := "test"
		p := &mocks.Prompter{}
		c := &mocks.Client{}
		s := &mocks.Store{}

		p.On("Input", mock.Anything).Return(show, nil)
		s.On("Create", mock.Anything).Return(nil)
		c.On("Search", show).Return([]*common.Show{&common.Show{ID: 1}}, nil)

		err := add([]string{}, p, c, s)

		assertErr(t, err, false)
		mock.AssertExpectationsForObjects(t, p, c, s)
	})

	t.Run("name empty", func(t *testing.T) {
		p := &mocks.Prompter{}
		c := &mocks.Client{}
		s := &mocks.Store{}

		p.On("Input", mock.Anything).Return("", nil)

		err := add([]string{}, p, c, s)

		assertErr(t, err, true)
		mock.AssertExpectationsForObjects(t, p, c, s)
	})

	t.Run("no shows found", func(t *testing.T) {
		show := "test"
		p := &mocks.Prompter{}
		c := &mocks.Client{}
		s := &mocks.Store{}

		c.On("Search", show).Return([]*common.Show{}, nil)

		err := add([]string{show}, p, c, s)

		assertErr(t, err, true)
		mock.AssertExpectationsForObjects(t, p, c, s)
	})

	t.Run("multiple shows found", func(t *testing.T) {
		show := "test"
		p := &mocks.Prompter{}
		c := &mocks.Client{}
		s := &mocks.Store{}

		p.On("Select", mock.Anything, mock.Anything).Return(0, nil)
		s.On("Create", mock.Anything).Return(nil)
		c.On("Search", show).Return([]*common.Show{&common.Show{ID: 1}, &common.Show{ID: 2}}, nil)

		err := add([]string{show}, p, c, s)

		assertErr(t, err, false)
		mock.AssertExpectationsForObjects(t, p, c, s)
	})

	t.Run("prompter input err", func(t *testing.T) {
		p := &mocks.Prompter{}
		c := &mocks.Client{}
		s := &mocks.Store{}

		p.On("Input", mock.Anything, mock.Anything).Return("", errors.New("test"))

		err := add([]string{}, p, c, s)

		assertErr(t, err, true)
		mock.AssertExpectationsForObjects(t, p, c, s)
	})

	t.Run("prompter select err", func(t *testing.T) {
		show := "test"
		p := &mocks.Prompter{}
		c := &mocks.Client{}
		s := &mocks.Store{}

		p.On("Select", mock.Anything, mock.Anything).Return(0, errors.New("test"))
		c.On("Search", show).Return([]*common.Show{&common.Show{ID: 1}, &common.Show{ID: 2}}, nil)

		err := add([]string{show}, p, c, s)

		assertErr(t, err, true)
		mock.AssertExpectationsForObjects(t, p, c, s)
	})

	t.Run("client err", func(t *testing.T) {
		show := "test"
		p := &mocks.Prompter{}
		c := &mocks.Client{}
		s := &mocks.Store{}

		c.On("Search", show).Return([]*common.Show{}, errors.New("test"))

		err := add([]string{show}, p, c, s)

		assertErr(t, err, true)
		mock.AssertExpectationsForObjects(t, p, c, s)
	})

	t.Run("store err", func(t *testing.T) {
		show := "test"
		p := &mocks.Prompter{}
		c := &mocks.Client{}
		s := &mocks.Store{}

		s.On("Create", mock.Anything).Return(errors.New("test"))
		c.On("Search", show).Return([]*common.Show{&common.Show{ID: 1, Name: "test1"}}, nil)

		err := add([]string{show}, p, c, s)

		assertErr(t, err, true)
		mock.AssertExpectationsForObjects(t, p, c, s)
	})
}
