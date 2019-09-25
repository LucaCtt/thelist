package server

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/LucaCtt/thelist/common"
	"github.com/LucaCtt/thelist/common/client"
	"github.com/LucaCtt/thelist/common/store"
	"github.com/LucaCtt/thelist/common/testutils"
	"github.com/LucaCtt/thelist/mocks"
	"github.com/go-chi/chi"
	"github.com/stretchr/testify/mock"
)

func assertListsEqual(t *testing.T, got, want []*common.Show) {
	t.Helper()

	testutils.AssertLenEqual(t, len(got), len(want))
	for i, s := range got {
		if !reflect.DeepEqual(s, want[i]) {
			t.Errorf("got[%d] %+v, want[%d] %+v", i, s, i, want[i])
		}
	}
}

func TestServer_handleAPIGetAll(t *testing.T) {
	tests := []struct {
		name   string
		items  []*store.Item
		movie  *client.Movie
		want   []*common.Show
		status int
	}{
		{
			name: "items found",
			items: []*store.Item{
				&store.Item{ID: 1, Type: store.MovieType, ShowID: 1},
				&store.Item{ID: 2, Type: store.MovieType, ShowID: 2},
			},
			movie: &client.Movie{ID: 1, Title: "test1"},
			want: []*common.Show{
				&common.Show{ID: 1, Type: store.MovieType, Name: "test1"},
				&common.Show{ID: 1, Type: store.MovieType, Name: "test1"},
			},
			status: http.StatusOK,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := &mocks.Store{}
			client := &mocks.Client{}

			store.On("All").Return(tt.items, nil)
			client.On("GetMovie", mock.Anything).Return(tt.movie, nil)

			s := &Server{store, client, chi.NewRouter()}
			s.router.HandleFunc("/", s.handleAPIGetAll())

			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodGet, "/", nil)
			s.ServeHTTP(w, r)

			testutils.AssertEqual(t, w.Result().StatusCode, tt.status)
			var got []*common.Show
			err := json.NewDecoder(w.Body).Decode(&got)
			testutils.AssertErr(t, err, false)
			assertListsEqual(t, got, tt.want)
		})
	}
}
