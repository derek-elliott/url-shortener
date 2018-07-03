package db

// Store represents a generic database store for URL shorteners
type Store interface {
	InitDB(user, pass, name, host string, port int) error
	GetShortURL(token string) (*ShortURL, error)
	GetAllURLTokens() ([]string, error)
	CreateShortURL(shortURL *ShortURL) error
	UpdateShortURL(shortURL *ShortURL) error
	DeleteShortURL(token string) error
	CollectStats() (*Stats, error)
}

// Stats holds the overall stats for the service
type Stats struct {
	TotalURLs      int `json:"total_urls"`
	TotalRedirects int `json:"total_redirects"`
}

// ShortURL represents the shortened url and all related metadata
type ShortURL struct {
	ID           uint   `json:"-"`
	URL          string `json:"url"`
	Token        string `json:"token"`
	ShortenedURL string `json:"shortened_url"`
	Expiration   string `json:"expiration"`
	Redirects    int    `json:"redirects"`
}

// ShortURLS represents multiple ShortURL
type ShortURLS []ShortURL
