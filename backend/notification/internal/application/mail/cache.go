package mail

// email
func (s *MailService) GetCompanyEmail() string {
	return s.cache.GetCompanyEmail()
}

func (s *MailService) SetCompanyEmail(email string) {
	s.cache.SetCompanyEmail(email)
}

// instagram
func (s *MailService) GetCompanyInstagram() string {
	return s.cache.GetCompanyInstagram()
}

func (s *MailService) SetCompanyInstagram(insta string) {
	s.cache.SetCompanyEmail(insta)
}

// legal address
func (s *MailService) GetCompanyLegalAddress() string {
	return s.cache.GetCompanyLegalAddress()
}

func (s *MailService) SetCompanyLegalAddress(insta string) {
	s.cache.SetCompanyEmail(insta)
}

// logo
func (s *MailService) GetLogo() []byte {
	return s.cache.GetLogo()
}

func (s *MailService) SetLogo(logo []byte) {
	s.cache.SetLogo(logo)
}
