package cache

import (
	"time"
)

// Cache defines a generic in-memory cache for holding shortened URLs
type Cache interface {
	InitCache(pass, host string, port int) error
	SetURL(token, url string, ttl time.Duration) error
	GetURL(token string) (*Shortener, error)
}

// Shortener holds the token, url map entry
type Shortener struct {
	Token string
	URL   string
}
