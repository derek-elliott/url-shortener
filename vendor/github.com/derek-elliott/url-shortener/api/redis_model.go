package api

import (
	"time"

	"github.com/go-redis/redis"
)

// Shortener holds the token, url map entry
type Shortener struct {
	Token string
	URL   string
}

func (s *Shortener) setURL(client *redis.Client, TTL time.Duration) error {
	if err := client.Set(s.Token, s.URL, TTL).Err(); err != nil {
		return err
	}
	return nil
}

func (s *Shortener) getURL(client *redis.Client) error {
	val, err := client.Get(s.Token).Result()
	if err != nil {
		return err
	}
	s.URL = val
	return nil
}
