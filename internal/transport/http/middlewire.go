package http

import (
	"net/http"

	log "github.com/sirupsen/logrus"
)

func JSONMiddlewire(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		next.ServeHTTP(w, r)
	})
}

func LogMiddlewire(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.WithFields(
			log.Fields{
				"method": r.Method,
				"path":   r.URL.Path,
				"remote": r.RemoteAddr,
				"agent":  r.UserAgent(),
			}).Info("request completed")
		next.ServeHTTP(w, r)
	})
}
