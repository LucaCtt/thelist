//go:generate mockgen -destination=../mocks/mock_client.go -package=mocks github.com/lucactt/thelist/util Client

package common

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/pkg/errors"
)

// Client allows to retrieve info about shows.
type Client interface {
	Search(name string) ([]*Show, error)
}

// BaseURL is the base url of the TMDb API.
const BaseURL = "https://api.themoviedb.org/3"

// TMDbClient allows to communicate with the TMDb API.
type TMDbClient struct {
	client  *http.Client
	baseURL string
	key     string
}

type tmdbMovieSearchResult struct {
	Results      []*tmdbMovieInfo `json:"results"`
	TotalResults int              `json:"total_results"`
}

type tmdbError struct {
	StatusMessage string `json:"status_message"`
}

type tmdbMovieInfo struct {
	ID          int     `json:"id"`
	Title       string  `json:"title"`
	ReleaseDate string  `json:"release_date"`
	Popularity  float32 `json:"popularity"`
	VoteAverage float32 `json:"vote_average"`
}

// convertToShowList wraps the results of the TMDb API library into the
// ShowSearchResult struct.
func convertToShowList(result *tmdbMovieSearchResult) []*Show {
	movies := result.Results
	shows := make([]*Show, result.TotalResults)

	for i := 0; i < len(movies); i++ {
		movie := movies[i]
		shows[i] = &Show{
			ID:          movie.ID,
			Name:        movie.Title,
			ReleaseDate: movie.ReleaseDate,
			Popularity:  movie.Popularity,
			VoteAverage: movie.VoteAverage,
		}
	}
	/*
		for i := 0; i < len(tv); i++ {
			tv := tv[i]
			shows[i+len(movies)] = &Show{
				ID:          tv.ID,
				Name:        tv.Name,
				ReleaseDate: tv.FirstAirDate,
				Popularity:  tv.Popularity,
				VoteAverage: tv.VoteAverage,
			}
		}*/

	return shows
}

// NewTMDbClient creates a new api client using the given API authentication key.
func NewTMDbClient(k string, baseURL string, c *http.Client) *TMDbClient {
	return &TMDbClient{
		client:  c,
		baseURL: baseURL,
		key:     k,
	}
}

// DefaultTMDbClient creates a new TMDb client with the default base URL and http client.
// This is the recommended way to create a TMDb client.
func DefaultTMDbClient(k string) *TMDbClient {
	return NewTMDbClient(k, BaseURL, &http.Client{
		Timeout: 10 * time.Second,
	})
}

// Search searches for the given show name in both movies an tv series.
func (c *TMDbClient) Search(name string) ([]*Show, error) {
	url := fmt.Sprintf("%s/search/movie?api_key=%s&query=%s,", c.baseURL, c.key, name)

	r, err := c.client.Get(url)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("open connection to %s failed", url))
	}
	defer r.Body.Close()

	decoder := json.NewDecoder(r.Body)

	if r.StatusCode < 200 || r.StatusCode >= 300 {
		var error tmdbError
		err = decoder.Decode(&error)
		if err != nil {
			return nil, errors.Wrap(err, "decode error body failed")
		}
		return nil, errors.New(error.StatusMessage)
	}

	var result tmdbMovieSearchResult
	err = decoder.Decode(&result)
	if err != nil {
		return nil, errors.Wrap(err, "decode result body failed")
	}

	return convertToShowList(&result), nil
}
