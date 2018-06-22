package api

import (
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"
)

// Logger wraps a request and logs out the request information
func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		next.ServeHTTP(w, r)

		log.WithFields(log.Fields{
			"method":        r.Method,
			"request_uri":   r.RequestURI,
			"response_time": time.Since(start),
		}).Info("Request received")
	})
}
