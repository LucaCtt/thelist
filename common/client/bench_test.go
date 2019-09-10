package client

import (
	"encoding/json"
	"math/rand"
	"net/http"
	"testing"
)

func BenchmarkTMDbClient_SearchMovie(b *testing.B) {
	const len = 100

	data := make([]*Movie, len)
	for i := 0; i < len; i++ {
		id := rand.Int()
		data[i] = &Movie{
			ID:    id,
			Title: randomString(10),
		}
	}

	c := getTestClient(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(&movieSearchResult{
			Results: data,
		})
	})

	for i := 0; i < b.N; i++ {
		c.SearchMovie("test")
	}
}

func BenchmarkTMDbClient_SearchTvShow(b *testing.B) {
	const len = 100

	data := make([]*TvShow, len)
	for i := 0; i < len; i++ {
		id := rand.Int()
		data[i] = &TvShow{
			ID:   id,
			Name: randomString(10),
		}
	}

	c := getTestClient(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(&tvSearchResult{
			Results: data,
		})
	})

	for i := 0; i < b.N; i++ {
		c.SearchTvShow("test")
	}
}
