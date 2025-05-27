package mailService

import "sync"

type cache struct {
	mu                  sync.RWMutex
	companyEmail        string
	companyInstagram    string
	companyLegalAddress string
	logo                []byte
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

func (s *service) GetCompanyLegalAddress(addr string) {
	s.cache.setCompanyLegalAddress(addr)
}

func (s *service) GetCompanyInstagram(insta string) {
	s.cache.setCompanyInstagram(insta)
}

// SetCompanyEmail updates the cached legal address
func (c *cache) setCompanyInstagram(addr string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.companyInstagram = addr
}

// GetCompanyEmail retrieves the cached email
func (c *cache) getCompanyInstagram() string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.companyInstagram
}

// SetCompanyEmail updates the cached legal address
func (c *cache) setCompanyLegalAddress(addr string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.companyLegalAddress = addr
}

// GetCompanyEmail retrieves the cached email
func (c *cache) getCompanyLegalAddress() string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.companyLegalAddress
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
