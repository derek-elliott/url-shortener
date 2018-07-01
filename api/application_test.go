package api

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/derek-elliott/url-shortener/cache"
	"github.com/derek-elliott/url-shortener/db"
	"github.com/stretchr/testify/assert"
)

type Test struct {
	description       string
	dbClient          *db.TestStore
	cacheClient       *cache.TestCache
	url               string
	payload           string
	expectedStatus    int
	expectedRedirects int
}

type Tests []Test

func TestRegisterShortener(t *testing.T) {
	assert := assert.New(t)

	tests := Tests{
		{
			description: "successful request to RegisterShortener",
			dbClient: &db.TestStore{
				ShortURL: &db.ShortURL{},
				Stats:    &db.Stats{},
				Error:    nil,
			},
			cacheClient: &cache.TestCache{
				Shortener: &cache.Shortener{},
				Error:     nil,
			},
			url:            "/",
			payload:        "{\"url\": \"http://www.example.com\", \"ttl\": \"10m\"}",
			expectedStatus: http.StatusCreated,
		},
		{
			description: "bad request for RegisterShortener",
			dbClient: &db.TestStore{
				ShortURL: &db.ShortURL{},
				Stats:    &db.Stats{},
				Error:    nil,
			},
			cacheClient: &cache.TestCache{
				Shortener: &cache.Shortener{},
				Error:     nil,
			},
			url:            "/",
			payload:        "Test request, please ignore",
			expectedStatus: http.StatusInternalServerError,
		},
		{
			description: "bad ttl on RegisterShortener",
			dbClient: &db.TestStore{
				ShortURL: &db.ShortURL{},
				Stats:    &db.Stats{},
				Error:    nil,
			},
			cacheClient: &cache.TestCache{
				Shortener: &cache.Shortener{},
				Error:     nil,
			},
			url:            "/",
			payload:        "{\"url\": \"http://www.example.com\", \"ttl\": \"Boo to you\"}",
			expectedStatus: http.StatusBadRequest,
		}, {
			description: "bad URL on RegisterShortener",
			dbClient: &db.TestStore{
				ShortURL: &db.ShortURL{},
				Stats:    &db.Stats{},
				Error:    nil,
			},
			cacheClient: &cache.TestCache{
				Shortener: &cache.Shortener{},
				Error:     nil,
			},
			url:            "/",
			payload:        "{\"url\": \"www.example.com\", \"ttl\": \"Boo to you\"}",
			expectedStatus: http.StatusBadRequest,
		},
		{
			description: "database error",
			dbClient: &db.TestStore{
				ShortURL: &db.ShortURL{},
				Stats:    &db.Stats{},
				Error:    errors.New("something is wrong with your database"),
			},
			cacheClient: &cache.TestCache{
				Shortener: &cache.Shortener{},
				Error:     nil,
			},
			url:            "/",
			payload:        "{\"url\": \"http://www.example.com\", \"ttl\": \"10m\"}",
			expectedStatus: http.StatusInternalServerError,
		},
		{
			description: "cache error",
			dbClient: &db.TestStore{
				ShortURL: &db.ShortURL{},
				Stats:    &db.Stats{},
				Error:    nil,
			},
			cacheClient: &cache.TestCache{
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
			dbClient: &db.TestStore{
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
			cacheClient: &cache.TestCache{
				Shortener: &cache.Shortener{
					Token: "testurl",
					URL:   "www.example.com",
				},
				Error: nil,
			},
			url:               "/testurl",
			payload:           "",
			expectedStatus:    http.StatusFound,
			expectedRedirects: 1,
		},
		{
			description: "cache error in RedirectToURL",
			dbClient: &db.TestStore{
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
			cacheClient: &cache.TestCache{
				Shortener: &cache.Shortener{},
				Error:     errors.New("something's wrong with your cache"),
			},
			url:               "/testurl",
			payload:           "",
			expectedStatus:    http.StatusInternalServerError,
			expectedRedirects: 0,
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
