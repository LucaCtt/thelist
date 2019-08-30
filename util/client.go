//go:generate mockgen -destination=../mocks/mock_client.go -package=mocks github.com/lucactt/thelist/util Client

package util

import (
	"github.com/lucactt/thelist/data"
	"github.com/ryanbradynd05/go-tmdb"
)

// Client allows to retrieve info about shows.
type Client interface {
	SearchShow(name string) (*data.ShowList, error)
}

// APIClient allows to communicate with the TMDb API.
// A new one should be created by using the NewClient method.
type APIClient struct {
	client *tmdb.TMDb
}

// convertToShowList wraps the results of the TMDb API library into the
// ShowSearchResult struct.
func convertToShowList(result *tmdb.MultiSearchResults) *data.ShowList {
	movies := result.GetMoviesResults()
	tv := result.GetTvResults()
	shows := make([]*data.Show, result.TotalResults)

	for i := 0; i < len(movies); i++ {
		movie := movies[i]
		shows[i] = &data.Show{
			ID:          movie.ID,
			Name:        movie.Title,
			ReleaseDate: movie.ReleaseDate,
			Popularity:  movie.Popularity,
			VoteAverage: movie.VoteAverage,
		}
	}

	for i := 0; i < len(tv); i++ {
		tv := tv[i]
		shows[i+len(movies)] = &data.Show{
			ID:          tv.ID,
			Name:        tv.Name,
			ReleaseDate: tv.FirstAirDate,
			Popularity:  tv.Popularity,
			VoteAverage: tv.VoteAverage,
		}
	}

	return &data.ShowList{
		Results:      shows,
		TotalResults: result.TotalResults,
	}
}

// NewAPIClient creates a new api client using the given API authentication key.
func NewAPIClient(apiKey string) *APIClient {
	return &APIClient{tmdb.Init(tmdb.Config{
		APIKey: apiKey,
	})}
}

// SearchShow searches for the given show name in both movies an tv series.
func (c *APIClient) SearchShow(name string) (*data.ShowList, error) {
	result, err := c.client.SearchMulti(name, nil)
	if err != nil {
		return nil, err
	}

	return convertToShowList(result), nil
}
