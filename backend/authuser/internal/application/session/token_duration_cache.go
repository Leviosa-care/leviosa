package session

import (
	"context"
	"sync"
	"time"
)

// TokenDurationCache manages cached token durations with thread-safe access
type TokenDurationCache struct {
	mu                         sync.RWMutex
	accessTokenDurationMinutes int
	refreshTokenDurationHours  int
	defaultAccessDuration      time.Duration
	defaultRefreshDuration     time.Duration
}

// NewTokenDurationCache creates a new token duration cache with default values
func NewTokenDurationCache() *TokenDurationCache {
	return &TokenDurationCache{
		accessTokenDurationMinutes: 30,  // 30 minutes default
		refreshTokenDurationHours:  168, // 7 days default
		defaultAccessDuration:      30 * time.Minute,
		defaultRefreshDuration:     168 * time.Hour,
	}
}

// GetAccessTokenDuration returns the current access token duration
func (c *TokenDurationCache) GetAccessTokenDuration() time.Duration {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return time.Duration(c.accessTokenDurationMinutes) * time.Minute
}

// GetRefreshTokenDuration returns the current refresh token duration
func (c *TokenDurationCache) GetRefreshTokenDuration() time.Duration {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return time.Duration(c.refreshTokenDurationHours) * time.Hour
}

// UpdateAccessTokenDuration updates the cached access token duration
func (c *TokenDurationCache) UpdateAccessTokenDuration(ctx context.Context, durationMinutes int) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if durationMinutes <= 0 {
		// Use default if invalid value
		c.accessTokenDurationMinutes = 30
	} else {
		c.accessTokenDurationMinutes = durationMinutes
	}

	return nil
}

// UpdateRefreshTokenDuration updates the cached refresh token duration
func (c *TokenDurationCache) UpdateRefreshTokenDuration(ctx context.Context, durationHours int) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if durationHours <= 0 {
		// Use default if invalid value
		c.refreshTokenDurationHours = 168
	} else {
		c.refreshTokenDurationHours = durationHours
	}

	return nil
}

// GetDurations returns both durations in a single call for efficiency
func (c *TokenDurationCache) GetDurations() (accessDuration, refreshDuration time.Duration) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return time.Duration(c.accessTokenDurationMinutes) * time.Minute,
		time.Duration(c.refreshTokenDurationHours) * time.Hour
}

// Reset resets both durations to their default values
func (c *TokenDurationCache) Reset() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.accessTokenDurationMinutes = 30
	c.refreshTokenDurationHours = 168
}

