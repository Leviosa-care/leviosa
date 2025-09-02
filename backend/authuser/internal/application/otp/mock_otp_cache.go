package otp

import "github.com/Leviosa-care/authuser/internal/ports"

// MockOTPCache is a mock implementation of ports.OTPCache for testing
type MockOTPCache struct {
	length      int
	duration    int
	maxAttempts int
}

// NewMockOTPCache creates a new mock OTP cache with default values
func NewMockOTPCache(length, duration, maxAttempts int) *MockOTPCache {
	return &MockOTPCache{
		length:      length,
		duration:    duration,
		maxAttempts: maxAttempts,
	}
}

func (m *MockOTPCache) GetOTPLength() int {
	return m.length
}

func (m *MockOTPCache) SetOTPLength(length int) {
	m.length = length
}

func (m *MockOTPCache) GetOTPDuration() int {
	return m.duration
}

func (m *MockOTPCache) SetOTPDuration(duration int) {
	m.duration = duration
}

func (m *MockOTPCache) GetOTPMaxAttempts() int {
	return m.maxAttempts
}

func (m *MockOTPCache) SetOTPMaxAttempts(maxAttempts int) {
	m.maxAttempts = maxAttempts
}

// Compile-time check that MockOTPCache implements ports.OTPCache
var _ ports.OTPCache = (*MockOTPCache)(nil)