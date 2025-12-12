package domain

import (
	"sync"
	"time"
)

// CompanyCache is a simple in-memory cache with TTL support
type CompanyCache struct {
	mu   sync.RWMutex
	data map[string]cacheEntry
}

type cacheEntry struct {
	value     any
	expiresAt time.Time
}

func NewCompanyCache() *CompanyCache {
	return &CompanyCache{
		data: make(map[string]cacheEntry),
	}
}

// Set stores a value with TTL (0 = no expiration)
func (c *CompanyCache) Set(key string, value any, ttl time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()

	var expiresAt time.Time
	if ttl > 0 {
		expiresAt = time.Now().Add(ttl)
	}

	c.data[key] = cacheEntry{
		value:     value,
		expiresAt: expiresAt,
	}
}

// Get retrieves a value if it exists and hasn't expired
func (c *CompanyCache) Get(key string) (any, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	entry, exists := c.data[key]
	if !exists {
		return nil, false
	}

	// Check expiration (zero time means no expiration)
	if !entry.expiresAt.IsZero() && time.Now().After(entry.expiresAt) {
		return nil, false
	}

	return entry.value, true
}

// Invalidate removes a specific key
func (c *CompanyCache) Invalidate(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.data, key)
}

// InvalidateAll clears the entire cache
func (c *CompanyCache) InvalidateAll() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.data = make(map[string]cacheEntry)
}
