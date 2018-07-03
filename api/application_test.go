package api

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/derek-elliott/url-shortener/cache"
	"github.com/derek-elliott/url-shortener/db"
	"github.com/stretchr/testify/assert"
)

type TestStore struct {
	ShortURL  *db.ShortURL
	Stats     *db.Stats
	AllTokens []string
	Error     error
}

func (t *TestStore) InitDB(user, pass, name, host string, port int) error {
	return t.Error
}

func (t *TestStore) GetShortURL(token string) (*db.ShortURL, error) {
	return t.ShortURL, t.Error
}

func (t *TestStore) GetAllURLTokens() ([]string, error) {
	return t.AllTokens, t.Error
}

func (t *TestStore) CreateShortURL(shortURL *db.ShortURL) error {
	t.ShortURL = shortURL
	return t.Error
}

func (t *TestStore) UpdateShortURL(shortURL *db.ShortURL) error {
	t.ShortURL = shortURL
	return t.Error
}

func (t *TestStore) DeleteShortURL(token string) error {
	return t.Error
}

func (t *TestStore) CollectStats() (*db.Stats, error) {
	return t.Stats, t.Error
}

type TestCache struct {
	Shortener *cache.Shortener
	Error     error
}

func (t *TestCache) InitCache(pass, host string, port int) error {
	return t.Error
}

func (t *TestCache) SetURL(token, url string, ttl time.Duration) error {
	t.Shortener = &cache.Shortener{
		Token: token,
		URL:   url,
	}
	return t.Error
}

func (t *TestCache) GetURL(token string) (*cache.Shortener, error) {
	return t.Shortener, t.Error
}

func (t *TestCache) DeleteURL(token string) error {
	return t.Error
}

type Test struct {
	description    string
	dbClient       *TestStore
	cacheClient    *TestCache
	url            string
	payload        string
	expectedStatus int
}

type Tests []Test

func TestRegisterShortener(t *testing.T) {
	assert := assert.New(t)

	tests := Tests{
		{
			description: "successful request to RegisterShortener",
			dbClient: &TestStore{
				ShortURL: &db.ShortURL{},
				Stats:    &db.Stats{},
				Error:    nil,
			},
			cacheClient: &TestCache{
				Shortener: &cache.Shortener{},
				Error:     nil,
			},
			url:            "/",
			payload:        "{\"url\": \"http://www.example.com\", \"ttl\": \"10m\"}",
			expectedStatus: http.StatusCreated,
		},
		{
			description: "bad request for RegisterShortener",
			dbClient: &TestStore{
				ShortURL: &db.ShortURL{},
				Stats:    &db.Stats{},
				Error:    nil,
			},
			cacheClient: &TestCache{
				Shortener: &cache.Shortener{},
				Error:     nil,
			},
			url:            "/",
			payload:        "Test request, please ignore",
			expectedStatus: http.StatusInternalServerError,
		},
		{
			description: "bad ttl on RegisterShortener",
			dbClient: &TestStore{
				ShortURL: &db.ShortURL{},
				Stats:    &db.Stats{},
				Error:    nil,
			},
			cacheClient: &TestCache{
				Shortener: &cache.Shortener{},
				Error:     nil,
			},
			url:            "/",
			payload:        "{\"url\": \"http://www.example.com\", \"ttl\": \"Boo to you\"}",
			expectedStatus: http.StatusBadRequest,
		}, {
			description: "bad URL on RegisterShortener",
			dbClient: &TestStore{
				ShortURL: &db.ShortURL{},
				Stats:    &db.Stats{},
				Error:    nil,
			},
			cacheClient: &TestCache{
				Shortener: &cache.Shortener{},
				Error:     nil,
			},
			url:            "/",
			payload:        "{\"url\": \"www.example.com\", \"ttl\": \"Boo to you\"}",
			expectedStatus: http.StatusBadRequest,
		},
		{
			description: "database error",
			dbClient: &TestStore{
				ShortURL: &db.ShortURL{},
				Stats:    &db.Stats{},
				Error:    errors.New("something is wrong with your database"),
			},
			cacheClient: &TestCache{
				Shortener: &cache.Shortener{},
				Error:     nil,
			},
			url:            "/",
			payload:        "{\"url\": \"http://www.example.com\", \"ttl\": \"10m\"}",
			expectedStatus: http.StatusInternalServerError,
		},
		{
			description: "cache error",
			dbClient: &TestStore{
				ShortURL: &db.ShortURL{},
				Stats:    &db.Stats{},
				Error:    nil,
			},
			cacheClient: &TestCache{
				Shortener: &cache.Shortener{},
				Error:     errors.New("something is wrong with your cache"),
			},
			url:            "/",
			payload:        "{\"url\": \"http://www.example.com\", \"ttl\": \"10m\"}",
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, test := range tests {
		app := &App{
			DB:       test.dbClient,
			Cache:    test.cacheClient,
			Hostname: "test.com",
		}

		request, err := http.NewRequest("POST", test.url, strings.NewReader(test.payload))
		assert.NoError(err)

		w := httptest.NewRecorder()
		app.RegisterShortener(w, request)

		assert.Equal(test.expectedStatus, w.Code, test.description)
	}
}

func TestRedirectToURL(t *testing.T) {
	assert := assert.New(t)

	tests := Tests{
		{
			description: "successful request to RedirectToURL",
			dbClient: &TestStore{
				ShortURL: &db.ShortURL{
					URL:          "www.example.com",
					Token:        "testurl",
					ShortenedURL: "test.com/testurl",
					Expiration:   "10m",
					Redirects:    0,
				},
				Stats: &db.Stats{},
				Error: nil,
			},
			cacheClient: &TestCache{
				Shortener: &cache.Shortener{
					Token: "testurl",
					URL:   "www.example.com",
				},
				Error: nil,
			},
			url:            "/testurl",
			payload:        "",
			expectedStatus: http.StatusFound,
		},
		{
			description: "cache error in RedirectToURL",
			dbClient: &TestStore{
				ShortURL: &db.ShortURL{
					URL:          "www.example.com",
					Token:        "testurl",
					ShortenedURL: "test.com/testurl",
					Expiration:   "10m",
					Redirects:    0,
				},
				Stats: &db.Stats{},
				Error: nil,
			},
			cacheClient: &TestCache{
				Shortener: &cache.Shortener{},
				Error:     errors.New("something's wrong with your cache"),
			},
			url:            "/testurl",
			payload:        "",
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, test := range tests {
		app := &App{
			DB:       test.dbClient,
			Cache:    test.cacheClient,
			Hostname: "test.com",
		}

		request, err := http.NewRequest("GET", test.url, nil)
		assert.NoError(err)

		w := httptest.NewRecorder()
		app.RedirectToURL(w, request)

		assert.Equal(test.expectedStatus, w.Code, test.description)
	}
}

func TestGetStats(t *testing.T) {
	assert := assert.New(t)

	tests := Tests{
		{
			description: "successful request to GetStats",
			dbClient: &TestStore{
				ShortURL: &db.ShortURL{},
				Stats: &db.Stats{
					TotalURLs:      10,
					TotalRedirects: 100,
				},
				Error: nil,
			},
			cacheClient: &TestCache{
				Shortener: &cache.Shortener{},
				Error:     nil,
			},
			url:            "/stats",
			payload:        "",
			expectedStatus: http.StatusOK,
		},
		{
			description: "database error in GetStats",
			dbClient: &TestStore{
				ShortURL: &db.ShortURL{},
				Stats:    &db.Stats{},
				Error:    errors.New("something went wrong"),
			},
			cacheClient: &TestCache{
				Shortener: &cache.Shortener{},
				Error:     nil,
			},
			url:            "/stats",
			payload:        "",
			expectedStatus: http.StatusNotFound,
		},
	}

	for _, test := range tests {
		app := &App{
			DB:       test.dbClient,
			Cache:    test.cacheClient,
			Hostname: "test.com",
		}

		request, err := http.NewRequest("GET", test.url, nil)
		assert.NoError(err)

		w := httptest.NewRecorder()
		app.GetStats(w, request)

		assert.Equal(test.expectedStatus, w.Code, test.description)
	}
}

func TestGetURLStats(t *testing.T) {
	assert := assert.New(t)

	tests := Tests{
		{
			description: "successful request to GetURLStats",
			dbClient: &TestStore{
				ShortURL: &db.ShortURL{},
				Stats: &db.Stats{
					TotalURLs:      10,
					TotalRedirects: 100,
				},
				Error: nil,
			},
			cacheClient: &TestCache{
				Shortener: &cache.Shortener{},
				Error:     nil,
			},
			url:            "/stats/testurl",
			payload:        "",
			expectedStatus: http.StatusOK,
		},
		{
			description: "database error in GetURLStats",
			dbClient: &TestStore{
				ShortURL: &db.ShortURL{},
				Stats:    &db.Stats{},
				Error:    errors.New("something went wrong"),
			},
			cacheClient: &TestCache{
				Shortener: &cache.Shortener{},
				Error:     nil,
			},
			url:            "/stats/testurl",
			payload:        "",
			expectedStatus: http.StatusNotFound,
		},
	}

	for _, test := range tests {
		app := &App{
			DB:       test.dbClient,
			Cache:    test.cacheClient,
			Hostname: "test.com",
		}

		request, err := http.NewRequest("GET", test.url, nil)
		assert.NoError(err)

		w := httptest.NewRecorder()
		app.GetURLStats(w, request)

		assert.Equal(test.expectedStatus, w.Code, test.description)
	}
}

func TestAllDelete(t *testing.T) {
	assert := assert.New(t)

	tests := Tests{
		{
			description: "successful request to DeleteAll",
			dbClient: &TestStore{
				ShortURL: &db.ShortURL{},
				Stats:    &db.Stats{},
				Error:    nil,
			},
			cacheClient: &TestCache{
				Shortener: &cache.Shortener{},
				Error:     nil,
			},
			url:            "/",
			payload:        "",
			expectedStatus: http.StatusNoContent,
		},
		{
			description: "database error in DeleteURL",
			dbClient: &TestStore{
				ShortURL: &db.ShortURL{},
				Stats:    &db.Stats{},
				Error:    errors.New("something went wrong"),
			},
			cacheClient: &TestCache{
				Shortener: &cache.Shortener{},
				Error:     nil,
			},
			url:            "/",
			payload:        "",
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, test := range tests {
		app := &App{
			DB:       test.dbClient,
			Cache:    test.cacheClient,
			Hostname: "test.com",
		}

		request, err := http.NewRequest("DELETE", test.url, nil)
		assert.NoError(err)

		w := httptest.NewRecorder()
		app.DeleteAll(w, request)

		assert.Equal(test.expectedStatus, w.Code, test.description)
	}
}

func TestDeleteURL(t *testing.T) {
	assert := assert.New(t)

	tests := Tests{
		{
			description: "successful request to DeleteURL",
			dbClient: &TestStore{
				ShortURL: &db.ShortURL{},
				Stats:    &db.Stats{},
				Error:    nil,
			},
			cacheClient: &TestCache{
				Shortener: &cache.Shortener{},
				Error:     nil,
			},
			url:            "/testurl",
			payload:        "",
			expectedStatus: http.StatusNoContent,
		},
		{
			description: "database error in DeleteURL",
			dbClient: &TestStore{
				ShortURL: &db.ShortURL{},
				Stats:    &db.Stats{},
				Error:    errors.New("something went wrong"),
			},
			cacheClient: &TestCache{
				Shortener: &cache.Shortener{},
				Error:     nil,
			},
			url:            "/testurl",
			payload:        "",
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, test := range tests {
		app := &App{
			DB:       test.dbClient,
			Cache:    test.cacheClient,
			Hostname: "test.com",
		}

		request, err := http.NewRequest("DELETE", test.url, nil)
		assert.NoError(err)

		w := httptest.NewRecorder()
		app.DeleteURL(w, request)

		assert.Equal(test.expectedStatus, w.Code, test.description)
	}
}
