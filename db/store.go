package db

// Store represents a generic database store for URL shorteners
type Store interface {
	InitDB(user, pass, name, host string, port int) error
	GetShortURL(token string) (*ShortURL, error)
	CreateShortURL(shortURL *ShortURL) error
	UpdateShortURL(shortURL *ShortURL) error
	DeleteShortURL(shortURL *ShortURL) error
	CollectStats() (*Stats, error)
}

// Stats holds the overall stats for the service
type Stats struct {
	TotalUrls       int `json:"total_urls"`
	TotalRedirects  int `json:"total_redirects"`
	AverageResponse int `json:"average_response"`
}

// URLData represents a URL
type URLData struct {
	URL          string `json:"url"`
	Token        string `json:"token"`
	ShortenedURL string `json:"shortened_url"`
	Expiration   string `json:"expiration"`
}

// URLStats hold the stats for a given URL
type URLStats struct {
	Redirects       int `json:"redirects"`
	AverageResponse int `json:"average_response"`
}

// ShortURL represents the shortened url and all related metadata
type ShortURL struct {
	ID int
	URLData
	URLStats
}

// ShortURLS represents multiple ShortURL
type ShortURLS []ShortURL
