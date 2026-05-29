package middleware

import (
	"net/http"
	"sync"
)

var (
	allowedOrigin     = "http://localhost:5173"
	allowedOriginOnce sync.Once
)

// SetAllowedOrigin sets the origin used by EnableCORS.
// It must be called once before the HTTP server starts accepting connections.
// Subsequent calls are ignored.
func SetAllowedOrigin(origin string) {
	if origin == "" {
		return
	}
	allowedOriginOnce.Do(func() {
		allowedOrigin = origin
	})
}

func EnableCORS(next Handler) Handler {
	return func(w http.ResponseWriter, r *http.Request) {
		// Set headers
		w.Header().Set("Access-Control-Allow-Origin", allowedOrigin)
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, X-Service-Key, X-Service-Name")
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		// Handle preflight
		if r.Method == "OPTIONS" {
			w.Header().Set("Access-Control-Max-Age", "86400")
			w.WriteHeader(http.StatusOK)
			return
		}

		next(w, r)
	}
}
