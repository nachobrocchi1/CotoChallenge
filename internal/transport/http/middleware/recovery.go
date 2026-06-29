package middleware

import (
	"log"
	"net/http"
	"runtime/debug"
)

// Recovery recovers any call after a panic
func Recovery(logger *log.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					logger.Printf("[PANIC RECOVERED] %v\n%s", err, debug.Stack())
					http.Error(w, "internal server error", http.StatusInternalServerError)
				}
			}()
			next.ServeHTTP(w, r)
		})
	}
}
