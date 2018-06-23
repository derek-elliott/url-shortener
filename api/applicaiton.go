package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/derek-elliott/url-shortener/cache"
	"github.com/derek-elliott/url-shortener/db"
	"github.com/gorilla/mux"
	// Blank import for Postgres support
	_ "github.com/jinzhu/gorm/dialects/postgres"
	log "github.com/sirupsen/logrus"
)

// App holds the router, db and cache connections
type App struct {
	Router   *mux.Router
	DB       db.Store
	Cache    cache.Cache
	Hostname string
}

// Route holds all the information about a route registered with our service.
type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

// Routes holds a list of Routes
type Routes []Route

// RegisterPayload represents a payload to register a URL with our shortener.
type RegisterPayload struct {
	URL string `json:"url"`
	TTL string `json:"ttl"`
}

// InitRouter initializes the router
func (a *App) InitRouter() {
	routes := Routes{
		Route{
			"Register",
			"POST",
			"/",
			a.RegisterShortener,
		},
		Route{
			"Redirect",
			"GET",
			"/{token}",
			a.RedirectToURL,
		},
		Route{
			"Stats",
			"GET",
			"/stats",
			a.GetStats,
		},
		Route{
			"URLStats",
			"GET",
			"/stats/{token}",
			a.GetURLStats,
		},
		Route{
			"DeleteAll",
			"DELETE",
			"/",
			a.DeleteAll,
		},
		Route{
			"DeleteURL",
			"DELETE",
			"/{token}",
			a.DeleteURL,
		},
	}

	a.Router = mux.NewRouter().StrictSlash(true)
	for _, route := range routes {
		a.Router.Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(route.HandlerFunc)
	}
	a.Router.Use(Logger)
}

// Run runs the application
func (a *App) Run(port int) error {
	a.InitRouter()
	bindAddress := fmt.Sprintf(":%d", port)
	err := http.ListenAndServe(bindAddress, a.Router)
	return err
}

// RegisterShortener registeres a shortened url with the service
func (a *App) RegisterShortener(w http.ResponseWriter, r *http.Request) {
	var payload RegisterPayload
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.WithError(err).Error("Unable to read request body")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer r.Body.Close()
	if err = json.Unmarshal(body, &payload); err != nil {
		log.WithError(err).Error("Unable to deserialize request body")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	duration, err := time.ParseDuration(payload.TTL)
	if err != nil {
		log.WithError(err).Error("Unable to parse TTL from request body")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	shortURL := db.ShortURL{}
	shortURL.URL = payload.URL
	shortURL.Expiration = time.Now().Add(duration).String()

	if err = a.DB.CreateShortURL(&shortURL); err != nil {
		log.WithError(err).Error("Database Error")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	shortURL.Token = ConvertIDToToken(int(shortURL.ID))
	shortURL.ShortenedURL = fmt.Sprintf("%s/%s", a.Hostname, shortURL.Token)

	if err = a.DB.UpdateShortURL(&shortURL); err != nil {
		log.WithError(err).Error("Database Error")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err = a.Cache.SetURL(shortURL.Token, shortURL.ShortenedURL, duration); err != nil {
		log.WithError(err).Error("Cache Error")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	response := db.URLData{URL: shortURL.URL, Token: shortURL.Token, ShortenedURL: shortURL.ShortenedURL, Expiration: shortURL.Expiration}
	if err = json.NewEncoder(w).Encode(response); err != nil {
		log.WithError(err).Error("Unable to serialize response")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

// RedirectToURL redirects a request to the specified URL
func (a *App) RedirectToURL(w http.ResponseWriter, r *http.Request) {

}

// GetStats retrieves all stats for the service
func (a *App) GetStats(w http.ResponseWriter, r *http.Request) {

}

// GetURLStats retrieves stats for the specified shortener
func (a *App) GetURLStats(w http.ResponseWriter, r *http.Request) {

}

// DeleteAll removes all shorteners from the service
func (a *App) DeleteAll(w http.ResponseWriter, r *http.Request) {

}

// DeleteURL removes the specified shortener form the service
func (a *App) DeleteURL(w http.ResponseWriter, r *http.Request) {

}
