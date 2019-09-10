// Package client implements a REST API client for retrieving
// information on shows.
package client

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

// TMDbClient allows to communicate with the TMDb API.
type TMDbClient struct {
	client  *http.Client
	baseURL string
	key     string
}

// TMDb api constants.
const (
	BaseURL    = "https://api.themoviedb.org/3"
	DateFormat = "2006-01-02"
)

// Custom creates a new custom TMDb API client.
func Custom(k string, baseURL string, c *http.Client) *TMDbClient {
	return &TMDbClient{
		client:  c,
		baseURL: baseURL,
		key:     k,
	}
}

// New creates a new TMDb API client with the default base URL and http client.
// This is the recommended way to create a TMDb client.
func New(k string) *TMDbClient {
	return Custom(k, BaseURL, &http.Client{
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
