package otpService

import (
	"sync"
)

type cache struct {
	mu             sync.RWMutex
	otpDuration    int
	otpLength      int
	otpMaxAttempts int
}

// Constructor
func newCache(duration, length, maxAttempts int) *cache {
	return &cache{
		otpDuration:    duration,
		otpLength:      length,
		otpMaxAttempts: maxAttempts,
	}
}

func (s *service) SetOTPDuration(duration int) {
	s.cache.setOTPDuration(duration)
}

func (c *cache) setOTPDuration(duration int) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.otpDuration = duration
}

func (s *service) SetOTPLength(length int) {
	s.cache.setOTPLength(length)
}

func (c *cache) setOTPLength(length int) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.otpLength = length
}

func (s *service) SetOTPMaxAttempts(maxAttempts int) {
	s.cache.setOTPMaxAttempts(maxAttempts)
}

func (c *cache) setOTPMaxAttempts(maxAttempts int) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.otpMaxAttempts = maxAttempts
}
