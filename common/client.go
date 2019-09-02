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

type tmdbTvSearchResult struct {
	Results      []*tmdbTvInfo `json:"results"`
	TotalResults int           `json:"total_results"`
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

type tmdbTvInfo struct {
	ID           int     `json:"id"`
	Name         string  `json:"name"`
	FirstAirDate string  `json:"first_air_date"`
	Popularity   float32 `json:"popularity"`
	VoteAverage  float32 `json:"vote_average"`
}

type tmdbSearchResult struct {
	MovieSearchResult *tmdbMovieSearchResult
	TvSearchResult    *tmdbTvSearchResult
	TotalResults      int
}

// convertToShowList wraps the results of the TMDb API library into the
// ShowSearchResult struct.
func convertToShowList(result *tmdbSearchResult) []*Show {
	shows := make([]*Show, result.TotalResults)

	for i, movie := range result.MovieSearchResult.Results {
		shows[i] = &Show{
			ID:          movie.ID,
			Name:        movie.Title,
			ReleaseDate: movie.ReleaseDate,
			Popularity:  movie.Popularity,
			VoteAverage: movie.VoteAverage,
		}
	}

	for i, tv := range result.TvSearchResult.Results {
		shows[i+result.MovieSearchResult.TotalResults] = &Show{
			ID:          tv.ID,
			Name:        tv.Name,
			ReleaseDate: tv.FirstAirDate,
			Popularity:  tv.Popularity,
			VoteAverage: tv.VoteAverage,
		}
	}

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

func (c *TMDbClient) doRequest(url string, result interface{}) error {
	r, err := c.client.Get(url)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("open connection to %s failed", url))
	}
	defer r.Body.Close()

	decoder := json.NewDecoder(r.Body)

	if r.StatusCode < 200 || r.StatusCode >= 300 {
		var error tmdbError
		err = decoder.Decode(&error)
		if err != nil {
			return errors.Wrap(err, "decode error body failed")
		}
		return errors.New(error.StatusMessage)
	}

	err = decoder.Decode(&result)
	if err != nil {
		return errors.Wrap(err, "decode result body failed")
	}

	return nil
}

// Search searches for the given show name in both movies and tv series.
func (c *TMDbClient) Search(name string) ([]*Show, error) {
	moviesURL := fmt.Sprintf("%s/search/movie?api_key=%s&query=%s,", c.baseURL, c.key, name)
	tvURL := fmt.Sprintf("%s/search/tv?api_key=%s&query=%s,", c.baseURL, c.key, name)

	var movies tmdbMovieSearchResult
	err := c.doRequest(moviesURL, &movies)
	if err != nil {
		return nil, errors.Wrap(err, "get movies failed")
	}

	var tv tmdbTvSearchResult
	err = c.doRequest(tvURL, &tv)
	if err != nil {
		return nil, errors.Wrap(err, "get tv shows failed")
	}

	return convertToShowList(&tmdbSearchResult{
		MovieSearchResult: &movies,
		TvSearchResult:    &tv,
		TotalResults:      movies.TotalResults + tv.TotalResults,
	}), nil
}
