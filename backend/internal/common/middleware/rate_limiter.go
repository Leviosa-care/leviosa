package middleware

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

// RateLimitConfig defines the per-endpoint rate-limit parameters.
type RateLimitConfig struct {
	// MaxRequests is the maximum number of requests allowed within the window.
	MaxRequests int
	// Window is the sliding window duration.
	Window time.Duration
}

// rateLimitPipeline is the minimal interface needed from a Redis pipeline.
type rateLimitPipeline interface {
	ZRemRangeByScore(ctx context.Context, key, min, max string) *redis.IntCmd
	ZRangeWithScores(ctx context.Context, key string, start, stop int64) *redis.ZSliceCmd
	ZCard(ctx context.Context, key string) *redis.IntCmd
	ZAdd(ctx context.Context, key string, members ...redis.Z) *redis.IntCmd
	Expire(ctx context.Context, key string, ttl time.Duration) *redis.BoolCmd
	Exec(ctx context.Context) ([]redis.Cmder, error)
}

// rateLimitClient is the minimal Redis interface the rate limiter needs.
type rateLimitClient interface {
	Pipeline() rateLimitPipeline
}

// RateLimiter tracks request counts in Redis using a sliding-window counter.
// Each unique (client IP, key) pair gets its own Redis key with a TTL equal to
// the window duration, so keys expire automatically.
type RateLimiter struct {
	client rateLimitClient
	logger *slog.Logger
}

// NewRateLimiter creates a RateLimiter backed by the given Redis client.
// It accepts redis.Cmdable (which *redis.Client satisfies) and adapts it.
func NewRateLimiter(client redis.Cmdable, logger *slog.Logger) *RateLimiter {
	return &RateLimiter{client: &redisClientAdapter{Cmdable: client}, logger: logger}
}

// redisClientAdapter adapts redis.Cmdable to the rateLimitClient interface.
type redisClientAdapter struct {
	redis.Cmdable
}

func (a *redisClientAdapter) Pipeline() rateLimitPipeline {
	return a.Cmdable.Pipeline()
}

// Middleware returns an http.Handler middleware that rate-limits by client IP
// using the provided config. The keyPrefix is used to namespace the Redis key
// (e.g. "rate_limit:signin").
func (rl *RateLimiter) Middleware(keyPrefix string, cfg RateLimitConfig) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ip := extractIP(r)
			redisKey := fmt.Sprintf("rate_limit:%s:%s", keyPrefix, ip)

			ctx := r.Context()

			allowed, retryAfter, err := rl.check(ctx, redisKey, cfg)
			if err != nil {
				// On Redis errors, allow the request through rather than blocking everyone.
				rl.logger.ErrorContext(ctx, "rate limiter redis error, allowing request",
					"error", err, "key", redisKey)
				next.ServeHTTP(w, r)
				return
			}

			if !allowed {
				w.Header().Set("Retry-After", strconv.Itoa(int(retryAfter.Seconds())))
				http.Error(w, "too many requests", http.StatusTooManyRequests)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// check increments the counter for the given key. It returns (allowed, retryAfter, error).
// On the (N+1)th request exceeding MaxRequests within Window, it returns false.
func (rl *RateLimiter) check(ctx context.Context, key string, cfg RateLimitConfig) (allowed bool, retryAfter time.Duration, err error) {
	now := time.Now()
	windowStart := now.Add(-cfg.Window)

	pipe := rl.client.Pipeline()
	// Remove entries outside the window.
	pipe.ZRemRangeByScore(ctx, key, "-inf", fmt.Sprintf("%d", windowStart.UnixMilli()))
	// Fetch the oldest remaining entry so we can compute a precise Retry-After.
	oldestCmd := pipe.ZRangeWithScores(ctx, key, 0, 0)
	// Count entries in the current window.
	countCmd := pipe.ZCard(ctx, key)
	// Add current request to the sorted set (score = timestamp, member = unique ID).
	pipe.ZAdd(ctx, key, redis.Z{
		Score:  float64(now.UnixMilli()),
		Member: now.UnixNano(),
	})
	// Set expiry on the key so it auto-cleans.
	pipe.Expire(ctx, key, cfg.Window+time.Second)

	if _, err = pipe.Exec(ctx); err != nil {
		return false, 0, fmt.Errorf("redis pipeline: %w", err)
	}

	count := countCmd.Val()

	if int(count) >= cfg.MaxRequests {
		// Default to the full window; refine to the precise time-to-expiry of the
		// oldest entry if available.
		retryAfter = cfg.Window
		if entries := oldestCmd.Val(); len(entries) > 0 {
			oldestTime := time.UnixMilli(int64(entries[0].Score))
			if precise := oldestTime.Add(cfg.Window).Sub(now); precise > 0 {
				retryAfter = precise
			}
		}
		return false, retryAfter, nil
	}

	return true, 0, nil
}

// extractIP returns the client IP from the X-Client-IP header (set by Caddy)
// or falls back to RemoteAddr.
func extractIP(r *http.Request) string {
	if ip := r.Header.Get("X-Client-IP"); ip != "" {
		if parsed := net.ParseIP(ip); parsed != nil {
			return ip
		}
	}
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return host
}
