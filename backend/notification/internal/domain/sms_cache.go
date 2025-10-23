package domain

import "sync"

type SMSCache struct {
	mu           sync.RWMutex
	CompanyPhone string
}

// Constructor
func NewSMSCache(phone string) *SMSCache {
	return &SMSCache{
		CompanyPhone: phone,
	}
}

// SetCompanyPhone updates the cached phone number
func (c *SMSCache) SetCompanyPhone(phone string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.CompanyPhone = phone
}

// GetCompanyPhone retrieves the cached phone number
func (c *SMSCache) GetCompanyPhone() string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.CompanyPhone
}
