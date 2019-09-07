package common

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"
)

var testTimeStr = "2000-01-01"
var testTime, _ = time.Parse(tmdbDateFormat, testTimeStr)

func assertShowsEqual(t *testing.T, got, want []*Show) {
	t.Helper()

	if len(got) != len(want) {
		t.Errorf("got has len %d, want %d", len(got), len(want))
		return
	}

	for i, s := range got {
		if !reflect.DeepEqual(*s, *want[i]) {
			t.Errorf("got[%d] %+v, want[%d] %+v", i, *s, i, *want[i])
		}
	}
}

func makeTvResult(t *testing.T, names ...string) *tmdbTvSearchResult {
	t.Helper()

	tv := make([]*tmdbTvInfo, len(names))
	for i, name := range names {
		tv[i] = &tmdbTvInfo{Name: name, FirstAirDate: testTimeStr}
	}

	return &tmdbTvSearchResult{
		Results: tv,
	}
}

func makeMovieResult(t *testing.T, names ...string) *tmdbMovieSearchResult {
	t.Helper()

	movies := make([]*tmdbMovieInfo, len(names))
	for i, name := range names {
		movies[i] = &tmdbMovieInfo{Title: name, ReleaseDate: testTimeStr}
	}

	return &tmdbMovieSearchResult{
		Results: movies,
	}
}

func makeShows(t *testing.T, showType ShowType, names ...string) []*Show {
	t.Helper()

	result := make([]*Show, len(names))
	for i, name := range names {
		result[i] = &Show{Name: name, Type: showType, ReleaseDate: testTime}
	}
	return result
}

func getTestClient(t *testing.T, movieHandler, tvHandler http.HandlerFunc) *TMDbClient {
	t.Helper()

	mux := http.NewServeMux()
	mux.HandleFunc("/search/movie", movieHandler)
	mux.HandleFunc("/search/tv", tvHandler)
	server := httptest.NewServer(mux)
	return NewTMDbClient("test", server.URL, server.Client())
}

func TestTMDbClient_Search(t *testing.T) {
	tests := []struct {
		name         string
		movieHandler http.HandlerFunc
		tvHandler    http.HandlerFunc
		want         []*Show
		wantErr      bool
	}{
		{
			name: "multiple movie results",
			movieHandler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				json.NewEncoder(w).Encode(makeMovieResult(t, "test1", "test2"))
			},
			tvHandler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				json.NewEncoder(w).Encode(makeTvResult(t, "test3", "test4"))
			},
			want:    append(makeShows(t, MovieType, "test1", "test2"), makeShows(t, TvShowType, "test3", "test4")...),
			wantErr: false,
		},
		{
			name: "no results",
			movieHandler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				json.NewEncoder(w).Encode(&tmdbMovieSearchResult{})
			},
			tvHandler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				json.NewEncoder(w).Encode(&tmdbTvSearchResult{})
			},
			want:    []*Show{},
			wantErr: false,
		},
		{
			name: "invalid movie response",
			movieHandler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			},
			tvHandler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				json.NewEncoder(w).Encode(&tmdbTvSearchResult{})
			},
			want:    []*Show{},
			wantErr: true,
		},
		{
			name: "invalid tv response",
			movieHandler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				json.NewEncoder(w).Encode(&tmdbMovieSearchResult{})
			},
			tvHandler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			},
			want:    []*Show{},
			wantErr: true,
		},
		{
			name: "valid error",
			movieHandler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusNotFound)
				json.NewEncoder(w).Encode(&tmdbError{
					StatusMessage: "test message 1",
				})
			},
			tvHandler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusNotFound)
				json.NewEncoder(w).Encode(&tmdbError{
					StatusMessage: "test message 2",
				})
			},
			want:    []*Show{},
			wantErr: true,
		},
		{
			name: "invalid error",
			movieHandler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusNotFound)
			},
			tvHandler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusNotFound)
			},
			want:    []*Show{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := getTestClient(t, tt.movieHandler, tt.tvHandler)
			got, err := c.Search("test")
			assertErr(t, err, tt.wantErr)
			assertShowsEqual(t, got, tt.want)
		})
	}
	t.Run("invalid baseurl", func(t *testing.T) {
		c := NewTMDbClient("test", "localhost:999999", &http.Client{})
		got, err := c.Search("test")
		assertErr(t, err, true)
		assertShowsEqual(t, got, []*Show{})
	})
}

func BenchmarkTMDbClient_Search(b *testing.B) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(&tmdbMovieSearchResult{
			Results: []*tmdbMovieInfo{
				&tmdbMovieInfo{
					ID:    1,
					Title: "test1",
				},
				&tmdbMovieInfo{
					ID:    1,
					Title: "test2",
				},
			},
		})
	}))

	c := NewTMDbClient("test", server.URL, server.Client())
	for i := 0; i < b.N; i++ {
		c.Search("test")
	}
}
