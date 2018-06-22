package cache

import (
	"fmt"
	"time"

	"github.com/go-redis/redis"
)

// RedisCache implements Cache for Redis
type RedisCache struct {
	client *redis.Client
}

// InitCache initializes the Redis client
func (c *RedisCache) InitCache(pass, host string, port int) error {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", host, port),
		Password: pass,
		DB:       0,
	})
	if _, err := client.Ping().Result(); err != nil {
		return err
	}
	c.client = client
	return nil
}

// SetURL sets the URL for a given token in Redis
func (c *RedisCache) SetURL(token, url string, ttl time.Duration) error {
	if err := c.client.Set(token, url, ttl).Err(); err != nil {
		return err
	}
	return nil
}

// GetURL gets the URL for the given token from Redis
func (c *RedisCache) GetURL(token string) (*Shortener, error) {
	url, err := c.client.Get(token).Result()
	if err != nil {
		return nil, err
	}
	shortener := Shortener{token, url}
	return &shortener, nil
}
