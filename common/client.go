package common

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// Client allows to retrieve info about shows.
type Client interface {
	SearchMovie(name string) ([]*Movie, error)
	SearchTvShow(name string) ([]*TvShow, error)
	GetMovie(id int) (*Movie, error)
	GetTvShow(id int) (*TvShow, error)
}

// BaseURL is the base url of the TMDb API.
const BaseURL = "https://api.themoviedb.org/3"

// DateFormat is the date format used by the TMDb API.
const DateFormat = "2006-01-02"

// TMDbClient allows to communicate with the TMDb API.
type TMDbClient struct {
	client  *http.Client
	baseURL string
	key     string
}

type movieSearchResult struct {
	Results []*Movie `json:"results"`
}

type tvSearchResult struct {
	Results []*TvShow `json:"results"`
}

type errorResult struct {
	StatusCode    int    `json:"status_code"`
	StatusMessage string `json:"status_message"`
}

// Movie represents a movie as returned by the TMDb API.
type Movie struct {
	ID          int     `json:"id"`
	Title       string  `json:"title"`
	ReleaseDate string  `json:"release_date"`
	Popularity  float32 `json:"popularity"`
	VoteAverage float32 `json:"vote_average"`
}

// TvShow represents a tv show as returned by the TMDb API.
type TvShow struct {
	ID           int     `json:"id"`
	Name         string  `json:"name"`
	FirstAirDate string  `json:"first_air_date"`
	Popularity   float32 `json:"popularity"`
	VoteAverage  float32 `json:"vote_average"`
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

func (c *TMDbClient) get(url string, result interface{}) error {
	r, err := c.client.Get(url)
	if err != nil {
		return fmt.Errorf("open connection to %s failed: %w", url, err)
	}
	defer r.Body.Close()

	decoder := json.NewDecoder(r.Body)

	if r.StatusCode < 200 || r.StatusCode >= 300 {
		var errRes errorResult
		err = decoder.Decode(&errRes)
		if err != nil {
			return fmt.Errorf("decode error body failed: %w", err)
		}
		return fmt.Errorf("error %d: %q", errRes.StatusCode, errRes.StatusMessage)
	}

	err = decoder.Decode(&result)
	if err != nil {
		return fmt.Errorf("decode result body failed: %w", err)
	}

	return nil
}

// SearchMovie searches for the given movie name in the API and returns a list
// of matching movies.
func (c *TMDbClient) SearchMovie(name string) ([]*Movie, error) {
	url := fmt.Sprintf("%s/search/movie?api_key=%s&query=%s", c.baseURL, c.key, name)

	var movies movieSearchResult
	err := c.get(url, &movies)
	if err != nil {
		return nil, fmt.Errorf("search movies failed: %w", err)
	}

	return movies.Results, nil
}

// SearchTvShow searches for the given tv show name in the API and returns a list
// of matching tv shows.
func (c *TMDbClient) SearchTvShow(name string) ([]*TvShow, error) {
	url := fmt.Sprintf("%s/search/tv?api_key=%s&query=%s", c.baseURL, c.key, name)

	var tv tvSearchResult
	err := c.get(url, &tv)
	if err != nil {
		return nil, fmt.Errorf("search tv shows failed: %w", err)
	}

	return tv.Results, nil
}

// GetMovie gets info about a movie with the given id.
func (c *TMDbClient) GetMovie(id int) (*Movie, error) {
	url := fmt.Sprintf("%s/movie/%d?api_key=%s", c.baseURL, id, c.key)

	var movie *Movie
	err := c.get(url, &movie)
	if err != nil {
		return nil, fmt.Errorf("get movie failed: %w", err)
	}

	return movie, nil
}

// GetTvShow gets info about a tv show with the given id.
func (c *TMDbClient) GetTvShow(id int) (*TvShow, error) {
	url := fmt.Sprintf("%s/tv/%d?api_key=%s", c.baseURL, id, c.key)

	var tv *TvShow
	err := c.get(url, &tv)
	if err != nil {
		return nil, fmt.Errorf("get tv show failed: %w", err)
	}

	return tv, nil
}
