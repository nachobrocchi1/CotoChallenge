package middleware

import (
	"log"
	"net/http"
	"time"
)

// TimingMiddleware logs the duration of the request in seconds
func TimingMiddleware(logger *log.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			next.ServeHTTP(w, r)
			duration := time.Since(start).Seconds()
			logger.Printf("[%s] %s %fs", r.Method, r.URL.Path, duration)
		})
	}
}
