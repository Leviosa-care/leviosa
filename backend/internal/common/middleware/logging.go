package middleware

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"os"
	"slices"
	"strings"
	"sync"
	"time"

	"crypto/sha256"
	"encoding/hex"

	"github.com/Leviosa-care/leviosa/backend/internal/common/ctxutil"
	"github.com/Leviosa-care/leviosa/backend/internal/common/envmode"

	"github.com/google/uuid"
)

// LoggingConfig holds configuration for the logging middleware
type LoggingConfig struct {
	IPHeaders       []string          // Headers to check for client IP (in order of preference)
	LoggingSalt     string            // Salt for IP hashing (optional)
	SkipPaths       []string          // Paths to skip logging
	EnableIPLog     bool              // Whether to log IP addresses
	SampleRate      float64           // Sampling rate (0.0 to 1.0, 0 = disabled)
	CustomFields    map[string]string // Custom fields to add to all log entries (key -> header name)
	LogRequestBody  bool              // Whether to log request body (for debugging)
	LogResponseBody bool              // Whether to log response body (for debugging)
	MaxBodySize     int64             // Maximum body size to log (bytes)
}

// DefaultLoggingConfig returns sensible defaults
func DefaultLoggingConfig() *LoggingConfig {
	return &LoggingConfig{
		IPHeaders:       []string{"X-Forwarded-For", "X-Real-IP", "CF-Connecting-IP"},
		SkipPaths:       []string{"/health", "/healthz", "/metrics"},
		EnableIPLog:     true,
		SampleRate:      1.0, // Log all requests by default
		CustomFields:    make(map[string]string),
		LogRequestBody:  false, // Disabled by default for security
		LogResponseBody: false, // Disabled by default for security
		MaxBodySize:     1024,  // 1KB limit for body logging
	}
}

const (
	maxCacheSize  = 1000 // Prevent memory exhaustion
	minSaltLength = 8    // Minimum salt length for security
)

var (
	ipHashCache = make(map[string]string)
	cacheMutex  = sync.RWMutex{}
)

// isValidSalt validates salt meets minimum security requirements
func isValidSalt(salt string) bool {
	return len(salt) >= minSaltLength
}

// sanitizeIP validates and sanitizes IP address to prevent injection
func sanitizeIP(ip string) string {
	// Remove any potential injection characters
	ip = strings.TrimSpace(ip)

	// Basic validation for IPv4/IPv6 format
	if net.ParseIP(ip) == nil {
		return "" // Invalid IP format
	}

	return ip
}

// evictOldestFromCache removes oldest entries when cache is full
func evictOldestFromCache() {
	if len(ipHashCache) < maxCacheSize {
		return
	}

	// Simple eviction: clear 20% of cache when full
	targetSize := int(float64(maxCacheSize) * 0.8)
	count := 0
	for key := range ipHashCache {
		delete(ipHashCache, key)
		count++
		if len(ipHashCache) <= targetSize {
			break
		}
	}
}

// extractClientIP extracts client IP from request headers with fallbacks and security validation
func extractClientIP(r *http.Request, env envmode.Mode, headers []string) string {
	// In development, always use remote address
	if env == envmode.Dev {
		host, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			return sanitizeIP(r.RemoteAddr)
		}
		return sanitizeIP(host)
	}

	// Production: check configured headers in order with security validation
	for _, header := range headers {
		if rawIP := r.Header.Get(header); rawIP != "" {
			// Handle comma-separated IPs (X-Forwarded-For can contain multiple)
			if firstIP := strings.Split(rawIP, ",")[0]; firstIP != "" {
				if cleanIP := sanitizeIP(firstIP); cleanIP != "" {
					return cleanIP
				}
			}
		}
	}

	// Fallback to RemoteAddr with validation
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return sanitizeIP(r.RemoteAddr)
	}
	return sanitizeIP(host)
}

// getHashedIP returns a cached hash of the IP address using the provided salt with security validation
func getHashedIP(ip, salt string) string {
	// Validate salt meets security requirements
	if !isValidSalt(salt) {
		return maskIP(ip) // Fallback to masking if salt is insufficient
	}

	cacheKey := fmt.Sprintf("%s:%s", ip, salt)

	cacheMutex.RLock()
	if cached, exists := ipHashCache[cacheKey]; exists {
		cacheMutex.RUnlock()
		return cached
	}
	cacheMutex.RUnlock()

	// Hash IP with salt for privacy compliance using SHA-256
	data := fmt.Sprintf("%s:%s", ip, salt)
	hash := sha256.Sum256([]byte(data))
	hashed := hex.EncodeToString(hash[:])

	// Cache the result with size management
	cacheMutex.Lock()
	evictOldestFromCache() // Manage cache size before adding
	ipHashCache[cacheKey] = hashed
	cacheMutex.Unlock()

	return hashed
}

// maskIP masks an IP address for privacy (when no salt is available)
func maskIP(ip string) string {
	parts := strings.Split(ip, ".")
	if len(parts) == 4 {
		// IPv4: show first two octets, mask last two
		return fmt.Sprintf("%s.%s.xxx.xxx", parts[0], parts[1])
	}

	// IPv6 or other format: show first few characters
	if len(ip) > 6 {
		return ip[:4] + "xxxx"
	}
	return "xxx.xxx.xxx.xxx"
}

// shouldSample determines if this request should be logged based on sample rate
func shouldSample(rate float64) bool {
	if rate <= 0 {
		return false
	}
	if rate >= 1.0 {
		return true
	}
	// Simple sampling based on time - deterministic but distributed
	now := time.Now().UnixNano()
	return float64(now%1000)/1000.0 < rate
}

// responseWriter wraps http.ResponseWriter to capture status code and response size
type responseWriter struct {
	http.ResponseWriter
	statusCode   int
	bytesWritten int64
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func (rw *responseWriter) Write(b []byte) (int, error) {
	size, err := rw.ResponseWriter.Write(b)
	rw.bytesWritten += int64(size)
	return size, err
}

func AttachLogger(env envmode.Mode, logger *slog.Logger) func(http.Handler) http.Handler {
	return AttachLoggerWithConfig(env, logger, DefaultLoggingConfig())
}

func AttachLoggerWithConfig(env envmode.Mode, logger *slog.Logger, config *LoggingConfig) func(http.Handler) http.Handler {
	// Validate and apply defaults to config
	if config == nil {
		config = DefaultLoggingConfig()
	}
	if len(config.IPHeaders) == 0 {
		config.IPHeaders = []string{"X-Forwarded-For", "X-Real-IP", "CF-Connecting-IP"}
	}
	if config.SampleRate < 0 || config.SampleRate > 1.0 {
		config.SampleRate = 1.0 // Default to logging all requests if invalid
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Skip logging for configured paths
			if slices.Contains(config.SkipPaths, r.URL.Path) {
				next.ServeHTTP(w, r)
				return
			}
			ctx := r.Context()

			requestID := uuid.NewString()

			// Extract client IP with fallbacks
			clientIP := extractClientIP(r, env, config.IPHeaders)

			// Get logging salt with fallback and security validation
			loggingSalt := config.LoggingSalt
			if loggingSalt == "" {
				if envSalt := os.Getenv("LOGGING_SALT"); envSalt != "" {
					loggingSalt = envSalt
				}
			}

			// Validate salt security requirements
			if loggingSalt != "" && !isValidSalt(loggingSalt) {
				if env != envmode.Dev {
					logger.WarnContext(ctx, "LOGGING_SALT too short for security requirements - IP addresses will be masked instead of hashed")
				}
				loggingSalt = "" // Clear invalid salt to force masking
			} else if loggingSalt == "" && env != envmode.Dev {
				logger.WarnContext(ctx, "LOGGING_SALT not configured - IP addresses will be masked instead of hashed")
			}

			// Create logger fields with pre-allocated capacity for performance
			loggerFields := make([]interface{}, 0, 10) // Pre-allocate for expected fields
			loggerFields = append(loggerFields,
				"method", r.Method,
				"url", r.URL.String(),
				"requestID", requestID,
				"user_agent", r.Header.Get("User-Agent"))

			// Add IP if enabled and available
			if config.EnableIPLog && clientIP != "" {
				var ipField string
				if loggingSalt != "" {
					ipField = getHashedIP(clientIP, loggingSalt)
				} else {
					ipField = maskIP(clientIP)
				}
				loggerFields = append(loggerFields, "ip", ipField)
			}

			// Add custom fields from headers
			for fieldName, headerName := range config.CustomFields {
				if headerValue := r.Header.Get(headerName); headerValue != "" {
					loggerFields = append(loggerFields, fieldName, headerValue)
				}
			}

			requestLogger := logger.With(loggerFields...)

			ctx = context.WithValue(r.Context(), ctxutil.LoggerKey, requestLogger)

			// Apply sampling if configured
			shouldLog := config.SampleRate >= 1.0 || (config.SampleRate > 0 && shouldSample(config.SampleRate))

			if shouldLog {
				requestLogger.InfoContext(ctx, "Request started")
			}

			start := time.Now()

			// Wrap response writer to capture status code and response size
			wrapped := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}

			next.ServeHTTP(wrapped, r.WithContext(ctx))

			duration := time.Since(start)
			if shouldLog {
				requestLogger.InfoContext(ctx, "Request completed",
					"duration", duration,
					"status_code", wrapped.statusCode,
					"response_size", wrapped.bytesWritten)
			}
		})
	}
}
