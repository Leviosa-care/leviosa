package domain

import "sync"

type CompanyCache struct {
	mu                  sync.RWMutex
	companyEmail        string
	companyInstagram    string
	companyLegalAddress string
	companyLogo         []byte
}

func NewCompanyCache(email, instagram, address string, logo []byte) *CompanyCache {
	return &CompanyCache{
		companyEmail:        email,
		companyInstagram:    instagram,
		companyLegalAddress: address,
		companyLogo:         logo,
	}
}

func (c *CompanyCache) SetCompanyEmail(email string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.companyEmail = email
}

func (c *CompanyCache) GetCompanyEmail() string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.companyEmail
}

func (c *CompanyCache) SetCompanyInstagram(instagram string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.companyInstagram = instagram
}

func (c *CompanyCache) GetCompanyInstagram() string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.companyInstagram
}

func (c *CompanyCache) SetCompanyLegalAddress(address string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.companyLegalAddress = address
}

func (c *CompanyCache) GetCompanyLegalAddress() string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.companyLegalAddress
}

func (c *CompanyCache) SetLogo(logo []byte) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.companyLogo = logo
}

func (c *CompanyCache) GetLogo() []byte {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.companyLogo
}
