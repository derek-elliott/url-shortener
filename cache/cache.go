package cache

import (
	"time"
)

// Cache defines a generic remote cache for holding shortened URLs
type Cache interface {
	InitCache(pass, host string, port int) error
	SetURL(token, url string, ttl time.Duration) error
	GetURL(token string) (*Shortener, error)
	DeleteURL(token string) error
}

// Shortener holds the token, url map entry
type Shortener struct {
	Token string
	URL   string
}

// TestCache is for testing purposes
type TestCache struct {
	Shortener *Shortener
	Error     error
}

// InitCache is for testing purposes
func (t TestCache) InitCache(pass, host string, port int) error {
	return t.Error
}

//SetURL is for testing purposes
func (t TestCache) SetURL(token, url string, ttl time.Duration) error {
	t.Shortener = &Shortener{
		Token: token,
		URL:   url,
	}
	return t.Error
}

// GetURL is for testing purposes
func (t TestCache) GetURL(token string) (*Shortener, error) {
	return t.Shortener, t.Error
}

// DeleteURL is for testing purposes
func (t TestCache) DeleteURL(token string) error {
	return t.Error
}
