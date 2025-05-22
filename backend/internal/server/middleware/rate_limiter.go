package middleware

import (
	"fmt"
	"net/http"
	"os"

	"github.com/hengadev/leviosa/pkg/ctxutil"
	"github.com/hengadev/leviosa/pkg/domainutil"
	"github.com/hengadev/leviosa/pkg/jsonio"

	"golang.org/x/time/rate"
)

// Map that holds hashed IP address and their respective rate limiter using the token bucket algorithm.
var ipLimiters = map[string]*rate.Limiter{}

// PerIPRateLimit rates limit an endpoint at the rate of r attempt per ID per second. Client can allow up to b attempts at most.
func PerIPRateLimit(lim, burst int) Middleware {
	return func(next Handlerfunc) Handlerfunc {
		return func(w http.ResponseWriter, r *http.Request) {
			// TODO: need a system to remove keys after a certain time
			ctx := r.Context()
			logger, err := ctxutil.GetLoggerFromContext(ctx)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			// I just make a fake IP for now, I know my function to work:
			IP := "127.0.0.1"
			// IP := getClientIP(r)
			if IP == "" || IP == "unknown" {
				logger.WarnContext(ctx, "Invalid IP address in per IP rate limiter")
				http.Error(w, "Invalid IP address", http.StatusBadRequest)
				return
			}
			// Check if the IP address is already present in the hash map
			limiter, exists := ipLimiters[IP]
			if !exists {
				// If the IP address is not present in the hash map, create a new rate limiter
				// limiter = rate.NewLimiter(2, 4)
				limiter = rate.NewLimiter(rate.Limit(lim), burst)
				ipLimiters[IP] = limiter
			}

			// Use the rate limiter to check if the request should be allowed
			if !limiter.Allow() {
				message := struct {
					Status string
					Body   string
				}{
					Status: "Request failed",
					Body:   "The API is at capacity, try again later",
				}
				if err := jsonio.Encode(w, http.StatusTooManyRequests, message); err != nil {
					rateLimitSalt := os.Getenv("RATE_LIMIT_SALT")
					logger.WarnContext(ctx, fmt.Sprintf("Rate limited %q for %s to %s", domainutil.HashWithSalt(IP, rateLimitSalt), r.Method, r.URL.String()))
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				return
			}

			next(w, r)
		}
	}
}

// // getClientIP retrieves the client IP set in the custom header 'X-Client-IP'. It returns an empty string if the IP address sent is invalid, unknown or could ne be parsed.
// func getClientIP(r *http.Request) string {
// 	clientIP := r.Header.Get("X-Client-IP")
// 	if clientIP == "" || clientIP == "unknown" {
// 		return ""
// 	}
// 	if ip := net.ParseIP(clientIP); ip == nil {
// 		return ""
// 	}
// 	return clientIP
// }
