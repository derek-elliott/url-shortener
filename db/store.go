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
	TotalUrls      int `json:"total_urls"`
	TotalRedirects int `json:"total_redirects"`
}

// ShortURL represents the shortened url and all related metadata
type ShortURL struct {
	URL          string `json:"url"`
	Token        string `json:"token"`
	ShortenedURL string `json:"shortened_url"`
	Expiration   string `json:"expiration"`
	Redirects    int    `json:"redirects"`
}

// ShortURLS represents multiple ShortURL
type ShortURLS []ShortURL

// TestStore is for testing purposes
type TestStore struct {
	ShortURL *ShortURL
	Stats    *Stats
	Error    error
}

// InitDB is for testing purposes
func (t TestStore) InitDB(user, pass, name, host string, port int) error {
	return t.Error
}

// GetShortURL is for testing purposes
func (t TestStore) GetShortURL(token string) (*ShortURL, error) {
	return t.ShortURL, t.Error
}

// CreateShortURL is for testing purposes
func (t TestStore) CreateShortURL(shortURL *ShortURL) error {
	t.ShortURL = shortURL
	return t.Error
}

// UpdateShortURL is for testing purposes
func (t TestStore) UpdateShortURL(shortURL *ShortURL) error {
	t.ShortURL = shortURL
	return t.Error
}

// DeleteShortURL is for testing purposes
func (t TestStore) DeleteShortURL(shortURL *ShortURL) error {
	return t.Error
}

// CollectStats is for testing purposes
func (t TestStore) CollectStats() (*Stats, error) {
	return t.Stats, t.Error
}
