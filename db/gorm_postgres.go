package db

import (
	"fmt"

	"github.com/jinzhu/gorm"
	// Blank import for postgres support
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

// GormStore implements Store for Gorm Postgres
type GormStore struct {
	client *gorm.DB
}

// InitDB initializes the database
func (s *GormStore) InitDB(user, pass, name, host string, port int) error {
	connStr := fmt.Sprintf("host=%s port=%d user=%s dbname=%s sslmode=disable password=%s", host, port, user, name, pass)
	db, err := gorm.Open("postgres", connStr)
	if err != nil {
		return err
	}
	defer db.Close()
	db.AutoMigrate(&ShortURL{})
	s.client = db
	return nil
}

func (s *GormStore) getShortURL(token string) (*ShortURL, error) {
	shortURL := ShortURL{}
	if err := s.client.Where("token = ?", token).First(&shortURL).Error; err != nil {
		return &shortURL, err
	}
	return &shortURL, nil
}

func (s *GormStore) createShortURL(shortURL *ShortURL) error {
	if err := s.client.Create(shortURL).Error; err != nil {
		return err
	}
	return nil
}

func (s *GormStore) updateShortURL(shortURL *ShortURL) error {
	if err := s.client.Model(shortURL).Updates(shortURL).Error; err != nil {
		return err
	}
	return nil
}

func (s *GormStore) deleteShortURL(shortURL *ShortURL) error {
	if err := s.client.Delete(shortURL).Error; err != nil {
		return err
	}
	return nil
}

// CollectStats collects the overall stats of the service
func (s *GormStore) CollectStats() (*Stats, error) {
	stats := Stats{}
	if err := s.client.Table("short_url").Count(&stats.TotalUrls).Error; err != nil {
		return &stats, err
	}
	allUrls := ShortURLS{}
	if err := s.client.Find(&allUrls).Error; err != nil {
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
