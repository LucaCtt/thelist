//go:generate mockgen -destination=../mocks/mock_client.go -package=mocks github.com/lucactt/thelist/util Client

package common

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

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

func makeMovieResult(t *testing.T, names ...string) *tmdbMovieSearchResult {
	t.Helper()

	movies := make([]*tmdbMovieInfo, len(names))
	for i, name := range names {
		movies[i] = &tmdbMovieInfo{ID: i, Title: name}
	}

	return &tmdbMovieSearchResult{
		Results:      movies,
		TotalResults: len(movies),
	}
}

func makeShows(t *testing.T, names ...string) []*Show {
	t.Helper()

	result := make([]*Show, len(names))
	for i, name := range names {
		result[i] = &Show{ID: i, Name: name}
	}
	return result
}

func getTestClient(t *testing.T, handler http.HandlerFunc) *TMDbClient {
	t.Helper()

	server := httptest.NewServer(http.HandlerFunc(handler))
	return NewTMDbClient("test", server.URL, server.Client())
}

func TestTMDbClient_Search(t *testing.T) {
	tests := []struct {
		name    string
		handler http.HandlerFunc
		want    []*Show
		wantErr bool
	}{
		{
			name: "multiple results",
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				json.NewEncoder(w).Encode(makeMovieResult(t, "test1", "test2"))
			},
			want:    makeShows(t, "test1", "test2"),
			wantErr: false,
		},
		{
			name: "no results",
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				json.NewEncoder(w).Encode(&tmdbMovieSearchResult{})
			},
			want:    []*Show{},
			wantErr: false,
		},
		{
			name: "invalid response",
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			},
			want:    []*Show{},
			wantErr: true,
		},
		{
			name: "valid error",
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusNotFound)
				json.NewEncoder(w).Encode(&tmdbError{
					StatusMessage: "test message",
				})
			},
			want:    []*Show{},
			wantErr: true,
		},
		{
			name: "invalid error",
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusNotFound)
			},
			want:    []*Show{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := getTestClient(t, tt.handler)
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
			TotalResults: 2,
		})
	}))

	c := NewTMDbClient("test", server.URL, server.Client())
	for i := 0; i < b.N; i++ {
		c.Search("test")
	}
}
