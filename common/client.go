package common

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
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
	Results []*tmdbMovieInfo `json:"results"`
}

type tmdbTvSearchResult struct {
	Results []*tmdbTvInfo `json:"results"`
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

// convertToShowList wraps the results of the TMDb API library into the
// ShowSearchResult struct.
func convertToShowsList(movies []*tmdbMovieInfo, tv []*tmdbTvInfo) []*Show {
	shows := make([]*Show, len(movies)+len(tv))

	for i, movie := range movies {
		shows[i] = &Show{
			ID:          movie.ID,
			Name:        movie.Title,
			ReleaseDate: movie.ReleaseDate,
			Popularity:  movie.Popularity,
			VoteAverage: movie.VoteAverage,
		}
	}

	for i, tv := range tv {
		shows[i+len(movies)] = &Show{
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
		return fmt.Errorf("open connection to %s failed: %w", url, err)
	}
	defer r.Body.Close()

	decoder := json.NewDecoder(r.Body)

	if r.StatusCode < 200 || r.StatusCode >= 300 {
		var error tmdbError
		err = decoder.Decode(&error)
		if err != nil {
			return fmt.Errorf("decode error body failed: %w", err)
		}
		return fmt.Errorf(error.StatusMessage)
	}

	err = decoder.Decode(&result)
	if err != nil {
		return fmt.Errorf("decode result body failed: %w", err)
	}

	return nil
}

// Search searches for the given show name in both movies and tv series.
func (c *TMDbClient) Search(name string) ([]*Show, error) {
	moviesURL := fmt.Sprintf("%s/search/movie?api_key=%s&query=%s,", c.baseURL, c.key, name)
	tvURL := fmt.Sprintf("%s/search/tv?api_key=%s&query=%s,", c.baseURL, c.key, name)
	errChan := make(chan error, 2)

	var movies tmdbMovieSearchResult
	go func() {
		errChan <- c.doRequest(moviesURL, &movies)
	}()

	var tv tmdbTvSearchResult
	go func() {
		errChan <- c.doRequest(tvURL, &tv)
	}()

	for i := 0; i < 2; i++ {
		err := <-errChan
		if err != nil {
			return nil, fmt.Errorf("get shows failed: %w", err)
		}
	}

	return convertToShowsList(movies.Results, tv.Results), nil
}
