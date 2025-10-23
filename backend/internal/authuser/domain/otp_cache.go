package domain

import "sync"

type OTPCache struct {
	mu             sync.RWMutex
	otpDuration    int
	otpLength      int
	otpMaxAttempts int
}

func NewOTPCache(duration, length, maxAttempts int) *OTPCache {
	return &OTPCache{
		otpDuration:    duration,
		otpLength:      length,
		otpMaxAttempts: maxAttempts,
	}
}

// SetOTPDuration updates the cached OTP duration setting
func (c *OTPCache) SetOTPDuration(duration int) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.otpDuration = duration
}

// GetOTPDuration retrieves the cached OTP duration setting
func (c *OTPCache) GetOTPDuration() int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.otpDuration
}

// SetOTPLength updates the cached OTP length setting
func (c *OTPCache) SetOTPLength(length int) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.otpLength = length
}

// GetOTPLength retrieves the cached OTP length setting
func (c *OTPCache) GetOTPLength() int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.otpLength
}

// SetOTPMaxAttempts updates the cached OTP max attempts setting
func (c *OTPCache) SetOTPMaxAttempts(maxAttempts int) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.otpMaxAttempts = maxAttempts
}

// GetOTPMaxAttempts retrieves the cached OTP max attempts setting
func (c *OTPCache) GetOTPMaxAttempts() int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.otpMaxAttempts
}
