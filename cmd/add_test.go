package cmd

import (
	"testing"

	"github.com/LucaCtt/thelist/data"
	"github.com/LucaCtt/thelist/mocks"
	"github.com/golang/mock/gomock"
)

func TestAdd_ShowInArgsMultipleSearchResults_CreatesShow(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	p := mocks.NewMockPrompt(ctrl)
	c := mocks.NewMockClient(ctrl)
	s := mocks.NewMockStore(ctrl)
	result := &data.ShowList{TotalResults: 2, Results: []*data.Show{
		&data.Show{
			ID: 1,
		},
		&data.Show{
			ID: 2,
		},
	}}

	c.EXPECT().SearchShow(gomock.Eq("a")).Return(result, nil)
	p.EXPECT().SelectShow(gomock.Eq(result)).Return(result.Results[0], nil)
	s.EXPECT().CreateItem(gomock.Any()).Return(nil)

	err := add([]string{"a"}, p, c, s)

	if err != nil {
		t.Errorf("Unexpected error %v", err)
	}
}

func TestAdd_ShowInArgsNoSearchResults_Exits(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	p := mocks.NewMockPrompt(ctrl)
	c := mocks.NewMockClient(ctrl)
	s := mocks.NewMockStore(ctrl)
	result := &data.ShowList{TotalResults: 0, Results: []*data.Show{}}

	c.EXPECT().SearchShow(gomock.Eq("a")).Return(result, nil)

	err := add([]string{"a"}, p, c, s)

	if err != nil {
		t.Errorf("Unexpected error %v", err)
	}
}
