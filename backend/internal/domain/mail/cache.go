package mailService

import "sync"

type cache struct {
	mu           sync.RWMutex
	companyEmail string
	logo         []byte
}

// Constructor
func newCache(email string, logo []byte) *cache {
	return &cache{
		companyEmail: email,
		logo:         logo,
	}
}

func (s *service) SetCompanyEmail(email string) {
	s.cache.setCompanyEmail(email)
}

func (s *service) SetLogo(logo []byte) {
	s.cache.setLogo(logo)
}

// SetCompanyEmail updates the cached email
func (c *cache) setCompanyEmail(email string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.companyEmail = email
}

// GetCompanyEmail retrieves the cached email
func (c *cache) getCompanyEmail() string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.companyEmail
}

// SetLogo updates the cached logo
func (c *cache) setLogo(logo []byte) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.logo = logo
}

// GetLogo retrieves the cached logo
func (c *cache) getLogo() []byte {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.logo
}
