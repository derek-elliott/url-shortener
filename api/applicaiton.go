package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/go-redis/redis"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	// Blank import for Postgres support
	_ "github.com/jinzhu/gorm/dialects/postgres"
	log "github.com/sirupsen/logrus"
)

// App holds the router, db and Redis connections
type App struct {
	Router   *mux.Router
	DB       *gorm.DB
	Cache    *redis.Client
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

var routes Routes

// InitDB initializes the database
func (a *App) InitDB(user, pass, name, host string, port int) error {
	connStr := fmt.Sprintf("host=%s port=%d user=%s dbname=%s sslmode=disable password=%s", host, port, user, name, pass)
	db, err := gorm.Open("postgres", connStr)
	if err != nil {
		return err
	}
	defer db.Close()
	db.AutoMigrate(&ShortURL{})
	a.DB = db
	return nil
}

// InitCache initializes the Redis client
func (a *App) InitCache(pass, host string, port int) error {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", host, port),
		Password: pass,
		DB:       0,
	})
	if _, err := client.Ping().Result(); err != nil {
		return err
	}
	a.Cache = client
	return nil
}

// InitRouter initializes the router
func (a *App) InitRouter() {
	routes = Routes{
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
		var handler http.Handler

		a.Router.Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(handler)
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
	if err := json.Unmarshal(body, &payload); err != nil {
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
	shortURL := ShortURL{}
	shortURL.URL = payload.URL
	shortURL.Expiration = time.Now().Add(duration).String()

	shortURL.createShortURL(a.DB)

	tokenURL := ShortURL{}
	tokenURL.Token = ConvertIDToToken(int(shortURL.ID))
	tokenURL.ShortenedURL = fmt.Sprintf("%s/%s", a.Hostname, shortURL.Token)

	shortURL.updateShortURL(a.DB, tokenURL)

	shortener := Shortener{shortURL.Token, shortURL.ShortenedURL}

	shortener.setURL(a.Cache, duration)
	w.WriteHeader(http.StatusCreated)

	response := URLData{shortURL.URL, shortURL.Token, shortURL.ShortenedURL, shortURL.Expiration}
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
