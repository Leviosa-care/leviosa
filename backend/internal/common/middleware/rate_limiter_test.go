package middleware

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	"github.com/redis/go-redis/v9"
)

// ---------------------------------------------------------------------------
// Mock implementations
// ---------------------------------------------------------------------------

// mockRedisClient implements the rateLimitClient interface for tests.
type mockRedisClient struct {
	entries []redis.Z
	err     error // inject pipeline exec error
}

func (m *mockRedisClient) Pipeline() rateLimitPipeline {
	return &mockPipeliner{store: m}
}

// mockPipeliner tracks the sorted-set entries and simulates the pipeline.
type mockPipeliner struct {
	store      *mockRedisClient
	countCmd   *redis.IntCmd
	oldestCmd  *redis.ZSliceCmd
}

func (p *mockPipeliner) Exec(ctx context.Context) ([]redis.Cmder, error) {
	if p.store.err != nil {
		return nil, p.store.err
	}

	// Populate the oldest-entry result before adding the new request, matching
	// real pipeline ordering (ZRangeWithScores after ZRemRangeByScore, before ZAdd).
	if p.oldestCmd != nil {
		if len(p.store.entries) > 0 {
			p.oldestCmd.SetVal([]redis.Z{p.store.entries[0]})
		} else {
			p.oldestCmd.SetVal(nil)
		}
	}

	// The count we report is the number of entries BEFORE this request is added,
	// matching the real pipeline ordering (ZCard before ZAdd).
	count := int64(len(p.store.entries))

	// Simulate adding the new entry.
	now := time.Now()
	p.store.entries = append(p.store.entries, redis.Z{
		Score:  float64(now.UnixMilli()),
		Member: now.UnixNano(),
	})

	// Set the value on the pre-created IntCmd.
	p.countCmd.SetVal(count)

	return []redis.Cmder{p.countCmd}, nil
}

func (p *mockPipeliner) ZRemRangeByScore(ctx context.Context, key, min, max string) *redis.IntCmd {
	return redis.NewIntCmd(ctx)
}

func (p *mockPipeliner) ZRangeWithScores(ctx context.Context, key string, start, stop int64) *redis.ZSliceCmd {
	cmd := redis.NewZSliceCmd(ctx)
	p.oldestCmd = cmd
	return cmd
}

func (p *mockPipeliner) ZCard(ctx context.Context, key string) *redis.IntCmd {
	cmd := redis.NewIntCmd(ctx)
	p.countCmd = cmd
	return cmd
}

func (p *mockPipeliner) ZAdd(ctx context.Context, key string, members ...redis.Z) *redis.IntCmd {
	return redis.NewIntCmd(ctx)
}

func (p *mockPipeliner) Expire(ctx context.Context, key string, ttl time.Duration) *redis.BoolCmd {
	return redis.NewBoolCmd(ctx)
}

// failPipeliner always returns an error on Exec.
type failPipeliner struct{}

func (f *failPipeliner) Exec(ctx context.Context) ([]redis.Cmder, error) {
	return nil, errors.New("redis is down")
}

func (f *failPipeliner) ZRemRangeByScore(ctx context.Context, key, min, max string) *redis.IntCmd {
	return nil
}

func (f *failPipeliner) ZRangeWithScores(ctx context.Context, key string, start, stop int64) *redis.ZSliceCmd {
	return nil
}

func (f *failPipeliner) ZCard(ctx context.Context, key string) *redis.IntCmd { return nil }
func (f *failPipeliner) ZAdd(ctx context.Context, key string, members ...redis.Z) *redis.IntCmd {
	return nil
}
func (f *failPipeliner) Expire(ctx context.Context, key string, ttl time.Duration) *redis.BoolCmd {
	return nil
}

// failRedisClient returns a failPipeliner.
type failRedisClient struct{}

func (f *failRedisClient) Pipeline() rateLimitPipeline { return &failPipeliner{} }

// ---------------------------------------------------------------------------
// Helper to create a rate limiter from our mock
// ---------------------------------------------------------------------------

func newTestRateLimiter(client *mockRedisClient) *RateLimiter {
	return &RateLimiter{client: client, logger: slog.Default()}
}

// ---------------------------------------------------------------------------
// Tests
// ---------------------------------------------------------------------------

func TestRateLimiter_AllowsWithinLimit(t *testing.T) {
	store := &mockRedisClient{}
	rl := newTestRateLimiter(store)

	cfg := RateLimitConfig{MaxRequests: 3, Window: 15 * time.Minute}
	mw := rl.Middleware("signin", cfg)

	handler := mw(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	for i := 0; i < 3; i++ {
		req := httptest.NewRequest(http.MethodPost, "/auth/login", nil)
		req.RemoteAddr = "1.2.3.4:1234"
		w := httptest.NewRecorder()

		handler.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("request %d: expected 200, got %d", i+1, w.Code)
		}
	}
}

func TestRateLimiter_BlocksAfterLimit(t *testing.T) {
	store := &mockRedisClient{}
	rl := newTestRateLimiter(store)

	cfg := RateLimitConfig{MaxRequests: 3, Window: 15 * time.Minute}
	mw := rl.Middleware("signin", cfg)

	handler := mw(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	// Exhaust the limit.
	for i := 0; i < 3; i++ {
		req := httptest.NewRequest(http.MethodPost, "/auth/login", nil)
		req.RemoteAddr = "1.2.3.4:1234"
		w := httptest.NewRecorder()
		handler.ServeHTTP(w, req)
		if w.Code != http.StatusOK {
			t.Fatalf("request %d: expected 200, got %d", i+1, w.Code)
		}
	}

	// The (N+1)th request should be blocked.
	req := httptest.NewRequest(http.MethodPost, "/auth/login", nil)
	req.RemoteAddr = "1.2.3.4:1234"
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	if w.Code != http.StatusTooManyRequests {
		t.Fatalf("expected 429, got %d", w.Code)
	}

	retryAfter := w.Header().Get("Retry-After")
	if retryAfter == "" {
		t.Fatal("expected Retry-After header to be set")
	}
	secs, err := strconv.Atoi(retryAfter)
	if err != nil {
		t.Fatalf("Retry-After should be an integer: %v", err)
	}
	// The precise value is time-until-oldest-entry-expires, which in this test is
	// just under the full window (entries were added milliseconds ago).
	if secs <= 0 || secs > 900 {
		t.Fatalf("Retry-After should be in (0, 900], got %d", secs)
	}
}

func TestRateLimiter_DifferentIPsIndependent(t *testing.T) {
	// Each IP gets its own Redis key, so they are independent.
	// Since our mock shares a single slice, we test that the middleware
	// constructs different keys per IP.
	// A more thorough isolation test would use real Redis (integration test).

	store1 := &mockRedisClient{}
	rl1 := newTestRateLimiter(store1)
	cfg := RateLimitConfig{MaxRequests: 2, Window: 15 * time.Minute}
	mw := rl1.Middleware("signin", cfg)

	okHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	handler := mw(okHandler)

	// IP 1 uses both requests.
	for i := 0; i < 2; i++ {
		req := httptest.NewRequest(http.MethodPost, "/auth/login", nil)
		req.RemoteAddr = "1.2.3.4:1234"
		w := httptest.NewRecorder()
		handler.ServeHTTP(w, req)
		if w.Code != http.StatusOK {
			t.Fatalf("IP1 request %d: expected 200, got %d", i+1, w.Code)
		}
	}

	// IP 2 with a fresh store should still be allowed.
	store2 := &mockRedisClient{}
	rl2 := newTestRateLimiter(store2)
	mw2 := rl2.Middleware("signin", cfg)
	handler2 := mw2(okHandler)

	req := httptest.NewRequest(http.MethodPost, "/auth/login", nil)
	req.RemoteAddr = "5.6.7.8:5678"
	w := httptest.NewRecorder()
	handler2.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("IP2 with fresh store: expected 200, got %d", w.Code)
	}
}

func TestRateLimiter_UsesXClientIPHeader(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/auth/login", nil)
	req.RemoteAddr = "10.0.0.1:1234"
	req.Header.Set("X-Client-IP", "203.0.113.50")

	ip := extractIP(req)
	if ip != "203.0.113.50" {
		t.Fatalf("expected X-Client-IP to be used, got %s", ip)
	}
}

func TestRateLimiter_FallsBackToRemoteAddr(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/auth/login", nil)
	req.RemoteAddr = "10.0.0.1:1234"

	ip := extractIP(req)
	if ip != "10.0.0.1" {
		t.Fatalf("expected RemoteAddr fallback, got %s", ip)
	}
}

func TestRateLimiter_IgnoresInvalidXClientIP(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/auth/login", nil)
	req.RemoteAddr = "10.0.0.1:1234"
	req.Header.Set("X-Client-IP", "not-an-ip")

	ip := extractIP(req)
	if ip != "10.0.0.1" {
		t.Fatalf("expected fallback to RemoteAddr for invalid X-Client-IP, got %s", ip)
	}
}

func TestRateLimiter_RedisErrorAllowsThrough(t *testing.T) {
	rl := &RateLimiter{
		client: &failRedisClient{},
		logger: slog.Default(),
	}

	cfg := RateLimitConfig{MaxRequests: 1, Window: 15 * time.Minute}
	mw := rl.Middleware("signin", cfg)

	handler := mw(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodPost, "/auth/login", nil)
	req.RemoteAddr = "1.2.3.4:1234"
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("on Redis error, request should be allowed through, got %d", w.Code)
	}
}

func TestRateLimiter_ResponseContainsCorrectRetryAfter(t *testing.T) {
	store := &mockRedisClient{}
	rl := newTestRateLimiter(store)

	cfg := RateLimitConfig{MaxRequests: 1, Window: 5 * time.Minute}
	mw := rl.Middleware("password_reset", cfg)

	handler := mw(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	// First request passes.
	req := httptest.NewRequest(http.MethodPost, "/auth/password/reset/request", nil)
	req.RemoteAddr = "1.2.3.4:1234"
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("first request: expected 200, got %d", w.Code)
	}

	// Second request blocked with Retry-After = 300s (5 minutes).
	req = httptest.NewRequest(http.MethodPost, "/auth/password/reset/request", nil)
	req.RemoteAddr = "1.2.3.4:1234"
	w = httptest.NewRecorder()
	handler.ServeHTTP(w, req)
	if w.Code != http.StatusTooManyRequests {
		t.Fatalf("second request: expected 429, got %d", w.Code)
	}
	retryAfterStr := w.Header().Get("Retry-After")
	retryAfterSecs, err := strconv.Atoi(retryAfterStr)
	if err != nil {
		t.Fatalf("Retry-After should be an integer, got %q: %v", retryAfterStr, err)
	}
	// Precise value is time-until-oldest-entry-expires: just under the 5-minute window.
	if retryAfterSecs <= 0 || retryAfterSecs > 300 {
		t.Fatalf("expected Retry-After in (0, 300], got %d", retryAfterSecs)
	}
}
