package middleware

import (
	"net/http"
	"time"
)

// PathRateLimit defines a rate limit rule for a specific HTTP method + path.
type PathRateLimit struct {
	Method      string
	Path        string
	KeyPrefix   string
	MaxRequests int
	Window      time.Duration
}

// RateLimitByPath returns an http.Handler middleware that applies rate limiting
// only to requests matching the given rules. Non-matching requests pass through.
func (rl *RateLimiter) RateLimitByPath(rules []PathRateLimit) func(http.Handler) http.Handler {
	// Pre-build per-rule middleware for efficiency.
	type ruleMiddleware struct {
		method string
		path   string
		mw     func(http.Handler) http.Handler
	}
	middlewares := make([]ruleMiddleware, len(rules))
	for i, r := range rules {
		middlewares[i] = ruleMiddleware{
			method: r.Method,
			path:   r.Path,
			mw:     rl.Middleware(r.KeyPrefix, RateLimitConfig{MaxRequests: r.MaxRequests, Window: r.Window}),
		}
	}

	return func(next http.Handler) http.Handler {
		// Build the chained handlers once per next, not per request.
		type ruleHandler struct {
			method  string
			path    string
			handler http.Handler
		}
		handlers := make([]ruleHandler, len(middlewares))
		for i, rm := range middlewares {
			handlers[i] = ruleHandler{
				method:  rm.method,
				path:    rm.path,
				handler: rm.mw(next),
			}
		}
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			for _, rh := range handlers {
				if r.Method == rh.method && r.URL.Path == rh.path {
					rh.handler.ServeHTTP(w, r)
					return
				}
			}
			next.ServeHTTP(w, r)
		})
	}
}
