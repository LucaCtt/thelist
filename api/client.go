package api

import "github.com/ryanbradynd05/go-tmdb"

// Client allows to communicate with the TMDb API.
// A new one should be created by using the NewClient method.
type Client struct {
	client *tmdb.TMDb
}

// ShowSearchResult represents the result of a show search.
type ShowSearchResult struct {
	Results      []*ShowSearchInfo
	Names        []string
	TotalResults int
}

// ShowSearchInfo represents the info of a show returned by a search query.
type ShowSearchInfo struct {
	ID          int
	Name        string
	ReleaseDate string
	Popularity  float32
	VoteAverage float32
}

// convertToShowSearchResult wraps the results of the TMDb API library into the
// ShowSearchResult struct.
func convertToShowSearchResult(result *tmdb.MultiSearchResults) *ShowSearchResult {
	movies := result.GetMoviesResults()
	tv := result.GetTvResults()
	shows := make([]*ShowSearchInfo, result.TotalResults)

	for i := 0; i < len(movies); i++ {
		movie := movies[i]
		shows[i] = &ShowSearchInfo{
			ID:          movie.ID,
			Name:        movie.Title,
			ReleaseDate: movie.ReleaseDate,
			Popularity:  movie.Popularity,
			VoteAverage: movie.VoteAverage,
		}
	}

	for i := 0; i < len(tv); i++ {
		tv := tv[i]
		shows[i+len(movies)] = &ShowSearchInfo{
			ID:          tv.ID,
			Name:        tv.Name,
			ReleaseDate: tv.FirstAirDate,
			Popularity:  tv.Popularity,
			VoteAverage: tv.VoteAverage,
		}
	}

	return &ShowSearchResult{
		Results:      shows,
		TotalResults: result.TotalResults,
	}
}

// NewClient creates a new api client using the given API authentication key.
func NewClient(apiKey string) *Client {
	return &Client{tmdb.Init(tmdb.Config{
		APIKey: apiKey,
	})}
}

// SearchShow searches for the given show name in both movies an tv series.
func (c *Client) SearchShow(name string) (*ShowSearchResult, error) {
	result, err := c.client.SearchMulti(name, nil)
	if err != nil {
		return nil, err
	}

	return convertToShowSearchResult(result), nil
}
