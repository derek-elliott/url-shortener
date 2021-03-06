package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	"github.com/derek-elliott/url-shortener/cache"
	"github.com/derek-elliott/url-shortener/db"
	"github.com/gorilla/mux"
	// Blank import for Postgres support
	_ "github.com/jinzhu/gorm/dialects/postgres"
	log "github.com/sirupsen/logrus"
)

const tokenLength = 6

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
			"/admin/stats",
			a.GetStats,
		},
		Route{
			"URLStats",
			"GET",
			"/admin/stats/{token}",
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

	a.Router = mux.NewRouter()
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
	go func() {
		for {
			a.cleanExpiredRecords()
			time.Sleep(30 * time.Second)
		}
	}()
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
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	shortURL := db.ShortURL{}
	if _, err = url.ParseRequestURI(payload.URL); err != nil {
		log.WithField("url", payload.URL).WithError(err).Error("Unable to parse URL from request body")
		w.WriteHeader(http.StatusBadRequest)
	}
	shortURL.URL = payload.URL
	shortURL.Expiration = time.Now().Add(duration).Format(time.RFC3339)
	shortURL.Token, err = generateToken(tokenLength)
	if err != nil {
		log.WithError(err).Error("Error generating URL token")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	shortURL.ShortenedURL = fmt.Sprintf("%s/%s", a.Hostname, shortURL.Token)

	if err = a.DB.CreateShortURL(&shortURL); err != nil {
		log.WithField("short_url", shortURL).WithError(err).Error("Database Error")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err = a.Cache.SetURL(shortURL.Token, shortURL.URL, duration); err != nil {
		log.WithFields(log.Fields{"token": shortURL.Token, "url": shortURL.URL, "duration": duration}).WithError(err).Error("Cache Error")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err = json.NewEncoder(w).Encode(shortURL); err != nil {
		log.WithField("response", shortURL).WithError(err).Error("Unable to serialize RegisterShortener response")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

// RedirectToURL redirects a request to the specified URL
func (a *App) RedirectToURL(w http.ResponseWriter, r *http.Request) {
	token := mux.Vars(r)["token"]
	url, err := a.Cache.GetURL(token)
	if err != nil {
		log.WithField("token", token).WithError(err).Error("Unable to obtain URL from cache")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	go a.incrementRedirects(token)
	http.Redirect(w, r, url.URL, http.StatusFound)
	return
}

// GetStats retrieves all stats for the service
func (a *App) GetStats(w http.ResponseWriter, r *http.Request) {
	stats, err := a.DB.CollectStats()
	if err != nil {
		log.WithError(err).Error("Unable to collect stats from database")
		w.WriteHeader(http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err = json.NewEncoder(w).Encode(&stats); err != nil {
		log.WithField("response", stats).WithError(err).Error("Unable to serialize GetStats response")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

// GetURLStats retrieves stats for the specified shortener
func (a *App) GetURLStats(w http.ResponseWriter, r *http.Request) {
	token := mux.Vars(r)["token"]
	shortURL, err := a.DB.GetShortURL(token)
	if err != nil {
		log.WithField("token", token).WithError(err).Error("Unable to retrieve ShortURL from database")
		w.WriteHeader(http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err = json.NewEncoder(w).Encode(&shortURL); err != nil {
		log.WithField("response", shortURL).WithError(err).Error("Unable to seralize GetURLStats response")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

// DeleteAll removes all shorteners from the service
func (a *App) DeleteAll(w http.ResponseWriter, r *http.Request) {
	tokens, err := a.DB.GetAllURLTokens()
	if err != nil {
		log.WithError(err).Error("Unable to get all tokens from database in DeleteAll")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	for _, token := range tokens {
		if err := a.DB.DeleteShortURL(token); err != nil {
			log.WithError(err).WithField("token", token).Error("Unable to delete from database")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
	w.WriteHeader(http.StatusNoContent)
	return
}

// DeleteURL removes the specified shortener form the service
func (a *App) DeleteURL(w http.ResponseWriter, r *http.Request) {
	token := mux.Vars(r)["token"]
	if err := a.DB.DeleteShortURL(token); err != nil {
		log.WithField("token", token).WithError(err).Error("Unable to delete ShortURL in DeleteURL")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
	return
}

func (a *App) incrementRedirects(token string) {
	shortURL, err := a.DB.GetShortURL(token)
	if err != nil {
		log.WithField("token", token).WithError(err).Error("Unable to retrieve ShortURL from database")
		return
	}
	shortURL.Redirects++
	if err := a.DB.UpdateShortURL(shortURL); err != nil {
		log.WithField("short_url", shortURL).WithError(err).Error("Unable to update ShortURL in database")
	}
}

func (a *App) cleanExpiredRecords() {
	tokens, err := a.DB.GetAllURLTokens()
	if err != nil {
		log.WithError(err).Error("Unable to get all tokens in cleanExpiredRecords")
	}
	now := time.Now()
	count := 0
	for _, token := range tokens {
		shortURL, err := a.DB.GetShortURL(token)
		if err != nil {
			log.WithField("token", token).WithError(err).Error("Unable to retrieve ShortURL from database in cleanExpiredRecords")
			continue
		}
		expireTime, err := time.Parse(time.RFC3339, shortURL.Expiration)
		if err != nil {
			log.WithField("time", expireTime).WithError(err).Error("Unable to parse expire time from database")
			continue
		}
		if now.After(expireTime) {
			a.DB.DeleteShortURL(token)
			count++
		}
	}
	if count > 0 {
		log.WithField("deleted_urls", count).Info("Expired URLs removed from database")
	}
}
