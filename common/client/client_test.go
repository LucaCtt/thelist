package client

import (
	"encoding/json"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/LucaCtt/thelist/common/testutils"
)

func assertMoviesListEqual(t *testing.T, got, want []*Movie) {
	t.Helper()

	testutils.AssertLenEqual(t, len(got), len(want))
	for i, s := range got {
		if !reflect.DeepEqual(s, want[i]) {
			t.Errorf("got[%d] %+v, want[%d] %+v", i, s, i, want[i])
		}
	}
}

func assertTvShowsListEqual(t *testing.T, got, want []*TvShow) {
	t.Helper()

	testutils.AssertLenEqual(t, len(got), len(want))
	for i, s := range got {
		if !reflect.DeepEqual(s, want[i]) {
			t.Errorf("got[%d] %+v, want[%d] %+v", i, s, i, want[i])
		}
	}
}

func makeTvResult(names ...string) *tvSearchResult {
	tv := make([]*TvShow, len(names))
	for i, name := range names {
		tv[i] = &TvShow{Name: name}
	}

	return &tvSearchResult{
		Results: tv,
	}
}

func makeMovieResult(names ...string) *movieSearchResult {
	movies := make([]*Movie, len(names))
	for i, name := range names {
		movies[i] = &Movie{Title: name}
	}

	return &movieSearchResult{
		Results: movies,
	}
}

func makeMovies(names ...string) []*Movie {
	result := make([]*Movie, len(names))
	for i, name := range names {
		result[i] = &Movie{Title: name}
	}
	return result
}

func makeTvShows(names ...string) []*TvShow {
	result := make([]*TvShow, len(names))
	for i, name := range names {
		result[i] = &TvShow{Name: name}
	}
	return result
}

func getTestClient(handler http.HandlerFunc) *TMDbClient {
	server := httptest.NewServer(handler)
	return Custom("test", server.URL, server.Client())
}

func TestTMDbClient_SearchMovie(t *testing.T) {
	tests := []struct {
		name    string
		handler http.HandlerFunc
		want    []*Movie
		wantErr bool
	}{
		{
			name: "multiple results",
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				json.NewEncoder(w).Encode(makeMovieResult("test1", "test2"))
			},
			want:    makeMovies("test1", "test2"),
			wantErr: false,
		},
		{
			name: "no results",
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				json.NewEncoder(w).Encode(&movieSearchResult{})
			},
			want:    nil,
			wantErr: false,
		},
		{
			name: "invalid response",
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "valid error",
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusNotFound)
				json.NewEncoder(w).Encode(&errorResult{
					StatusMessage: "test message 1",
				})
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "invalid error",
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusNotFound)
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := getTestClient(tt.handler)
			got, err := c.SearchMovie("test")
			testutils.AssertErr(t, err, tt.wantErr)
			assertMoviesListEqual(t, got, tt.want)
		})
	}
	t.Run("invalid baseurl", func(t *testing.T) {
		c := Custom("test", "localhost:999999", &http.Client{})
		got, err := c.SearchMovie("test")
		testutils.AssertErr(t, err, true)
		assertMoviesListEqual(t, got, []*Movie{})
	})
}

func TestTMDbClient_SearchTvShow(t *testing.T) {
	tests := []struct {
		name    string
		handler http.HandlerFunc
		want    []*TvShow
		wantErr bool
	}{
		{
			name: "multiple movie results",
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				json.NewEncoder(w).Encode(makeTvResult("test3", "test4"))
			},
			want:    makeTvShows("test3", "test4"),
			wantErr: false,
		},
		{
			name: "no results",
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				json.NewEncoder(w).Encode(&tvSearchResult{})
			},
			want:    nil,
			wantErr: false,
		},
		{
			name: "invalid response",
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "valid error",
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusNotFound)
				json.NewEncoder(w).Encode(&errorResult{
					StatusMessage: "test message 2",
				})
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "invalid error",
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusNotFound)
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := getTestClient(tt.handler)
			got, err := c.SearchTvShow("test")
			testutils.AssertErr(t, err, tt.wantErr)
			assertTvShowsListEqual(t, got, tt.want)
		})
	}
	t.Run("invalid baseurl", func(t *testing.T) {
		c := Custom("test", "localhost:999999", &http.Client{})
		got, err := c.SearchTvShow("test")
		testutils.AssertErr(t, err, true)
		assertTvShowsListEqual(t, got, []*TvShow{})
	})
}

func TestTMDbClient_GetMovie(t *testing.T) {
	tests := []struct {
		name    string
		id      int
		handler http.HandlerFunc
		want    *Movie
		wantErr bool
	}{
		{
			name: "movie exists",
			id:   1,
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				json.NewEncoder(w).Encode(Movie{ID: 1, Title: "test"})
			},
			want:    &Movie{ID: 1, Title: "test"},
			wantErr: false,
		},
		{
			name: "movie does not exists",
			id:   1,
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusNotFound)
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := getTestClient(tt.handler)
			got, err := c.GetMovie(tt.id)
			testutils.AssertErr(t, err, tt.wantErr)
			testutils.AssertEqual(t, got, tt.want)
		})
	}
}

func TestTMDbClient_GetTvShow(t *testing.T) {
	tests := []struct {
		name    string
		id      int
		handler http.HandlerFunc
		want    *TvShow
		wantErr bool
	}{
		{
			name: "tv show exists",
			id:   1,
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				json.NewEncoder(w).Encode(TvShow{ID: 1, Name: "test"})
			},
			want:    &TvShow{ID: 1, Name: "test"},
			wantErr: false,
		},
		{
			name: "tv show does not exists",
			id:   1,
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusNotFound)
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := getTestClient(tt.handler)
			got, err := c.GetTvShow(tt.id)
			testutils.AssertErr(t, err, tt.wantErr)
			testutils.AssertEqual(t, got, tt.want)
		})
	}
}

func randomString(len int) string {
	bytes := make([]byte, len)
	for i := 0; i < len; i++ {
		bytes[i] = byte(rand.Int())
	}
	return string(bytes)
}
