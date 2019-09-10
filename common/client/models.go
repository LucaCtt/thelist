package client

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
