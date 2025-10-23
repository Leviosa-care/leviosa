package domain

import "sync"

type MailCache struct {
	mu                  sync.RWMutex
	CompanyEmail        string
	CompanyInstagram    string
	CompanyLegalAddress string
	CompanyLogo         []byte
}

// Constructor
func NewMailCache(email, insta, address string, logo []byte) *MailCache {
	return &MailCache{
		CompanyEmail:        email,
		CompanyInstagram:    insta,
		CompanyLegalAddress: address,
		CompanyLogo:         logo,
	}
}

// SetCompanyEmail updates the cached legal address
func (c *MailCache) SetCompanyInstagram(addr string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.CompanyInstagram = addr
}

// GetCompanyEmail retrieves the cached email
func (c *MailCache) GetCompanyInstagram() string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.CompanyInstagram
}

// SetCompanyLegalAddress updates the cached legal address
func (c *MailCache) SetCompanyLegalAddress(addr string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.CompanyLegalAddress = addr
}

// GetCompanyLegalAddress retrieves the cached email
func (c *MailCache) GetCompanyLegalAddress() string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.CompanyLegalAddress
}

// SetCompanyEmail updates the cached email
func (c *MailCache) SetCompanyEmail(email string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.CompanyEmail = email
}

// GetCompanyEmail retrieves the cached email
func (c *MailCache) GetCompanyEmail() string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.CompanyEmail
}

// SetLogo updates the cached logo
func (c *MailCache) SetLogo(logo []byte) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.CompanyLogo = logo
}

// GetLogo retrieves the cached logo
func (c *MailCache) GetLogo() []byte {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.CompanyLogo
}
