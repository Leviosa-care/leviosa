package notification

import "sync"

type cache struct {
	mu                  sync.RWMutex
	companyEmail        string
	companyInstagram    string
	companyLegalAddress string
	companyLogo         []byte
}

// Constructor
func newCache(email, insta, address string, logo []byte) *cache {
	return &cache{
		companyEmail:        email,
		companyInstagram:    insta,
		companyLegalAddress: address,
		companyLogo:         logo,
	}
}

func (s *mailService) SetCompanyEmail(email string) {
	s.cache.setCompanyEmail(email)
}

func (s *mailService) SetLogo(logo []byte) {
	s.cache.setLogo(logo)
}

func (s *mailService) GetCompanyLegalAddress(addr string) {
	s.cache.setCompanyLegalAddress(addr)
}

func (s *mailService) GetCompanyInstagram(insta string) {
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

// GetCompanyLegalAddress retrieves the cached email
func (c *cache) getCompanyLegalAddress() string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.companyLegalAddress
}

// setCompanyEmail updates the cached email
func (c *cache) setCompanyEmail(email string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.companyEmail = email
}

// getCompanyEmail retrieves the cached email
func (c *cache) getCompanyEmail() string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.companyEmail
}

// setLogo updates the cached logo
func (c *cache) setLogo(logo []byte) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.companyLogo = logo
}

// getLogo retrieves the cached logo
func (c *cache) getLogo() []byte {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.companyLogo
}
