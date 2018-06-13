package api

import (
	"github.com/jinzhu/gorm"
	// Blank import for postgres support
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

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
	gorm.Model
	URLData
	URLStats
}

// ShortURLS represents multiple ShortURL
type ShortURLS []ShortURL

func (s *ShortURL) getShortURL(db *gorm.DB) error {
	if err := db.Where(s).First(&s).Error; err != nil {
		return err
	}
	return nil
}

func (s *ShortURL) createShortURL(db *gorm.DB) error {
	if err := db.Create(s).Error; err != nil {
		return err
	}
	return nil
}

func (s *ShortURL) updateShortURL(db *gorm.DB, shortURL ShortURL) error {
	if err := db.Model(s).Updates(shortURL).Error; err != nil {
		return err
	}
	return nil
}

func (s *ShortURL) deleteShortURL(db *gorm.DB) error {
	if err := db.Delete(s).Error; err != nil {
		return err
	}
	return nil
}

// CollectStats collects the overall stats of the service
func CollectStats(db *gorm.DB) (*Stats, error) {
	stats := Stats{}
	if err := db.Table("short_url").Count(&stats.TotalUrls).Error; err != nil {
		return &stats, err
	}
	allUrls := ShortURLS{}
	if err := db.Find(&allUrls).Error; err != nil {
		return &stats, err
	}
	stats.TotalRedirects = 0
	stats.AverageResponse = 0
	for _, url := range allUrls {
		stats.TotalRedirects += url.Redirects
	}
	for _, url := range allUrls {
		stats.AverageResponse += (url.Redirects / stats.TotalRedirects) * url.AverageResponse
	}
	return &stats, nil
}
