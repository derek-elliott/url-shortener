package api

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/derek-elliott/url-shortener/cache"
	"github.com/derek-elliott/url-shortener/db"
	"github.com/derek-elliott/url-shortener/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestSuccessfulRegisterShortener(t *testing.T) {
	assert := assert.New(t)

	testDB := &mocks.Store{}
	testDB.On("CreateShortURL", mock.Anything).Return(nil)

	testCache := &mocks.Cache{}
	testCache.On("SetURL", mock.Anything, mock.Anything, mock.Anything).Return(nil)

	app := &App{
		DB:       testDB,
		Cache:    testCache,
		Hostname: "test.com",
	}

	payload := "{\"url\": \"http://www.example.com\", \"ttl\": \"10m\"}"

	request, err := http.NewRequest("POST", "/", strings.NewReader(payload))
	assert.NoError(err)

	w := httptest.NewRecorder()
	app.RegisterShortener(w, request)

	testDB.AssertExpectations(t)
	testCache.AssertExpectations(t)
	assert.Equal(http.StatusCreated, w.Code, "successful request")

}

func TestBadPayloadRegiterShortener(t *testing.T) {
	assert := assert.New(t)

	testDB := &mocks.Store{}
	testDB.On("CreateShortURL", mock.Anything).Return(nil)

	testCache := &mocks.Cache{}
	testCache.On("SetURL", mock.Anything, mock.Anything, mock.Anything).Return(nil)

	app := &App{
		DB:       testDB,
		Cache:    testCache,
		Hostname: "test.com",
	}

	payload := "Test payload, please ignore"

	request, err := http.NewRequest("POST", "/", strings.NewReader(payload))
	assert.NoError(err)

	w := httptest.NewRecorder()
	app.RegisterShortener(w, request)

	assert.Equal(http.StatusInternalServerError, w.Code, "bad payload")

}

type errReader int

func (errReader) Read(p []byte) (n int, err error) {
	return 0, errors.New("test error")
}

func TestBadRequestRegiterShortener(t *testing.T) {

	assert := assert.New(t)

	testDB := &mocks.Store{}
	testDB.On("CreateShortURL", mock.Anything).Return(nil)

	testCache := &mocks.Cache{}
	testCache.On("SetURL", mock.Anything, mock.Anything, mock.Anything).Return(nil)

	app := &App{
		DB:       testDB,
		Cache:    testCache,
		Hostname: "test.com",
	}

	request, err := http.NewRequest("POST", "/", errReader(0))
	assert.NoError(err)

	w := httptest.NewRecorder()
	app.RegisterShortener(w, request)

	assert.Equal(http.StatusBadRequest, w.Code, "Body read error")

}

func TestBadTTLRegisterShortener(t *testing.T) {
	assert := assert.New(t)

	testDB := &mocks.Store{}
	testDB.On("CreateShortURL", mock.Anything).Return(nil)

	testCache := &mocks.Cache{}
	testCache.On("SetURL", mock.Anything, mock.Anything, mock.Anything).Return(nil)

	app := &App{
		DB:       testDB,
		Cache:    testCache,
		Hostname: "test.com",
	}

	payload := "{\"url\": \"http://www.example.com\", \"ttl\": \"Boo\"}"

	request, err := http.NewRequest("POST", "/", strings.NewReader(payload))
	assert.NoError(err)

	w := httptest.NewRecorder()
	app.RegisterShortener(w, request)

	assert.Equal(http.StatusBadRequest, w.Code, "bad ttl")

}

func TestBadURLRegisterShortener(t *testing.T) {
	assert := assert.New(t)

	testDB := &mocks.Store{}
	testDB.On("CreateShortURL", mock.Anything).Return(nil)

	testCache := &mocks.Cache{}
	testCache.On("SetURL", mock.Anything, mock.Anything, mock.Anything).Return(nil)

	app := &App{
		DB:       testDB,
		Cache:    testCache,
		Hostname: "test.com",
	}

	payload := "{\"url\": \"www.example.com\", \"ttl\": \"10m\"}"

	request, err := http.NewRequest("POST", "/", strings.NewReader(payload))
	assert.NoError(err)

	w := httptest.NewRecorder()
	app.RegisterShortener(w, request)

	assert.Equal(http.StatusBadRequest, w.Code, "bad url")

}

func TestDBErrorCreateShortener(t *testing.T) {
	assert := assert.New(t)

	testDB := &mocks.Store{}
	testDB.On("CreateShortURL", mock.Anything).Return(errors.New("test db error"))

	testCache := &mocks.Cache{}
	testCache.On("SetURL", mock.Anything, mock.Anything, mock.Anything).Return(nil)

	app := &App{
		DB:       testDB,
		Cache:    testCache,
		Hostname: "test.com",
	}

	payload := "{\"url\": \"https://www.example.com\", \"ttl\": \"10m\"}"

	request, err := http.NewRequest("POST", "/", strings.NewReader(payload))
	assert.NoError(err)

	w := httptest.NewRecorder()
	app.RegisterShortener(w, request)

	assert.Equal(http.StatusInternalServerError, w.Code, "CreateShortURL error")

}

func TestCacheErrorCreateShortener(t *testing.T) {
	assert := assert.New(t)

	testDB := &mocks.Store{}
	testDB.On("CreateShortURL", mock.Anything).Return(nil)

	testCache := &mocks.Cache{}
	testCache.On("SetURL", mock.AnythingOfType("string"), mock.AnythingOfType("string"), mock.Anything).Return(errors.New("test db error"))

	app := &App{
		DB:       testDB,
		Cache:    testCache,
		Hostname: "test.com",
	}

	payload := "{\"url\": \"https://www.example.com\", \"ttl\": \"10m\"}"

	request, err := http.NewRequest("POST", "/", strings.NewReader(payload))
	assert.NoError(err)

	w := httptest.NewRecorder()
	app.RegisterShortener(w, request)

	assert.Equal(http.StatusInternalServerError, w.Code, "SetURL error")

}

func TestSuccessfulRedirectToURL(t *testing.T) {
	assert := assert.New(t)

	testDB := &mocks.Store{}
	testDB.On("GetShortURL", mock.AnythingOfType("string")).Return(&db.ShortURL{URL: "https://www.example.com", Token: "testurl", ShortenedURL: "test.com/testurl", Expiration: "", Redirects: 0}, nil)
	testDB.On("UpdateShortURL", mock.Anything).Return(nil)
	testCache := &mocks.Cache{}
	testCache.On("GetURL", mock.AnythingOfType("string")).Return(&cache.Shortener{Token: "testurl", URL: "https://www.example.com"}, nil)

	app := &App{
		DB:       testDB,
		Cache:    testCache,
		Hostname: "test.com",
	}

	request, err := http.NewRequest("GET", "/testurl", nil)
	assert.NoError(err)

	w := httptest.NewRecorder()
	app.RedirectToURL(w, request)

	testCache.AssertExpectations(t)
	assert.Equal(http.StatusFound, w.Code, "successful RedirectToURL")

}

func TestBadCacheRedirectToURL(t *testing.T) {
	assert := assert.New(t)

	testDB := &mocks.Store{}
	testDB.On("GetShortURL", mock.AnythingOfType("string")).Return(&db.ShortURL{}, nil)
	testDB.On("UpdateShortURL", mock.Anything).Return(nil)
	testCache := &mocks.Cache{}
	testCache.On("GetURL", mock.AnythingOfType("string")).Return(&cache.Shortener{}, errors.New("test cache error"))

	app := &App{
		DB:       testDB,
		Cache:    testCache,
		Hostname: "test.com",
	}

	request, err := http.NewRequest("GET", "/testurl", nil)
	assert.NoError(err)

	w := httptest.NewRecorder()
	app.RedirectToURL(w, request)

	testCache.AssertExpectations(t)
	assert.Equal(http.StatusInternalServerError, w.Code, "bad cache in RedirectToURL")

}

func TestSuccessfulGetStats(t *testing.T) {
	assert := assert.New(t)

	testDB := &mocks.Store{}
	testDB.On("CollectStats").Return(&db.Stats{TotalURLs: 0, TotalRedirects: 0}, nil)
	testCache := &mocks.Cache{}

	app := &App{
		DB:       testDB,
		Cache:    testCache,
		Hostname: "test.com",
	}

	request, err := http.NewRequest("GET", "/admin/stats", nil)
	assert.NoError(err)

	w := httptest.NewRecorder()
	app.GetStats(w, request)

	testDB.AssertExpectations(t)
	assert.Equal(http.StatusOK, w.Code, "successful GetStats")
}

func TestBadDBGetStats(t *testing.T) {
	assert := assert.New(t)

	testDB := &mocks.Store{}
	testDB.On("CollectStats").Return(&db.Stats{}, errors.New("test db error"))
	testCache := &mocks.Cache{}

	app := &App{
		DB:       testDB,
		Cache:    testCache,
		Hostname: "test.com",
	}

	request, err := http.NewRequest("GET", "/admin/stats", nil)
	assert.NoError(err)

	w := httptest.NewRecorder()
	app.GetStats(w, request)

	testDB.AssertExpectations(t)
	assert.Equal(http.StatusNotFound, w.Code, "db error in GetStats")
}

func TestSuccessfulGetURLStats(t *testing.T) {
	assert := assert.New(t)

	testDB := &mocks.Store{}
	testDB.On("GetShortURL", mock.AnythingOfType("string")).Return(&db.ShortURL{}, nil)
	testCache := &mocks.Cache{}

	app := &App{
		DB:       testDB,
		Cache:    testCache,
		Hostname: "test.com",
	}

	request, err := http.NewRequest("GET", "/admin/stats/testurl", nil)
	assert.NoError(err)

	w := httptest.NewRecorder()
	app.GetURLStats(w, request)

	testDB.AssertExpectations(t)
	assert.Equal(http.StatusOK, w.Code, "successful GetURLStats")
}

func TestDBErrorGetURLStats(t *testing.T) {
	assert := assert.New(t)

	testDB := &mocks.Store{}
	testDB.On("GetShortURL", mock.AnythingOfType("string")).Return(&db.ShortURL{}, errors.New("test db error"))
	testCache := &mocks.Cache{}

	app := &App{
		DB:       testDB,
		Cache:    testCache,
		Hostname: "test.com",
	}

	request, err := http.NewRequest("GET", "/admin/stats/testurl", nil)
	assert.NoError(err)

	w := httptest.NewRecorder()
	app.GetURLStats(w, request)

	testDB.AssertExpectations(t)
	assert.Equal(http.StatusNotFound, w.Code, "db error in GetURLStats")
}

func TestSuccessfulDeleteAll(t *testing.T) {
	assert := assert.New(t)

	testDB := &mocks.Store{}
	testDB.On("GetAllURLTokens").Return([]string{}, nil)
	testDB.On("DeleteShortURL", mock.AnythingOfType("string")).Return(nil)
	testCache := &mocks.Cache{}

	app := &App{
		DB:       testDB,
		Cache:    testCache,
		Hostname: "test.com",
	}

	request, err := http.NewRequest("DELETE", "/", nil)
	assert.NoError(err)

	w := httptest.NewRecorder()
	app.DeleteAll(w, request)

	assert.Equal(http.StatusNoContent, w.Code, "successful DeleteAll")
}

func TestDBGetTokenErrorDeleteAll(t *testing.T) {
	assert := assert.New(t)

	testDB := &mocks.Store{}
	testDB.On("GetAllURLTokens").Return([]string{}, errors.New("test GetAllURLTokens error"))
	testDB.On("DeleteShortURL", mock.AnythingOfType("string")).Return(nil)
	testCache := &mocks.Cache{}

	app := &App{
		DB:       testDB,
		Cache:    testCache,
		Hostname: "test.com",
	}

	request, err := http.NewRequest("DELETE", "/", nil)
	assert.NoError(err)

	w := httptest.NewRecorder()
	app.DeleteAll(w, request)

	assert.Equal(http.StatusInternalServerError, w.Code, "GetAllURLTOkens error in DeleteAll")
}

func TestDBDeleteErrorDeleteAll(t *testing.T) {
	assert := assert.New(t)

	testDB := &mocks.Store{}
	testDB.On("GetAllURLTokens").Return([]string{"testurl"}, nil)
	testDB.On("DeleteShortURL", mock.AnythingOfType("string")).Return(errors.New("test DeleteShortURL error"))
	testCache := &mocks.Cache{}

	app := &App{
		DB:       testDB,
		Cache:    testCache,
		Hostname: "test.com",
	}

	request, err := http.NewRequest("DELETE", "/", nil)
	assert.NoError(err)

	w := httptest.NewRecorder()
	app.DeleteAll(w, request)

	assert.Equal(http.StatusInternalServerError, w.Code, "GetAllURLTOkens error in DeleteAll")
}

func TestSuccessfulDeleteURL(t *testing.T) {
	assert := assert.New(t)

	testDB := &mocks.Store{}
	testDB.On("DeleteShortURL", mock.AnythingOfType("string")).Return(nil)
	testCache := &mocks.Cache{}

	app := &App{
		DB:       testDB,
		Cache:    testCache,
		Hostname: "test.com",
	}

	request, err := http.NewRequest("DELETE", "/testurl", nil)
	assert.NoError(err)

	w := httptest.NewRecorder()
	app.DeleteURL(w, request)

	assert.Equal(http.StatusNoContent, w.Code, "successful DeleteURL")
}

func TestDBErrorDeleteURL(t *testing.T) {
	assert := assert.New(t)

	testDB := &mocks.Store{}
	testDB.On("DeleteShortURL", mock.AnythingOfType("string")).Return(errors.New("test DB error"))
	testCache := &mocks.Cache{}

	app := &App{
		DB:       testDB,
		Cache:    testCache,
		Hostname: "test.com",
	}

	request, err := http.NewRequest("DELETE", "/testurl", nil)
	assert.NoError(err)

	w := httptest.NewRecorder()
	app.DeleteURL(w, request)

	assert.Equal(http.StatusInternalServerError, w.Code, "db error in DeleteURL")
}

func TestIncrementRedirects(t *testing.T) {
	testDB := &mocks.Store{}
	testDB.On("GetShortURL", "testurl").Return(
		&db.ShortURL{
			URL:          "https://www.example.com",
			Token:        "testurl",
			ShortenedURL: "test.com/testurl",
			Expiration:   "",
			Redirects:    0},
		nil)
	testDB.On("UpdateShortURL",
		&db.ShortURL{
			URL:          "https://www.example.com",
			Token:        "testurl",
			ShortenedURL: "test.com/testurl",
			Expiration:   "",
			Redirects:    1}).Return(nil)
	testCache := &mocks.Cache{}

	app := &App{
		DB:       testDB,
		Cache:    testCache,
		Hostname: "test.com",
	}

	app.incrementRedirects("testurl")

	testDB.AssertExpectations(t)
	testDB.AssertCalled(t, "UpdateShortURL",
		&db.ShortURL{
			URL:          "https://www.example.com",
			Token:        "testurl",
			ShortenedURL: "test.com/testurl",
			Expiration:   "",
			Redirects:    1})
}

func TestDBErrorIncrementRedirects(t *testing.T) {
	testDB := &mocks.Store{}
	testDB.On("GetShortURL", "testurl").Return(&db.ShortURL{}, errors.New("test db error"))
	testDB.On("UpdateShortURL", &db.ShortURL{}).Return(nil)
	testCache := &mocks.Cache{}

	app := &App{
		DB:       testDB,
		Cache:    testCache,
		Hostname: "test.com",
	}

	app.incrementRedirects("testurl")

	testDB.AssertCalled(t, "GetShortURL", "testurl")
	testDB.AssertNotCalled(t, "UpdateShortURL")
}

func TestDBErrorUpdateIncrementRedirects(t *testing.T) {
	testDB := &mocks.Store{}
	testDB.On("GetShortURL", "testurl").Return(
		&db.ShortURL{
			URL:          "https://www.example.com",
			Token:        "testurl",
			ShortenedURL: "test.com/testurl",
			Expiration:   "",
			Redirects:    0},
		nil)
	testDB.On("UpdateShortURL",
		&db.ShortURL{
			URL:          "https://www.example.com",
			Token:        "testurl",
			ShortenedURL: "test.com/testurl",
			Expiration:   "",
			Redirects:    1}).Return(errors.New("test db error"))
	testCache := &mocks.Cache{}

	app := &App{
		DB:       testDB,
		Cache:    testCache,
		Hostname: "test.com",
	}

	app.incrementRedirects("testurl")

	testDB.AssertCalled(t, "GetShortURL", "testurl")
	testDB.AssertCalled(t, "UpdateShortURL",
		&db.ShortURL{
			URL:          "https://www.example.com",
			Token:        "testurl",
			ShortenedURL: "test.com/testurl",
			Expiration:   "",
			Redirects:    1})
}

func TestCleanExpiredRecords(t *testing.T) {
	testDB := &mocks.Store{}
	testDB.On("GetAllURLTokens").Return([]string{"testurl"}, nil)
	testDB.On("GetShortURL", "testurl").Return(
		&db.ShortURL{
			URL:          "https://www.example.com",
			Token:        "testurl",
			ShortenedURL: "test.com/testurl",
			Expiration:   "2018-07-03T11:10:33-04:00",
			Redirects:    0}, nil)
	testDB.On("DeleteShortURL", "testurl").Return(nil)
	testCache := &mocks.Cache{}

	app := &App{
		DB:       testDB,
		Cache:    testCache,
		Hostname: "test.com",
	}

	app.cleanExpiredRecords()

	testDB.AssertExpectations(t)
}

func TestNoURLSCleanExpiredRecords(t *testing.T) {
	testDB := &mocks.Store{}
	testDB.On("GetAllURLTokens").Return([]string{"testurl"}, errors.New("test db error"))
	testDB.On("GetShortURL", "testurl").Return(&db.ShortURL{}, nil)
	testDB.On("DeleteShortURL", "testurl").Return(nil)
	testCache := &mocks.Cache{}

	app := &App{
		DB:       testDB,
		Cache:    testCache,
		Hostname: "test.com",
	}

	app.cleanExpiredRecords()

	testDB.AssertCalled(t, "GetAllURLTokens")
}

func TestDBErrorCleanExpiredRecords(t *testing.T) {
	testDB := &mocks.Store{}
	testDB.On("GetAllURLTokens").Return([]string{"testurl"}, nil)
	testDB.On("GetShortURL", "testurl").Return(&db.ShortURL{}, errors.New("test db error"))
	testDB.On("DeleteShortURL", "testurl").Return(nil)
	testCache := &mocks.Cache{}

	app := &App{
		DB:       testDB,
		Cache:    testCache,
		Hostname: "test.com",
	}

	app.cleanExpiredRecords()

	testDB.AssertCalled(t, "GetAllURLTokens")
	testDB.AssertCalled(t, "GetShortURL", "testurl")
}
